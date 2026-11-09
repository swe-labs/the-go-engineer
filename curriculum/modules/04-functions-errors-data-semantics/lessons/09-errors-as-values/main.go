package main

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrInvalidAge      = errors.New("invalid age")
	ErrNegativeAge     = errors.New("age must be positive")
	ErrUnreasonableAge = errors.New("age exceeds 150")
)

func main() {
	inputs := []string{"25", "-1", "abc", "200"}
	for _, s := range inputs {
		age, err := parseAge(s)
		if err != nil {
			fmt.Printf("parseAge(%q) error: %v\n", s, err)
		} else {
			fmt.Printf("parseAge(%q) = %d\n", s, age)
		}
	}
}

func parseAge(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, ErrInvalidAge
	}
	if n < 0 {
		return 0, ErrNegativeAge
	}
	if n > 150 {
		return 0, ErrUnreasonableAge
	}
	return n, nil
}
