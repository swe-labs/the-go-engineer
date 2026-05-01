// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: If / Else
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how a Go program chooses one path or another based on a condition.
//
// WHY THIS MATTERS:
//   - Branching is the ability to ask a question and choose a path. With `if`, `else if`, and `else`: - one condition is checked - one branch runs - the ...
//
// RUN:
//   go run ./02-language-basics/03-control-flow/1-if-else
//
// KEY TAKEAWAY:
//   - Go's 'if' statements do not require parentheses around the condition.
//     Branching logic executes exactly one matching block top-to-bottom.
// ============================================================================

package main

import "fmt"

func main() {
	temperature := 25

	// Go does not require parentheses around the condition.
	// The curly braces {} are strictly required.
	if temperature > 30 {
		fmt.Println("Temperature is above 30C.")
	} else {
		fmt.Println("Temperature is 30C or below.")
	}

	score := 85

	// Multiple conditions are evaluated top-to-bottom.
	// As soon as one condition evaluates to true, its block runs and the
	// rest of the chain is skipped.
	if score >= 90 {
		fmt.Println("Grade: A")
	} else if score >= 80 {
		fmt.Println("Grade: B")
	} else if score >= 70 {
		fmt.Println("Grade: C")
	} else {
		fmt.Println("Grade: F")
	}

	username := ""

	if username == "" {
		fmt.Println("Username is missing.")
	} else {
		fmt.Println("Username is present.")
	}

	// Forward reference:
	// If/else is great for simple conditions, but when evaluating the same
	// variable against many distinct values, a 'switch' statement is often
	// cleaner. See: ../4-switch/README.md
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: CF.2 -> 02-language-basics/03-control-flow/2-for-basics")
	fmt.Println("Current: CF.1 (if / else)")
	fmt.Println("Previous: LB.4 (application-logger)")
	fmt.Println("---------------------------------------------------")
}
