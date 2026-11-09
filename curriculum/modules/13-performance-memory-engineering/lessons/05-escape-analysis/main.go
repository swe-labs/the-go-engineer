package main

import "fmt"

type User struct {
	ID   int
	Name string
}

// stackAllocated creates a User that stays on the stack.
func stackAllocated(id int, name string) User {
	return User{ID: id, Name: name}
}

// heapAllocated creates a User that escapes to the heap.
func heapAllocated(id int, name string) *User {
	return &User{ID: id, Name: name}
}

func main() {
	a := stackAllocated(1, "Alice")
	b := heapAllocated(2, "Bob")

	fmt.Println(a, b)

	// To see escape analysis decisions, run:
	// go build -gcflags='-m -m' ./curriculum/modules/13-performance-memory-engineering/lessons/05-escape-analysis
}
