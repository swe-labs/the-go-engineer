# Pagination and filtering

## Learning objective

Implement offset-based and cursor-based pagination for RESTful list endpoints, parse filter parameters from query strings, and apply sorting with configurable sort keys.

## Why this matters

Real APIs return thousands or millions of records. Returning all of them in one response is impossible -- it exhausts memory, clogs bandwidth, and creates a terrible client experience. Pagination, filtering, and sorting are the core mechanisms that make list endpoints usable at scale. Every production API you use (GitHub Issues, Stripe Charges, Slack Messages) implements these patterns.

## Mental model

Think of pagination like reading a book:

- **Offset pagination**: "Turn to page 42." You know the page number and the number of items per page. Simple, but if items are inserted/deleted between requests, pages shift (phantom reads).
- **Cursor pagination**: "Start from the last word on page 42." The cursor is a bookmark. You always continue from where you left off. Stable despite insertions, but no random access to pages.

Filtering is like asking for "all books by this author, published after 2020, under $30." Each query parameter narrows the result set.

Sorting is like saying "sort by price, cheapest first" or "by date, newest first."

## Core idea

**Offset pagination** uses `page` and `limit` parameters:

```
GET /products?page=2&limit=20
```

Response includes metadata so clients know the total and can compute pages:

```json
{
  "data": [...],
  "page": 2,
  "limit": 20,
  "total": 95,
  "total_pages": 5
}
```

**Cursor pagination** uses a `cursor` (opaque value, often the last item's ID) and `limit`:

```
GET /products?cursor=42&limit=20
```

Response indicates whether more results exist:

```json
{
  "data": [...],
  "cursor": 62,
  "has_more": true
}
```

**Filtering** uses query parameters:

| Parameter | Example | Effect |
|---|---|---|
| `q` | `?q=laptop` | Search by name |
| `category` | `?category=electronics` | Filter by category |
| `min_price` | `?min_price=10` | Minimum price |
| `max_price` | `?max_price=100` | Maximum price |

**Sorting** uses a `sort` parameter:

| Value | Effect |
|---|---|
| `price_asc` | Sort by price, low to high |
| `price_desc` | Sort by price, high to low |
| `name` | Sort alphabetically by name |

## Under the hood

Offset pagination translates to SQL like:

```sql
SELECT * FROM products LIMIT 20 OFFSET 40;
```

Cursor pagination translates to:

```sql
SELECT * FROM products WHERE id > 42 ORDER BY id LIMIT 20;
```

Cursor is faster on large datasets because it can use the index directly, while `OFFSET` scans and skips rows. However, cursor pagination doesn't support random page access (you can't jump to page 5).

In Go, you parse query parameters, apply filters by iterating the slice (or by building SQL queries), apply sorting, and slice the result based on offset/limit.

## How Go uses it

The Go community uses both pagination strategies. `chi` and `gin` have pagination middleware/helpers. The standard library's `net/url` package provides `Values.Get()` for parsing query parameters.

Most Go ORMs (GORM, Ent) support `.Offset()` and `.Limit()` for pagination and `.Where()` for filtering.

## Go example

```go
package main

import (
	"encoding/json"
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

type PaginatedResponse struct {
	Data       []Product `json:"data"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	Total      int       `json:"total"`
	TotalPages int       `json:"total_pages"`
}

