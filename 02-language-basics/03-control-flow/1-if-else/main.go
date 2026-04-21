// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import "fmt"

func main() {
	temperature := 25

	if temperature > 30 {
		fmt.Println("Temperature is above 30C.")
	} else {
		fmt.Println("Temperature is 30C or below.")
	}

	score := 85

	if score >= 90 {
		fmt.Println("Grade: A")
	} else if score >= 80 {
		fmt.Println("Grade: B")
	} else if score >= 70 {
		fmt.Println("Grade: C")
	} else {
		fmt.Println("Grade: F")
	}

	username := ""

	if username == "" {
		fmt.Println("Username is missing.")
	} else {
		fmt.Println("Username is present.")
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: CF.2 for-basics")
	fmt.Println("Current: CF.1 (if / else)")
	fmt.Println("---------------------------------------------------")
}
