package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Clear any overrides
	os.Clearenv()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected LoadConfig to succeed with defaults, got error: %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("expected default Port 8080, got %d", cfg.Port)
	}
	if cfg.Environment != "development" {
		t.Errorf("expected default Environment 'development', got %q", cfg.Environment)
	}
	if cfg.WorkerCount != 4 {
		t.Errorf("expected default WorkerCount 4, got %d", cfg.WorkerCount)
	}
	if cfg.RequestTimeout != 30*time.Second {
		t.Errorf("expected default RequestTimeout 30s, got %v", cfg.RequestTimeout)
	}
}

func TestLoadConfig_EnvironmentOverrides(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("ENV", "staging")
	os.Setenv("WORKER_COUNT", "8")
	os.Setenv("JWT_SECRET", "custom-super-secure-secret-key-32b")
	os.Setenv("REQUEST_TIMEOUT", "15s")
	defer os.Clearenv()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected LoadConfig to succeed, got: %v", err)
	}

	if cfg.Port != 9090 {
		t.Errorf("expected Port 9090, got %d", cfg.Port)
	}
	if cfg.Environment != "staging" {
		t.Errorf("expected Environment 'staging', got %q", cfg.Environment)
	}
	if cfg.WorkerCount != 8 {
		t.Errorf("expected WorkerCount 8, got %d", cfg.WorkerCount)
	}
	if cfg.RequestTimeout != 15*time.Second {
		t.Errorf("expected RequestTimeout 15s, got %v", cfg.RequestTimeout)
	}
}

func TestLoadConfig_InvalidPort(t *testing.T) {
	os.Setenv("PORT", "70000") // Out of TCP range
	defer os.Clearenv()

	_, err := LoadConfig()
	if err == nil {
		t.Fatalf("expected error for invalid PORT 70000, got nil")
	}
}

func TestLoadConfig_ProductionInsecureSecret(t *testing.T) {
	os.Setenv("ENV", "production")
	// Using default insecure JWT_SECRET
	defer os.Clearenv()

	_, err := LoadConfig()
	if err == nil {
		t.Fatalf("expected error when running production with default JWT secret, got nil")
	}
}

func TestApplication_HealthCheck(t *testing.T) {
	cfg, _ := LoadConfig()
	app := NewApplication(cfg)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	app.Routes().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from /healthz, got %d", rr.Code)
	}

	var payload map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode health check response: %v", err)
	}

	if payload["status"] != "ok" || payload["service"] != "opslane-backend" {
		t.Errorf("unexpected health check payload: %+v", payload)
	}
}

func TestApplication_ReadinessCheck(t *testing.T) {
	cfg, _ := LoadConfig()
	app := NewApplication(cfg)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rr := httptest.NewRecorder()

	app.Routes().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from /readyz, got %d", rr.Code)
	}
}
