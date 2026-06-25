package main

import "testing"

func TestAdd(t *testing.T) {
	got := Add(2, 3)
	want := 5
	if got != want {
		t.Errorf("Add(2, 3) = %d; want %d", got, want)
	}
}

func TestIsEven(t *testing.T) {
	if !IsEven(4) {
		t.Error("IsEven(4) = false; want true")
	}
	if IsEven(5) {
		t.Error("IsEven(5) = true; want false")
	}
}

func TestDivide(t *testing.T) {
	got, err := Divide(10, 3)
	if err != nil {
		t.Fatalf("Divide(10, 3) unexpected error: %v", err)
	}
	if got != 3 {
		t.Errorf("Divide(10, 3) = %d; want 3", got)
	}
}

func TestDivideByZero(t *testing.T) {
	_, err := Divide(5, 0)
	if err == nil {
		t.Fatal("Divide(5, 0) expected error, got nil")
	}
	if err.Error() != "division by zero" {
		t.Errorf("Divide(5, 0) error = %q; want %q", err.Error(), "division by zero")
	}
}
