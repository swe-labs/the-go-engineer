// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: Switch
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Multi-way branching with `switch`.
//   - Matching multiple values in a single `case`.
//   - "Tagless" switch for complex conditional chains.
//
// WHY THIS MATTERS:
//   - Deeply nested `if/else` ladders are difficult to read and audit. `switch`
//     provides a tabular format that is easier for humans to scan and for the
//     compiler to optimize.
//
// RUN:
//   go run ./02-language-basics/03-control-flow/4-switch
//
// KEY TAKEAWAY:
//   - Go's `switch` does not require explicit `break` at the end of every case;
//     it stops by default, eliminating a major source of bugs.
// ============================================================================

package main

import "fmt"

func main() {
	day := "Monday"

	// A value-based switch evaluates a single variable ('day') against multiple
	// distinct case values.
	// Notice how cases can be grouped using a comma ("Saturday", "Sunday").
	switch day {
	case "Saturday", "Sunday":
		fmt.Println("Weekend mode.")
	case "Monday":
		fmt.Println("Start-of-week mode.")
	default: // The default case runs if no other case matches.
		fmt.Println("Regular workday mode.")
	}

	score := 82

	// A condition-based switch (switch without a target value) acts exactly
	// like an 'if / else if / else' chain. It evaluates expressions top-to-bottom
	// and executes the first one that evaluates to true.
	switch {
	case score >= 90:
		fmt.Println("Excellent result.")
	case score >= 80:
		fmt.Println("Strong result.")
	case score >= 70:
		fmt.Println("Passing result.")
	default:
		fmt.Println("Needs more work.")
	}

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: CF.5 -> 02-language-basics/03-control-flow/5-defer-basics")
	fmt.Println("Run    : go run ./02-language-basics/03-control-flow/5-defer-basics")
	fmt.Println("Current: CF.4 (switch)")
	fmt.Println("---------------------------------------------------")
}
