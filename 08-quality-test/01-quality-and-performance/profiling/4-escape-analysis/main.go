package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 08: Quality & Testing - Escape analysis
//
// Run: go run ./08-quality-test/01-quality-and-performance/profiling/4-escape-analysis

func main() {
	fmt.Println("=== PR.4 Escape analysis ===")
	fmt.Println("Learn how the compiler decides whether values stay on the stack or escape to the heap.")
	fmt.Println()
	fmt.Println("- Stack values are cheaper to allocate and reclaim.")
	fmt.Println("- Escapes often happen because a value must outlive the current function.")
	fmt.Println("- Compiler diagnostics explain where heap pressure starts.")
	fmt.Println()
	fmt.Println("Escape analysis explains many allocation surprises, especially in helper-heavy code or tight loops.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: PR.5")
	fmt.Println("Current: PR.4 (escape analysis)")
	fmt.Println("---------------------------------------------------")
}
