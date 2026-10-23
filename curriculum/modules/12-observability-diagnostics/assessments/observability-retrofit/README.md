# Project Assessment — Observability Retrofit

## Purpose

Validate mastery for the target curriculum items.

## Targets

- `project-observability-retrofit`

## Evidence

Submit answers, code evidence, test output, and explanation where required.

---

## 1. Assessment Overview & Industrial Context

In production distributed systems, deploying microservices as opaque black boxes is unacceptable. When an e-commerce checkout flow spikes in error rates or degrades in tail latency (p99), operations teams require real-time correlated telemetry across the three foundational pillars of observability:
1. **Structured Contextual Logging**: Emitting JSON-formatted log lines with consistent key/value attributes (`order_id`, `customer`, `amount_cents`) and context trace identifiers via Go 1.21+'s standard `log/slog`.
2. **Dimensional Telemetry Metrics**: Exposing counters (`order_requests_total`), histograms (`order_duration_seconds`), and gauges (`active_orders_count`) via Prometheus `client_golang` with low-cardinality labels for time-series aggregation and alerting.
3. **Distributed Tracing**: Creating hierarchical spans (`OrderService.CreateOrder`, `OrderService.GetOrder`, `OrderService.CancelOrder`) using OpenTelemetry (`go.opentelemetry.io/otel`), capturing execution context, recording runtime errors, and propagating trace context across boundaries.

Your challenge in this assessment is to retrofit production telemetry into an existing uninstrumented Go HTTP microservice without altering its public HTTP contracts or breaking core domain invariants.

---

## 2. Core Competencies Tested

- **Zero-Allocation Context Logging**: Replacing `log.Println` with `slog.InfoContext`, `slog.WarnContext`, and `slog.ErrorContext` carrying explicit typed attributes.
- **Dimensional Metric Modeling**: Correctly choosing between Counter, Gauge, and Histogram; avoiding high-cardinality label explosions (e.g., never using `order_id` as a Prometheus label).
- **Latency Distribution Tracking**: Configuring histogram latency buckets suitable for microservice response times (e.g., exponential or defined sub-second boundaries).
- **OpenTelemetry Span Lifecycle**: Instantiating spans via `tracer.Start(ctx, name)`, recording errors with `span.RecordError(err)`, setting semantic status codes (`codes.Error`), and guaranteeing `defer span.End()`.
- **Telemetry Exposure**: Registering and serving the standard Prometheus `/metrics` HTTP scrape endpoint alongside business routes.

---

## 3. Zero-Magic Architectural Breakdown

```
[ Inbound HTTP Request ]
          │
          ▼
┌─────────────────────────────────────────────────────────────────┐
│ HTTP Handler Layer (handler.go)                                │
│   1. Extracts or initializes OpenTelemetry trace context       │
│   2. Invokes OrderService method with contextual propagation   │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│ Domain Service Layer (OrderService)                             │
│   1. tracer.Start(ctx, "OrderService.<Op>") -> defer span.End() │
│   2. timer := prometheus.NewTimer(durationHist.WithLabelValues) │
│      defer timer.ObserveDuration()                              │
│   3. slog.InfoContext(ctx, "event", "order_id", id, ...)        │
│   4. Mutates state & increments active_orders_count Gauge       │
│   5. On error: span.RecordError(err), counter status="error"    │
└─────────────────────────────────────────────────────────────────┘
```

### Why Unstructured Logging Fails at Scale
`log.Printf("order %s failed: %v", id, err)` produces raw strings that require expensive regex indexing at log ingestion time. In contrast, `slog` structures logs into key-value pairs stored in binary or JSON formats, allowing sub-second indexing across billions of log lines in systems like Elasticsearch, Loki, or Datadog.

### Metric Cardinality Rules
In Prometheus, each unique combination of key-value label pairs generates a distinct time-series stored in memory. Adding high-cardinality data (such as `user_id`, `order_id`, or `uuid`) to label sets causes unbounded memory growth, degrading Prometheus scrapers. Labels must always be bounded enumerations (e.g., `operation="create"`, `status="success"`).

---

## 4. Implementation Requirements & Task Breakdown

### Task 1: Structured Logging with `log/slog`
- Replace all legacy `log.Print*` calls with `slog` methods.
- For business operations (`CreateOrder`, `GetOrder`, `CancelOrder`), log using `slog.InfoContext` or `slog.WarnContext`.
- Ensure each log record includes contextual keys:
  - `"order_id"`
  - `"customer"`
  - `"amount_cents"`

