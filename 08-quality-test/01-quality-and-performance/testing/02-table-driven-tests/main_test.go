// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidateEmail (Test): demonstrates table-driven testing for a validation function.
func TestValidateEmail(t *testing.T) {
	cases := []struct {
		desc  string
		input string
		want  bool
	}{
		{desc: "valid email", input: "user@example.com", want: true},
		{desc: "no @ symbol", input: "userexample.com", want: false},
		{desc: "empty local part", input: "@example.com", want: false},
		{desc: "no dot in domain", input: "user@example", want: false},
		{desc: "dot before @", input: "user@.example", want: false},
		{desc: "empty string", input: "", want: false},
		{desc: "multiple dots", input: "user@sub.example.com", want: true},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := ValidateEmail(tc.input)
			assert.Equal(t, tc.want, got, "ValidateEmail(%q)", tc.input)
		})
	}
}

// TestCalculateScore (Test): demonstrates table-driven testing with integer results.
func TestCalculateScore(t *testing.T) {
	cases := []struct {
		desc  string
		input string
		want  int
	}{
		{desc: "empty", input: "", want: 0},
		{desc: "short", input: "abc", want: 30},
		{desc: "longer", input: "hello world", want: 100},
		{desc: "contains bonus", input: "bonus", want: 100},
		{desc: "capped at 100", input: "this is a very long input string that exceeds ten chars", want: 100},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := CalculateScore(tc.input)
			assert.Equal(t, tc.want, got, "CalculateScore(%q)", tc.input)
		})
	}
}
