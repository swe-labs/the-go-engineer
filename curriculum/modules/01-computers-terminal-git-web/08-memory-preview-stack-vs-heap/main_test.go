package main

import "testing"

func TestStackValue(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"returns local value 42", 42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StackValue(); got != tt.want {
				t.Errorf("StackValue() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestHeapPointer(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"returns pointer to value 42", 42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HeapPointer()
			if got == nil {
				t.Fatal("HeapPointer() returned nil")
			}
			if *got != tt.want {
				t.Errorf("HeapPointer() = %d, want %d", *got, tt.want)
			}
		})
	}
}

func TestSumNumbers(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{"zero elements sum to 0", 0, 0},
		{"one element sums to 0", 1, 0},
		{"five elements sum to 10", 5, 10},
		{"ten elements sum to 45", 10, 45},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SumNumbers(tt.n); got != tt.want {
				t.Errorf("SumNumbers(%d) = %d, want %d", tt.n, got, tt.want)
			}
		})
	}
}

func TestStackValuesAreIndependent(t *testing.T) {
	a := StackValue()
	b := StackValue()
	_ = a
	if b != 42 {
		t.Errorf("stack values should be independent copies, got b = %d", b)
	}
}

func TestHeapPointersAreIndependent(t *testing.T) {
	p := HeapPointer()
	q := HeapPointer()
	*p = 100
	if *q != 42 {
		t.Errorf("heap pointers should be independent allocations, got *q = %d", *q)
	}
}
