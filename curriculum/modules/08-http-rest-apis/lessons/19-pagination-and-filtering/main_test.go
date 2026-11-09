package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParsePagination(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		wantPage  int
		wantLimit int
	}{
		{"default values", "", 1, 20},
		{"custom page", "page=3", 3, 20},
		{"custom limit", "limit=10", 1, 10},
		{"limit too high", "limit=200", 1, 20},
		{"limit zero", "limit=0", 1, 20},
		{"page zero", "page=0", 1, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/products?"+tt.query, nil)
			page, limit := parsePagination(req)

			if page != tt.wantPage {
				t.Errorf("page = %d, want %d", page, tt.wantPage)
			}
			if limit != tt.wantLimit {
				t.Errorf("limit = %d, want %d", limit, tt.wantLimit)
			}
		})
	}
}

func TestFilterProducts(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		wantLen int
	}{
		{"no filter", "", 10},
		{"by category", "category=electronics", 5},
		{"by search", "q=laptop", 1},
		{"by min price", "min_price=500", 1},
		{"by max price", "max_price=10", 1},
		{"combined", "category=electronics&min_price=100", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/products?"+tt.query, nil)
			result := filterProducts(req)

			if len(result) != tt.wantLen {
				t.Errorf("got %d products, want %d", len(result), tt.wantLen)
			}
		})
	}
}

func TestSortProducts(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/products?sort=price_asc", nil)
	filtered := filterProducts(req)
	sorted := sortProducts(filtered, req.URL.Query().Get("sort"))

	if len(sorted) < 2 {
		t.Fatal("need at least 2 products")
	}
	// price_asc: first should be cheapest
	if sorted[0].Price > sorted[1].Price {
		t.Error("not sorted by price ascending")
	}

	// price_desc
	req2 := httptest.NewRequest(http.MethodGet, "/api/products?sort=price_desc", nil)
	sorted2 := sortProducts(filterProducts(req2), "price_desc")
	if sorted2[0].Price < sorted2[1].Price {
		t.Error("not sorted by price descending")
	}
}

func TestProductsHandlerPagination(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		wantStatus int
		wantLen    int
	}{
		{"first page default", "", http.StatusOK, 10},
		{"first page limit 3", "limit=3", http.StatusOK, 3},
		{"second page", "page=2&limit=3", http.StatusOK, 3},
		{"page beyond end", "page=10&limit=3", http.StatusOK, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/products?"+tt.query, nil)
			rec := httptest.NewRecorder()
			productsHandler(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			var pr PaginatedResponse
			if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
				t.Fatalf("json: %v", err)
			}
			if len(pr.Data) != tt.wantLen {
				t.Errorf("data length = %d, want %d", len(pr.Data), tt.wantLen)
			}
		})
	}
}

func TestProductsHandlerWithFilter(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/products?category=electronics&limit=2", nil)
	rec := httptest.NewRecorder()
	productsHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var pr PaginatedResponse
	json.NewDecoder(resp.Body).Decode(&pr)

	if pr.Total != 5 {
		t.Errorf("total = %d, want 5", pr.Total)
	}
	if pr.Page != 1 {
		t.Errorf("page = %d, want 1", pr.Page)
	}
}

func TestCursorPagination(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/products/cursor?limit=3", nil)
	rec := httptest.NewRecorder()
	cursorProductsHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var cr CursorResponse
	json.NewDecoder(resp.Body).Decode(&cr)

	if len(cr.Data) != 3 {
		t.Errorf("data length = %d, want 3", len(cr.Data))
	}
	if !cr.HasMore {
		t.Error("expected has_more = true")
	}

	// Fetch next page
	req2 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/products/cursor?limit=3&cursor=%d", cr.Cursor), nil)
	rec2 := httptest.NewRecorder()
	cursorProductsHandler(rec2, req2)

	resp2 := rec2.Result()
	defer resp2.Body.Close()

	var cr2 CursorResponse
	json.NewDecoder(resp2.Body).Decode(&cr2)

	if len(cr2.Data) == 0 {
		t.Error("expected data on second page")
	}
}

func TestProductsHandlerMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/products", nil)
	rec := httptest.NewRecorder()
	productsHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

func TestLessonCompiles(t *testing.T) {
}
