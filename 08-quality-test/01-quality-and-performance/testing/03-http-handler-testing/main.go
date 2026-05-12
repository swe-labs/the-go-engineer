// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 08: Quality & Testing
// Title: HTTP Handler Testing
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Testing HTTP handlers without starting a real web server
//   - Using httptest.NewRecorder and httptest.NewRequest
//   - Testing response status codes, headers, and body content
//
// WHY THIS MATTERS:
//   Starting a real HTTP server in tests is slow and brittle. httptest gives
//   you a fake request/response cycle that runs in the same process, making
//   handler tests fast, deterministic, and portable.
//
// RUN:
//   go test ./08-quality-test/01-quality-and-performance/testing/03-http-handler-testing
// ============================================================================

package main

import (
	"fmt"
	"io"
	"net/http"
)

// HelloWorldHandler (Function): returns a static greeting response.
func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World!"))
}

// EchoHandler (Function): reads the request body and echoes it back in the response.
// Boundary: accepts only POST method; returns 405 for other methods.
func EchoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// KEY TAKEAWAY:
// - httptest.NewRecorder captures handler output without a real server.
// - httptest.NewRequest creates synthetic requests for any method/path.
// - Handler tests are fast, deterministic, and portable.
//
// NEXT UP: TE.4 -> 08-quality-test/01-quality-and-performance/testing/benchmarks
func main() {
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: TE.4 -> 08-quality-test/01-quality-and-performance/testing/benchmarks")
	fmt.Println("Run    : go test -bench=. -benchmem ./08-quality-test/01-quality-and-performance/testing/benchmarks")
	fmt.Println("Current: TE.3 (http-handler-testing)")
	fmt.Println("---------------------------------------------------")
}
