package main

import "fmt"

func main() {
	msg := greeting("Alice")
	fmt.Println(msg)

	result := compute(10, 5)
	fmt.Println("Result:", result)

	fmt.Println("Done.")
}

func greeting(name string) string {
	return "Hello, " + name + "!"
}

func compute(a, b int) int {
	sum := a + b
	diff := a - b
	product := a * b
	return sum + diff + product
}
