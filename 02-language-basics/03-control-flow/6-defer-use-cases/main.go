// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: Defer in real use cases
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Idiomatic "Pair Patterns" (Open/Close, Lock/Unlock).
//   - Why deferring cleanup immediately prevents resource leaks.
//   - How `defer` handles early returns and errors safely.
//
// WHY THIS MATTERS:
//   - Resource leaks (unclosed files or database connections) are a primary
//     cause of system instability and crashes. In Go, the `defer` keyword is
//     your primary defense against these "leaks" in production.
//
// RUN:
//   go run ./02-language-basics/03-control-flow/6-defer-use-cases
//
// KEY TAKEAWAY:
//   - Clean code is safe code. Defer the cleanup right after the acquisition.
// ============================================================================

package main

import (
	"fmt"
)

func main() {
	fmt.Println("CF.6: Real-World Defer Patterns")
	fmt.Println("--------------------------------")

	simulateFileOperation()

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: CF.7 -> 02-language-basics/03-control-flow/7-pricing-checkout")
	fmt.Println("Run    : go run ./02-language-basics/03-control-flow/7-pricing-checkout")
	fmt.Println("Current: CF.6 (defer-use-cases)")
	fmt.Println("---------------------------------------------------")
}

// simulateFileOperation (Function): runs the simulate file operation step and keeps its inputs, outputs, or errors visible.
func simulateFileOperation() {
	fmt.Println("[Step 1] Opening simulated file 'data.txt'...")

	// Idiomatic pattern: defer the cleanup right after the "open" succeeds
	defer fmt.Println("[Step 4] File 'data.txt' closed automatically by defer.")

	fmt.Println("[Step 2] Reading data from file...")

	// Simulate some work or potential early return
	fmt.Println("[Step 3] Processing data complete.")
}
