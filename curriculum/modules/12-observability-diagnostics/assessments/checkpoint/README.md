# Module 12 Assessment — Observability, Diagnostics, and Reliability

## Purpose

Validate mastery for the target curriculum items.

## Targets

- `core-12-01`
- `core-12-02`
- `core-12-03`
- `core-12-04`
- `core-12-05`
- `core-12-06`
- `core-12-07`
- `core-12-08`
- `core-12-09`
- `core-12-10`
- `core-12-11`
- `core-12-12`

## Evidence

Submit answers, code evidence, test output, and explanation where required.

---

## 1. Theoretical & Architecture Scenarios

Answer each scenario in depth, citing Go runtime mechanics, standard library specifications (`log/slog`, `net/http/pprof`), and telemetry standards (Prometheus, OpenTelemetry).

### Scenario 1: Structured Logging vs. Legacy Formatted Strings (`core-12-01`, `core-12-02`)
A legacy Go service logs high-volume messages using `log.Printf("processed order %s for customer %s in %v", orderID, customer, duration)`. In an incident retrospective, engineering discovers that indexing and querying these logs in Elasticsearch takes minutes and frequently times out.
- Explain the mechanical differences between `fmt.Sprintf`/`log.Printf` and `slog.InfoContext(ctx, "processed order", "order_id", orderID, "customer", customer, "duration", duration)`.
- How does `slog.Attr` minimize heap allocation using value representations (`slog.Value`)?
- How does context propagation in `slog.InfoContext` enable distributed trace correlation across disparate services?

### Scenario 2: Metric Types & Cardinality Explosions (`core-12-03`, `core-12-04`)
An engineer instruments an HTTP API with Prometheus using:
```go
httpRequestsTotal = prometheus.NewCounterVec(
    prometheus.CounterOpts{Name: "http_requests_total"},
    []string{"method", "path", "user_id", "status"},
)
```
Within two days, the Prometheus scraper encounters Out-Of-Memory (OOM) errors and service latency increases noticeably.
- Define "cardinality" in time-series metric databases and explain why `user_id` caused this incident.
- When should an engineer use a `Counter`, a `Gauge`, or a `Histogram`?
- How should histogram buckets (`prometheus.LinearBuckets` vs `prometheus.ExponentialBuckets`) be chosen to accurately measure tail latency (p95, p99)?

### Scenario 3: OpenTelemetry Tracing & Context Propagation (`core-12-05`, `core-12-06`)
In a microservice call chain (Client -> Gateway -> Order Service -> Payment Service), trace spans generated in Payment Service appear as isolated traces detached from the Gateway request.
- Explain how W3C TraceContext headers (`traceparent`, `tracestate`) must be injected into HTTP outbound headers using OpenTelemetry propagators (`propagation.TraceContext`).
- Explain the role of `span.RecordError(err)` versus `span.SetStatus(codes.Error, msg)`. Does recording an error automatically mark the span as failed in APM dashboards?
- Why must `span.End()` be called in a `defer` block immediately following span initiation?

### Scenario 4: Health Probes, Liveness vs Readiness (`core-12-07`, `core-12-08`)
In a Kubernetes deployment, a service's database dependency experiences temporary 5-second connection pool exhaustion. The liveness probe immediately fails, causing Kubernetes to terminate the pod. While restarting, cascading failure takes down all other replicas.
- Differentiate between **Liveness** (`/healthz`), **Readiness** (`/readyz`), and **Startup** probes.
- Should a downstream database outage fail the liveness probe, the readiness probe, or neither? Justify your architectural design.
- How do you implement a graceful shutdown lifecycle in Go using `os/signal`, `syscall.SIGTERM`, and `server.Shutdown(ctx)`?

### Scenario 5: Diagnostic Profiling in Production (`core-12-09` through `core-12-12`)
A production Go service experiences high CPU utilization and latency spikes during peak hours. An operator wants to mount `net/http/pprof`.
- What are the security risks of exposing `net/http/pprof` on public API ports? How should it be isolated?
- Explain the sampling mechanism behind Go's CPU profiler (SIGPROF signals at 100 Hz). What is the observed CPU overhead?
- How does Go's memory/heap profiler track allocations without intercepting every pointer on the stack?

---

## 2. Practical Diagnostic Challenge

Write a self-contained Go demonstration or code snippet illustrating:
1. Setting up a root `slog.Logger` configured with a custom JSON handler that extracts OpenTelemetry `trace_id` and `span_id` from `context.Context` and appends them to all emitted log records.
2. An HTTP middleware that:
   - Starts an OpenTelemetry span for incoming requests.
   - Measures duration using a Prometheus histogram vector partitioned by `path` and `status_code`.
   - Increments an active request gauge at entry and decrements it at completion.
   - Logs an error with stack trace if status code >= 500.

---

## 3. Invariant Operational Rules

1. **Strict Cardinality Bounds**: Prometheus metric label sets must NEVER include unbounded identifiers (UUIDs, timestamps, user IDs, credit card numbers, email addresses).
2. **Context Continuity**: Never discard `context.Context` when invoking logging or tracing methods in HTTP or gRPC request paths.
3. **Fail Safe Scrapes**: `/metrics` and `/healthz` endpoints must never deadlock, block indefinitely, or allocate excessive memory.
4. **Graceful Draining**: During SIGTERM reception, stop accepting new readiness traffic before shutting down in-flight requests and background workers.

---

## 4. 4-Tier Assessment Rubric

| Level | Criteria & Observable Evidence |
|---|---|
| **Level 1: Novice** | Confuses metrics with logs; creates Prometheus metrics with UUID labels; passes `context.TODO()` or `context.Background()` ignoring request context; unaware of liveness vs readiness differences. |
| **Level 2: Advanced Beginner** | Uses `slog` but without contextual attributes; understands metric types but sets improper histogram buckets; uses OpenTelemetry but forgets context propagation across HTTP boundaries. |
| **Level 3: Competent** | Explains telemetry mechanics accurately; distinguishes counter, gauge, and histogram correctly; understands traceparent headers; designs isolated health endpoints; handles graceful shutdown with `server.Shutdown`. |
| **Level 4: Proficient** | Masters deep runtime diagnostics; explains pprof sampling rates and signal handling; designs zero-allocation custom `slog` handlers with trace correlation; articulates production incident mitigation strategies with clear root-cause telemetry. |

---

## 5. Submission & Verification

Submit a markdown document containing your complete technical analyses for Scenarios 1–5 and your tested code snippet for the practical challenge.
