// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 09: Architecture & Security
// Title: Domain-Driven Design basics
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn the boundary language behind domains, aggregates, and bounded contexts.
//
// WHY THIS MATTERS:
//   - DDD is mainly a naming and boundary discipline, not a folder structure gimmick.
//
// RUN:
//   go run ./09-architecture/03-architecture-patterns/2-ddd-basics
//
// KEY TAKEAWAY:
//   - Learn the boundary language behind domains, aggregates, and bounded contexts.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== ARCH.2 Domain-Driven Design basics ===")
	fmt.Println("Learn the boundary language behind domains, aggregates, and bounded contexts.")
	fmt.Println()
	fmt.Println("- Bounded contexts separate meanings that look similar but behave differently.")
	fmt.Println("- Entities, value objects, and aggregates organize domain rules.")
	fmt.Println("- Ubiquitous language keeps code and business conversations aligned.")
	fmt.Println()
	fmt.Println("DDD becomes useful when the domain is complex enough that language mismatches create real implementation mistakes.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: ARCH.3")
	fmt.Println("Current: ARCH.2 (domain-driven design basics)")
	fmt.Println("---------------------------------------------------")
}
