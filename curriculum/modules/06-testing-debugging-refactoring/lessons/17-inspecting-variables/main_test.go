package main

import "testing"

func TestComputeBonusHigh(t *testing.T) {
	got := computeBonus(95.5)
	if got != 5.0 {
		t.Errorf("computeBonus(95.5) = %f, want 5.0", got)
	}
}

func TestComputeBonusLow(t *testing.T) {
	got := computeBonus(70.0)
	if got != 1.0 {
		t.Errorf("computeBonus(70.0) = %f, want 1.0", got)
	}
}

func TestComputeBonusBoundary(t *testing.T) {
	got := computeBonus(90.0)
	if got != 1.0 {
		t.Errorf("computeBonus(90.0) = %f, want 1.0 (not > 90)", got)
	}
}

func TestProcessUser(t *testing.T) {
	// processUser prints to stdout; just verify no panic.
	u := User{ID: 2, Name: "Bob", Score: 80.0, Tags: []string{"test"}}
	processUser(u)
}
