package main

import "fmt"

func main() {
	users := []string{"Alice", "Bob", "Charlie", "Diana"}

	for i, name := range users {
		processUser(i, name)
	}

	fmt.Println("All users processed.")
}

func processUser(id int, name string) {
	score := computeScore(id)
	fmt.Printf("User %d (%s): score=%d\n", id, name, score)
}

func computeScore(id int) int {
	base := id * 10
	bonus := id * id
	return base + bonus
}

func multiply(x, y int) int {
	return x * y
}
