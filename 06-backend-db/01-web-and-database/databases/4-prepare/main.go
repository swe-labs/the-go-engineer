// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Prepared Statements
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to 'pre-compile' SQL queries using 'db.Prepare'.
//   - The performance benefits of reusing prepared statements.
//   - How to manage the lifecycle of a statement with 'stmt.Close'.
//
// WHY THIS MATTERS:
//   - Prepared statements improve performance by allowing the database to
//     parse and plan a query once, then execute it many times with
//     different data. They also provide an extra layer of protection
//     against SQL injection.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/databases/4-prepare
//
// KEY TAKEAWAY:
//   - Prepare once, execute many. Always close your statements.
// ============================================================================

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// Stage 06: Databases - Prepared Statements
//
//   - db.Prepare: Compiling the SQL template
//   - stmt.Exec: Running the compiled template with data
//   - Resource Management: Closing statements
//
// ENGINEERING DEPTH:
//   When you call `db.Query()`, the database must parse the SQL string,
//   validate the table names, check permissions, and create an execution
//   plan. By using `db.Prepare()`, you move all that work to a single
//   "Preparation" step. Subsequent calls only send the raw data values
//   and the ID of the prepared statement. This reduces CPU load on the
//   database and decreases network overhead.

func main() {
	dbPath := "example.db"

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("=== Prepared Statements ===")
	fmt.Println()

	// 1. Setup table
	db.Exec(`CREATE TABLE IF NOT EXISTS inventory (id INTEGER PRIMARY KEY, item TEXT, quantity INTEGER);`)

	// 2. PREPARE the statement (The "Compilation" step)
	// We define the structure with '?' placeholders.
	stmt, err := db.Prepare(`INSERT INTO inventory (item, quantity) VALUES (?, ?)`)
	if err != nil {
		log.Fatalf("  Failed to prepare statement: %v", err)
	}
	// CRITICAL: A prepared statement holds a resource. You must close it!
	defer stmt.Close()

	fmt.Println("  Statement prepared. Executing in a loop...")

	// 3. EXECUTE the statement multiple times
	items := []string{"Apples", "Bananas", "Cherries"}
	for i, name := range items {
		_, err := stmt.Exec(name, (i+1)*10)
		if err != nil {
			log.Printf("  Error inserting %s: %v", name, err)
			continue
		}
		fmt.Printf("  ✔ Inserted %s\n", name)
	}

	fmt.Println("\n  Batch insertion complete.")

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: DB.5 -> 06-backend-db/01-web-and-database/databases/5-transactions")
	fmt.Println("Current: DB.4 (prepare-statements)")
	fmt.Println("Previous: DB.3 (select-queries)")
	fmt.Println("---------------------------------------------------")
}
