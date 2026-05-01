// Copyright (c) 2026 Rasel Hossen

// ============================================================================
// Section 04: Types and Design
// Title: Generic Data Structures
// Level: Stretch
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn to build type-safe generic data structures like Stack, Queue, and Set using Go's generics.
//
// WHY THIS MATTERS:
//   - Think of a reusable storage box. Without generics, you'd need separate boxes for books, clothes, and electronics. With generics, one "Box<T>" works...
//
// RUN:
//   go run ./04-types-design/15-generic-data-structures
//
// KEY TAKEAWAY:
//   - Learn to build type-safe generic data structures like Stack, Queue, and Set using Go's generics.
// ============================================================================

// See LICENSE for usage terms.

package main

import "fmt"

//
//   - Generic Stack implementation
//   - Generic Queue implementation
//   - Generic Set using map
//

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack[T]) Size() int {
	return len(s.items)
}

type Queue[T any] struct {
	items []T
}

func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

func (q *Queue[T]) Size() int {
	return len(q.items)
}

type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

func (s Set[T]) Add(item T) {
	s[item] = struct{}{}
}

func (s Set[T]) Contains(item T) bool {
	_, ok := s[item]
	return ok
}

func (s Set[T]) Remove(item T) {
	delete(s, item)
}

func (s Set[T]) Size() int {
	return len(s)
}

func (s Set[T]) Items() []T {
	result := make([]T, 0, len(s))
	for item := range s {
		result = append(result, item)
	}
	return result
}

func main() {
	fmt.Println("=== Generic Data Structures ===")
	fmt.Println()

	fmt.Println("--- Generic Stack ---")
	stack := Stack[int]{}
	stack.Push(10)
	stack.Push(20)
	stack.Push(30)
	fmt.Printf("  Stack size: %d\n", stack.Size())
	val, ok := stack.Pop()
	fmt.Printf("  Popped: %d, ok: %t\n", val, ok)
	fmt.Printf("  Stack size after pop: %d\n", stack.Size())

	fmt.Println()
	fmt.Println("--- Generic Stack with Strings ---")
	strStack := Stack[string]{}
	strStack.Push("first")
	strStack.Push("second")
	strStack.Push("third")
	for !strStack.IsEmpty() {
		val, _ := strStack.Pop()
		fmt.Printf("  Popped: %s\n", val)
	}

	fmt.Println()
	fmt.Println("--- Generic Queue ---")
	queue := Queue[string]{}
	queue.Enqueue("one")
	queue.Enqueue("two")
	queue.Enqueue("three")
	if val, ok := queue.Dequeue(); ok {
		fmt.Printf("  Dequeued: %s\n", val)
	}
	if val, ok := queue.Dequeue(); ok {
		fmt.Printf("  Dequeued: %s\n", val)
	}
	fmt.Printf("  Remaining: %d\n", queue.Size())

	fmt.Println()
	fmt.Println("--- Generic Set ---")
	names := NewSet[string]()
	names.Add("alice")
	names.Add("bob")
	names.Add("alice")
	fmt.Printf("  Contains 'alice': %t\n", names.Contains("alice"))
	fmt.Printf("  Contains 'charlie': %t\n", names.Contains("charlie"))
	fmt.Printf("  Set size: %d\n", names.Size())
	names.Remove("bob")
	fmt.Printf("  After removing bob, size: %d\n", names.Size())
	fmt.Printf("  All items: %v\n", names.Items())

	fmt.Println()
	fmt.Println("--- Set with Integers ---")
	intSet := NewSet[int]()
	intSet.Add(1)
	intSet.Add(2)
	intSet.Add(3)
	intSet.Add(2)
	fmt.Printf("  Int set: %v\n", intSet.Items())

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Generic data structures provide type safety at compile time")
	fmt.Println("  - Stack: LIFO - Push/Pop")
	fmt.Println("  - Queue: FIFO - Enqueue/Dequeue")
	fmt.Println("  - Set: unique values using map[T]struct{}")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: CO.1 -> 04-types-design/composition/1-composition")
	fmt.Println("Current: TI.15 (generic-data-structures)")
	fmt.Println("Previous: TI.14 (complex-generic-constraints)")
	fmt.Println("---------------------------------------------------")
}
