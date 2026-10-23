package main

import (
	"cmp"
	"fmt"
)

// Map transforms a slice of type T into a slice of type U using a mapping function.
func Map[T, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

// Filter returns elements of a slice that satisfy the predicate.
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Contains checks if a target value exists in a slice of comparable types.
func Contains[T comparable](slice []T, target T) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

// Min returns the smaller of two values that satisfy the cmp.Ordered constraint.
func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Stack is a type-safe generic LIFO stack data structure.
type Stack[T any] struct {
	items []T
}

// Push adds an item to the top of the stack.
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top item from the stack.
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	lastIdx := len(s.items) - 1
	item := s.items[lastIdx]
	s.items = s.items[:lastIdx]
	return item, true
}

// Len returns the current number of items in the stack.
func (s *Stack[T]) Len() int {
	return len(s.items)
}

func main() {
	fmt.Println("=== 1. Generic Slice Transformations (Map & Filter) ===")
	nums := []int{1, 2, 3, 4, 5, 6}
	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	squared := Map(evens, func(n int) int { return n * n })
	fmt.Printf("Original: %v\nEvens:    %v\nSquared:  %v\n", nums, evens, squared)

	fmt.Println("\n=== 2. Comparable Constraint (Contains) ===")
	roles := []string{"admin", "member", "auditor"}
	fmt.Printf("Contains 'admin': %t\n", Contains(roles, "admin"))
	fmt.Printf("Contains 'guest': %t\n", Contains(roles, "guest"))

	fmt.Println("\n=== 3. Ordered Constraint (cmp.Ordered) ===")
	fmt.Printf("Min(42, 17):     %d\n", Min(42, 17))
	fmt.Printf("Min(3.14, 2.71): %.2f\n", Min(3.14, 2.71))
	fmt.Printf("Min('b', 'a'):   %c\n", Min('b', 'a'))

	fmt.Println("\n=== 4. Generic Type (Stack[T]) ===")
	strStack := &Stack[string]{}
	strStack.Push("first")
	strStack.Push("second")
	strStack.Push("third")

	for strStack.Len() > 0 {
		val, _ := strStack.Pop()
		fmt.Printf("Popped: %s\n", val)
	}
}
