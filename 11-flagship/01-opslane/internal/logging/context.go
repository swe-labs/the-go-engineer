// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package logging

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
)

// correlationIDKey is the context key for the request correlation ID.
// Using an unexported struct type prevents collisions with other packages.
// correlationIDKey (Struct): unexported context key type for the request correlation ID
type correlationIDKey struct{}

// WithCorrelationID (Function): stores a correlation ID on the context for request tracing through logs
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationIDKey{}, id)
}

// CorrelationID (Function): retrieves the correlation ID from the context; returns empty string if absent
func CorrelationID(ctx context.Context) string {
	id, _ := ctx.Value(correlationIDKey{}).(string)
	return id
}

// GenerateCorrelationID (Function): produces a 16-byte hex-encoded random correlation ID using crypto/rand
func GenerateCorrelationID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// CorrelationIDFromRequest (Function): extracts or generates a correlation ID from an HTTP request header
func CorrelationIDFromRequest(r *http.Request) string {
	if id := r.Header.Get("X-Correlation-ID"); id != "" {
		return id
	}
	return GenerateCorrelationID()
}

// ContextAttrs (Function): extracts structured slog attributes (correlation_id) from a context
func ContextAttrs(ctx context.Context) []slog.Attr {
	var attrs []slog.Attr
	if id := CorrelationID(ctx); id != "" {
		attrs = append(attrs, slog.String("correlation_id", id))
	}
	return attrs
}
