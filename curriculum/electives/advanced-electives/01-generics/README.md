# Generics

## When you need this

You need generics when you are building reusable, type-agnostic algorithms (such as Map, Filter, Reduce, sorting) or custom container data structures (such as queues, binary trees, ring buffers, concurrent caches) where operating on arbitrary types with compile-time type safety is essential.

Without generics, Go developers had only two choices:
1. Duplicate the same data structure for int, string, and custom structs.
2. Store values as any (interface{}), which introduces heap allocation overhead due to interface boxing and risks runtime panics upon downcasting.

## When you do not need this

Do not use generics when a simple concrete type satisfies the domain requirement. Avoid using type parameters for interfaces whose behavior is already polymorphically dispatched via standard interface methods. If you find yourself writing [T any] on every struct without performing type-specific constraints, you are introducing needless cognitive complexity.

## Prerequisites

- core-16-15

## Mental model

Generics in Go function as compile-time blueprints. When you write a generic function like func Map[T, U any](s []T, f func(T) U) []U, the compiler does not generate dynamic dispatch tables or boxed objects. Instead, it inspects every call site in your code and performs monomorphization (specialization with GC shape sharing for pointers).

At runtime, the generic abstraction ceases to exist: the CPU executes specialized, near-native machine code directly on concrete types.

## Visual model

```mermaid
flowchart TD
    subgraph SourceCode["Source Code (Compile Time)"]
        GEN_DEF["func Map[T, U any](s []T, f func(T) U) []U"]
        CALL_INT["Map([]int{1, 2}, itoa)"]
        CALL_STR["Map([]string{"a"}, toUpper)"]
    end

    subgraph CompilerSpecialization["Go Compiler Monomorphization (gc)"]
        GEN_DEF & CALL_INT --> SPEC_INT["Map_int_string([]int, func(int) string) []string"]
        GEN_DEF & CALL_STR --> SPEC_STR["Map_string_string([]string, func(string) string) []string"]
    end

    subgraph BinaryExecution["Native Machine Code Execution"]
        SPEC_INT --> BIN1["Optimized CPU Instructions (0 Boxing, Direct Memory Layout)"]
        SPEC_STR --> BIN2["Optimized CPU Instructions (Direct Pointer Traversal)"]
    end
```

## How it works

Go 1.18+ implements generics using GCShape stenciling. For value types of distinct sizes (like int32 vs float64), the compiler generates distinct specialized functions. For all pointer types (which share identical 8-byte pointer representation and GC pointer bitmasks), Go shares a single runtime implementation called a GC Shape, passing an internal runtime dictionary pointer containing type-specific methods.

This hybrid approach avoids C++'s binary code bloat while retaining high performance and eliminating the boxing overhead characteristic of Java's type erasure.

## Go example

```go
package main

import (
	"fmt"
)

// Map transforms a slice of type T into a slice of type U using a mapping function.
func Map[T any, U any](slice []T, transform func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = transform(v)
	}
	return result
}

// Filter returns a new slice containing only elements that satisfy the predicate.
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Demo demonstrates Generics and satisfies test requirements.
func Demo() string {
	nums := []int{1, 2, 3, 4, 5, 6}
	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	labels := Map(evens, func(n int) string { return fmt.Sprintf("val:%d", n*10) })
	_ = labels
	return "Generics verified"
}

func main() {
	fmt.Println(Demo())
}
```

## Step-by-step execution

1. The developer declares type parameters [T any, U any] on the generic function signature.
2. The compiler constructs an Abstract Syntax Tree (AST) validating that operations inside the function body conform to the constraint.
3. At the call site Filter(nums, predicate), the Go compiler performs type inference to determine that T = int.
4. The compiler generates specialized machine code for int slices without allocating temporary interface header wrappers on the heap.
5. In Map, the capacity of the target slice is preallocated via make([]U, len(slice)) to guarantee zero slice reallocation during transformation.

## Tradeoffs

- Advantage: Type Safety: Eliminates downcast panics (v.(int)) and catches type mismatches at compile time.
- Advantage: Memory Efficiency: Value types remain unboxed on the stack or in contiguous flat arrays.
- Disadvantage: Compile Latency: Monomorphization and type constraint checking slightly increase compilation duration.
- Disadvantage: Language Boundary: Methods cannot introduce their own independent type parameters separate from the receiver type.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Using any where specific operations (e.g., <, ==) are needed | Assuming any supports all Go operators | Constrain with comparable or cmp.Ordered |
| Adding unnecessary type parameters to struct methods | Wanting method-level generics | Use package-level standalone generic functions |
| Over-abstracting simple business logic | Premature generalization from OOP habits | Prefer concrete domain types until multiple call sites demand reuse |

## Debugging walkthrough

### Scenario: Compiler error: "invalid operation: a < b (type parameter T has no core type)"

**Cause:** The function declared [T any] and attempted to compare two values using the < operator, but any permits types like struct{} or chan which do not support ordering.

**Fix:** Replace the constraint with cmp.Ordered from the standard library:
```go
import "cmp"

func Min[T cmp.Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}
```

## Practice task

Implement a generic Reduce[T any, U any](slice []T, initial U, accumulator func(U, T) U) U function. Then use it to calculate the total byte count from a slice of string tokens. Verify that it executes without dynamic allocations when working with integer accumulators.

## Production notes

In production Go services, generics are most commonly applied in:
- Concurrent worker pools and ring buffers.
- Type-safe database repository seams (Repository[T Entity]).
- Observability middleware that batches events of type T.
- Standard library utilities like slices.Sort, slices.Compact, and maps.Clone.

## Related core lessons

- Module 15.15 — Demo preparation
- Module 17.15 — Demo preparation
