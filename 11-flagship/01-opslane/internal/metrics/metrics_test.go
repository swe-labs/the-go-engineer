// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package metrics

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestCounterIncrement(t *testing.T) {
	t.Parallel()

	var c Counter
	c.Inc()
	c.Inc()
	c.Add(3)

	if got := c.Value(); got != 5 {
		t.Fatalf("Counter.Value() = %d, want 5", got)
	}
}

func TestCounterConcurrentAccess(t *testing.T) {
	t.Parallel()

	var c Counter
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Inc()
		}()
	}
	wg.Wait()

	if got := c.Value(); got != 100 {
		t.Fatalf("Counter.Value() = %d, want 100 after concurrent access", got)
	}
}

func TestHistogramBucketDistribution(t *testing.T) {
	t.Parallel()

	h := NewHistogram([]float64{1, 5, 10})

	h.Observe(0.5)  // bucket 0 (<=1)
	h.Observe(3)    // bucket 1 (<=5)
	h.Observe(7)    // bucket 2 (<=10)
	h.Observe(15)   // overflow bucket
	h.Observe(1)    // bucket 0 (<=1)

	snap := h.Snapshot()

	if snap.Count != 5 {
		t.Fatalf("count = %d, want 5", snap.Count)
	}
	if snap.Sum != 26.5 {
		t.Fatalf("sum = %f, want 26.5", snap.Sum)
	}

	// buckets: [<=1, <=5, <=10, >10]
	want := []int64{2, 1, 1, 1}
	for i, w := range want {
		if snap.Buckets[i] != w {
			t.Errorf("bucket[%d] = %d, want %d", i, snap.Buckets[i], w)
		}
	}
}

func TestRegistryCounterLazyCreation(t *testing.T) {
	t.Parallel()

	r := NewRegistry()
	c1 := r.Counter("test_counter")
	c2 := r.Counter("test_counter")

	if c1 != c2 {
		t.Fatal("expected same counter instance for same name")
	}

	c1.Inc()
	if got := c2.Value(); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestRegistryHistogramLazyCreation(t *testing.T) {
	t.Parallel()

	r := NewRegistry()
	h1 := r.Histogram("test_hist", []float64{1, 5})
	h2 := r.Histogram("test_hist", nil)

	if h1 != h2 {
		t.Fatal("expected same histogram instance for same name")
	}
}

func TestRegistrySnapshot(t *testing.T) {
	t.Parallel()

	r := NewRegistry()
	r.Counter("requests").Add(42)
	r.Histogram("latency", []float64{0.1}).Observe(0.05)

	snap := r.Snapshot()

	if v, ok := snap["requests"].(int64); !ok || v != 42 {
		t.Fatalf("requests = %v, want 42", snap["requests"])
	}
	if _, ok := snap["latency"].(HistogramSnapshot); !ok {
		t.Fatal("expected HistogramSnapshot for latency")
	}
}

func TestAppMetricsClassifyStatus(t *testing.T) {
	t.Parallel()

	m := NewAppMetrics()
	m.ClassifyStatus(200)
	m.ClassifyStatus(201)
	m.ClassifyStatus(404)
	m.ClassifyStatus(500)
	m.ClassifyStatus(502)

	if got := m.HTTPResponses2xx.Value(); got != 2 {
		t.Fatalf("2xx = %d, want 2", got)
	}
	if got := m.HTTPResponses4xx.Value(); got != 1 {
		t.Fatalf("4xx = %d, want 1", got)
	}
	if got := m.HTTPResponses5xx.Value(); got != 2 {
		t.Fatalf("5xx = %d, want 2", got)
	}
}

func TestHTTPMetricsMiddleware(t *testing.T) {
	t.Parallel()

	m := NewAppMetrics()

	handler := HTTPMetrics(m)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got := m.HTTPRequestsTotal.Value(); got != 1 {
		t.Fatalf("requests_total = %d, want 1", got)
	}
	if got := m.HTTPResponses2xx.Value(); got != 1 {
		t.Fatalf("responses_2xx = %d, want 1", got)
	}

	snap := m.HTTPRequestDuration.Snapshot()
	if snap.Count != 1 {
		t.Fatalf("duration count = %d, want 1", snap.Count)
	}
}
