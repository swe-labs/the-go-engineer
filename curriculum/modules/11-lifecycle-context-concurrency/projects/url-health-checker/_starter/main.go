package main

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// EndpointStatus represents health status of a monitored target.
type EndpointStatus string

const (
	StatusUp   EndpointStatus = "UP"
	StatusDown EndpointStatus = "DOWN"
)

// CheckResult records an individual probe measurement.
type CheckResult struct {
	URL        string         `json:"url"`
	Status     EndpointStatus `json:"status"`
	StatusCode int            `json:"status_code"`
	Latency    time.Duration  `json:"latency"`
	Error      string         `json:"error,omitempty"`
	CheckedAt  time.Time      `json:"checked_at"`
}

// TransitionEvent logs changes in health status (e.g. UP -> DOWN).
type TransitionEvent struct {
	URL        string         `json:"url"`
	FromStatus EndpointStatus `json:"from_status"`
	ToStatus   EndpointStatus `json:"to_status"`
	OccurredAt time.Time      `json:"occurred_at"`
}

// CheckerConfig holds parameters for the concurrent probe engine.
type CheckerConfig struct {
	Workers int
	Timeout time.Duration
	Client  *http.Client
}

// HealthChecker coordinates concurrent URL probing and transition tracking.
type HealthChecker struct {
	workers int
	timeout time.Duration
	client  *http.Client

	mu          sync.RWMutex
	lastStatus  map[string]EndpointStatus
	transitions []TransitionEvent
}

// NewHealthChecker initializes a new HealthChecker.
func NewHealthChecker(cfg CheckerConfig) *HealthChecker {
	if cfg.Workers <= 0 {
		cfg.Workers = 4
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 3 * time.Second
	}
	if cfg.Client == nil {
		cfg.Client = &http.Client{Timeout: cfg.Timeout}
	}

	return &HealthChecker{
		workers:    cfg.Workers,
		timeout:    cfg.Timeout,
		client:     cfg.Client,
		lastStatus: make(map[string]EndpointStatus),
	}
}

// ProbeTarget performs a single HTTP probe against a URL with deadline context.
func (h *HealthChecker) ProbeTarget(ctx context.Context, targetURL string) CheckResult {
	// TODO: Create child context with timeout (h.timeout).
	// TODO: Create http.RequestWithContext and execute with h.client.Do.
	// TODO: Measure latency and record StatusUp for 2xx/3xx, StatusDown for error/5xx.
	// TODO: Call recordTransition to detect state changes.
	return CheckResult{URL: targetURL, Status: StatusDown}
}

func (h *HealthChecker) recordTransition(url string, currentStatus EndpointStatus) {
	// TODO: Under h.mu.Lock(), check if last status differs from currentStatus.
	// TODO: If changed, append TransitionEvent to h.transitions.
	// TODO: Update h.lastStatus[url].
}

// CheckAll checks a list of URLs concurrently using a bounded worker pool.
func (h *HealthChecker) CheckAll(ctx context.Context, urls []string) []CheckResult {
	// TODO: Use a buffered channel to distribute work across min(h.workers, len(urls)) goroutines.
	// TODO: Coordinate completion using sync.WaitGroup.
	// TODO: Save results to output slice indexed by original position.
	return nil
}

// Transitions returns all recorded status transition events.
func (h *HealthChecker) Transitions() []TransitionEvent {
	h.mu.RLock()
	defer h.mu.RUnlock()
	copied := make([]TransitionEvent, len(h.transitions))
	copy(copied, h.transitions)
	return copied
}

func main() {
	checker := NewHealthChecker(CheckerConfig{
		Workers: 3,
		Timeout: 2 * time.Second,
	})

	urls := []string{
		"https://google.com",
		"https://github.com",
		"https://go.dev",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Info("starting concurrent health check", "total_targets", len(urls))
	results := checker.CheckAll(ctx, urls)

	for _, res := range results {
		slog.Info("probe result",
			"url", res.URL,
			"status", res.Status,
			"status_code", res.StatusCode,
			"latency", res.Latency.String(),
		)
	}
}
