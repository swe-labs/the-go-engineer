// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 09: Architecture & Security
// Title: SQL injection prevention
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn why parameterized queries are the baseline defense against SQL injection in Go.
//
// WHY THIS MATTERS:
//   - SQL injection happens when untrusted input is treated as part of the query syntax instead of as data.
//
// RUN:
//   go run ./09-architecture/04-security/02-sql-injection-prevention
//
// KEY TAKEAWAY:
//   - Learn why parameterized queries are the baseline defense against SQL injection in Go.
// ============================================================================

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// vulnerableLogin (Function): demonstrates SQL injection by concatenating user input into a query string.
func vulnerableLogin(db *sql.DB, email, password string) (bool, error) {
	query := fmt.Sprintf("SELECT id FROM users WHERE email = '%s' AND password = '%s'", email, password)
	fmt.Println("  [VULNERABLE] Executing:", query)

	var id int
	err := db.QueryRow(query).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	fmt.Printf("  [VULNERABLE] Logged in as user id=%d\n", id)
	return true, nil
}

// safeLogin (Function): uses parameterized queries to bind user input as data, never as SQL syntax.
func safeLogin(db *sql.DB, email, password string) (bool, error) {
	query := "SELECT id FROM users WHERE email = ? AND password = ?"
	fmt.Println("  [SAFE] Executing: SELECT id FROM users WHERE email = ? AND password = ?")

	var id int
	err := db.QueryRow(query, email, password).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	fmt.Printf("  [SAFE] Logged in as user id=%d\n", id)
	return true, nil
}

// seedDB (Function): creates the users table and inserts sample data.
func seedDB(db *sql.DB) {
	schema := `CREATE TABLE users (
		id    INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL,
		password TEXT NOT NULL
	)`
	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("create table: %v", err)
	}

	users := []struct {
		email    string
		password string
	}{
		{"alice@example.com", "hunter2"},
		{"bob@example.com", "p@ssw0rd"},
		{"admin@example.com", "supersecret"},
	}
	for _, u := range users {
		if _, err := db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", u.email, u.password); err != nil {
			log.Fatalf("insert user: %v", err)
		}
	}
	fmt.Println("  Seeded 3 users.")
}

func main() {
	fmt.Println("=== SEC.2 SQL injection prevention ===")
	fmt.Println()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	seedDB(db)
	fmt.Println()

	// ── 1. Normal login — both approaches work ──
	fmt.Println("── Normal login: alice@example.com / hunter2 ──")
	vulnOK, _ := vulnerableLogin(db, "alice@example.com", "hunter2")
	safeOK, _ := safeLogin(db, "alice@example.com", "hunter2")
	fmt.Printf("  vulnerableLogin=%t  safeLogin=%t\n\n", vulnOK, safeOK)

	// ── 2. Wrong password — both reject ──
	fmt.Println("── Wrong password: alice@example.com / wrong ──")
	vulnOK, _ = vulnerableLogin(db, "alice@example.com", "wrong")
	safeOK, _ = safeLogin(db, "alice@example.com", "wrong")
	fmt.Printf("  vulnerableLogin=%t  safeLogin=%t\n\n", vulnOK, safeOK)

	// ── 3. SQL injection — vulnerable succeeds, safe rejects ──
	inj := "' OR '1'='1"
	fmt.Printf("── SQL injection: email=%q / password=%q ──\n", inj, inj)
	vulnOK, _ = vulnerableLogin(db, inj, inj)
	safeOK, _ = safeLogin(db, inj, inj)
	fmt.Printf("  vulnerableLogin=%t  safeLogin=%t\n\n", vulnOK, safeOK)

	fmt.Println("── Key insight ──")
	fmt.Println("  The vulnerable function treated user input as SQL syntax,")
	fmt.Println("  so  \"' OR '1'='1\"  turned the WHERE clause into a tautology.")
	fmt.Println("  The safe function bound the same string as data, so it")
	fmt.Println("  searched for a literal email/password — and found nothing.")
	fmt.Println()
	fmt.Println("──────────────────────────────────────────────────────")
	fmt.Println("NEXT UP: SEC.3 -> 09-architecture/04-security/03-xss-and-csrf")
	fmt.Println("Current: SEC.2 (sql injection prevention)")
	fmt.Println("──────────────────────────────────────────────────────")
}
