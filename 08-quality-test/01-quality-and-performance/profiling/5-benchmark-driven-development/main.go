// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 08: Quality & Testing
// Title: Benchmark-driven development
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how benchmarks protect performance discussions from guesswork by measuring before and after changes.
//
// WHY THIS MATTERS:
//   - Benchmarks turn performance from opinion into evidence.
//
// RUN:
//
//
// KEY TAKEAWAY:
//   - Learn how benchmarks protect performance discussions from guesswork by measuring before and after changes.
// ============================================================================

package main

import (
	"fmt"
	"strings"
)

//
// NEXT UP: PR.6

func pr_5Summary(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func main() {
	fmt.Println("summary:", pr_5Summary("  Benchmark Driven Development  "))
}
