// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"log"
	"net/http"
)

// ============================================================================
// Section 13: Web Masterclass — Routing Exercise
// Level: Beginner
// ============================================================================
//
// EXERCISE: Build a multi-route web server
//
// Your assignment is to expand this basic web server to handle multiple
// API endpoints for a bookstore.
//
// REQUIREMENTS:
//  1. [ ] The root path ("/") should welcome users to the bookstore.
//         Ensure any unknown path (like "/foo") returns a 404 Not Found
//         (the current handleHome function does this).
//  2. [ ] Add a "/books" route that returns a list of 3 fake books.
//         (Hint: use mux.HandleFunc)
//  3. [ ] Add a "/health" route that returns "status: ok" for monitoring.
//  4. [ ] The server must run on port 8080.
//
// ENGINEERING DEPTH:
//   In Go 1.22+, `http.ServeMux` received a massive upgrade. It now supports
//   method-based routing directly!
//   Example: mux.HandleFunc("GET /users", handleGetUsers)
//
//   Why use a custom `http.NewServeMux()` instead of `http.HandleFunc` (the DefaultServeMux)?
//   Using the DefaultServeMux (a global variable) is a security risk if you
//   import third-party packages, as they can silently register malicious routes
//   on your public server. ALWAYS create your own mux!
//
// RUN: go run ./13-web-masterclass/1-routing/exercise
// ============================================================================

func main() {
	// Create a new, isolated request multiplexer (router).
	mux := http.NewServeMux()

	// Register routes.
	// You need to add "/books" and "/health" here.
	mux.HandleFunc("/", handleHome)

	// TODO: Add mux.HandleFunc("GET /books", handleBooks)
	// TODO: Add mux.HandleFunc("GET /health", handleHealth)

	// Start the server
	port := 8080
	fmt.Printf("🚀 Bookstore server starting exactly on http://localhost:%d\n", port)

	// http.ListenAndServe blocks forever unless there's a fatal error.
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Fatalf("Server crashed: %v", err)
	}
}

// handleHome handles the root path.
func handleHome(w http.ResponseWriter, r *http.Request) {
	// The pattern "/" matches EVERYTHING that doesn't have a more specific route.
	// This means "/unknown-path" will hit handleHome!
	// We must manually reject paths that aren't EXACTLY "/" with a 404.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Welcome to The Go Engineer Bookstore!")
}

// TODO: Create handleBooks function here

// TODO: Create handleHealth function here
