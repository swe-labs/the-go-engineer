# Why observability exists

## Learning objective

Explain why observability is a prerequisite for operating production Go services, distinguish observability from monitoring, and identify the three pillars (logs, metrics, traces) with their specific use cases.

## Why this matters

A production Go service fails in ways you cannot predict. A goroutine leaks under high concurrency, a database query slows after a schema change, a dependency returns corrupt data only for certain users. Without observability, every incident starts with "reproduce it locally" — and many production bugs cannot be reproduced locally because they depend on traffic patterns, data volume, or timing. Observability means you can answer any question about your system's internal state by examining data already being emitted, without shipping new code. Teams that invest in observability resolve incidents in minutes instead of hours.

## Mental model

Observability is the property that lets you understand a system's internal state from its external outputs. Monitoring tells you "something is wrong" (CPU > 90%, error rate > 5%). Observability tells you "what is wrong and where" (the /api/orders endpoint returns 500 for requests with user_id containing special characters because the validation regex does not escape the backslash).

The three pillars each answer a different question:

- **Logs**: "What happened?" — discrete events with a timestamp and structured data.
- **Metrics**: "How many times did it happen?" — numeric aggregations over time.
- **Traces**: "Where exactly did it happen?" — the end-to-end path of a single request.

When all three are present and correlated (via a correlation ID), an operator can start with a metric alert (error rate spike), drill into logs (all errors have status=500, handler=/api/payments), and open a trace (the span for processPayment takes 5s because the database is under memory pressure).

## Core idea

Observability is designed in, not retrofitted. A system is observable when you can determine its internal state from the data it emits — without guessing, without reproducing, without redeploying with new log lines.

The opposite of observability is "blind debugging": adding a log line, redeploying, waiting for the error to recur, and repeating. This cycle takes 20–60 minutes per iteration and fails entirely for intermittent issues.

## Under the hood

At the machine level, every observability signal occupies CPU, memory, and I/O:

- **Log I/O**: Each structured log line serializes key-value pairs into a buffer, acquires a write lock on the output stream (or enqueues to a buffer), and flushes to disk or network. The cost is ~1–10μs per line depending on the number of attributes.
- **Metric aggregation**: A counter increment is a single atomic CPU instruction (LOCK XADD on x86). A histogram bucket update requires a mutex lock on the bucket slice. Prometheus client_golang uses a sync.Mutex for histograms and atomic.AddUint64 for counters.
- **Trace spans**: Each span allocates a struct with a span ID, parent ID, timestamps, and attributes. Span lifecycle: start → add attributes → end → enqueue for export. The allocation and queuing cost is ~100–500ns per span.

Every observability library adds a code path to the hot path of your handler. The key is keeping the cost below 1% of request latency so it is negligible at the p99 level.

## How Go uses it

Go's standard library provides the building blocks for all three pillars:

- **Logs**: `log/slog` (Go 1.21+) for structured logging with levels, handlers (text, JSON), and context-scoped loggers.
- **Metrics**: `expvar` for exporting runtime and custom metrics as JSON over HTTP at `/debug/vars`. The `runtime` package exposes MemStats, NumGoroutine, NumCPU, and GC stats.
- **Traces**: No built-in tracer, but Go uses the OpenTelemetry SDK for distributed tracing. The `net/http` package wraps HTTP handlers and transports for span creation and propagation.

Production services layer on top: `prometheus/client_golang` for metrics, `otel/otel` for traces, `log/slog` for structured logs. The standard library's net/http supports middleware chaining, making it straightforward to add observability to every handler.

## Go example

```go
package main

import (
	"fmt"
	"log/slog"
	"time"
)

// BlindLog prints request info without structured context.
func BlindLog(path string, status int, dur time.Duration) {
	fmt.Printf("request handled: %s %d %v\n", path, status, dur)
}

// ObservableLog logs a request with structured context.
func ObservableLog(path string, status int, dur time.Duration) {
	slog.Info("request handled",
		"path", path,
		"status_code", status,
		"duration_ms", dur.Milliseconds(),
	)
}

func main() {
	BlindLog("/api/orders", 500, 5200*time.Millisecond)
	ObservableLog("/api/orders", 500, 5200*time.Millisecond)
}
```

Run with `go run .` to see the difference. The unstructured output `request handled: /api/orders 500 5.2s` cannot be filtered, aggregated, or alerted on in a log aggregator. The structured output `level=INFO msg="request handled" path=/api/orders status_code=500 duration_ms=5200` can be parsed by any log aggregator: search for `status_code=500`, sum `duration_ms` by path, or alert when `duration_ms > 1000`.

## Step-by-step execution

For `ObservableLog("/api/orders", 500, 5200*time.Millisecond)` with a default text handler:

1. `slog.Info` creates a `slog.Record` with current time, `LevelInfo`, and the message `"request handled"`.
2. The variadic args are grouped into attribute pairs: `"path"` → `"/api/orders"`, `"status_code"` → `500`, `"duration_ms"` → `5200`.
3. The handler's `Handle` method receives the record and serializes it: `time=2026-01-15T10:30:00Z level=INFO msg="request handled" path=/api/orders status_code=500 duration_ms=5200`.
4. The handler writes the serialized bytes to stderr (default output).
5. The record is discarded. No memory is retained after the call returns.

