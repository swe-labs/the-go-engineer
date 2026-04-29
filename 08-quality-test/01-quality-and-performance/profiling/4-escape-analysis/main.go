// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 08: Quality & Testing
// Title: Escape analysis
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how the compiler decides whether values stay on the stack or escape to the heap.
//
// WHY THIS MATTERS:
//   - Escape analysis is the compiler's answer to one question: does this value need to outlive the current stack frame?
//
// RUN:
//   go run ./08-quality-test/01-quality-and-performance/profiling/4-escape-analysis
//
// KEY TAKEAWAY:
//   - Learn how the compiler decides whether values stay on the stack or escape to the heap.
// ============================================================================

package main

import "fmt"

//

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
