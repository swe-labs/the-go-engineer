// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package tracing

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/logging"
)

// SpanContext holds timing and identity data for a named operation.
// This is a lightweight teaching implementation - production systems
// would use OpenTelemetry spans with trace/span IDs and sampling.
type SpanContext struct {
	Name          string
	CorrelationID string
	StartTime     time.Time
}

// spanContextKey is the context key for the current span.
type spanContextKey struct{}

// StartSpan begins a named operation span. It reads the correlation ID
// from the context and records the start time. The returned context
// carries the span so EndSpan can access it later.
//
// This pattern teaches the concept of span propagation without
// requiring an external tracing library.
func StartSpan(ctx context.Context, name string) (context.Context, *SpanContext) {
	span := &SpanContext{
		Name:          name,
		CorrelationID: logging.CorrelationID(ctx),
		StartTime:     time.Now(),
	}

	return context.WithValue(ctx, spanContextKey{}, span), span
}

// EndSpan logs the span's duration. If logger is nil, the span is
// silently completed (useful in tests).
func EndSpan(span *SpanContext, logger *slog.Logger) {
	if span == nil || logger == nil {
		return
	}

	duration := time.Since(span.StartTime)
	logger.Info("span",
		slog.String("span_name", span.Name),
		slog.Duration("duration", duration),
		slog.String("correlation_id", span.CorrelationID),
	)
}

// SpanFromContext retrieves the current span from the context.
// Returns nil if no span was started.
func SpanFromContext(ctx context.Context) *SpanContext {
	span, _ := ctx.Value(spanContextKey{}).(*SpanContext)
	return span
}

// InjectCorrelationHeader sets the X-Correlation-ID header on an
// outbound HTTP request. Use this when making calls to downstream
// services so the correlation ID propagates across service boundaries.
func InjectCorrelationHeader(r *http.Request, correlationID string) {
	if correlationID != "" {
		r.Header.Set("X-Correlation-ID", correlationID)
	}
}

// ExtractCorrelationHeader reads the X-Correlation-ID from an inbound
// HTTP request. Returns empty string if absent.
func ExtractCorrelationHeader(r *http.Request) string {
	return r.Header.Get("X-Correlation-ID")
}
