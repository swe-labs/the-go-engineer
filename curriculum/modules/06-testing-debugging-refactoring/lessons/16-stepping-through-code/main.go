package main

import "fmt"

func main() {
	result := factorial(5)
	fmt.Println("Factorial:", result)
}

func factorial(n int) int {
	fmt.Println("factorial called with:", n)
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}
