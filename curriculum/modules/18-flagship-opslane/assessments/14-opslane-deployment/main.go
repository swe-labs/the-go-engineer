package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type ComponentHealth struct {
	Name      string        `json:"name"`
	Healthy   bool          `json:"healthy"`
	Latency   time.Duration `json:"latency_ns"`
	LatencyMs float64       `json:"latency_ms"`
	Error     string        `json:"error,omitempty"`
}

type ReadinessReport struct {
	Status     string            `json:"status"` // "UP", "DEGRADED"
	Timestamp  time.Time         `json:"timestamp"`
	DurationMs float64           `json:"duration_ms"`
	Components []ComponentHealth `json:"components"`
}

// HealthChecker verifies dependent services.
type HealthChecker interface {
	CheckHealth(ctx context.Context) ComponentHealth
}

type GenericDependencyChecker struct {
	ComponentName string
	SimulateError error
	SimulateDelay time.Duration
}

func (c *GenericDependencyChecker) CheckHealth(ctx context.Context) ComponentHealth {
	start := time.Now()
	if c.SimulateDelay > 0 {
		time.Sleep(c.SimulateDelay)
	}
	latency := time.Since(start)

	health := ComponentHealth{
		Name:      c.ComponentName,
		Healthy:   c.SimulateError == nil,
		Latency:   latency,
		LatencyMs: float64(latency.Microseconds()) / 1000.0,
	}
	if c.SimulateError != nil {
		health.Error = c.SimulateError.Error()
	}
	return health
}

// DeploymentHealthService manages Kubernetes/Docker liveness and readiness probes.
type DeploymentHealthService struct {
	checkers []HealthChecker
}

func NewDeploymentHealthService(checkers []HealthChecker) *DeploymentHealthService {
	return &DeploymentHealthService{
		checkers: checkers,
	}
}

func (s *DeploymentHealthService) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", s.handleLiveness)
	mux.HandleFunc("GET /health/ready", s.handleReadiness)
	return mux
}

// handleLiveness answers Kubernetes liveness probe (indicates container is running).
func (s *DeploymentHealthService) handleLiveness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ALIVE"}`))
}

// handleReadiness answers Kubernetes readiness probe (verifies DB, Cache, Workers).
func (s *DeploymentHealthService) handleReadiness(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// 50ms deadline for readiness checks
	ctx, cancel := context.WithTimeout(r.Context(), 50*time.Millisecond)
	defer cancel()

	results := make([]ComponentHealth, len(s.checkers))
	var wg sync.WaitGroup
	wg.Add(len(s.checkers))

	for i, checker := range s.checkers {
		go func(idx int, ch HealthChecker) {
			defer wg.Done()
			results[idx] = ch.CheckHealth(ctx)
		}(i, checker)
	}

	wg.Wait()

	allHealthy := true
	for _, res := range results {
		if !res.Healthy {
			allHealthy = false
			break
		}
	}

	totalDuration := time.Since(start)
	report := ReadinessReport{
		Status:     "UP",
		Timestamp:  time.Now().UTC(),
		DurationMs: float64(totalDuration.Microseconds()) / 1000.0,
		Components: results,
	}

	w.Header().Set("Content-Type", "application/json")

	if !allHealthy {
		report.Status = "DEGRADED"
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	_ = json.NewEncoder(w).Encode(report)
}

func main() {
	checkers := []HealthChecker{
		&GenericDependencyChecker{ComponentName: "postgresql"},
		&GenericDependencyChecker{ComponentName: "redis_cache"},
		&GenericDependencyChecker{ComponentName: "worker_pool"},
	}

	svc := NewDeploymentHealthService(checkers)

	fmt.Println("Opslane Deployment Health Probes listening on :8080...")
	_ = http.ListenAndServe(":8080", svc.Routes())
}
