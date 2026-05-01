// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 01: Getting Started
// Title: Hello World
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn the smallest useful shape of an executable Go program.
//
// WHY THIS MATTERS:
//   - A minimal Go executable has a stable shape: 1. Declare the package. 2. Import what the file needs. 3. Define `main`. 4. Execute statements inside `...
//
// RUN:
//   go run ./01-getting-started/2-hello-world
//
// KEY TAKEAWAY:
//   - Learn the smallest useful shape of an executable Go program.
// ============================================================================

package main

import "fmt"

func main() {
	fmt.Println("Hello, World! Welcome to The Go Engineer.")
	fmt.Println("Go was created at", "Google", "in", 2009)

	language := "Go"
	year := 2009
	fmt.Printf("%s was created in %d\n", language, year)

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: GT.3 -> 01-getting-started/3-how-go-works")
	fmt.Println("Current: GT.2 (hello-world)")
	fmt.Println("---------------------------------------------------")
}
