// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: Variables and Types
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how Go declares variables and why every type has a predictable zero value.
//
// WHY THIS MATTERS:
//   - A variable is a named slot that holds a value while the program runs. Go gives you three common declaration shapes: 1. `var name string` 2. `var na...
//
// RUN:
//   go run ./02-language-basics/1-variables
//
// KEY TAKEAWAY:
//   - Go enforces explicit type safety and guarantees predictable zero-values
//     upon allocation, eliminating undefined memory bugs by default.
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

	// Forward reference:
	// We are printing to the console using the 'fmt' package.
	// You will learn how the 'fmt' package formats text in detail in:
	// ../../05-packages-io/02-io-and-cli/cli-tools/1-args/README.md
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: LB.2 constants")
	fmt.Println("Current: LB.1 (variables)")
	fmt.Println("Previous: none")
	fmt.Println("---------------------------------------------------")
}
