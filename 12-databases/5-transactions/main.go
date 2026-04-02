// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// ============================================================================
// Section 12: Databases — Transactions
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - ACID Properties in relational databases
//   - db.BeginTx: initializing a transaction
//   - tx.Commit() vs tx.Rollback()
//   - The `defer tx.Rollback()` safety pattern
//
// HOW TRANSACTIONS WORK:
//   A transaction groups multiple SQL operations into a single "all-or-nothing"
//   unit. For example:
//     1. User creates an account (INSERT INTO users)
//     2. Create a profile for the user (INSERT INTO profiles)
//   If step 1 succeeds but step 2 fails, we MUST undo step 1! Otherwise, we have
//   a fragmented database.
//
// ENGINEERING DEPTH:
//   Calling `db.BeginTx()` forces the Go database driver to acquire an EXCLUSIVE
//   TCP socket connection from the Connection Pool and HOLD it. The Go driver
//   will literally block other queries from using this physical socket until you
//   call `Commit()` or `Rollback()`. This means if you start a transaction, make
//   an HTTP request to a slow external API, and then commit... you are hogging a
//   database connection for absolutely no reason. NEVER do slow non-DB work
//   inside an active SQL Transaction.
//
// RUN: go run ./12-databases/5-transactions
// ============================================================================

var profileSchema = `
CREATE TABLE IF NOT EXISTS profiles (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    avatar TEXT NOT NULL,
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

	// 1. Ensure the profiles table exists
	_, err = db.Exec(profileSchema)
	if err != nil {
		log.Fatal("Failed to create profile schema:", err)
	}

	fmt.Println("=== Executing SQL Transaction ===")

	// 2. Execute a transaction
	userID, err := createUserWithProfile(db, "TxUser", "tx@example.com", "hunter2", "avatar_01.png")
	if err != nil {
		log.Println("❌ Transaction Failed (Rolled Back):", err)
	} else {
		fmt.Printf("✅ Transaction Succeeded! Created User %d + Profile\n", userID)
	}
}

// createUserWithProfile demonstrates the canonical Go Transaction pattern.
func createUserWithProfile(db *sql.DB, name, email, hashedPassword, avatar string) (int64, error) {
	ctx := context.Background()

	// 1. BEGIN THE TRANSACTION
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// 2. THE SAFETY CATCH (defer Rollback)
	// If the function returns before calling tx.Commit(), this deferred Rollback
	// executes and undoes any partial SQL writes.
	// If tx.Commit() is successfully called later, this Rollback safely becomes a no-op!
	defer tx.Rollback()

	// 3. FIRST OPERATION (Insert User)
	// Note: We use `tx.ExecContext` instead of `db.Exec`.
	// We MUST bind all queries to the `tx` object, otherwise they execute outside
	// the transaction!
	userQuery := `INSERT INTO users (name, email, hashed_password) VALUES (?, ?, ?)`
	userResult, err := tx.ExecContext(ctx, userQuery, name, email, hashedPassword)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user (tx aborting): %w", err)
	}

	userID, err := userResult.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get user ID: %w", err)
	}

	// 4. SECOND OPERATION (Insert Profile)
	profileQuery := `INSERT INTO profiles (user_id, avatar) VALUES (?, ?)`
	_, err = tx.ExecContext(ctx, profileQuery, userID, avatar)
	if err != nil {
		// If this fails, the deferred tx.Rollback() undoes the user creation!
		return 0, fmt.Errorf("failed to insert profile (tx aborting): %w", err)
	}

	// 5. COMMIT THE TRANSACTION
	// Both operations succeeded. We tell the database to persist the changes.
	// Once committed, the deferred Rollback() ignores the command harmlessly.
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userID, nil
}
