// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Middleware - The Pattern
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The architecture of the middleware pattern in Go.
//   - How to wrap handlers to add cross-cutting concerns (logging, headers).
//   - How to chain multiple middleware layers together.
//
// WHY THIS MATTERS:
//   - Middleware allows you to keep your core business logic clean by
//     separating it from infrastructure concerns like auth, logging, and metrics.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/http-servers/3-middleware-pattern
//
// KEY TAKEAWAY:
//   - Middleware is just a function that wraps an 'http.Handler'.
// ============================================================================

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Stage 06: Backend - Middleware Pattern
//
//   - type Middleware func(http.Handler) http.Handler
//   - The "Decorator" pattern for web servers
//   - Executing code BEFORE and AFTER the main handler
//
// ENGINEERING DEPTH:
//   Middleware in Go works because the `http.Handler` interface is a single
//   method. By creating a new handler that calls the original one, we can
//   intercept the request flow. This is much more flexible than "Hooks" or
//   "Interceptors" in other languages because it relies on simple function
//   composition rather than complex reflection or magic.

func main() {
	fmt.Println("=== Middleware Pattern ===")
	fmt.Println()

	mux := http.NewServeMux()

	// Our core business logic handler
	helloHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, Go Engineer! The logic is running.")
	})

	// Wrap the handler in multiple layers of middleware
	// The execution order is: Logger -> SetJSONHeader -> helloHandler
	wrappedHandler := Logger(SetJSONHeader(helloHandler))

	mux.Handle("GET /", wrappedHandler)

	fmt.Println("  Server starting on :8082...")
	fmt.Println("  Try visiting: http://localhost:8082")
	fmt.Println("  Check the console for logs!")
	fmt.Println()

	err := http.ListenAndServe(":8082", mux)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HS.4 -> 06-backend-db/01-web-and-database/http-servers/4-request-parsing-and-validation")
	fmt.Println("Current: HS.3 (middleware-pattern)")
	fmt.Println("Previous: HS.2 (routing-patterns)")
	fmt.Println("---------------------------------------------------")
}

// Logger is a middleware that logs the details of every request.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 1. Code here runs BEFORE the main handler
		log.Printf("  Started %s %s", r.Method, r.URL.Path)

		// 2. Call the next handler in the chain
		next.ServeHTTP(w, r)

		// 3. Code here runs AFTER the main handler
		log.Printf("  Completed in %v", time.Since(start))
	})
}

// SetJSONHeader is a middleware that sets the Content-Type to application/json.
func SetJSONHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
