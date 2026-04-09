// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// ============================================================================
// Section 14: Application Architecture - Structured Logging: Context-Keyed Logger
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Storing a logger in context so every function in a request chain can log
//     with request-scoped fields (request_id, user_id, trace_id)
//   - Writing HTTP middleware that injects the logger into the request context
//   - Extracting the logger from context in handler functions
//   - Why this pattern is superior to a global logger singleton
//
// ENGINEERING DEPTH:
//   The biggest production logging problem is correlation — connecting all the
//   log lines from a single HTTP request so you can replay what happened.
//   This requires every log line in the chain (handler → service → repository)
//   to carry the same request_id. The options are:
//     1. Global logger (bad): no per-request fields, thread-unsafe
//     2. Pass logger as function parameter (verbose): every function signature changes
//     3. Store logger in context (idiomatic): one middleware line, works everywhere
//
//   Google's internal style guide mandates option 3. It integrates naturally
//   with the context-first pattern every I/O function already follows.
//
// RUN: go run ./14-application-architecture/structured-logging/2-context-logger
//   Then: curl http://localhost:8080/api/orders/42
// ============================================================================

// loggerKey is a private type to prevent key collisions in context.
// Using a plain string like "logger" is a bug — any package can write
// context.WithValue(ctx, "logger", something_else).
type loggerKey struct{}

// FromContext extracts the logger from context.
// If no logger was stored, it returns the global default logger.
// This safe fallback means functions can always call FromContext without
// checking for nil — the program never panics due to a missing logger.
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// WithLogger stores a logger in context.
// Always returns a new context — context values are immutable.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// LoggingMiddleware is the HTTP middleware that wires the pattern together.
// It runs before every request handler, creating a child logger loaded with
// request-scoped fields and injecting it into the request context.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Build a child logger with request-scoped fields.
		// Every log line for this request will automatically include these.
		requestLogger := slog.Default().With(
			slog.String("request_id", generateRequestID()),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
		)

		// Inject the logger into the request context.
		ctx := WithLogger(r.Context(), requestLogger)
		r = r.WithContext(ctx)

		requestLogger.Info("request started")

		// Wrap the ResponseWriter to capture the status code.
		rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		// Log AFTER the handler returns — we now know the final status.
		requestLogger.Info("request completed",
			slog.Int("status", rw.status),
			slog.Duration("latency", time.Since(start)),
		)
	})
}

// statusRecorder captures the HTTP status code set by the handler.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

var counter int

func generateRequestID() string {
	counter++
	return fmt.Sprintf("req_%06d", counter)
}

// ============================================================================
// Handler functions — they receive the logger from context
// ============================================================================

func handleGetOrder(w http.ResponseWriter, r *http.Request) {
	// Extract the logger anywhere in the call chain.
	// It already has request_id, method, and path pre-loaded.
	log := FromContext(r.Context())

	orderID := r.PathValue("id")
	log.Info("fetching order", slog.String("order_id", orderID))

	// Simulate calling downstream services
	order, err := fetchOrderFromDB(r.Context(), orderID)
	if err != nil {
		log.Error("database error", slog.Any("error", err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	log.Info("order retrieved", slog.String("status", order))
	fmt.Fprintf(w, `{"order_id":"%s","status":"%s"}`, orderID, order)
}

// fetchOrderFromDB is a downstream function that also uses the context logger.
// Notice it receives ctx (not the logger directly). This is the idiomatic way:
// functions never take *slog.Logger as a parameter — they use FromContext.
func fetchOrderFromDB(ctx context.Context, id string) (string, error) {
	log := FromContext(ctx) // Gets the same request-scoped logger
	log.Debug("executing SQL query",
		slog.String("table", "orders"),
		slog.String("id", id),
	)
	time.Sleep(5 * time.Millisecond) // Simulate DB latency
	return "shipped", nil
}

func main() {
	// Set up a JSON logger as the application default.
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/orders/{id}", handleGetOrder)

	// Wrap the entire mux with the logging middleware.
	handler := LoggingMiddleware(mux)

	slog.Info("server starting", slog.String("addr", ":8080"))
	http.ListenAndServe(":8080", handler)

	// KEY TAKEAWAY:
	// - Use a private context key type to prevent collisions
	// - FromContext() has a safe fallback — callers never need nil checks
	// - Middleware injects the request-scoped logger into context ONCE
	// - All downstream functions use FromContext — no logger parameters needed
	// - This pattern gives every log line the same request_id automatically
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: SL.3 custom slog.Handler")
	fmt.Println("   Current: SL.2 (context-keyed logger)")
	fmt.Println("---------------------------------------------------")
}
