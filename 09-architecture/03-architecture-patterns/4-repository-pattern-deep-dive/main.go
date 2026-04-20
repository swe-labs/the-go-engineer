package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 09: Architecture & Security - Repository pattern - deep dive
//
// Run: go run ./09-architecture/03-architecture-patterns/4-repository-pattern-deep-dive

func main() {
	fmt.Println("=== ARCH.4 Repository pattern - deep dive ===")
	fmt.Println("Learn what the repository pattern is for and where it becomes over-abstraction instead of useful design.")
	fmt.Println()
	fmt.Println("- Expose domain-oriented methods, not generic CRUD wrappers for everything.")
	fmt.Println("- Keep transactions and query intent clear at the boundary.")
	fmt.Println("- Do not add repository layers where direct storage calls would be simpler and clearer.")
	fmt.Println()
	fmt.Println("Repositories earn their keep when storage choices or complex mapping concerns would otherwise leak everywhere.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: ARCH.5")
	fmt.Println("Current: ARCH.4 (repository pattern - deep dive)")
	fmt.Println("---------------------------------------------------")
}
