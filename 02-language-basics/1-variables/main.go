// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: Variables and Types
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Three common declaration shapes (`var name type`, `var name = val`, `name := val`).
//   - Why every type has a predictable zero value.
//   - How the compiler enforces code quality via unused variable checks.
//
// WHY THIS MATTERS:
//   - In many languages, uninitialized variables contain "garbage" data from previous
//     memory usage. Go eliminates this entire class of bugs by guaranteeing a
//     predictable zero state for every allocated variable.
//
// RUN:
//   go run ./02-language-basics/1-variables
//
// KEY TAKEAWAY:
//   - Predictable zero-values and strict unused-variable checks make Go code
//     stable and easier to audit under pressure.
// ============================================================================

package main

import "fmt"

func main() {
	// Explicit declaration using 'var'.
	// Go guarantees 'greeting' starts with the zero value for a string (""),
	// rather than pointing to random, uninitialized memory.
	var greeting string
	fmt.Printf("Initial zero value: '%s'\n", greeting)

	// Assignment updates the value of the already declared variable.
	greeting = "Hello, world!"
	fmt.Println(greeting)

	var count int
	fmt.Printf("Initial zero value: %d\n", count)
	count = 10
	fmt.Println(count)

	// The zero value depends on the type. For 'bool', it is 'false'.
	var isActive bool
	fmt.Printf("Initial zero value: %t\n", isActive)
	isActive = true
	fmt.Println(isActive)

	// Short variable declaration (:=) infers the type and declares the variable in one step.
	// This is the most common way to declare local variables in Go.
	firstName, lastName := "John", "Doe"
	fmt.Println(firstName, lastName)

	email := "test@test.com"
	fmt.Println(email)

	age := 24
	fmt.Println(age)

	var year = 2025
	fmt.Println(year)

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: LB.2 -> 02-language-basics/2-constants")
	fmt.Println("Run    : go run ./02-language-basics/2-constants")
	fmt.Println("Current: LB.1 (variables)")
	fmt.Println("---------------------------------------------------")
}
