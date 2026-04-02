// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package repository

import (
	"context"
	"database/sql"

	"github.com/rasel9t6/the-go-engineer/12-databases/6-repository/models"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository is the Contract.
// High-level business logic depends ONLY on this interface.
// It operates exclusively on domain models (models.User), completely hiding
// the underlying database queries and SQL drivers.
type UserRepository interface {
	CreateUser(name, email, hashedPassword, avatar string) (int64, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUsers() ([]models.User, error)
}

// SQLUserRepository is the Concrete Implementation.
// It wraps a physical *sql.DB connection pool.
type SQLUserRepository struct {
	db *sql.DB
}

// NewSQLUserRepository is the Constructor.
// Note that it returns the `UserRepository` interface, not the concrete struct!
// This guarantees that callers can't accidentally use `sql.DB` directly.
func NewSQLUserRepository(db *sql.DB) UserRepository {
	return &SQLUserRepository{
		db: db,
	}
}

// CreateUser executes an ACID transaction to write to both users and profiles tables.
func (r *SQLUserRepository) CreateUser(name, email, rawPassword, avatar string) (int64, error) {
	ctx := context.Background()

	// 1. Begin Transaction
	// `BeginTx` acquires a dedicated connection from the `*sql.DB` connection pool
	// and locks it for our exclusive use. All following queries on `tx` will execute
	// within this single ACID transaction.
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	// 2. The Safety Catch
	// If the function returns early due to an error, `Rollback()` will abort the
	// transaction. If `Commit()` succeeds later, calling `Rollback()` does nothing!
	defer tx.Rollback()

	// 2. Hash Password
	hp, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// 3. Prevent SQL Injection
	// We NEVER use `fmt.Sprintf` to build queries!
	// Using `?` parameters forces the database driver to treat the inputs strictly
	// as literal bytes rather than executable SQL commands, completely neutralizing
	// SQL injection attacks.
	userQuery := `INSERT INTO users (name, email, hashed_password) VALUES (?, ?, ?)`
	result, err := tx.ExecContext(ctx, userQuery, name, email, string(hp))
	if err != nil {
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// 4. Insert Profile
	profileQuery := `INSERT INTO profiles (user_id, avatar) VALUES( ?, ?)`
	_, err = tx.ExecContext(ctx, profileQuery, userID, avatar)
	if err != nil {
		return 0, err // Rollback handles cleanup
	}

	// 5. Commit
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return userID, nil
}

// GetUserByEmail joins the underlying tables and maps them to the single User domain model.
func (r *SQLUserRepository) GetUserByEmail(email string) (*models.User, error) {
	// JOIN query to combine user data and profile data
	stmt := `
		SELECT u.id, u.name, u.email, u.hashed_password, u.created_at, p.avatar 
		FROM users u 
		INNER JOIN profiles p ON u.id = p.user_id 
		WHERE u.email = ?
	`

	row := r.db.QueryRow(stmt, email)

	// 4. Memory Pointer Scanning
	// `.Scan()` reads the raw bytes returned by the database driver and
	// casts them into the specific Go types, writing the data directly into
	// the memory addresses we provide (`&`).
	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.HashedPassword,
		&user.CreatedAt,
		&user.Profile.Avatar,
	)

	if err != nil {
		return nil, err
	}
	user.Profile.UserID = user.ID

	return &user, nil
}

// GetUsers returns all users from the database.
func (r *SQLUserRepository) GetUsers() ([]models.User, error) {
	stmt := `SELECT id, name, email, hashed_password, created_at FROM users`
	rows, err := r.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // ALWAYS close your rows

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.HashedPassword,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
