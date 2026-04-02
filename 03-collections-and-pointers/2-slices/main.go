// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 3: Collections & Pointers — Slices
// Level: Beginner → Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Slices are Go's primary collection type (used 100x more than arrays)
//   - Slice internals: pointer + length + capacity (the slice header)
//   - make() for pre-allocating slices
//   - append() for growing slices dynamically
//   - Slice expressions (sub-slicing)
//   - How capacity growth works under the hood
//
// ENGINEERING DEPTH:
//   A Slice is purely a 24-byte struct (known as a Slice Header) containing:
//   a Pointer to an underlying contiguous array, an `int` for the Length, and an
//   `int` for the Capacity. This allows slices to be passed into functions incredibly
//   cheaply (copying just 24 bytes), while still allowing the function to mutate
//   the underlying shared memory. This is Go's secret to lightning-fast
//   data processing.
//
// RUN: go run ./03-collections-and-pointers/2-slices
// ============================================================================

func main() {

	// --- SLICE LITERAL ---
	// A slice looks like an array literal but WITHOUT a size in the brackets.
	// []string (slice) vs [3]string (array) — the absence of a number is key.
	//
	// Under the hood, a slice is a HEADER with 3 fields:
	//   1. Pointer — to the first element of the underlying array
	//   2. Length  — number of elements currently in the slice
	//   3. Capacity — total space allocated in the underlying array
	names := []string{"Alice", "John", "Mark"}
	fmt.Println("Names:", names)
	fmt.Printf("  len=%d, cap=%d\n", len(names), cap(names))

	// --- MAKE: Pre-allocate with length and capacity ---
	// Syntax: make([]Type, length, capacity)
	//
	// - length: how many elements exist RIGHT NOW (initialized to zero values)
	// - capacity: how much space is pre-allocated (avoids re-allocations)
	//
	// WHY PRE-ALLOCATE?
	// If you know you'll store ~1000 items, pre-allocate with make([]int, 0, 1000).
	// Without it, Go will re-allocate and copy the data every time the slice grows,
	// which is expensive for large collections.
	items := make([]int, 3, 5) // 3 elements exist (all 0), space for 5 total
	fmt.Printf("\nItems: %+v, Len: %d, Cap: %d\n", items, len(items), cap(items))

	// --- APPEND: Growing a slice ---
	// append() adds elements to the end of a slice.
	// IMPORTANT: append() returns a NEW slice header. Always reassign:
	//   items = append(items, value)   ← correct
	//   append(items, value)           ← WRONG! returned slice is discarded
	items = append(items, 1) // len=4, cap=5 (fits within capacity)
	items = append(items, 2) // len=5, cap=5 (fits within capacity)
	items = append(items, 3) // len=6, cap=10! Capacity DOUBLED

	// CAPACITY GROWTH: When append exceeds capacity, Go allocates a NEW
	// underlying array (roughly 2x the size) and copies all elements over.
	// This is O(n) when it happens, but amortized O(1) over many appends.
	items = append(items, 4) // len=7, cap=10 (still fits)
	fmt.Printf("Items: %+v, Len: %d, Cap: %d\n", items, len(items), cap(items))

	// --- SLICE EXPRESSIONS (Sub-slicing) ---
	// Syntax: slice[low:high]
	//   - low: starting index (inclusive)
	//   - high: ending index (exclusive)
	//   - Omit low = 0, Omit high = len(slice)
	//
	// CRITICAL: Sub-slices SHARE the same underlying array.
	// Modifying the sub-slice modifies the original!
	sub := items[3:7]
	fmt.Printf("\nSub-slice items[3:7]: %+v\n", sub)

	// KEY TAKEAWAY:
	// - Slices are dynamic, reference-based collections — USE THESE, not arrays
	// - The slice header has 3 fields: pointer, length, capacity
	// - Always reassign when using append: s = append(s, value)
	// - Pre-allocate with make() when you know the approximate size
	// - Sub-slices share memory — be careful with mutations
}
