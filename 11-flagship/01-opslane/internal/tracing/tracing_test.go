// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package tracing

import (
	"bytes"
	"context"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/logging"
)

func TestStartSpanRecordsNameAndCorrelation(t *testing.T) {
	t.Parallel()

	ctx := logging.WithCorrelationID(context.Background(), "trace-id-42")
	ctx, span := StartSpan(ctx, "db.query")

	if span.Name != "db.query" {
		t.Fatalf("span name = %q, want %q", span.Name, "db.query")
	}
	if span.CorrelationID != "trace-id-42" {
		t.Fatalf("correlation_id = %q, want %q", span.CorrelationID, "trace-id-42")
	}
	if span.StartTime.IsZero() {
		t.Fatal("start time should not be zero")
	}

	// Verify span is retrievable from context.
	got := SpanFromContext(ctx)
	if got != span {
		t.Fatal("expected same span from context")
	}
}

func TestSpanFromContextReturnsNilWhenNoSpan(t *testing.T) {
	t.Parallel()

	got := SpanFromContext(context.Background())
	if got != nil {
		t.Fatal("expected nil span from empty context")
	}
}

func TestEndSpanLogsCorrectFields(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	ctx := logging.WithCorrelationID(context.Background(), "log-span-id")
	_, span := StartSpan(ctx, "service.process")

	// Simulate some work.
	time.Sleep(1 * time.Millisecond)

	EndSpan(span, logger)

	output := buf.String()
	if !strings.Contains(output, "service.process") {
		t.Fatalf("expected span_name in log, got %q", output)
	}
	if !strings.Contains(output, "log-span-id") {
		t.Fatalf("expected correlation_id in log, got %q", output)
	}
	if !strings.Contains(output, "duration") {
		t.Fatalf("expected duration in log, got %q", output)
	}
}

func TestEndSpanNilSpanDoesNotPanic(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	// Should not panic.
	EndSpan(nil, logger)
	EndSpan(nil, nil)
}

func TestInjectCorrelationHeader(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("GET", "/downstream", nil)
	InjectCorrelationHeader(r, "outbound-id-7")

	got := r.Header.Get("X-Correlation-ID")
	if got != "outbound-id-7" {
		t.Fatalf("expected %q, got %q", "outbound-id-7", got)
	}
}

func TestInjectCorrelationHeaderEmptyIsNoOp(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("GET", "/downstream", nil)
	InjectCorrelationHeader(r, "")

	got := r.Header.Get("X-Correlation-ID")
	if got != "" {
		t.Fatalf("expected empty header, got %q", got)
	}
}

func TestExtractCorrelationHeader(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-Correlation-ID", "inbound-id-9")

	got := ExtractCorrelationHeader(r)
	if got != "inbound-id-9" {
		t.Fatalf("expected %q, got %q", "inbound-id-9", got)
	}
}

func TestExtractCorrelationHeaderMissing(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("GET", "/", nil)
	got := ExtractCorrelationHeader(r)

	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}
