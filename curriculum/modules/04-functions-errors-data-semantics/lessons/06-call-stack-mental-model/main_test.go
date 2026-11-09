package main

import "testing"

func TestFib(t *testing.T) {
	tests := []struct {
		n        int
		expected int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 5},
		{10, 55},
	}
	for _, tc := range tests {
		got := fib(tc.n)
		if got != tc.expected {
			t.Errorf("fib(%d) = %d; want %d", tc.n, got, tc.expected)
		}
	}
}
