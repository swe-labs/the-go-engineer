package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInfoHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	rec := httptest.NewRecorder()

	infoHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var info Info
	if err := json.NewDecoder(rec.Body).Decode(&info); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if info.Service != "docker-basics" {
		t.Errorf("got service %q, want %q", info.Service, "docker-basics")
	}
	if info.Version != "1.0.0" {
		t.Errorf("got version %q, want %q", info.Version, "1.0.0")
	}
}

func TestInfoHandlerContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	rec := httptest.NewRecorder()

	infoHandler(rec, req)

	ct := rec.Header().Get("Content-Type")
	want := "application/json"
	if ct != want {
		t.Errorf("got Content-Type %q, want %q", ct, want)
	}
}
