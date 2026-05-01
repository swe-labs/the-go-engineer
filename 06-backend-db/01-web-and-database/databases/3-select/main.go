// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Reading Data (SELECT)
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to use 'db.Query' to fetch multiple rows.
//   - How to use 'db.QueryRow' for single record retrieval.
//   - The importance of 'rows.Scan' and 'rows.Close'.
//   - How to handle 'sql.ErrNoRows'.
//
// WHY THIS MATTERS:
//   - Reading data is the most common database operation. Understanding
//     how to safely map database columns to Go structs while managing
//     resources (like open cursors) is a fundamental backend skill.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/databases/3-select
//
// KEY TAKEAWAY:
//   - Always defer 'rows.Close()' and check 'rows.Err()' after iteration.
// ============================================================================

package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

// Stage 06: Databases - SELECT Queries
//
//   - Query: Fetching multiple records
//   - QueryRow: Fetching a single record
//   - Scan: Mapping SQL columns to Go fields
//   - rows.Next(): Iterating over results
//
// ENGINEERING DEPTH:
//   When you call `db.Query()`, Go keeps a network connection open and
//   a cursor active on the database server. If you forget to call
//   `rows.Close()`, that connection is "leaked" and cannot be used
//   by other parts of your app. In a high-traffic system, this can
//   exhaust the connection pool in seconds, causing your app to hang.

type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

func main() {
	dbPath := "example.db"

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("=== Reading Data (SELECT) ===")
	fmt.Println()

	// 1. Setup: Ensure table and data exist
	setup(db)

	// 2. Fetching multiple rows
	fmt.Println("  Fetching all users...")
	users, err := getAllUsers(db)
	if err != nil {
		log.Fatalf("  Error fetching users: %v", err)
	}

	for _, u := range users {
		fmt.Printf("  - [ID: %d] %-15s (%s)\n", u.ID, u.Name, u.Email)
	}

	// 3. Fetching a single row
	fmt.Println("\n  Fetching user with ID 1...")
	user, err := getUserByID(db, 1)
	if err != nil {
		log.Fatalf("  Error fetching user: %v", err)
	}
	fmt.Printf("  ✔ Found: %s\n", user.Name)

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: DB.4 -> 06-backend-db/01-web-and-database/databases/4-prepare")
	fmt.Println("Current: DB.3 (select-queries)")
	fmt.Println("Previous: DB.2 (query-insert)")
	fmt.Println("---------------------------------------------------")
}

func getAllUsers(db *sql.DB) ([]User, error) {
	query := `SELECT id, name, email, created_at FROM users`

	// db.Query returns an iterator (*sql.Rows)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	// CRITICAL: Always defer closing the rows!
	defer rows.Close()

	var users []User
	// Iterate through the result set
	for rows.Next() {
		var u User
		// Scan copies columns into Go variables.
		// Order MUST match the SELECT statement.
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	// IMPORTANT: Check for errors that happened DURING iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func getUserByID(db *sql.DB, id int) (*User, error) {
	query := `SELECT id, name, email, created_at FROM users WHERE id = ?`

	// QueryRow is used when you expect exactly one result.
	row := db.QueryRow(query, id)

	var u User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &u, nil
}

func setup(db *sql.DB) {
	schema := `CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, email TEXT UNIQUE, created_at DATETIME DEFAULT CURRENT_TIMESTAMP);`
	db.Exec(schema)

	// Add dummy data if empty
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count == 0 {
		db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", "Alice", "alice@example.com")
		db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", "Bob", "bob@example.com")
	}
}
