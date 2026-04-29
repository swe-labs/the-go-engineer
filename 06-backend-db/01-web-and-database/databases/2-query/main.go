// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Executing Queries (INSERT)
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to insert data into a table using 'db.Exec'.
//   - How to use parameterized queries ('?') to prevent SQL Injection.
//   - How to retrieve the ID of the newly inserted row.
//
// WHY THIS MATTERS:
//   - SQL Injection is one of the most common security vulnerabilities.
//     Learning how to safely handle user input in your queries is
//     essential for every backend engineer.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/databases/2-query
//
// KEY TAKEAWAY:
//   - NEVER concatenate strings to build SQL. Always use placeholders.
// ============================================================================

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// Stage 06: Databases - Executing Queries (INSERT)
//
//   - SQL Parameterization (?)
//   - db.Exec for mutations
//   - LastInsertId()
//
// ENGINEERING DEPTH:
//   When you use `db.Exec("INSERT ... VALUES (?)", value)`, the Go
//   database driver sends the query and the data to the database
//   separately. This means the database treats the data as a literal
//   value, not as part of the SQL command. This is the primary defense
//   against SQL Injection attacks.

func main() {
	dbPath := "example.db"

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("=== Executing INSERT Queries ===")
	fmt.Println()

	// 1. Ensure the table exists (from Lesson 1)
	schema := `CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, email TEXT UNIQUE);`
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalf("  Failed to create table: %v", err)
	}

	// 2. Insert a new user safely using placeholders
	name := "Rasel Hossen"
	email := "rasel@example.com"

	// The '?' is a placeholder for a single value.
	query := `INSERT INTO users (name, email) VALUES (?, ?)`

	// db.Exec returns a Result object and an error.
	result, err := db.Exec(query, name, email)
	if err != nil {
		// Handle specific errors like unique constraint violations
		log.Fatalf("  Failed to insert user: %v", err)
	}

	// 3. Get the auto-incremented ID
	id, _ := result.LastInsertId()
	// Get the number of rows affected (should be 1)
	rows, _ := result.RowsAffected()

	fmt.Printf("  ✔ Inserted user: %s (ID: %d)\n", name, id)
	fmt.Printf("  ✔ Rows affected: %d\n", rows)

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: DB.3 select-queries")
	fmt.Println("Current: DB.2 (query-insert)")
	fmt.Println("Previous: DB.1 (connecting-to-sqlite)")
	fmt.Println("---------------------------------------------------")
}
