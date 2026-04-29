// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: Switch
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how to choose among several possible paths without building long, hard-to-scan branch chains.
//
// WHY THIS MATTERS:
//   - `switch` is a multi-way branch. It is useful when: - one value may match several known cases - several conditions need a clean top-to-bottom table ...
//
// RUN:
//   go run ./02-language-basics/03-control-flow/4-switch
//
// KEY TAKEAWAY:
//   - 'switch' is a cleaner alternative to deep 'if/else' chains. It implicitly
//     breaks after matching a case, preventing accidental fallthrough bugs
//     common in other languages.
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

	// Forward reference:
	// So far we have learned standard synchronous control flow. Go also provides
	// a unique keyword called 'defer' to schedule cleanup work. We will cover
	// that next: ../5-defer-basics/README.md
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: CF.5 defer-basics")
	fmt.Println("Current: CF.4 (switch)")
	fmt.Println("Previous: CF.3 (break-continue)")
	fmt.Println("---------------------------------------------------")
}
