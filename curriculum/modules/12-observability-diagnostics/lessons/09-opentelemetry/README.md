# OpenTelemetry

## Learning objective

Instrument Go code with OpenTelemetry tracing: create a tracer provider, start and end spans, add attributes and events, and export spans to a backend.

## Why this matters

Before OpenTelemetry, every observability backend had its own SDK. Instrumenting with Jaeger's SDK locked you into Jaeger. Switching to Datadog meant rewriting all instrumentation. OpenTelemetry (OTel) provides a single vendor-neutral API that works with any backend. It is the CNCF standard for observability instrumentation, adopted by Kubernetes, AWS, Azure, Google Cloud, and every major observability vendor. Learning OTel means you instrument once and ship anywhere.

## Mental model

OpenTelemetry is a universal adapter for observability — like USB-C for instrumentation. One connector (OTel API) works with any device (backend). You write code against the OTel API (`go.opentelemetry.io/otel`), and at startup you plug in the SDK with the appropriate exporter (Jaeger, Datadog, Grafana Tempo).

The OTel Collector is like a USB hub: it receives OTLP (OpenTelemetry Protocol) data from many services, processes it (filter, sample, add attributes), and routes it to one or more backends. You can send traces to both Datadog and S3-backed Tempo simultaneously without changing any application code.

The analogy breaks because OTel also handles metrics and logs, not just traces, and the API/SDK separation requires explicit initialization rather than plug-and-play.

## Core idea

OpenTelemetry is a collection of APIs, SDKs, and tools for generating, collecting, and exporting telemetry data (traces, metrics, logs). Key concepts:

- **TracerProvider**: the entry point. Creates `Tracer` instances. Configured once at application startup.
- **Tracer**: creates `Span` instances. Named (e.g., `"orders-service"`) to identify the instrumentation source.
- **Span**: represents a unit of work. Has name, context, attributes, events, status, and timestamps.
- **Context Propagation**: spans are linked through `context.Context`. A span's context is stored in the Go context and propagated to downstream calls.
- **Exporter**: sends completed spans to a backend. Common exporters: OTLP gRPC, OTLP HTTP, stdout (for debugging).
- **SpanProcessor**: hooks into span start/end events. `BatchSpanProcessor` queues spans and exports them in batches.

OTel uses W3C TraceContext for propagation (the `traceparent` header) and Baggage for propagating arbitrary key-value pairs across service boundaries.

## Under the hood

When `tracer.Start(ctx, "operation")` is called:

1. The tracer calls `propagators.TextMapPropagator.Extract` to get the parent span context from the incoming context (parsed from HTTP headers at the service boundary).
2. A new `SpanContext` is created: same trace ID, new span ID, parent span ID set to the extracted span ID.
3. A new `Span` struct is allocated with start timestamp and the new context.
4. The span is stored in the returned context using a private context key.
5. When `span.End()` is called, `BatchSpanProcessor.OnEnd(span)` enqueues the span.
6. A background goroutine in `BatchSpanProcessor` wakes up every `batchTimeout` (default 1s) or when the batch size reaches `maxQueueSize` (default 512), serializes spans as OTLP protobuf, and calls `exporter.ExportSpans(ctx, spans)`.
7. The OTLP exporter sends the protobuf payload to the configured endpoint via gRPC or HTTP.

The OTLP protobuf schema (`opentelemetry-proto`) defines `TraceService/Export` for traces, `MetricService/Export` for metrics, and `LogsService/Export` for logs. Each uses the same resource and instrumentation library metadata.

## How Go uses it

The Go OTel implementation follows the API/SDK split:

- `go.opentelemetry.io/otel` — API package (interfaces for TracerProvider, Tracer, Span).
- `go.opentelemetry.io/otel/sdk` — SDK package (implementations of the API).
- `go.opentelemetry.io/otel/exporters/otlp/otlptrace` — OTLP trace exporter.
- `go.opentelemetry.io/otel/exporters/stdout/stdouttrace` — stdout exporter for debugging.

