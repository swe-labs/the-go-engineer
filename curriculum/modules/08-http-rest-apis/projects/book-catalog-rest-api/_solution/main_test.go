package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestServer() (http.Handler, *BookStore) {
	store := NewBookStore()
	handler := NewAPIHandler(store)
	return handler, store
}

func TestHealthz(t *testing.T) {
	handler, _ := setupTestServer()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp JSONResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	data, ok := resp.Data.(map[string]any)
	if !ok || data["status"] != "ok" {
		t.Fatalf("unexpected data: %v", resp.Data)
	}
}

func TestCreateBook_Success(t *testing.T) {
	handler, _ := setupTestServer()
	payload := CreateBookRequest{
		Title:      "The Go Programming Language",
		Author:     "Alan Donovan",
		ISBN:       "978-0134190440",
		PriceCents: 4500,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp JSONResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse json: %v", err)
	}

	rawBook, _ := json.Marshal(resp.Data)
	var book Book
	_ = json.Unmarshal(rawBook, &book)

	if book.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if book.Title != payload.Title {
		t.Fatalf("expected title %q, got %q", payload.Title, book.Title)
	}
	if book.PriceCents != payload.PriceCents {
		t.Fatalf("expected price %d, got %d", payload.PriceCents, book.PriceCents)
	}
}

func TestCreateBook_ValidationErrors(t *testing.T) {
	handler, _ := setupTestServer()

	testCases := []struct {
		name    string
		payload CreateBookRequest
	}{
		{
			name: "empty title",
			payload: CreateBookRequest{
				Title:      "",
				Author:     "Rob Pike",
				ISBN:       "0123456789",
				PriceCents: 3000,
			},
		},
		{
			name: "empty author",
			payload: CreateBookRequest{
				Title:      "Go Concurrency",
				Author:     "",
				ISBN:       "0123456789",
				PriceCents: 3000,
			},
		},
		{
			name: "invalid ISBN",
			payload: CreateBookRequest{
				Title:      "Go Concurrency",
				Author:     "Katherine Cox-Buday",
				ISBN:       "invalid-isbn",
				PriceCents: 3000,
			},
		},
		{
			name: "negative or zero price",
			payload: CreateBookRequest{
				Title:      "Go Concurrency",
				Author:     "Katherine Cox-Buday",
				ISBN:       "0123456789",
				PriceCents: 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.payload)
			req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusUnprocessableEntity {
				t.Fatalf("expected 422 Unprocessable Entity, got %d", rec.Code)
			}

			var resp JSONResponse
			_ = json.Unmarshal(rec.Body.Bytes(), &resp)
			if resp.Error == nil || resp.Error.Code != "VALIDATION_FAILED" {
				t.Fatalf("expected VALIDATION_FAILED error code, got %v", resp.Error)
			}
		})
	}
}

func TestGetBook_FoundAndNotFound(t *testing.T) {
	handler, store := setupTestServer()
	book, _ := store.Create(t.Context(), CreateBookRequest{
		Title:      "Concurrency in Go",
		Author:     "Katherine Cox-Buday",
		ISBN:       "1491941197",
		PriceCents: 3800,
	})

	// Found
	req := httptest.NewRequest(http.MethodGet, "/books/"+book.ID, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	// Not Found
	reqNotFound := httptest.NewRequest(http.MethodGet, "/books/non-existent-id", nil)
	recNotFound := httptest.NewRecorder()
	handler.ServeHTTP(recNotFound, reqNotFound)

	if recNotFound.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", recNotFound.Code)
	}
}

