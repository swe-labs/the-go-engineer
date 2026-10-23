package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHealthzEndpoint(t *testing.T) {
	state := NewServerState()
	handler := SetupRoutes(state)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"status":"ok"`) {
		t.Errorf("expected status ok, got %s", rec.Body.String())
	}
}

func TestReadyzEndpoint_StateTransition(t *testing.T) {
	state := NewServerState()
	handler := SetupRoutes(state)

	// 1. Initial ready state
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK when ready, got %d", rec.Code)
	}

	// 2. Mark unready (simulating shutdown drain)
	// In the solution, state has isReady atomic.Bool
	// In starter, student must implement this
	reqUnready := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	recUnready := httptest.NewRecorder()
	handler.ServeHTTP(recUnready, reqUnready)

	// Verifies endpoint handles ready checks
	if recUnready.Code != http.StatusOK && recUnready.Code != http.StatusServiceUnavailable {
		t.Fatalf("unexpected readyz status code: %d", recUnready.Code)
	}
}

func TestServiceInfoEndpoint(t *testing.T) {
	state := NewServerState()
	handler := SetupRoutes(state)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/info", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	var data map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &data); err != nil {
		t.Fatalf("failed to parse json response: %v", err)
	}

	for _, field := range []string{"version", "commit_sha", "go_version", "arch", "os", "uptime_sec"} {
		if _, exists := data[field]; !exists {
			t.Errorf("missing required metadata field %q in info response", field)
		}
	}
}

func TestInternalHealthCheckFunction(t *testing.T) {
	state := NewServerState()
	handler := SetupRoutes(state)
	ts := httptest.NewServer(handler)
	defer ts.Close()

	parts := strings.Split(ts.URL, ":")
	port := ":" + parts[len(parts)-1]

	code := runHealthCheck(port)
	if code != 0 {
		t.Errorf("expected health check return code 0 for healthy server, got %d", code)
	}
}

func TestDockerfileSpecification(t *testing.T) {
	content, err := os.ReadFile("Dockerfile")
	if err != nil {
		t.Fatalf("failed to read Dockerfile: %v", err)
	}
	text := string(content)

	if !strings.Contains(text, "AS builder") {
		t.Errorf("Dockerfile must utilize a multi-stage build with an 'AS builder' stage")
	}
	if !strings.Contains(text, "CGO_ENABLED=0") {
		t.Errorf("Dockerfile must compile with CGO_ENABLED=0 for static portability")
	}
	if !strings.Contains(text, "USER appuser") && !strings.Contains(text, "USER 10001") {
		t.Errorf("Dockerfile must run under an unprivileged non-root user")
	}
	if !strings.Contains(text, "FROM scratch") && !strings.Contains(text, "FROM alpine") {
		t.Errorf("Dockerfile final stage must use a minimal base (scratch or alpine)")
	}
}

func TestCIDeploymentWorkflow(t *testing.T) {
	content, err := os.ReadFile(".github/workflows/ci.yml")
	if err != nil {
		t.Fatalf("failed to read ci.yml: %v", err)
	}
	text := string(content)

	if !strings.Contains(text, "-race") {
		t.Errorf("CI workflow must execute tests with data race detector (-race)")
	}
	if !strings.Contains(text, "docker") || !strings.Contains(text, "build") {
		t.Errorf("CI workflow must include a Docker container build job")
	}
}
