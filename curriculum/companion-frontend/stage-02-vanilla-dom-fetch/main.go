package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:embed static/*
var defaultStaticFS embed.FS

// Book models a publication in the catalog.
type Book struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	ISBN      string    `json:"isbn"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

// ValidationError holds field-level error messages.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// APIError is the standardized JSON response for client errors.
type APIError struct {
	Error   string            `json:"error"`
	Details []ValidationError `json:"details,omitempty"`
}

// BookStore is a thread-safe in-memory repository for books.
type BookStore struct {
	mu     sync.RWMutex
	books  map[int]Book
	nextID int
}

// NewBookStore initializes a seeded BookStore.
func NewBookStore() *BookStore {
	bs := &BookStore{
		books:  make(map[int]Book),
		nextID: 1,
	}

	// Seed with classic Go engineering titles
	seeds := []Book{
		{Title: "The Go Programming Language", Author: "Alan Donovan & Brian Kernighan", ISBN: "9780134190440", Price: 39.99},
		{Title: "Concurrency in Go", Author: "Katherine Cox-Buday", ISBN: "9781491941195", Price: 44.99},
		{Title: "Cloud Native Go", Author: "Matthew Titmus", ISBN: "9781492076339", Price: 49.99},
		{Title: "100 Go Mistakes and How to Avoid Them", Author: "Teiva Harsanyi", ISBN: "9781617299599", Price: 45.00},
	}

	for _, s := range seeds {
		_ = bs.Create(s.Title, s.Author, s.ISBN, s.Price)
	}

	return bs
}

// List returns all books matching optional query filter.
func (bs *BookStore) List(query string) []Book {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	var result []Book
	q := strings.ToLower(strings.TrimSpace(query))

	for _, b := range bs.books {
		if q == "" || strings.Contains(strings.ToLower(b.Title), q) ||
			strings.Contains(strings.ToLower(b.Author), q) ||
			strings.Contains(b.ISBN, q) {
			result = append(result, b)
		}
	}
	return result
}

// GetByID finds a book by its integer ID.
func (bs *BookStore) GetByID(id int) (Book, bool) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	b, ok := bs.books[id]
	return b, ok
}

// Create validates and inserts a new book.
func (bs *BookStore) Create(title, author, isbn string, price float64) Book {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	b := Book{
		ID:        bs.nextID,
		Title:     strings.TrimSpace(title),
		Author:    strings.TrimSpace(author),
		ISBN:      strings.TrimSpace(isbn),
		Price:     price,
		CreatedAt: time.Now().UTC(),
	}
	bs.books[bs.nextID] = b
	bs.nextID++
	return b
}

// Delete removes a book by ID.
func (bs *BookStore) Delete(id int) bool {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if _, exists := bs.books[id]; !exists {
		return false
	}
	delete(bs.books, id)
	return true
}

// AppServer wraps the HTTP routing and dependencies.
type AppServer struct {
	store    *BookStore
	staticFS fs.FS
}

// NewAppServer constructs a configured AppServer.
func NewAppServer(store *BookStore, staticFS fs.FS) *AppServer {
	return &AppServer{
		store:    store,
		staticFS: staticFS,
	}
}

// Handler configures all API and static file routes.
func (s *AppServer) Handler() http.Handler {
	mux := http.NewServeMux()

	// REST API Endpoints
	mux.HandleFunc("/api/books", s.handleBooksCollection)
	mux.HandleFunc("/api/books/", s.handleBookItem)

	// Static Files Frontend
	staticSub, err := fs.Sub(s.staticFS, "static")
	if err == nil {
		fileServer := http.FileServer(http.FS(staticSub))
		mux.Handle("/", fileServer)
	}

	return s.loggingMiddleware(mux)
}

func (s *AppServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// handleBooksCollection processes GET /api/books and POST /api/books.
func (s *AppServer) handleBooksCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		query := r.URL.Query().Get("q")
		books := s.store.List(query)
		s.respondJSON(w, http.StatusOK, books)

	case http.MethodPost:
		var req struct {
			Title  string  `json:"title"`
			Author string  `json:"author"`
			ISBN   string  `json:"isbn"`
			Price  float64 `json:"price"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.respondError(w, http.StatusBadRequest, "Invalid JSON payload format", nil)
			return
		}

		// Comprehensive input validation
		var errs []ValidationError
		if strings.TrimSpace(req.Title) == "" {
			errs = append(errs, ValidationError{Field: "title", Message: "Title is required"})
		}
		if strings.TrimSpace(req.Author) == "" {
			errs = append(errs, ValidationError{Field: "author", Message: "Author name is required"})
		}

		isbnClean := strings.ReplaceAll(strings.ReplaceAll(req.ISBN, "-", ""), " ", "")
		isbnRegex := regexp.MustCompile(`^[0-9]{10}$|^[0-9]{13}$`)
		if !isbnRegex.MatchString(isbnClean) {
			errs = append(errs, ValidationError{Field: "isbn", Message: "ISBN must be 10 or 13 digits"})
		}

		if req.Price <= 0 {
			errs = append(errs, ValidationError{Field: "price", Message: "Price must be greater than $0.00"})
		}

		if len(errs) > 0 {
			s.respondError(w, http.StatusBadRequest, "Validation failed", errs)
			return
		}

		created := s.store.Create(req.Title, req.Author, isbnClean, req.Price)
		s.respondJSON(w, http.StatusCreated, created)

	default:
		w.Header().Set("Allow", "GET, POST")
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
	}
}

// handleBookItem processes GET /api/books/{id} and DELETE /api/books/{id}.
func (s *AppServer) handleBookItem(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid book identifier", nil)
		return
	}

	switch r.Method {
	case http.MethodGet:
		book, found := s.store.GetByID(id)
		if !found {
			s.respondError(w, http.StatusNotFound, fmt.Sprintf("Book with ID %d not found", id), nil)
			return
		}
		s.respondJSON(w, http.StatusOK, book)

	case http.MethodDelete:
		if !s.store.Delete(id) {
			s.respondError(w, http.StatusNotFound, fmt.Sprintf("Book with ID %d not found", id), nil)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		w.Header().Set("Allow", "GET, DELETE")
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
	}
}

func (s *AppServer) respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func (s *AppServer) respondError(w http.ResponseWriter, status int, message string, details []ValidationError) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(APIError{
		Error:   message,
		Details: details,
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	store := NewBookStore()
	server := NewAppServer(store, defaultStaticFS)

	addr := ":" + port
	log.Printf("[Stage 2 Starter] Book Catalog REST Server listening at http://localhost%s", addr)
	srv := &http.Server{
		Addr:         addr,
		Handler:      server.Handler(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server terminated: %v", err)
	}
}
