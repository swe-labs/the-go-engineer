package main

import (
	"fmt"
)

func countdown() {
	fmt.Println("start")
	for i := 1; i <= 3; i++ {
		defer fmt.Println(i)
	}
	fmt.Println("launch")
}

func addLogger(fn func(int) int) func(int) (int, error) {
	return func(x int) (result int, err error) {
		defer func() {
			fmt.Printf("called with %d, result %d\n", x, result)
		}()
		result = fn(x)
		return result, nil
	}
}

func square(x int) int {
	return x * x
}

func main() {
	fmt.Println("=== countdown ===")
	countdown()

	fmt.Println("\n=== addLogger ===")
	loggedSquare := addLogger(square)
	res, err := loggedSquare(5)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("returned:", res)
	}

	res, err = loggedSquare(-3)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("returned:", res)
	}
}
