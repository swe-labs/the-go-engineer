package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
)

// Book represents a catalog item.
type Book struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Author     string    `json:"author"`
	ISBN       string    `json:"isbn"`
	PriceCents int       `json:"price_cents"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateBookRequest holds payload for new books.
type CreateBookRequest struct {
	Title      string `json:"title"`
	Author     string `json:"author"`
	ISBN       string `json:"isbn"`
	PriceCents int    `json:"price_cents"`
}

// UpdateBookRequest holds payload for updates.
type UpdateBookRequest struct {
	Title      string `json:"title"`
	Author     string `json:"author"`
	PriceCents int    `json:"price_cents"`
}

// JSONResponse defines the standardized API response structure.
type JSONResponse struct {
	Data  any       `json:"data,omitempty"`
	Error *APIError `json:"error,omitempty"`
	Meta  any       `json:"meta,omitempty"`
}

// APIError defines standard error envelope.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// BookStore provides thread-safe in-memory book storage.
type BookStore struct {
	mu    sync.RWMutex
	books map[string]Book
}

// NewBookStore initializes a new empty BookStore.
func NewBookStore() *BookStore {
	return &BookStore{
		books: make(map[string]Book),
	}
}

// GenerateID produces a random hexadecimal string identifier.
func GenerateID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// ValidateISBN validates ISBN-10 or ISBN-13 format.
func ValidateISBN(isbn string) bool {
	clean := strings.ReplaceAll(isbn, "-", "")
	match, _ := regexp.MatchString(`^(\d{10}|\d{13})$`, clean)
	return match
}

// Create adds a new book to the store.
func (s *BookStore) Create(ctx context.Context, req CreateBookRequest) (Book, error) {
	// TODO: Validate Title non-empty, Author non-empty, ISBN valid, PriceCents > 0
	// TODO: Generate ID, store Book thread-safely with CreatedAt and UpdatedAt
	return Book{}, errors.New("not implemented")
}

// Get finds a book by its ID.
func (s *BookStore) Get(ctx context.Context, id string) (Book, error) {
	// TODO: Lookup book by ID with RLock, return ErrNotFound if absent
	return Book{}, errors.New("not implemented")
}

// List returns a paginated slice of books, optionally filtered by author.
func (s *BookStore) List(ctx context.Context, author string, limit, offset int) ([]Book, int, error) {
	// TODO: Filter books matching author case-insensitively, apply pagination limit & offset
	return nil, 0, errors.New("not implemented")
}

// Update updates an existing book's details.
func (s *BookStore) Update(ctx context.Context, id string, req UpdateBookRequest) (Book, error) {
	// TODO: Validate inputs, update book under Lock, update UpdatedAt timestamp
	return Book{}, errors.New("not implemented")
}

// Delete removes a book from the store.
func (s *BookStore) Delete(ctx context.Context, id string) error {
	// TODO: Remove book under Lock, return ErrNotFound if absent
	return errors.New("not implemented")
}

type contextKey string

const correlationIDKey contextKey = "correlation_id"

func CorrelationIDFromContext(ctx context.Context) string {
	if val, ok := ctx.Value(correlationIDKey).(string); ok {
		return val
	}
	return ""
}

// Middleware: Correlation ID
func CorrelationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Extract or generate correlation ID, set in header and request context
		next.ServeHTTP(w, r)
	})
}

// Middleware: Recovery
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Recover from handler panic, log with slog, return 500 JSON envelope
		next.ServeHTTP(w, r)
	})
}

// Middleware: Logging
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		corrID := CorrelationIDFromContext(r.Context())
		slog.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start).String(),
			"correlation_id", corrID,
		)
	})
}

// APIHandler contains handlers and routes.
type APIHandler struct {
	store  *BookStore
	router *http.ServeMux
}

// NewAPIHandler configures endpoints.
func NewAPIHandler(store *BookStore) http.Handler {
	api := &APIHandler{
		store:  store,
		router: http.NewServeMux(),
	}

	api.router.HandleFunc("GET /healthz", api.handleHealthz)
	api.router.HandleFunc("POST /books", api.handleCreateBook)
	api.router.HandleFunc("GET /books", api.handleListBooks)
	api.router.HandleFunc("GET /books/{id}", api.handleGetBook)
	api.router.HandleFunc("PUT /books/{id}", api.handleUpdateBook)
	api.router.HandleFunc("DELETE /books/{id}", api.handleDeleteBook)

	return RecoveryMiddleware(CorrelationIDMiddleware(LoggingMiddleware(api.router)))
}

func (a *APIHandler) handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, JSONResponse{Data: map[string]string{"status": "ok"}})
}

func (a *APIHandler) handleCreateBook(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse json body, call store.Create, handle validation and return 201 Created
	writeJSON(w, http.StatusNotImplemented, JSONResponse{Error: &APIError{Code: "NOT_IMPLEMENTED"}})
}

func (a *APIHandler) handleListBooks(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse author, limit, offset query params, call store.List, return 200 with meta
	writeJSON(w, http.StatusNotImplemented, JSONResponse{Error: &APIError{Code: "NOT_IMPLEMENTED"}})
}

func (a *APIHandler) handleGetBook(w http.ResponseWriter, r *http.Request) {
	// TODO: Extract id path param, call store.Get, return 200 or 404
	writeJSON(w, http.StatusNotImplemented, JSONResponse{Error: &APIError{Code: "NOT_IMPLEMENTED"}})
}

func (a *APIHandler) handleUpdateBook(w http.ResponseWriter, r *http.Request) {
	// TODO: Extract id, decode body, call store.Update, return 200 or 404/422
	writeJSON(w, http.StatusNotImplemented, JSONResponse{Error: &APIError{Code: "NOT_IMPLEMENTED"}})
}

func (a *APIHandler) handleDeleteBook(w http.ResponseWriter, r *http.Request) {
	// TODO: Extract id, call store.Delete, return 204 No Content or 404
	writeJSON(w, http.StatusNotImplemented, JSONResponse{Error: &APIError{Code: "NOT_IMPLEMENTED"}})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func main() {
	store := NewBookStore()
	handler := NewAPIHandler(store)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("server starting", "port", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-shutdown
	slog.Info("server shutting down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}
	slog.Info("server exiting")
}
