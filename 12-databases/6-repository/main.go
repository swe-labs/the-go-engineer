// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================
// Section 12: Databases — Repository Pattern
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The Repository Design Pattern (Domain-Driven Design)
//   - Decoupling database logic from business logic via Interfaces
//   - Dependency Injection in Go
//
// ENGINEERING DEPTH:
//   If you pass `*sql.DB` directly into your business logic (e.g. your HTTP
//   handlers), you tightly couple your HTTP server to an operational SQLite
//   database. This makes unit testing impossible without spinning up a real
//   database. The Repository Pattern solves this by defining an Interface
//   (`UserRepository`). Your handlers only know about the Interface, never the
//   database. In production, you inject `SQLUserRepository`. In tests, you
//   inject a `MockUserRepository` that just returns hardcoded structs instantly.
//
// RUN: go run ./12-databases/6-repository
// ============================================================================

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/rasel9t6/the-go-engineer/12-databases/6-repository/repository"
)

// The startup script creates the DB table so this demo works cleanly
var profileSchema = `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL, 
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS profiles (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    avatar TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

func main() {

	// 1. Database Initialization
	dbName := "users_database.db"
	db, err := connectToDatabase(dbName)
	checkErr(err)
	defer db.Close()

	// (Demo setup: Ensure tables exist)
	_, _ = db.Exec(profileSchema)

	fmt.Println("✅ Database connection established")

	// 2. Dependency Injection
	// We instantiate the concrete SQL repository, passing the *sql.DB dependency.
	// But notice the return type of NewSQLUserRepository? It returns the Interface!
	repo := repository.NewSQLUserRepository(db)

	// 3. Application Logic
	// The application logic now uses the interface (`repo`), oblivious to the
	// fact that SQLite is powering it under the hood.
	fmt.Println("\n=== Calling Repository Methods ===")

	// Create a user
	id, err := repo.CreateUser("Alice Repo", "alice.repo@example.com", "secret", "alice_avatar.png")
	if err != nil {
		fmt.Println("⚠️ Insert failed (User might already exist):", err)
	} else {
		fmt.Printf("✅ Created User ID: %d\n", id)
	}

	// Fetch users
	printUsers(repo)
}

// printUsers ONLY accepts the `repository.UserRepository` interface.
// It has absolutely no dependency on `database/sql`.
func printUsers(repo repository.UserRepository) {
	users, err := repo.GetUsers()
	checkErr(err)
	fmt.Println("\n--- User List ---")
	for _, user := range users {
		fmt.Printf("- ID: %d | Email: %s\n", user.ID, user.Email)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func connectToDatabase(name string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
