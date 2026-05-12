// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 08: Quality & Testing
// Title: Table-Driven Tests
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The idiomatic Go pattern for testing multiple cases without duplication
//   - Using anonymous structs as test case tables
//   - Sub-tests with t.Run for granular failure reporting
//
// WHY THIS MATTERS:
//   Table-driven tests are the standard Go testing idiom. They make adding
//   new test cases a one-line change and keep your test code DRY.
//
// RUN:
//   go test ./08-quality-test/01-quality-and-performance/02-testing/02-table-driven-tests
// ============================================================================

package main

import (
	"fmt"
	"strings"
)

// ValidateEmail checks whether an email address has a basic valid shape.
// ValidateEmail (Function): checks whether an email address has a basic valid shape.
// Boundary: accepts only addresses containing exactly one @ and at least one dot after.
func ValidateEmail(email string) bool {
	atIndex := strings.Index(email, "@")
	if atIndex < 1 {
		return false
	}
	dotIndex := strings.LastIndex(email, ".")
	if dotIndex <= atIndex+1 {
		return false
	}
	return true
}

// CalculateScore returns a numeric score based on input length and content.
// CalculateScore (Function): returns a numeric score based on input length and content.
func CalculateScore(input string) int {
	if len(input) == 0 {
		return 0
	}
	score := len(input) * 10
	if strings.Contains(strings.ToLower(input), "bonus") {
		score += 50
	}
	if score > 100 {
		score = 100
	}
	return score
}

// KEY TAKEAWAY:
// - Table-driven tests are the idiomatic Go pattern for multiple test cases.
// - Anonymous structs define the table; t.Run creates sub-tests per case.
// - Adding a new case is a one-line change: no new function needed.
//
// NEXT UP: TE.3 -> 08-quality-test/01-quality-and-performance/02-testing/03-http-handler-testing
func main() {
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: TE.3 -> 08-quality-test/01-quality-and-performance/02-testing/03-http-handler-testing")
	fmt.Println("Run    : go test ./08-quality-test/01-quality-and-performance/02-testing/03-http-handler-testing")
	fmt.Println("Current: TE.2 (table-driven-tests)")
	fmt.Println("---------------------------------------------------")
}
