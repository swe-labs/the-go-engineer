// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package user

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// Test Files
// ============================================================================
// 1. Must be named `*_test.go`
// 2. Must be in the same package (usually) to test unexported functions,
//    or `package user_test` for strict black-box testing.
// 3. Test functions must start with `Test` and take `t *testing.T`.
// ============================================================================

// --- 1. Basic Unit Test ---
// The simplest form of testing. Good for single, obvious cases.
func TestCheckUsername_Basic(t *testing.T) {
	// Arrange
	username := "rasel9t6"

	// Act
	result := CheckUsername(username)

	// Assert
	// We use "github.com/stretchr/testify/assert" because Go's standard library
	// doesn't have built-in assertions. It forces you to write `if result != true`.
	// Testify makes test logs much cleaner.
	assert.True(t, result, "rasel9t6 should be a valid username")
}

// --- 2. Table-Driven Tests ---
// THE IDIOMATIC GO WAY to write tests.
// Instead of writing 10 different Test functions for 10 conditions, we define
// a slice of anonymous structs (a "table") and loop over it.
func TestCheckUsername_TableDriven(t *testing.T) {
	// 1. Define the table
	testCases := []struct {
		desc     string // Description of the exact case
		input    string // The input to the function
		expected bool   // The expected result
	}{
		{"Valid username", "rasel9t6", true},
		{"Too short", "bob", false},
		{"Exact minimum length", "abcdef", true},
		{"Contains admin (lowercase)", "admin_user", false},
		{"Contains admin (uppercase)", "AdMiN_user", false},
	}

	// 2. Loop over the table
	for _, tc := range testCases {
		// 3. Use t.Run to create a "Sub-test" for each case.
		// This means if one test fails, the OTHERS KEEP RUNNING.
		// It also prints beautifully in the terminal: `go test -v ./...`
		t.Run(tc.desc, func(t *testing.T) {
			actual := CheckUsername(tc.input)

			// Assert
			// The custom message helps debugging when table tests fail
			assert.Equal(t, tc.expected, actual, "Input: %q", tc.input)
		})
	}
}

// --- 3. Testing Multiple Returns (Error Handling) ---
func TestLogin(t *testing.T) {
	t.Run("Valid Login", func(t *testing.T) {
		success, err := Login("rasel9t6")

		assert.NoError(t, err)  // Error should be nil
		assert.True(t, success) // Success should be true
	})

	t.Run("Invalid Login", func(t *testing.T) {
		success, err := Login("admin")

		assert.Error(t, err) // Error should NOT be nil
		assert.False(t, success)

		// We can even assert the exact error string
		assert.Equal(t, errors.New("invalid username: must be 6+ chars and not contain admin"), err)
	})
}
