// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"slices"
)

// ============================================================================
// Section 3: Collections & Pointers — Advanced Slicing
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Advanced slice expressions and capacity behavior
//   - How sub-slices share memory (and when they stop sharing)
//   - "Slicing a slice" to isolate data
//   - Introduction to the new Go 1.21 `slices` package
//
// ENGINEERING DEPTH:
//   When you sub-slice (e.g., `s[2:5]`), the Go runtime does NOT allocate any
//   new memory. It simply creates a new 24-byte Slice Header with a pointer
//   that acts as an offset to the original contiguous array. This makes
//   window-based algorithms (like parsing network packets) virtually instant,
//   but introduces memory leak risks if a tiny sub-slice keeps a gigantic
//   10GB underlying array alive in memory.
//
// RUN: go run ./03-collections-and-pointers/5-slices-2
// ============================================================================

func main() {

	fmt.Println("=== Advanced Slicing Operations ===")

	// Create a base slice. This allocates a backing array of 10 integers.
	original := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Printf("Original: %v, len: %d, cap: %d\n", original, len(original), cap(original))

	// --- SLICING SYNTAX EXPLAINED ---
	// Form: slice[low:high]
	// - low is INCLUSIVE (starts pulling at this index)
	// - high is EXCLUSIVE (stops pulling right before this index)

	// Pull from index 2 up to (but not including) index 5.
	// Elements taken: original[2], original[3], original[4].
	// CAPACITY behavior: The capacity of a sub-slice is the capacity of the
	// original slice MINUS the starting offset (`low`).
	// So, cap = 10 - 2 = 8.
	s1 := original[2:5]
	fmt.Printf("\n  s1 (original[2:5]): %v\n", s1)
	fmt.Printf("  s1 len: %d, cap: %d\n", len(s1), cap(s1))

	// --- OMITTING BOUNDARIES ---

	// If you omit `low`, it defaults to 0. (Starts from the beginning)
	// Grabs elements at index 0, 1, 2, 3.
	s2 := original[:4]
	fmt.Printf("\n  s2 (original[:4]): %v\n", s2)

	// If you omit `high`, it defaults to the `len(original)`.
	// Grabs elements from index 6 all the way to the end.
	s3 := original[6:]
	fmt.Printf("  s3 (original[6:]): %v\n", s3)

	// If you omit BOTH, it copies the entire slice header.
	// It points to the exact same memory, with the exact same length and capacity.
	s4 := original[:]
	fmt.Printf("  s4 (original[:]):  %v\n", s4)

	// --- THE SLICES STANDARD PACKAGE ---
	// Introduced in Go 1.21, the "slices" package provides generic algorithms for slices.
	// slices.Contains iterates through the slice to find an exact match.
	found := slices.Contains(s4, 8)
	fmt.Printf("\n  s4 contains 8? %v\n", found)

	// --- THE "APPEND" SHARED MEMORY TRAP ---
	// This is the #1 bug written by intermediate Go developers.
	// Append adds an element. If there is enough CAPACITY in the backing array,
	// append overwrites the next memory slot and does NOT allocate new memory.
	// Because `s4` shares memory with `original` (and len=10, cap=10),
	// appending to s4 exceeds capacity, FORCING A RE-ALLOCATION!
	s4 = append(s4, 1000)

	fmt.Println("\n=== The Re-Allocation Trap ===")
	fmt.Printf("s4 (after append): %v\n", s4)
	fmt.Printf("Original:          %v\n", original)

	// Notice that `original` was NOT modified!
	// When s4 needed space for the 11th element, capacity was maxed (10/10).
	// Go automatically allocated a brand new array in memory (cap=20),
	// copied the data over, and appended 1000. `s4` and `original` are no
	// longer sharing the same memory.

	// KEY TAKEAWAY:
	// - Sub-slicing is fast and cheap, but shares memory.
	// - The capacity of a sub-slice extends to the end of the backing array.
	// - Appending to a sub-slice CAN overwrite data in the original slice if capacity allows it.
	// - If capacity is exceeded, append allocates new memory, breaking the link.
}
