package main

import "fmt"

type User struct {
	ID    int
	Name  string
	Score float64
	Tags  []string
}

func main() {
	u := User{
		ID:    1,
		Name:  "Alice",
		Score: 95.5,
		Tags:  []string{"go", "debugging", "delve"},
	}
	processUser(u)
}

func processUser(u User) {
	bonus := computeBonus(u.Score)
	fmt.Printf("User %d (%s): final score %.1f\n", u.ID, u.Name, u.Score+bonus)
}

func computeBonus(score float64) float64 {
	if score > 90 {
		return 5.0
	}
	return 1.0
}
