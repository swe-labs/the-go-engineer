// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: Slice Sharing and Capacity
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Why sub-slices share the same underlying backing array.
//   - How `append` can accidentally overwrite original data.
//   - Breaking the shared link using the `copy` pattern.
//   - Reasoning about capacity-triggered reallocations.
//
// WHY THIS MATTERS:
//   - Shared state is the #1 source of "mysterious" bugs in Go. A sub-slice
//     that modifies its data might unintentionally corrupt the original
//     collection. Engineering-grade Go requires constant awareness of whether
//     your slice is a "view" or an "independent copy."
//
// RUN:
//   go run ./02-language-basics/04-data-structures/5-slices-2
//
// KEY TAKEAWAY:
//   - Slicing creates a view; `make` + `copy` creates an independent replica.
// ============================================================================

package main

import "fmt"

// 04 Data Structures - Slice Sharing and Capacity
//
// Mental model:
// A sub-slice is usually just another view over the same backing array.
// That keeps slicing cheap, but it also means two slices can still affect the
// same stored data.
//

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

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: DS.6 -> 02-language-basics/04-data-structures/6-contact-manager")
	fmt.Println("Run    : go run ./02-language-basics/04-data-structures/6-contact-manager")
	fmt.Println("Current: DS.5 (slices-in-depth)")
	fmt.Println("---------------------------------------------------")
}
