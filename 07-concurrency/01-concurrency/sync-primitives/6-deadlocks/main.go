package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 07: Concurrency - Deadlocks
//
// Run: go run ./07-concurrency/01-concurrency/sync-primitives/6-deadlocks

func main() {
	fmt.Println("=== SY.6 Deadlocks ===")
	fmt.Println("Learn how circular waits and unbalanced channel or lock usage can freeze a concurrent program.")
	fmt.Println()
	fmt.Println("- Keep lock ordering consistent.")
	fmt.Println("- Match every send with a reachable receiver and every wait with a reachable done path.")
	fmt.Println("- Prefer simple coordination patterns before stacking multiple synchronization tools together.")
	fmt.Println()
	fmt.Println("Deadlocks are design bugs. They disappear when ownership, lock ordering, and channel direction are made explicit.")
}
