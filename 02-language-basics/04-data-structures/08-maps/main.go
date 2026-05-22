// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: The maps package
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Use maps.Clone to create an independent copy of a map.
//   - Use maps.Copy to merge keys from one map into another.
//   - Use maps.Keys and maps.Values to extract map contents as slices.
//   - Use maps.DeleteFunc to remove entries conditionally.
//
// WHY THIS MATTERS:
//   Before Go 1.21, copying a map required a manual loop. The maps package
//   provides generic, type-safe operations that eliminate boilerplate.
//
// RUN:
//   go run ./02-language-basics/04-data-structures/08-maps
//
// KEY TAKEAWAY:
//   - The maps package provides generic operations for all map types.
//   - Use maps.Clone for independent copies and maps.Copy for merging.
//   - Keys() and Values() are not ordered — sort if you need stable output.
// ============================================================================

package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	// 1. Creating and cloning.
	scores := map[string]int{"Alice": 95, "Bob": 87, "Charlie": 92}
	backup := maps.Clone(scores)
	fmt.Println("Original:", scores)
	fmt.Println("Clone   :", backup)

	// Modifying the original doesn't affect the clone.
	scores["David"] = 78
	fmt.Println("After add: original has David:", scores["David"] > 0, "clone does not:", backup["David"] > 0)
	fmt.Println()

	// 2. Copy merges all keys from src into dst.
	moreScores := map[string]int{"Eve": 88, "Frank": 91}
	maps.Copy(scores, moreScores)
	fmt.Println("After Copy:", scores)
	fmt.Println()

	// 3. Keys and Values extract the contents (returns iterators in Go 1.23+).
	names := slices.Collect(maps.Keys(scores))
	allScores := slices.Collect(maps.Values(scores))
	slices.Sort(names)
	slices.Sort(allScores)
	fmt.Println("Sorted keys  :", names)
	fmt.Println("Sorted values:", allScores)
	fmt.Println()

	// 4. DeleteFunc removes entries where the callback returns true.
	honorRoll := maps.Clone(scores)
	maps.DeleteFunc(honorRoll, func(_ string, v int) bool { return v < 90 })
	fmt.Println("Honor roll (score >= 90):", honorRoll)
	fmt.Println()

	// 5. Equal checks if two maps have identical key-value pairs.
	fmt.Printf("Equal(scores, backup): %v\n", maps.Equal(scores, backup))
	equalCopy := maps.Clone(scores)
	fmt.Printf("Equal(scores, clone) : %v\n", maps.Equal(scores, equalCopy))

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: FE.1 -> 03-functions-errors/01-functions-basics")
	fmt.Println("Run    : go run ./03-functions-errors/01-functions-basics")
	fmt.Println("Current: DS.8 (maps)")
	fmt.Println("---------------------------------------------------")
}
