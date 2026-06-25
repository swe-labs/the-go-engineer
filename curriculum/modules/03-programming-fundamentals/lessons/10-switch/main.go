package main

import (
	"errors"
	"fmt"
)

func calc(a, b float64, op string) (float64, error) {
	switch op {
	case "add", "sum", "+":
		return a + b, nil
	case "sub", "diff", "-":
		return a - b, nil
	case "mul", "prod", "*":
		return a * b, nil
	case "div", "/":
		if b == 0 {
			return 0, errors.New("division by zero")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("unknown operator: %s", op)
	}
}

func describe(v interface{}) string {
	switch v := v.(type) {
	case int:
		return fmt.Sprintf("int: %d", v)
	case float64:
		return fmt.Sprintf("float: %f", v)
	case string:
		return fmt.Sprintf("string: %s", v)
	case bool:
		return fmt.Sprintf("bool: %t", v)
	default:
		return fmt.Sprintf("unknown type: %T", v)
	}
}

func main() {
	ops := []string{"add", "sub", "mul", "div", "+", "-", "*", "/", "unknown"}
	for _, op := range ops {
		result, err := calc(10, 3, op)
		if err != nil {
			fmt.Printf("calc(10, 3, %q) error: %v\n", op, err)
		} else {
			fmt.Printf("calc(10, 3, %q) = %f\n", op, result)
		}
	}

	fmt.Println()
	fmt.Println(describe(42))
	fmt.Println(describe(3.14))
	fmt.Println(describe("hello"))
	fmt.Println(describe(true))
	fmt.Println(describe(struct{ name string }{"test"}))
}