Best practice: depend on the API package in your library code and wire the SDK at the application entry point (`main()`). This allows users of your library to choose any OTel-compatible backend.

## Go example

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

func newExporter() (*stdouttrace.Exporter, error) {
	return stdouttrace.New(stdouttrace.WithPrettyPrint())
}

func newResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("order-service"),
		semconv.ServiceVersionKey.String("1.0.0"),
		attribute.String("environment", "development"),
	)
}

func newTraceProvider(exp *stdouttrace.Exporter) *sdktrace.TracerProvider {
	bsp := sdktrace.NewBatchSpanProcessor(exp)
	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(newResource()),
		sdktrace.WithSpanProcessor(bsp),
	)
}

func main() {
	exp, err := newExporter()
	if err != nil {
		panic(err)
	}

	tp := newTraceProvider(exp)
	otel.SetTracerProvider(tp)

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()

	tracer := otel.Tracer("order-service")

	ctx := context.Background()
	ctx, span := tracer.Start(ctx, "create_order")
	defer span.End()

	span.SetAttributes(
		attribute.String("order.id", "ORD-12345"),
		attribute.Float64("order.amount", 99.50),
		attribute.String("order.currency", "USD"),
	)
	span.AddEvent("order.created", trace.WithAttributes(attribute.String("user", "alice")))

	if err := processPayment(ctx, tracer); err != nil {
		span.RecordError(err)
		span.SetStatus(sdktrace.StatusError, "payment failed")
	}

	span.SetStatus(sdktrace.StatusOk, "order created successfully")
}

func processPayment(ctx context.Context, tracer trace.Tracer) error {
	ctx, span := tracer.Start(ctx, "process_payment")
	defer span.End()

	span.SetAttributes(attribute.String("payment.method", "credit_card"))
	span.AddEvent("charging_card")

	time.Sleep(50 * time.Millisecond)

	if time.Now().Unix()%2 == 0 {
		return errors.New("insufficient funds")
	}
	return nil
}
```

## Step-by-step execution

1. `newExporter()` creates a stdout exporter that prints spans as formatted JSON.
2. `newTraceProvider()` creates a `TracerProvider` with batch span processing and always-sample sampling.
3. `otel.SetTracerProvider(tp)` sets the global tracer provider — all subsequent `otel.Tracer()` calls use this provider.
4. `tracer.Start(ctx, "create_order")` creates a root span named `create_order`.
5. Attributes and events are recorded on the root span.
6. `processPayment` calls `tracer.Start(ctx, "process_payment")` which extracts the parent span from the context and creates a child span.
7. After `processPayment` returns, the child span ends, and `BatchSpanProcessor` enqueues it.
8. When the root span ends (deferred `span.End()`), both spans are exported via stdout as pretty-printed JSON.
9. `tp.Shutdown()` flushes any remaining spans and stops the exporter.

## Common mistakes

- Mistake: Importing and using the SDK package directly in library code instead of the API.
  - Why it happens: `go.opentelemetry.io/otel/sdk/trace` provides the concrete `TracerProvider`, but importing it in a library ties users to that specific SDK.
  - Fix: Library code should only import `go.opentelemetry.io/otel` and `go.opentelemetry.io/otel/trace`. The SDK is wired at the application entry point.

- Mistake: Creating a new `TracerProvider` on every request instead of once at startup.
  - Why it happens: Developers treat `sdktrace.NewTracerProvider()` like a lightweight constructor, not realizing it starts goroutines and opens network connections.
  - Fix: Create the `TracerProvider` once in `main()`, store it as a global or inject it via dependency injection.

- Mistake: Not calling `TracerProvider.Shutdown()` on application exit.
  - Why it happens: Forgetting the deferred shutdown call. Unflushed spans in the batch processor are lost.
  - Fix: Always defer `tp.Shutdown(ctx)` after creating the provider. The shutdown flushes pending spans and stops background goroutines.

## Debugging walkthrough

A service has OpenTelemetry tracing configured, but no spans appear in the backend.

**Symptom**: The Jaeger or Grafana Tempo UI shows no traces for the instrumented service.

**Investigation**:
1. Check if `otel.SetTracerProvider(tp)` is called before any instrumentation runs — without it, the default no-op provider silently drops all spans.
2. Verify the exporter is correctly configured by using `stdouttrace` exporter first — if spans appear on stdout, the issue is the exporter connection.
3. Check the collector/backend endpoint URL: is it `localhost:4317` (gRPC) or `localhost:4318` (HTTP)?
4. Verify network connectivity: `telnet collector 4317` from the service pod.
5. If using HTTP exporter, check for TLS errors or missing certificates.

**Root cause**: The service was deployed with `OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4317` but the collector was listening on port 4318 (HTTP) instead of 4317 (gRPC).

**Fix**: Change the endpoint to `http://otel-collector:4318` or configure the collector to listen on port 4317 for gRPC.

