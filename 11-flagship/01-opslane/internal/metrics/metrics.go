// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

// Package metrics provides in-memory metrics collection for the Opslane backend.
// It includes thread-safe counters and histograms for HTTP requests, cache hits/misses,
// and worker job processing.
package metrics

import (
	"math"
	"sort"
	"sync"
	"sync/atomic"
)

// Counter (Struct): thread-safe monotonically increasing counter using sync/atomic
type Counter struct {
	value atomic.Int64
}

// Inc (Method): increments the counter by one
func (c *Counter) Inc() { c.value.Add(1) }

// Add (Method): increments the counter by n
func (c *Counter) Add(n int64) { c.value.Add(n) }

// Value (Method): returns the current counter value
func (c *Counter) Value() int64 { return c.value.Load() }

// Histogram (Struct): records observations into fixed buckets for latency distribution analysis
type Histogram struct {
	mu      sync.Mutex
	bounds  []float64 // sorted upper bounds
	buckets []int64   // count per bucket + 1 overflow
	count   int64
	sum     float64
}

// NewHistogram (Constructor): creates a histogram with automatically sorted bucket boundaries
func NewHistogram(bounds []float64) *Histogram {
	sorted := make([]float64, len(bounds))
	copy(sorted, bounds)
	sort.Float64s(sorted)

	return &Histogram{
		bounds:  sorted,
		buckets: make([]int64, len(sorted)+1), // +1 for overflow
	}
}

// Observe (Method): records a single value into the appropriate bucket
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

// HistogramSnapshot (Struct): point-in-time copy of a histogram for serialization
type HistogramSnapshot struct {
	Bounds  []float64 `json:"bounds"`
	Buckets []int64   `json:"buckets"`
	Count   int64     `json:"count"`
	Sum     float64   `json:"sum"`
}

// Snapshot (Method): returns a point-in-time copy of the histogram data
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

// Registry (Struct): holds named counters and histograms as a single metrics collection point
type Registry struct {
	mu         sync.RWMutex
	counters   map[string]*Counter
	histograms map[string]*Histogram
}

// NewRegistry (Constructor): creates an empty metric registry
func NewRegistry() *Registry {
	return &Registry{
		counters:   make(map[string]*Counter),
		histograms: make(map[string]*Histogram),
	}
}

// Counter (Method): returns the named counter, creating it if necessary (thread-safe)
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

// Histogram (Method): returns the named histogram, creating it with the given bounds if necessary
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

// Snapshot (Method): returns all metrics as a serializable map
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

// DefaultHTTPBuckets (Slice): latency bucket boundaries in seconds for HTTP request duration tracking
var DefaultHTTPBuckets = []float64{
	0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10,
}

// AppMetrics (Struct): holds pre-registered application-level metrics for HTTP, cache, and workers
type AppMetrics struct {
	Registry *Registry

	// HTTP
	HTTPRequestsTotal   *Counter
	HTTPRequestDuration *Histogram
	HTTPResponses2xx    *Counter
	HTTPResponses4xx    *Counter
	HTTPResponses5xx    *Counter

	// Cache
	CacheHits   *Counter
	CacheMisses *Counter

	// Workers
	WorkerJobsProcessed *Counter
	WorkerJobsFailed    *Counter
}

// NewAppMetrics (Constructor): creates the application metrics with pre-registered counters and histograms
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

// ClassifyStatus (Method): increments the appropriate HTTP status class counter (2xx, 4xx, 5xx)
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

// RoundToMs (Function): rounds a float64 seconds value to millisecond precision
func RoundToMs(seconds float64) float64 {
	return math.Round(seconds*1000) / 1000
}
