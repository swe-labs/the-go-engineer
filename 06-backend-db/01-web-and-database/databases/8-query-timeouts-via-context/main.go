// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Query Timeouts via Context
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to use 'context.Context' to bound database operations.
//   - The importance of 'ExecContext' and 'QueryContext'.
//   - How to handle 'context.DeadlineExceeded' errors.
//
// WHY THIS MATTERS:
//   - In production, a database might become locked or slow due to load.
//     Without timeouts, your goroutines will hang indefinitely, leading
//     to a "Cascading Failure" where your entire server becomes unresponsive.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/databases/8-query-timeouts-via-context
//
// KEY TAKEAWAY:
//   - Every database query should have a deadline.
// ============================================================================

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

// Stage 06: Databases - Query Timeouts
//
//   - context.WithTimeout: Setting the budget
//   - db.QueryRowContext: Passing the budget to the driver
//   - Deadline Enforcement: Stopping the work
//
// ENGINEERING DEPTH:
//   When you pass a context to a database call, the Go driver doesn't
//   just wait; it actively monitors the context. If the timeout
//   expires, the driver sends a cancellation signal to the database
//   server (if supported) and immediately closes the socket. This
//   ensures that your app stays responsive and doesn't waste CPU
//   cycles waiting for a query that the user has already given up on.

func main() {
	dbPath := "timeout_demo.db"
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("=== Query Timeouts via Context ===")
	fmt.Println()

	// 1. Success Case: Fast query within 2 seconds
	fmt.Println("  Executing a fast query (2s budget)...")
	ctxSuccess, cancelSuccess := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelSuccess()

	err = db.PingContext(ctxSuccess)
	if err != nil {
		fmt.Printf("  ❌ Error: %v\n", err)
	} else {
		fmt.Println("  ✔ Successfully connected within budget.")
	}

	fmt.Println()

	// 2. Timeout Case: Simulated slow operation with a 10ms budget
	// Note: We use a tiny budget to ensure it fails even on a fast local disk.
	fmt.Println("  Executing a query with an impossible budget (10ms)...")
	ctxFail, cancelFail := context.WithTimeout(context.Background(), 10*time.Microsecond)
	defer cancelFail()

	// We'll try to run a query. Even a simple one should fail with such a tiny budget.
	var result int
	err = db.QueryRowContext(ctxFail, "SELECT 1").Scan(&result)

	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("  ✔ Caught expected error: Database operation timed out!")
		} else {
			fmt.Printf("  ❌ Unexpected Error: %v\n", err)
		}
	} else {
		fmt.Println("  ❌ Unexpected Success (The query was too fast!)")
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: GC.0 -> 07-concurrency/01-concurrency/goroutines/0-why-concurrency-exists")
	fmt.Println("   Current: DB.8 (query timeouts via context)")
	fmt.Println("Previous: DB.7 (n-plus-one-query-detection)")
	fmt.Println("---------------------------------------------------")
}
