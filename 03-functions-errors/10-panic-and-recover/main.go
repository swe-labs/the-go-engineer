// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 03: Functions and Errors
// Title: panic and recover
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn when panic is appropriate, when it is not, and how recover turns a crash into an explicit boundary decision.
//
// WHY THIS MATTERS:
//   - Errors describe expected failure. Panic describes broken assumptions. Recover belongs at process or request boundaries, not in ordinary business flow.
//
// RUN:
//   go run ./03-functions-errors/10-panic-and-recover
//
// KEY TAKEAWAY:
//   - [TODO: Summarize the core takeaway]
// ============================================================================

package main

import "fmt"

//

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

// ---------------------------------------------------
// NEXT UP: TI.1
// ---------------------------------------------------
