// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 03: Data Structures — Contact Manager (Exercise Starter)
// Level: Beginner
// ============================================================================
//
// EXERCISE: Build an In-Memory Contact Manager
//
// REQUIREMENTS:
//  1. [ ] Define a `Contact` struct with ID, Name, Email, and Phone fields
//  2. [ ] Store contacts in a slice and keep a `map[string]int` for name lookup
//  3. [ ] Implement `addContact(name, email, phone string)`
//  4. [ ] Implement `findContact(name string) *Contact`
//  5. [ ] Show at least one update in `main()` that persists through the returned pointer
//
// HINTS:
//   - Use `append()` to grow the master contact slice
//   - Store each contact's slice index in a map for fast lookup
//   - Returning `*Contact` lets you mutate the real stored value
//
// RUN: go run ./03-data-structures/6-contact-manager/_starter
// SOLUTION: See the main.go file in the parent directory
// ============================================================================

// TODO: Define your Contact struct here

// TODO: Add your contact slice, lookup map, and nextID state here

// TODO: Implement addContact

// TODO: Implement findContact

// TODO: In main(), prove that an update through the returned pointer persists

func main() {
	fmt.Println("=== Contact Manager Exercise ===")
	fmt.Println()
	fmt.Println("TODO: Implement your contact manager!")
	fmt.Println("See the REQUIREMENTS above for what to build, including the persistent update step.")
	fmt.Println()
	fmt.Println("When finished, compare your solution with ../main.go")
}
