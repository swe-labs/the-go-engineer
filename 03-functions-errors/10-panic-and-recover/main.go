package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 03: Functions and Errors - panic and recover
//
// Run: go run ./03-functions-errors/10-panic-and-recover

func main() {
	fmt.Println("=== FE.10 panic and recover ===")
	fmt.Println("Learn when panic is appropriate, when it is not, and how recover turns a crash into an explicit boundary decision.")
	fmt.Println()
	fmt.Println("- Reserve panic for broken invariants, not validation errors.")
	fmt.Println("- Recover works only from a deferred function on the panicking goroutine.")
	fmt.Println("- Middleware-style recovery keeps one bad request from crashing the whole server.")
	fmt.Println()
	fmt.Println("Use panic sparingly for programmer bugs or impossible states, then recover only at boundaries where you can translate the crash into logging and containment.")
}
