// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 09: Architecture & Security
// Title: When to split services
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn the signals that justify service boundaries and the anti-signals that only imitate architecture maturity.
//
// WHY THIS MATTERS:
//   - Split services because one boundary is costing too much, not because distributed systems sound advanced.
//
// RUN:
//   go run ./09-architecture/03-architecture-patterns/8-when-to-split-services
//
// KEY TAKEAWAY:
//   - Learn the signals that justify service boundaries and the anti-signals that only imitate architecture maturity.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== ARCH.8 When to split services ===")
	fmt.Println("Learn the signals that justify service boundaries and the anti-signals that only imitate architecture maturity.")
	fmt.Println()
	fmt.Println("- Split when the current boundary is the problem, not when the word 'microservice' sounds appealing.")
	fmt.Println("- Look for sustained pressure in deploy cadence, ownership, or scaling needs.")
	fmt.Println("- Every split creates new network and operational cost immediately.")
	fmt.Println()
	fmt.Println("The best split is one backed by clear pressure: team autonomy, scaling profile, fault isolation, or compliance.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: ARCH.9 -> 09-architecture/03-architecture-patterns/9-modular-refactor-exercise")
	fmt.Println("Current: ARCH.8 (when to split services)")
	fmt.Println("---------------------------------------------------")
}
