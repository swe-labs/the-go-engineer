// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

// Package tracing provides lightweight distributed tracing for the Opslane backend.
// It implements a simple span-based pattern for tracking operation duration and
// correlation IDs across service boundaries.
package tracing

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/logging"
)

// SpanContext (Struct): holds timing and identity data for a named tracing span
type SpanContext struct {
	Name          string
	CorrelationID string
	StartTime     time.Time
}

// spanContextKey is the context key for the current span.
// spanContextKey (Struct): unexported context key type for the current span
type spanContextKey struct{}

// StartSpan (Function): begins a named operation span with correlation ID propagation
func StartSpan(ctx context.Context, name string) (context.Context, *SpanContext) {
	span := &SpanContext{
		Name:          name,
		CorrelationID: logging.CorrelationID(ctx),
		StartTime:     time.Now(),
	}

	return context.WithValue(ctx, spanContextKey{}, span), span
}

// EndSpan (Function): logs the span's duration; silently completes if logger is nil (useful in tests)
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

// SpanFromContext (Function): retrieves the current span from context; returns nil if absent
func SpanFromContext(ctx context.Context) *SpanContext {
	span, _ := ctx.Value(spanContextKey{}).(*SpanContext)
	return span
}

// InjectCorrelationHeader (Function): sets X-Correlation-ID on outbound HTTP requests for trace propagation
func InjectCorrelationHeader(r *http.Request, correlationID string) {
	if correlationID != "" {
		r.Header.Set("X-Correlation-ID", correlationID)
	}
}

// ExtractCorrelationHeader (Function): reads X-Correlation-ID from an inbound HTTP request
func ExtractCorrelationHeader(r *http.Request) string {
	return r.Header.Get("X-Correlation-ID")
}