### Task 2: Prometheus Telemetry Instrumentation
- Implement `NewMetricsRegistry(reg prometheus.Registerer) *MetricsRegistry`:
  - `order_requests_total`: `*prometheus.CounterVec` with labels `[]string{"operation", "status"}`.
  - `order_duration_seconds`: `*prometheus.HistogramVec` with label `[]string{"operation"}`.
  - `active_orders_count`: `prometheus.Gauge`.
  - Register all metrics with `reg` (or `prometheus.DefaultRegisterer` if `nil`).
- In `CreateOrder`:
  - Start duration timer using `prometheus.NewTimer(s.metrics.orderDuration.WithLabelValues("create"))`.
  - On validation failure: increment counter with `operation="create", status="error"`.
  - On success: increment counter with `operation="create", status="success"`, and increment `activeOrders.Inc()`.
- In `GetOrder` and `CancelOrder`:
  - Observe execution duration.
  - Update status counters (`success` vs `error`).
  - When `CancelOrder` succeeds, decrement `activeOrders.Dec()`.
- Expose the HTTP metrics route:
  - `GET /metrics` routed to `promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{})`.

### Task 3: OpenTelemetry Distributed Tracing
- Obtain a named tracer via `otel.GetTracerProvider().Tracer("order-service")`.
- In `CreateOrder`, `GetOrder`, and `CancelOrder`:
  - Start span named `OrderService.<Operation>`.
  - Guarantee `defer span.End()`.
  - Set semantic attributes: `attribute.String("order.id", ...)`, `attribute.String("order.customer", ...)`, `attribute.Int("order.amount_cents", ...)`.
  - If operation fails, invoke `span.RecordError(err)` and set status code to `codes.Error`.

---

## 5. Verification & Acceptance Criteria

Run all tests in the project directory:

```bash
go test -v -race ./curriculum/modules/12-observability-diagnostics/projects/observability-retrofit/_solution
```

### Acceptance Checklist
- [ ] `TestOrderService_MetricsAndTracing` passes without race conditions.
- [ ] `TestOrderService_InvalidOrder_ErrorMetrics` confirms `status="error"` increments on bad input.
- [ ] `TestOrderService_Cancellation_Gauge` verifies `active_orders_count` increments on creation and decrements on cancellation.
- [ ] `GET /metrics` returns HTTP 200 with valid Prometheus exposition format containing registered metric families.
- [ ] OpenTelemetry in-memory span exporter captures all expected spans with matching semantic attributes.

---

## 6. Common Pitfalls & Anti-Patterns

1. **Leaking Unbounded Labels to Prometheus**:
   *Anti-pattern*: Adding `order_id` as a label to `order_requests_total`.
   *Correction*: Keep `order_id` strictly in `slog` contextual logs and OpenTelemetry span attributes; keep Prometheus labels restricted to finite enumerations (`operation`, `status`).

2. **Forgetting `defer span.End()`**:
   *Anti-pattern*: Starting an OpenTelemetry span without ending it on early return paths, causing span leaks and unexported trace telemetry.
   *Correction*: Call `defer span.End()` immediately after `tracer.Start(ctx, ...)`.

3. **Ignoring Context Propagation**:
   *Anti-pattern*: Passing `context.Background()` into downstream services rather than the request's `r.Context()`.
   *Correction*: Always propagate `r.Context()` to retain active trace IDs and cancelation signals.

4. **Re-registering Existing Metrics on Handler Reload**:
   *Anti-pattern*: Calling `prometheus.MustRegister` inside request handlers or constructors executed multiple times.
   *Correction*: Encapsulate metric registration in a singleton or factory (`NewMetricsRegistry`) attached to the application dependency graph.

---

## 7. 4-Tier Assessment Rubric

| Level | Criteria & Observable Evidence |
|---|---|
| **Level 1: Novice** | Code fails to compile or imports wrong OpenTelemetry packages. Unstructured `log.Println` remains in critical paths. Metrics are registered globally inside request handlers causing panic on re-registration. |
| **Level 2: Advanced Beginner** | Logging switched to `slog` but attributes are formatted as raw strings rather than key/value pairs. Metrics counter increments on success but fails to record errors. Spans are opened but missing `defer span.End()`. |
| **Level 3: Competent** | `slog`, Prometheus metrics, and OpenTelemetry spans are correctly implemented across all service methods. Prometheus counter distinguishes `success` from `error`. Active orders gauge decrements upon cancellation. All tests pass. |
| **Level 4: Proficient** | Flawless telemetry implementation. Histograms define tuned latency boundaries. Error spans record error details and update span status. Context propagation is maintained across all layers. Zero allocations in logging fast paths. |

---

## 8. Submission & Reflection

1. Provide output from `go test -v -race ./...`.
2. Provide a curl output snippet from `GET /metrics` showing `order_requests_total` and `active_orders_count`.
3. Answer: *Why is distributed tracing necessary in a microservice mesh even if every service has comprehensive Prometheus metrics and structured logging?*
