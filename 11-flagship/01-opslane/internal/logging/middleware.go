// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package logging

import (
	"log/slog"
	"net/http"
	"time"
)

// statusRecorder wraps http.ResponseWriter to capture the status code
// written by downstream handlers. This is necessary because the
// standard ResponseWriter does not expose the status after WriteHeader
// is called.
type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (sr *statusRecorder) WriteHeader(code int) {
	if !sr.wroteHeader {
		sr.status = code
		sr.wroteHeader = true
	}
	sr.ResponseWriter.WriteHeader(code)
}

func (sr *statusRecorder) Write(b []byte) (int, error) {
	if !sr.wroteHeader {
		sr.status = http.StatusOK
		sr.wroteHeader = true
	}
	return sr.ResponseWriter.Write(b)
}

// RequestLogger returns HTTP middleware that:
//
//  1. Extracts or generates a correlation ID from the X-Correlation-ID header.
//  2. Injects the correlation ID into the request context.
//  3. Sets the correlation ID on the response header so clients can reference it.
//  4. Wraps the ResponseWriter to capture the status code.
//  5. Logs a structured request-completion line with method, path, status,
//     latency, and correlation_id.
//
// This middleware replaces ad-hoc request logging with a single structured
// entry that contains enough context to trace a request through the system.
func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Extract or generate correlation ID.
			correlationID := CorrelationIDFromRequest(r)
			ctx := WithCorrelationID(r.Context(), correlationID)

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
