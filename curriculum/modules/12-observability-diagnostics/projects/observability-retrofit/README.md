# Observability Retrofit

## Learning objective

Retrofit production-grade observability into an existing, uninstrumented Go HTTP microservice without altering its business logic or public HTTP API contracts. Master structured contextual logging with `log/slog`, dimensional telemetry metrics with Prometheus (`client_golang`), distributed tracing with OpenTelemetry (`go.opentelemetry.io/otel`), and standard health/metrics scrape endpoints.

## Why this matters

In modern distributed microservice environments, applications cannot be operated as opaque black boxes. When customer transactions fail or experience degraded latencies, operators need correlation across logs, metrics, and traces (the three pillars of observability) to locate root causes within seconds. Retrofitting telemetry onto legacy services without destabilizing core business contracts is one of the most common and valuable real-world engineering responsibilities.

## Prerequisites

Complete all Module 12 lessons (01–21) on structured logging with `log/slog`, Prometheus metrics instrumentation (Counters, Histograms, Gauges), OpenTelemetry distributed tracing spans and context propagation, and graceful shutdown.

## Available concepts

This project draws on `log/slog`, `context.Context` propagation, `github.com/prometheus/client_golang/prometheus`, `promhttp.Handler`, `go.opentelemetry.io/otel`, `go.opentelemetry.io/otel/trace`, OpenTelemetry attribute conventions, and idiomatic Go HTTP servers.

## Project overview

The `_starter/` directory provides an existing order management service with in-memory persistence and HTTP endpoints (`POST /orders`, `GET /orders/{id}`, `POST /orders/{id}/cancel`). However, the existing code uses unstructured `log.Println`, has zero metrics, and does not produce OpenTelemetry trace spans.

Your mission is to retrofit comprehensive observability:
- **Structured Contextual Logging**: Replace all unstructured `log.Println` calls with `log/slog` methods (`InfoContext`, `WarnContext`, `ErrorContext`) carrying contextual attributes (`order_id`, `customer`, `amount_cents`).
- **Prometheus Telemetry**:
  - `order_requests_total`: CounterVec tracking operation counts partitioned by `operation` (`create`, `get`, `cancel`) and `status` (`success`, `error`).
  - `order_duration_seconds`: HistogramVec tracking latency distribution across operations.
  - `active_orders_count`: Gauge tracking the current number of non-cancelled orders.
  - Expose `/metrics` scrape endpoint.
- **OpenTelemetry Distributed Tracing**:
  - Instrument `CreateOrder`, `GetOrder`, and `CancelOrder` with trace spans (`OrderService.<Operation>`).
  - Set semantic span attributes (`order.id`, `order.customer`, `order.amount_cents`).
  - Record error events on failures using `span.RecordError`.

## Architecture & Metric Contracts

```
[ HTTP Client ]
      │
      ▼
[ HTTP Handler (POST /orders, GET /orders/{id}, /metrics) ]
      │
      ├─► Context Propagation & Trace Span (otel)
      ├─► Latency Observation (Prometheus Histogram)
      ├─► Operation Count & Status Increment (Prometheus Counter)
      │
      ▼
[ OrderService Core Logic ]
      │
      ├─► Structured slog (JSON / Contextual Attrs)
      └─► State Tracking (Prometheus Gauge active_orders_count)
```

## Tasks

### Task 1 — Structured Logging with `log/slog`

In `_starter/main.go`:
- Initialize a JSON or Text handler configured with `slog.SetDefault`.
- Replace unstructured `log.Println` calls with `slog.InfoContext` and `slog.WarnContext`.
- Ensure log records include explicit attributes: `order_id`, `customer`, and `amount_cents`.

### Task 2 — Prometheus Metrics & Scrape Endpoint

- In `NewMetricsRegistry(reg prometheus.Registerer)`:
  - Define `order_requests_total` (`CounterVec` with labels `["operation", "status"]`).
  - Define `order_duration_seconds` (`HistogramVec` with labels `["operation"]`).
  - Define `active_orders_count` (`Gauge`).
  - Register metrics against the provided `reg` (or `prometheus.DefaultRegisterer` if `nil`).
- In `CreateOrder`:
  - Observe execution duration with `prometheus.NewTimer`.
  - Increment counter on success (`status="success"`) and on invalid payload (`status="error"`).
  - Increment `active_orders_count` on successful creation.
- In `GetOrder` & `CancelOrder`:
  - Observe execution duration.
  - Increment success and error counters accordingly.
  - Decrement `active_orders_count` when an order transition to `CANCELLED`.
- In `NewHTTPHandler`:
  - Expose `GET /metrics` using `promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{})`.

### Task 3 — OpenTelemetry Distributed Tracing

- Acquire a tracer via `otel.GetTracerProvider().Tracer("order-service")`.
- In `CreateOrder`, start span `"OrderService.CreateOrder"`. Defer `span.End()`.
- Record error with `span.RecordError(err)` when validation fails.
- Attach semantic attributes to the span: `order.id`, `order.customer`, `order.amount_cents`.
- Instrument `GetOrder` and `CancelOrder` with corresponding span names and `order.id` attributes.

## Verification

Run the automated test suite:

```bash
go test -v ./curriculum/modules/12-observability-diagnostics/projects/observability-retrofit/_solution
```

Expected output:
```
=== RUN   TestOrderService_MetricsAndTracing
--- PASS: TestOrderService_MetricsAndTracing (0.00s)
=== RUN   TestOrderService_InvalidOrder_ErrorMetrics
--- PASS: TestOrderService_InvalidOrder_ErrorMetrics (0.00s)
=== RUN   TestHTTPHandler_Endpoints
--- PASS: TestHTTPHandler_Endpoints (0.00s)
PASS
ok  	github.com/rasel9t6/the-go-engineer/curriculum/modules/12-observability-diagnostics/projects/observability-retrofit/_solution	0.443s
```

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|---|---|---|---|
| **Prometheus Metrics** | Counter, histogram, and gauge accurately track operations, latencies, and active counts without leaks | Metrics recorded but missing label dimensions or gauge decrement | Metrics missing or cause registration panics |
| **OpenTelemetry Tracing** | Spans created for all operations, errors recorded, and semantic attributes attached | Spans created but lack attributes or error recording | No tracing instrumentation |
| **Structured Logging** | `slog` with context-aware methods and structured attributes; zero legacy log calls | Replaced some logs but left unstructured strings or missing fields | Uses unstructured standard library `log` |
| **API Contract Preservation** | All existing HTTP endpoints function identically without regressions | Endpoints work but response codes or payload formats modified | Service broken or fails to start |

## NEXT UP

[Profile and Fix Lab](../../13-performance-memory-engineering/projects/profile-and-fix-lab/README.md)
