package main

import "testing"

func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{"no args", nil, 0},
		{"single", []int{5}, 5},
		{"multiple", []int{3, 7, 2}, 7},
		{"negatives", []int{-5, -1, -10}, -1},
		{"all equal", []int{4, 4, 4}, 4},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := max(tc.nums...)
			if got != tc.expected {
				t.Errorf("max(%v) = %d; want %d", tc.nums, got, tc.expected)
			}
		})
	}
}

func TestConcat(t *testing.T) {
	tests := []struct {
		sep   string
		parts []string
		ex    string
	}{
		{", ", []string{"a", "b", "c"}, "a, b, c"},
		{"-", []string{"x"}, "x"},
		{" ", nil, ""},
		{"...", []string{"hello", "world"}, "hello...world"},
	}
	for _, tc := range tests {
		got := concat(tc.sep, tc.parts...)
		if got != tc.ex {
			t.Errorf("concat(%q, %v) = %q; want %q", tc.sep, tc.parts, got, tc.ex)
		}
	}
}
