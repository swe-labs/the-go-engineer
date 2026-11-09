package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuildInfoHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/build", nil)
	rec := httptest.NewRecorder()

	buildInfoHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var info BuildInfo
	if err := json.NewDecoder(rec.Body).Decode(&info); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if info.Version != "1.0.0" {
		t.Errorf("got version %q, want %q", info.Version, "1.0.0")
	}
	if !info.StaticBuild {
		t.Error("expected static_build to be true")
	}
	if info.GoVersion == "" {
		t.Error("expected non-empty go_version")
	}
}

func TestBuildInfoGoArch(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/build", nil)
	rec := httptest.NewRecorder()

	buildInfoHandler(rec, req)

	var info BuildInfo
	json.NewDecoder(rec.Body).Decode(&info)

	if info.GoArch == "" {
		t.Error("expected non-empty go_arch")
	}
	if info.GoOS == "" {
		t.Error("expected non-empty go_os")
	}
}
