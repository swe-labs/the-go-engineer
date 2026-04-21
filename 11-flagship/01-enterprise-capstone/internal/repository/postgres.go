// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-enterprise-capstone/internal/models"
)

// ============================================================================
// Internal Package: Repository
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Isolating Database logic from HTTP logic (Dependency Inversion)
//   - Defining interfaces for clean mocking
//
// ENGINEERING DEPTH:
//   A 'Handler' (the HTTP layer) should never run raw SQL queries.
//   Instead, the Handler calls an Interface (the Repository).
//   This gives you the power to swap out PostgreSQL for MongoDB, or mock it perfectly
//   in your tests, without touching the HTTP layer at all!
// ============================================================================

// --- The Interface (The Contract) ---

type PostRepository interface {
	Create(ctx context.Context, post *models.Post) (int, error)
	GetAll(ctx context.Context) ([]models.Post, error)
}

// --- The Concrete Implementation (PostgreSQL) ---

type PostgresRepo struct {
	db *sql.DB
}

// NewPostgresRepo acts as an explicit constructor for Postgres functionality.
func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Create(ctx context.Context, p *models.Post) (int, error) {
	query := `
		INSERT INTO posts (title, content, author_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	// Postgres uses $1, $2, $3 for positional arguments, NOT `?` like SQLite!
	var newID int

	// Database Context Timeout Integration
	// By using `QueryRowContext(ctx, ...)`, the database driver becomes aware of
	// HTTP constraints. If the HTTP Request drops out or times out, the Context
	// signals the SQL Driver to instantly ABORT the query, saving database CPU overhead.
	err := r.db.QueryRowContext(ctx, query, p.Title, p.Content, p.AuthorID).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("repository.create_post: %w", err)
	}

	return newID, nil
}

func (r *PostgresRepo) GetAll(ctx context.Context) ([]models.Post, error) {
	query := `
		SELECT id, title, content, author_id, created_at 
		FROM posts 
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository.get_all_posts: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		// Pointers must perfectly match the columns selected above
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.AuthorID, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
