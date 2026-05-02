// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 04: Types and Design
// Title: Generic Data Structures
// Level: Stretch
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Building type-safe collections (Stack, Queue, Set) using generics.
//   - The mechanics of LIFO (Last-In-First-Out) and FIFO (First-In-First-Out).
//   - Implementing unique sets using `map[T]struct{}` with `comparable` constraints.
//   - Managing generic state in receiver-based data structures.
//
// WHY THIS MATTERS:
//   - Before generics, Go collections either required redundant code
//     generation or the unsafe use of `interface{}`/`any` with runtime
//     type assertions. Generic data structures allow for reusable,
//     performant collections that are verified by the compiler,
//     preventing class-wide categories of runtime type errors.
//
// RUN:
//   go run ./04-types-design/15-generic-data-structures
//
// KEY TAKEAWAY:
//   - Generics enable type-safe collection logic without runtime overhead.
// ============================================================================

// See LICENSE for usage terms.

package main

import "fmt"

// Section 04: Types & Design - Generic Data Structures

// Stack implements a Last-In-First-Out (LIFO) collection for any type T.
// Stack (Type): implements a Last-In-First-Out (LIFO) collection for any type T.
type Stack[T any] struct {
	items []T
}

// Push adds an item to the top of the stack.
// Stack.Push (Method): adds an item to the top of the stack.
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top item from the stack.
// Returns the zero value and false if the stack is empty.
// Stack.Pop (Method): removes and returns the top item from the stack.
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}

// IsEmpty returns true if the stack contains no elements.
// Stack.IsEmpty (Method): returns true if the stack contains no elements.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Size returns the total number of elements in the stack.
// Stack.Size (Method): returns the total number of elements in the stack.
func (s *Stack[T]) Size() int {
	return len(s.items)
}

// Queue implements a First-In-First-Out (FIFO) collection for any type T.
// Queue (Type): implements a First-In-First-Out (FIFO) collection for any type T.
type Queue[T any] struct {
	items []T
}

// Enqueue adds an item to the end of the queue.
// Queue.Enqueue (Method): adds an item to the end of the queue.
func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

// Dequeue removes and returns the front item from the queue.
// Returns the zero value and false if the queue is empty.
// Queue.Dequeue (Method): removes and returns the front item from the queue.
func (q *Queue[T]) Dequeue() (T, bool) {
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// IsEmpty returns true if the queue contains no elements.
// Queue.IsEmpty (Method): returns true if the queue contains no elements.
func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

// Size returns the total number of elements in the queue.
// Queue.Size (Method): returns the total number of elements in the queue.
func (q *Queue[T]) Size() int {
	return len(q.items)
}

// Set implements a unique collection of elements for any comparable type T.
// Set (Type): implements a unique collection of elements for any comparable type T.
type Set[T comparable] map[T]struct{}

// NewSet initializes a new empty set.
// NewSet (Function): initializes a new empty set.
func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

// Add inserts an item into the set.
// Set.Add (Method): inserts an item into the set.
func (s Set[T]) Add(item T) {
	s[item] = struct{}{}
}

// Contains returns true if the item exists in the set.
// Set.Contains (Method): returns true if the item exists in the set.
func (s Set[T]) Contains(item T) bool {
	_, ok := s[item]
	return ok
}

// Remove deletes an item from the set.
// Set.Remove (Method): deletes an item from the set.
func (s Set[T]) Remove(item T) {
	delete(s, item)
}

// Size returns the number of elements in the set.
// Set.Size (Method): returns the number of elements in the set.
func (s Set[T]) Size() int {
	return len(s)
}

// Items returns all elements in the set as a slice.
// Set.Items (Method): returns all elements in the set as a slice.
func (s Set[T]) Items() []T {
	result := make([]T, 0, len(s))
	for item := range s {
		result = append(result, item)
	}
	return result
}

func main() {
	fmt.Println("=== Generic Data Structures: Type-Safe Collections ===")
	fmt.Println()

	// 1. Generic Stack (LIFO).
	// The Stack handles type safety at compile time. No assertions needed.
	fmt.Println("--- Stack Behavior (int) ---")
	stack := Stack[int]{}
	stack.Push(10)
	stack.Push(20)
	val, ok := stack.Pop()
	fmt.Printf("  Popped value: %d (success: %t)\n", val, ok)
	fmt.Println()

	// 2. Generic Queue (FIFO).
	// We can reuse the same queue logic for strings without code duplication.
	fmt.Println("--- Queue Behavior (string) ---")
	queue := Queue[string]{}
	queue.Enqueue("job_1")
	queue.Enqueue("job_2")
	if v, ok := queue.Dequeue(); ok {
		fmt.Printf("  Processing: %s\n", v)
	}
	fmt.Println()

	// 3. Generic Set (comparable).
	// The Set uses the 'comparable' constraint to ensure elements can be hashed.
	fmt.Println("--- Set Behavior (comparable) ---")
	names := NewSet[string]()
	names.Add("alice")
	names.Add("bob")
	names.Add("alice") // Duplicate ignored
	fmt.Printf("  Unique Items: %v (Size: %d)\n", names.Items(), names.Size())

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: TI.12 -> 04-types-design/12-functional-options")
	fmt.Println("Run    : go run ./04-types-design/12-functional-options")
	fmt.Println("Current: TI.15 (generic-data-structures)")
	fmt.Println("---------------------------------------------------")
}
