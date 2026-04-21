# TI.15 Generic Data Structures

## Mission

Learn to build type-safe generic data structures like `Stack`, `Queue`, and `Set` using Go's generics.

## Why This Lesson Exists Now

You've learned generic functions. Now learn to build generic data structures that are type-safe at compile time - no runtime type assertions needed.

## Prerequisites

- `TI.9` generics
- `TI.14` complex constraints

## Mental Model

Think of a reusable storage box. Without generics, you'd need separate boxes for books, clothes, and electronics. With generics, one `Box[T]` works for all - type-safe and efficient.

## Visual Model

```mermaid
graph TD
    A["data"] --> B["type definition"]
    B --> C["methods or interface behavior"]
```

```go
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T)
func (s *Stack[T]) Pop() (T, bool)
```

## Machine View

These generic structures are still built from the same tools you already know: slices, maps, methods, and type parameters. Generics remove type duplication, but the underlying storage behavior is still slice growth and map lookup.

## Run Instructions

```bash
go run ./04-types-design/15-generic-data-structures
```

## Code Walkthrough

### Stack implementation

LIFO data structure with type-safe push and pop behavior.

### Queue implementation

FIFO data structure.

### Set implementation

Unique element collection using a map.

## Try It

1. Add a `Peek` method to `Stack` that returns the top element without removing it.
2. Implement a generic linked list.
3. Add a `Remove` method to `Set`.

## In Production

Generic data structures are used throughout Go codebases for type-safe collections without runtime overhead.

## Thinking Questions

1. What problem is this lesson trying to solve?
2. What would change if you removed this idea from the program?
3. Where do you expect to see this pattern again in real Go code?

## Next Step

The optional stretch path is complete. Move to **Composition** next, or return to the Section 04 map whenever you want to review the core path again.
