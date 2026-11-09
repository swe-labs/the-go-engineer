package main

import (
	"fmt"
	"time"
)

type task func(int) int

func sequential(tasks []task, input int) int {
	result := input
	for _, t := range tasks {
		result = t(result)
	}
	return result
}

func slowSquare(n int) int {
	time.Sleep(10 * time.Millisecond)
	return n * n
}

func slowDouble(n int) int {
	time.Sleep(10 * time.Millisecond)
	return n * 2
}

func main() {
	tasks := []task{slowSquare, slowDouble, slowSquare}

	start := time.Now()
	result := sequential(tasks, 3)
	elapsed := time.Since(start)
	fmt.Printf("Sequential: %d (took %v)\n", result, elapsed)
}
