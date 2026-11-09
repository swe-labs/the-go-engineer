package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthzInitiallyHealthy(t *testing.T) {
	healthy.Store(true)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	healthzHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var status HealthStatus
	json.NewDecoder(rec.Body).Decode(&status)
	if status.Status != "healthy" {
		t.Errorf("got status %q, want %q", status.Status, "healthy")
	}
}

func TestHealthzUnhealthy(t *testing.T) {
	healthy.Store(false)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	healthzHandler(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}

	var status HealthStatus
	json.NewDecoder(rec.Body).Decode(&status)
	if status.Status != "unhealthy" {
		t.Errorf("got status %q, want %q", status.Status, "unhealthy")
	}

	healthy.Store(true)
}

func TestReadyzNotReady(t *testing.T) {
	ready.Store(false)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	readyzHandler(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
}

func TestReadyzWhenReady(t *testing.T) {
	ready.Store(true)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	readyzHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var body map[string]string
	json.NewDecoder(rec.Body).Decode(&body)
	if body["status"] != "ready" {
		t.Errorf("got status %q, want %q", body["status"], "ready")
	}
}

func TestHealthzResponseHasTimestamp(t *testing.T) {
	healthy.Store(true)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	healthzHandler(rec, req)

	var status HealthStatus
	json.NewDecoder(rec.Body).Decode(&status)
	if status.Timestamp == "" {
		t.Error("expected non-empty timestamp")
	}
	if status.Uptime == "" {
		t.Error("expected non-empty uptime")
	}
}

func TestSimulatedBecomeReady(t *testing.T) {
	healthy.Store(true)
	ready.Store(false)

	initReq := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	initRec := httptest.NewRecorder()
	readyzHandler(initRec, initReq)
	if initRec.Code != http.StatusServiceUnavailable {
		t.Error("expected not ready initially")
	}

	ready.Store(true)
	time.Sleep(10 * time.Millisecond)

	finalReq := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	finalRec := httptest.NewRecorder()
	readyzHandler(finalRec, finalReq)
	if finalRec.Code != http.StatusOK {
		t.Error("expected ready after setting ready=true")
	}
}
