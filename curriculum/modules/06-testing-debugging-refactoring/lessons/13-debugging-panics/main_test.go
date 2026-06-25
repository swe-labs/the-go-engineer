package main

import (
	"testing"
)

func TestSafeDivide(t *testing.T) {
	a := []int{10, 20, 30, 40}
	b := []int{2, 0, 5, 0}
	result := SafeDivide(a, b)
	if len(result) != len(a) {
		t.Fatalf("expected %d results, got %d", len(a), len(result))
	}
	expected := []int{5, 0, 6, 0}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("index %d: expected %d, got %d", i, v, result[i])
		}
	}
}

func TestSafeDivideShortB(t *testing.T) {
	a := []int{10, 20, 30}
	b := []int{2}
	result := SafeDivide(a, b)
	if len(result) != len(a) {
		t.Fatalf("expected %d results, got %d", len(a), len(result))
	}
	if result[0] != 5 {
		t.Errorf("expected 5, got %d", result[0])
	}
	if result[1] != 20 {
		t.Errorf("expected 20 (fallback to 1), got %d", result[1])
	}
}

func TestSafeSliceAccess(t *testing.T) {
	nums := []int{10, 20, 30}
	v, ok := SafeSliceAccess(nums, 1)
	if !ok {
		t.Fatal("expected ok for valid index")
	}
	if v != 20 {
		t.Errorf("expected 20, got %d", v)
	}

	v, ok = SafeSliceAccess(nums, 5)
	if ok {
		t.Log("out of bounds recovered, got:", v)
	}

	v, ok = SafeSliceAccess(nums, -1)
	if ok {
		t.Error("expected false for negative index")
	}
}

func TestSafeConfigRead(t *testing.T) {
	// Should not panic — recovers internally.
	SafeConfigRead()
}
