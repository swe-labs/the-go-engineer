package main

import "fmt"

// Section 00: How Computers Work

func main() {
	// We create a variable (data)
	score := 100

	// We print the value
	fmt.Printf("The value of score is: %d\n", score)

	// We print the exact physical RAM address where 'score' lives
	fmt.Printf("The memory address of score (Stack) is: %p\n", &score)

	// We create another variable explicitly requesting permanent memory
	heapData := new(int)
	*heapData = 500
	fmt.Printf("The memory address of heapData (Heap) is: %p\n", heapData)

	fmt.Println("NEXT UP: HC.4 terminal-confidence")
}
