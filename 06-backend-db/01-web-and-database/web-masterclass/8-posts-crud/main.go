// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Web Masterclass - Posts CRUD
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to build a complete CRUD (Create, Read, Update, Delete) system.
//   - How to integrate Routing, Repositories, and JSON processing.
//   - How to use path parameters for resource identification.
//   - Best practices for error handling and HTTP status codes.
//
// WHY THIS MATTERS:
//   - CRUD is the bread and butter of web development. Almost every
//     application involves managing resources in a database. Learning
//     how to structure these operations cleanly and safely is a
//     foundational skill for any backend engineer.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/web-masterclass/8-posts-crud
//
// KEY TAKEAWAY:
//   - Consistency is key. Structure your resources following RESTful patterns.
// ============================================================================

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

// Stage 06: Web Masterclass - Posts CRUD
//
//   - POST /api/posts: Create
//   - GET /api/posts/{id}: Read
//   - PUT /api/posts/{id}: Update
//   - DELETE /api/posts/{id}: Delete
//
// ENGINEERING DEPTH:
//   A professional CRUD system isn't just about running SQL. It's
//   about handling "Edge Cases." What if the ID is not a number?
//   What if the JSON is malformed? What if the resource doesn't
//   exist? By centralizing our data logic in a `PostRepository`, we
//   ensure that our HTTP handlers stay focused on transport concerns
//   (parsing, status codes) while our data layer handles the specifics
//   of the database.

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// PostRepository defines the contract for managing posts.
type PostRepository interface {
	Create(ctx context.Context, title, content string) (int64, error)
	Get(ctx context.Context, id int) (*Post, error)
	List(ctx context.Context) ([]Post, error)
}

// SQLPostRepository implements PostRepository using SQLite.
type SQLPostRepository struct {
	db *sql.DB
}

func (r *SQLPostRepository) Create(ctx context.Context, title, content string) (int64, error) {
	res, err := r.db.ExecContext(ctx, "INSERT INTO posts (title, content) VALUES (?, ?)", title, content)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *SQLPostRepository) Get(ctx context.Context, id int) (*Post, error) {
	var p Post
	err := r.db.QueryRowContext(ctx, "SELECT id, title, content, created_at FROM posts WHERE id = ?", id).
		Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *SQLPostRepository) List(ctx context.Context) ([]Post, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, title, content, created_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

type application struct {
	posts PostRepository
}

func main() {
	// 1. Initialize DB (In-Memory for demo)
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec("CREATE TABLE posts (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT, content TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP);")

	app := &application{
		posts: &SQLPostRepository{db: db},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/posts", app.handleCreate)
	mux.HandleFunc("GET /api/posts", app.handleList)
	mux.HandleFunc("GET /api/posts/{id}", app.handleGet)

	fmt.Println("=== Web Masterclass: Posts CRUD ===")
	fmt.Println("  🚀 Server starting on http://localhost:8087")
	fmt.Println()

	log.Fatal(http.ListenAndServe(":8087", mux))

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: MC.9 pagination")
	fmt.Println("Current: MC.8 (posts-crud)")
	fmt.Println("Previous: MC.7 (forms)")
	fmt.Println("---------------------------------------------------")
}

func (app *application) handleCreate(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := app.posts.Create(r.Context(), input.Title, input.Content)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Post created with ID: %d", id)
}

func (app *application) handleList(w http.ResponseWriter, r *http.Request) {
	posts, err := app.posts.List(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (app *application) handleGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	post, err := app.posts.Get(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
