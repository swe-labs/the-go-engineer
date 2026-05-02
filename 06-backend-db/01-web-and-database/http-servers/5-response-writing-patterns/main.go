// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Response Writing Patterns
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to send structured JSON responses using 'json.NewEncoder'.
//   - How to set HTTP status codes and headers correctly.
//   - Best practices for consistent API response structures.
//
// WHY THIS MATTERS:
//   - Clear and consistent responses make your API easy for clients to
//     consume and debug.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/http-servers/5-response-writing-patterns
//
// KEY TAKEAWAY:
//   - Use 'w.Header().Set' BEFORE 'w.WriteHeader' or writing the body.
// ============================================================================

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// User represents a domain model.
// User (Struct): represents a domain model.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// APIResponse is a standard wrapper for all our API responses.
// APIResponse (Struct): is a standard wrapper for all our API responses.
type APIResponse struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func main() {
	fmt.Println("=== Response Writing Patterns ===")
	fmt.Println()

	mux := http.NewServeMux()

	// 1. Standard JSON Success Response
	mux.HandleFunc("GET /user", getUserHandler)

	// 2. Structured Error Response
	mux.HandleFunc("GET /error", errorHandler)

	fmt.Println("  Server starting on :8084...")
	fmt.Println("  Test Success Response:")
	fmt.Println("    curl http://localhost:8084/user")
	fmt.Println("  Test Error Response:")
	fmt.Println("    curl http://localhost:8084/error")
	fmt.Println()

	err := http.ListenAndServe(":8084", mux)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HS.6 -> 06-backend-db/01-web-and-database/http-servers/6-error-handling-middleware")
	fmt.Println("Current: HS.5 (response-writing-patterns)")
	fmt.Println("Previous: HS.4 (request-parsing-and-validation)")
	fmt.Println("---------------------------------------------------")
}

// getUserHandler (Function): runs the get user handler step and keeps its inputs, outputs, or errors visible.
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := User{ID: 1, Username: "gopher", Email: "gopher@golang.org"}

	// 1. Set the Content-Type header FIRST.
	w.Header().Set("Content-Type", "application/json")

	// 2. Set the status code (default is 200 OK, but being explicit is good).
	w.WriteHeader(http.StatusOK)

	// 3. Encode the data directly to the ResponseWriter.
	// NewEncoder is preferred over Marshal for performance and streaming.
	json.NewEncoder(w).Encode(APIResponse{
		Data:    user,
		Message: "User retrieved successfully",
	})
}

// errorHandler (Function): runs the error handler step and keeps its inputs, outputs, or errors visible.
func errorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Set an error status code.
	w.WriteHeader(http.StatusNotFound)

	json.NewEncoder(w).Encode(APIResponse{
		Error: "User not found",
	})
}
