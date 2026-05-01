// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: Maps
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Associative data storage using keys and values.
//   - The "comma-ok" idiom for safe lookups.
//   - Zero-value behavior for missing keys.
//   - Using `make` vs. literals for map initialization.
//
// WHY THIS MATTERS:
//   - While slices are great for order, maps are unbeatable for speed. A map
//     allows you to find a single item among millions in almost constant time
//     (O(1)). In production, maps are used for caches, registries, and
//     configuration lookup tables.
//
// RUN:
//   go run ./02-language-basics/04-data-structures/3-maps
//
// KEY TAKEAWAY:
//   - Maps provide fast lookup; the comma-ok pattern provides lookup safety.
// ============================================================================

package main

import "fmt"

// 04 Data Structures - Maps
//
// Mental model:
// A map connects keys to values. You use it when lookup by name matters more
// than keeping items in order.
//

func main() {
	fmt.Println("=== Maps ===")

	studentGrades := map[string]int{
		"Alice": 90,
		"James": 85,
		"Dan":   60,
	}
	fmt.Printf("initial: %v\n", studentGrades)

	studentGrades["Alice"] = 95
	studentGrades["Mary"] = 88
	fmt.Printf("after updates: %v\n", studentGrades)

	// Missing keys return the zero value for the value type.
	fmt.Printf("\nMissing score without comma-ok: %d\n", studentGrades["Zack"])

	aliceScore, aliceExists := studentGrades["Alice"]
	fmt.Printf("Alice exists? %v, score: %d\n", aliceExists, aliceScore)

	zackScore, zackExists := studentGrades["Zack"]
	fmt.Printf("Zack exists? %v, score: %d\n", zackExists, zackScore)

	delete(studentGrades, "Dan")
	fmt.Printf("\nafter delete(\"Dan\"): %v\n", studentGrades)

	settings := make(map[string]string)
	settings["theme"] = "dark"
	settings["timezone"] = "UTC"
	fmt.Printf("settings: %v\n", settings)

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: DS.4 -> 02-language-basics/04-data-structures/4-pointers")
	fmt.Println("Run    : go run ./02-language-basics/04-data-structures/4-pointers")
	fmt.Println("Current: DS.3 (maps)")
	fmt.Println("---------------------------------------------------")
}
