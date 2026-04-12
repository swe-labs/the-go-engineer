// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 03: Data Structures - Contact Directory (Exercise)
// Level: Beginner -> Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Combining slices, maps, and pointers in one small program
//   - Using a map to find slice positions quickly
//   - Taking a pointer to a slice element so an update persists
//   - Keeping the exercise inside Section 03 concepts only
//
// ENGINEERING DEPTH:
//   This exercise intentionally does NOT use helper functions, methods, or a
//   struct-heavy design yet. The goal is to prove today's tools directly:
//   slices for ordered storage, a map for O(1) lookup of slice positions, and
//   pointers for updates that must stick.
//
// RUN: go run ./03-data-structures/6-contact-manager
// ============================================================================

func main() {
	fmt.Println("=== Contact Directory ===")

	// Parallel slices keep the exercise inside Section 03 concepts only.
	// Each contact uses the same index across these three slices.
	names := make([]string, 0, 4)
	emails := make([]string, 0, 4)
	phones := make([]string, 0, 4)

	// The map gives us O(1) lookup from a name to the slice index.
	indexByName := make(map[string]int)

	// Add Alice.
	names = append(names, "Alice Wonderland")
	emails = append(emails, "alice@example.com")
	phones = append(phones, "111-2222")
	indexByName["Alice Wonderland"] = len(names) - 1

	// Add Bob.
	names = append(names, "Bob The Builder")
	emails = append(emails, "bob@example.com")
	phones = append(phones, "333-4444")
	indexByName["Bob The Builder"] = len(names) - 1

	// Add Charlie.
	names = append(names, "Charlie Brown")
	emails = append(emails, "charlie@example.com")
	phones = append(phones, "555-6666")
	indexByName["Charlie Brown"] = len(names) - 1

	// Duplicate check stays simple: if the key exists, skip the add.
	if _, exists := indexByName["Alice Wonderland"]; exists {
		fmt.Println("Duplicate add skipped for Alice Wonderland.")
	}

	fmt.Println("\n--- Listing Contacts ---")
	for i := 0; i < len(names); i++ {
		fmt.Printf("%d. %s | %s | %s\n", i+1, names[i], emails[i], phones[i])
	}

	// Lookup with the map, then take a pointer to the stored phone number.
	fmt.Println("\n--- Lookup and Update ---")
	bobIndex, ok := indexByName["Bob The Builder"]
	if !ok {
		fmt.Println("Bob lookup failed.")
		return
	}

	phonePtr := &phones[bobIndex]
	fmt.Printf("Found Bob at index %d with phone %s\n", bobIndex, *phonePtr)

	*phonePtr = "333-9999"
	fmt.Printf("Updated Bob through pointer: %s\n", *phonePtr)

	// Re-read the stored slice value to prove the update persisted.
	updatedBobIndex := indexByName["Bob The Builder"]
	fmt.Printf("Persisted Bob phone: %s\n", phones[updatedBobIndex])

	// Missing lookup still uses the comma-ok pattern from maps.
	if _, exists := indexByName["Zack"]; !exists {
		fmt.Println("Zack not found.")
	}

	// KEY TAKEAWAY:
	// - Slices hold the ordered contact data.
	// - The map tells us where each contact lives in those slices.
	// - A pointer to a slice element lets an update persist.
	// - This milestone stays intentionally simple so the data-structure choices stay visible.
}
