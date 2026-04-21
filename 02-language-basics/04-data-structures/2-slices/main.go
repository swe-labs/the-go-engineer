// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import "fmt"

// 04 Data Structures - Slices
//
// Mental model:
// A slice is a small descriptor that points at an underlying array. It tracks
// a current length and a capacity for growth.
//
// Run: go run ./02-language-basics/04-data-structures/2-slices

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
	fmt.Println("---------------------------------------------------")
}