## Common mistakes

- **Observability as an afterthought**: Adding logging, metrics, and tracing after the first production incident. The team is blind during the incident and must deploy an observability fix to understand the outage. By the time the data arrives, the incident may have self-resolved with no root cause identified.

- **Logging everything at INFO level**: A service with 1000 requests/second logging 10 fields per request produces 10,000 log lines/second. At this volume, operators cannot distinguish a routine request from an error without complex queries. The signal is drowned in noise.

- **Using fmt.Println for production diagnostics**: Unstructured text cannot be queried by a log aggregator (Datadog, Grafana Loki, CloudWatch). To find all errors for a specific endpoint, operators must grep raw files — which is impossible in ephemeral container environments.

- **Monitoring only infrastructure metrics**: CPU, memory, and disk usage look healthy while the application is deadlocked. Goroutine leaks, channel deadlocks, and request queue backpressure are invisible to infrastructure metrics. Application-level observability (request rate, error rate, latency percentile) catches these.

- **Confusing monitoring with observability**: Monitoring tells you the error rate is 5%. Observability tells you the error rate is 5% on /api/payments because the SSN validation returns `false` for users with hyphens in their SSN — and which specific users are affected.

## Debugging walkthrough

Consider a Go HTTP handler with no observability:

```go
func handler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	err := processUser(userID)
	if err != nil {
		fmt.Printf("error processing user: %v\n", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}
```

**Symptom**: Production monitoring shows a 5% error rate on the endpoint, but the logs contain only `error processing user: something went wrong` with no detail about which user_id failed or what the actual error was.

**Investigation**: Add structured context to the error path:

```go
slog.Error("process user failed",
	"user_id", userID,
	"error", err,
	"path", r.URL.Path,
)
```

With this change, the log aggregator shows that errors correlate with `user_id` values containing the string `null` — a downstream bug in data ingestion. Without structured context, this correlation requires manually cross-referencing application logs with database query logs, taking 30+ minutes.

**Root cause**: The team treated logging as a low priority and used unstructured print statements. A single structured log line with three key-value pairs would have reduced debugging time from 30 minutes to 30 seconds.

**Fix**: Always log structured context (user_id, error type, request path) for every error. Use `slog.Error` instead of `fmt.Printf`.

## Production notes

Every production Go service at scale depends on observability. Kubernetes uses structured logging throughout, exposing metrics via a `/metrics` endpoint scraped by Prometheus. Docker Hub uses correlation IDs across its microservices. Uber's observability pipeline processes billions of spans per day.

Without observability, a 10-query database slowdown is invisible until users report errors — with observability, the p99 latency chart shows the degradation in real time. SRE teams define SLOs (service level objectives) based on metrics: latency < 200ms for 99.9% of requests over 30 days. These SLOs are impossible to measure without a metrics pipeline.

Production observability also feeds into incident response runbooks: "If error rate > 5% for /api/payments, check database connection pool saturation. If latency > 1000ms for /api/search, check Elasticsearch cluster health." These runbooks are only possible when the relevant signals are already instrumented.

## Performance implications

Observability always adds CPU and memory overhead. The engineering challenge is making that overhead negligible:

- **Single log line**: ~1μs extra latency per handler call with 5–10 attributes. For a handler that takes 50ms, this is 0.002% overhead — invisible.
- **Metric counter increment**: ~50ns (atomic add). This is free for all practical purposes.
- **Expensive operations**: Stack traces (`runtime.Caller`), large attribute values (serializing a 10KB payload into a log attribute), and synchronous log writes to disk all add measurable latency. Use `slog.HandlerOptions` with `AddSource: false` in production, and batch log writes with an async handler.

The rule: never let observability add more than 1% to your p99 latency. If it does, move the expensive operation (log serialization, metric export, trace batching) off the request goroutine.

## Practice task

Complete the function `LogRequest` that writes a structured log entry with all request context:

```go
func LogRequest(logger *slog.Logger, handler, method, path string, statusCode int, latencyMs int64) {
	// Log a structured entry with all parameters as attributes
}
```

Use `logger.Info` and pass each parameter as a key-value pair. Run `go run .` to see the output, then `go test .` to verify your implementation against the table-driven tests.

## Tests / verification

```bash
go run ./curriculum/modules/12-observability-diagnostics/lessons/01-why-observability-exists
go test ./curriculum/modules/12-observability-diagnostics/lessons/01-why-observability-exists
```

## Review questions

1. What is the difference between observability and monitoring? Give a concrete example where a system is monitored but not observable.
2. Name the three pillars of observability and state what question each pillar answers.
3. Why is `fmt.Println` inadequate for production diagnostics in Go?
4. A production incident causes a 50ms latency spike on one endpoint but not others. Without observability, how would you investigate? With observability, what data would you look at first?
5. What does "observability must be designed in, not retrofitted" mean? Describe a scenario where retrofitting observability during an incident would fail.

## NEXT UP

Structured logging with slog — the first pillar of observability, with Go's log/slog package.
