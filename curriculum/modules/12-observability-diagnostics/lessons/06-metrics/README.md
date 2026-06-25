# Metrics

## Learning objective

Instrument a Go service with custom metrics (counter, gauge) using `expvar`, and implement a thread-safe in-memory metrics store that tracks request count, error count, and average latency.

## Why this matters

Logs tell you what happened to a specific request. Metrics tell you what is happening to your system right now and how it is changing over time. A single "database timeout" log line is a data point. A metric showing "database timeout rate: 5% and rising over the last 10 minutes" is an actionable signal that triggers an alert, a dashboard update, and an incident response. Without metrics, you are flying blind: you cannot answer "is the system healthy?", "is the deploy causing problems?", or "is latency increasing?" with log search alone — you need aggregated numeric data.

## Mental model

Metrics are the vital signs of a service. Like heart rate (request rate), blood pressure (latency percentiles), temperature (error rate), and oxygen level (saturation). Each vital sign is a time series: a sequence of (timestamp, value) pairs. A single number tells you nothing meaningful; the trend over time tells you everything. A spike in heart rate during a code deploy is expected; a spike in temperature is not.

Metrics answer "is the service healthy?" and "what changed?". They do not answer "why?" (that is logs) or "where exactly?" (that is traces). The three together give a complete picture.

## Core idea

A metric is a numeric value that changes over time, exposed for external collection (scraping) by a monitoring system. The four fundamental metric types:

- **Counter**: A monotonically increasing value that only goes up (request count, bytes sent, errors occurred). Reset to zero on process restart. Used to compute rates over time.
- **Gauge**: A value that can go up or down (active connections, memory usage, queue depth). Represents a snapshot at a point in time.
- **Histogram**: A set of counters for pre-defined value buckets (request duration in ms: bucket ≤10, ≤50, ≤100, ≤500, ≤1000). Used to compute percentiles (p50, p95, p99).
- **Summary**: Like a histogram but calculates quantiles on the client side. Less common in Go; histograms are preferred.

Go's standard library provides `expvar` for exporting counters and gauges as JSON at `/debug/vars`. For histograms and advanced metric types, use `prometheus/client_golang` (lesson 07).

## Under the hood

`expvar` stores metrics in a global registry (`expvar.Do` iterates all registered vars). The key types:

- `expvar.Int`: An atomic int64. `Add(delta)` calls `atomic.AddInt64`. `Set(val)` calls `atomic.StoreInt64`. String() returns the JSON-encoded integer.
- `expvar.Map`: A concurrent map from string to expvar.Var. Keys are sorted during serialization. `Add(key, delta)` creates or updates an Int entry.
- `expvar.Func`: A function that returns a value at scrape time. Used for runtime metrics like `runtime.NumGoroutine()`.

When the `/debug/vars` HTTP endpoint is requested (registered by `expvar.Handler()`), expvar calls `expvar.Do` which iterates all registered vars and calls `String()` on each, producing a JSON object like:

```json
{
  "cmdline": ["..."],
  "memstats": {"Alloc": 123456, ...},
  "http_requests_total": 42,
  "http_errors_total": 3
}
```

The memstats block is automatically registered by `expvar.NewExpvar` in the `runtime` package, exposing `MemStats` fields.

## How Go uses it

`expvar` is the standard library approach for metrics in Go:

```go
import "expvar"

var (
    requestsTotal = expvar.NewInt("http_requests_total")
    errorsTotal   = expvar.NewInt("http_errors_total")
    durationMap   = expvar.NewMap("http_request_duration_ms")
)

func handler(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    requestsTotal.Add(1)
    defer func() {
        durationMap.Add("all", time.Since(start).Milliseconds())
        if r.URL.Path == "/api/orders" {
            durationMap.Add("orders", time.Since(start).Milliseconds())
        }
    }()
    // ... handler logic
}
```

The metrics are available at `/debug/vars` when the default `expvar.Handler()` is registered (it is automatically registered on the default ServeMux at `/debug/vars`).

