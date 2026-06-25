package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <name> [age]")
		return
	}
	name := os.Args[1]
	if len(os.Args) >= 3 {
		age, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Invalid age %q: must be a number\n", os.Args[2])
			return
		}
		fmt.Printf("Hello %s, age %d\n", name, age)
		return
	}
	fmt.Printf("Hello %s\n", name)
}
