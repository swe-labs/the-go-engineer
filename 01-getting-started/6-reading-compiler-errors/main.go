// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 01: Getting Started — Lesson GT.6: Reading Compiler Errors
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to read Go compiler error messages
//   - How to locate the line causing the error
//   - Common error types like "unused variable" or "syntax error"
//
// WHY THIS MATTERS:
//   The compiler is your first line of defense against bugs.
//
// RUN: go run ./01-getting-started/6-reading-compiler-errors
// ============================================================================

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

	// KEY TAKEAWAY:
	// - Compiler errors are your friends, not enemies.
	// - Go enforces clean code through strict compilation rules.

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: LB.1 variables")
	fmt.Println("Current: GT.6 (reading-compiler-errors)")
	fmt.Println("-----------------------------------")
}