For production, most Go services use `prometheus/client_golang` (Prometheus) instead of raw expvar, because Prometheus provides histograms, summaries, labels, and a standardized text format scraped by Prometheus server. But expvar is useful for simple use cases and for understanding the underlying concepts.

## Go example

```go
package main

import (
	"expvar"
	"fmt"
	"sync/atomic"
)

type MetricsStore struct {
	requestCount  atomic.Int64
	errorCount    atomic.Int64
	totalDuration atomic.Int64
}

func NewMetricsStore() *MetricsStore {
	return &MetricsStore{}
}

func (m *MetricsStore) RecordRequest(durationMs int64, isError bool) {
	m.requestCount.Add(1)
	m.totalDuration.Add(durationMs)
	if isError {
		m.errorCount.Add(1)
	}
}

func (m *MetricsStore) Snapshot() (requests, errors, avgDurationMs int64) {
	reqs := m.requestCount.Load()
	errs := m.errorCount.Load()
	total := m.totalDuration.Load()
	if reqs > 0 {
		avgDurationMs = total / reqs
	}
	return reqs, errs, avgDurationMs
}

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
```

Run with `go run .` to see the metrics snapshot. The `MetricsStore` is thread-safe (uses atomic operations) and requires no external dependencies.

## Step-by-step execution

For `store.RecordRequest(150, false)`:

1. `m.requestCount.Add(1)` — atomic increment of the request count. Current count is now 2 (two successful requests).
2. `m.totalDuration.Add(150)` — atomic add of 150ms to the running total. Total is now 192 (42 + 150).
3. `isError` is false, so `m.errorCount.Add(1)` is NOT called. Error count stays at 0.
4. On `Snapshot()`: `requestCount.Load()` returns 2, `errorCount.Load()` returns 0, `totalDuration.Load()` returns 192. Average = 192 / 2 = 96ms.
5. Function returns (2, 0, 96).

## Common mistakes

- **Using gauge for monotonically increasing values**: A gauge can go up and down (CPU, memory). A counter only goes up (request count). Using a gauge for request count means a process restart or metric reset sets the value to 0, making rate calculations incorrect. Always use a counter for cumulative totals.

- **No label dimensions**: A single metric `http_requests_total` with no labels tells you the total request count but not which endpoint, method, or status. Without labels, a 5xx spike on `/api/login` is invisible in the aggregate. Always add meaningful dimensions: endpoint, method, status_code. (Note: expvar does not support labels natively; use Prometheus for labeled metrics.)

- **Calling metric inc/decr on the request hot path**: If the metric backend uses a mutex (like Prometheus histograms), calling `Observe()` on every request creates contention under high concurrency. For hot-path metrics, use an atomic counter (like `atomic.Int64`) and batch-export periodically.

- **Forgetting to export metrics**: Creating `expvar.NewInt("name")` registers the metric in the global `expvar` registry, but it is only visible if the `/debug/vars` HTTP endpoint is served. The default ServeMux does not automatically expose this — developers must start an HTTP server or register `expvar.Handler()` on a custom mux.

- **Not monitoring goroutine count**: `expvar` automatically exports `memstats`, but goroutine count is only exported if explicitly registered via `expvar.NewFunc("num_goroutine", func() interface{} { return runtime.NumGoroutine() })`. A goroutine leak is invisible without this metric.

## Debugging walkthrough

Consider a Go HTTP service with no metrics:

```go
func main() {
    http.HandleFunc("/api/orders", ordersHandler)
    http.ListenAndServe(":8080", nil)
}
```

**Symptom**: Users report that the service feels slow. Operations team sees CPU at 30%, memory at 40%, no alerts firing.

**Investigation**: Without application metrics, you have no data about request rate, error rate, or latency. You guess: maybe it is a database issue? A dependency slow-down? A goroutine leak? You add a log line to every handler, redeploy, wait for the slow period, grep the logs, compute average latency manually, and repeat for each endpoint. This takes 2 hours.

