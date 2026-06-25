package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseCelsius(input string) (float64, error) {
	s := strings.TrimSpace(input)
	if s == "" {
		return 0, fmt.Errorf("empty input")
	}
	if strings.HasSuffix(s, "C") {
		s = strings.TrimSuffix(s, "C")
	} else if strings.HasSuffix(s, "F") {
		return 0, fmt.Errorf("expected Celsius, got Fahrenheit: %q", input)
	} else if strings.HasSuffix(s, "K") {
		return 0, fmt.Errorf("expected Celsius, got Kelvin: %q", input)
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing %q: %w", input, err)
	}
	return v, nil
}

func fahrenheitToCelsius(f float64) float64 {
	return (f - 32) * 5 / 9
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <temperature>")
		fmt.Println("Examples: 23.5C, 98.6F, -10")
		return
	}

	input := os.Args[1]
	if strings.HasSuffix(input, "F") {
		f, err := strconv.ParseFloat(strings.TrimSuffix(input, "F"), 64)
		if err != nil {
			fmt.Printf("Error parsing Fahrenheit value: %v\n", err)
			return
		}
		c := fahrenheitToCelsius(f)
		fmt.Printf("%s = %.2fC\n", input, c)
		return
	}

	c, err := parseCelsius(input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2fC\n", c)
}
