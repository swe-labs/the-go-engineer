// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 07: Concurrency
// Title: Deadlocks
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how circular waits and unbalanced channel or lock usage can freeze a concurrent program.
//
// WHY THIS MATTERS:
//   - Deadlocks happen when each waiting actor needs another waiting actor to move first.
//
// RUN:
//   go run ./07-concurrency/01-concurrency/sync-primitives/6-deadlocks
//
// KEY TAKEAWAY:
//   - [TODO: Summarize the core takeaway]
// ============================================================================

package main

import "fmt"

//

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

// ---------------------------------------------------
// NEXT UP: CT.1
// ---------------------------------------------------
