# Prometheus

## Learning objective

Instrument a Go HTTP service with Prometheus counters, gauges, and histograms, expose a `/metrics` endpoint, and verify metric values programmatically.

## Why this matters

Metrics are the backbone of production observability. Prometheus is the CNCF-graduated standard for metrics collection in cloud-native environments, installed in over 80% of Kubernetes clusters. Without metrics, you deploy blind — you cannot detect slowdowns, error spikes, or resource exhaustion until users complain. Instrumenting Go services with Prometheus lets you answer questions like "how many requests per second?", "what is the p99 latency?", and "how many goroutines are running?" — in real time, with historical data.

## Mental model

Prometheus is like a meter reader who walks a route every 15 seconds (scrape interval), reading each household's meter (HTTP `/metrics` endpoint) and recording the numbers (time series) in a ledger (TSDB). The meter reader decides when to read — the household never calls to report a reading (pull model, not push). Each household has multiple meters: electricity (request count), water pressure (latency), temperature (memory usage). If a meter is unreadable (service down), the reader notes it. An analyst (Grafana) can query the ledger: "what was the average request rate in the last hour?" The reader only records what exists; if the household installs a new meter (new metric), the reader automatically discovers it at the next scrape.

The analogy breaks for histograms: a histogram is not a single meter but a set of 15+ sub-meters, each counting how many readings fell into a specific bucket (≤10ms, ≤25ms, ≤50ms, ≤100ms, ≤250ms, etc.). Prometheus never computes percentiles server-side; it stores bucket counts, and the query tool (PromQL) estimates percentiles by linear interpolation within the bucket.

## Core idea

Prometheus defines four core metric types:

- **Counter**: a monotonically increasing value that only ever goes up (resets on restart). Use for request counts, error counts, bytes served.
- **Gauge**: a value that can go up or down arbitrarily. Use for current memory usage, goroutine count, queue depth, temperature.
- **Histogram**: samples observations (like request latency) into configurable buckets, and counts observations in each bucket. Also tracks total sum and count of all observations. Use for latency and response sizes.
- **Summary**: similar to histogram but computes configurable quantiles directly on the client side. Less flexible than histograms; prefer histograms unless you need pre-computed quantiles.

Each metric is identified by its name and a set of label key-value pairs. A unique combination of name and labels defines a _time series_.

## Under the hood

The `prometheus` client library stores every metric in a `Registry`. Each metric implements a `Collector` interface with a `Collect(chan<- Metric)` method. When `/metrics` is scraped, `promhttp.Handler` iterates all registered collectors, calls `Collect` on each, and serializes the metrics in Prometheus text exposition format.

Internally, a `Counter` is an `atomic.Uint64` plus a mutex-protected map of label values. Incrementing a counter is a single atomic add on the hot path. A `Gauge` is similar but supports both increment and decrement. A `Histogram` maintains an array of `uint64` bucket counters, a total sum as `float64`, and an observation count — each `Observe` call finds the right bucket via binary search on the bucket boundaries and atomically increments the matching bucket counter and the total count.

The exposition format looks like:

```
http_requests_total{method="GET",status="200"} 1024
http_request_duration_seconds_bucket{le="0.1"} 512
http_request_duration_seconds_bucket{le="0.5"} 890
http_request_duration_seconds_sum 142.7
http_request_duration_seconds_count 1024
```

Prometheus server scrapes this text, parses it, and stores each time series in its TSDB (time series database). The TSDB stores data in 2-hour blocks, each containing an index (label-to-series mapping via postings lists), chunk files (compressed sample data), and metadata. Queries like `rate(http_requests_total[5m])` are evaluated by the PromQL engine which iterates matching series, reads raw samples from chunks, and computes the per-second rate.

## How Go uses it

The `prometheus/client_golang` library is the de facto standard for Go metrics. Most Go infrastructure projects (etcd, Kubernetes, Istio, Grafana itself) expose Prometheus metrics. The typical integration pattern:

1. Define metrics as package-level variables using `promauto` or manual registration.
2. Record metric values at relevant code points.
3. Expose `/metrics` via `promhttp.Handler()` on a dedicated HTTP server or endpoint.
4. Configure Prometheus server to scrape the endpoint.

Best practice: serve `/metrics` on a separate port (e.g., `:9091`) to avoid conflating application and metrics traffic, and use `promhttp.HandlerFor` with a custom registry to isolate metrics per-component.

## Go example

```go
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "path", "status"})

	inFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "http_in_flight_requests",
		Help: "Current number of in-flight requests",
	})

	requestDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Request latency distribution",
		Buckets: prometheus.DefBuckets,
	})
)

func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/api/hello", instrumentedHandler(http.HandlerFunc(helloHandler)))
	http.Handle("/api/echo", instrumentedHandler(http.HandlerFunc(echoHandler)))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func instrumentedHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inFlight.Inc()
		start := time.Now()

		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)

		duration := time.Since(start).Seconds()
		inFlight.Dec()
		requestsTotal.WithLabelValues(r.Method, r.URL.Path, fmt.Sprint(lrw.statusCode)).Inc()
		requestDuration.Observe(duration)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"hello"}`))
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"echo":"ok"}`))
}
```

## Step-by-step execution

1. `main` registers three metrics: a counter vector (labels: method, path, status), a gauge, and a histogram.
2. `promauto` auto-registers each metric with the default registry.
3. `instrumentedHandler` wraps every handler: increments gauge before, decrements after.
4. The histogram records the request duration in seconds.
5. The counter records the request with its method, path, and response status code labels.
6. `promhttp.Handler` at `/metrics` serializes all registered metrics in Prometheus text format.
7. A Prometheus server scrapes `http://service:8080/metrics` every 15s, parses the text, and stores samples in TSDB blocks.
8. Grafana queries PromQL to render dashboards: `rate(http_requests_total[5m])` for request rate, `histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))` for p99 latency.