**With metrics**, the `/debug/vars` endpoint shows:

```json
{
  "http_requests_total": 15000,
  "http_errors_total": 0,
  "http_request_duration_ms": {"all": 750000}
}
```

You compute average latency: 750000 / 15000 = 50ms. But the p99 latency (which requires a histogram, not available in expvar) would show 5000ms — one endpoint is slow but others are fast. You add a Prometheus histogram in lesson 07.

**Root cause**: No metrics means no visibility into which endpoint is slow, which status codes are failing, or how many requests are affected. Adding a simple counter per endpoint reduces investigation time from 2 hours to 2 minutes.

**Fix**: Add at minimum: request count by endpoint, error count by endpoint, and request duration. Use atomic operations for counters and a histogram implementation for latency percentiles.

## Production notes

Metrics are the foundation of production observability. Every Go service exports metrics:

- **Kubernetes**: Uses metrics for horizontal pod autoscaling (HPA) based on CPU/memory and custom metrics (requests per second).
- **Prometheus ecosystem**: Metrics are scraped every 15–60 seconds. Alerts fire when metrics violate SLOs (error rate > 1% for 5 minutes, p99 latency > 500ms for 10 minutes).
- **SRE practice**: SLOs are defined in terms of metrics. Example: "99.9% of HTTP requests complete in < 200ms over a 30-day rolling window." This is impossible to measure without metrics.

Production metrics minimum set (the "Golden Signals"):

1. **Latency** — time to serve a request (by endpoint, method, status).
2. **Traffic** — requests per second (by endpoint).
3. **Errors** — rate of failed requests (by endpoint, error type).
4. **Saturation** — how "full" the service is (goroutine count, queue depth, connection pool utilization).

## Performance implications

Different metric types have different costs:

- **Atomic counter (`expvar.Int`)**: `Add` is a single `atomic.AddInt64` instruction (~50ns). No allocations. The fastest metric type.
- **Atomic gauge (`expvar.Int` with Set)**: `Set` is a single `atomic.StoreInt64` instruction (~20ns). Also zero allocations.
- **expvar.Map**: `Add(key, delta)` involves a map lookup (amortized O(1)) plus an atomic add on the inner Int. If the key does not exist, it creates a new Int, which allocates.
- **Histogram**: Requires multiple bucket counters and a mutex for concurrent safety. A histogram with 10 buckets costs ~200ns per observe plus ~200 bytes allocation for the buckets.

For the critical request path, use atomic counters and gauges. Use histograms on a sampled subset of requests (1 in 100) for p99 latency tracking, or use Prometheus client's histograms which are optimized for concurrent access.

Exporting metrics (scraping) is separate from recording. The export path iterates all metrics and serializes them, which blocks the exporter but not the request handler.

## Practice task

Implement a `MetricsStore` that tracks HTTP request metrics:

1. `NewMetricsStore()` creates a new store.
2. `RecordRequest(durationMs int64, isError bool)` atomically records a request with its duration and error status.
3. `Snapshot()` returns the current (requests, errors, avgDurationMs) tuple.

The tests verify correct counting, error tracking, and average duration calculation across multiple operations, including edge cases like no operations (zero values).

## Tests / verification

```bash
go run ./curriculum/modules/12-observability-diagnostics/lessons/06-metrics
go test ./curriculum/modules/12-observability-diagnostics/lessons/06-metrics
```

## Review questions

1. What are the four Golden Signals of production monitoring? Give a Go metric name for each.
2. What is the difference between a counter and a gauge? Give an example use case for each in a Go HTTP service.
3. Why is `atomic.Int64` preferred over `sync.Mutex` for metrics counters in high-throughput Go services?
4. What does `expvar.Int.String()` return? How is this used by the `/debug/vars` endpoint?
5. A service has 5 endpoints. Without labels, what is missing from a single `http_requests_total` counter? How would you fix this with expvar?

## NEXT UP

Prometheus — a production-grade metrics system with histograms, summaries, labels, and a standardized scrape protocol.
