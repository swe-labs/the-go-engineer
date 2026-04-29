// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: N+1 Query Detection
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to spot the N+1 query problem in your code.
//   - Why repeated database calls inside a loop kill performance.
//   - How to use SQL JOINS to solve N+1 and reduce database chatter.
//
// WHY THIS MATTERS:
//   - The N+1 problem is one of the most common performance bottlenecks in
//     backend systems. A small code change can turn 1 fast query into
//     1,001 slow queries, crashing your database under load.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/databases/7-n-plus-one-query-detection
//
// KEY TAKEAWAY:
//   - 1 JOIN is almost always faster than N queries.
// ============================================================================

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

// Stage 06: Databases - N+1 Query Detection
//
//   - The Problem: Loop-driven database calls
//   - The Symptom: Unnecessary network latency
//   - The Solution: Eager Loading via JOINs
//
// ENGINEERING DEPTH:
//   Every database query has a "Network Round-Trip" cost. Even if the
//   query takes 1ms, the network latency might be 10ms. If you have
//   100 users and you query their profiles in a loop, you are adding
//   1,000ms (1 second) of latency to your request! By using a JOIN,
//   you pay the network cost only once.

func main() {
	dbPath := "nplusone.db"
	os.Remove(dbPath)

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 1. Setup Data (5 users, each with a profile)
	setup(db)

	fmt.Println("=== N+1 Query Detection ===")
	fmt.Println()

	// 2. THE NAIVE WAY (N+1)
	// Query 1: Get all users
	// Query N: Get profile for each user
	fmt.Println("  [Strategy] Naive Loop (N+1):")
	queryCount := 0

	rows, _ := db.Query("SELECT id, name FROM users")
	queryCount++ // First query

	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)

		// BAD: Querying inside a loop!
		var avatar string
		db.QueryRow("SELECT avatar FROM profiles WHERE user_id = ?", id).Scan(&avatar)
		queryCount++ // N inner queries

		fmt.Printf("    - User: %s | Avatar: %s\n", name, avatar)
	}
	rows.Close()
	fmt.Printf("  >> Total Queries Executed: %d\n", queryCount)

	fmt.Println()

	// 3. THE OPTIMIZED WAY (1 Query)
	// Query 1: Get users and profiles in one go using a JOIN
	fmt.Println("  [Strategy] Optimized JOIN (1 Query):")
	queryCount = 0

	optimizedQuery := `
		SELECT u.name, p.avatar
		FROM users u
		JOIN profiles p ON u.id = p.user_id
	`
	rows, _ = db.Query(optimizedQuery)
	queryCount++ // Only one query!

	for rows.Next() {
		var name, avatar string
		rows.Scan(&name, &avatar)
		fmt.Printf("    - User: %s | Avatar: %s\n", name, avatar)
	}
	rows.Close()
	fmt.Printf("  >> Total Queries Executed: %d\n", queryCount)

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: DB.8 query-timeouts-via-context")
	fmt.Println("Current: DB.7 (n-plus-one-query-detection)")
	fmt.Println("Previous: DB.6 (repository-pattern)")
	fmt.Println("---------------------------------------------------")
}

func setup(db *sql.DB) {
	db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT);")
	db.Exec("CREATE TABLE profiles (user_id INTEGER, avatar TEXT);")

	for i := 1; i <= 5; i++ {
		name := fmt.Sprintf("User_%d", i)
		avatar := fmt.Sprintf("avatar_%d.png", i)
		db.Exec("INSERT INTO users (id, name) VALUES (?, ?)", i, name)
		db.Exec("INSERT INTO profiles (user_id, avatar) VALUES (?, ?)", i, avatar)
	}
}
