# Generics Fundamentals

## Learning objective

Understand and apply Go generics (type parameters) introduced in Go 1.18+, including generic function signatures, standard constraints (`any`, `comparable`, and `cmp.Ordered`), generic types (containers/data structures), compile-time monomorphization and GC shape stenciling, and the engineering decision framework for choosing between generics and interfaces.

## Why this matters

Prior to Go 1.18, writing type-agnostic code required either duplicating logic across every concrete type, relying on `interface{}` with runtime type assertions (which bypasses compile-time type safety and incurs heap boxing allocations), or using code generation tools. Generics solve this by allowing functions and data structures to accept type parameters checked at compile time. However, generics can easily be abused: using generics where standard interfaces suffice creates needlessly complex APIs. Understanding when to reach for generics and when to stick to interfaces is a crucial software design skill.

## Mental model

A generic function or type is a **type-safe blueprint**. When you supply concrete type arguments (or let Go infer them), the compiler verifies that the concrete types fulfill the specified constraint.

```
Type Parameter: [T comparable]
       │
       ├─ If T = int      ──► Valid (integers support == and !=)
       ├─ If T = string   ──► Valid (strings support == and !=)
       └─ If T = []byte   ──► Compile Error! (slices cannot be compared with ==)
```

## Visual model

```mermaid
flowchart TD
    subgraph GenericDeclaration["Generic Function Declaration: [T constraints.Ordered]"]
        SIG["func Min[T constraints.Ordered](a, b T) T"]
    end

    subgraph Instantiations["Compile-Time Monomorphization"]
        C1["Min[int](10, 20) --> Generates concrete int machine instructions"]
        C2["Min[string]("a", "b") --> Generates concrete string machine instructions"]
    end

    GenericDeclaration --> C1
    GenericDeclaration --> C2
```

## Core idea

### 1. Type Parameters & Constraints
- `any`: An alias for `interface{}`, permitting any type.
- `comparable`: Built-in constraint matching any type that supports the `==` and `!=` operators (primitives, pointers, structs of comparable fields). Excludes slices, maps, and functions.
- `cmp.Ordered`: Defined in package `cmp` (Go 1.21+), matching types that support ordering operators (`<`, `<=`, `>`, `>=`): all integer, float, and string types.

```go
func Min[T cmp.Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}
```

### 2. Generics vs. Interfaces: The Decision Framework

| Requirement | Use Generics `[T any]` | Use Interfaces (`io.Reader`) |
| :--- | :---: | :---: |
| Type-safe data structures (sets, stacks, queues, trees) | **Yes** | No (requires boxing/assertions) |
| Slice/map utility algorithms (`slices.Sort`, `Map`, `Filter`) | **Yes** | No |
| Behavioral abstraction (calling `.Read()` or `.Close()`) | No | **Yes** (idiomatic dynamic dispatch) |
| Decoupling components and mocking in unit tests | No | **Yes** |
| Heterogeneous collections (slice holding different types) | No (all items must match `T`) | **Yes** (`[]Speaker`) |

## Under the hood

- **Monomorphization & GC Shape Stenciling:** Unlike Java (which uses type erasure and boxes primitives into objects), Go generates specialized machine code. To prevent binary size explosion from infinite instantiations, Go groups types with identical runtime representations (all pointer types share a single "GC shape") into shared instantiations, while generating specialized code for distinct primitive shapes (`int`, `float64`).
- **Zero Runtime Boxing:** Operating on generic values does not allocate interface fat pointers or trigger runtime type assertions.

## How Go uses it

- Packages `slices`, `maps`, and `cmp` in the standard library are entirely built on generics.
- Modern ORMs and repository patterns use generics (`Repository[T]`) to provide type-safe CRUD operations without repetitive boilerplate.

## Go example

```go
package main

import (
	"cmp"
	"fmt"
)

// Clamp restricts a value to the range [low, high].
func Clamp[T cmp.Ordered](val, low, high T) T {
	if val < low {
		return low
	}
	if val > high {
		return high
	}
	return val
}

func main() {
	fmt.Println(Clamp(15, 0, 10))      // 10 (int)
	fmt.Println(Clamp(3.14, 0.0, 5.0))  // 3.14 (float64)
	fmt.Println(Clamp("zoo", "a", "m")) // "m" (string)
}
```

## Failure modes & common mistakes

1. **Over-engineering with Type Parameters:** Introducing `[T any]` when an `io.Reader` or `error` interface would be cleaner and more readable. If your function only calls methods on `T`, use an interface instead.
2. **Method Type Parameters Misconception:** In Go, methods cannot have their own type parameters separate from their receiver struct. You cannot write `func (s MyStruct) Process[T any](v T)`. The type parameter must be declared on the struct itself: `type MyStruct[T any] struct`.
3. **Attempting Equality on `any`:** `func Equal[T any](a, b T) bool { return a == b }` fails to compile because `any` does not guarantee comparability. Use `[T comparable]`.

## Practice task

1. Open `main.go` and inspect the implementations of `Map`, `Filter`, `Contains`, and `Stack[T]`.
2. Add a generic `Set[T comparable]` struct with `Add(item T)`, `Remove(item T)`, and `Contains(item T) bool` methods backed by a `map[T]struct{}`.
3. Run the test suite:
   ```bash
   go test -v ./curriculum/modules/05-types-interfaces-packages-modules/lessons/24-generics-fundamentals
   ```

## Review questions

1. What is the difference between the `any` constraint and the `comparable` constraint? Why can't `[]byte` satisfy `comparable`?
2. Why does Go use GC shape stenciling instead of pure monomorphization for every single pointer type?
3. In what scenario would you choose an `interface` over a generic type parameter?
4. Can a method on a non-generic struct declare its own type parameters? How do you work around this limitation?
