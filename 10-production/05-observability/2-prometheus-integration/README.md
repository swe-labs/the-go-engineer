# OPS.2 Prometheus Integration

## Mission

Master the "Pull-Based" monitoring model. Learn how to use the **Prometheus client library** for Go to expose an HTTP endpoint (usually `/metrics`) that allows a Prometheus server to "scrape" your application's internal metrics. Understand how to register custom metrics and how they are transformed into **Time Series** data for visualization in Grafana.

## Prerequisites

- OPS.1 Metrics Basics
- Section 06: Backend & APIs (Basic HTTP server knowledge)

## Mental Model

Think of Prometheus Integration as **An Open Guest Book**.

1. **The Push Model (Traditional)**: Every time something happens, you pick up the phone and call the monitoring server. If you have 10,000 servers, the monitoring server gets overwhelmed by phone calls.
2. **The Pull Model (Prometheus)**: You keep a guest book on your front porch (The `/metrics` endpoint). You update it internally whenever you want.
3. **The Scrape**: Once a minute, the Prometheus server walks by, reads your guest book, and writes down the numbers in its own notebook.
4. **The Advantage**: Your application doesn't care if the monitoring server is slow or down; it just keeps writing in its own book.

## Visual Model

```mermaid
graph LR
    App[Go App] -->|Internal Update| Register[Metric Registry]
    Register -->|Serve HTTP| Endpoint[/metrics]
    Prom[Prometheus Server] -->|Scrape / Pull| Endpoint
    Prom -->|Query| Grafana[Grafana Dashboard]
```

## Machine View

- **`promhttp.Handler()`**: The standard HTTP handler provided by Prometheus that converts your Go metric objects into the Prometheus text exposition format.
- **`prometheus.MustRegister()`**: The global registry where you "Check in" your metrics so the handler knows about them.
- **Scrape Interval**: The frequency at which Prometheus pulls data (e.g., every 15s). This determines the resolution of your graphs.

## Run Instructions

```bash
# Start the service
go run ./10-production/05-observability/2-prometheus-integration

# In another terminal, view the raw metrics format:
# curl http://localhost:8080/metrics
```

## Code Walkthrough

### Defining the Registry
Shows how to create a custom registry instead of using the global one (recommended for testing).

### The Metrics Endpoint
Demonstrates attaching the `promhttp` handler to your HTTP mux.

### Practical Instrumenting
Shows how to wrap an existing HTTP handler with a middleware that automatically tracks request counts and latencies.

## Try It

1. Run the service and `curl` the `/metrics` endpoint. Can you find the `http_requests_total` metric?
2. Refresh the browser a few times and `curl` again. Does the number increase?
3. Add a `HELP` and `TYPE` comment to a custom metric using the `prometheus.Opts` struct.
4. Discuss: Why is the Prometheus format simple plain text instead of JSON?

## In Production
**Don't put the metrics endpoint on the public internet.** Anyone who can access `/metrics` can see internal details about your system's performance and usage. Always put your metrics on a **Private Port** (e.g., `9090`) or protect them with internal network rules (IP allow-listing) or basic authentication.

## Thinking Questions
1. What are the benefits of a Pull model over a Push model?
2. How does Prometheus handle a target that is temporarily offline?
3. Why should you avoid using the "Global" registry in a large, modular application?

## Next Step

Next: `OPS.3` -> `10-production/05-observability/3-distributed-tracing-basics`

Open `10-production/05-observability/3-distributed-tracing-basics/README.md` to continue.
