package main

import "testing"

func TestEvaluate(t *testing.T) {
	tests := []struct {
		a, b     int
		op       rune
		expected int
	}{
		{10, 5, '+', 15},
		{10, 5, '-', 5},
		{10, 5, '*', 50},
		{10, 5, '/', 2},
		{10, 5, '%', 0},
		{7, 3, '+', 10},
		{7, 3, '-', 4},
		{7, 3, '*', 21},
		{7, 3, '/', 2},
		{7, 3, '%', 1},
	}

	for _, tc := range tests {
		got := evaluate(tc.a, tc.b, tc.op)
		if got != tc.expected {
			t.Errorf("evaluate(%d, %d, %q) = %d; want %d", tc.a, tc.b, tc.op, got, tc.expected)
		}
	}
}
