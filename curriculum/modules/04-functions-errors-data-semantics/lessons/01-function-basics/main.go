package main

import "fmt"

func main() {
	fmt.Println("add(3, 5)     =", add(3, 5))
	fmt.Println("subtract(10, 4) =", subtract(10, 4))
	fmt.Println("operate(6, 2, \"mul\") =", operate(6, 2, "mul"))
	fmt.Println("operate(10, 3, \"div\") =", operate(10, 3, "div"))

	double := func(x int) int {
		return x * 2
	}
	fmt.Println("double(9)     =", double(9))
}

func add(a int, b int) int {
	return a + b
}

func subtract(a int, b int) int {
	return a - b
}

func operate(a, b int, op string) int {
	if op == "add" {
		return a + b
	}
	if op == "sub" {
		return a - b
	}
	if op == "mul" {
		return a * b
	}
	if op == "div" {
		return a / b
	}
	panic("unknown operator: " + op)
}
