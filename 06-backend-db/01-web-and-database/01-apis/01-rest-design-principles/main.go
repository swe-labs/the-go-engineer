// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: REST Design Principles
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The philosophy of Representational State Transfer (REST).
//   - How to use HTTP verbs (GET, POST, PUT, DELETE) as actions.
//   - The discipline of resource-oriented naming.
//
// WHY THIS MATTERS:
//   - REST isn't a framework; it's a set of constraints that makes APIs
//     predictable, scalable, and easy to integrate for humans and machines.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/01-apis/01-rest-design-principles
//
// KEY TAKEAWAY:
//   - Resources are Nouns. HTTP Verbs are Actions.
// ============================================================================

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// User (Struct): represents a resource in this REST API demonstration.
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// userStore (Struct): in-memory store demonstrating resource state management.
type userStore struct {
	mu    sync.RWMutex
	users map[int]User
	next  int
}

// newUserStore (Constructor): initializes a new in-memory user store with auto-incrementing IDs.
func newUserStore() *userStore {
	return &userStore{users: make(map[int]User), next: 1}
}

// userStore.create (Method): inserts a new User and returns it with an assigned ID.
func (s *userStore) create(name, email string) User {
	s.mu.Lock()
	defer s.mu.Unlock()
	u := User{ID: s.next, Name: name, Email: email, CreatedAt: time.Now()}
	s.users[u.ID] = u
	s.next++
	return u
}

// userStore.list (Method): returns a snapshot of all stored users.
func (s *userStore) list() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]User, 0, len(s.users))
	for _, u := range s.users {
		result = append(result, u)
	}
	return result
}

// userStore.get (Method): retrieves a user by ID; the bool indicates whether the user was found.
func (s *userStore) get(id int) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	return u, ok
}

// userStore.delete (Method): removes a user by ID and reports whether the user existed.
func (s *userStore) delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.users[id]
	delete(s.users, id)
	return ok
}

// writeJSON (Function): serializes v as JSON and writes it with the given HTTP status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func main() {
	store := newUserStore()

	// Seed with sample data
	store.create("Alice", "alice@example.com")
	store.create("Bob", "bob@example.com")

	mux := http.NewServeMux()

	// RESTful endpoints demonstrating resource-oriented design
	mux.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		// GET /users -> List all users (200 OK)
		writeJSON(w, http.StatusOK, store.list())
	})

	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		// POST /users -> Create a user (201 Created)
		var input struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
			return
		}
		user := store.create(input.Name, input.Email)
		writeJSON(w, http.StatusCreated, user)
	})

	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		// GET /users/{id} -> Get one user (200 OK) or 404
		id, _ := strconv.Atoi(r.PathValue("id"))
		user, ok := store.get(id)
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
			return
		}
		writeJSON(w, http.StatusOK, user)
	})

	mux.HandleFunc("DELETE /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		// DELETE /users/{id} -> Delete a user (200 OK) or 404
		id, _ := strconv.Atoi(r.PathValue("id"))
		if !store.delete(id) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	})

	fmt.Println("=== REST Design Principles ===")
	fmt.Println()
	fmt.Println("  1. RESOURCE NAMING (Nouns)")
	fmt.Println("     GET  /users        -> List all users (Collection)")
	fmt.Println("     POST /users        -> Create a new user")
	fmt.Println("     GET  /users/{id}   -> Get one user by ID")
	fmt.Println("     DEL  /users/{id}   -> Remove a user")
	fmt.Println()
	fmt.Println("  2. HTTP VERBS AS ACTIONS")
	fmt.Println("     Resources are Nouns. HTTP Verbs are Actions.")
	fmt.Println("     Uses standard HTTP status codes: 200, 201, 400, 404")
	fmt.Println("     Uses consistent JSON payload structures")
	fmt.Println()
	fmt.Println("  3. Try these curl commands in another terminal:")
	fmt.Println("     curl http://localhost:8081/users")
	fmt.Println("     curl http://localhost:8081/users/1")
	fmt.Println("     curl -X POST -H 'Content-Type: application/json' -d '{\"name\":\"Charlie\",\"email\":\"charlie@example.com\"}' http://localhost:8081/users")
	fmt.Println("     curl -X DELETE http://localhost:8081/users/1")
	fmt.Println()
	fmt.Println("    Server running on :8081. Press Ctrl+C to stop.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: API.2 -> 06-backend-db/01-web-and-database/01-apis/02-api-versioning-strategies")
	fmt.Println("Current: API.1 (rest-design-principles)")
	fmt.Println("---------------------------------------------------")

	log.Fatal(http.ListenAndServe(":8081", mux))
}
