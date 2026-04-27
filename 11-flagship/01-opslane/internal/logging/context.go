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
type correlationIDKey struct{}

// WithCorrelationID stores a correlation ID on the context. Every log
// line produced from this context can include the ID, making it possible
// to trace one request through HTTP handlers, services, and workers.
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationIDKey{}, id)
}

// CorrelationID retrieves the correlation ID from the context.
// Returns an empty string if no ID was set.
func CorrelationID(ctx context.Context) string {
	id, _ := ctx.Value(correlationIDKey{}).(string)
	return id
}

// GenerateCorrelationID produces a 16-byte hex-encoded random string.
// Using crypto/rand avoids adding an external UUID dependency while still
// providing collision-resistant identifiers.
func GenerateCorrelationID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// CorrelationIDFromRequest extracts the X-Correlation-ID header from
// an HTTP request. If the header is absent or empty, it generates a
// new correlation ID.
func CorrelationIDFromRequest(r *http.Request) string {
	if id := r.Header.Get("X-Correlation-ID"); id != "" {
		return id
	}
	return GenerateCorrelationID()
}

// ContextAttrs extracts structured log attributes from a context.
// It reads correlation_id from the logging context. Auth identity
// (tenant_id, user_id) should be added by the caller using
// auth.IdentityFromContext, keeping this package free of auth imports.
func ContextAttrs(ctx context.Context) []slog.Attr {
	var attrs []slog.Attr
	if id := CorrelationID(ctx); id != "" {
		attrs = append(attrs, slog.String("correlation_id", id))
	}
	return attrs
}
