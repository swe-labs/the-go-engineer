// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: Arrays
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to declare and initialize fixed-size arrays.
//   - Why the array size is part of its type.
//   - Array "value-copy" behavior (pass-by-value).
//
// WHY THIS MATTERS:
//   - Arrays are the low-level foundation for memory storage in Go. Even
//     though you'll use slices more often, understanding arrays is key to
//     mastering memory layout and avoiding accidental data duplication.
//
// RUN:
//   go run ./02-language-basics/04-data-structures/1-array
//
// KEY TAKEAWAY:
//   - Copying an array copies every single element.
// ============================================================================

package main

import "fmt"

// 04 Data Structures - Arrays
//
// Mental model:
// Arrays are fixed-size values. If you copy one array into another, Go copies
// all of the elements.
//

func main() {
	fmt.Println("=== Arrays ===")

	// The size is part of the array type.
	var numbers [2]int
	fmt.Printf("Zero value array: %v\n", numbers)

	numbers[0] = 1
	numbers[1] = 2
	fmt.Printf("After assignment: %v\n", numbers)

	primes := [4]int{2, 3, 5, 7}
	fmt.Printf("\nLiteral: %v\n", primes)

	fmt.Println("Range iteration:")
	for i, value := range primes {
		fmt.Printf("  primes[%d] = %d\n", i, value)
	}

	original := [3]int{10, 20, 30}
	copied := original
	copied[0] = 99

	fmt.Printf("\nOriginal: %v\n", original)
	fmt.Printf("Copy:     %v\n", copied)

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: DS.2 -> 02-language-basics/04-data-structures/2-slices")
	fmt.Println("Run    : go run ./02-language-basics/04-data-structures/2-slices")
	fmt.Println("Current: DS.1 (arrays)")
	fmt.Println("---------------------------------------------------")
}
