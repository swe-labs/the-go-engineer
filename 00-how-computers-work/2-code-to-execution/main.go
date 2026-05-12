// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 00: How Computers Work
// Title: How code becomes execution
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Understand the journey from Go source code to a running program.
//   - The stages: tokens, AST, type checking, IR, optimization, and code generation.
//
// WHY THIS MATTERS:
//   - The compiler is a translation pipeline. It turns human-readable text into
//     machine-optimized instructions. Knowing this explains why Go catches errors
//     early and why binaries are the unit of deployment.
//
// RUN:
//   go run ./00-how-computers-work/2-code-to-execution
//
// KEY TAKEAWAY:
//   - The CPU never sees your .go files; it only sees the machine code output
//     of the build pipeline.
// ============================================================================

package main

import "fmt"

// main (Function): entry point for the program. It walks through the Go
// compiler pipeline stages — tokens, AST, IR, binary — to demonstrate how
// source code becomes a running executable.
func main() {
	fmt.Println("Source code : x := 42 + y")
	fmt.Println("Tokens      : IDENT(x) DEFINE(:=) NUMBER(42) PLUS(+) IDENT(y)")
	fmt.Println("AST         : Assign(x, Add(42, y))")
	fmt.Println("IR          : t0 = load y; t1 = add 42, t0; store x, t1")
	fmt.Println("Binary      : machine instructions + linked packages")
	fmt.Println("Runtime     : the OS loads the binary and starts executing it")

	// - Go source code is transformed through several compiler stages.
	// - The thing that runs is the final binary, not the original text file.
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HC.3 -> 00-how-computers-work/3-memory-basics")
	fmt.Println("Run    : go run ./00-how-computers-work/3-memory-basics")
	fmt.Println("Current: HC.2 (code-to-execution)")
	fmt.Println("---------------------------------------------------")
}
