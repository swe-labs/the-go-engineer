// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Request Parsing and Validation
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to decode JSON request bodies safely.
//   - How to read and parse URL query parameters.
//   - Patterns for validating incoming data before processing it.
//
// WHY THIS MATTERS:
//   - A server is only as secure as its input validation.
//     Parsing data correctly prevents crashes and security vulnerabilities.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/http-servers/4-request-parsing-and-validation
//
// KEY TAKEAWAY:
//   - Never trust client input. Always decode, then validate.
// ============================================================================

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// UserRequest represents the expected payload for creating a user.
type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

// Validate checks if the request data meets our business rules.
func (u *UserRequest) Validate() error {
	if u.Username == "" {
		return fmt.Errorf("username is required")
	}
	if !strings.Contains(u.Email, "@") {
		return fmt.Errorf("invalid email format")
	}
	if u.Age < 18 {
		return fmt.Errorf("user must be at least 18 years old")
	}
	return nil
}

func main() {
	fmt.Println("=== Request Parsing and Validation ===")
	fmt.Println()

	mux := http.NewServeMux()

	// 1. Parsing JSON Body
	mux.HandleFunc("POST /users", createUserHandler)

	// 2. Parsing Query Parameters
	mux.HandleFunc("GET /search", searchHandler)

	fmt.Println("  Server starting on :8083...")
	fmt.Println("  Test JSON parsing:")
	fmt.Println("    curl -X POST -d '{\"username\":\"jane\", \"email\":\"jane@example.com\", \"age\":25}' http://localhost:8083/users")
	fmt.Println("  Test Query parsing:")
	fmt.Println("    curl \"http://localhost:8083/search?q=golang&limit=10\"")
	fmt.Println()

	err := http.ListenAndServe(":8083", mux)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HS.5 -> 06-backend-db/01-web-and-database/http-servers/5-response-writing-patterns")
	fmt.Println("Current: HS.4 (request-parsing-and-validation)")
	fmt.Println("Previous: HS.3 (middleware-pattern)")
	fmt.Println("---------------------------------------------------")
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var req UserRequest

	// Use json.NewDecoder to stream the body into the struct.
	// This is more efficient than ReadAll + Unmarshal for large payloads.
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate the business rules
	if err := req.Validate(); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	fmt.Fprintf(w, "User %q created successfully!\n", req.Username)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// r.URL.Query() returns a url.Values map containing the query parameters.
	query := r.URL.Query()

	q := query.Get("q")
	limit := query.Get("limit")

	if q == "" {
		http.Error(w, "Search query 'q' is required", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Searching for: %s (Limit: %s)\n", q, limit)
}
