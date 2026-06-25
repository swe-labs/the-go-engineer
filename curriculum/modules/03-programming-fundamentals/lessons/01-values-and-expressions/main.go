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
	switch op {
	case '+':
		return a + b
	case '-':
		return a - b
	case '*':
		return a * b
	case '/':
		return a / b
	case '%':
		return a % b
	default:
		panic("unknown operator")
	}
}
