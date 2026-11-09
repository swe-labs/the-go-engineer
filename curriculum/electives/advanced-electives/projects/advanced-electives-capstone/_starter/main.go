package main

import (
	"fmt"
)

// RingBuffer defines a generic, thread-safe circular buffer.
// TODO: Implement circular ring buffer with generic constraint [T any].
type RingBuffer[T any] struct {
	capacity int
	items    []T
}

func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	return &RingBuffer[T]{
		capacity: capacity,
		items:    make([]T, 0, capacity),
	}
}

func (r *RingBuffer[T]) Push(item T) bool {
	if len(r.items) >= r.capacity {
		return false
	}
	r.items = append(r.items, item)
	return true
}

func main() {
	fmt.Println("Advanced Electives Capstone Starter")
}
