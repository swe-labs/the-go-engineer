// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Query timeouts via context
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how context deadlines stop slow database operations from owning a request forever.
//
// WHY THIS MATTERS:
//   - A query timeout is a budget, not a guess. It bounds how long one dependency can hold the caller.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/databases/8-query-timeouts-via-context
//
// KEY TAKEAWAY:
//   - [TODO: Summarize the core takeaway]
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== DB.8 Query timeouts via context ===")
	fmt.Println("Learn how context deadlines stop slow database operations from owning a request forever.")
	fmt.Println()
	fmt.Println("- Pass context all the way down to the query call.")
	fmt.Println("- Separate request budgets from database budgets intentionally.")
	fmt.Println("- Observe timeout failures so you know whether the budget or the query needs to change.")
	fmt.Println()
	fmt.Println("Deadline discipline keeps one bad dependency call from consuming the whole request budget or worker pool.")
}

// ---------------------------------------------------
// NEXT UP: GC.0
// ---------------------------------------------------
