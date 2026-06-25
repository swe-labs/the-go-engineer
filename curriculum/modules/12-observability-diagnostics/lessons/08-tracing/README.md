# Tracing

## Learning objective

Implement manual distributed tracing in Go using context propagation, create parent-child span relationships, and extract trace context from HTTP headers.

## Why this matters

Metrics tell you that your p99 latency is 2 seconds. Logs tell you that a specific request failed. Only tracing tells you _which_ service call caused the 2-second latency: is it the database query, the external API call, or the payment processing step? In a monolith, you can add timers. In a distributed system with 10+ services, you need propagated trace context that follows a single request across every service boundary. Without tracing, debugging a slow request in a microservice architecture is like finding a needle in 10 haystacks.

## Mental model

A trace is a complete itinerary of a package delivery. The trace ID is the tracking number. Each stop (span) records: what happened (HTTP GET /orders), who handled it (Order Service), how long it took (8ms), and what was done inside (SELECT from orders table, then POST to Payment Service). The trace waterfall shows the full journey from the client's first request to the final response, with every nested call visible.

A span is like a single stop on the itinerary. It has a start time, end time, and optional child spans. Child spans represent sub-operations — the payment service call is a child of the order handler span. The critical insight: spans are connected via context propagation. The parent passes its span context to the child through HTTP headers, gRPC metadata, or message queue headers. The child creates a new span with the parent's trace ID and a new span ID, linking them in the trace tree.

The analogy breaks for sampling: unlike a physical delivery where every package is tracked, distributed tracing typically samples only 1-5% of requests to keep storage and CPU costs manageable.

## Core idea

A _trace_ represents the full path of a single request as it travels through a distributed system. A _span_ represents a single unit of work within a trace. Key span properties:

- **Trace ID**: a 16-byte (128-bit) identifier shared by all spans in the same trace.
- **Span ID**: an 8-byte (64-bit) identifier unique to this span.
- **Parent Span ID**: the span ID of the parent span, or empty for root spans.
- **Operation name**: a human-readable name like `POST /api/orders`.
- **Timestamps**: start and end time.
- **Attributes**: key-value pairs (e.g., `http.method=GET`, `db.statement=SELECT`).
- **Status**: OK or Error with optional description.
- **Events**: timestamped annotations within a span (e.g., "cache miss").

Trace context is propagated via the W3C TraceContext standard: the `traceparent` header contains `version-trace_id-parent_span_id-trace_flags` (e.g., `00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01`).

## Under the hood

A manual tracing library works through several layers:

1. **Tracer**: creates new spans. `tracer.Start(ctx, "operation")` extracts the parent span from the context, creates a new span with a new span ID, and stores the new span in the returned context.
2. **Span**: records the operation. `span.End()` computes duration and finalizes the span for export.
3. **Propagator**: injects span context into HTTP headers (for the caller) and extracts it from headers (for the callee).
4. **Exporter**: sends completed spans to a backend (Jaeger, Zipkin, Datadog).

Context propagation in Go uses `context.Context` as the carrier. The span is stored in the context using a private context key type to avoid collisions. When a service calls another service, it injects the trace context into the outgoing request headers. The receiving service extracts the trace context from the incoming headers and creates a child span.

The propagation format (W3C TraceContext) encodes the trace ID, span ID, and trace flags as a comma-separated hex string. The `traceparent` header format is:

```
00-<trace_id_32_hex>-<parent_span_id_16_hex>-<trace_flags_2_hex>
```

## How Go uses it

Before OpenTelemetry became the standard, Go tracing was fragmented: OpenTracing provided a vendor-neutral API, OpenCensus provided both API and implementation, and each vendor (Jaeger, Zipkin, Datadog) had its own SDK. OpenTelemetry merged both efforts in 2019 and is now the CNCF standard.

For manual tracing without external dependencies, Go's standard library provides everything needed: `context.Context` for propagation, `crypto/rand` for ID generation, and `net/http` for header injection/extraction. Production services use OpenTelemetry, but understanding manual tracing reveals exactly what OpenTelemetry does internally.

## Go example

```go
package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

type SpanContext struct {
	TraceID    string
	SpanID     string
	ParentSpanID string
}

type Span struct {
	Context    SpanContext
	Operation  string
	StartTime  time.Time
	EndTime    time.Time
	Attributes map[string]string
	Children   []*Span
}

func generateID(bytes int) string {
	b := make([]byte, bytes)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func NewTrace() SpanContext {
	return SpanContext{
		TraceID: generateID(16),
		SpanID:  generateID(8),
	}
}

func NewChildSpan(parent SpanContext) SpanContext {
	return SpanContext{
		TraceID:      parent.TraceID,
		SpanID:       generateID(8),
		ParentSpanID: parent.SpanID,
	}
}

type Tracer struct {
	root *Span
}

func (t *Tracer) StartSpan(ctx SpanContext, operation string) (*Span, *Tracer) {
	span := &Span{
		Context:    ctx,
		Operation:  operation,
		StartTime:  time.Now(),
		Attributes: make(map[string]string),
	}
	return span, &Tracer{root: span}
}

func (s *Span) End() {
	s.EndTime = time.Now()
}

func (s *Span) SetAttribute(key, value string) {
	s.Attributes[key] = value
}

func (s *Span) AddChild(child *Span) {
	s.Children = append(s.Children, child)
}

func simulateServiceCall(tracer *Tracer, traceCtx SpanContext, serviceName string, duration time.Duration) {
	childCtx := NewChildSpan(traceCtx)
	span, _ := tracer.StartSpan(childCtx, serviceName)
	span.SetAttribute("service", serviceName)
	time.Sleep(duration)
	span.End()
	fmt.Printf("[%s] Span %s: %s (%v)\n",
		traceCtx.TraceID[:8], span.Context.SpanID[:8], span.Operation, span.EndTime.Sub(span.StartTime))
}

func main() {
	traceCtx := NewTrace()
	tracer := &Tracer{}
	root, _ := tracer.StartSpan(traceCtx, "GET /api/orders")
	root.SetAttribute("http.method", "GET")
	root.SetAttribute("http.path", "/api/orders")

	simulateServiceCall(tracer, traceCtx, "auth_service", 10*time.Millisecond)
	simulateServiceCall(tracer, traceCtx, "order_service", 30*time.Millisecond)
	simulateServiceCall(tracer, traceCtx, "payment_service", 50*time.Millisecond)

	root.End()
	log.Printf("Trace %s completed: %d spans, total duration %v",
		traceCtx.TraceID[:8], 4, root.EndTime.Sub(root.StartTime))
}
```

