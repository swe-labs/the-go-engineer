// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"runtime/debug"
	"sync"
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

// CORS allows browser-based clients to call the API during local development.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RateLimit bounds how many requests one client IP can make inside a fixed window.
func RateLimit(maxRequests int, window time.Duration) func(http.Handler) http.Handler {
	type clientWindow struct {
		count   int
		resetAt time.Time
	}

	var mu sync.Mutex
	clients := make(map[string]clientWindow)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if maxRequests <= 0 || window <= 0 {
				next.ServeHTTP(w, r)
				return
			}

			clientIP := clientAddress(r)
			now := time.Now()

			mu.Lock()
			state := clients[clientIP]
			if state.resetAt.IsZero() || now.After(state.resetAt) {
				state = clientWindow{resetAt: now.Add(window)}
			}
			state.count++
			clients[clientIP] = state
			allowed := state.count <= maxRequests
			mu.Unlock()

			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(`{"error":{"code":"rate_limited","message":"rate limit exceeded"}}`))
				return
			}

			next.ServeHTTP(w, r)
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

func clientAddress(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return r.RemoteAddr
}
