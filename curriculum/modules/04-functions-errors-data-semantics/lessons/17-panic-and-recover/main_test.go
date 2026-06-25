package main

import (
	"testing"
)

func TestSafeDivide_Success(t *testing.T) {
	res, err := safeDivide(10, 2)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if res != 5 {
		t.Errorf("expected 5, got %d", res)
	}
}

func TestSafeDivide_ByZero(t *testing.T) {
	_, err := safeDivide(10, 0)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "panic: division by zero" {
		t.Errorf("unexpected message: %v", err)
	}
}

func TestSafeBatch_Success(t *testing.T) {
	results, err := safeBatch([]int{10, 20, 30}, []int{2, 4, 5})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	expected := []int{5, 5, 6}
	for i, v := range results {
		if v != expected[i] {
			t.Errorf("index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestSafeBatch_WithDivisionByZero(t *testing.T) {
	results, err := safeBatch([]int{10, 20, 30}, []int{2, 0, 3})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	// The first and third should succeed, second should error
	if results[0] != 5 {
		t.Errorf("expected results[0]=5, got %d", results[0])
	}
	if results[2] != 10 {
		t.Errorf("expected results[2]=10, got %d", results[2])
	}
}

func TestSafeBatch_MismatchedLengths(t *testing.T) {
	_, err := safeBatch([]int{1, 2}, []int{1})
	if err == nil {
		t.Fatal("expected error for mismatched lengths")
	}
}

func TestSafeBatch_Empty(t *testing.T) {
	results, err := safeBatch([]int{}, []int{})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty results, got %v", results)
	}
}

func TestRecoverOnlyWorksInDefer(t *testing.T) {
	// Calling recover() directly returns nil
	r := recover()
	if r != nil {
		t.Error("expected nil from recover outside defer")
	}
}
