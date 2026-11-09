package main

import (
	"testing"
	"time"
)

func TestDurationArithmetic(t *testing.T) {
	tests := []struct {
		name string
		a, b time.Duration
		op   string
		want time.Duration
	}{
		{"add seconds", 5 * time.Second, 3 * time.Second, "+", 8 * time.Second},
		{"subtract", 10 * time.Second, 4 * time.Second, "-", 6 * time.Second},
		{"add mixed", 1*time.Minute + 30*time.Second, 15 * time.Second, "+", 1*time.Minute + 45*time.Second},
		{"negative result", 3 * time.Second, 10 * time.Second, "-", -7 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got time.Duration
			switch tt.op {
			case "+":
				got = tt.a + tt.b
			case "-":
				got = tt.a - tt.b
			}
			if got != tt.want {
				t.Errorf("%v %s %v = %v, want %v", tt.a, tt.op, tt.b, got, tt.want)
			}
		})
	}
}

func TestDurationComparison(t *testing.T) {
	tests := []struct {
		name string
		a, b time.Duration
		less bool
	}{
		{"seconds less", 2 * time.Second, 5 * time.Second, true},
		{"equal", 3 * time.Second, 3 * time.Second, false},
		{"minutes vs seconds", 1 * time.Minute, 30 * time.Second, false},
		{"negative vs positive", -1 * time.Second, 1 * time.Second, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a < tt.b; got != tt.less {
				t.Errorf("%v < %v = %v, want %v", tt.a, tt.b, got, tt.less)
			}
		})
	}
}

func TestDurationConversion(t *testing.T) {
	d := 90 * time.Second
	if d.Minutes() != 1.5 {
		t.Errorf("Minutes() = %v, want 1.5", d.Minutes())
	}
	if d.Truncate(time.Minute) != 1*time.Minute {
		t.Errorf("Truncate(1m) = %v, want 1m", d.Truncate(time.Minute))
	}
}
