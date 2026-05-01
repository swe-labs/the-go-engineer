// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: If / Else
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Branching logic using `if`, `else if`, and `else`.
//   - Comparison operators and boolean evaluation.
//   - Why Go requires curly braces even for single-line branches.
//
// WHY THIS MATTERS:
//   - Branching is the fundamental tool for decision-making in code. From
//     validating inputs to choosing business rules, `if/else` is the primary
//     way programs handle multiple scenarios.
//
// RUN:
//   go run ./02-language-basics/03-control-flow/1-if-else
//
// KEY TAKEAWAY:
//   - Go's `if` statements are concise (no parentheses required) but strict
//     (curly braces are mandatory), enforcing a consistent and readable style.
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

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: CF.2 -> 02-language-basics/03-control-flow/2-for-basics")
	fmt.Println("Run    : go run ./02-language-basics/03-control-flow/2-for-basics")
	fmt.Println("Current: CF.1 (if / else)")
	fmt.Println("---------------------------------------------------")
}
