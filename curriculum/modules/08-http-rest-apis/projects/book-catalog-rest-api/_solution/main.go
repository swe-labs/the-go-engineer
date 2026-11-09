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

// ErrNotFound indicates a requested entity was not found.
var ErrNotFound = errors.New("not found")

// ErrInvalidInput indicates validation failure.
var ErrInvalidInput = errors.New("invalid input")

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
	if strings.TrimSpace(req.Title) == "" {
		return Book{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}
	if strings.TrimSpace(req.Author) == "" {
		return Book{}, fmt.Errorf("%w: author is required", ErrInvalidInput)
	}
	if !ValidateISBN(req.ISBN) {
		return Book{}, fmt.Errorf("%w: isbn must be 10 or 13 digits", ErrInvalidInput)
	}
	if req.PriceCents <= 0 {
		return Book{}, fmt.Errorf("%w: price_cents must be positive", ErrInvalidInput)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	b := Book{
		ID:         GenerateID(),
		Title:      strings.TrimSpace(req.Title),
		Author:     strings.TrimSpace(req.Author),
		ISBN:       strings.TrimSpace(req.ISBN),
		PriceCents: req.PriceCents,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	s.books[b.ID] = b
	return b, nil
}

// Get finds a book by its ID.
func (s *BookStore) Get(ctx context.Context, id string) (Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	b, ok := s.books[id]
	if !ok {
		return Book{}, ErrNotFound
	}
	return b, nil
}

// List returns a paginated slice of books, optionally filtered by author.
func (s *BookStore) List(ctx context.Context, author string, limit, offset int) ([]Book, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var matched []Book
	targetAuthor := strings.ToLower(strings.TrimSpace(author))

	for _, b := range s.books {
		if targetAuthor != "" && !strings.Contains(strings.ToLower(b.Author), targetAuthor) {
			continue
		}
		matched = append(matched, b)
	}

	total := len(matched)
	if offset < 0 {
		offset = 0
	}
	if offset >= total {
		return []Book{}, total, nil
	}

	end := offset + limit
	if limit <= 0 || end > total {
		end = total
	}

	return matched[offset:end], total, nil
}

// Update updates an existing book's details.
func (s *BookStore) Update(ctx context.Context, id string, req UpdateBookRequest) (Book, error) {
	if strings.TrimSpace(req.Title) == "" {
		return Book{}, fmt.Errorf("%w: title cannot be empty", ErrInvalidInput)
	}
	if strings.TrimSpace(req.Author) == "" {
		return Book{}, fmt.Errorf("%w: author cannot be empty", ErrInvalidInput)
	}
	if req.PriceCents <= 0 {
		return Book{}, fmt.Errorf("%w: price_cents must be positive", ErrInvalidInput)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	b, ok := s.books[id]
	if !ok {
		return Book{}, ErrNotFound
	}

	b.Title = strings.TrimSpace(req.Title)
	b.Author = strings.TrimSpace(req.Author)
	b.PriceCents = req.PriceCents
	b.UpdatedAt = time.Now().UTC()

	s.books[id] = b
	return b, nil
}

// Delete removes a book from the store.
func (s *BookStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.books[id]; !ok {
		return ErrNotFound
	}
	delete(s.books, id)
	return nil
}

// Context key for correlation ID
type contextKey string

const correlationIDKey contextKey = "correlation_id"

// CorrelationIDFromContext retrieves the correlation ID from context.
func CorrelationIDFromContext(ctx context.Context) string {
	if val, ok := ctx.Value(correlationIDKey).(string); ok {
		return val
	}
	return ""
}

// Middleware: Correlation ID
func CorrelationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corrID := r.Header.Get("X-Correlation-ID")
		if corrID == "" {
			corrID = GenerateID()
		}
		w.Header().Set("X-Correlation-ID", corrID)
		ctx := context.WithValue(r.Context(), correlationIDKey, corrID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Middleware: Recovery
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered in handler", "panic", rec, "path", r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(JSONResponse{
					Error: &APIError{
						Code:    "INTERNAL_ERROR",
						Message: "An internal server error occurred",
					},
				})
			}
		}()
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

	// Chain middlewares: Recovery -> CorrelationID -> Logging -> ServeMux
	return RecoveryMiddleware(CorrelationIDMiddleware(LoggingMiddleware(api.router)))
}

func (a *APIHandler) handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, JSONResponse{
		Data: map[string]string{"status": "ok"},
	})
}

func (a *APIHandler) handleCreateBook(w http.ResponseWriter, r *http.Request) {
	var req CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, JSONResponse{
			Error: &APIError{Code: "BAD_REQUEST", Message: "Malformed JSON payload"},
		})
		return
	}

	book, err := a.store.Create(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			writeJSON(w, http.StatusUnprocessableEntity, JSONResponse{
				Error: &APIError{Code: "VALIDATION_FAILED", Message: err.Error()},
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			Error: &APIError{Code: "INTERNAL_ERROR", Message: "Failed to create book"},
		})
		return
	}

	writeJSON(w, http.StatusCreated, JSONResponse{Data: book})
}

func (a *APIHandler) handleListBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	author := query.Get("author")
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset, _ := strconv.Atoi(query.Get("offset"))

	books, total, err := a.store.List(r.Context(), author, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			Error: &APIError{Code: "INTERNAL_ERROR", Message: "Failed to list books"},
		})
		return
	}

	writeJSON(w, http.StatusOK, JSONResponse{
		Data: books,
		Meta: map[string]int{
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	})
}

func (a *APIHandler) handleGetBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	book, err := a.store.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeJSON(w, http.StatusNotFound, JSONResponse{
				Error: &APIError{Code: "NOT_FOUND", Message: fmt.Sprintf("Book %s not found", id)},
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			Error: &APIError{Code: "INTERNAL_ERROR", Message: "Failed to retrieve book"},
		})
		return
	}

	writeJSON(w, http.StatusOK, JSONResponse{Data: book})
}

func (a *APIHandler) handleUpdateBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, JSONResponse{
			Error: &APIError{Code: "BAD_REQUEST", Message: "Malformed JSON payload"},
		})
		return
	}

	book, err := a.store.Update(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeJSON(w, http.StatusNotFound, JSONResponse{
				Error: &APIError{Code: "NOT_FOUND", Message: fmt.Sprintf("Book %s not found", id)},
			})
			return
		}
		if errors.Is(err, ErrInvalidInput) {
			writeJSON(w, http.StatusUnprocessableEntity, JSONResponse{
				Error: &APIError{Code: "VALIDATION_FAILED", Message: err.Error()},
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			Error: &APIError{Code: "INTERNAL_ERROR", Message: "Failed to update book"},
		})
		return
	}

	writeJSON(w, http.StatusOK, JSONResponse{Data: book})
}

func (a *APIHandler) handleDeleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := a.store.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeJSON(w, http.StatusNotFound, JSONResponse{
				Error: &APIError{Code: "NOT_FOUND", Message: fmt.Sprintf("Book %s not found", id)},
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			Error: &APIError{Code: "INTERNAL_ERROR", Message: "Failed to delete book"},
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
