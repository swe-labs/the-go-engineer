// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// ============================================================================
// Section 12: Databases — Prepared Statements & Context
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Prepared Statements (`db.Prepare()`)
//   - Passing Context to queries (`ExecContext`) for cancellations
//   - When to actually use Prepared Statements (and when to avoid them)
//
// ENGINEERING DEPTH:
//   When you run `db.Query()`, the database parses, compiles, and optimizes
//   the SQL string from scratch every single time. A "Prepared Statement"
//   compiles the SQL layout ONCE and stores it in the database's cache. You
//   then execute it repeatedly just by passing arguments to the cache reference.
//   HOWEVER, Go implicitly prepares and caches standard `db.Query()` queries
//   for you under the hood automatically anyway! In Go, manually calling
//   `db.Prepare` is only beneficial if you are executing the EXACT same
//   statement hundreds of times in a tight `for` loop.
//
// RUN: go run ./12-databases/4-prepare
// ============================================================================

type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
}

func main() {
	dbName := "users_database.db"
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create context with a timeout.
	// We mandate that this database operation MUST complete within 2 seconds.
	// If the database is locked or slow, it will abort early rather than
	// hanging the server indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Release the context resources when done

	fmt.Println("=== Executing Prepared Statement with Context ===")

	// Executing the prepared statement
	id, err := createUserWithCtx(ctx, db, "Context User", "ctx@go.dev", "hunter2")
	if err != nil {
		log.Println("❌ Failed:", err)
	} else {
		fmt.Printf("✅ Inserted user via Prepared Statement! (ID: %d)\n", id)
	}
}

// createUserWithCtx demonstrates explicit Prepared Statements and Context execution.
func createUserWithCtx(ctx context.Context, db *sql.DB, name, email, hashedPassword string) (int64, error) {

	// PREPARE: The database parses and compiles this structure ONCE.
	// In a real application, you would prepare this statement ONCE globally on
	// startup, and reuse the `stmt` variable across millions of requests.
	// In this demo, since we are inside a function, we have to close it.
	stmt, err := db.Prepare(`INSERT INTO users (name, email, hashed_password) VALUES (?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	// A prepared statement consumes a connection pool socket. It MUST be closed.
	defer stmt.Close()

	hp, err := bcrypt.GenerateFromPassword([]byte(hashedPassword), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// EXECUTE WITH CONTEXT:
	// We pass the context to `ExecContext`. If the context timeout expires
	// before the operation finishes, the Go driver automatically severs the
	// database socket and returns a `context deadline exceeded` error.
	result, err := stmt.ExecContext(ctx, name, email, string(hp))
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
