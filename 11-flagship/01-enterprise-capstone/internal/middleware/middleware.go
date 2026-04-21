// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

// ============================================================================
// Internal Package: Middleware
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Wrapping Handlers in structured logging (`slog`)
//   - Defending the server from Panics
//   - Adding critical Content-Type security headers to prevent Cross-site scripting (XSS)
//
// ENGINEERING DEPTH:
//   "Layering" middleware wraps your core application logic in an onion-like shell.
//   Every request *must* pass through these layers first before hitting `ServeMux`.
// ============================================================================

// RecoverPanic acts as the outermost onion shell.
// By wrapping `next.ServeHTTP` inside a function containing a `defer recover()`,
// we guarantee that if ANY downstream handler throws a fatal panic, this
// middleware will instantly catch it, write a 500 status code, and prevent
// the entire server process from crashing!
func RecoverPanic(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Hard-close the TCP connection since the state is unstable.
					w.Header().Set("Connection", "close")
					logger.Error("panic caught by middleware",
						slog.Any("error", err),
						slog.String("stack", string(debug.Stack())),
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r) // Execute downstream chain
		})
	}
}

func LogRequest(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Duration("latency", time.Since(start)),
			)
		})
	}
}

// SecureHeaders sets mandatory security directives for browsers.
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent Clickjacking (disallow being loaded in an iframe)
		w.Header().Set("X-Frame-Options", "deny")
		// Prevent MIME-sniffing bypasses
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// Stop execution if an XSS attack is detected by the browser natively
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		next.ServeHTTP(w, r)
	})
}
