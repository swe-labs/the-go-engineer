// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 03: Data Structures — Contact Manager (Exercise)
// Level: Beginner → Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Synthesizing structs, slices, maps, and pointers into a real application
//   - Using `init()` functions for setup
//   - Indexing slices using a map for O(1) lookups
//   - Returning pointers from functions to allow nil-returns on failure
//
// ENGINEERING DEPTH:
//   This exercise demonstrates a fundamental backend architectural pattern:
//   Secondary Indexing. Searching an array/slice for a specific user takes
//   O(n) linear time. By creating a parallel `map[string]int` that tracks
//   exactly which slice index a name lives at, we achieve O(1) instant
//   lookups. This is precisely how database table indexes work under the hood!
//
// RUN: go run ./03-data-structures/6-contact-manager
// ============================================================================

// Contact structure. Maps real-world data to a Go object.
type Contact struct {
	ID    int
	Name  string
	Email string
	Phone string
}

// --- GLOBAL STATE ---
// `contactList` holds the actual Contact data contiguously in memory.
var contactList []Contact

// `contactIndexByName` is our Secondary Index.
// It maps a Name (string) directly to the integer index in the `contactList` slice.
var contactIndexByName map[string]int

// `nextID` tracks the auto-incrementing primary key for new contacts.
var nextID = 1

// --- THE INIT FUNCTION ---
// The `init()` function is a special Go feature.
// It executes automatically exactly ONCE per package, before `main()` begins.
// It is perfectly designed for setting up memory allocations, config parsing, or
// establishing database connections.
//
// Warning: Misusing `init()` for complex logic makes code hard to test.
// Use it strictly for basic memory/system initialization.
func init() {
	// Pre-allocate the slice with a length of 0.
	contactList = make([]Contact, 0)

	// Maps MUST be initialized with `make` before use.
	// Writing to an uninitialized (nil) map causes a runtime panic.
	contactIndexByName = make(map[string]int)
}

// addContact inserts a new record into our mini-database.
func addContact(name, email, phone string) {

	// 1. O(1) Check: Does this contact already exist?
	// The comma-ok idiom safely checks the map without triggering panics.
	if _, exists := contactIndexByName[name]; exists {
		fmt.Printf("❌ Contact already exists: %s\n", name)
		return
	}

	// 2. Create the struct
	newContact := Contact{
		ID:    nextID,
		Name:  name,
		Email: email,
		Phone: phone,
	}
	nextID++ // Increment ID for the next user

	// 3. Append to the master list.
	// Remember: append returns a new slice header, so we MUST reassign it!
	contactList = append(contactList, newContact)

	// 4. Update the secondary index.
	// The new contact is at the very end of the slice, so its index is length - 1.
	contactIndexByName[name] = len(contactList) - 1

	fmt.Printf("✅ Contact added: %s\n", name)
}

// findContact retrieves a Contact by name.
//
// RETURN TYPE: *Contact (Pointer to a Contact)
// Why return a pointer instead of a value?
// Because if the user doesn't exist, we can return `nil`. If we returned a value,
// we would have to return an empty `Contact{}` struct, which the caller might
// mistakenly think is a real user! Returning nil allows for a clean `if c == nil` check.
func findContact(name string) *Contact {
	// Instantly lookup the slice index using our map.
	index, exists := contactIndexByName[name]
	if exists {
		// Return the memory address (&) of the contact inside the slice.
		return &contactList[index]
	}
	return nil
}

// ListContacts iterates across the entire memory layout.
func ListContacts() {
	fmt.Println("\n--- Listing Contacts ---")

	if len(contactList) == 0 {
		fmt.Println("No contacts found.")
		return
	}

	// Iterate using `range`. `i` is the array index, `contact` is the copied struct.
	for i, contact := range contactList {
		fmt.Printf("%d. ID: %d | Name: %s | Email: %s | Phone: %s\n",
			i+1, contact.ID, contact.Name, contact.Email, contact.Phone)
	}
	fmt.Println("------------------------")
}

func main() {
	fmt.Println("=== Starting Contact Manager ===")

	// Adding initial contacts
	addContact("Alice Wonderland", "alice@example.com", "111-2222")
	addContact("Bob The Builder", "bob@example.com", "333-4444")
	addContact("Charlie Brown", "charlie@example.com", "555-6666")

	// Try to add a duplicate to prove the map logic works.
	addContact("Alice Wonderland", "alice.new@example.com", "777-8888")

	// Print the database
	ListContacts()

	// Perform a lookup
	fmt.Println("--- Lookup Testing ---")
	bob := findContact("Bob The Builder")
	if bob == nil {
		fmt.Println("❌ Bob lookup failed.")
	} else {
		// Notice how we use `bob.Name` directly?
		// Go automatically dereferences pointers when accessing struct fields.
		// You don't have to write `(*bob).Name` like you do in C.
		fmt.Printf("✅ Found Bob: ID %d, Phone: %s\n", bob.ID, bob.Phone)

		// Update Bob through the returned pointer so the change persists.
		bob.Phone = "333-9999"
		fmt.Printf("✅ Updated Bob's phone to: %s\n", bob.Phone)
	}

	// Lookup a missing user exactly
	missing := findContact("Zack")
	if missing == nil {
		fmt.Println("❌ Zack not found (expected nil return).")
	}

	// Re-check Bob to prove the update changed the stored slice entry.
	updatedBob := findContact("Bob The Builder")
	if updatedBob != nil {
		fmt.Printf("✅ Persisted Bob update: ID %d, Phone: %s\n", updatedBob.ID, updatedBob.Phone)
	}

	// KEY TAKEAWAY:
	// - `init()` runs automatically before main for setups
	// - Combining a Slice (data storage) with a Map (index lookup) is incredibly fast
	// - Returning pointers lets updates persist without copying the stored value
	// - Go auto-dereferences pointer structs (e.g. ptr.Field instead of (*ptr).Field)
}
