# Opslane observability

## Mission

Integrate enterprise-grade observability across Opslane, combining Prometheus dimensional telemetry metrics, sub-second latency histogram distribution tracking, contextual structured logging with `log/slog`, and distributed request correlation via `X-Correlation-ID`.

## Prerequisites

- `opslane-11`

## Mental Model

Observability enables operators to answer questions about the internal state of a production system without deploying new code. The telemetry stack connects three complementary pillars:
1. **Metrics**: Real-time dimensional counters (`opslane_orders_processed_total`) and latency histograms (`opslane_order_processing_duration_seconds`) scraped by Prometheus for alerts and Grafana dashboards.
2. **Logs**: Context-rich JSON records emitted via `log/slog` carrying the active `correlation_id` and `tenant_id` for granular root-cause investigation.
3. **Correlation Tracking**: Passing `X-Correlation-ID` across HTTP middleware and context trees, ensuring that every log line and metric sample can be correlated to a single client transaction.

## Visual Model

```text
HTTP Client (X-Correlation-ID: "corr-101")
  │
  ▼
┌─────────────────────────────────────────────────────────────┐
│ correlationMiddleware                                       │
│   Injects correlation_id into context.Context               │
│   Echoes X-Correlation-ID into response header              │
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ handleProcessOrder Execution                                │
│   ├─► timer := time.Now()                                   │
│   ├─► Prometheus: OrdersProcessedTotal.Inc()                │
│   ├─► Prometheus: OrderProcessingDuration.Observe(duration) │
│   └─► slog.Info("...", "correlation_id", "tenant_id", ...)  │
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ Prometheus Scraping: GET /metrics                           │
│   # HELP opslane_orders_processed_total ...                 │
│   opslane_orders_processed_total{status="success"} 42       │
└─────────────────────────────────────────────────────────────┘
```

## Machine View

Under the hood:
1. `NewMetricsRegistry` instantiates an isolated `prometheus.NewRegistry()` rather than relying on the default global registry, avoiding metric naming collisions during parallel testing.
2. `orderDuration` defines exponential sub-second buckets (`0.005s` to `2.5s`), tuned specifically for microservice SLA boundaries (p50, p90, p99).
3. The `GET /metrics` route mounts `promhttp.HandlerFor(s.metrics.Registry, promhttp.HandlerOpts{})`, outputting standard text-based Prometheus exposition format.
4. Structured logging emits JSON lines with key-value pairs (`"correlation_id"`, `"tenant_id"`, `"duration_sec"`), allowing log aggregators to index events in real time.

## Run Instructions

```bash
# Run observability test suite
go test -v ./curriculum/modules/18-flagship-opslane/assessments/12-opslane-observability

# Launch server and scrape metrics
go run ./curriculum/modules/18-flagship-opslane/assessments/12-opslane-observability
```

## Code Walkthrough

1. **Correlation Context (lines 17–32)**: Implements `CorrelationIDFromContext` and `ContextWithCorrelationID` using unexported context key.
2. **`MetricsRegistry` (lines 35–79)**: Configures isolated registry, counter vector with `tenant_id` and `status`, histogram vector with tuned duration buckets, and active workers gauge.
3. **`ObservabilityServer` & Router (lines 82–104)**: Sets up `GET /metrics` endpoint and instrumented business route wrapped in correlation middleware.
4. **`correlationMiddleware` (lines 106–119)**: Extracts or generates 8-byte hexadecimal correlation ID; updates context and response headers.
5. **`handleProcessOrder` (lines 121–147)**: Measures execution duration; increments Prometheus counters; records latency histogram; emits structured `slog.Info`.

## Try It

1. Run the test suite: verify that processing an order increments the counter and records latency in the histogram.
2. Query `GET /metrics`: verify that `opslane_orders_processed_total` and histogram buckets appear in Prometheus exposition text.
3. Send an order with custom `X-Correlation-ID: my-trace-123`: verify that response header and log output contain the exact correlation ID.

Common mistakes:
- Adding unbounded high-cardinality labels (like order IDs or user emails) to Prometheus metrics, causing Prometheus OOM crashes.
- Emitting unstructured log lines with `fmt.Sprintf` instead of structured attributes in `slog`.
- Using the global default Prometheus registry in packages, causing panics on re-registration in unit tests.

## In Production

In production Kubernetes clusters, Prometheus scrapers poll `/metrics` every 15–30 seconds. Alerts are configured on error rates (`rate(opslane_orders_processed_total{status="error"}[5m])`) and p99 latency (`histogram_quantile(0.99, sum(rate(opslane_order_processing_duration_seconds_bucket[5m])) by (le))`).

## Thinking Questions

1. Why should `order_id` NEVER be used as a Prometheus label, while `tenant_id` or `status` is acceptable?
2. What is the mathematical difference between an average latency calculation versus p99 histogram quantiles?
3. How does `X-Correlation-ID` differ from distributed trace spans in OpenTelemetry?

## Next Step

Proceed to [Milestone 13 — Graceful Shutdown](../13-opslane-graceful-shutdown/README.md) to implement signal handling, drain phases, and timeout-bounded termination.
