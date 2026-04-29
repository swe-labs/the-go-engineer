// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Web Masterclass - Routing
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to use the Go 1.22+ 'http.ServeMux' for advanced routing.
//   - How to define method-based routes (GET, POST).
//   - How to extract path parameters using the {param} syntax.
//
// WHY THIS MATTERS:
//   - Routing is the entry point for every web request. Understanding the
//     native Go routing capabilities allows you to build clean, efficient
//     APIs without the need for external frameworks.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/web-masterclass/1-routing
//
// KEY TAKEAWAY:
//   - Pattern matching is now a first-class citizen in the Go standard library.
// ============================================================================

package main

import (
	"fmt"
	"log"
	"net/http"
)

// Stage 06: Web Masterclass - Routing
//
//   - Method-based routing: GET /about
//   - Path parameters: /posts/{id}
//   - Multiplexing with http.NewServeMux
//
// ENGINEERING DEPTH:
//   Prior to Go 1.22, developers had to use third-party libraries for
//   even basic path parameter support. The new ServeMux implementation
//   is highly optimized and supports complex pattern matching while
//   maintaining the "Zero-Allocation" philosophy of the Go standard
//   library's network stack.

func main() {
	// 1. Create a new multiplexer
	mux := http.NewServeMux()

	// 2. Define routes
	mux.HandleFunc("GET /", handleHome)
	mux.HandleFunc("GET /about", handleAbout)

	// {id} is a wildcard parameter
	mux.HandleFunc("GET /posts/{id}", handleGetPost)

	mux.HandleFunc("GET /api/health", handleHealth)
	mux.HandleFunc("POST /api/posts", handleCreatePost)

	fmt.Println("=== Web Masterclass: Routing ===")
	fmt.Println("  🚀 Server starting on http://localhost:8080")
	fmt.Println()
	fmt.Println("  Try these routes:")
	fmt.Println("    - http://localhost:8080/")
	fmt.Println("    - http://localhost:8080/about")
	fmt.Println("    - http://localhost:8080/posts/42")

	// 3. Start the server
	log.Fatal(http.ListenAndServe(":8080", mux))

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: MC.2 dependency-injection")
	fmt.Println("Current: MC.1 (routing)")
	fmt.Println("Previous: DM.1 (embedded-migrations)")
	fmt.Println("---------------------------------------------------")
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	// "/" matches everything, so we check for exact match for the home page
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintln(w, "Welcome to the Go Web Masterclass!")
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "About: Building production-grade Go backends.")
}

func handleGetPost(w http.ResponseWriter, r *http.Request) {
	// Extract the {id} parameter using PathValue
	id := r.PathValue("id")
	fmt.Fprintf(w, "Viewing post with ID: %s\n", id)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status":"ok"}`)
}

func handleCreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, `{"message":"post created"}`)
}
