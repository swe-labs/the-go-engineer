package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 09: Architecture & Security - Domain-Driven Design basics
//
// Run: go run ./09-architecture/03-architecture-patterns/2-ddd-basics

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
