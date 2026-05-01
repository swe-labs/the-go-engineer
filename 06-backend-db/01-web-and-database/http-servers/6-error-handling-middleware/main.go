// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Error Handling Middleware
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to centralize error reporting using a custom handler type.
//   - How to build a Recovery middleware to catch and handle panics.
//   - How to simplify your business logic by moving error responses
//     to a shared layer.
//
// WHY THIS MATTERS:
//   - Handling errors in every single handler leads to "Boilerplate Bloat".
//     Centralized handling ensures consistent error formats and easier logging.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/http-servers/6-error-handling-middleware
//
// KEY TAKEAWAY:
//   - A custom handler type that returns 'error' makes your code more readable.
// ============================================================================

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

// Stage 06: Backend - Error Handling Middleware
//
//   - AppHandler: A handler type that returns error
//   - Recovery middleware: Safeguarding against panics
//   - Unified error JSON structure
//
// ENGINEERING DEPTH:
//   By creating a type like `AppHandler func(...) error` and implementing
//   the `http.Handler` interface on it, we can transform standard Go
//   errors into HTTP responses in a single location. This pattern is
//   preferred over "Exception Handlers" in other languages because it
//   preserves Go's "errors are values" philosophy while reducing repetition.

// AppError represents a business error with an associated HTTP status code.
type AppError struct {
	Code    int
	Message string
	Err     error
}

func (ae AppError) Error() string {
	return ae.Message
}

// AppHandler is a custom handler type that returns an error.
type AppHandler func(w http.ResponseWriter, r *http.Request) error

// ServeHTTP makes AppHandler satisfy the http.Handler interface.
// This is where we centralize our error responding logic!
func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		var appErr AppError
		if errors.As(err, &appErr) {
			log.Printf("  [APP ERROR] %d: %v", appErr.Code, appErr.Err)
			http.Error(w, appErr.Message, appErr.Code)
		} else {
			log.Printf("  [INTERNAL ERROR]: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func main() {
	fmt.Println("=== Error Handling and Recovery ===")
	fmt.Println()

	mux := http.NewServeMux()

	// 1. Success case
	mux.Handle("GET /ok", AppHandler(okHandler))

	// 2. Expected error case
	mux.Handle("GET /fail", AppHandler(failHandler))

	// 3. Panic case (handled by Recovery middleware)
	mux.Handle("GET /panic", AppHandler(panicHandler))

	// Wrap everything in Recovery middleware
	handler := Recovery(mux)

	fmt.Println("  Server starting on :8085...")
	fmt.Println("  Try these:")
	fmt.Println("    curl http://localhost:8085/ok")
	fmt.Println("    curl http://localhost:8085/fail")
	fmt.Println("    curl http://localhost:8085/panic")
	fmt.Println()

	err := http.ListenAndServe(":8085", handler)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HS.7 -> 06-backend-db/01-web-and-database/http-servers/7-server-timeouts")
	fmt.Println("Current: HS.6 (error-handling-middleware)")
	fmt.Println("Previous: HS.5 (response-writing-patterns)")
	fmt.Println("---------------------------------------------------")
}

func okHandler(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintln(w, "Everything is fine.")
	return nil
}

func failHandler(w http.ResponseWriter, r *http.Request) error {
	// Return a structured business error
	return AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid input provided",
		Err:     errors.New("db record not found"),
	}
}

func panicHandler(w http.ResponseWriter, r *http.Request) error {
	// Simulate an unexpected crash
	panic("something went horribly wrong!")
}

// Recovery is a middleware that recovers from panics and logs them.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("  [RECOVERED] Panic: %v", err)
				http.Error(w, "Critical Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
