package main

import "testing"

func TestSwap(t *testing.T) {
	tests := []struct {
		a, b         int
		wantA, wantB int
	}{
		{1, 2, 2, 1},
		{0, 0, 0, 0},
		{-5, 5, 5, -5},
		{100, -100, -100, 100},
	}

	for _, tc := range tests {
		a, b := tc.a, tc.b
		swap(&a, &b)
		if a != tc.wantA || b != tc.wantB {
			t.Errorf("swap(%d, %d) = (%d, %d); want (%d, %d)", tc.a, tc.b, a, b, tc.wantA, tc.wantB)
		}
	}
}

func TestApply(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	apply(nums, func(p *int) { *p *= 2 })
	want := []int{2, 4, 6, 8, 10}
	for i := range nums {
		if nums[i] != want[i] {
			t.Errorf("apply[%d] = %d; want %d", i, nums[i], want[i])
		}
	}
}

func TestApplyEmpty(t *testing.T) {
	nums := []int{}
	apply(nums, func(p *int) { *p = 0 })
}

func TestApplyIdentity(t *testing.T) {
	nums := []int{5, 10, 15}
	apply(nums, func(p *int) {})
	if nums[0] != 5 || nums[1] != 10 || nums[2] != 15 {
		t.Errorf("identity apply modified slice: %v", nums)
	}
}
