// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCheckUsername_Basic (Test): demonstrates a basic unit test — one input, one assertion.
func TestCheckUsername_Basic(t *testing.T) {
	result := CheckUsername("rasel9t6")
	assert.True(t, result, "rasel9t6 should be a valid username")
}

// TestCheckUsername_TableDriven (Test): demonstrates the idiomatic Go table-driven pattern.
// Instead of writing separate test functions for each case, we define a table of cases
// and loop over them. This example is a preview — see TE.2 for the full treatment.
func TestCheckUsername_TableDriven(t *testing.T) {
	cases := []struct {
		desc  string
		input string
		want  bool
	}{
		{"valid username", "rasel9t6", true},
		{"too short", "bob", false},
		{"exact minimum", "abcdef", true},
		{"contains admin", "admin_user", false},
		{"contains admin uppercase", "AdMiN_user", false},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := CheckUsername(tc.input)
			assert.Equal(t, tc.want, got, "CheckUsername(%q)", tc.input)
		})
	}
}

// TestLogin (Test): tests a function that returns multiple values, including an error.
func TestLogin(t *testing.T) {
	t.Run("valid login", func(t *testing.T) {
		success, err := Login("rasel9t6")
		assert.NoError(t, err)
		assert.True(t, success)
	})

	t.Run("invalid login", func(t *testing.T) {
		success, err := Login("admin")
		assert.Error(t, err)
		assert.False(t, success)
		assert.Equal(t, errors.New("invalid username: must be 6+ chars and not contain admin"), err)
	})
}
