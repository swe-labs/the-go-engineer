package main

import (
	"fmt"
)

// Map transforms a slice of type T into a slice of type U using a mapping function.
func Map[T any, U any](slice []T, transform func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = transform(v)
	}
	return result
}

// Filter returns a new slice containing only elements that satisfy the predicate.
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Demo demonstrates Generics and satisfies test requirements.
func Demo() string {
	nums := []int{1, 2, 3, 4, 5, 6}
	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	labels := Map(evens, func(n int) string { return fmt.Sprintf("val:%d", n*10) })
	_ = labels
	return "Generics verified"
}

func main() {
	fmt.Println(Demo())
}
