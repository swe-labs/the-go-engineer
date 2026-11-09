package main

import (
	"reflect"
	"testing"
)

func TestRotate(t *testing.T) {
	tests := []struct {
		nums     []int
		n        int
		expected []int
	}{
		{[]int{1, 2, 3, 4, 5}, 2, []int{3, 4, 5, 1, 2}},
		{[]int{1, 2, 3, 4, 5}, 0, []int{1, 2, 3, 4, 5}},
		{[]int{1, 2, 3, 4, 5}, 5, []int{1, 2, 3, 4, 5}},
		{[]int{1}, 3, []int{1}},
		{[]int{}, 2, []int{}},
		{[]int{10, 20, 30}, 1, []int{20, 30, 10}},
	}
	for _, tc := range tests {
		original := make([]int, len(tc.nums))
		copy(original, tc.nums)

		got := rotate(tc.nums, tc.n)
		if !reflect.DeepEqual(got, tc.expected) {
			t.Errorf("rotate(%v, %d) = %v; want %v", tc.nums, tc.n, got, tc.expected)
		}
		if !reflect.DeepEqual(tc.nums, original) {
			t.Errorf("rotate modified original: %v → %v", original, tc.nums)
		}
	}
}
