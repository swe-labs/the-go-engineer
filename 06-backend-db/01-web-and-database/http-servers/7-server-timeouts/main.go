// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Server Timeouts
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to configure production-grade timeouts on 'http.Server'.
//   - The difference between Read, Write, and Idle timeouts.
//   - How to use 'http.TimeoutHandler' to protect against hanging handlers.
//
// WHY THIS MATTERS:
//   - A server without timeouts is a "Sitting Duck". Slow clients or
//     hanging database calls can consume all server resources and cause
//     a total outage.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/http-servers/7-server-timeouts
//
// KEY TAKEAWAY:
//   - Never use 'http.ListenAndServe' in production. Always use an
//     explicit 'http.Server' with timeouts.
// ============================================================================

package main

import (
	"fmt"
	"net/http"
	"time"
)

// Stage 06: Backend - Server Timeouts
//
//   - ReadTimeout: Max duration for reading the entire request.
//   - WriteTimeout: Max duration before timing out writes of the response.
//   - IdleTimeout: Max amount of time to wait for the next request.
//   - TimeoutHandler: A wrapper that returns 503 if a handler takes too long.
//
// ENGINEERING DEPTH:
//   In Go, every request runs in its own goroutine. If a request never
//   finishes, the goroutine never exits. On a busy server, thousands of
//   hanging goroutines will eventually exhaust the system's memory and
//   file descriptors. Timeouts are the "Self-Preservation" mechanism of
//   a Go backend.

func main() {
	fmt.Println("=== Server Timeouts and Protection ===")
	fmt.Println()

	mux := http.NewServeMux()

	// 1. A normal, fast handler
	mux.HandleFunc("GET /fast", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Fast response!")
	})

	// 2. A slow handler that will be cut off by TimeoutHandler
	slowHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second) // Simulate a very slow database call
		fmt.Fprintln(w, "Slow response (if you see this, timeout failed!)")
	})

	// Wrap the slow handler with a 2-second limit
	// TimeoutHandler returns a handler that runs h with the given time limit.
	mux.Handle("GET /slow", http.TimeoutHandler(slowHandler, 2*time.Second, "Error: Request Timed Out"))

	// 3. Configure the HTTP Server
	// Instead of the helper http.ListenAndServe, we define the server
	// parameters explicitly.
	server := &http.Server{
		Addr:         ":8086",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,   // Protects against slow header/body sends
		WriteTimeout: 10 * time.Second,  // Protects against slow clients reading
		IdleTimeout:  120 * time.Second, // Keeps connections alive for re-use
	}

	fmt.Println("  Server starting on :8086...")
	fmt.Println("  Try these:")
	fmt.Println("    curl http://localhost:8086/fast")
	fmt.Println("    curl http://localhost:8086/slow (will time out after 2s)")
	fmt.Println()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Printf("  Error: %v\n", err)
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HS.8 -> 06-backend-db/01-web-and-database/http-servers/8-graceful-http-shutdown")
	fmt.Println("Current: HS.7 (server-timeouts)")
	fmt.Println("Previous: HS.6 (error-handling-middleware)")
	fmt.Println("---------------------------------------------------")
}
