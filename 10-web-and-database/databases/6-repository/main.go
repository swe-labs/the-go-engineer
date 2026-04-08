// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================
// Section 10: Databases - Repository Pattern
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The Repository Design Pattern (Domain-Driven Design)
//   - Decoupling database logic from business logic via interfaces
//   - Dependency injection in Go
//
// ENGINEERING DEPTH:
//   If you pass `*sql.DB` directly into your business logic (for example HTTP
//   handlers), you tightly couple your service layer to a live SQLite database.
//   That makes focused testing much harder. The Repository Pattern solves this
//   by defining an interface (`UserRepository`). Higher-level code depends on
//   the interface, not the concrete database implementation.
//
// RUN: go run ./10-web-and-database/databases/6-repository
// ============================================================================

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/rasel9t6/the-go-engineer/10-web-and-database/databases/6-repository/repository"
)

// repositorySchema creates the database tables needed by the repository demo.
var repositorySchema = `
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
	_ = os.Remove(dbName)

	db, err := connectToDatabase(dbName)
	checkErr(err)
	defer db.Close()

	// Demo setup: ensure tables exist so the example is deterministic.
	_, err = db.Exec(repositorySchema)
	checkErr(err)

	fmt.Println("Database connection established")

	// 2. Dependency Injection
	repo := repository.NewSQLUserRepository(db)

	// 3. Application Logic
	fmt.Println("\n=== Calling Repository Methods ===")

	id, err := repo.CreateUser("Alice Repo", "alice.repo@example.com", "secret", "alice_avatar.png")
	if err != nil {
		fmt.Println("Warning: insert failed (user might already exist):", err)
	} else {
		fmt.Printf("Created User ID: %d\n", id)
	}

	printUsers(repo)

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("Milestone complete: DB.6 repository pattern")
	fmt.Println("Next section: Section 11 concurrency")
	fmt.Println("---------------------------------------------------")
}

// printUsers ONLY accepts the repository.UserRepository interface.
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

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
