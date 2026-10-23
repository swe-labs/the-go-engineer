# Opslane Milestone 12 Assessment — Observability: Questions

Total: 10 questions (6 multiple choice × 1 pt, 4 short answer × 2 pt = 14 points). Passing score: 10/14 (71%).

---

## Multiple Choice (1 point each)

**1.** Why must high-cardinality fields like `order_id` or user email NEVER be used as Prometheus label values?
- a) Prometheus cannot serialize strings with hyphens.
- b) Each unique combination of label key-value pairs creates a new distinct time-series in memory; unbounded cardinality causes exponential memory consumption and crashes the Prometheus scraper.
- c) Label strings are limited to 4 characters.
- d) Prometheus only supports integer labels.

**2.** What is the purpose of histogram buckets in `opslane_order_processing_duration_seconds`?
- a) To round order amounts to the nearest dollar.
- b) To partition observed execution times into duration ranges, enabling Prometheus to calculate accurate percentile latencies (p50, p90, p99) across time windows.
- c) To throttle slow orders.
- d) To limit CPU usage during logging.

**3.** In `ObservabilityServer`, why is `prometheus.NewRegistry()` created locally instead of calling `prometheus.MustRegister()` on the default global registry?
- a) Local registries avoid metric collision panics (`duplicate metrics collector registration attempted`) when tests run in parallel or instantiate multiple server instances.
- b) The global registry cannot serialize JSON.
- c) Default registries do not support histograms.
- d) Go standard library forbids global registries.

**4.** What does `correlationMiddleware` do when an incoming HTTP request lacks an `X-Correlation-ID` header?
- a) Drops the connection immediately.
- b) Generates an 8-byte cryptographically random hex string, injects it into the request context, and writes it to the response header.
- c) Returns an HTTP 400 Bad Request.
- d) Defaults the correlation ID to `"unknown"`.

**5.** Which metric type is most appropriate for tracking the current number of background worker goroutines actively executing?
- a) Counter.
- b) Gauge.
- c) Histogram.
- d) Summary.

**6.** Where are Prometheus metrics exposed for scraping in this milestone?
- a) `GET /api/v1/metrics`.
- b) `GET /metrics`.
- c) `POST /debug/metrics`.
- d) On a raw TCP socket port 9100.

---

## Short Answer (2 points each)

**7.** Explain the difference between Counter, Gauge, and Histogram metric types in Prometheus. Give an example use case for each in an e-commerce platform.

**8.** What is Cardinality in time-series databases? How does Prometheus calculate the number of active time-series for a metric?

**9.** Why is average latency (mean) a misleading metric for evaluating user experience in microservices? Why do Site Reliability Engineers (SREs) prioritize 95th and 99th percentiles (p95, p99)?

**10.** Explain how structured contextual logging with `log/slog` and `X-Correlation-ID` allows an engineer to trace a single request's lifecycle across distributed log streams.
