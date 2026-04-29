// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Connecting to SQLite
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to use the 'database/sql' package to manage database connections.
//   - The role of database drivers and blank imports ('_').
//   - The difference between 'sql.Open' and 'db.Ping'.
//
// WHY THIS MATTERS:
//   - Most Go applications need a way to persist data. 'database/sql' provides
//     a standard, driver-agnostic way to talk to any SQL database.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/databases/1-connecting-to-db
//
// KEY TAKEAWAY:
//   - 'sql.Open' is lazy. It doesn't actually connect until you 'Ping'.
// ============================================================================

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// The underscore "_" is a blank import.
	// It registers the SQLite driver with the database/sql package
	// without us needing to use any of the driver's functions directly.
	_ "modernc.org/sqlite"
)

// Stage 06: Databases - Connecting to SQLite
//
//   - Blank imports (_) and driver registration
//   - sql.Open vs db.Ping
//   - The Connection Pool
//
// ENGINEERING DEPTH:
//   `sql.Open()` does not immediately open a connection. Instead, it
//   validates the arguments and prepares a connection pool. The actual
//   connection to the file (or server) is deferred until it's actually
//   needed. This "Lazy Connection" pattern allows your app to start up
//   quickly even if the database is temporarily slow to respond.

var schema = `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

func main() {
	dbPath := "example.db"

	// Cleanup for demonstration
	os.Remove(dbPath)

	fmt.Println("=== Connecting to SQLite ===")
	fmt.Println()

	// 1. Initialize the Connection Pool
	// "sqlite" matches the driver registered by the blank import.
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("  Failed to open config: %v", err)
	}

	// 2. ALWAYS defer db.Close()
	// This ensures the connection pool releases its resources when main exits.
	defer func() {
		fmt.Println("  Closing database connection...")
		db.Close()
	}()

	// 3. Verify the connection (Force a physical connection)
	// Because sql.Open is lazy, we use Ping to ensure the file is accessible.
	err = db.Ping()
	if err != nil {
		log.Fatalf("  Failed to connect to database: %v", err)
	}

	fmt.Println("  ✔ Database connection established.")

	// 4. Create a table
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalf("  Failed to create schema: %v", err)
	}

	fmt.Println("  ✔ Table 'users' created successfully.")

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: DB.2 query-insert")
	fmt.Println("Current: DB.1 (connecting-to-sqlite)")
	fmt.Println("Previous: API.9 (grpc-service-exercise)")
	fmt.Println("---------------------------------------------------")
}
