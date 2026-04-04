// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package user

import (
	"errors"
	"strings"
)

// ============================================================================
// Section 14: Testing — Foundational Testing concepts
// Level: Beginner → Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Writing standard unit tests
//   - Table-driven tests (the Go way)
//   - Sub-tests (t.Run)
//   - The testing package vs assert libraries (testify)
//
// ENGINEERING DEPTH:
//   Go has a built-in testing framework: `go test`.
//   It expects test files to end with `_test.go` and functions to start
//   with `TestXxx(t *testing.T)`.
//   Unlike Java/C#, Go tests are just normal Go code. There are no
//   decorators or magic syntax.
// ============================================================================

// CheckUsername validates if a username meets our platform requirements.
// A valid username must be at least 6 characters and cannot contain "admin".
func CheckUsername(username string) bool {
	if len(username) < 6 {
		return false
	}

	// Case-insensitive check to prevent "AdMin" loopholes
	if strings.Contains(strings.ToLower(username), "admin") {
		return false
	}
	return true
}

// Login attempts to authenticate a user.
// Returns a boolean success flag, and an error if validation fails.
func Login(username string) (bool, error) {
	if !CheckUsername(username) {
		return false, errors.New("invalid username: must be 6+ chars and not contain admin")
	}

	// In a real application, we would check a database and compare password hashes
	// See Section 12 (Databases) and Section 13 (Web) for the full implementation.
	return true, nil
}
