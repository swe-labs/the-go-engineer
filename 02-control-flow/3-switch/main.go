// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"time"
)

// ============================================================================
// Section 2: Control Flow — Switch Statements
// Level: Beginner → Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Basic switch with multiple case values
//   - Switch with init statement (like if with init)
//   - Tagless switch (switch without a value — acts like if/else-if)
//   - Type switch — inspecting the dynamic type of an interface value
//   - Why Go's switch does NOT fall through by default (safer than C)
//
// ENGINEERING DEPTH:
//   In C/Java, if you forget a `break` statement inside a switch case, the program
//   silently falls through to the next case, causing devastating bugs (famously
//   crashing the AT&T phone network in 1990). Go flips this: `break` is implicitly
//   injected by the compiler at the end of every case. You must explicitly use the
//   `fallthrough` keyword if you actually want that dangerous behavior.
//
// RUN: go run ./02-control-flow/3-switch
// ============================================================================

func main() {

	// --- BASIC SWITCH ---
	// Go's switch is MUCH safer than C's switch:
	//   - NO fallthrough by default — each case breaks automatically
	//   - Multiple values per case (comma-separated)
	//   - No parentheses needed around the switch expression
	day := "Monday"
	fmt.Println("Today is", day)

	switch day {
	case "Sunday", "Saturday": // Multiple values in one case
		fmt.Println("Weekend! No work")
	case "Monday", "Tuesday":
		fmt.Println("Work days. Lots of meetings")
	default: // Runs if no other case matches
		fmt.Println("Mid-week")
	}

	fmt.Println()

	// --- SWITCH WITH INIT STATEMENT ---
	// Like "if init; condition", switch can initialize a variable.
	// The variable (hour) is scoped to the switch block only.
	switch hour := time.Now().Hour(); {
	case hour < 12:
		fmt.Println("Good morning")
	case hour < 17:
		fmt.Println("Good afternoon")
	default:
		fmt.Println("Good evening")
	}

	fmt.Println()

	// --- TYPE SWITCH ---
	// Type switches inspect the DYNAMIC TYPE of an interface value.
	// The "any" type (also written as "interface{}") can hold any value.
	//
	// Syntax: switch v := x.(type) { case int: ... }
	//
	// This is crucial when working with:
	//   - JSON parsing (json.Unmarshal returns map[string]any)
	//   - Generic data processing
	//   - Plugin systems
	checkType := func(i any) {
		switch v := i.(type) {
		case int:
			fmt.Printf("Integer: %d\n", v)
		case string:
			fmt.Printf("String: %s\n", v)
		case bool:
			fmt.Printf("Boolean: %t\n", v)
		default:
			fmt.Printf("Unknown type: %T with value: %v\n", v, v)
		}
	}

	checkType(21)     // Integer: 21
	checkType("Test") // String: Test
	checkType(true)   // Boolean: true
	checkType(312.23) // Unknown type: float64 with value: 312.23

	// KEY TAKEAWAY:
	// - Go switch does NOT fall through (no need for "break" in each case)
	// - Use comma-separated values for multiple matches: case "a", "b":
	// - Tagless switch (no expression) works like if/else-if chains
	// - Type switches reveal the actual type behind an interface — essential
	//   for working with JSON, databases, and generic data handling
}
