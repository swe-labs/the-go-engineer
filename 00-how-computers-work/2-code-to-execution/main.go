// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 00: How Computers Work — How Code Becomes Execution
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Source code passes through multiple compiler stages
//   - The compiler builds machine-friendly representations before machine code
//   - The final binary is what the OS actually launches
//
// WHY THIS MATTERS:
//   Compile-time reasoning explains why Go catches type errors early and why
//   binaries can be deployed as concrete build artifacts.
//
// RUN: go run ./00-how-computers-work/2-code-to-execution
// ============================================================================

package main

import "fmt"

func main() {
	fmt.Println("Source code : x := 42 + y")
	fmt.Println("Tokens      : IDENT(x) DEFINE(:=) NUMBER(42) PLUS(+) IDENT(y)")
	fmt.Println("AST         : Assign(x, Add(42, y))")
	fmt.Println("IR          : t0 = load y; t1 = add 42, t0; store x, t1")
	fmt.Println("Binary      : machine instructions + linked packages")
	fmt.Println("Runtime     : the OS loads the binary and starts executing it")

	// KEY TAKEAWAY:
	// - Go source code is transformed through several compiler stages.
	// - The thing that runs is the final binary, not the original text file.
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HC.3 memory-basics")
	fmt.Println("Run    : go run ./00-how-computers-work/3-memory-basics")
	fmt.Println("Current: HC.2 (code-to-execution)")
	fmt.Println("---------------------------------------------------")
}
