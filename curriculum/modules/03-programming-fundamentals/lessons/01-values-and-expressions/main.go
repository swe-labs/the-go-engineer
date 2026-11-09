package main

import "fmt"

func main() {
	fmt.Println(evaluate(10, 5, '+'))
	fmt.Println(evaluate(10, 5, '-'))
	fmt.Println(evaluate(10, 5, '*'))
	fmt.Println(evaluate(10, 5, '/'))
	fmt.Println(evaluate(10, 5, '%'))
}

func evaluate(a, b int, op rune) int {
	if op == '+' {
		return a + b
	}
	if op == '-' {
		return a - b
	}
	if op == '*' {
		return a * b
	}
	if op == '/' {
		return a / b
	}
	if op == '%' {
		return a % b
	}
	return 0
}
