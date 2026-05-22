// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: The slices package
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Use slices.Sort to sort any slice of ordered types.
//   - Use slices.Compact to remove adjacent duplicates.
//   - Use slices.Contains and slices.Index for searching.
//   - Use slices.Delete and slices.Insert to modify slices efficiently.
//
// WHY THIS MATTERS:
//   Before Go 1.21, sorting slices required the reflect-heavy sort.Slice or
//   writing boilerplate. The slices package provides generic, type-safe
//   operations that are faster and cleaner.
//
// RUN:
//   go run ./02-language-basics/04-data-structures/07-slices
//
// KEY TAKEAWAY:
//   - The slices package is the modern Go way to search, sort, and manipulate slices.
//   - All functions are generic and work with any comparable or ordered type.
// ============================================================================

package main

import (
	"fmt"
	"slices"
)

func main() {
	// 1. Sorting — slices.Sort modifies the slice in place.
	nums := []int{42, 7, 13, 99, 1, 7, 42}
	fmt.Println("Before sort:", nums)
	slices.Sort(nums)
	fmt.Println("After sort :", nums)
	fmt.Println()

	// 2. Deduplication — Compact removes adjacent duplicates (sort first).
	uniq := slices.Compact(nums)
	fmt.Println("Compact    :", uniq)
	fmt.Println()

	// 3. Searching — BinarySearch requires a sorted slice.
	pos, found := slices.BinarySearch(uniq, 42)
	fmt.Printf("BinarySearch(42): found=%v position=%d\n", found, pos)

	pos, found = slices.BinarySearch(uniq, 50)
	fmt.Printf("BinarySearch(50): found=%v insertAt=%d\n", found, pos)
	fmt.Println()

	// 4. Contains and Index for linear search.
	words := []string{"apple", "banana", "cherry", "date"}
	fmt.Printf("Contains(banana): %v\n", slices.Contains(words, "banana"))
	fmt.Printf("Index(date)    : %d\n", slices.Index(words, "date"))
	fmt.Printf("Contains(zebra): %v\n", slices.Contains(words, "zebra"))
	fmt.Println()

	// 5. Delete and Insert for slice modification.
	// Delete removes elements i..j-1. Insert grows the slice.
	clipped := slices.Delete(words, 1, 3)
	fmt.Println("After Delete(1,3):", clipped)
	extended := slices.Insert(clipped, 1, "blueberry", "cranberry")
	fmt.Println("After Insert     :", extended)
	fmt.Println()

	// 6. Clone creates an independent copy.
	original := []int{1, 2, 3}
	copy := slices.Clone(original)
	copy[0] = 99
	fmt.Printf("Clone test — original[0]=%d copy[0]=%d (independent)\n", original[0], copy[0])
	fmt.Println()

	// 7. Equal compares element-by-element.
	a := []int{1, 2, 3}
	b := []int{1, 2, 3}
	c := []int{1, 2, 4}
	fmt.Printf("Equal(a, b): %v\n", slices.Equal(a, b))
	fmt.Printf("Equal(a, c): %v\n", slices.Equal(a, c))

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: DS.8 -> 02-language-basics/04-data-structures/08-maps")
	fmt.Println("Run    : go run ./02-language-basics/04-data-structures/08-maps")
	fmt.Println("Current: DS.7 (slices)")
	fmt.Println("---------------------------------------------------")
}
