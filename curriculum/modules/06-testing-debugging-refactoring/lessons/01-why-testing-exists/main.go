package main

import (
	"errors"
	"fmt"
)

func Add(a, b int) int {
	return a + b
}

func IsEven(n int) bool {
	return n%2 == 0
}

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func main() {
	fmt.Println("Add(2, 3) =", Add(2, 3))
	fmt.Println("IsEven(4) =", IsEven(4))
	fmt.Println("IsEven(5) =", IsEven(5))
	fmt.Println("Divide(10, 3) =", func() string {
		r, err := Divide(10, 3)
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("%d", r)
	}())
	fmt.Println("Divide(5, 0) =", func() string {
		_, err := Divide(5, 0)
		if err != nil {
			return err.Error()
		}
		return "ok"
	}())
}
