package main

import "testing"

func TestSumAndAverage(t *testing.T) {
	tests := []struct {
		nums    [8]float64
		wantSum float64
		wantAvg float64
	}{
		{[8]float64{1, 2, 3, 4, 5, 6, 7, 8}, 36, 4.5},
		{[8]float64{0, 0, 0, 0, 0, 0, 0, 0}, 0, 0},
		{[8]float64{1.5, 2.5, 3.0, 4.0, 5.5, 6.0, 7.5, 8.0}, 38, 4.75},
		{[8]float64{-1, -2, -3, -4, -5, -6, -7, -8}, -36, -4.5},
	}

	for _, tc := range tests {
		gotSum, gotAvg := sumAndAverage(tc.nums)
		if gotSum != tc.wantSum || gotAvg != tc.wantAvg {
			t.Errorf("sumAndAverage(%v) = (%f, %f); want (%f, %f)", tc.nums, gotSum, gotAvg, tc.wantSum, tc.wantAvg)
		}
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		input [6]int
		want  [6]int
	}{
		{[6]int{1, 2, 3, 4, 5, 6}, [6]int{6, 5, 4, 3, 2, 1}},
		{[6]int{6, 5, 4, 3, 2, 1}, [6]int{1, 2, 3, 4, 5, 6}},
		{[6]int{0, 0, 0, 0, 0, 0}, [6]int{0, 0, 0, 0, 0, 0}},
		{[6]int{1, 1, 1, 1, 1, 1}, [6]int{1, 1, 1, 1, 1, 1}},
		{[6]int{-1, -2, -3, -4, -5, -6}, [6]int{-6, -5, -4, -3, -2, -1}},
	}

	for _, tc := range tests {
		var copy [6]int
		copy = tc.input
		reverse(&copy)
		if copy != tc.want {
			t.Errorf("reverse(&%v) = %v; want %v", tc.input, copy, tc.want)
		}
	}
}

func TestModifyAttemptDoesNotChangeOriginal(t *testing.T) {
	original := [6]int{10, 20, 30, 40, 50, 60}
	saved := original
	modifyAttempt(original)
	if original != saved {
		t.Errorf("original changed after modifyAttempt: got %v, want %v", original, saved)
	}
}
