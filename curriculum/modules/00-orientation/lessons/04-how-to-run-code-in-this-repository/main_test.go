package main

import (
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{a: 2, b: 3, want: 5},
		{a: -1, b: 1, want: 0},
		{a: 0, b: 0, want: 0},
		{a: 100, b: 200, want: 300},
	}
	for _, tc := range tests {
		got := Add(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("Add(%d, %d) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{a: 10, b: 3, want: 7},
		{a: 5, b: 10, want: -5},
		{a: 0, b: 0, want: 0},
	}
	for _, tc := range tests {
		got := Subtract(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("Subtract(%d, %d) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{a: 3, b: 4, want: 12},
		{a: 0, b: 5, want: 0},
		{a: -2, b: 3, want: -6},
	}
	for _, tc := range tests {
		got := Multiply(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("Multiply(%d, %d) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		a, b    int
		want    int
		wantErr bool
	}{
		{a: 10, b: 2, want: 5, wantErr: false},
		{a: 7, b: 3, want: 2, wantErr: false},
		{a: 5, b: 0, want: 0, wantErr: true},
	}
	for _, tc := range tests {
		got, err := Divide(tc.a, tc.b)
		if tc.wantErr {
			if err == nil {
				t.Errorf("Divide(%d, %d) expected error", tc.a, tc.b)
			}
			continue
		}
		if err != nil {
			t.Errorf("Divide(%d, %d) unexpected error: %v", tc.a, tc.b, err)
		}
		if got != tc.want {
			t.Errorf("Divide(%d, %d) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestAtoi(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{input: "42", want: 42},
		{input: "0", want: 0},
		{input: "abc", want: 0},
		{input: "", want: 0},
	}
	for _, tc := range tests {
		got := atoi(tc.input)
		if got != tc.want {
			t.Errorf("atoi(%q) = %d, want %d", tc.input, got, tc.want)
		}
	}
}
