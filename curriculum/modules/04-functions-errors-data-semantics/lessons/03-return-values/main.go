package main

import "fmt"

func main() {
	fmt.Println("factorial(5) =", factorial(5))
	fmt.Println("factorial(0) =", factorial(0))
	fmt.Println("factorial(1) =", factorial(1))
	fmt.Println("safeDivide(10, 3) =", safeDivide(10, 3))
	fmt.Println("safeDivide(10, 0) =", safeDivide(10, 0))
}

func factorial(n int) (result int) {
	if n <= 1 {
		return 1
	}
	result = n * factorial(n-1)
	return
}

func safeDivide(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}
