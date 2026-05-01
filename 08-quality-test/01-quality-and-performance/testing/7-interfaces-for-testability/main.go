// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 08: Quality & Testing
// Title: Interfaces for testability
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how narrow interfaces create seams that let code be tested without real dependencies.
//
// WHY THIS MATTERS:
//   - A testing seam is a boundary where production code depends on behavior, not a concrete implementation.
//
// RUN:
//
//
// KEY TAKEAWAY:
//   - Learn how narrow interfaces create seams that let code be tested without real dependencies.
// ============================================================================

package main

import (
	"fmt"
	"strings"
)

//
// NEXT UP: TE.8 -> 08-quality-test/01-quality-and-performance/testing/8-mocking-with-interfaces

func te_7Summary(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func main() {
	fmt.Println("summary:", te_7Summary("  Interface Seams  "))
}
