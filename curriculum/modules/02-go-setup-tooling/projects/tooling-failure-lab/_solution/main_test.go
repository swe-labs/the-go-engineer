package main

import "testing"

func TestMultiply(t *testing.T) {
	got := Multiply(3, 4)
	want := 12
	if got != want {
		t.Errorf("Multiply(3, 4) = %d; want %d", got, want)
	}
}

func TestMultiplyZero(t *testing.T) {
	got := Multiply(0, 7)
	want := 0
	if got != want {
		t.Errorf("Multiply(0, 7) = %d; want %d", got, want)
	}
}

func TestMultiplyNegative(t *testing.T) {
	got := Multiply(-2, 5)
	want := -10
	if got != want {
		t.Errorf("Multiply(-2, 5) = %d; want %d", got, want)
	}
}
