// Package calc provides basic arithmetic accumulation.
//
// Use [Sum] for a one-shot addition of multiple values, or [Calculator]
// for stateful accumulation with step-by-step operations.
package calc

import "fmt"

// Sum returns the sum of all provided values.
//
// Example:
//
//	total := calc.Sum(1, 2, 3, 4)
//	fmt.Println(total) // 10
func Sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Calculator maintains a running total and provides methods
// to manipulate it. Use [NewCalculator] to create one.
type Calculator struct {
	total int
}

// NewCalculator creates a Calculator with an initial value of zero.
func NewCalculator() *Calculator {
	return &Calculator{}
}

// Add adds n to the running total.
func (c *Calculator) Add(n int) {
	c.total += n
}

// Result returns the current running total.
func (c *Calculator) Result() int {
	return c.total
}

// OldSum sums numbers using the outdated API.
//
// Deprecated: Use [Sum] instead.
func OldSum(nums ...int) int {
	fmt.Println("OldSum called — this function is deprecated")
	return Sum(nums...)
}
