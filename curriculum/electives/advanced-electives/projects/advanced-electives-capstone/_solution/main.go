package main

import (
	"fmt"
	"sync"
)

// RingBuffer defines a generic, thread-safe circular buffer.
type RingBuffer[T any] struct {
	mu       sync.RWMutex
	capacity int
	items    []T
	head     int
	tail     int
	size     int
}

func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	return &RingBuffer[T]{
		capacity: capacity,
		items:    make([]T, capacity),
	}
}

func (r *RingBuffer[T]) Push(item T) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.size == r.capacity {
		return false
	}
	r.items[r.tail] = item
	r.tail = (r.tail + 1) % r.capacity
	r.size++
	return true
}

func (r *RingBuffer[T]) Pop() (T, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var zero T
	if r.size == 0 {
		return zero, false
	}
	val := r.items[r.head]
	r.items[r.head] = zero
	r.head = (r.head + 1) % r.capacity
	r.size--
	return val, true
}

func (r *RingBuffer[T]) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.size
}

func main() {
	rb := NewRingBuffer[string](3)
	rb.Push("task-1")
	rb.Push("task-2")
	val, _ := rb.Pop()
	fmt.Printf("Popped from generic ring buffer: %s (remaining: %d)\n", val, rb.Len())
}
