// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 03: Data Structures - Contact Directory (Exercise Starter)
// Level: Beginner
// ============================================================================
//
// EXERCISE: Build an In-Memory Contact Directory
//
// REQUIREMENTS:
//  1. [ ] Create parallel `names`, `emails`, and `phones` slices
//  2. [ ] Keep a `map[string]int` for name lookup
//  3. [ ] Add at least three contacts in `main()`
//  4. [ ] Use the map to find one contact's slice index
//  5. [ ] Take a pointer to a stored phone number and prove the update persists
//
// HINTS:
//   - Use `append()` to grow each slice
//   - Keep the same index for one contact across all three slices
//   - `phones[index]` is a real slice element, so `&phones[index]` gives you a pointer
//
// RUN: go run ./03-data-structures/6-contact-manager/_starter
// SOLUTION: See the main.go file in the parent directory
// ============================================================================

func main() {
	fmt.Println("=== Contact Directory Exercise ===")
	fmt.Println()
	fmt.Println("TODO: Build the contact directory using only Section 03 concepts.")
	fmt.Println("Start with slices plus a map, then prove one pointer-based update persists.")
	fmt.Println()
	fmt.Println("When finished, compare your solution with ../main.go")
}
