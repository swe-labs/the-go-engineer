// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Routing Patterns
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to use modern Go 1.22+ routing patterns.
//   - How to match specific HTTP methods (GET, POST, etc.).
//   - How to extract path parameters using the '{name}' syntax.
//
// WHY THIS MATTERS:
//   - Clean routing makes your API predictable and easy to maintain.
//   - Path parameters allow you to create dynamic resource-based URLs.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/http-servers/2-routing-patterns
//
// KEY TAKEAWAY:
//   - Modern 'http.ServeMux' handles methods and variables without extra libraries.
// ============================================================================

package main

import (
	"fmt"
	"net/http"
)

// Stage 06: Backend - Routing Patterns
//
//   - Method matching: "GET /path"
//   - Path parameters: "/users/{id}"
//   - Wildcards: "/static/{filepath...}"
//   - Exact matching vs. Prefix matching
//
// ENGINEERING DEPTH:
//   Before Go 1.22, the standard library router was very limited. Engineers
//   often reached for external libraries like `gorilla/mux` or `chi`.
//   With the new enhanced routing patterns, the standard library is now
//   sufficient for most RESTful API designs, reducing dependency bloat
//   and keeping your binary small and secure.

func main() {
	fmt.Println("=== Routing Patterns (Go 1.22+) ===")
	fmt.Println()

	mux := http.NewServeMux()

	// 1. Method matching
	// We can prefix the pattern with an HTTP method to only match that verb.
	mux.HandleFunc("GET /items", listItemsHandler)
	mux.HandleFunc("POST /items", createItemHandler)

	// 2. Path Parameters
	// We can define variables in the path using curly braces {}.
	// Use r.PathValue("name") to retrieve the value in the handler.
	mux.HandleFunc("GET /items/{id}", getItemHandler)

	// 3. Wildcards (Tail matching)
	// The {name...} syntax matches any number of path segments.
	mux.HandleFunc("GET /files/{path...}", serveFileHandler)

	fmt.Println("  Server starting on :8081...")
	fmt.Println("  Try these URLs:")
	fmt.Println("  - GET  http://localhost:8081/items")
	fmt.Println("  - POST http://localhost:8081/items (requires curl or Postman)")
	fmt.Println("  - GET  http://localhost:8081/items/42")
	fmt.Println("  - GET  http://localhost:8081/files/images/logo.png")
	fmt.Println()

	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HS.3 -> 06-backend-db/01-web-and-database/http-servers/3-middleware-pattern")
	fmt.Println("Current: HS.2 (routing-patterns)")
	fmt.Println("Previous: HS.1 (net/http-basics)")
	fmt.Println("---------------------------------------------------")
}

// listItemsHandler (Function): runs the list items handler step and keeps its inputs, outputs, or errors visible.
func listItemsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Listing all items...")
}

// createItemHandler (Function): runs the create item handler step and keeps its inputs, outputs, or errors visible.
func createItemHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Item created!")
}

// getItemHandler (Function): runs the get item handler step and keeps its inputs, outputs, or errors visible.
func getItemHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the {id} parameter from the URL path.
	id := r.PathValue("id")
	fmt.Fprintf(w, "Fetching details for item ID: %s\n", id)
}

// serveFileHandler (Function): runs the serve file handler step and keeps its inputs, outputs, or errors visible.
func serveFileHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the full tail path.
	path := r.PathValue("path")
	fmt.Fprintf(w, "Serving file from virtual path: %s\n", path)
}
