// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 03: Data Structures - Slice Sharing and Capacity
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How sub-slices share a backing array
//   - How len and cap change when you slice a slice
//   - How append can still mutate the original data when capacity remains
//   - How to break the shared link by copying into a new slice
//
// ENGINEERING DEPTH:
//   When you create a sub-slice like `s[2:5]`, Go usually allocates no new
//   memory. It creates a new slice header that still points into the same
//   backing array. That is why slicing is cheap and also why a small view can
//   accidentally mutate or keep alive a much larger block of data.
//
// RUN: go run ./03-data-structures/5-slices-2
// ============================================================================

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
	independent := append([]int(nil), original[2:4]...)
	fmt.Printf("\nIndependent copy before change: %v\n", independent)

	independent[0] = 500
	fmt.Printf("Independent copy after change:  %v\n", independent)
	fmt.Printf("Original after copy change:     %v\n", original)

	// KEY TAKEAWAY:
	// - Sub-slicing is fast because it usually shares the original backing array.
	// - len tells you how much of the view you can read now; cap tells you how far the shared array extends.
	// - append can still mutate the original if the sub-slice has spare capacity.
	// - Copy into a new slice when you need an isolated view.
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: DS.6 contact-manager")
	fmt.Println("Current: DS.5 (slice sharing and capacity)")
	fmt.Println("---------------------------------------------------")
}
