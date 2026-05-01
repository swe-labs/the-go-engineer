// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: REST API Exercise
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to combine all previous lessons into a complete, functional REST API.
//   - Best practices for CRUD (Create, Read, Update, Delete) operations.
//   - How to manage state safely in a concurrent HTTP environment.
//
// WHY THIS MATTERS:
//   - This is the cumulative proof of your skills. Building a complete
//     API shows that you understand how transport, logic, and state
//     work together in Go.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/http-servers/10-rest-api-exercise
//
// KEY TAKEAWAY:
//   - A great API is the result of consistent patterns applied to every endpoint.
// ============================================================================

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// --- Domain Model ---

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// --- In-Memory Database ---

type TaskStore struct {
	mu     sync.RWMutex
	tasks  map[int]Task
	nextID int
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks:  make(map[int]Task),
		nextID: 1,
	}
}

func (s *TaskStore) Create(title string) Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := Task{
		ID:        s.nextID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}
	s.tasks[task.ID] = task
	s.nextID++
	return task
}

func (s *TaskStore) List() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		list = append(list, t)
	}
	return list
}

func (s *TaskStore) Get(id int) (Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	return t, ok
}

func (s *TaskStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.tasks[id]
	if ok {
		delete(s.tasks, id)
	}
	return ok
}

// --- Middleware ---

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("  %s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// --- Handlers ---

type TaskAPI struct {
	store *TaskStore
}

func (api *TaskAPI) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /tasks", api.listTasks)
	mux.HandleFunc("POST /tasks", api.createTask)
	mux.HandleFunc("GET /tasks/{id}", api.getTask)
	mux.HandleFunc("DELETE /tasks/{id}", api.deleteTask)
}

func (api *TaskAPI) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks := api.store.List()
	respondJSON(w, http.StatusOK, tasks)
}

func (api *TaskAPI) createTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if input.Title == "" {
		respondError(w, http.StatusUnprocessableEntity, "Title is required")
		return
	}

	task := api.store.Create(input.Title)
	respondJSON(w, http.StatusCreated, task)
}

func (api *TaskAPI) getTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseTaskID(w, r)
	if !ok {
		return
	}

	task, ok := api.store.Get(id)
	if !ok {
		respondError(w, http.StatusNotFound, "Task not found")
		return
	}
	respondJSON(w, http.StatusOK, task)
}

func (api *TaskAPI) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseTaskID(w, r)
	if !ok {
		return
	}

	if ok := api.store.Delete(id); !ok {
		respondError(w, http.StatusNotFound, "Task not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Helpers ---

func parseTaskID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid task id")
		return 0, false
	}

	return id, true
}

func respondJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("  response encode error: %v", err)
	}
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

// --- Main ---

func main() {
	fmt.Println("=== Task Management REST API ===")
	fmt.Println()

	store := NewTaskStore()
	api := &TaskAPI{store: store}
	mux := http.NewServeMux()

	api.RegisterRoutes(mux)

	// Wrap mux in logging middleware
	handler := Logger(mux)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Graceful shutdown setup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Println("  Server starting on http://localhost:8080")
		fmt.Println("  Use curl to interact with the API:")
		fmt.Println("    - GET /tasks")
		fmt.Println("    - POST /tasks -d '{\"title\": \"Learn Go\"}'")
		fmt.Println()

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("  Fatal: %v", err)
		}
	}()

	<-stop
	fmt.Println("\n  Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("  Shutdown error: %v\n", err)
	}

	fmt.Println("  Server stopped.")

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: API.1 -> 06-backend-db/01-web-and-database/apis/1-rest-design-principles")
	fmt.Println("Current: HS.10 (rest-api-exercise)")
	fmt.Println("Previous: HS.9 (health-and-readiness-probes)")
	fmt.Println("---------------------------------------------------")
}
