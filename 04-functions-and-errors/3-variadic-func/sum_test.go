// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "testing"

// ============================================================================
// Tests for: Variadic Functions
// ============================================================================
//
// These tests verify the sum() variadic function handles
// zero args, single args, multiple args, and spread slices.
//
// RUN: go test -v ./04-functions-and-errors/3-variadic-func
// ============================================================================

func TestSum(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		expect int
	}{
		{name: "no arguments", input: []int{}, expect: 0},
		{name: "single value", input: []int{5}, expect: 5},
		{name: "multiple values", input: []int{1, 2, 3, 4}, expect: 10},
		{name: "negative values", input: []int{-5, 10, -3}, expect: 2},
		{name: "large slice", input: []int{10, 20, 30}, expect: 60},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use the spread operator to pass a slice to a variadic function.
			got := sum(tt.input...)
			if got != tt.expect {
				t.Errorf("sum(%v) = %d, want %d", tt.input, got, tt.expect)
			}
		})
	}
}
