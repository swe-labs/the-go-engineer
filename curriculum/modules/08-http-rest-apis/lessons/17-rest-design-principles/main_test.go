package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func resetStore() {
	bookStore = newStore()
}

func TestBooksHandlerGetAll(t *testing.T) {
	resetStore()
	bookStore.create(Book{Title: "Test Book", Author: "Test Author", ISBN: "123"})

	req := httptest.NewRequest(http.MethodGet, "/api/books", nil)
	rec := httptest.NewRecorder()
	booksHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var books []Book
	if err := json.NewDecoder(resp.Body).Decode(&books); err != nil {
		t.Fatalf("json decode: %v", err)
	}
	if len(books) != 1 {
		t.Errorf("expected 1 book, got %d", len(books))
	}
}

func TestBooksHandlerCreate(t *testing.T) {
	resetStore()

	body := `{"title":"New Book","author":"New Author","isbn":"456"}`
	req := httptest.NewRequest(http.MethodPost, "/api/books", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	booksHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	var b Book
	json.NewDecoder(resp.Body).Decode(&b)
	if b.Title != "New Book" {
		t.Errorf("title = %q, want %q", b.Title, "New Book")
	}
	if b.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestBooksHandlerCreateValidation(t *testing.T) {
	resetStore()

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"missing title", `{"author":"A"}`, http.StatusBadRequest},
		{"missing author", `{"title":"T"}`, http.StatusBadRequest},
		{"empty body", `{}`, http.StatusBadRequest},
		{"invalid json", `not json`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/books", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			booksHandler(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

func TestBooksHandlerMethodNotAllowed(t *testing.T) {
	resetStore()

	req := httptest.NewRequest(http.MethodDelete, "/api/books", nil)
	rec := httptest.NewRecorder()
	booksHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

func TestBookByIDHandlerGet(t *testing.T) {
	resetStore()
	created := bookStore.create(Book{Title: "Test", Author: "Author", ISBN: "789"})

	req := httptest.NewRequest(http.MethodGet, "/api/books/"+fmt.Sprint(created.ID), nil)
	rec := httptest.NewRecorder()
	bookByIDHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var b Book
	json.NewDecoder(resp.Body).Decode(&b)
	if b.Title != "Test" {
		t.Errorf("title = %q, want %q", b.Title, "Test")
	}
}

func TestBookByIDHandlerNotFound(t *testing.T) {
	resetStore()

	req := httptest.NewRequest(http.MethodGet, "/api/books/999", nil)
	rec := httptest.NewRecorder()
	bookByIDHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestBookByIDHandlerDelete(t *testing.T) {
	resetStore()
	created := bookStore.create(Book{Title: "Delete Me", Author: "Author", ISBN: "000"})

	req := httptest.NewRequest(http.MethodDelete, "/api/books/"+fmt.Sprint(created.ID), nil)
	rec := httptest.NewRecorder()
	bookByIDHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", resp.StatusCode)
	}
}

func TestBookByIDHandlerUpdate(t *testing.T) {
	resetStore()
	created := bookStore.create(Book{Title: "Original", Author: "Author", ISBN: "111"})

	body := `{"title":"Updated","author":"Author","isbn":"222"}`
	req := httptest.NewRequest(http.MethodPut, "/api/books/"+fmt.Sprint(created.ID), strings.NewReader(body))
	rec := httptest.NewRecorder()
	bookByIDHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var b Book
	json.NewDecoder(resp.Body).Decode(&b)
	if b.Title != "Updated" {
		t.Errorf("title = %q, want %q", b.Title, "Updated")
	}
}

func TestValidateAuthor(t *testing.T) {
	if err := validateAuthor("Valid"); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if err := validateAuthor(""); err == nil {
		t.Error("expected error for empty author")
	}
	if err := validateAuthor("  "); err == nil {
		t.Error("expected error for whitespace author")
	}
}

func TestLessonCompiles(t *testing.T) {
}
