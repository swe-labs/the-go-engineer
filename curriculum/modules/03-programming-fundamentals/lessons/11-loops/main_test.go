package main

import "testing"

func TestSumUntil(t *testing.T) {
	tests := []struct {
		nums               []int
		target             int
		wantSum, wantCount int
	}{
		{[]int{3, -1, 5, 2, -3, 4}, 10, 10, 3}, // 3 + 5 + 2 = 10, count=3
		{[]int{1, 2, 3, 4, 5}, 100, 15, 5},     // all summed, never reaches target
		{[]int{10, -1, 20}, 5, 10, 1},          // first element reaches target
		{[]int{-1, -2, -3}, 1, 0, 0},           // all skipped, sum=0, count=0
		{[]int{}, 10, 0, 0},                    // empty slice
		{[]int{5, -5, 5, -5, 5}, 10, 10, 2},    // stops at sum>=10: 5+5=10, count=2
		{[]int{0, 0, 0}, 0, 0, 1},              // zero reaches target immediately
	}

	for _, tc := range tests {
		gotSum, gotCount := sumUntil(tc.nums, tc.target)
		if gotSum != tc.wantSum || gotCount != tc.wantCount {
			t.Errorf("sumUntil(%v, %d) = (%d, %d); want (%d, %d)",
				tc.nums, tc.target, gotSum, gotCount, tc.wantSum, tc.wantCount)
		}
	}
}

func TestForRange(t *testing.T) {
	// Verify range over slice produces correct elements
	nums := []int{10, 20, 30}
	var got []int
	for _, v := range nums {
		got = append(got, v)
	}
	if len(got) != 3 || got[0] != 10 || got[1] != 20 || got[2] != 30 {
		t.Errorf("range over slice failed: %v", got)
	}
}

func TestLabeledBreak(t *testing.T) {
	var result []int
outer:
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if i == 2 && j == 2 {
				break outer
			}
			result = append(result, i*10+j)
		}
	}
	// Should stop at i=2, j=2
	if len(result) != 12 { // 0,0-0,4 (5) + 1,0-1,4 (5) + 2,0-2,1 (2) = 12
		t.Errorf("expected 12 elements, got %d: %v", len(result), result)
	}
}

func TestForLoopContinue(t *testing.T) {
	var result []int
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue
		}
		result = append(result, i)
	}
	if len(result) != 5 {
		t.Errorf("expected 5 odd numbers, got %d: %v", len(result), result)
	}
	for _, v := range result {
		if v%2 == 0 {
			t.Errorf("unexpected even number %d in result", v)
		}
	}
}
