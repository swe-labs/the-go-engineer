package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 03: Functions and Errors - Closures - mechanics
//
// Run: go run ./03-functions-errors/9-closures-mechanics

func main() {
	fmt.Println("=== FE.9 Closures - mechanics ===")
	fmt.Println("Learn how closures capture variables, why that extends lifetimes, and where the loop-variable trap comes from.")
	fmt.Println()
	fmt.Println("- Closures capture variables from outer scope.")
	fmt.Println("- Captured state stays alive as long as the closure can still use it.")
	fmt.Println("- Loop variables must be rebound when each closure needs its own copy.")
	fmt.Println()
	fmt.Println("Closure bugs are usually state bugs. The most common one is reusing the same loop variable across multiple callbacks or goroutines.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: FE.10")
	fmt.Println("Current: FE.9 (closures - mechanics)")
	fmt.Println("---------------------------------------------------")
}