var products = []Product{
	{ID: 1, Name: "Laptop", Category: "electronics", Price: 999.99},
	{ID: 2, Name: "Mouse", Category: "electronics", Price: 29.99},
	// ...more products
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Filter
	category := strings.ToLower(r.URL.Query().Get("category"))
	filtered := products
	if category != "" {
		var f []Product
		for _, p := range products {
			if strings.ToLower(p.Category) == category {
				f = append(f, p)
			}
		}
		filtered = f
	}

	// Pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	total := len(filtered)
	totalPages := (total + limit - 1) / limit
	start := (page - 1) * limit
	if start >= total {
		json.NewEncoder(w).Encode(PaginatedResponse{
			Data: []Product{}, Page: page, Limit: limit,
			Total: total, TotalPages: totalPages,
		})
		return
	}
	end := start + limit
	if end > total {
		end = total
	}

	json.NewEncoder(w).Encode(PaginatedResponse{
		Data: filtered[start:end], Page: page, Limit: limit,
		Total: total, TotalPages: totalPages,
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/products", productsHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## Step-by-step execution

Client requests `GET /products?category=electronics&page=2&limit=2`:

1. `productsHandler` parses query params: `category=electronics`, `page=2`, `limit=2`.
2. Filtering: iterates all products, keeps only `category == "electronics"`. Result: 5 products.
3. Pagination: `page=2`, `limit=2`. `start = (2-1)*2 = 2`, `end = 4`.
4. Slices `filtered[2:4]` -- returns the 3rd and 4th electronic products.
5. Response: `{"data": [...], "page": 2, "limit": 2, "total": 5, "total_pages": 3}`.

For cursor pagination, client requests `GET /products/cursor?limit=3`:

1. No cursor in query, so `cursor = 0` (start from beginning).
2. Filtering applied. For each product, check `product.ID > 0` -- all pass.
3. First 3 products returned. `next_cursor = 3`, `has_more = true`.
4. Client follows up: `GET /products/cursor?cursor=3&limit=3`.
5. Products with ID > 3 are found and returned.

## Common mistakes

- Not validating `page` and `limit`. Negative values can cause crashes. Always clamp to sensible defaults.
- Returning the full dataset without pagination. This is the most common production pagination bug.
- Using offset pagination for real-time feeds where items are constantly added. Cursor is better.
- Not returning pagination metadata (`total`, `total_pages`, `has_more`). Clients cannot build UI without it.
- Exposing raw database IDs as cursors (security concern). Use a hash or opaque token for public APIs.
- Filtering after pagination (pagination slice before filter). Apply filter first, then paginate.

## Debugging walkthrough

A client reports page 2 returns the same items as page 1:

```go
func badPagination(w http.ResponseWriter, r *http.Request) {
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    start := (page - 1) * limit
    end := start + limit
    // BUG: operating on original slice, not filtered
    json.NewEncoder(w).Encode(products[start:end])
}
```

**Symptom**: `GET /products?category=electronics&page=2&limit=2` returns the same as page 1 because filter is applied to the full list but pagination uses a different slice.

**Investigation**: The handler never applies the filter before pagination. It paginates the raw `products` slice.

**Fix**: Apply filter first, store in `filtered`, then paginate `filtered`.

## Production notes

- Set a maximum limit (e.g., 100) to prevent abuse. Return a 400 if exceeded.
- For cursor pagination, use a base64-encoded opaque token rather than exposing internal IDs.
- Include total count only on the first page or when explicitly requested (`?include_total=true`) to avoid expensive COUNT queries.
- Cache filtered results for popular queries (e.g., "all electronics sorted by price") using Redis or similar.
- Document pagination, filtering, and sorting parameters in your API docs.

## Performance implications

- Offset pagination becomes slow on large datasets (OFFSET 100000 scans 100k rows). Cursor pagination is O(log n).
- Filtering by iterating a Go slice is O(n). For in-memory datasets under 100k items, this is fine. For larger datasets, push filtering to the database.
- Sorting with bubble sort (as shown in the example) is O(n^2). Use `sort.Slice` in production.
- Include `Cache-Control` headers on paginated responses if data changes infrequently.

## Practice task

1. Create a list endpoint for `Order` with fields: `id`, `customer`, `total`, `status`, `created_at`.
2. Implement offset pagination with `page` and `limit` (max 50).
3. Add filtering by `status` (pending, shipped, delivered) and `min_total` / `max_total`.
4. Add sorting by `total_asc`, `total_desc`, `date_asc`, `date_desc`.
5. Implement cursor pagination as an alternative endpoint.
6. Write table-driven tests for: default pagination, custom page/limit, filters, sorting, and cursor traversal.

## Tests / verification

```bash
go test ./curriculum/modules/08-http-rest-apis/lessons/19-pagination-and-filtering -v
```

Tests cover parsing, filtering, sorting, offset pagination, cursor pagination edge cases (page beyond end, zero values), and method validation.

## Review questions

1. What is the difference between offset pagination and cursor pagination? When should you use each?
2. Why is cursor pagination often faster on large datasets than offset pagination?
3. What metadata should a paginated response include so clients can build page navigation?
4. Why should filtering be applied before pagination rather than after?
5. How would you handle a request with `limit=100000`? What's the maximum reasonable limit?

## NEXT UP

OpenAPI basics -- documenting your REST API with OpenAPI 3.0, defining paths, methods, and schemas, and generating Go code from the spec.
