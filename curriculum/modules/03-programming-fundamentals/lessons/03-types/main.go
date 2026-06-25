package main

import (
	"errors"
	"fmt"
)

type Celsius float64
type Fahrenheit float64

func main() {
	fmt.Println("=== basic types ===")
	var a int = 42
	var b float64 = 3.14
	var c string = "hello"
	var d bool = true
	fmt.Printf("int: %d, float64: %.2f, string: %s, bool: %t\n", a, b, c, d)

	fmt.Println("=== type inference ===")
	x := 42
	y := 3.14
	z := "hello"
	ok := true
	fmt.Printf("inferred: %T %T %T %T\n", x, y, z, ok)

	fmt.Println("=== byte and rune ===")
	var char byte = 'A'
	var code rune = '世'
	fmt.Printf("byte: %c (%d), rune: %c (U+%04x)\n", char, char, code, code)

	fmt.Println("=== type definition ===")
	var temp Celsius = 100.0
	fmt.Printf("Celsius: %.1f°C\n", temp)

	fmt.Println("=== type alias ===")
	var b1 byte = 255
	var b2 uint8 = b1
	fmt.Printf("byte == uint8: %d == %d\n", b1, b2)

	fmt.Println("=== temperature conversions ===")
	val, err := convertUnits(100, "C", "F")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("100 C = %.2f F\n", val)
	}

	val, err = convertUnits(32, "F", "C")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("32 F = %.2f C\n", val)
	}

	val, err = convertUnits(0, "C", "K")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("0 C = %.2f K\n", val)
	}

	val, err = convertUnits(300, "K", "C")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("300 K = %.2f C\n", val)
	}

	val, err = convertUnits(212, "F", "K")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("212 F = %.2f K\n", val)
	}

	_, err = convertUnits(100, "C", "X")
	if err != nil {
		fmt.Println("Expected error:", err)
	}
}

func convertUnits(value float64, from, to string) (float64, error) {
	if from == to {
		return value, nil
	}

	var celsius float64
	switch from {
	case "C":
		celsius = value
	case "F":
		celsius = (value - 32) * 5 / 9
	case "K":
		celsius = value - 273.15
	default:
		return 0, errors.New("unknown unit: " + from)
	}

	switch to {
	case "C":
		return celsius, nil
	case "F":
		return celsius*9/5 + 32, nil
	case "K":
		return celsius + 273.15, nil
	default:
		return 0, errors.New("unknown unit: " + to)
	}
}
