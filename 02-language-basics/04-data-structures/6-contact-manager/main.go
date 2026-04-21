// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import "fmt"

// 04 Data Structures - Contact Directory (Exercise)
//
// Mental model:
// This milestone is intentionally plain. It proves slices, maps, and pointers
// directly before the curriculum moves on to functions, structs, and more
// layered design.
//
// Run: go run ./02-language-basics/04-data-structures/6-contact-manager

func main() {
	fmt.Println("=== Contact Directory ===")

	// Parallel slices keep the exercise inside 04 Data Structures concepts only.
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

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("SECTION COMPLETE: DS.6 contact-directory")
	fmt.Println("NEXT UP: FE.1 functions")
	fmt.Println("---------------------------------------------------")
}
