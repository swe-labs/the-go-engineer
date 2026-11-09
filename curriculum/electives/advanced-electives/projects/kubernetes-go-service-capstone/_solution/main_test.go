package main

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
	"testing"
	"time"
)

func sampleBuildInfo() *debug.BuildInfo {
	return &debug.BuildInfo{
		Main: debug.Module{Version: "v1.4.0"},
		Settings: []debug.BuildSetting{
			{Key: "vcs.revision", Value: "9f2c1ab"},
			{Key: "vcs.time", Value: "2026-11-01T12:00:00Z"},
			{Key: "vcs.modified", Value: "true"},
		},
	}
}

func TestBuildVersionFromSettings(t *testing.T) {
	got := BuildVersion("production", sampleBuildInfo())
	if got.Version != "v1.4.0" {
		t.Fatalf("version = %q, want v1.4.0", got.Version)
	}
	if got.Commit != "9f2c1ab" {
		t.Fatalf("commit = %q, want 9f2c1ab", got.Commit)
	}
	if !got.Modified {
		t.Fatal("modified = false, want true")
	}
	if got.Schema != currentSchema {
		t.Fatalf("schema = %q, want %q", got.Schema, currentSchema)
	}
}

func TestBuildVersionDefaultsWhenInfoMissing(t *testing.T) {
	got := BuildVersion("", nil)
	if got.Version != "(devel)" || got.Commit != "unknown" || got.BuildTime != "unknown" {
		t.Fatalf("defaults not applied: %+v", got)
	}
	if got.Environment != "local" {
		t.Fatalf("environment = %q, want local", got.Environment)
	}
}

func TestHealthzAlwaysOK(t *testing.T) {
	s := NewServer("test", nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("healthz = %d, want 200", rec.Code)
	}
}

func TestReadyzReflectsReadiness(t *testing.T) {
	s := NewServer("test", nil)
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)

	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("before Ready: status = %d, want 503", rec.Code)
	}

	s.Ready()
	rec = httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("after Ready: status = %d, want 200", rec.Code)
	}

	s.NotReady()
	rec = httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("after NotReady: status = %d, want 503", rec.Code)
	}
}

func TestVersionHandlerServesJSON(t *testing.T) {
	s := NewServer("staging", sampleBuildInfo())
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/version", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("content-type = %q, want application/json", ct)
	}
	var got VersionInfo
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Version != "v1.4.0" || got.Environment != "staging" {
		t.Fatalf("unexpected payload: %+v", got)
	}
}

func TestLoadConfigFailsFast(t *testing.T) {
	env := map[string]string{"APP_DATABASE_URL_FILE": "/run/secrets/db"}
	read := func(string) ([]byte, error) { return []byte("postgres://db\n"), nil }
	getenv := func(k string) string { return env[k] }

	_, err := LoadConfig(getenv, read)
	var ce *ConfigError
	if !errors.As(err, &ce) || ce.Field != "APP_MAX_CONNS" {
		t.Fatalf("err = %v, want config error for APP_MAX_CONNS", err)
	}
}

func TestLoadConfigReadsSecretFile(t *testing.T) {
	env := map[string]string{
		"APP_ENV":                "production",
		"APP_MAX_CONNS":          "25",
		"APP_DATABASE_URL_FILE":  "/run/secrets/db",
		"APP_MEMORY_LIMIT_BYTES": "1073741824",
	}
	read := func(p string) ([]byte, error) {
		if p != "/run/secrets/db" {
			t.Fatalf("read %q, want the secret path", p)
		}
		return []byte("postgres://db\n"), nil
	}
	getenv := func(k string) string { return env[k] }

	cfg, err := LoadConfig(getenv, read)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DatabaseURL != "postgres://db" {
		t.Fatalf("database url = %q, trailing newline not trimmed", cfg.DatabaseURL)
	}
	if cfg.MaxConns != 25 || cfg.MemLimit != 1073741824 {
		t.Fatalf("unexpected config: %+v", cfg)
	}
	if got := ComputeBudget(cfg.MemLimit); got != 939524096 {
		t.Fatalf("budget = %d, want 939524096", got)
	}
}

func TestShutdownDrainsInFlightRequests(t *testing.T) {
	s := NewServer("test", nil)
	ln := newListener(t)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- s.Run(ctx, ln, 5*time.Second) }()

	waitReady(t, s)

	// Start a request against /work, which takes 100ms.
	respCh := make(chan int, 1)
	go func() {
		resp, err := http.Get("http://" + ln.Addr().String() + "/work")
		if err != nil {
			respCh <- -1
			return
		}
		defer resp.Body.Close()
		respCh <- resp.StatusCode
	}()

	// Give the request time to reach the handler, then trigger the drain.
	time.Sleep(30 * time.Millisecond)
	cancel()

	select {
	case code := <-respCh:
		if code != http.StatusOK {
			t.Fatalf("in-flight request status = %d, want 200", code)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("in-flight request was dropped during shutdown")
	}

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after drain")
	}
}

func TestVersionResponseIsRollbackCompatible(t *testing.T) {
	// A legacy client knows only a subset of fields; unknown additions must be
	// ignored so an older client keeps decoding a newer response.
	type legacyClient struct {
		Version     string `json:"version"`
		Environment string `json:"environment"`
	}
	future := []byte(`{"schema":"version.v1","version":"v1.5.0","environment":"prod","build_meta":{"go":"1.27"}}`)
	var got legacyClient
	if err := json.Unmarshal(future, &got); err != nil {
		t.Fatalf("legacy client failed to decode a newer response: %v", err)
	}
	if got.Version != "v1.5.0" || got.Environment != "prod" {
		t.Fatalf("legacy client misread the response: %+v", got)
	}
}

func waitReady(t *testing.T, s *Server) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if s.ready.Load() {
			return
		}
		time.Sleep(2 * time.Millisecond)
	}
	t.Fatal("server did not become ready")
}

func newListener(t *testing.T) net.Listener {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { _ = ln.Close() })
	return ln
}
