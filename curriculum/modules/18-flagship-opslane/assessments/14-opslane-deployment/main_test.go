package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDeployment_ReadinessWithin50msAllHealthy(t *testing.T) {
	// Core practice task requirement:
	// Write a health check probe test that verifies the deployment health endpoint
	// returns 200 OK within 50ms and reports all dependent services (database, cache, workers) as healthy.
	checkers := []HealthChecker{
		&GenericDependencyChecker{ComponentName: "database", SimulateDelay: 5 * time.Millisecond},
		&GenericDependencyChecker{ComponentName: "cache", SimulateDelay: 3 * time.Millisecond},
		&GenericDependencyChecker{ComponentName: "workers", SimulateDelay: 2 * time.Millisecond},
	}

	service := NewDeploymentHealthService(checkers)
	router := service.Routes()

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	rr := httptest.NewRecorder()

	start := time.Now()
	router.ServeHTTP(rr, req)
	elapsed := time.Since(start)

	// 1. Verify 200 OK
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from readiness probe, got %d: %s", rr.Code, rr.Body.String())
	}

	// 2. Verify probe completed strictly within 50ms SLA
	if elapsed > 50*time.Millisecond {
		t.Fatalf("PROBE SLA BREACH: Readiness check took %v, must be < 50ms", elapsed)
	}

	var report ReadinessReport
	if err := json.NewDecoder(rr.Body).Decode(&report); err != nil {
		t.Fatalf("failed to decode readiness report: %v", err)
	}

	if report.Status != "UP" {
		t.Errorf("expected status 'UP', got %q", report.Status)
	}

	// 3. Verify all 3 components are reported healthy
	if len(report.Components) != 3 {
		t.Fatalf("expected 3 component statuses, got %d", len(report.Components))
	}

	for _, comp := range report.Components {
		if !comp.Healthy {
			t.Errorf("expected component %q to be healthy, but reported error: %s", comp.Name, comp.Error)
		}
	}
}

func TestDeployment_ReadinessDegradedWhenDependencyFails(t *testing.T) {
	checkers := []HealthChecker{
		&GenericDependencyChecker{ComponentName: "database", SimulateError: errors.New("connection refused to postgres:5432")},
		&GenericDependencyChecker{ComponentName: "cache"},
		&GenericDependencyChecker{ComponentName: "workers"},
	}

	service := NewDeploymentHealthService(checkers)
	router := service.Routes()

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	// Must return 503 Service Unavailable when any dependency is degraded
	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 Service Unavailable on degraded dependency, got %d", rr.Code)
	}

	var report ReadinessReport
	_ = json.NewDecoder(rr.Body).Decode(&report)

	if report.Status != "DEGRADED" {
		t.Errorf("expected status 'DEGRADED', got %q", report.Status)
	}

	var dbHealth *ComponentHealth
	for _, c := range report.Components {
		if c.Name == "database" {
			dbHealth = &c
			break
		}
	}

	if dbHealth == nil || dbHealth.Healthy || dbHealth.Error == "" {
		t.Errorf("expected database to report unhealthy with error message, got: %+v", dbHealth)
	}
}

func TestDeployment_LivenessFast(t *testing.T) {
	service := NewDeploymentHealthService(nil)
	router := service.Routes()

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from liveness probe, got %d", rr.Code)
	}
}
