package main

import "fmt"

func main() {
	fmt.Println("fib(5) =", fib(5))
}

func fib(n int) int {
	indent := ""
	for i := 0; i < 5-n; i++ {
		indent += "  "
	}
	fmt.Printf("%sfib(%d) called\n", indent, n)

	if n <= 1 {
		fmt.Printf("%sfib(%d) returns %d\n", indent, n, n)
		return n
	}

	a := fib(n - 1)
	b := fib(n - 2)
	result := a + b
	fmt.Printf("%sfib(%d) returns %d\n", indent, n, result)
	return result
}
