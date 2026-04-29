// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package repository

import (
	"context"
	"database/sql"

	"github.com/rasel9t6/the-go-engineer/06-backend-db/01-web-and-database/databases/6-repository/models"
)

// UserRepository defines the contract for user data access.
// By depending on this interface, our business logic stays decoupled
// from the underlying database implementation.
type UserRepository interface {
	Create(ctx context.Context, name, email, password, avatar string) (int64, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	List(ctx context.Context) ([]models.User, error)
}

// SQLUserRepository implements the UserRepository interface using a SQL database.
type SQLUserRepository struct {
	db *sql.DB
}

// NewSQLUserRepository creates a new SQL-backed repository.
func NewSQLUserRepository(db *sql.DB) *SQLUserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) Create(ctx context.Context, name, email, password, avatar string) (int64, error) {
	// Start a transaction since we are writing to two tables
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// 1. Insert User
	res, err := tx.ExecContext(ctx, `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`, name, email, password)
	if err != nil {
		return 0, err
	}

	userID, _ := res.LastInsertId()

	// 2. Insert Profile
	_, err = tx.ExecContext(ctx, `INSERT INTO profiles (user_id, avatar) VALUES (?, ?)`, userID, avatar)
	if err != nil {
		return 0, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *SQLUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.password, u.created_at, p.avatar, p.created_at
		FROM users u
		LEFT JOIN profiles p ON u.id = p.user_id
		WHERE u.email = ?
	`
	var u models.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt,
		&u.Profile.Avatar, &u.Profile.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	u.Profile.UserID = u.ID
	return &u, nil
}

func (r *SQLUserRepository) List(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, name, email, created_at FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}
