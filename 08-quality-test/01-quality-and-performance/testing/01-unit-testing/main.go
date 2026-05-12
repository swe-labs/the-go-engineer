// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 08: Quality & Testing
// Title: Unit Testing
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Writing standard unit tests with Go's built-in testing package
//   - Test file naming conventions and function signatures
//   - Using testify/assert for cleaner test output
//
// WHY THIS MATTERS:
//   Go's testing package is built into the toolchain. There are no magic
//   frameworks — tests are just Go code. This makes them fast, readable,
//   and easy to debug.
//
// RUN:
//   go test ./08-quality-test/01-quality-and-performance/testing/01-unit-testing
// ============================================================================

package main

import (
	"errors"
	"fmt"
	"strings"
)

// CheckUsername validates that a username meets platform requirements.
// CheckUsername (Function): validates that a username meets platform requirements.
// Role: input validation boundary — rejects short or admin-impersonating names.
func CheckUsername(username string) bool {
	if len(username) < 6 {
		return false
	}
	if strings.Contains(strings.ToLower(username), "admin") {
		return false
	}
	return true
}

// Login attempts to authenticate a user and returns success and error.
// Login (Function): attempts to authenticate a user.
// Failure mode: returns false with a descriptive error for invalid input.
func Login(username string) (bool, error) {
	if !CheckUsername(username) {
		return false, errors.New("invalid username: must be 6+ chars and not contain admin")
	}
	return true, nil
}

// KEY TAKEAWAY:
// - Go tests are just Go code: no magic frameworks, no decorators.
// - test files end with `_test.go` and functions start with `Test`.
// - testify/assert makes test output cleaner than raw `if` checks.
//
// NEXT UP: TE.2 -> 08-quality-test/01-quality-and-performance/testing/02-table-driven-tests
func main() {
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: TE.2 -> 08-quality-test/01-quality-and-performance/testing/02-table-driven-tests")
	fmt.Println("Run    : go test ./08-quality-test/01-quality-and-performance/testing/02-table-driven-tests")
	fmt.Println("Current: TE.1 (unit-testing)")
	fmt.Println("---------------------------------------------------")
}
