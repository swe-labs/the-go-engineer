// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

// ============================================================================
// Section 13: Pagination
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Computing pagination metadata (total pages, has next/prev)
//   - SQL LIMIT/OFFSET queries for paginated data
//   - Building dynamic pagination links
//   - API response structure for paginated data
//
// ENGINEERING DEPTH:
//   While `LIMIT ? OFFSET ?` is standard, you must understand its algorithmic cost.
//   `OFFSET 50000 LIMIT 10` forces the database engine to locate, load, and scan
//   50,010 sequential rows off the disk, just to discard the first 50,000 and
//   return the final 10. For massive Google-scale systems, this O(n) scan creates
//   deadly latency ("The Offset Penalty"). In those scenarios, engineers abandon
//   Offsets entirely and use "Cursor Pagination" (`WHERE id > ? LIMIT 10`), taking
//   advantage of sub-millisecond O(log N) B-Tree Primary Key index seeks!
//
// RUN: go run ./13-web-masterclass/9-pagination
// ============================================================================

// Metadata holds pagination information for API responses.
type Metadata struct {
	CurrentPage  int  `json:"current_page"`
	PageSize     int  `json:"page_size"`
	TotalRecords int  `json:"total_records"`
	TotalPages   int  `json:"total_pages"`
	HasNext      bool `json:"has_next"`
	HasPrev      bool `json:"has_prev"`
	FirstPage    int  `json:"first_page"`
	LastPage     int  `json:"last_page"`
}

// computeMetadata calculates all pagination fields from the total record count.
func computeMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		HasNext:      page < totalPages,
		HasPrev:      page > 1,
		FirstPage:    1,
		LastPage:     totalPages,
	}
}

// PaginatedResponse is a generic API response wrapper for paginated data.
type PaginatedResponse struct {
	Data     any      `json:"data"`
	Metadata Metadata `json:"metadata"`
}

// generatePageLinks builds navigation links for an API response.
func generatePageLinks(baseURL string, meta Metadata) map[string]string {
	links := map[string]string{
		"self": fmt.Sprintf("%s?page=%d&page_size=%d", baseURL, meta.CurrentPage, meta.PageSize),
	}

	if meta.HasNext {
		links["next"] = fmt.Sprintf("%s?page=%d&page_size=%d", baseURL, meta.CurrentPage+1, meta.PageSize)
	}
	if meta.HasPrev {
		links["prev"] = fmt.Sprintf("%s?page=%d&page_size=%d", baseURL, meta.CurrentPage-1, meta.PageSize)
	}

	links["first"] = fmt.Sprintf("%s?page=%d&page_size=%d", baseURL, meta.FirstPage, meta.PageSize)
	links["last"] = fmt.Sprintf("%s?page=%d&page_size=%d", baseURL, meta.LastPage, meta.PageSize)

	return links
}

func main() {
	// Simulate 95 total items
	totalItems := 95

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/items", func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

		if page < 1 {
			page = 1
		}
		if pageSize < 1 || pageSize > 100 {
			pageSize = 10
		}

		// Compute pagination metadata
		meta := computeMetadata(totalItems, page, pageSize)
		links := generatePageLinks("/api/items", meta)

		// Generate fake items for this page
		start := (page - 1) * pageSize
		end := start + pageSize
		if end > totalItems {
			end = totalItems
		}

		items := make([]map[string]any, 0, end-start)
		for i := start; i < end; i++ {
			items = append(items, map[string]any{
				"id":   i + 1,
				"name": fmt.Sprintf("Item #%d", i+1),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data":     items,
			"metadata": meta,
			"links":    links,
		})
	})

	fmt.Println("Pagination demo: http://localhost:8080/api/items?page=1&page_size=10")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
