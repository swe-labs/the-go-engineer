package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

// ============================================================================
// Section 13: Web Masterclass — Middleware
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The middleware signature: func(http.Handler) http.Handler
//   - Middleware chaining (wrapping handlers like layers of an onion)
//   - Common middleware: logging, security headers, panic recovery
//   - How to compose middleware in a clean, readable way
//
// ENGINEERING DEPTH:
//   A Middleware chain is natively just an execution Call Stack. When a request
//   comes in, the outermost middleware executes until it hits `next.ServeHTTP(w,r)`.
//   At that exact moment, the current function pauses, pushes its state onto the
//   Call Stack, and jumps execution into the next nested middleware. It keeps
//   drilling down until it hits the final route handler. As the final handler
//   returns, the Call Stack legally unwinds back "up" the chain layer by layer.
//   This means code BEFORE `next` runs on the inbound request, and code AFTER
//   `next` runs exactly on the outbound response.
//
// RUN: go run ./13-web-masterclass/4-middleware
// ============================================================================

// secureHeaders adds security-related HTTP headers to every response.
// This is a MUST for production web applications.
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "deny")
		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// Enable XSS protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		// Control referrer information
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// logRequest logs every incoming HTTP request with structured logging.
func logRequest(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Call the next handler
			next.ServeHTTP(w, r)

			// Log after the request completes
			logger.Info("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote", r.RemoteAddr),
				slog.Duration("latency", time.Since(start)),
			)
		})
	}
}

// recoverPanic recovers from panics inside handlers and returns a 500.
// Without this, a panic would crash the entire server.
func recoverPanic(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Deferred function runs even if panic occurs
			defer func() {
				if err := recover(); err != nil {
					// Set Connection: close to tell client to stop sending
					w.Header().Set("Connection", "close")
					logger.Error("panic recovered",
						slog.Any("error", err),
						slog.String("stack", string(debug.Stack())),
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from a middleware-protected handler!")
	})

	mux.HandleFunc("GET /panic", func(w http.ResponseWriter, r *http.Request) {
		// This panic will be caught by recoverPanic middleware
		panic("something went terribly wrong!")
	})

	// Chain middleware: outermost runs first
	// Request flow: recoverPanic → secureHeaders → logRequest → handler
	handler := recoverPanic(logger)(secureHeaders(logRequest(logger)(mux)))

	log.Println("Starting server with middleware on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
