// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: net/http basics
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The fundamental components of a Go HTTP server.
//   - How to define handlers using 'http.HandlerFunc'.
//   - How to route requests using 'http.ServeMux'.
//
// WHY THIS MATTERS:
//   - Go's standard library is powerful enough to build production-grade
//     web services without any external frameworks.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/http-servers/1-net-http-basics
//
// KEY TAKEAWAY:
//   - A server is just a listener that passes requests to handlers.
// ============================================================================

package main

import (
	"fmt"
	"net/http"
)

// Stage 06: Backend - net/http basics
//
//   - http.HandlerFunc - turning functions into handlers
//   - http.ServeMux - the request router
//   - http.ResponseWriter - writing the response
//   - *http.Request - reading the incoming data
//
// ENGINEERING DEPTH:
//   The `http.Handler` interface is the core of Go's web ecosystem:
//   type Handler interface { ServeHTTP(ResponseWriter, *Request) }
//   By keeping this interface minimal, Go ensures that middleware, routers,
//   and business logic can all be composed together seamlessly.

func main() {
	fmt.Println("=== net/http basics ===")
	fmt.Println()

	// 1. Define a router (ServeMux)
	// A ServeMux matches the URL of each incoming request against a list of
	// registered patterns and calls the handler for the pattern that
	// most closely matches the URL.
	mux := http.NewServeMux()

	// 2. Register handlers
	// http.HandlerFunc is an adapter that allows the use of ordinary functions
	// as HTTP handlers.
	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/health", healthHandler)

	fmt.Println("  Server starting on :8080...")
	fmt.Println("  Try visiting: http://localhost:8080")
	fmt.Println("  Try visiting: http://localhost:8080/health")
	fmt.Println("  (Press Ctrl+C to stop)")

	// 3. Start the server
	// ListenAndServe listens on the TCP network address addr and then calls
	// Serve with handler to handle requests on incoming connections.
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("  Error starting server: %v\n", err)
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HS.2 -> 06-backend-db/01-web-and-database/http-servers/2-routing-patterns")
	fmt.Println("Current: HS.1 (net/http basics)")
	fmt.Println("Previous: FS.8 (fs-testing-seam)")
	fmt.Println("---------------------------------------------------")
}

// helloHandler responds with a simple greeting.
// helloHandler (Function): responds with a simple greeting.
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// ResponseWriter is used to construct the HTTP response.
	fmt.Fprintln(w, "Hello, Go Engineer!")
	fmt.Fprintf(w, "You requested: %s\n", r.URL.Path)
}

// healthHandler is a common pattern for monitoring.
// healthHandler (Function): is a common pattern for monitoring.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
