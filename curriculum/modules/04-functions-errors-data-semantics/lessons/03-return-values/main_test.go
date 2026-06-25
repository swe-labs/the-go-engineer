package main

import "testing"

func TestFactorial(t *testing.T) {
	tests := []struct {
		n        int
		expected int
	}{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 6},
		{5, 120},
		{7, 5040},
	}
	for _, tc := range tests {
		got := factorial(tc.n)
		if got != tc.expected {
			t.Errorf("factorial(%d) = %d; want %d", tc.n, got, tc.expected)
		}
	}
}

func TestSafeDivide(t *testing.T) {
	if got := safeDivide(10, 3); got != 3 {
		t.Errorf("safeDivide(10, 3) = %d; want 3", got)
	}
	if got := safeDivide(10, 0); got != 0 {
		t.Errorf("safeDivide(10, 0) = %d; want 0", got)
	}
}
