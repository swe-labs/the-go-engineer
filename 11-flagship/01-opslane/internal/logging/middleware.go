// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package logging

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/otel"
)

// statusRecorder (Struct): wraps http.ResponseWriter to capture the status code for logging
type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

// WriteHeader (Method): captures the HTTP status code before delegating to the underlying writer
func (sr *statusRecorder) WriteHeader(code int) {
	if !sr.wroteHeader {
		sr.status = code
		sr.wroteHeader = true
	}
	sr.ResponseWriter.WriteHeader(code)
}

// Write (Method): captures the default 200 OK status if WriteHeader wasn't explicitly called
func (sr *statusRecorder) Write(b []byte) (int, error) {
	if !sr.wroteHeader {
		sr.status = http.StatusOK
		sr.wroteHeader = true
	}
	return sr.ResponseWriter.Write(b)
}

// RequestLogger (Function): returns HTTP middleware that logs structured request completions with correlation ID
func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Extract or generate correlation ID.
			correlationID := CorrelationIDFromRequest(r)
			ctx := WithCorrelationID(r.Context(), correlationID)
			ctx = otel.WithTraceParent(ctx, r.Header.Get("traceparent"))

			// Echo the correlation ID back to the client.
			w.Header().Set("X-Correlation-ID", correlationID)

			// Wrap the writer to capture the status code.
			recorder := &statusRecorder{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			// Serve the request with the enriched context.
			next.ServeHTTP(recorder, r.WithContext(ctx))

			// Log the completed request.
			logger.Info("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", recorder.status),
				slog.Duration("latency", time.Since(start)),
				slog.String("correlation_id", correlationID),
				slog.String("remote_addr", r.RemoteAddr),
			)
		})
	}
}
