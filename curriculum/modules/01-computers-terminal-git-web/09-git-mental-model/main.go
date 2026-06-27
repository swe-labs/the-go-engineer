package main

import "fmt"

func HashContent(content string) string {
	sum := 0
	for i := 0; i < len(content); i++ {
		sum += int(content[i])
	}
	return fmt.Sprintf("%x", sum)
}

func main() {
	h1 := HashContent("hello world")
	fmt.Printf("Blob hash: %s\n", h1)

	h2 := HashContent("Initial commit")
	fmt.Printf("Commit 1 hash: %s\n", h2)

	h3 := HashContent("Add feature")
	fmt.Printf("Commit 2 hash: %s\n", h3)

	same := HashContent("hello world")
	fmt.Printf("Same content: %s == %s: %v\n", h1, same, h1 == same)
}
