package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestApp(t *testing.T) (*AppServer, *BookStore) {
	t.Helper()
	store := NewBookStore()
	server := NewAppServer(store, defaultStaticFS)
	return server, store
}

func TestBookCatalog_ListAndSearch(t *testing.T) {
	server, _ := setupTestApp(t)

	// 1. List all seeded books
	req := httptest.NewRequest(http.MethodGet, "/api/books", nil)
	rr := httptest.NewRecorder()
	server.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var books []Book
	if err := json.NewDecoder(rr.Body).Decode(&books); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}
	if len(books) != 4 {
		t.Errorf("expected 4 seeded books, got %d", len(books))
	}

	// 2. Query search
	reqSearch := httptest.NewRequest(http.MethodGet, "/api/books?q=concurrency", nil)
	rrSearch := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrSearch, reqSearch)

	if rrSearch.Code != http.StatusOK {
		t.Fatalf("expected status 200 for search query, got %d", rrSearch.Code)
	}

	var matched []Book
	_ = json.NewDecoder(rrSearch.Body).Decode(&matched)
	if len(matched) != 1 || matched[0].Title != "Concurrency in Go" {
		t.Errorf("expected 1 match for 'concurrency', got %v", matched)
	}
}

func TestBookCatalog_CreateValidation(t *testing.T) {
	server, _ := setupTestApp(t)

	// Missing title, bad ISBN, zero price
	invalidPayload := map[string]any{
		"title":  "",
		"author": "",
		"isbn":   "123", // invalid length
		"price":  -5.0,
	}
	body, _ := json.Marshal(invalidPayload)

	req := httptest.NewRequest(http.MethodPost, "/api/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	server.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got %d", rr.Code)
	}

	var apiErr APIError
	if err := json.NewDecoder(rr.Body).Decode(&apiErr); err != nil {
		t.Fatalf("failed to decode APIError: %v", err)
	}

	if len(apiErr.Details) != 4 {
		t.Errorf("expected 4 validation errors, got %d: %v", len(apiErr.Details), apiErr.Details)
	}
}

func TestBookCatalog_CreateSuccess(t *testing.T) {
	server, _ := setupTestApp(t)

	validPayload := map[string]any{
		"title":  "Building Microservices with Go",
		"author": "Nic Jackson",
		"isbn":   "9781786468666",
		"price":  34.99,
	}
	body, _ := json.Marshal(validPayload)

	req := httptest.NewRequest(http.MethodPost, "/api/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	server.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d: %s", rr.Code, rr.Body.String())
	}

	var created Book
	if err := json.NewDecoder(rr.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode created Book: %v", err)
	}

	if created.ID <= 0 || created.Title != "Building Microservices with Go" {
		t.Errorf("unexpected created book data: %+v", created)
	}
}

func TestBookCatalog_Delete(t *testing.T) {
	server, _ := setupTestApp(t)

	// Delete existing book ID 1
	req := httptest.NewRequest(http.MethodDelete, "/api/books/1", nil)
	rr := httptest.NewRecorder()
	server.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status 204 No Content, got %d", rr.Code)
	}

	// Try deleting again -> 404
	req2 := httptest.NewRequest(http.MethodDelete, "/api/books/1", nil)
	rr2 := httptest.NewRecorder()
	server.Handler().ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 Not Found on second delete, got %d", rr2.Code)
	}
}

func TestBookCatalog_StaticFrontend(t *testing.T) {
	server, _ := setupTestApp(t)

	// HTML index
	reqIndex := httptest.NewRequest(http.MethodGet, "/", nil)
	rrIndex := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrIndex, reqIndex)

	if rrIndex.Code != http.StatusOK {
		t.Fatalf("expected status 200 for root index, got %d", rrIndex.Code)
	}
	if !strings.Contains(rrIndex.Body.String(), "id=\"search-input\"") {
		t.Errorf("expected index.html to contain search input element")
	}

	// JS asset
	reqJS := httptest.NewRequest(http.MethodGet, "/app.js", nil)
	rrJS := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrJS, reqJS)

	if rrJS.Code != http.StatusOK {
		t.Fatalf("expected status 200 for /app.js, got %d", rrJS.Code)
	}
	if !strings.Contains(rrJS.Body.String(), "debounce") {
		t.Errorf("expected app.js to contain debounce function")
	}
}
