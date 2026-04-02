// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 3: Collections & Pointers — Arrays
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Arrays are FIXED-SIZE, value-type collections
//   - Declaring arrays with explicit size
//   - Array literals and zero values
//   - Why arrays are rarely used directly (slices are preferred)
//   - Multi-dimensional arrays
//
// ENGINEERING DEPTH:
//   An Array in Go is a contiguous block of memory allocated strictly at compile
//   time. Because the size is baked into the type itself (`[3]int`), the Go scheduler
//   can optimize memory layout on the Stack rather than the Heap. However, because
//   they are "Value Types", passing a `[1000]int` array to a function copies all
//   1000 integers sequentially in RAM—which is why Slices are used 99% of the time.
//
// RUN: go run ./03-collections-and-pointers/1-array
// ============================================================================

func main() {

	// --- DECLARING AN ARRAY ---
	// Syntax: var name [SIZE]Type
	//
	// The SIZE is part of the TYPE. [2]int and [3]int are DIFFERENT TYPES.
	// You cannot assign a [2]int to a [3]int variable.
	//
	// Zero values: all elements are initialized to the type's zero value.
	// For int, that's 0. For string, that's "".
	var numbers [2]int
	fmt.Printf("Zero value array: %+v\n", numbers) // [0 0]

	// Set individual elements by index (0-based)
	numbers[0] = 1
	numbers[1] = 2
	fmt.Printf("After assignment: %+v\n", numbers) // [1 2]

	// --- ARRAY LITERALS ---
	// Initialize with values directly. Go counts the elements for you.
	primes := [4]int{2, 3, 5, 7}
	fmt.Printf("Primes: %+v\n", primes)

	// You can modify individual elements
	primes[3] = 11                        // Replace 7 with 11
	fmt.Printf("Modified: %+v\n", primes) // [2 3 5 11]

	// --- ITERATING AN ARRAY ---
	// Method 1: Classic for loop with len()
	// len(array) returns the number of elements.
	fmt.Println("\nIterating with index:")
	for i := 0; i < len(primes); i++ {
		fmt.Printf("  primes[%d] = %d\n", i, primes[i])
	}

	// Method 2: range (preferred — cleaner and safer)
	fmt.Println("\nIterating with range:")
	for i, v := range primes {
		fmt.Printf("  primes[%d] = %d\n", i, v)
	}

	// --- MULTI-DIMENSIONAL ARRAYS ---
	// Arrays of arrays — useful for grids, matrices, game boards.
	// [2][3]int = 2 rows, 3 columns.
	var matrix [2][3]int
	matrix[0][0] = 1
	matrix[0][1] = 2
	matrix[1][2] = 3

	fmt.Printf("\nMatrix: %+v\n", matrix) // [[1 2 0] [0 0 3]]

	// --- ARRAYS ARE VALUE-TYPES ---
	// When you assign an array to another variable, Go makes a FULL COPY.
	// Modifying the copy does NOT affect the original.
	original := [3]int{10, 20, 30}
	copied := original // Full copy — NOT a reference
	copied[0] = 99
	fmt.Printf("\nOriginal: %v (unchanged)\n", original) // [10 20 30]
	fmt.Printf("Copy:     %v (modified)\n", copied)      // [99 20 30]

	// KEY TAKEAWAY:
	// - Arrays have FIXED size — the size is part of the type
	// - Arrays are VALUE TYPES — assignment creates a full copy
	// - In practice, Go programmers almost always use SLICES (next lesson)
	//   because slices are dynamic and pass by reference to the underlying data
}
