package main

import (
	"expvar"
	"fmt"
	"sync/atomic"
)

// MetricsStore tracks request metrics in memory.
type MetricsStore struct {
	requestCount  atomic.Int64
	errorCount    atomic.Int64
	totalDuration atomic.Int64
}

// NewMetricsStore returns an initialized MetricsStore.
func NewMetricsStore() *MetricsStore {
	return &MetricsStore{}
}

// RecordRequest records a request with its duration and error status.
func (m *MetricsStore) RecordRequest(durationMs int64, isError bool) {
	m.requestCount.Add(1)
	m.totalDuration.Add(durationMs)
	if isError {
		m.errorCount.Add(1)
	}
}

// Snapshot returns the current metrics values.
func (m *MetricsStore) Snapshot() (requests, errors, avgDurationMs int64) {
	reqs := m.requestCount.Load()
	errs := m.errorCount.Load()
	total := m.totalDuration.Load()
	if reqs > 0 {
		avgDurationMs = total / reqs
	}
	return reqs, errs, avgDurationMs
}

// expvar metrics (global, for production exposure).
var (
	expRequests = expvar.NewInt("http_requests_total")
	expErrors   = expvar.NewInt("http_errors_total")
	expDuration = expvar.NewMap("http_request_duration_ms")
)

func main() {
	store := NewMetricsStore()

	store.RecordRequest(42, false)
	store.RecordRequest(150, false)
	store.RecordRequest(5200, true)

	reqs, errs, avg := store.Snapshot()
	fmt.Printf("Requests: %d, Errors: %d, Avg duration: %dms\n", reqs, errs, avg)

	expRequests.Set(3)
	expErrors.Set(1)
	expDuration.Add("p50", 42)
	fmt.Println("expvar exported at /debug/vars")
}
