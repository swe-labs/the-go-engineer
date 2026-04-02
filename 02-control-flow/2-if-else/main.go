// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 2: Control Flow — If/Else Statements
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Basic if/else and if/else-if chains
//   - Go's unique "if with init statement" pattern
//   - The comma-ok idiom for map lookups
//   - Why Go doesn't need parentheses around conditions
//
// ENGINEERING DEPTH:
//   The "If with Init" statement (`if err := doSomething(); err != nil`) is the
//   linchpin of Go's famous error handling pattern. Structurally, it creates a
//   micro-lexical scope. Variables declared in the init statement are destroyed
//   the millisecond the `if` block terminates. This prevents namespace pollution
//   and ensures variables never accidentally leak into the outer function scope.
//
// RUN: go run ./02-control-flow/2-if-else
// ============================================================================

func main() {

	// --- BASIC IF/ELSE ---
	// Go's if/else looks like C/Java but WITHOUT parentheses around the condition.
	// The braces { } are ALWAYS required, even for single-line bodies.
	// This prevents a whole class of "goto fail" style bugs.
	tmp := 25
	if tmp > 30 {
		fmt.Println("Temperature is above 30°C — hot!")
	} else {
		fmt.Println("Temperature is 30°C or below — comfortable")
	}

	// --- IF/ELSE-IF CHAIN ---
	// Multiple conditions are checked top-to-bottom. The FIRST match wins.
	// Always order from most specific to least specific.
	score := 85
	if score >= 90 {
		fmt.Println("Grade: A — Excellent")
	} else if score >= 80 {
		fmt.Println("Grade: B — Good") // ← This matches (85 >= 80)
	} else if score >= 70 {
		fmt.Println("Grade: C — Average")
	} else {
		fmt.Println("Grade: F — Failed")
	}

	// --- IF WITH INIT STATEMENT (Go's Signature Pattern) ---
	// Go allows you to declare and initialize a variable INSIDE the if statement.
	// Syntax: if <init>; <condition> { ... }
	//
	// The variable created in the init statement is SCOPED to the if/else block.
	// It doesn't leak into the surrounding function — keeping scope tight.
	//
	// This is extremely common with map lookups (the "comma-ok" pattern):
	userAccess := map[string]bool{
		"jane": true,
		"john": false,
	}

	// hasAccess and ok are created IN the if statement.
	// They only exist inside this if/else block.
	if hasAccess, ok := userAccess["john"]; ok && hasAccess {
		fmt.Println("John can access the system")
	} else if ok && !hasAccess {
		fmt.Println("John exists but access is denied") // ← This matches
	} else {
		fmt.Println("User not found")
	}

	// Trying to use hasAccess here would be a COMPILE ERROR:
	// fmt.Println(hasAccess) ← "undefined: hasAccess"

	// --- THE COMMA-OK PATTERN EXPLAINED ---
	// When reading from a map: value, ok := myMap[key]
	//   - ok is true if the key exists
	//   - ok is false if the key doesn't exist (value will be the zero value)
	//
	// ALWAYS check ok when reading from maps. Without it, you can't tell
	// the difference between "key doesn't exist" and "key exists with value 0".
	if _, exists := userAccess["admin"]; !exists {
		fmt.Println("admin user not found in the system")
	}

	// KEY TAKEAWAY:
	// - No parentheses around conditions (Go enforces this)
	// - Braces are mandatory (Go enforces this too)
	// - "if init; condition { }" keeps variables scoped tightly
	// - The comma-ok pattern (value, ok := map[key]) is used everywhere in Go
}
