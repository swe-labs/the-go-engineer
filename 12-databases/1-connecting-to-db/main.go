// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// The underscore "_" is a blank import.
	// It tells Go: "Load this package into memory and run its init() function,
	// but I won't use any of its variables directly."
	// The SQLite driver's init() function registers itself with the database/sql package.
	_ "github.com/mattn/go-sqlite3"
)

// ============================================================================
// Section 12: Databases — Connecting to SQLite
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Blank imports (`_`) and why drivers require them
//   - sql.Open vs db.Ping
//   - Handling connection pools automatically
//
// ENGINEERING DEPTH:
//   Contrary to its name, `sql.Open()` does NOT physically open a network
//   connection to the database. It merely validates the connection string
//   and sets up a Connection Pool in memory. The connection is established
//   LAZILY the very first time you execute a query or call `db.Ping()`.
//   Go manages this connection pool for you implicitly in the background.
//
// RUN: go run ./12-databases/1-connecting-to-db
// ============================================================================

var schema = `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    hashed_password BLOB NOT NULL, -- Storing as BLOB for byte slice
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

func main() {

	dbName := "data.db"

	// Clean up previous runs (for demo purposes only!)
	_ = os.Remove(dbName)

	fmt.Println("=== Connecting to SQLite ===")
	// 1. Initialize the Connection Pool
	// "sqlite3" matches the driver registered by the blank import.
	// dbName is the path to the sqlite file.
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal("Failed to parse config:", err)
	}

	// 2. ALWAYS defer db.Close()
	// This ensures the connection pool releases all sockets when main() exits.
	defer func() {
		fmt.Println("closing database connection")
		if err := db.Close(); err != nil {
			log.Printf("error closing database connection: %v", err)
		}
	}()

	// 3. Force a physical connection
	// Because sql.Open is lazy, we use db.Ping() to force Go to actually
	// connect to the database right now. If the database is down, this will fail.
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("✅ database connection established")

	// 4. Execute raw SQL (Create Table)
	// db.Exec is used for queries that do not return rows (INSERT, UPDATE, CREATE)
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	fmt.Println("✅ table was created successfully")

	// KEY TAKEAWAY:
	// - Import DB drivers with `_` to register them invisibly.
	// - `sql.Open` is lazy. It doesn't connect until you `Ping()`.
	// - `db.Exec` runs queries that don't return data.
}
