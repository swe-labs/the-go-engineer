// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Web Masterclass - Pagination
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to slice large datasets using SQL LIMIT and OFFSET.
//   - How to calculate pagination metadata (Current Page, Total Pages).
//   - How to build a standard JSON response for paginated data.
//   - The performance trade-offs of offset-based pagination.
//
// WHY THIS MATTERS:
//   - No production system should ever return thousands of records in
//     a single request. Pagination keeps your application fast, reduces
//     memory usage, and provides a better experience for your users.
//     Mastering this pattern is essential for building scalable APIs.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/web-masterclass/9-pagination
//
// KEY TAKEAWAY:
//   - Divide and conquer. Never load more data than you need to display.
// ============================================================================

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

// Stage 06: Web Masterclass - Pagination
//
//   - Limit: The number of items per page
//   - Offset: The number of items to skip
//   - Metadata: Info about the overall collection
//
// ENGINEERING DEPTH:
//   While `LIMIT ? OFFSET ?` is simple to implement, it has a
//   "Scaling Penalty." To calculate `OFFSET 10000`, the database must
//   physically scan through the first 10,000 rows just to find where
//   to start. For extremely large datasets, engineers use "Cursor
//   Pagination" (e.g., `WHERE id > ? LIMIT 20`), which uses database
//   indexes to jump directly to the next page in O(log N) time.

// Metadata holds information about the paginated response.
type Metadata struct {
	CurrentPage  int  `json:"current_page"`
	PageSize     int  `json:"page_size"`
	FirstPage    int  `json:"first_page"`
	LastPage     int  `json:"last_page"`
	TotalRecords int  `json:"total_records"`
	HasNext      bool `json:"has_next"`
}

func main() {
	// 1. Simulate a large collection of items
	const totalItems = 105

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/items", func(w http.ResponseWriter, r *http.Request) {
		// 2. Extract pagination parameters from the query string
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}

		pageSize := 10 // Hardcoded for simplicity

		// 3. Calculate offset and metadata
		offset := (page - 1) * pageSize
		totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))

		meta := Metadata{
			CurrentPage:  page,
			PageSize:     pageSize,
			FirstPage:    1,
			LastPage:     totalPages,
			TotalRecords: totalItems,
			HasNext:      page < totalPages,
		}

		// 4. Generate dummy items for the current "page"
		items := []string{}
		start := offset
		end := offset + pageSize
		if end > totalItems {
			end = totalItems
		}

		for i := start; i < end; i++ {
			items = append(items, fmt.Sprintf("Item %d", i+1))
		}

		// 5. Send paginated JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"metadata": meta,
			"items":    items,
		})
	})

	fmt.Println("=== Web Masterclass: Pagination ===")
	fmt.Println("  🚀 Server starting on http://localhost:8088")
	fmt.Println()
	fmt.Println("  Try these links:")
	fmt.Println("    - http://localhost:8088/api/items?page=1")
	fmt.Println("    - http://localhost:8088/api/items?page=2")
	fmt.Println("    - http://localhost:8088/api/items?page=11 (Last page)")

	log.Fatal(http.ListenAndServe(":8088", mux))

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: MC.10 -> 06-backend-db/01-web-and-database/web-masterclass/10-comments")
	fmt.Println("Current: MC.9 (pagination)")
	fmt.Println("Previous: MC.8 (posts-crud)")
	fmt.Println("---------------------------------------------------")
}
