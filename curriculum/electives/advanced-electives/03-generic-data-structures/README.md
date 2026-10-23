# Generic data structures

## When you need this

Before Go 1.18, type-safe data structures (stacks, queues, trees, sets) required either code generation per type or interface{} with runtime type assertions. Both approaches add maintenance burden or lose type safety. Generic data structures solve this by allowing a single implementation that works with any type while maintaining compile-time type checking.

## When you do not need this

When simpler concrete types or standard core library primitives suffice without adding abstraction overhead.

## Prerequisites

- elective-02

## Mental model

A generic data structure is a struct with a type parameter. The type parameter propagates to all methods on the struct. When you write Stack[int], the compiler monomorphizes the struct and all its methods for int. At runtime, Stack[int] and Stack[string] are completely distinct types with no shared code — each has its own set of concrete methods operating on the correct concrete type.

## How it works

A generic struct is defined with type parameters: type Stack[T any] struct { items []T }. Each method on the struct uses the same type parameter: func (s *Stack[T]) Push(item T). When instantiated (Stack[int]), the compiler creates a concrete struct with all T replaced by int. This means Push on Stack[int] takes int, and Push on Stack[string] takes string. The monomorphized struct is identical to hand-written code — no boxing, no type assertions, no indirection.

## Go example

```go
package main

import "fmt"

func main() {
    fmt.Println("Elective: Generic data structures")
}
```

## Tradeoffs

- Generic data structures have zero runtime overhead — each instantiation is a concrete type with no boxing or dispatch
- Binary size grows with each distinct (struct type, method) instantiation pair
- Methods on generic structs are inlinable — the compiler can inline monomorphized methods like any concrete method
- Large generic structs with many methods and many instantiations can significantly increase compile time

## Common mistakes

- Using interface{} for data structure elements — loses type safety and requires type assertions on every Pop/Get call.
- Not using pointer receivers for methods that modify the structure — value receivers operate on a copy, and modifications are lost.
- Forgetting that methods on generic types cannot have additional type parameters — use standalone generic functions for operations that need a different type parameter.
- Storing different concrete instantiations in the same variable — Stack[int] and Stack[string] are different types and cannot be assigned to each other.
- Using zero-value initialization incorrectly — var s Stack[T] initializes items to nil, and Push must handle nil slices by allocating.

## Practice task

The learner must implement a generic data structure (Stack[T] or Queue[T]) with Push, Pop, Peek, and Len methods, write tests that verify type safety (Stack[int] rejects string values at compile time), and explain how monomorphization produces concrete code for each instantiation.

## Production notes

Generic data structures are used in Go for type-safe collections: Set[T] for deduplication, Stack[T] for undo/redo, Queue[T] for work queues, Heap[T] for priority queues, and LRU[T] for caches. Production Go code uses generic structures to avoid interface{} boxing in hot paths where type assertions would add overhead and risk panics.

## Related core lessons

- Module 16.02 — Complex generic constraints
- Module 16.15 — Demo preparation
