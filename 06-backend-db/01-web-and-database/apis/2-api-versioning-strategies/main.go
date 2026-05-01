// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: API Versioning Strategies
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Why API versioning is critical for long-lived services.
//   - How to implement URL-based versioning (/v1, /v2).
//   - How to handle deprecated endpoints gracefully.
//
// WHY THIS MATTERS:
//   - Once an API is public, you cannot change it without breaking
//     existing clients. Versioning allows you to evolve your service
//     while keeping old mobile apps or integrations working.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/apis/2-api-versioning-strategies
//
// KEY TAKEAWAY:
//   - "The first version of your API is your last version of flexibility."
//     Plan for v2 before you launch v1.
// ============================================================================

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Stage 06: APIs - Versioning Strategies
//
//   - URL Versioning: /v1/users (Explicit, easy to cache)
//   - Header Versioning: X-API-Version: 2 (Clean URLs, harder to test)
//   - Supporting multiple versions in the same binary
//
// ENGINEERING DEPTH:
//   Versioning is about managing the "Contract" between you and your
//   consumers. A breaking change includes: renaming a field, changing
//   a data type (int -> string), or changing the status code.
//   A non-breaking change includes: adding a new optional field or
//   adding a new endpoint.

type UserV1 struct {
	Name string `json:"full_name"`
}

type UserV2 struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func main() {
	fmt.Println("=== API Versioning Strategies ===")
	fmt.Println()

	mux := http.NewServeMux()

	// 1. Version 1: Simple full name
	mux.HandleFunc("GET /v1/user", func(w http.ResponseWriter, r *http.Request) {
		user := UserV1{Name: "Rasel Hossen"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})

	// 2. Version 2: Split first and last name (Breaking change!)
	mux.HandleFunc("GET /v2/user", func(w http.ResponseWriter, r *http.Request) {
		user := UserV2{FirstName: "Rasel", LastName: "Hossen"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})

	fmt.Println("  Server starting on :8089...")
	fmt.Println("  Compare the responses:")
	fmt.Println("    - curl http://localhost:8089/v1/user")
	fmt.Println("    - curl http://localhost:8089/v2/user")
	fmt.Println()

	err := http.ListenAndServe(":8089", mux)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: API.3 -> 06-backend-db/01-web-and-database/apis/3-pagination-and-filtering")
	fmt.Println("Current: API.2 (api-versioning-strategies)")
	fmt.Println("Previous: API.1 (rest-design-principles)")
	fmt.Println("---------------------------------------------------")
}
