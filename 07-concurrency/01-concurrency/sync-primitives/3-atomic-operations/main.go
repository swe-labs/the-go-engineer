package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 07: Concurrency - Atomic operations
//
// Run: go run ./07-concurrency/01-concurrency/sync-primitives/3-atomic-operations

func main() {
	fmt.Println("=== SY.3 Atomic operations ===")
	fmt.Println("Learn when lock-free atomic reads and writes are the right tool for counters, flags, and single-word state.")
	fmt.Println()
	fmt.Println("- Use atomics for flags and counters, not arbitrary structs.")
	fmt.Println("- Compare-and-swap coordinates changes around one current value.")
	fmt.Println("- Prefer clarity over clever lock-free code when the workload does not demand it.")
	fmt.Println()
	fmt.Println("Atomics are precise tools. If the state spans multiple fields or invariants, a mutex is usually the clearer choice.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: SY.4")
	fmt.Println("Current: SY.3 (atomic operations)")
	fmt.Println("---------------------------------------------------")
}
