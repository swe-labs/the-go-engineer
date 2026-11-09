package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestConfigHandler(t *testing.T) {
	os.Setenv("SERVICE_NAME", "test-api")
	os.Setenv("DATABASE_URL", "postgres://user:pass@db:5432/test")
	defer os.Unsetenv("SERVICE_NAME")
	defer os.Unsetenv("DATABASE_URL")

	req := httptest.NewRequest(http.MethodGet, "/config", nil)
	rec := httptest.NewRecorder()

	configHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var cfg AppConfig
	if err := json.NewDecoder(rec.Body).Decode(&cfg); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if cfg.ServiceName != "test-api" {
		t.Errorf("got service_name %q, want %q", cfg.ServiceName, "test-api")
	}
	if cfg.DatabaseURL != "postgres://user:pass@db:5432/test" {
		t.Errorf("got database_url %q, want %q", cfg.DatabaseURL, "postgres://user:pass@db:5432/test")
	}
}

func TestConfigHandlerDefaults(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/config", nil)
	rec := httptest.NewRecorder()

	configHandler(rec, req)

	var cfg AppConfig
	json.NewDecoder(rec.Body).Decode(&cfg)

	if cfg.Version != "1.0.0" {
		t.Errorf("got version %q, want %q", cfg.Version, "1.0.0")
	}
	if cfg.ListenAddr != ":8080" {
		t.Errorf("got listen_addr %q, want %q", cfg.ListenAddr, ":8080")
	}
}

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	healthHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var body map[string]string
	json.NewDecoder(rec.Body).Decode(&body)
	if body["status"] != "ok" {
		t.Errorf("got status %q, want %q", body["status"], "ok")
	}
}
