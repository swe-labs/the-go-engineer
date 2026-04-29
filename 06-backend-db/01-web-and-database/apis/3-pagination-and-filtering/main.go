// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Pagination and Filtering
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to implement offset-based pagination.
//   - How to filter results using query parameters.
//   - Best practices for returning metadata (total count, next page).
//
// WHY THIS MATTERS:
//   - Sending thousands of records in a single response is slow, consumes
//     huge memory, and can crash both your server and the client's browser.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/apis/3-pagination-and-filtering
//
// KEY TAKEAWAY:
//   - Always set a maximum 'limit' to protect your server.
// ============================================================================

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Stage 06: APIs - Pagination and Filtering
//
//   - Limit/Offset: The standard SQL-like approach
//   - Query Filters: ?status=active&category=tools
//   - Pagination Metadata: { "page": 1, "total": 100 }
//
// ENGINEERING DEPTH:
//   Offset-based pagination is easy to implement but gets slower as the
//   offset increases (the database still has to scan through the
//   skipped records). For very large datasets, "Cursor-based" pagination
//   (using an ID or timestamp) is more performant and avoids issues
//   when records are added/deleted while a user is scrolling.

type Product struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

var allProducts []Product

func init() {
	// Generate 100 dummy products
	for i := 1; i <= 100; i++ {
		category := "Electronics"
		if i%3 == 0 {
			category = "Books"
		} else if i%5 == 0 {
			category = "Home"
		}
		allProducts = append(allProducts, Product{ID: i, Name: fmt.Sprintf("Product %d", i), Category: category})
	}
}

func main() {
	fmt.Println("=== Pagination and Filtering ===")
	fmt.Println()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /products", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		// 1. Get and Parse Pagination Params
		limit, _ := strconv.Atoi(query.Get("limit"))
		if limit <= 0 || limit > 50 {
			limit = 10 // Default and Max limit
		}

		offset, _ := strconv.Atoi(query.Get("offset"))
		if offset < 0 {
			offset = 0
		}

		// 2. Get Filtering Params
		category := query.Get("category")

		// 3. Apply Filtering
		var filtered []Product
		for _, p := range allProducts {
			if category == "" || p.Category == category {
				filtered = append(filtered, p)
			}
		}

		// 4. Apply Pagination
		start := offset
		if start > len(filtered) {
			start = len(filtered)
		}
		end := start + limit
		if end > len(filtered) {
			end = len(filtered)
		}

		results := filtered[start:end]

		// 5. Build Response with Metadata
		response := struct {
			Data   []Product `json:"data"`
			Total  int       `json:"total"`
			Limit  int       `json:"limit"`
			Offset int       `json:"offset"`
		}{
			Data:   results,
			Total:  len(filtered),
			Limit:  limit,
			Offset: offset,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("  Server starting on :8090...")
	fmt.Println("  Try these URLs:")
	fmt.Println("    - Default:     http://localhost:8090/products")
	fmt.Println("    - Paginated:   http://localhost:8090/products?limit=5&offset=10")
	fmt.Println("    - Filtered:    http://localhost:8090/products?category=Books")
	fmt.Println("    - Both:        http://localhost:8090/products?category=Books&limit=2")
	fmt.Println()

	err := http.ListenAndServe(":8090", mux)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: API.4 protobuf-basics")
	fmt.Println("Current: API.3 (pagination-and-filtering)")
	fmt.Println("Previous: API.2 (api-versioning-strategies)")
	fmt.Println("---------------------------------------------------")
}
