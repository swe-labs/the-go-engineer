package main

import (
	"context"
	"fmt"
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
	probeCtx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	start := time.Now()
	req, err := http.NewRequestWithContext(probeCtx, http.MethodGet, targetURL, nil)
	if err != nil {
		return CheckResult{
			URL:       targetURL,
			Status:    StatusDown,
			Error:     err.Error(),
			CheckedAt: time.Now().UTC(),
		}
	}

	resp, err := h.client.Do(req)
	latency := time.Since(start)

	res := CheckResult{
		URL:       targetURL,
		Latency:   latency,
		CheckedAt: time.Now().UTC(),
	}

	if err != nil {
		res.Status = StatusDown
		res.Error = err.Error()
	} else {
		defer resp.Body.Close()
		res.StatusCode = resp.StatusCode
		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			res.Status = StatusUp
		} else {
			res.Status = StatusDown
			res.Error = fmt.Sprintf("non-2xx HTTP status: %d", resp.StatusCode)
		}
	}

	h.recordTransition(res.URL, res.Status)
	return res
}

func (h *HealthChecker) recordTransition(url string, currentStatus EndpointStatus) {
	h.mu.Lock()
	defer h.mu.Unlock()

	last, exists := h.lastStatus[url]
	if exists && last != currentStatus {
		h.transitions = append(h.transitions, TransitionEvent{
			URL:        url,
			FromStatus: last,
			ToStatus:   currentStatus,
			OccurredAt: time.Now().UTC(),
		})
	}
	h.lastStatus[url] = currentStatus
}

// CheckAll checks a list of URLs concurrently using a bounded worker pool.
func (h *HealthChecker) CheckAll(ctx context.Context, urls []string) []CheckResult {
	if len(urls) == 0 {
		return nil
	}

	results := make([]CheckResult, len(urls))
	urlChan := make(chan struct {
		index int
		url   string
	}, len(urls))

	for idx, u := range urls {
		urlChan <- struct {
			index int
			url   string
		}{index: idx, url: u}
	}
	close(urlChan)

	var wg sync.WaitGroup
	workers := h.workers
	if workers > len(urls) {
		workers = len(urls)
	}

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range urlChan {
				select {
				case <-ctx.Done():
					results[item.index] = CheckResult{
						URL:       item.url,
						Status:    StatusDown,
						Error:     ctx.Err().Error(),
						CheckedAt: time.Now().UTC(),
					}
				default:
					results[item.index] = h.ProbeTarget(ctx, item.url)
				}
			}
		}()
	}

	wg.Wait()
	return results
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
