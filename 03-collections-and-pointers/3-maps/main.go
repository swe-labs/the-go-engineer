// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

/*
=============================================================================
ENGINEERING DEPTH: Maps Internal Mechanics
=============================================================================
- TIME COMPLEXITY: $O(1)$ amortized lookup/insert.
- MEMORY MODEL: Under the hood, a Go map is a pointer to an `hmap` struct.
  It manages an array of buckets (each holding up to 8 key-value pairs).
- HASH COLLISIONS: If a bucket overflows, Go chains a new bucket to it.
- GROWING: When the load factor hits 6.5 (6.5 items per bucket), the map
  allocates a new underlying array double the size and incrementally migrates
  keys to avoid latency spikes (evacuation).
- CONCURRENCY: Maps are NOT thread-safe. A concurrent read/write will trigger
  a fatal panic. Use `sync.Mutex` or `sync.Map` for concurrent access.
=============================================================================
*/

func main() {

	// Idiomatic pattern: make(map[KeyType]ValueType, capacityHint)
	// Providing a capacity hint prevents expensive runtime re-allocations
	// when you know roughly how many items will be stored.
	// 1. Map Literal Initialization
	// This directly allocates the `hmap` struct and hashes the initial keys.
	studentGrades := map[string]int{
		"Alice": 90,
		"James": 85,
		"Dan":   60,
	}
	fmt.Printf("%+v\n", studentGrades)

	// 2. Map Assignment
	// Go hashes "Alice", finds her bucket, and overwrites the integer bytes.
	studentGrades["Alice"] = 95
	fmt.Printf("%+v\n", studentGrades)

	// 3. The "Comma-Ok" Idiom
	// Accessing a missing key returns the zero-value (0 for int).
	// To distinguish between "Alice scored a 0" and "Alice is not in the map",
	// Go returns a secondary boolean.
	alice, ok := studentGrades["Alice"]
	if ok {
		fmt.Printf("Alice: %+v\n", alice)
	}

	// 4. Inline Comma-Ok Check
	// This scopes the `value` and `ok` variables strictly to the `if` block,
	// preventing namespace pollution.
	key := "James"
	if value, ok := studentGrades[key]; ok {
		fmt.Printf("%s: %+v\n", key, value)
	}

	// 5. The Builtin delete() Function
	// This doesn't shrink the map's capacity! It merely flags the bucket
	// slot as "empty", keeping the memory allocated for future inserts.
	delete(studentGrades, "Alice")

	fmt.Printf("%+v\n", studentGrades)

	configs := make(map[string]int)
	fmt.Printf("%+v %T\n", configs, configs)

}
