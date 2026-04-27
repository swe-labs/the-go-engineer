// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package metrics

import (
	"math"
	"sort"
	"sync"
	"sync/atomic"
)

// Counter is a thread-safe monotonically increasing counter.
// It uses sync/atomic for lock-free concurrent access.
type Counter struct {
	value atomic.Int64
}

// Inc increments the counter by one.
func (c *Counter) Inc() { c.value.Add(1) }

// Add increments the counter by n.
func (c *Counter) Add(n int64) { c.value.Add(n) }

// Value returns the current counter value.
func (c *Counter) Value() int64 { return c.value.Load() }

// Histogram records observations into fixed buckets for latency
// distribution analysis. Each bucket counts how many observations
// were less than or equal to its upper bound.
//
// This is a simplified teaching implementation. Production systems
// use more efficient representations (e.g., HDR histograms or
// Prometheus-style cumulative buckets).
type Histogram struct {
	mu      sync.Mutex
	bounds  []float64 // sorted upper bounds
	buckets []int64   // count per bucket + 1 overflow
	count   int64
	sum     float64
}

// NewHistogram creates a histogram with the given bucket boundaries.
// Boundaries are sorted automatically.
func NewHistogram(bounds []float64) *Histogram {
	sorted := make([]float64, len(bounds))
	copy(sorted, bounds)
	sort.Float64s(sorted)

	return &Histogram{
		bounds:  sorted,
		buckets: make([]int64, len(sorted)+1), // +1 for overflow
	}
}

// Observe records a single value. The value is placed in the first
// bucket whose upper bound is greater than or equal to the value.
func (h *Histogram) Observe(v float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.count++
	h.sum += v

	for i, bound := range h.bounds {
		if v <= bound {
			h.buckets[i]++
			return
		}
	}
	// Overflow bucket: value exceeds all bounds.
	h.buckets[len(h.bounds)]++
}

// HistogramSnapshot holds a point-in-time copy of a histogram.
type HistogramSnapshot struct {
	Bounds  []float64 `json:"bounds"`
	Buckets []int64   `json:"buckets"`
	Count   int64     `json:"count"`
	Sum     float64   `json:"sum"`
}

// Snapshot returns a point-in-time copy of the histogram data.
func (h *Histogram) Snapshot() HistogramSnapshot {
	h.mu.Lock()
	defer h.mu.Unlock()

	bounds := make([]float64, len(h.bounds))
	copy(bounds, h.bounds)

	buckets := make([]int64, len(h.buckets))
	copy(buckets, h.buckets)

	return HistogramSnapshot{
		Bounds:  bounds,
		Buckets: buckets,
		Count:   h.count,
		Sum:     h.sum,
	}
}

// Registry holds named counters and histograms. It provides a single
// collection point for all application metrics.
type Registry struct {
	mu         sync.RWMutex
	counters   map[string]*Counter
	histograms map[string]*Histogram
}

// NewRegistry creates an empty metric registry.
func NewRegistry() *Registry {
	return &Registry{
		counters:   make(map[string]*Counter),
		histograms: make(map[string]*Histogram),
	}
}

// Counter returns the named counter, creating it if necessary.
func (r *Registry) Counter(name string) *Counter {
	r.mu.RLock()
	if c, ok := r.counters[name]; ok {
		r.mu.RUnlock()
		return c
	}
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()

	// Double-check after acquiring write lock.
	if c, ok := r.counters[name]; ok {
		return c
	}

	c := &Counter{}
	r.counters[name] = c
	return c
}

// Histogram returns the named histogram, creating it with the given
// bounds if necessary. If the histogram already exists, bounds are
// ignored.
func (r *Registry) Histogram(name string, bounds []float64) *Histogram {
	r.mu.RLock()
	if h, ok := r.histograms[name]; ok {
		r.mu.RUnlock()
		return h
	}
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()

	if h, ok := r.histograms[name]; ok {
		return h
	}

	h := NewHistogram(bounds)
	r.histograms[name] = h
	return h
}

// Snapshot returns all metrics as a serializable map.
func (r *Registry) Snapshot() map[string]any {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]any, len(r.counters)+len(r.histograms))

	for name, c := range r.counters {
		result[name] = c.Value()
	}
	for name, h := range r.histograms {
		result[name] = h.Snapshot()
	}

	return result
}

// DefaultHTTPBuckets are latency bucket boundaries in seconds,
// suitable for HTTP request duration tracking.
var DefaultHTTPBuckets = []float64{
	0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10,
}

// AppMetrics holds pre-registered application-level metrics.
type AppMetrics struct {
	Registry *Registry

	// HTTP
	HTTPRequestsTotal    *Counter
	HTTPRequestDuration  *Histogram
	HTTPResponses2xx     *Counter
	HTTPResponses4xx     *Counter
	HTTPResponses5xx     *Counter

	// Cache
	CacheHits   *Counter
	CacheMisses *Counter

	// Workers
	WorkerJobsProcessed *Counter
	WorkerJobsFailed    *Counter
}

// NewAppMetrics creates the application metrics with pre-registered
// counters and histograms.
func NewAppMetrics() *AppMetrics {
	r := NewRegistry()

	return &AppMetrics{
		Registry:            r,
		HTTPRequestsTotal:   r.Counter("http_requests_total"),
		HTTPRequestDuration: r.Histogram("http_request_duration_seconds", DefaultHTTPBuckets),
		HTTPResponses2xx:    r.Counter("http_responses_2xx"),
		HTTPResponses4xx:    r.Counter("http_responses_4xx"),
		HTTPResponses5xx:    r.Counter("http_responses_5xx"),
		CacheHits:           r.Counter("cache_hits_total"),
		CacheMisses:         r.Counter("cache_misses_total"),
		WorkerJobsProcessed: r.Counter("worker_jobs_processed_total"),
		WorkerJobsFailed:    r.Counter("worker_jobs_failed_total"),
	}
}

// ClassifyStatus increments the appropriate HTTP status class counter.
func (m *AppMetrics) ClassifyStatus(status int) {
	switch {
	case status >= 200 && status < 300:
		m.HTTPResponses2xx.Inc()
	case status >= 400 && status < 500:
		m.HTTPResponses4xx.Inc()
	case status >= 500:
		m.HTTPResponses5xx.Inc()
	}
}

// RoundToMs rounds a float64 seconds value to millisecond precision.
// Useful for human-readable duration display.
func RoundToMs(seconds float64) float64 {
	return math.Round(seconds*1000) / 1000
}
