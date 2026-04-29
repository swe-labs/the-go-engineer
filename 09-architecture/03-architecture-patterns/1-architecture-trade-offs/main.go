// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 09: Architecture & Security
// Title: Monolith vs Modular Monolith vs Microservices
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Compare the default service shapes teams choose as systems and organizations grow.
//
// WHY THIS MATTERS:
//   - Architecture is a cost allocation decision. Every boundary buys one kind of freedom by adding another kind of overhead.
//
// RUN:
//   go run ./09-architecture/03-architecture-patterns/1-architecture-trade-offs
//
// KEY TAKEAWAY:
//   - Compare the default service shapes teams choose as systems and organizations grow.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== ARCH.1 Monolith vs Modular Monolith vs Microservices ===")
	fmt.Println("Compare the default service shapes teams choose as systems and organizations grow.")
	fmt.Println()
	fmt.Println("- Monolith keeps change local.")
	fmt.Println("- Modular monolith keeps one deploy while improving internal boundaries.")
	fmt.Println("- Microservices add independence but also network, deployment, and coordination cost.")
	fmt.Println()
	fmt.Println("The right default is the one that keeps change cheap for your current team and system size.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: ARCH.2")
	fmt.Println("Current: ARCH.1 (monolith vs modular monolith vs microservices)")
	fmt.Println("---------------------------------------------------")
}
