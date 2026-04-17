package main

import "fmt"

// Section 00: How Computers Work

func main() {
	// We create a variable (data)
	score := 100

	// We print the value
	fmt.Printf("The value of score is: %d\n", score)

	// We print the exact physical RAM address where 'score' lives
	fmt.Printf("The memory address of score is: %p\n", &score)
	fmt.Println("NEXT UP: HC.4 terminal-confidence")
}
