package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 06: Backend, APIs & Databases - Query timeouts via context
//
// Run: go run ./06-backend-db/01-web-and-database/databases/8-query-timeouts-via-context

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