## Step-by-step execution

1. `NewTrace()` generates a 128-bit trace ID and 64-bit root span ID using `crypto/rand`.
2. `Tracer.StartSpan` creates a span with the given context and operation name, recording the start time.
3. `simulateServiceCall` creates a child span context (same trace ID, new span ID, parent span ID set) and calls `StartSpan`.
4. The child span simulates work with `time.Sleep` and records attributes.
5. `span.End()` records the end time; the duration is computed when printed.
6. Each span prints its trace ID prefix, span ID prefix, operation name, and duration.
7. The root span shows the total trace duration.

## Common mistakes

- Mistake: Not propagating the trace context through the call chain — passing `context.Background()` instead of the parent's context to downstream calls.
  - Why it happens: Developers create a new context for each function call instead of passing the received context.
  - Fix: Always pass `r.Context()` (which contains the trace context) to downstream functions and services.

- Mistake: Setting attributes with high-cardinality values like `user_id` or `session_id`.
  - Why it happens: Tracing backends index span attributes for search. High-cardinality attributes explode the index size and slow down queries.
  - Fix: Use attributes for low-cardinality metadata (service name, HTTP method, result status). Put high-cardinality data in span events or logs.

- Mistake: Creating too many spans by adding one for every function call.
  - Why it happens: Developers treat spans as debug logs, creating a span for every helper function.
  - Fix: Create spans at service boundaries (HTTP handler → service → database → external API). A typical request should have 5-15 spans, not 500.

## Debugging walkthrough

A service shows intermittent 5-second timeouts. Metrics show p99 latency is normal. Logs show no errors. The team cannot reproduce the issue locally.

**Symptom**: Some requests take 5+ seconds but most complete in 200ms.

**Investigation**: Add tracing to the request path. The trace waterfall reveals that the `payment_service` span sometimes takes 4.8 seconds. The `db_query` span (child of payment_service) is fast (20ms), but the `external_gateway` span shows a 4.5-second gap between start and end.

**Root cause**: The external payment gateway has a timeout of 5 seconds configured in the HTTP client, but no retry logic. When the gateway is slow, the request blocks for nearly the full timeout.

**Fix**: Reduce the HTTP client timeout to 2 seconds, add retry with exponential backoff (max 3 retries), and add a circuit breaker to stop calling a failing gateway.

```go
// Before: no tracing reveals the slow hop
// After: trace waterfall shows payment_service -> external_gateway taking 4.5s
```

## Production notes

- Sample aggressively in production. Tracing 100% of requests on a 1000 req/s service generates 1M spans/second — too expensive for storage and CPU. Use probabilistic sampling (1%) combined with rate-limited sampling. Always sample error traces at 100%.
- Use tail-based sampling: store all spans temporarily, then decide which traces to keep based on the complete trace (e.g., keep traces with any error span).
- Set a span processor that batches exports (e.g., every 1 second or every 512 spans) to avoid per-span network calls.
- Production tracing backends: Grafana Tempo (self-hosted, S3-backed), Jaeger (legacy), Honeycomb, Datadog, AWS X-Ray.

## Performance implications

- Span creation: allocating a span struct and generating random IDs takes ~500ns.
- Context propagation: storing and retrieving the span context from `context.Context` is O(1) but requires an interface lookup (~10ns).
- Exporting spans: serialization and network I/O are the main costs. Batching reduces per-span overhead from ~1ms to ~2µs.
- Sampling at 1% reduces tracing overhead by 99% while still capturing representative data.
- Without sampling, tracing adds 1-3% CPU overhead to a typical Go HTTP service. With 1% sampling, overhead drops below 0.1%.

## Practice task

Implement a `Trace` struct that supports:
1. Creating a root trace with `NewTrace(operation string) *Span`.
2. Creating child spans with `parent.StartChild(operation string) *Span`.
3. Adding string attributes and events to spans.
4. Printing the trace tree as indented text showing operation name, duration, and attributes.

Then write a function `simulateRequest() *Span` that creates spans for: HTTP handler → authentication → database query → response. Verify the trace tree has the correct parent-child relationships.

## Tests / verification

```bash
go test ./curriculum/modules/12-observability-diagnostics/lessons/08-tracing -v
go run ./curriculum/modules/12-observability-diagnostics/lessons/08-tracing
```

## Review questions

1. What is the difference between a trace ID and a span ID?
2. How does the W3C TraceContext `traceparent` header encode the parent-child relationship between spans?
3. Why should you sample only 1-5% of requests in production tracing?
4. If a trace has 4 spans across 3 services, how many unique span IDs and trace IDs exist?
5. What happens to the trace if a downstream service does not propagate the trace context in its outgoing HTTP requests?

## NEXT UP

OpenTelemetry — the unified API for traces, metrics, and logs that replaced OpenTracing and OpenCensus.
