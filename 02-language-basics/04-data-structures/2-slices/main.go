// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: Slices
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how Go represents dynamic collections through slices, and why `len`, `cap`, `make`, and `append` are all part of one connected idea.
//
// WHY THIS MATTERS:
//   - A slice is a small view over an underlying array. It tracks: - which array data it points to - how many elements are currently in the slice - how m...
//
// RUN:
//   go run ./02-language-basics/04-data-structures/2-slices
//
// KEY TAKEAWAY:
//   - Learn how Go represents dynamic collections through slices, and why `len`, `cap`, `make`, and `append` are all part of one connected idea.
// ============================================================================

package main

import "fmt"

// 04 Data Structures - Slices
//
// Mental model:
// A slice is a small descriptor that points at an underlying array. It tracks
// a current length and a capacity for growth.
//

func main() {
	fmt.Println("=== Slices ===")

	names := []string{"Alice", "John", "Mark"}
	fmt.Printf("names: %v\n", names)
	fmt.Printf("len=%d cap=%d\n", len(names), cap(names))

	// make creates the backing array and the first slice view.
	items := make([]int, 0, 3)
	fmt.Printf("\nitems: %v len=%d cap=%d\n", items, len(items), cap(items))

	items = append(items, 10)
	items = append(items, 20)
	items = append(items, 30)
	fmt.Printf("after three appends: %v len=%d cap=%d\n", items, len(items), cap(items))

	// This append forces growth because the original capacity was 3.
	items = append(items, 40)
	fmt.Printf("after growth append: %v len=%d cap=%d\n", items, len(items), cap(items))

	// Slicing creates a smaller view.
	firstTwo := items[:2]
	fmt.Printf("\nfirstTwo := items[:2] -> %v\n", firstTwo)

	lastTwo := items[2:]
	fmt.Printf("lastTwo := items[2:] -> %v\n", lastTwo)

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: DS.3 maps")
	fmt.Println("Current: DS.2 (slices)")
	fmt.Println("Previous: DS.1 (arrays)")
	fmt.Println("---------------------------------------------------")
}
