package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 07: Concurrency - sync.Once and sync.Map
//
// Run: go run ./07-concurrency/01-concurrency/sync-primitives/2-once-and-sync-map

func main() {
	fmt.Println("=== SY.2 sync.Once and sync.Map ===")
	fmt.Println("Learn the standard-library helpers for one-time initialization and specific concurrent map workloads.")
	fmt.Println()
	fmt.Println("- sync.Once protects one-time initialization.")
	fmt.Println("- sync.Map is useful for read-heavy or disjoint-key workloads.")
	fmt.Println("- Regular maps plus explicit locks are still the default choice.")
	fmt.Println()
	fmt.Println("Pick these tools for the workload they are meant for; overusing them usually hides a simpler design.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: SY.3")
	fmt.Println("Current: SY.2 (sync.once and sync.map)")
	fmt.Println("---------------------------------------------------")
}
