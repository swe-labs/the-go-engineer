// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics — Lesson CF.6: Defer in Real Use Cases
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to use defer for resource cleanup (simulated)
//   - The pattern of deferring immediately after success
//   - Why defer reduces boilerplate in error handling
//
// WHY THIS MATTERS:
//   Proper resource management is critical for production stability.
//
// RUN: go run ./02-language-basics/03-control-flow/6-defer-use-cases
// ============================================================================

package main

import (
	"fmt"
)

func main() {
	fmt.Println("CF.6: Real-World Defer Patterns")
	fmt.Println("--------------------------------")

	simulateFileOperation()

	// KEY TAKEAWAY:
	// - Defer cleanup immediately after a resource is successfully acquired.
	// - This ensures that even if subsequent logic fails, resources are released.

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: CF.7 pricing-checkout")
	fmt.Println("Current: CF.6 (defer-use-cases)")
	fmt.Println("---------------------------------------------------")
}

func simulateFileOperation() {
	fmt.Println("[Step 1] Opening simulated file 'data.txt'...")

	// Idiomatic pattern: defer the cleanup right after the "open" succeeds
	defer fmt.Println("[Step 4] File 'data.txt' closed automatically by defer.")

	fmt.Println("[Step 2] Reading data from file...")

	// Simulate some work or potential early return
	fmt.Println("[Step 3] Processing data complete.")
}
