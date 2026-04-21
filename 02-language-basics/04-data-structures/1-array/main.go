// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import "fmt"

// 04 Data Structures - Arrays
//
// Mental model:
// Arrays are fixed-size values. If you copy one array into another, Go copies
// all of the elements.
//
// Run: go run ./02-language-basics/04-data-structures/1-array

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

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: DS.2 slices")
	fmt.Println("Current: DS.1 (arrays)")
	fmt.Println("---------------------------------------------------")
}
