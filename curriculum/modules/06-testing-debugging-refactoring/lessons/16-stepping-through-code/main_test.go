package main

import "testing"

func TestFactorial(t *testing.T) {
	tests := []struct {
		n    int
		want int
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
		if got != tc.want {
			t.Errorf("factorial(%d) = %d, want %d", tc.n, got, tc.want)
		}
	}
}

func TestFactorialNegative(t *testing.T) {
	// factorial(-1): n <= 1 is true, returns 1
	got := factorial(-1)
	if got != 1 {
		t.Errorf("factorial(-1) = %d, want 1", got)
	}
}
