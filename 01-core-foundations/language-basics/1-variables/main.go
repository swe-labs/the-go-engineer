// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 1: Language Basics — Variables & Types
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to declare variables in Go
//   - The difference between `var` and `:=` (short declaration)
//   - Go's ZERO VALUE system — every type has a guaranteed default
//   - Basic types: string, int, bool, float64
//   - Why Go refuses to compile unused variables
//
// ENGINEERING DEPTH:
//   In languages like C or C++, declaring a un-initialized variable reserves a
//   block of RAM, but whatever garbage was previously at that memory address
//   is retained. This causes unpredictable bugs. Go's compiler mathematically
//   guarantees memory safety by automatically sweeping newly declared variables
//   and forcefully setting every byte to 0. This is known as "Zero Values".
//
// RUN: go run ./01-language-basics/1-variables
// ============================================================================

func main() {

	// --- METHOD 1: var declaration ---
	// The "var" keyword explicitly declares a variable with a specific type.
	// Syntax: var <name> <type>
	//
	// ZERO VALUES: When you declare with "var" but don't assign,
	// Go automatically initializes the variable to its ZERO VALUE:
	//   string → ""   (empty string)
	//   int    → 0
	//   float  → 0.0
	//   bool   → false
	//   pointer → nil
	//
	// This is a SAFETY FEATURE. In C, uninitialized variables contain garbage
	// data. Go guarantees you never encounter undefined behavior.
	var greeting string // greeting is "" right now (zero value for string)
	fmt.Printf("Initial zero value:   '%s'\n", greeting)
	greeting = "Hello, world!"

	fmt.Println(greeting)

	// var with a numeric type — zero value is 0
	var count int
	fmt.Printf("Initial zero value:    %d\n", count)
	count = 10
	fmt.Println(count)

	// var with a boolean — zero value is false
	var isRunning bool
	fmt.Printf("Initial zero value:    %t\n", isRunning)
	isRunning = true
	fmt.Println(isRunning)

	// --- Multiple declaration on one line ---
	// You can declare multiple variables of the same type together.
	var firstName, lastName string
	firstName = "John"
	lastName = "Doe"
	fmt.Println(firstName, lastName)

	// --- METHOD 2: Short declaration with := ---
	// The `:=` operator declares AND assigns in one step.
	// Go infers the type from the right-hand side (type inference).
	//
	// `:=` can ONLY be used inside functions. You cannot use it at the
	// package level (outside of func). Use "var" for package-level variables.
	//
	// WHEN TO USE WHICH:
	//   := — Inside functions, when the type is obvious from the value
	//   var — When you want to be explicit, or when you need the zero value
	email := "test@test.com" // Go infers: email is a string
	fmt.Println(email)

	age := 24 // Go infers: age is an int
	fmt.Println(age)

	// --- METHOD 3: var with initial value ---
	// You can combine var with an initial value. Go infers the type.
	var year = 2025
	fmt.Println(year)

	// --- COMPILE ERROR: Unused variables ---
	// Uncomment the line below and try to compile:
	//   name := "unused"
	//
	// Go will refuse to compile with the error:
	//   "name declared and not used"
	//
	// This is intentional. Unused variables are dead code — they waste
	// memory and confuse readers. Go eliminates this problem at compile time.

	// KEY TAKEAWAY:
	// Go has exactly 3 ways to create variables:
	//   1. var name string          — Explicit type, zero value
	//   2. var name = "value"       — Type inferred from value
	//   3. name := "value"          — Short declaration (inside functions only)
}
