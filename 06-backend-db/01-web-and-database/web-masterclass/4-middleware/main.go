// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Web Masterclass - Middleware
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to implement the Middleware pattern in Go.
//   - The standard signature: func(http.Handler) http.Handler.
//   - How to chain multiple middleware together (The Onion Pattern).
//   - Essential production middleware: Logging, Recovery, and Security Headers.
//
// WHY THIS MATTERS:
//   - Middleware allows you to handle "Cross-Cutting Concerns"-tasks that
//     apply to every request-in a single, reusable location. This keeps
//     your main handler logic focused purely on business goals.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/web-masterclass/4-middleware
//
// KEY TAKEAWAY:
//   - Middleware is just a function that wraps another function.
// ============================================================================

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Stage 06: Web Masterclass - Middleware
//
//   - The HandlerFunc type: Casting logic to interfaces
//   - next.ServeHTTP: Passing control down the chain
//   - Panic Recovery: Ensuring server uptime
//
// ENGINEERING DEPTH:
//   In Go, a middleware chain is literally just a series of nested
//   function calls. When the innermost handler finishes, execution
//   "unwinds" back through the middleware. This allows you to perform
//   actions both BEFORE a request (like authentication) and AFTER a
//   request (like measuring latency or logging the status code).

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handleHome)
	mux.HandleFunc("GET /panic", handlePanic)

	// CHAINING MIDDLEWARE
	// We wrap our router (the 'mux') in multiple layers.
	// Order: Logger -> Recovery -> SecureHeaders -> Router
	handler := loggerMiddleware(recoveryMiddleware(secureHeadersMiddleware(mux)))

	fmt.Println("=== Web Masterclass: Middleware ===")
	fmt.Println("  🚀 Server starting on http://localhost:8083")
	fmt.Println()

	log.Fatal(http.ListenAndServe(":8083", handler))

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: MC.5 -> 06-backend-db/01-web-and-database/web-masterclass/5-sessions")
	fmt.Println("Current: MC.4 (middleware)")
	fmt.Println("Previous: MC.3 (templates)")
	fmt.Println("---------------------------------------------------")
}

// 1. Logger Middleware
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Pass the request to the next handler
		next.ServeHTTP(w, r)

		// Log the result
		log.Printf("  [LOG] %s %s took %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// 2. Recovery Middleware
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use defer to catch panics
		defer func() {
			if err := recover(); err != nil {
				log.Printf("  [RECOVERY] Caught panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// 3. Security Headers Middleware
func secureHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the secure, logged, and safe homepage!")
}

func handlePanic(w http.ResponseWriter, r *http.Request) {
	panic("Oops! Something went wrong.")
}
