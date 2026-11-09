package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
	"testing"
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
	if got.BuildTime != "2026-11-01T12:00:00Z" {
		t.Fatalf("build_time = %q", got.BuildTime)
	}
	if !got.Modified {
		t.Fatal("modified = false, want true")
	}
	if got.Environment != "production" {
		t.Fatalf("environment = %q, want production", got.Environment)
	}
	if got.Schema != currentSchema {
		t.Fatalf("schema = %q, want %q", got.Schema, currentSchema)
	}
}

func TestBuildVersionDefaultsWhenInfoMissing(t *testing.T) {
	got := BuildVersion("", nil)
	if got.Version != "(devel)" {
		t.Fatalf("version = %q, want (devel)", got.Version)
	}
	if got.Commit != "unknown" || got.BuildTime != "unknown" {
		t.Fatalf("missing metadata not defaulted: %+v", got)
	}
	if got.Environment != "local" {
		t.Fatalf("environment = %q, want local", got.Environment)
	}
}

func TestVersionHandlerServesJSON(t *testing.T) {
	h := versionHandler(sampleBuildInfo())
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/version", nil)

	h.ServeHTTP(rec, req)

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
	if got.Version != "v1.4.0" {
		t.Fatalf("version = %q, want v1.4.0", got.Version)
	}
}

// legacyClient mirrors the first release of the /version contract and knows only
// a subset of the fields. A newer server that adds fields must not break it.
type legacyClient struct {
	Version     string `json:"version"`
	Environment string `json:"environment"`
}

func TestResponseShapeIsRollbackCompatible(t *testing.T) {
	// The server added build_meta after the legacy client shipped. json.Unmarshal
	// ignores unknown fields, so the old client still reads what it needs.
	future := []byte(`{"schema":"version.v1","version":"v1.5.0","environment":"prod","build_meta":{"go":"1.27"}}`)
	var got legacyClient
	if err := json.Unmarshal(future, &got); err != nil {
		t.Fatalf("legacy client failed to decode a newer response: %v", err)
	}
	if got.Version != "v1.5.0" || got.Environment != "prod" {
		t.Fatalf("legacy client misread the response: %+v", got)
	}

	// And the current handler output must still satisfy the legacy client.
	rec := httptest.NewRecorder()
	versionHandler(sampleBuildInfo()).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/version", nil))
	var fromHandler legacyClient
	if err := json.Unmarshal(rec.Body.Bytes(), &fromHandler); err != nil {
		t.Fatalf("legacy client cannot read current handler output: %v", err)
	}
	if fromHandler.Version != "v1.4.0" {
		t.Fatalf("legacy client read version = %q, want v1.4.0", fromHandler.Version)
	}
}

// TestRealBuildInfoResolves proves the metadata path works against the metadata
// this test binary actually carries when built with -buildvcs.
func TestRealBuildInfoResolves(t *testing.T) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		t.Skip("build info not available")
	}
	got := BuildVersion("test", bi)
	if got.Schema != currentSchema {
		t.Fatalf("schema = %q", got.Schema)
	}
	if got.Version == "" {
		t.Fatal("version is empty")
	}
}
