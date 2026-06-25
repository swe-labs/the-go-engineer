package main

import (
	"errors"
	"fmt"
	"strconv"
)

func classifyTemperature(temp int) string {
	if temp > 40 {
		return "dangerously hot"
	} else if temp > 30 {
		return "hot"
	} else if temp > 20 {
		return "warm"
	} else if temp > 10 {
		return "mild"
	} else if temp >= 0 {
		return "cool"
	}
	return "freezing"
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func validateAge(s string) error {
	if age, err := strconv.Atoi(s); err != nil {
		return fmt.Errorf("invalid age: %w", err)
	} else if age < 0 || age > 150 {
		return fmt.Errorf("age %d out of range [0, 150]", age)
	} else if age < 18 {
		return fmt.Errorf("age %d is under 18", age)
	}
	return nil
}

func gradeClassification(score int) string {
	if score < 0 || score > 100 {
		return "invalid"
	} else if score >= 90 {
		return "A"
	} else if score >= 80 {
		return "B"
	} else if score >= 70 {
		return "C"
	} else if score >= 60 {
		return "D"
	}
	return "F"
}

func main() {
	for _, t := range []int{45, 32, 22, 15, 5, -5} {
		fmt.Printf("%3dC -> %s\n", t, classifyTemperature(t))
	}

	if result, err := divide(10, 3); err != nil {
		fmt.Println("divide error:", err)
	} else {
		fmt.Printf("10 / 3 = %.2f\n", result)
	}

	ages := []string{"25", "abc", "200", "16"}
	for _, s := range ages {
		if err := validateAge(s); err != nil {
			fmt.Printf("age %q invalid: %v\n", s, err)
		} else {
			fmt.Printf("age %q valid\n", s)
		}
	}

	scores := []int{-5, 0, 45, 70, 85, 95, 100, 101}
	for _, s := range scores {
		fmt.Printf("score %d -> %s\n", s, gradeClassification(s))
	}
}
