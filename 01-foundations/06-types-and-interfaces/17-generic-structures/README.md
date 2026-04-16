# TI.17 Generic Data Structures

## Mission

Learn to build type-safe generic data structures like Stack, Queue, and Set using Go's generics.

## Why This Lesson Exists Now

You've learned generic functions. Now learn to build generic data structures that are type-safe at compile time—no runtime type assertions needed.

## Prerequisites

- `TI.9` generics
- `TI.16` complex constraints

## Mental Model

Think of a reusable storage box. Without generics, you'd need separate boxes for books, clothes, and electronics. With generics, one "Box<T>" works for all—type-safe and efficient.

## Visual Model

```go
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T)
func (s *Stack[T]) Pop() (T, bool)
```

## Run Instructions

```bash
go run ./01-foundations/06-types-and-interfaces/17-generic-structures
```

## Code Walkthrough

### Stack implementation

LIFO data structure with type-safe push/pop.

### Queue implementation

FIFO data structure.

### Set implementation

Unique element collection using map.

## Try It

1. Add a Peek method to Stack that returns top element without removing.
2. Implement a generic LinkedList.
3. Add Remove method to Set.

## Production Relevance

Generic data structures are used throughout Go codebases for type-safe collections without runtime overhead.

## Next Step

Section 06 complete! Move to Composition (s06).