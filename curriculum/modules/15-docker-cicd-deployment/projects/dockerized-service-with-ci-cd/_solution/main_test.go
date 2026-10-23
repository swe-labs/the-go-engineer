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
	state.isReady.Store(false)

	reqUnready := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	recUnready := httptest.NewRecorder()
	handler.ServeHTTP(recUnready, reqUnready)

	if recUnready.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 Service Unavailable when unready, got %d", recUnready.Code)
	}
	if !strings.Contains(recUnready.Body.String(), `"status":"not_ready"`) {
		t.Errorf("expected not_ready body, got %s", recUnready.Body.String())
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

	// Extract port from test server URL
	// ts.URL is "http://127.0.0.1:XXXXX"
	parts := strings.Split(ts.URL, ":")
	port := ":" + parts[len(parts)-1]

	code := runHealthCheck(port)
	if code != 0 {
		t.Errorf("expected health check return code 0 for healthy server, got %d", code)
	}

	// Non-existent port should fail
	failCode := runHealthCheck(":59999")
	if failCode != 1 {
		t.Errorf("expected return code 1 for unreachable server, got %d", failCode)
	}
}

func TestDockerfileSpecification(t *testing.T) {
	content, err := os.ReadFile("Dockerfile")
	if err != nil {
		t.Fatalf("failed to read Dockerfile: %v", err)
	}
	text := string(content)

	// Check multi-stage build
	if !strings.Contains(text, "AS builder") {
		t.Errorf("Dockerfile must utilize a multi-stage build with an 'AS builder' stage")
	}

	// Check static compilation flags
	if !strings.Contains(text, "CGO_ENABLED=0") {
		t.Errorf("Dockerfile must compile with CGO_ENABLED=0 for static portability")
	}

	// Check non-root execution
	if !strings.Contains(text, "USER appuser") && !strings.Contains(text, "USER 10001") {
		t.Errorf("Dockerfile must run under an unprivileged non-root user")
	}

	// Check final stage minimalism
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
