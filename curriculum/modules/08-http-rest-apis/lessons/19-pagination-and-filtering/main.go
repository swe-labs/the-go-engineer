package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
}

var allProducts = []Product{
	{ID: 1, Name: "Laptop", Category: "electronics", Price: 999.99},
	{ID: 2, Name: "Mouse", Category: "electronics", Price: 29.99},
	{ID: 3, Name: "Keyboard", Category: "electronics", Price: 89.99},
	{ID: 4, Name: "Desk", Category: "furniture", Price: 299.99},
	{ID: 5, Name: "Chair", Category: "furniture", Price: 499.99},
	{ID: 6, Name: "Monitor", Category: "electronics", Price: 349.99},
	{ID: 7, Name: "Notebook", Category: "stationery", Price: 4.99},
	{ID: 8, Name: "Pen Set", Category: "stationery", Price: 12.99},
	{ID: 9, Name: "Lamp", Category: "furniture", Price: 79.99},
	{ID: 10, Name: "Tablet", Category: "electronics", Price: 399.99},
}

type PaginatedResponse struct {
	Data       []Product `json:"data"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	Total      int       `json:"total"`
	TotalPages int       `json:"total_pages"`
}

type CursorResponse struct {
	Data    []Product `json:"data"`
	Cursor  int       `json:"cursor"`
	HasMore bool      `json:"has_more"`
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	filtered := filterProducts(r)
	sorted := sortProducts(filtered, r.URL.Query().Get("sort"))

	page, limit := parsePagination(r)
	total := len(sorted)
	totalPages := (total + limit - 1) / limit

	start := (page - 1) * limit
	if start >= total {
		json.NewEncoder(w).Encode(PaginatedResponse{
			Data:       []Product{},
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		})
		return
	}

	end := start + limit
	if end > total {
		end = total
	}

	json.NewEncoder(w).Encode(PaginatedResponse{
		Data:       sorted[start:end],
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

func parsePagination(r *http.Request) (page, limit int) {
	page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return
}

func filterProducts(r *http.Request) []Product {
	q := strings.ToLower(r.URL.Query().Get("q"))
	category := strings.ToLower(r.URL.Query().Get("category"))
	minPrice, _ := strconv.ParseFloat(r.URL.Query().Get("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(r.URL.Query().Get("max_price"), 64)

	var result []Product
	for _, p := range allProducts {
		if q != "" && !strings.Contains(strings.ToLower(p.Name), q) {
			continue
		}
		if category != "" && strings.ToLower(p.Category) != category {
			continue
		}
		if minPrice > 0 && p.Price < minPrice {
			continue
		}
		if maxPrice > 0 && p.Price > maxPrice {
			continue
		}
		result = append(result, p)
	}
	return result
}

func sortProducts(products []Product, sort string) []Product {
	if sort == "" {
		return products
	}

	result := make([]Product, len(products))
	copy(result, products)

	switch sort {
	case "price_asc":
		for i := 0; i < len(result); i++ {
			for j := i + 1; j < len(result); j++ {
				if result[j].Price < result[i].Price {
					result[i], result[j] = result[j], result[i]
				}
			}
		}
	case "price_desc":
		for i := 0; i < len(result); i++ {
			for j := i + 1; j < len(result); j++ {
				if result[j].Price > result[i].Price {
					result[i], result[j] = result[j], result[i]
				}
			}
		}
	case "name":
		for i := 0; i < len(result); i++ {
			for j := i + 1; j < len(result); j++ {
				if result[j].Name < result[i].Name {
					result[i], result[j] = result[j], result[i]
				}
			}
		}
	}

	return result
}

func cursorProductsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	cursorStr := r.URL.Query().Get("cursor")
	cursor := 0
	if cursorStr != "" {
		cursor, _ = strconv.Atoi(cursorStr)
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 5
	}

	filtered := filterProducts(r)

	var start int
	for i, p := range filtered {
		if p.ID > cursor {
			start = i
			break
		}
	}

	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	hasMore := end < len(filtered)
	var nextCursor int
	if hasMore && len(filtered) > 0 {
		nextCursor = filtered[end-1].ID
	}

	json.NewEncoder(w).Encode(CursorResponse{
		Data:    filtered[start:end],
		Cursor:  nextCursor,
		HasMore: hasMore,
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/products", productsHandler)
	mux.HandleFunc("/api/products/cursor", cursorProductsHandler)

	log.Println("Products API with pagination on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

var _ = fmt.Sprintf
