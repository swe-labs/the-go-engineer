package main

import (
	"testing"
)

func TestStackValue(t *testing.T) {
	val := StackValue()
	if val != 42 {
		t.Errorf("expected StackValue to return 42, got %d", val)
	}
}

func TestHeapValue(t *testing.T) {
	ptr := HeapValue()
	if ptr == nil {
		t.Fatalf("expected non-nil pointer from HeapValue")
	}
	if *ptr != 100 {
		t.Errorf("expected pointer value 100, got %d", *ptr)
	}
}

func TestSliceOnHeap(t *testing.T) {
	n := 5
	slice := SliceOnHeap(n)
	if len(slice) != n {
		t.Fatalf("expected slice length %d, got %d", n, len(slice))
	}
	for i := 0; i < n; i++ {
		expected := (i + 1) * 10
		if slice[i] != expected {
			t.Errorf("slice[%d] = %d, expected %d", i, slice[i], expected)
		}
	}
}

func TestMainExecution(t *testing.T) {
	main()
}

