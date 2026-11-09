package main

import (
	"testing"
)

func TestDescribe(t *testing.T) {
	tests := []struct {
		input any
		want  string
	}{
		{nil, "nil"},
		{42, "integer 42"},
		{"hello", `string "hello" (len 5)`},
		{true, "true"},
		{false, "false"},
		{3.14, "unknown type float64"},
	}
	for _, tc := range tests {
		got := describe(tc.input)
		if got != tc.want {
			t.Errorf("describe(%v) = %q; want %q", tc.input, got, tc.want)
		}
	}
}

func TestMathSwitch(t *testing.T) {
	tests := []struct {
		input any
		want  float64
		err   bool
	}{
		{float64(2.5), 2.5, false},
		{int(10), 10.0, false},
		{int64(20), 20.0, false},
		{float32(3.5), 7.0, false},
		{"hello", 0, true},
		{true, 0, true},
	}
	for _, tc := range tests {
		got, err := mathSwitch(tc.input)
		if tc.err && err == nil {
			t.Errorf("mathSwitch(%v) expected error, got %v", tc.input, got)
		}
		if !tc.err && err != nil {
			t.Errorf("mathSwitch(%v) unexpected error: %v", tc.input, err)
		}
		if !tc.err && got != tc.want {
			t.Errorf("mathSwitch(%v) = %v; want %v", tc.input, got, tc.want)
		}
	}
}
