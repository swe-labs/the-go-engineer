// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 01: Getting Started
// Title: Reading compiler errors
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn to treat the compiler as a helpful partner instead of an obstacle by decoding its error messages.
//
// WHY THIS MATTERS:
//   - Think of the compiler as a "spell checker" for logic. It won't let you run a program that it knows will crash or behave unpredictably.
//
// RUN:
//   go run ./01-getting-started/6-reading-compiler-errors
//
// KEY TAKEAWAY:
//   - Learn to treat the compiler as a helpful partner instead of an obstacle by decoding its error messages.
// ============================================================================

//   The compiler is your first line of defense against bugs.

package main

import (
	"fmt"
)

func main() {
	fmt.Println("GT.6: Reading Compiler Errors")
	fmt.Println("--------------------------------")

	fmt.Println("To learn from the compiler:")
	fmt.Println("1. Look at the filename and line number.")
	fmt.Println("2. Read the error message carefully.")
	fmt.Println("3. Fix the code and run again.")

	// - Compiler errors are your friends, not enemies.
	// - Go enforces clean code through strict compilation rules.

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: LB.1 variables")
	fmt.Println("Current: GT.6 (reading-compiler-errors)")
	fmt.Println("-----------------------------------")
}
