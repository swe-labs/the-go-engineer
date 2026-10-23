# Generics

## When you need this

New Go developers write separate functions for each type (IntMap, StringMap, FloatMap) or use interface{} with runtime type assertions. Both approaches produce duplicated code or lose compile-time type safety. Generics solve this by allowing functions and data structures to operate on any type while maintaining compile-time type checking.

## When you do not need this

When simpler concrete types or standard core library primitives suffice without adding abstraction overhead.

## Prerequisites

- core-16-15

## Mental model

Generics are compile-time templates. When you write a generic function, the compiler creates a concrete copy of the function for each distinct set of type arguments at each call site. The generic function itself never exists at runtime — only the monomorphized concrete functions do. The type parameter is a placeholder that the compiler replaces with the actual type argument during compilation.

## How it works

Go generics use monomorphization: the compiler generates a separate concrete function for each unique set of type arguments. This means generics have zero runtime cost — the machine code for a generic function is identical to hand-written concrete code. The tradeoff is binary size and compile time: each distinct instantiation adds code to the binary. The compiler also checks that the type argument satisfies the constraint at compile time, preventing any runtime type assertion failures.

## Go example

```go
package main

import "fmt"

func main() {
    fmt.Println("Elective: Generics")
}
```

## Tradeoffs

- Generics are zero-cost at runtime: monomorphization produces the same machine code as hand-written concrete functions — no boxing, no dynamic dispatch
- Binary size increases linearly with the number of distinct instantiations — each unique (T, constraint) pair generates a separate function
- Compile time increases because the compiler must type-check and generate code for each instantiation — generics-heavy packages compile slower
- Inlining works normally with monomorphized generics — the compiler can inline specialized copies just like concrete functions

## Common mistakes

- Using interface{} instead of a type parameter — loses compile-time type safety and requires runtime type assertions that can panic.
- Declaring type parameters that are never used by the function body — adds complexity without benefit.
- Using too many type parameters — each type parameter increases compile time and makes the function signature harder to read.
- Writing constraint interfaces that are too restrictive — prevents valid type arguments from being used, and the function is less reusable than necessary.
- Forgetting that methods cannot have their own type parameters in Go 1.18-1.23 — standalone generic functions are needed instead.

## Practice task

The learner must write a generic function that works with at least two distinct concrete types, write a generic data structure (e.g., a Stack or Queue), and explain how monomorphization works at compile time. They must also demonstrate understanding of constraint interfaces by writing a constraint with type approximation (~T) and a type union.

## Production notes

Generics are used in Go libraries for type-safe data structures (slices, maps, sets, heaps), algorithm functions (sort, map, reduce, filter), and serialization libraries. The standard library's slices.Sort, slices.Compact, and maps.Clone are built on generics. Production Go code uses generics for type-safe repository patterns, event handlers, and middleware that operate on different types without type assertions.

## Related core lessons

- Module 15.15 — Demo preparation
- Module 17.15 — Demo preparation