## Common mistakes

- Mistake: Using `prometheus.MustRegister` in an init function that panics on duplicate registration.
  - Why it happens: When tests import the package, `init()` runs multiple times, re-registering the same metric and causing a panic.
  - Fix: Use `promauto` (safe for auto-registration) or check with `prometheus.Register` returning an `AlreadyRegisteredError`.

- Mistake: Forgetting to set a read/write timeout on the metrics server.
  - Why it happens: Generating the `/metrics` page for a service with 10K+ metrics takes seconds. A slow scrape client ties up a goroutine indefinitely.
  - Fix: Use `http.Server{ReadTimeout: 5s, WriteTimeout: 10s}` on the metrics mux.

- Mistake: Using high-cardinality labels like `user_id` or `request_id`.
  - Why it happens: Each unique label value creates a new time series. A label `user_id` with 10K users creates 10K time series per metric, exploding Prometheus memory and storage.
  - Fix: Label cardinality should be bounded (method, path, status_code). Never use unbounded values as labels.

## Debugging walkthrough

Consider a service where the gauge `http_in_flight_requests` is always 0, even during load testing.

```go
inFlight := promauto.NewGauge(prometheus.GaugeOpts{
	Name: "http_in_flight_requests",
	Help: "In-flight requests",
})

func handler(w http.ResponseWriter, r *http.Request) {
	inFlight.Inc()
	defer inFlight.Dec()
	// ... work ...
}
```

**Symptom**: `curl http://localhost:8080/metrics | grep in_flight` shows `http_in_flight_requests 0` even while `hey` (load tester) is sending requests.

**Investigation**:
1. Verify the handler is actually being called — add a `log.Println` statement.
2. Check if the handler is registered with `instrumentedHandler` wrapping.
3. Use `prometheus.MustRegister` check: verify the metric appears in `promhttp.Handler().ServeHTTP(w, r)` output.
4. Add a `time.Sleep(5 * time.Second)` in the handler and check `/metrics` during the sleep.

**Root cause**: The handler was registered directly (`http.HandleFunc("/api", handler)`) without the instrumentation wrapper. The `inFlight` metric was registered but never `Inc()`'d.

**Fix**: Wrap all handlers with `instrumentedHandler`.

## Production notes

- Serve `/metrics` on a separate port (e.g., `:9091`) to isolate metrics traffic from application traffic. Use `http.Server{Addr: ":9091", Handler: promhttp.Handler()}`.
- Use application-specific `prometheus.NewRegistry()` instead of the default registry. The default registry includes Go runtime metrics (goroutines, GC, memory) which are useful but should be managed separately.
- Set `promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}` to avoid returning errors that break the scrape.
- Consider `promhttp.HandlerOpts{EnableOpenMetrics: true}` for the OpenMetrics exposition format, which adds metadata and is the newer standard.
- Monitor scrape health: Prometheus's `up` metric is 1 when the target is reachable. Alert on `up == 0`.

## Performance implications

- Atomic counter increment: ~10ns per operation.
- Histogram `Observe`: ~50ns for default buckets (binary search on 12 buckets), scales with bucket count.
- Prometheus serialization is O(number of time series). A service with 10K time series takes ~100ms to serialize on each scrape.
- Label values are interned (shared across metrics) to reduce memory. But each unique label combination still creates a separate time series.
- For very high-throughput services (>100K req/s), use batch processing and periodically update metrics rather than incrementing per-request.

## Practice task

Write a function `recordRequest(method, path string, status int, duration time.Duration)` that records all three metrics (counter, gauge, histogram). Then write a test that:
1. Creates a new `prometheus.NewRegistry()`.
2. Registers the three metrics on the custom registry.
3. Calls `recordRequest` with sample data.
4. Gathers metrics from the registry using `registry.Gather()`.
5. Asserts that the counter value matches the number of calls and the gauge value matches expected in-flight requests.

## Tests / verification

```go
package main

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestRecordRequest(t *testing.T) {
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

	record := func(method, path string, status int, duration time.Duration) {
		gauge.Inc()
		counter.WithLabelValues(method, path, string(rune(status))).Inc()
		histogram.Observe(duration.Seconds())
		gauge.Dec()
	}

	record("GET", "/api/users", 200, 50*time.Millisecond)
	record("POST", "/api/users", 201, 120*time.Millisecond)
	record("GET", "/api/users", 500, 5*time.Millisecond)

	// Assert the counter value using testutil
	expected := 3
	got := testutil.CollectAndCount(counter)
	if got != expected {
		t.Errorf("expected %d metrics collected, got %d", expected, got)
	}
}
```

```bash
go test ./curriculum/modules/12-observability-diagnostics/lessons/07-prometheus -v
go run ./curriculum/modules/12-observability-diagnostics/lessons/07-prometheus
```

## Review questions

1. Why does Prometheus use a pull model instead of a push model for metrics collection?
2. A histogram with buckets `[.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10]` receives 100 observations all under 5ms. Which bucket counters are incremented and what does `histogram_quantile(0.95, ...)` return?
3. What happens to Prometheus storage if a metric with label `user_id="abc123"` records a new user ID every minute for 24 hours?
4. Why should you avoid registering metrics in an `init()` function when writing library code?
5. A gauge `http_in_flight_requests` shows negative values at scrape time. What could cause this?

## NEXT UP

Tracing — moving from aggregate metrics to per-request distributed traces.
