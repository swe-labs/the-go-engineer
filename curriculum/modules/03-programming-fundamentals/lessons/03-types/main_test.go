package main

import "testing"

func TestConvertUnits(t *testing.T) {
	tests := []struct {
		value    float64
		from, to string
		want     float64
		wantErr  bool
		epsilon  float64
	}{
		{100, "C", "F", 212, false, 0.01},
		{0, "C", "F", 32, false, 0.01},
		{-40, "C", "F", -40, false, 0.01},
		{32, "F", "C", 0, false, 0.01},
		{212, "F", "C", 100, false, 0.01},
		{0, "C", "K", 273.15, false, 0.01},
		{273.15, "K", "C", 0, false, 0.01},
		{100, "C", "C", 100, false, 0.01},
		{100, "C", "X", 0, true, 0},
		{100, "X", "C", 0, true, 0},
	}

	for _, tc := range tests {
		got, err := convertUnits(tc.value, tc.from, tc.to)
		if tc.wantErr {
			if err == nil {
				t.Errorf("convertUnits(%.2f, %q, %q) expected error but got result %.2f", tc.value, tc.from, tc.to, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("convertUnits(%.2f, %q, %q) unexpected error: %v", tc.value, tc.from, tc.to, err)
			continue
		}
		diff := got - tc.want
		if diff < 0 {
			diff = -diff
		}
		if diff > tc.epsilon {
			t.Errorf("convertUnits(%.2f, %q, %q) = %.4f; want %.4f (diff %.4f)", tc.value, tc.from, tc.to, got, tc.want, diff)
		}
	}
}