## Production notes

- Use `OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_SERVICE_NAME`, and `OTEL_RESOURCE_ATTRIBUTES` environment variables instead of hardcoding configuration. OTel Go SDK reads these automatically.
- Set `OTEL_TRACES_SAMPLER=parentbased_traceidratio` and `OTEL_TRACES_SAMPLER_ARG=0.01` for 1% probabilistic sampling.
- Deploy the OTel Collector as a sidecar or DaemonSet in Kubernetes. The collector handles batching, retries, filtering, and routing to multiple backends.
- For production, use the OTLP gRPC exporter (`otlptracegrpc`) instead of stdout. Configure TLS and authentication if the collector is remote.

## Performance implications

- Span creation with the SDK allocates a `Span` struct, generates a random span ID, and stores the span in context: ~1µs per span.
- Batch span processor adds ~500ns per span for enqueuing. Exporting a batch of 512 spans takes ~10-50ms (serialization + network).
- Always-sample strategy records every span: at 1000 req/s with 5 spans per request, that is 5000 spans/second. At 1% sampling, it drops to 50 spans/second — negligible overhead.
- The OTel Go SDK is designed for zero-allocation hot paths in the API layer. SDK allocations only happen when spans are created and ended.

## Practice task

Write a function `InstrumentedHandler(next http.Handler, tracer trace.Tracer) http.Handler` that:
1. Extracts the trace context from the incoming request headers.
2. Starts a span named after the HTTP method and path (e.g., `"GET /api/users"`).
3. Sets attributes for HTTP method, path, and status code.
4. Records an event for any error status code (>=400).
5. Ends the span after the handler completes.
6. Passes the span context to `next` via `r.WithContext(ctx)`.

Then write a test that creates a `TracerProvider` with a stdout exporter, wraps a handler, sends a request via `httptest`, and verifies the span appears in stdout.

## Tests / verification

```bash
go test ./curriculum/modules/12-observability-diagnostics/lessons/09-opentelemetry -v
go run ./curriculum/modules/12-observability-diagnostics/lessons/09-opentelemetry
```

## Review questions

1. What is the difference between the OpenTelemetry API (`go.opentelemetry.io/otel/trace`) and the SDK (`go.opentelemetry.io/otel/sdk/trace`)?
2. How does `TracerProvider.Shutdown()` ensure no spans are lost?
3. What happens if `otel.SetTracerProvider()` is never called in `main()`?
4. Why should you use environment variables (`OTEL_SERVICE_NAME`, `OTEL_EXPORTER_OTLP_ENDPOINT`) for OTel configuration?
5. Describe the path of a span from `tracer.Start()` to the backend storage.

## NEXT UP

Alerting mindset — translating observability data into actionable alerts that on-call engineers trust.
