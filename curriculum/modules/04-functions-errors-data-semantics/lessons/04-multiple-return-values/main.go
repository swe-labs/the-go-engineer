package main

import (
	"errors"
	"fmt"
	"strconv"
)

func main() {
	inputs := []string{"42", "-1", "abc", "0"}
	for _, s := range inputs {
		val, err := parsePositive(s)
		if err != nil {
			fmt.Printf("parsePositive(%q) error: %v\n", s, err)
		} else {
			fmt.Printf("parsePositive(%q) = %d\n", s, val)
		}
	}
}

func parsePositive(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if n <= 0 {
		return 0, errors.New("value must be positive")
	}
	return n, nil
}
