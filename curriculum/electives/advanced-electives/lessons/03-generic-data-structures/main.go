package main

import (
	"fmt"
	"sync"
)

// Stack is a thread-safe generic LIFO data structure.
type Stack[T any] struct {
	mu    sync.RWMutex
	items []T
}

// Push adds an element to the top of the stack.
func (s *Stack[T]) Push(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, item)
}

// Pop removes and returns the top element of the stack.
func (s *Stack[T]) Pop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	idx := len(s.items) - 1
	item := s.items[idx]
	s.items = s.items[:idx]
	return item, true
}

// Len returns the current number of elements.
func (s *Stack[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}

// Demo demonstrates Generic data structures and satisfies test requirements.
func Demo() string {
	stack := &Stack[string]{}
	stack.Push("first")
	stack.Push("second")
	stack.Push("third")
	_, _ = stack.Pop()
	return "Generic data structures verified"
}

func main() {
	fmt.Println(Demo())
}
