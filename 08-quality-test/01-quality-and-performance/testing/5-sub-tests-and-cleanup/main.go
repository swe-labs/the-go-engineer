// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 08: Quality & Testing
// Title: Sub-tests and t.Cleanup
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how sub-tests scope assertions and how t.Cleanup keeps test teardown tied to the test that created the resource.
//
// WHY THIS MATTERS:
//   - Sub-tests let one test file express many related cases without losing names or isolation.
//
// RUN:
//
//
// KEY TAKEAWAY:
//   - Learn how sub-tests scope assertions and how t.Cleanup keeps test teardown tied to the test that created the resource.
// ============================================================================

package main

import (
	"fmt"
	"strings"
)

//
// NEXT UP: TE.6 -> 08-quality-test/01-quality-and-performance/testing/6-fuzz-testing

func te_5Summary(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func main() {
	fmt.Println("summary:", te_5Summary("  Subtests With Cleanup  "))
}
