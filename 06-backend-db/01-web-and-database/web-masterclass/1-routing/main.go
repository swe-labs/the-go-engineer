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
// Stage 06: Web Masterclass — Routing
// Level: Beginner -> Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Creating an HTTP server with net/http
//   - Registering route handlers with http.NewServeMux
//   - Method-based routing (Go 1.22+)
//   - Path parameters with {param} syntax
//   - The difference between HandleFunc and Handle
//
// ENGINEERING DEPTH:
//   Prior to Go 1.22, `http.ServeMux` was a naive map that just matched exact
//   strings, forcing the entire industry to use 3rd-party routers like `chi`
//   and `gorilla/mux`. With Go 1.22, the core team completely rewrote `ServeMux`
//   using an advanced algorithmic Radix Tree that supports pattern matching
//   (`GET /posts/{id}`) directly in the standard library! This parses millions of
//   incoming URLs per second with almost zero memory allocations.
//
// RUN: go run ./06-backend-db/01-web-and-database/web-masterclass/1-routing
// VISIT: http://localhost:8080
// ============================================================================

func main() {
	// http.NewServeMux is the standard library's HTTP request multiplexer (router).
	// It matches incoming requests to registered patterns and calls handlers.
	mux := http.NewServeMux()

	// Basic route — handles ALL HTTP methods on "/"
	mux.HandleFunc("/", handleHome)

	// Go 1.22+ method-based routing: "METHOD /path"
	// Only GET requests will match this route
	mux.HandleFunc("GET /about", handleAbout)

	// Path parameters with {param} syntax (Go 1.22+)
	// The parameter is extracted with r.PathValue("id")
	mux.HandleFunc("GET /posts/{id}", handleGetPost)

	// Separate handlers for different methods on the same path
	mux.HandleFunc("GET /api/health", handleHealth)
	mux.HandleFunc("POST /api/posts", handleCreatePost)

	// Static file serving
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	fmt.Println("🚀 Server starting on http://localhost:8080")
	fmt.Println("   Routes:")
	fmt.Println("   GET  /            — Home page")
	fmt.Println("   GET  /about       — About page")
	fmt.Println("   GET  /posts/{id}  — Get post by ID")
	fmt.Println("   GET  /api/health  — Health check (JSON)")
	fmt.Println("   POST /api/posts   — Create a post (JSON)")

	// ListenAndServe starts an HTTP server on the given address.
	// It blocks until the server shuts down or encounters an error.
	log.Fatal(http.ListenAndServe(":8080", mux))
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: WM.2 dependency injection")
	fmt.Println("   Current: WM.1 (routing)")
	fmt.Println("---------------------------------------------------")
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	// "/" matches everything without a more specific handler,
	// so we explicitly check for exact match to avoid catch-all behavior.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintln(w, "Welcome to the Go Web Masterclass!")
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "About: Learning Go Web Development")
}

func handleGetPost(w http.ResponseWriter, r *http.Request) {
	// r.PathValue extracts named parameters from the URL pattern.
	// Pattern "GET /posts/{id}" with URL "/posts/42" → id = "42"
	id := r.PathValue("id")
	fmt.Fprintf(w, "Viewing post with ID: %s\n", id)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"ok"}`)
}

func handleCreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, `{"message":"post created"}`)
}
