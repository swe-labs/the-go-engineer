package main

import (
	"errors"
	"fmt"
	"os"
)

func Add(a, b int) int {
	return a + b
}

func Subtract(a, b int) int {
	return a - b
}

func Multiply(a, b int) int {
	return a * b
}

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run . <op> <a> <b>")
		fmt.Println("Ops: add, sub, mul, div")
		return
	}
	op := os.Args[1]
	a := atoi(os.Args[2])
	b := atoi(os.Args[3])

	switch op {
	case "add":
		fmt.Println(Add(a, b))
	case "sub":
		fmt.Println(Subtract(a, b))
	case "mul":
		fmt.Println(Multiply(a, b))
	case "div":
		result, err := Divide(a, b)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		fmt.Println(result)
	default:
		fmt.Fprintln(os.Stderr, "Unknown operation:", op)
		os.Exit(1)
	}
}

func atoi(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	return n
}
