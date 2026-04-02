// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// ============================================================================
// Section 12: Databases — Executing Queries (INSERT)
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Executing INSERT queries using db.Exec
//   - SQL Parameterization (`?`) to prevent SQL Injection
//   - Safely hashing passwords before storage
//   - Retrieving the LastInsertId
//
// ENGINEERING DEPTH:
//   NEVER use string concatenation (`"SELECT * FROM users WHERE name = " + user`)
//   to build SQL queries. This leaves you instantly vulnerable to SQL Injection.
//   Go's `database/sql` driver uses "Parameterized Queries" (the `?` symbol).
//   When you pass arguments to `db.Exec("query ?", arg)`, the Go driver sends the
//   query structure and the arguments to the database server SEPARATELY. The
//   database safely escapes all payloads internally before execution.
//
// RUN: go run ./12-databases/2-query
// ============================================================================

var schema = `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL, 
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

func main() {
	dbName := "users_database.db"
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ Database connection established")

	// 1. Create the table
	createTable(db)
	fmt.Println("✅ Users table ensured")

	// 2. Insert a secure user
	userId, err := createUser(db, "Alice", "alice@example.com", "supersecret123")
	if err != nil {
		log.Fatal("Failed to insert user:", err)
	}

	fmt.Printf("✅ Inserted user 'Alice' successfully! (ID: %d)\n", userId)
}

func createTable(db *sql.DB) {
	_, err := db.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}
}

// createUser demonstrates safe insertion using bcrypt and parameterized queries.
func createUser(db *sql.DB, name, email, rawPassword string) (int64, error) {

	// PREPARE THE DATA: Hash the password using bcrypt.
	// bcrypt.DefaultCost introduces intentional CPU delay to prevent brute-force attacks.
	hp, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("hashing password failed: %w", err)
	}

	// EXECUTE THE QUERY:
	// Notice the `?` symbols. These are positional placeholders.
	stmt := `INSERT INTO users (name, email, hashed_password) VALUES (?, ?, ?)`

	// Pass the arguments IN ORDER to replace the `?` parameters safely.
	result, err := db.Exec(stmt, name, email, string(hp))
	if err != nil {
		return 0, fmt.Errorf("insert failed: %w", err)
	}

	// result object contains useful metadata like the auto-incremented ID
	return result.LastInsertId()
}
