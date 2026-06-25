# Pointers as addresses

## Learning objective

Use Go's `&` and `*` operators to create, dereference, and pass pointers, and explain the differences between Go pointers and C pointers.

## Why this matters

Pointers enable reference semantics: modifying a variable through a function call, sharing large structs without copying, and building linked data structures. Go uses pointers sparingly compared to C, but every Go engineer must understand them because they appear in method receivers, interface values, concurrency, and performance optimization. Misunderstanding pointers leads to nil dereference panics, logic bugs, and data races.

## Mental model

A pointer is a memory address. It tells you where a value lives, not what it is. The `&` operator asks "where are you?" The `*` operator says "what's at this address?".

Think of it like a home address on a piece of paper. You can give someone the paper (pass a pointer), they can visit the house and change the furniture (modify the pointed-to value), and when you go home, you see the changes. Without a pointer, you photocopy the house's contents and hand over the photocopy — their changes don't affect the original.

## Core idea

- `*T` is the type "pointer to T". Zero value is `nil`.
- `&x` yields a pointer to `x` (address of `x`).
- `*p` dereferences pointer `p`, giving access to the value at that address.
- `p.Field` is shorthand for `(*p).Field` — Go automatically dereferences.
- No pointer arithmetic: `p++` is a compile error. Go pointers are safe.
- `new(T)` allocates a zero-valued `T` and returns `*T`.

## Under the hood

A pointer is an unsigned integer holding a virtual memory address. On 64-bit systems, it is 8 bytes. On 32-bit, 4 bytes.

When you write `p := &x`, the compiler emits a `LEA` (Load Effective Address) instruction. When you write `*p = 5`, it emits a `MOV` to the address stored in `p`.

Go's pointer model is a subset of C's: you can take the address of any addressable value (variable, array element, struct field), but you cannot:
- Perform arithmetic on pointers.
- Cast between unrelated pointer types.
- Take the address of a constant or a literal (except composite literals with `&`).

The escape analyzer decides whether a pointer must live on the heap. If a pointer to a local variable escapes (e.g., returned from a function), the variable is heap-allocated instead of stack-allocated.

## How Go uses it

- **Method receivers**: Value receiver `(s Struct)` copies the struct; pointer receiver `(s *Struct)` shares it.
- **Modifying function arguments**: `func update(p *int) { *p = 42 }`.
- **Optional values / nilability**: `*T` can be `nil`, signaling absence. Used for optional configuration, nullable fields.
- **Sharing large data**: Passing `*Struct` avoids copying megabytes of data.
- **Interface values**: An interface value stores a pointer to the concrete value when the concrete type is not a pointer.
- **Linked structures**: Tree nodes, linked lists use pointers: `type Node struct { next *Node }`.

## Go example

```go
package main

import "fmt"

func main() {
	x := 10
	p := &x

	fmt.Println("x:", x)
	fmt.Println("p:", p)   // memory address
	fmt.Println("*p:", *p) // value at address

	*p = 20
	fmt.Println("x after *p = 20:", x)

	// Pointer to struct
	type Point struct{ X, Y int }
	pt := Point{3, 4}
	q := &pt
	q.X = 5 // equivalent to (*q).X = 5
	fmt.Println("pt after q.X = 5:", pt)

	// new function
	n := new(int)
	*n = 99
	fmt.Println("*n:", *n)
}
```

## Step-by-step execution

For `x := 10; p := &x; *p = 20`:

1. `x := 10` allocates a local variable `x` on the stack with value `10`.
2. `p := &x` reads the address of `x`'s stack slot and stores it in `p`.
3. `*p = 20` dereferences `p`, writing `20` to the address pointed to by `p` — which is `x`'s stack slot.
4. Printing `x` shows `20`.

For `q := &Point{3, 4}; q.X = 5`:

1. `&Point{3, 4}` allocates a `Point` value (either on stack or heap depending on escape analysis) and returns its address.
2. `q` holds the address.
3. `q.X` is syntactic sugar for `(*q).X`. The compiler inserts the dereference automatically.
4. `= 5` writes to the `X` field at the pointed-to address.

## Common mistakes

- Mistake: Dereferencing a nil pointer.
  - Why it happens: `var p *int` is nil. `*p = 5` panics with "nil pointer dereference".
  - Fix: Always check `if p != nil` before dereferencing.

- Mistake: Confusing `*T` with `T` as function parameter.
  - Why it happens: `func f(v T)` copies; `func f(v *T)` shares. Passing without `&` copies the value, so modifications inside the function are invisible to the caller.
  - Fix: Match the signature. If the function modifies, use pointer. If not, consider value.

- Mistake: Using `new(T)` vs `&T{}`.
  - Why it happens: Both return `*T`. `new(T)` zeroes all fields; `&T{}` uses composite literal syntax with field values. They are equivalent for zero values.
  - Fix: Prefer `&T{}` when you know the initial values; `new(T)` is rarely used in idiomatic Go.

- Mistake: Pointers to pointers (`**int`) when not needed.
  - Why it happens: Ported from C where double pointers are common for modifying pointer arguments. Go can return the new pointer instead.
  - Fix: Use `return *T` instead of `**T` parameter. Go is not C.

## Debugging walkthrough

```go
package main

import "fmt"

func double(x int) {
	x = x * 2
}

func main() {
	n := 21
	double(n)
	fmt.Println(n)
}
```

**Symptom**: Prints `21`, not `42`.

**Investigation**: The function `double` receives a copy of `n`. Inside `double`, only the copy is modified.

**Root cause**: Go passes by value. `x` is a copy of `n`. `x = x * 2` modifies the copy.

**Fix**: Use a pointer:

```go
func double(x *int) {
	*x = *x * 2
}
// call: double(&n)
```

## Production notes

- **`*T` in public APIs**: Only use pointer parameters when the function needs to modify the value or when `T` is very large. Otherwise, use value types — they are safer and the compiler may optimize them.
- **Nil checks are cheap but mandatory**: A single comparison instruction. Missing one causes a panic.
- **Pointer equality**: Two pointers are equal if they point to the same address. This is useful for identity comparison (e.g., two values at the same memory location).
- **`unsafe.Pointer`**: Escapes Go's type system. Only use when interfacing with C via `cgo` or for certain low-level optimizations. Never use `unsafe` casually.

## Performance implications

- Passing `*T` (8 bytes) instead of `T` (e.g., 100 bytes) is cheaper in call overhead.
- Indirection costs: `p.Field` requires loading the pointer (memory read) before accessing the field. For hot loops, this can be slower than direct access.
- Escape to heap: If a pointer to a local variable escapes, the variable is heap-allocated, which adds GC pressure. Profile before optimizing.
- `new(T)` always allocates on the heap for most types (small types may be optimized away).

## Practice task

Write a function `swap(a, b *int)` that swaps the values at two pointers. Then write a function `apply(nums []int, fn func(*int))` that applies a function to each element of a slice via a pointer — the function should be able to modify the element. In `main()`, create a slice of `[3]int{1, 2, 3}`, swap the first and last via pointers, then use `apply` to double each element.

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/22-pointers-as-addresses
go test ./curriculum/modules/03-programming-fundamentals/lessons/22-pointers-as-addresses
```

## Review questions

1. What is the zero value of a pointer type?
2. What does `&x` produce? What type does it have?
3. Can you write `p++` on a pointer `p` in Go? Why or why not?
4. What is the difference between `new(int)` and `&int{}` in Go?
5. When would you use a pointer receiver instead of a value receiver?

## NEXT UP

Beginner pointer mistakes — the most common pitfalls when learning Go pointers and how to avoid them.
