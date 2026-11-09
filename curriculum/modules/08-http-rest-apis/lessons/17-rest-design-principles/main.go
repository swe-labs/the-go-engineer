package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	ISBN   string `json:"isbn"`
}

type store struct {
	mu    sync.RWMutex
	books map[int]Book
	next  int
}

func newStore() *store {
	return &store{
		books: make(map[int]Book),
		next:  1,
	}
}

func (s *store) create(b Book) Book {
	s.mu.Lock()
	defer s.mu.Unlock()
	b.ID = s.next
	s.books[b.ID] = b
	s.next++
	return b
}

func (s *store) getAll() []Book {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Book, 0, len(s.books))
	for _, b := range s.books {
		result = append(result, b)
	}
	return result
}

func (s *store) getByID(id int) (Book, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.books[id]
	return b, ok
}

func (s *store) update(b Book) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.books[b.ID]; !ok {
		return false
	}
	s.books[b.ID] = b
	return true
}

func (s *store) delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.books[id]; !ok {
		return false
	}
	delete(s.books, id)
	return true
}

var bookStore = newStore()

func booksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		books := bookStore.getAll()
		json.NewEncoder(w).Encode(books)

	case http.MethodPost:
		var b Book
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}
		if b.Title == "" || b.Author == "" {
			http.Error(w, `{"error":"title and author are required"}`, http.StatusBadRequest)
			return
		}
		b = bookStore.create(b)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(b)

	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func bookByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/api/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		b, ok := bookStore.getByID(id)
		if !ok {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(b)

	case http.MethodPut:
		var b Book
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}
		b.ID = id
		if !bookStore.update(b) {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(b)

	case http.MethodDelete:
		if !bookStore.delete(id) {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func main() {
	bookStore.create(Book{Title: "The Go Programming Language", Author: "Donovan & Kernighan", ISBN: "978-0134190440"})
	bookStore.create(Book{Title: "Design Patterns", Author: "Gamma et al.", ISBN: "978-0201633610"})

	mux := http.NewServeMux()
	mux.HandleFunc("/api/books", booksHandler)
	mux.HandleFunc("/api/books/", bookByIDHandler)

	log.Println("RESTful books API on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func validateAuthor(author string) error {
	if strings.TrimSpace(author) == "" {
		return fmt.Errorf("author is required")
	}
	return nil
}
