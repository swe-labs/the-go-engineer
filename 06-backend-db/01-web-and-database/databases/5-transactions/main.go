// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Transactions
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to use 'db.Begin' to start a transaction.
//   - The importance of 'tx.Commit' and 'tx.Rollback'.
//   - How to ensure multiple SQL operations succeed or fail as a single unit.
//
// WHY THIS MATTERS:
//   - Transactions are the bedrock of data integrity. They prevent
//     "Partial Writes"-scenarios where half of your data is updated
//     but the other half fails, leaving your database in a broken state.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/databases/5-transactions
//
// KEY TAKEAWAY:
//   - All or nothing. Never leave your data in an inconsistent state.
// ============================================================================

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// Stage 06: Databases - Transactions
//
//   - Atomicity: The "All or Nothing" rule
//   - tx.Commit(): Saving the changes permanently
//   - tx.Rollback(): Undoing the changes if something fails
//
// ENGINEERING DEPTH:
//   A transaction locks a database connection from the pool. While a
//   transaction is active, that specific connection cannot be used for
//   anything else. This is why you should keep your transactions as
//   short as possible. Never perform slow operations (like calling
//   an external API or sending an email) inside a transaction, as
//   it will "hog" the database resource and slow down your entire app.

func main() {
	dbPath := "example.db"
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("=== SQL Transactions ===")
	fmt.Println()

	// 1. Setup tables
	setup(db)

	// 2. Successful Transaction (Create User + Profile)
	fmt.Println("  Attempting a successful transaction...")
	err = createUserWithProfile(db, "Alice", "alice@example.com", "avatar_alice.png")
	if err != nil {
		fmt.Printf("  ❌ Failed: %v\n", err)
	} else {
		fmt.Println("  ✔ Successfully created user and profile.")
	}

	// 3. Failing Transaction (Simulated error)
	fmt.Println("\n  Attempting a failing transaction (Duplicate Email)...")
	err = createUserWithProfile(db, "Alice", "alice@example.com", "avatar_dup.png")
	if err != nil {
		fmt.Printf("  ✔ Correctly rolled back: %v\n", err)
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: DB.6 -> 06-backend-db/01-web-and-database/databases/6-repository")
	fmt.Println("Current: DB.5 (transactions)")
	fmt.Println("Previous: DB.4 (prepare-statements)")
	fmt.Println("---------------------------------------------------")
}

// createUserWithProfile (Function): runs the create user with profile step and keeps its inputs, outputs, or errors visible.
func createUserWithProfile(db *sql.DB, name, email, avatar string) error {
	// A. START THE TRANSACTION
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// B. DEFER ROLLBACK
	// This is a safety net. If the function returns before Commit is called,
	// Rollback will undo everything. If Commit is called successfully,
	// the Rollback call becomes a harmless no-op.
	defer tx.Rollback()

	// C. FIRST OPERATION: Insert User
	// Note: We use tx.Exec, NOT db.Exec!
	res, err := tx.Exec(`INSERT INTO users (name, email) VALUES (?, ?)`, name, email)
	if err != nil {
		return fmt.Errorf("user insert failed: %w", err)
	}

	userID, _ := res.LastInsertId()

	// D. SECOND OPERATION: Insert Profile
	_, err = tx.Exec(`INSERT INTO profiles (user_id, avatar) VALUES (?, ?)`, userID, avatar)
	if err != nil {
		return fmt.Errorf("profile insert failed: %w", err)
	}

	// E. COMMIT THE TRANSACTION
	// This makes all changes permanent in the database.
	return tx.Commit()
}

// setup (Function): runs the setup step and keeps its inputs, outputs, or errors visible.
func setup(db *sql.DB) {
	db.Exec(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, email TEXT UNIQUE);`)
	db.Exec(`CREATE TABLE IF NOT EXISTS profiles (id INTEGER PRIMARY KEY, user_id INTEGER, avatar TEXT);`)
}
