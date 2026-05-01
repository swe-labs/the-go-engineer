// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Repository Pattern
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to decouple your database logic from your business logic.
//   - How to define and implement a Repository interface.
//   - How to use Dependency Injection to pass the repository to your services.
//
// WHY THIS MATTERS:
//   - High-quality Go services don't scatter SQL strings throughout their
//     handlers. The Repository pattern creates a clean boundary that
//     makes your code easier to read, maintain, and unit test.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/databases/6-repository
//
// KEY TAKEAWAY:
//   - Program to an interface, not an implementation.
// ============================================================================

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"

	"github.com/swe-labs/the-go-engineer/06-backend-db/01-web-and-database/databases/6-repository/repository"
)

// Stage 06: Databases - Repository Pattern
//
//   - Interfaces as Boundaries
//   - Implementation Hiding
//   - Dependency Injection
//
// ENGINEERING DEPTH:
//   By using a `UserRepository` interface, your application logic (like
//   an HTTP handler) no longer needs to know if the data is stored
//   in SQLite, PostgreSQL, or even a Mock in-memory database during
//   testing. This separation of concerns is the hallmark of professional
//   software engineering and is essential for building scalable,
//   maintainable systems.

func main() {
	dbPath := "repository_demo.db"
	os.Remove(dbPath) // Start fresh

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 1. Setup the tables (normally handled by migrations)
	setup(db)

	fmt.Println("=== Repository Pattern ===")
	fmt.Println()

	// 2. Initialize the repository
	// Note: We inject the concrete 'db' into the repository.
	userRepo := repository.NewSQLUserRepository(db)

	// 3. Use the repository via its interface
	// In a real app, you would pass 'userRepo' to your handlers.
	ctx := context.Background()

	fmt.Println("  Creating user via Repository...")
	id, err := userRepo.Create(ctx, "Alice Repo", "alice@repo.com", "secret123", "alice.png")
	if err != nil {
		log.Fatalf("  Error creating user: %v", err)
	}
	fmt.Printf("  ✔ Created User ID: %d\n", id)

	fmt.Println("\n  Listing all users from Repository:")
	users, err := userRepo.List(ctx)
	if err != nil {
		log.Fatalf("  Error listing users: %v", err)
	}

	for _, u := range users {
		fmt.Printf("  - [%d] %s (%s)\n", u.ID, u.Name, u.Email)
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: DB.7 -> 06-backend-db/01-web-and-database/databases/7-n-plus-one-query-detection")
	fmt.Println("Current: DB.6 (repository-pattern)")
	fmt.Println("Previous: DB.5 (transactions)")
	fmt.Println("---------------------------------------------------")
}

func setup(db *sql.DB) {
	db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, email TEXT UNIQUE, password TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP);`)
	db.Exec(`CREATE TABLE profiles (user_id INTEGER PRIMARY KEY, avatar TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP);`)
}
