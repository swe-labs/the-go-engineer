package main

import (
	"bytes"
	"context"
	"io"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type mockExporter struct {
	spans []sdktrace.ReadOnlySpan
}

func (m *mockExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	m.spans = append(m.spans, spans...)
	return nil
}

func (m *mockExporter) Shutdown(ctx context.Context) error {
	return nil
}

func TestSpanCreation(t *testing.T) {
	exp := &mockExporter{}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(exp)),
	)
	otel.SetTracerProvider(tp)
	defer tp.Shutdown(context.Background())

	tracer := otel.Tracer("test-service")
	ctx := context.Background()

	ctx, span := tracer.Start(ctx, "test_operation")
	span.SetAttributes(attribute.String("test.key", "test.value"))
	span.End()

	// flush
	tp.ForceFlush(context.Background())

	if len(exp.spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(exp.spans))
	}

	got := exp.spans[0]
	if got.Name() != "test_operation" {
		t.Errorf("expected span name 'test_operation', got '%s'", got.Name())
	}
}

func TestSpanAttributesAndEvents(t *testing.T) {
	exp := &mockExporter{}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(exp)),
	)
	otel.SetTracerProvider(tp)
	defer tp.Shutdown(context.Background())

	tracer := otel.Tracer("test-service")
	ctx := context.Background()
	ctx, span := tracer.Start(ctx, "span_with_attrs")
	span.SetAttributes(
		attribute.String("service", "test"),
		attribute.Int("count", 42),
	)
	span.AddEvent("event_occurred")
	span.End()
	tp.ForceFlush(context.Background())

	span1 := exp.spans[0]
	if span1.Name() != "span_with_attrs" {
		t.Errorf("expected 'span_with_attrs', got '%s'", span1.Name())
	}
}

func TestParentChildSpan(t *testing.T) {
	exp := &mockExporter{}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(exp)),
	)
	otel.SetTracerProvider(tp)
	defer tp.Shutdown(context.Background())

	tracer := otel.Tracer("test-service")
	ctx := context.Background()

	ctx, parent := tracer.Start(ctx, "parent")
	_, child := tracer.Start(ctx, "child")
	child.End()
	parent.End()
	tp.ForceFlush(context.Background())

	if len(exp.spans) != 2 {
		t.Fatalf("expected 2 spans, got %d", len(exp.spans))
	}

	var parentSpan, childSpan sdktrace.ReadOnlySpan
	if exp.spans[0].Name() == "parent" {
		parentSpan = exp.spans[0]
		childSpan = exp.spans[1]
	} else {
		parentSpan = exp.spans[1]
		childSpan = exp.spans[0]
	}

	if childSpan.Parent().SpanID() != parentSpan.SpanContext().SpanID() {
		t.Error("child span should have parent's span ID as ParentSpanID")
	}
}

func TestStdoutExporterOutput(t *testing.T) {
	var buf bytes.Buffer
	exp, err := stdouttrace.New(
		stdouttrace.WithWriter(io.Writer(&buf)),
	)
	if err != nil {
		t.Fatalf("failed to create exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(exp)),
	)
	otel.SetTracerProvider(tp)
	defer tp.Shutdown(context.Background())

	tracer := otel.Tracer("test")
	_, span := tracer.Start(context.Background(), "stdout_test")
	span.End()
	tp.ForceFlush(context.Background())

	output := buf.String()
	if output == "" {
		t.Error("expected non-empty stdout output from exporter")
	}
	if !contains(output, "stdout_test") {
		t.Errorf("expected span name 'stdout_test' in output, got:\n%s", output)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsStr(s, substr)
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
