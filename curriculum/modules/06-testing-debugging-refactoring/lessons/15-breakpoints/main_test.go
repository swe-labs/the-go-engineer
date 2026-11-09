package main

import "testing"

func TestComputeScore(t *testing.T) {
	tests := []struct {
		id   int
		want int
	}{
		{0, 0},  // 0*10 + 0*0 = 0
		{1, 11}, // 10 + 1 = 11
		{2, 24}, // 20 + 4 = 24
		{3, 39}, // 30 + 9 = 39
		{5, 75}, // 50 + 25 = 75
	}
	for _, tc := range tests {
		got := computeScore(tc.id)
		if got != tc.want {
			t.Errorf("computeScore(%d) = %d, want %d", tc.id, got, tc.want)
		}
	}
}

func TestMultiply(t *testing.T) {
	got := multiply(6, 7)
	if got != 42 {
		t.Errorf("multiply(6,7) = %d, want 42", got)
	}
}

func TestProcessUser(t *testing.T) {
	// processUser prints to stdout; verify it doesn't panic.
	processUser(99, "test")
}
