package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 09: Architecture & Security - Input validation patterns
//
// Run: go run ./09-architecture/04-security/1-input-validation-patterns

func main() {
	fmt.Println("=== SEC.1 Input validation patterns ===")
	fmt.Println("Learn how allow-lists, normalization, and fail-fast checks turn raw input into trustworthy domain values.")
	fmt.Println()
	fmt.Println("- Validate early and explicitly.")
	fmt.Println("- Normalize data before deeper rules depend on it.")
	fmt.Println("- Prefer allow-lists to deny-lists when the acceptable shape is known.")
	fmt.Println()
	fmt.Println("Treat validation as a first-class engineering concern because every public boundary is a security boundary.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: SEC.2")
	fmt.Println("Current: SEC.1 (input validation patterns)")
	fmt.Println("---------------------------------------------------")
}
