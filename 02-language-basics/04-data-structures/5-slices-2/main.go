// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import "fmt"

// 04 Data Structures - Slice Sharing and Capacity
//
// Mental model:
// A sub-slice is usually just another view over the same backing array.
// That keeps slicing cheap, but it also means two slices can still affect the
// same stored data.
//
// Run: go run ./02-language-basics/04-data-structures/5-slices-2

func main() {
	fmt.Println("=== Slice Sharing and Capacity ===")

	original := []int{0, 1, 2, 3, 4, 5}
	fmt.Printf("Original: %v, len: %d, cap: %d\n", original, len(original), cap(original))

	// Sub-slices share the same backing array.
	shared := original[1:4]
	fmt.Printf("\nShared view original[1:4]: %v\n", shared)
	fmt.Printf("shared len: %d, cap: %d\n", len(shared), cap(shared))

	shared[0] = 100
	fmt.Printf("After shared[0] = 100, shared:   %v\n", shared)
	fmt.Printf("After shared[0] = 100, original: %v\n", original)

	// A sub-slice can still overwrite the original when append fits within cap.
	growth := original[2:4]
	fmt.Printf("\nGrowth view original[2:4]: %v, len: %d, cap: %d\n", growth, len(growth), cap(growth))

	growth = append(growth, 200)
	fmt.Printf("After append(growth, 200), growth:   %v\n", growth)
	fmt.Printf("After append(growth, 200), original: %v\n", original)

	// Copying into a new slice breaks the shared link.
	independent := make([]int, len(original[2:4]))
	for i, value := range original[2:4] {
		independent[i] = value
	}
	fmt.Printf("\nIndependent copy before change: %v\n", independent)

	independent[0] = 500
	fmt.Printf("Independent copy after change:  %v\n", independent)
	fmt.Printf("Original after copy change:     %v\n", original)

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: DS.6 contact-directory")
	fmt.Println("Current: DS.5 (slice sharing and capacity)")
	fmt.Println("---------------------------------------------------")
}
