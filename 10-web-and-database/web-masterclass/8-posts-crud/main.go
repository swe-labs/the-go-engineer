// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// ============================================================================
// Section 13: Posts CRUD with Repository Pattern
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Repository pattern for data access abstraction
//   - Interface-based design for testability
//   - SQLite database operations (CREATE, INSERT, SELECT)
//   - Context propagation for database operations
//   - Clean separation of models, repository, and handlers
//
// ENGINEERING DEPTH:
//   Why pass `context.Context` from the HTTP Handler directly into `ExecContext`?
//   If a client makes a `GET /posts` request that triggers a slow 10-second DB
//   scan, but the user immediately closes their browser tab, the HTTP server
//   detects the dropped TCP connection and cancels the `Request.Context()`.
//   Because we propagate this `ctx` into the `sql` package, the Go runtime instantly
//   fires a cancellation signal to the database engine via the socket, killing the
//   heavy query mid-flight! This prevents massive resource exhaustion attacks.
//
// RUN: go run ./13-web-masterclass/8-posts-crud
// ============================================================================

// Post represents a blog post.
type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  int       `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PostRepository defines the interface for post data access.
// Using an interface allows swapping implementations (SQLite → PostgreSQL)
// and makes testing easy with mock implementations.
type PostRepository interface {
	Create(ctx context.Context, post *Post) (int64, error)
	GetByID(ctx context.Context, id int) (*Post, error)
	GetAll(ctx context.Context, page, pageSize int) ([]Post, int, error)
}

// SQLitePostRepository implements PostRepository using SQLite.
type SQLitePostRepository struct {
	db *sql.DB
}

// NewSQLitePostRepository creates a new repository with dependency injection.
// The *sql.DB is injected — the repository doesn't create its own connection.
func NewSQLitePostRepository(db *sql.DB) *SQLitePostRepository {
	return &SQLitePostRepository{db: db}
}

// Create inserts a new post and returns the generated ID.
func (r *SQLitePostRepository) Create(ctx context.Context, post *Post) (int64, error) {
	query := `INSERT INTO posts (title, content, author_id, created_at, updated_at)
	           VALUES (?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		post.Title, post.Content, post.AuthorID, now, now)
	if err != nil {
		return 0, fmt.Errorf("create post: %w", err) // %w wraps the error for errors.Is/As
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id: %w", err)
	}
	return id, nil
}

// GetByID retrieves a single post by its ID.
func (r *SQLitePostRepository) GetByID(ctx context.Context, id int) (*Post, error) {
	query := `SELECT id, title, content, author_id, created_at, updated_at
	           FROM posts WHERE id = ?`

	var post Post
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID, &post.Title, &post.Content,
		&post.AuthorID, &post.CreatedAt, &post.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("post %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("get post %d: %w", id, err)
	}
	return &post, nil
}

// GetAll retrieves paginated posts and the total count.
// Returns (posts, totalCount, error) for pagination metadata.
func (r *SQLitePostRepository) GetAll(ctx context.Context, page, pageSize int) ([]Post, int, error) {
	// Get total count for pagination
	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM posts").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count posts: %w", err)
	}

	// Fetch paginated results
	offset := (page - 1) * pageSize
	query := `SELECT id, title, content, author_id, created_at, updated_at
	           FROM posts ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list posts: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content,
			&p.AuthorID, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("scan post: %w", err)
		}
		posts = append(posts, p)
	}

	return posts, total, rows.Err()
}

// Application wires handlers with their dependencies
type application struct {
	posts PostRepository
}

func main() {
	// Initialize SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create posts table
	_, err = db.Exec(`CREATE TABLE posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		author_id INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// Inject repository into application
	app := &application{
		posts: NewSQLitePostRepository(db),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/posts", app.handleCreatePost)
	mux.HandleFunc("GET /api/posts/{id}", app.handleGetPost)
	mux.HandleFunc("GET /api/posts", app.handleListPosts)

	fmt.Println("🚀 Posts CRUD server on http://localhost:8080")
	fmt.Println("   POST /api/posts        — Create a post (JSON body)")
	fmt.Println("   GET  /api/posts/{id}   — Get a post by ID")
	fmt.Println("   GET  /api/posts?page=1 — List posts (paginated)")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func (app *application) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := app.posts.Create(r.Context(), &post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"id": id, "message": "Post created"})
}

func (app *application) handleGetPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	post, err := app.posts.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (app *application) handleListPosts(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize := 10

	posts, total, err := app.posts.GetAll(r.Context(), page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"posts": posts,
		"total": total,
		"page":  page,
	})
}
