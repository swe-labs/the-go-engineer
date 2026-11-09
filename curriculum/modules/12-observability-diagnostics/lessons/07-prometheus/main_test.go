package main

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestRecordRequestMetrics(t *testing.T) {
	reg := prometheus.NewRegistry()

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total requests",
	}, []string{"method", "path", "status"})
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "http_in_flight_requests",
		Help: "In-flight requests",
	})
	histogram := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Request duration",
		Buckets: prometheus.DefBuckets,
	})
	reg.MustRegister(counter, gauge, histogram)

	tests := []struct {
		name     string
		method   string
		path     string
		status   int
		duration time.Duration
	}{
		{"successful get", "GET", "/api/users", 200, 50 * time.Millisecond},
		{"created post", "POST", "/api/users", 201, 120 * time.Millisecond},
		{"server error", "GET", "/api/users", 500, 5 * time.Millisecond},
		{"not found", "GET", "/api/unknown", 404, 2 * time.Millisecond},
		{"bad request", "POST", "/api/data", 400, 30 * time.Millisecond},
	}

	for _, tc := range tests {
		gauge.Inc()
		counter.WithLabelValues(tc.method, tc.path, statusText(tc.status)).Inc()
		histogram.Observe(tc.duration.Seconds())
		gauge.Dec()
	}

	expectedCount := len(tests)
	gotCount := testutil.CollectAndCount(counter)
	if gotCount != expectedCount {
		t.Errorf("expected %d metric label combinations, got %d", expectedCount, gotCount)
	}

	c, err := testutil.GatherAndCount(reg)
	if err != nil {
		t.Fatalf("failed to gather metrics: %v", err)
	}
	if c < 3 {
		t.Errorf("expected at least 3 metric families, got %d", c)
	}
}
