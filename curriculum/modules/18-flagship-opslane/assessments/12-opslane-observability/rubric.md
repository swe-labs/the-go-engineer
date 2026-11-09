# Opslane Milestone 12 Assessment — Rubric

## Assessment Dimensions

### 1. Prometheus Telemetry & Metric Modeling (30%)
- **Excellent (5 pts)**: Correct metric types chosen (Counter for totals, Histogram for latency, Gauge for workers); low-cardinality labels strictly enforced; sub-second histogram buckets configured; scrape route exposes standard exposition format.
- **Proficient (3 pts)**: Metrics registered and exposed, but histogram buckets uncalibrated or high-cardinality labels risked.
- **Developing (1 pt)**: Uses global registry causing collisions, or uses Gauge for cumulative counters.
- **Unacceptable (0 pts)**: Metrics do not increment or scrape route fails.

### 2. Contextual Structured Logging (25%)
- **Excellent (5 pts)**: `log/slog` emits JSON key-value pairs; correlation ID and tenant ID consistently attached; log levels used semantically (`Info` for success, `Error` for failures).
- **Proficient (3 pts)**: `slog` used, but attributes formatted as concatenated strings instead of structured attributes.
- **Developing (1 pt)**: Reverts to unstructured `log.Print` in request paths.
- **Unacceptable (0 pts)**: No logging implemented.

### 3. Correlation Tracking & Middleware Pipeline (25%)
- **Excellent (5 pts)**: `correlationMiddleware` checks incoming `X-Correlation-ID` header; generates random fallback; injects into context; echoes into response headers; safe context retrieval.
- **Proficient (3 pts)**: Middleware sets context, but fails to echo `X-Correlation-ID` into response headers.
- **Developing (1 pt)**: Middleware generates predictable sequence numbers instead of random entropy.
- **Unacceptable (0 pts)**: Correlation ID not propagated.

### 4. Conceptual & SRE Diagnostics Defense (20%)
- **Excellent (5 pts)**: Accurately explains cardinality explosion risks, p99 vs mean latency pitfalls, Prometheus metric type semantics, and distributed correlation tracking.
- **Proficient (3 pts)**: Answers conceptual questions accurately with minor gaps in histogram quantile calculations.
- **Developing (1 pt)**: Cannot explain why average latency hides tail latency spikes.
- **Unacceptable (0 pts)**: No explanation provided.

## Grading Thresholds
- **12–14 pts**: Exceeds expectations (Gold standard)
- **10–11 pts**: Meets expectations (Passing)
- **7–9 pts**: Needs revision
- **0–6 pts**: Unsatisfactory — review prerequisites and retake
