// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 08: Quality & Testing
// Title: Testable Design with io.Writer
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Why hardcoding stdout/stderr makes code untestable
//   - Refactoring to use io.Writer for dependency injection
//   - Testing output by passing a strings.Builder instead of os.Stdout
//
// WHY THIS MATTERS:
//   Functions that print directly to os.Stdout (via fmt.Println) cannot be
//   tested without hijacking the OS output stream — which is brittle and
//   not concurrent-safe. Accepting an io.Writer makes the same function
//   testable with a simple strings.Builder.
//
// RUN:
//   go test ./08-quality-test/01-quality-and-performance/02-testing/11-testable-design
// ============================================================================

package main

import (
	"fmt"
	"io"
	"os"
)

// GreetingHardcoded prints directly to stdout — untestable without OS hacks.
// GreetingHardcoded (Function): prints directly to stdout — untestable without OS hacks.
func GreetingHardcoded(prefix, name string) {
	fmt.Printf("Hello, %s %s!\n", prefix, name)
}

// GreetingTestable accepts any io.Writer — testable with any writer.
// GreetingTestable (Function): accepts any io.Writer — testable with any writer.
// Boundary: decouples output destination from business logic.
func GreetingTestable(out io.Writer, prefix, name string) {
	fmt.Fprintf(out, "Hello, %s %s!\n", prefix, name)
}

// KEY TAKEAWAY:
// - Accept io.Writer instead of writing to os.Stdout directly.
// - Pass os.Stdout in production, pass a buffer in tests.
// - This is the simplest form of dependency injection in Go.

func main() {
	fmt.Println()
	GreetingHardcoded("Mr.", "Hardcoded")
	GreetingTestable(os.Stdout, "Ms.", "Testable")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: PR.1 -> 08-quality-test/01-quality-and-performance/01-profiling/01-cpu-profile")
	fmt.Println("Run    : go run ./08-quality-test/01-quality-and-performance/01-profiling/01-cpu-profile")
	fmt.Println("Current: TE.11 (testable-design)")
	fmt.Println("---------------------------------------------------")
}
