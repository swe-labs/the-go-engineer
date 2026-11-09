package main

import (
	"math"
	"testing"
)

func TestParseCelsius(t *testing.T) {
	tests := []struct {
		input   string
		want    float64
		wantErr bool
	}{
		{"23.5", 23.5, false},
		{"-10", -10, false},
		{"37.2C", 37.2, false},
		{"0C", 0, false},
		{"-40C", -40, false},
		{"100C", 100, false},
		{"", 0, true},
		{"abc", 0, true},
		{"12.5F", 0, true},
		{"300K", 0, true},
		{"12.5.6", 0, true},
		{"  25.0  ", 25.0, false},
		{"-273.15C", -273.15, false},
	}

	for _, tc := range tests {
		got, err := parseCelsius(tc.input)
		if tc.wantErr {
			if err == nil {
				t.Errorf("parseCelsius(%q) expected error, got %f", tc.input, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("parseCelsius(%q) unexpected error: %v", tc.input, err)
			continue
		}
		if math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("parseCelsius(%q) = %f; want %f", tc.input, got, tc.want)
		}
	}
}

func TestFahrenheitToCelsius(t *testing.T) {
	tests := []struct {
		f    float64
		want float64
	}{
		{32, 0},
		{212, 100},
		{-40, -40},
		{98.6, 37},
		{0, -17.77777777777778},
	}

	for _, tc := range tests {
		got := fahrenheitToCelsius(tc.f)
		if math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("fahrenheitToCelsius(%f) = %f; want %f", tc.f, got, tc.want)
		}
	}
}

func TestIntConversions(t *testing.T) {
	var f1 float64 = 3.99
	if got := int(f1); got != 3 {
		t.Errorf("int(3.99) = %d; want 3", got)
	}
	var f2 float64 = -3.99
	if got := int(f2); got != -3 {
		t.Errorf("int(-3.99) = %d; want -3", got)
	}
	var x uint16 = 256
	if got := byte(x); got != 0 {
		t.Errorf("byte(256) = %d; want 0 (truncation)", got)
	}
	var y float64 = 9_007_199_254_740_993
	if got := int64(y); got != 9007199254740992 {
		t.Errorf("int64(float64(9007199254740993)) = %d; want 9007199254740992 (precision loss)", got)
	}
}

func TestStringConversions(t *testing.T) {
	var r rune = 65
	if got := string(r); got != "A" {
		t.Errorf("string(65) = %q; want \"A\"", got)
	}
	if got := string([]byte("Go")); got != "Go" {
		t.Errorf("string([]byte(\"Go\")) = %q; want \"Go\"", got)
	}
	if got := []byte("Hi"); len(got) != 2 || got[0] != 'H' || got[1] != 'i' {
		t.Errorf("[]byte(\"Hi\") = %v; want [72 105]", got)
	}
}
