package main

import "testing"

func TestStackValue(t *testing.T) {
	got := StackValue()
	if got != 42 {
		t.Errorf("StackValue() = %d, want 42", got)
	}
}

func TestAddNumbers(t *testing.T) {
	got := AddNumbers(10, 20)
	if got != 30 {
		t.Errorf("AddNumbers(10, 20) = %d, want 30", got)
	}
}

func TestAddNumbersNegative(t *testing.T) {
	got := AddNumbers(-5, 10)
	if got != 5 {
		t.Errorf("AddNumbers(-5, 10) = %d, want 5", got)
	}
}

func TestSumArray(t *testing.T) {
	arr := [5]int{0, 1, 2, 3, 4}
	got := SumArray(arr)
	if got != 10 {
		t.Errorf("SumArray([5]int{0,1,2,3,4}) = %d, want 10", got)
	}
}

func TestSumArrayEmpty(t *testing.T) {
	arr := [5]int{}
	got := SumArray(arr)
	if got != 0 {
		t.Errorf("SumArray([5]int{}) = %d, want 0", got)
	}
}

func TestStackValuesAreIndependent(t *testing.T) {
	a := StackValue()
	b := StackValue()
	if b != 42 {
		t.Errorf("stack values should be independent copies, got b = %d", b)
	}
	_ = a
}
