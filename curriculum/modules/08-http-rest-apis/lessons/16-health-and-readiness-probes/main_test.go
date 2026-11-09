package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthzHandlerHealthy(t *testing.T) {
	state.setDBReady(true)
	state.setCacheReady(true)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	healthzHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var hr healthResponse
	if err := json.NewDecoder(resp.Body).Decode(&hr); err != nil {
		t.Fatalf("json decode: %v", err)
	}
	if hr.Status != "ok" {
		t.Errorf("status = %q, want %q", hr.Status, "ok")
	}
}

func TestHealthzHandlerUnhealthy(t *testing.T) {
	state.setDBReady(false)
	state.setCacheReady(false)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	healthzHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", resp.StatusCode)
	}

	var hr healthResponse
	json.NewDecoder(resp.Body).Decode(&hr)
	if hr.Status != "unhealthy" {
		t.Errorf("status = %q, want %q", hr.Status, "unhealthy")
	}

	state.setDBReady(true)
	state.setCacheReady(true)
}

func TestHealthzMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	rec := httptest.NewRecorder()
	healthzHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

func TestReadyzHandlerReady(t *testing.T) {
	state.setDBReady(true)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	readyzHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var hr healthResponse
	json.NewDecoder(resp.Body).Decode(&hr)
	if hr.Status != "ready" {
		t.Errorf("status = %q, want %q", hr.Status, "ready")
	}
}

func TestReadyzHandlerNotReady(t *testing.T) {
	state.setDBReady(false)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	readyzHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", resp.StatusCode)
	}

	var hr healthResponse
	json.NewDecoder(resp.Body).Decode(&hr)
	if hr.Status != "not ready" {
		t.Errorf("status = %q, want %q", hr.Status, "not ready")
	}

	state.setDBReady(true)
}

func TestLivenessHandler(t *testing.T) {
	tests := []struct {
		name     string
		check    func() bool
		wantCode int
		wantBody string
	}{
		{
			name:     "healthy check returns 200",
			check:    func() bool { return true },
			wantCode: http.StatusOK,
			wantBody: "ok",
		},
		{
			name:     "unhealthy check returns 503",
			check:    func() bool { return false },
			wantCode: http.StatusServiceUnavailable,
			wantBody: "unhealthy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := livenessHandler(tt.check)
			req := httptest.NewRequest(http.MethodGet, "/livez", nil)
			rec := httptest.NewRecorder()
			handler(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantCode {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantCode)
			}

			var hr healthResponse
			json.NewDecoder(resp.Body).Decode(&hr)
			if !strings.Contains(hr.Status, tt.wantBody) {
				t.Errorf("status text = %q, want containing %q", hr.Status, tt.wantBody)
			}
		})
	}
}

func TestLessonCompiles(t *testing.T) {
}
