// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 08: Quality & Testing
// Title: Integration tests
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn where integration tests sit between unit tests and end-to-end tests and what they should actually prove.
//
// WHY THIS MATTERS:
//   - Integration tests verify that real components work together across a boundary.
//
// RUN:
//
//
// KEY TAKEAWAY:
//   - Learn where integration tests sit between unit tests and end-to-end tests and what they should actually prove.
// ============================================================================

package main

import (
	"fmt"
	"strings"
)

//
// NEXT UP: TE.10

func te_9Summary(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func main() {
	fmt.Println("summary:", te_9Summary("  Integration Boundaries  "))
}
