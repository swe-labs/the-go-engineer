package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 09: Architecture & Security - CQRS basics
//
// Run: go run ./09-architecture/03-architecture-patterns/7-cqrs-basics

func main() {
	fmt.Println("=== ARCH.7 CQRS basics ===")
	fmt.Println("Learn when separating write models from read models improves a system and when it is needless complexity.")
	fmt.Println()
	fmt.Println("- One model does not always fit both commands and queries.")
	fmt.Println("- Read models often optimize different access patterns than write models.")
	fmt.Println("- CQRS raises consistency and synchronization questions that simpler systems avoid.")
	fmt.Println()
	fmt.Println("CQRS is a useful tool only when the read/write mismatch is real enough to justify the extra moving parts.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: ARCH.8")
	fmt.Println("Current: ARCH.7 (cqrs basics)")
	fmt.Println("---------------------------------------------------")
}
