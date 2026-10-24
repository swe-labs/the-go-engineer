package main

import (
	"fmt"
)

// StackValue returns a copy of a local integer.
// The variable 'x' lives strictly inside StackValue's stack frame.
// When StackValue returns, its stack frame is instantly reused.
func StackValue() int {
	x := 42
	return x
}

// HeapValue returns a POINTER (memory address) to a local integer.
// Because the caller needs to access 'x' after HeapValue has returned,
// the Go compiler's escape analysis automatically moves 'x' to the HEAP.
func HeapValue() *int {
	x := 100
	return &x // &x means "the memory address of x"
}

// SliceOnHeap demonstrates dynamic heap allocation.
// The size 'n' is determined at runtime, so the backing array must live on the heap.
func SliceOnHeap(n int) []int {
	numbers := make([]int, n)
	for i := 0; i < n; i++ {
		numbers[i] = (i + 1) * 10
	}
	return numbers
}

func main() {
	fmt.Println("=== Stack vs Heap Memory Preview ===")

	// 1. Stack Allocation
	val := StackValue()
	fmt.Printf("1. StackValue returned integer: %d\n", val)

	// 2. Heap Allocation (Escaped Pointer)
	ptr := HeapValue()
	fmt.Printf("2. HeapValue returned pointer address: %p (Value: %d)\n", ptr, *ptr)

	// 3. Dynamic Heap Slice
	list := SliceOnHeap(3)
	fmt.Printf("3. SliceOnHeap returned %d items: %v\n", len(list), list)

	fmt.Println("\nTakeaway: Local function values live on the fast stack; values that outlive functions move to the heap.")
}
