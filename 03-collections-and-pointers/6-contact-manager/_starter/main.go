// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 3: Collections — Contact Manager (Exercise Starter)
// Level: Beginner
// ============================================================================
//
// EXERCISE: Build a Slice-based Contact Manager
//
// REQUIREMENTS:
//  1. [ ] Define a `Contact` struct with Name, Email, and Phone fields
//  2. [ ] Implement `addContact(contacts []Contact, c Contact) []Contact`
//  3. [ ] Implement `findContact(contacts []Contact, name string) (Contact, bool)`
//  4. [ ] Implement `deleteContact(contacts []Contact, name string) []Contact`
//  5. [ ] Test all operations in main() with sample data
//
// HINTS:
//   - Use `append()` to add to a slice
//   - Use `range` to iterate and find by name
//   - To delete, use append(slice[:i], slice[i+1:]...) to remove index i
//
// RUN: go run ./03-collections-and-pointers/6-contact-manager/_starter
// SOLUTION: See the main.go file in the parent directory
// ============================================================================

// TODO: Define your Contact struct here

// TODO: Implement addContact

// TODO: Implement findContact

// TODO: Implement deleteContact

func main() {
	fmt.Println("=== Contact Manager Exercise ===")
	fmt.Println()
	fmt.Println("TODO: Implement your contact manager!")
	fmt.Println("See the REQUIREMENTS above for what to build.")
	fmt.Println()
	fmt.Println("When finished, compare your solution with ../main.go")
}
