// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import "fmt"

// 04 Data Structures - Pointers
//
// Mental model:
// A pointer stores the address of a value. You use it when you need to update
// the original stored value instead of a copy.
//
// Run: go run ./02-language-basics/04-data-structures/4-pointers

func main() {
	fmt.Println("=== Pointers ===")

	score := 50
	scorePtr := &score

	fmt.Printf("score value:   %d\n", score)
	fmt.Printf("score address: %p\n", &score)
	fmt.Printf("scorePtr:      %p\n", scorePtr)
	fmt.Printf("*scorePtr:     %d\n", *scorePtr)

	// Changing a copied value does not affect the original.
	scoreCopy := score
	scoreCopy = 95
	fmt.Printf("\nAfter changing the copy: score=%d scoreCopy=%d\n", score, scoreCopy)

	// Changing through the pointer updates the original.
	*scorePtr = 95
	fmt.Printf("After changing through the pointer: score=%d\n", score)

	// Pointers work well with slice elements too.
	phones := []string{"111-2222", "333-4444", "555-6666"}
	bobPhone := &phones[1]
	*bobPhone = "333-9999"
	fmt.Printf("\nPhones after pointer update: %v\n", phones)

	var optionalScore *int
	if optionalScore == nil {
		fmt.Println("optionalScore is nil, so there is nothing to dereference yet.")
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: DS.5 slice-sharing")
	fmt.Println("Current: DS.4 (pointers)")
	fmt.Println("---------------------------------------------------")
}
