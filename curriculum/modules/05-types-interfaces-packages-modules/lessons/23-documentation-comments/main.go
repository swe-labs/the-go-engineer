package main

import (
	"fmt"

	"github.com/rasel9t6/the-go-engineer/curriculum/modules/05-types-interfaces-packages-modules/lessons/23-documentation-comments/calc"
)

func main() {
	// Test Sum.
	fmt.Println("Sum(1,2,3) =", calc.Sum(1, 2, 3))

	// Test Calculator.
	c := calc.NewCalculator()
	c.Add(10)
	c.Add(20)
	fmt.Println("Calculator result:", c.Result())

	// Test deprecated function.
	fmt.Println("OldSum(4,5) =", calc.OldSum(4, 5))
}
