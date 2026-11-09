package main

import "testing"

func TestMerge(t *testing.T) {
	tests := []struct {
		a, b []int
		want []int
	}{
		{[]int{1, 2}, []int{3, 4}, []int{1, 2, 3, 4}},
		{[]int{}, []int{1, 2}, []int{1, 2}},
		{[]int{1, 2}, []int{}, []int{1, 2}},
		{[]int{}, []int{}, []int{}},
		{[]int{1}, []int{2}, []int{1, 2}},
	}

	for _, tc := range tests {
		got := merge(tc.a, tc.b)
		if len(got) != len(tc.want) {
			t.Errorf("merge(%v, %v) len=%d; want len=%d", tc.a, tc.b, len(got), len(tc.want))
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("merge(%v, %v)[%d] = %d; want %d", tc.a, tc.b, i, got[i], tc.want[i])
			}
		}
	}
}

func TestDedup(t *testing.T) {
	tests := []struct {
		input []int
		want  []int
	}{
		{[]int{1, 1, 2, 3, 3, 3, 4, 5, 5}, []int{1, 2, 3, 4, 5}},
		{[]int{}, []int{}},
		{[]int{7}, []int{7}},
		{[]int{1, 1, 1, 1}, []int{1}},
		{[]int{1, 2, 3, 4}, []int{1, 2, 3, 4}},
	}

	for _, tc := range tests {
		got := dedup(tc.input)
		if len(got) != len(tc.want) {
			t.Errorf("dedup(%v) len=%d; want len=%d", tc.input, len(got), len(tc.want))
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("dedup(%v)[%d] = %d; want %d", tc.input, i, got[i], tc.want[i])
			}
		}
	}
}
