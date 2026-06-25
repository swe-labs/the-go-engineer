package main

import (
	"math"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		r, h       float64
		wantArea   float64
		wantVolume float64
	}{
		{3, 5, 2 * math.Pi * 3 * (3 + 5), math.Pi * 3 * 3 * 5},
		{1.5, 2, 2 * math.Pi * 1.5 * (1.5 + 2), math.Pi * 1.5 * 1.5 * 2},
		{10, 1, 2 * math.Pi * 10 * (10 + 1), math.Pi * 10 * 10 * 1},
		{0, 5, 2 * math.Pi * 0 * (0 + 5), math.Pi * 0 * 0 * 5},
	}

	for _, tc := range tests {
		gotArea, gotVolume := calculate(tc.r, tc.h)
		if math.Abs(gotArea-tc.wantArea) > 1e-9 {
			t.Errorf("calculate(%.1f, %.1f) area = %.10f; want %.10f", tc.r, tc.h, gotArea, tc.wantArea)
		}
		if math.Abs(gotVolume-tc.wantVolume) > 1e-9 {
			t.Errorf("calculate(%.1f, %.1f) volume = %.10f; want %.10f", tc.r, tc.h, gotVolume, tc.wantVolume)
		}
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		a, b     int
		wantQuot int
		wantRem  int
	}{
		{10, 3, 3, 1},
		{100, 2, 50, 0},
		{7, 5, 1, 2},
	}

	for _, tc := range tests {
		gotQuot, gotRem := divide(tc.a, tc.b)
		if gotQuot != tc.wantQuot || gotRem != tc.wantRem {
			t.Errorf("divide(%d, %d) = (%d, %d); want (%d, %d)", tc.a, tc.b, gotQuot, gotRem, tc.wantQuot, tc.wantRem)
		}
	}
}
