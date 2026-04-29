// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: Maps
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how Go performs keyed lookup with maps and why the comma-ok pattern matters whenever a missing key would otherwise be ambiguous.
//
// WHY THIS MATTERS:
//   - A map connects keys to values. Use it when finding something by name, ID, or label matters more than keeping items in order.
//
// RUN:
//   go run ./02-language-basics/04-data-structures/3-maps
//
// KEY TAKEAWAY:
//   - Learn how Go performs keyed lookup with maps and why the comma-ok pattern matters whenever a missing key would otherwise be ambiguous.
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

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: DS.4 pointers")
	fmt.Println("Current: DS.3 (maps)")
	fmt.Println("Previous: DS.2 (slices)")
	fmt.Println("---------------------------------------------------")
}
