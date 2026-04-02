// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// ============================================================================
// Section 12: Databases — SELECT Queries
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - QueryRow: fetching a single record
//   - Query: fetching multiple records
//   - rows.Scan: mapping SQL columns directly into Go struct fields
//   - Iterating over result sets with `rows.Next()`
//   - defer rows.Close(): the cardinal rule of database iteration
//
// ENGINEERING DEPTH:
//   When you run `db.Query()`, Go allocates a `sql.Rows` construct which
//   maintains an active, open TCP connection to the database holding the
//   cursor state. If you return early or forget to call `rows.Close()`,
//   that network connection stalls in the Connection Pool indefinitely. If
//   this happens often enough, you will exhaust your Database limits and crash
//   both your server and your Database.
//
// RUN: go run ./12-databases/3-select
// ============================================================================

type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"` // Omit password hash from JSON dumps!
	CreatedAt      time.Time `json:"created_at"`
}

func main() {
	dbName := "users_database.db"
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 1. Fetch ALL users
	fmt.Println("=== Fetching All Users ===")
	users, err := GetUsers(db)
	if err != nil {
		log.Fatal("Failed reading users:", err)
	}

	// Marshal to JSON to prettify the console output
	bs, _ := json.MarshalIndent(users, "", "  ")
	fmt.Println(string(bs))

	// 2. Fetch SINGLE user
	fmt.Println("\n=== Fetching Single User ===")
	// Note: this assumes "Alice" (from the previous exercise) exists in the database.
	// You may need to run `go run ./12-databases/2-query` first if the DB is empty.
	if len(users) > 0 {
		alice, err := GetUserByEmail(db, users[0].Email)
		if err != nil {
			log.Fatal("Failed to fetch user by email:", err)
		}
		fmt.Printf("✅ Found: ID=%d, Name=%s, Email=%s\n", alice.ID, alice.Name, alice.Email)
	} else {
		fmt.Println("⚠️ Database is empty. Run exercise `2-query` first to populate it.")
	}
}

// GetUserByEmail demonstrates reading a SINGLE row from the database.
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	stmt := `SELECT id, name, email, hashed_password, created_at FROM users WHERE email = ?`

	// QueryRow asks for exactly one row.
	row := db.QueryRow(stmt, email)

	var user User
	// `Scan` copies the columns returned from the DB into the memory addresses (&)
	// of our struct's fields. The columns MUST match the order sequentially.
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetUsers demonstrates reading MULTIPLE rows from the database.
func GetUsers(db *sql.DB) ([]User, error) {
	stmt := `SELECT id, name, email, hashed_password, created_at FROM users`

	// db.Query returns a `*sql.Rows` iterator.
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}

	// CRITICAL: Always defer rows.Close() to release the database connection connection!
	defer rows.Close()

	var users []User

	// rows.Next() moves the cursor forward one row.
	// It returns false when there are no more rows.
	for rows.Next() {
		var user User
		// Scan columns into struct pointers
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.HashedPassword,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}
		// Append populated struct to slice
		users = append(users, user)
	}

	// After the loop finishes, ALWAYS check rows.Err().
	// It's possible the network dropped mid-stream, breaking the loop early!
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