func TestListBooks_PaginationAndFilter(t *testing.T) {
	handler, store := setupTestServer()
	ctx := t.Context()

	_, _ = store.Create(ctx, CreateBookRequest{Title: "Book A", Author: "John Doe", ISBN: "1111111111", PriceCents: 1000})
	_, _ = store.Create(ctx, CreateBookRequest{Title: "Book B", Author: "Jane Smith", ISBN: "2222222222", PriceCents: 2000})
	_, _ = store.Create(ctx, CreateBookRequest{Title: "Book C", Author: "John Doe", ISBN: "3333333333", PriceCents: 3000})

	// Filter by author John
	req := httptest.NewRequest(http.MethodGet, "/books?author=John", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp JSONResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)

	rawBooks, _ := json.Marshal(resp.Data)
	var books []Book
	_ = json.Unmarshal(rawBooks, &books)

	if len(books) != 2 {
		t.Fatalf("expected 2 books for John Doe, got %d", len(books))
	}

	// Pagination limit=1, offset=1
	reqPage := httptest.NewRequest(http.MethodGet, "/books?limit=1&offset=1", nil)
	recPage := httptest.NewRecorder()
	handler.ServeHTTP(recPage, reqPage)

	var respPage JSONResponse
	_ = json.Unmarshal(recPage.Body.Bytes(), &respPage)
	rawPageBooks, _ := json.Marshal(respPage.Data)
	var pageBooks []Book
	_ = json.Unmarshal(rawPageBooks, &pageBooks)

	if len(pageBooks) != 1 {
		t.Fatalf("expected 1 book on page 2, got %d", len(pageBooks))
	}
}

func TestUpdateAndDeleteBook(t *testing.T) {
	handler, store := setupTestServer()
	ctx := t.Context()

	book, _ := store.Create(ctx, CreateBookRequest{
		Title:      "Original Title",
		Author:     "Original Author",
		ISBN:       "1234567890",
		PriceCents: 1500,
	})

	// Update
	updatePayload := UpdateBookRequest{
		Title:      "Updated Title",
		Author:     "Updated Author",
		PriceCents: 1999,
	}
	body, _ := json.Marshal(updatePayload)
	reqUpdate := httptest.NewRequest(http.MethodPut, "/books/"+book.ID, bytes.NewReader(body))
	recUpdate := httptest.NewRecorder()
	handler.ServeHTTP(recUpdate, reqUpdate)

	if recUpdate.Code != http.StatusOK {
		t.Fatalf("expected 200 on update, got %d", recUpdate.Code)
	}

	// Delete
	reqDelete := httptest.NewRequest(http.MethodDelete, "/books/"+book.ID, nil)
	recDelete := httptest.NewRecorder()
	handler.ServeHTTP(recDelete, reqDelete)

	if recDelete.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content on delete, got %d", recDelete.Code)
	}

	// Verify it's gone
	reqGet := httptest.NewRequest(http.MethodGet, "/books/"+book.ID, nil)
	recGet := httptest.NewRecorder()
	handler.ServeHTTP(recGet, reqGet)

	if recGet.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after deletion, got %d", recGet.Code)
	}
}

func TestCorrelationIDMiddleware(t *testing.T) {
	handler, _ := setupTestServer()

	// Provided correlation ID
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("X-Correlation-ID", "custom-trace-id-12345")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Header().Get("X-Correlation-ID") != "custom-trace-id-12345" {
		t.Fatalf("expected X-Correlation-ID header 'custom-trace-id-12345', got %q", rec.Header().Get("X-Correlation-ID"))
	}

	// Auto-generated correlation ID
	reqAuto := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	recAuto := httptest.NewRecorder()
	handler.ServeHTTP(recAuto, reqAuto)

	if recAuto.Header().Get("X-Correlation-ID") == "" {
		t.Fatal("expected non-empty auto-generated X-Correlation-ID")
	}
}

func TestPanicRecoveryMiddleware(t *testing.T) {
	panickingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("simulated fatal crash in handler")
	})

	wrapped := RecoveryMiddleware(panickingHandler)
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 Internal Server Error, got %d", rec.Code)
	}

	body, _ := io.ReadAll(rec.Body)
	var resp JSONResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("failed to parse json response: %v", err)
	}

	if resp.Error == nil || resp.Error.Code != "INTERNAL_ERROR" {
		t.Fatalf("expected INTERNAL_ERROR code, got %v", resp.Error)
	}
}
