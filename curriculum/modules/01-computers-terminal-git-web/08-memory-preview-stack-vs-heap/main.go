package main

import "fmt"

// StackValue creates a local variable that stays on the stack.
func StackValue() int {
	x := 42
	return x
}

// HeapPointer creates a value and returns its address (escapes to heap).
func HeapPointer() *int {
	x := 42
	return &x
}

// SumNumbers demonstrates heap allocation for dynamic-size data.
func SumNumbers(n int) int {
	nums := make([]int, n)
	for i := range nums {
		nums[i] = i
	}
	sum := 0
	for _, v := range nums {
		sum += v
	}
	return sum
}

func main() {
	s := StackValue()
	p := HeapPointer()
	sum := SumNumbers(5)
	fmt.Printf("Stack: %d, Heap: %d, Sum: %d\n", s, *p, sum)
}
