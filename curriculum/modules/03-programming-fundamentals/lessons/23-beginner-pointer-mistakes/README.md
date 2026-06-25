# Beginner pointer mistakes

## Learning objective

Identify and avoid the five most common pointer-related bugs in Go: nil dereference, returning pointer to local variable, pointer vs value receiver confusion, range variable aliasing with pointers, and dangling pointers via append.

## Why this matters

Pointer bugs are the most common source of panics and data corruption in beginner Go code. A single nil dereference takes down a production server. A value receiver used where a pointer was needed silently operates on a copy. A misused `&` in a range loop captures the same address every iteration. Learning these patterns prevents whole categories of bugs before they happen.

## Mental model

Every pointer-related bug falls into one of two categories: **pointing nowhere** (nil) or **pointing to the wrong thing** (aliasing, stale reference, copy confusion). The fix is always to ensure the pointer targets the correct, live memory with the correct semantics.

Go's pointer model is simpler than C's, but the common mistakes are more insidious because Go's automatic dereferencing and escape analysis can hide what is happening.

## Core idea

Five classic mistakes:

1. **Nil pointer dereference**: `var p *int; *p = 5` → panic.
2. **Returning pointer to local variable**: The variable escapes to heap (actually fine in Go, but the surprise is that it works — the real mistake is expecting _not_ to escape, leading to unintended heap allocation).
3. **Pointer vs value receiver confusion**: Value receiver on a method that needs to modify the struct.
4. **Range variable aliasing**: `for _, v := range slice { p := &v }` — all `p`s point to the same `v`.
5. **Dangling pointers via append**: Calling `append` on a shared slice, then using the original pointer to a stale backing array.

## Under the hood

**Nil pointer**: A pointer is a memory address stored as an unsigned integer. `nil` is address `0`. The MMU (Memory Management Unit) prevents any read or write to address 0. The Go runtime checks for nil before dereferencing in many cases, but not all — the hardware ultimately catches it and delivers a SIGSEGV, which Go translates to a panic.

**Escape analysis**: The compiler determines whether a variable can be stack-allocated. If a pointer to a local variable is returned or stored in a heap-allocated structure, the variable "escapes" to the heap. On Go playground, returning `&local` is safe — the compiler moves it to the heap.

**Range variable reuse**: In Go versions before 1.22, the loop variable `v` in `for _, v := range slice` is a single variable reused across iterations. Taking `&v` gives the same address every time. Go 1.22 changed this: each iteration gets its own `v`. But the old pattern of `v := v` inside the loop is still idiomatic for backward compatibility.

## How Go uses it

Go's design intentionally avoids the most dangerous C pointer patterns (pointer arithmetic, dangling pointers to stack frames). But Go's automatic escape analysis and garbage collection create their own confusion — beginners expect C-like rules and are surprised when Go handles things differently, or vice versa.

The standard library uses pointers extensively:
- `*http.Request`, `*http.Response`, `*os.File` — all are pointer types.
- `*sql.Rows`, `*sql.Stmt` — database handles are pointers.
- `*big.Int`, `*big.Float` — big numbers require pointer semantics.
- Method sets: `*T` implements all methods of `T` plus the pointer receiver methods.

## Go example

```go
package main

import "fmt"

func main() {
	// 1. Nil dereference
	var p *int
	if p != nil {
		*p = 5
	}
	// Without the nil check, *p = 5 panics.

	// 2. Returning pointer to local variable (fine in Go)
	pt := newPoint(3, 4)
	fmt.Println("point from function:", pt)

	// 3. Value vs pointer receiver
	d := &Dog{name: "Rex"}
	d.Bark() // works: *Dog implements Speaker
	// Dog{name: "Rex"}.Bark() // compile error: Dog does not implement Bark (value receiver)

	// 4. Range variable aliasing (pre-Go 1.22)
	items := []int{1, 2, 3}
	ptrs := []*int{}
	for _, v := range items {
		v := v // create new variable per iteration
		ptrs = append(ptrs, &v)
	}
	for i, p := range ptrs {
		fmt.Printf("ptrs[%d] = %d (addr=%p)\n", i, *p, p)
	}
}

type Dog struct {
	name string
}

func (d *Dog) Bark() {
	fmt.Println(d.name, "says woof!")
}

func newPoint(x, y int) *struct{ X, Y int } {
	return &struct{ X, Y int }{x, y}
}
```

## Step-by-step execution

For the nil dereference `var p *int; *p = 5`:

1. `var p *int` creates a pointer variable on the stack initialized to `nil` (address `0`).
2. `*p = 5` attempts to write `5` to address `0`.
3. The hardware raises a page fault. Go's runtime catches it and panics: "runtime error: invalid memory address or nil pointer dereference".

For the range variable bug `for _, v := range slice { ptrs = append(ptrs, &v) }` in Go <1.22:

1. Loop iteration 1: `v = 10`. `&v` points to the loop variable's stack slot. Append to `ptrs`.
2. Loop iteration 2: `v = 20`. Same `v`, same address. `&v` overwrites the previous value at that address.
3. All pointers in `ptrs` point to `v`, which holds `20` (the last value).

**Fix** (pre-1.22): `v := v` inside the loop shadows the outer `v` with a new variable per iteration.

## Common mistakes

All five are covered in the Core idea section. Here are additional tips for each:

- **Nil dereference**: Always check error returns and nil pointers from functions. Use `if p != nil` guards.
- **Returning pointer to local variable**: This is actually safe in Go (escape analysis moves it to heap). The mistake is expecting stack allocation. If you want to avoid heap allocation, don't take the address.
- **Pointer vs value receiver**: A value method is part of the method set of `T` and `*T`. A pointer method is only part of `*T`. If you call a pointer method on a non-addressable value, the compiler rejects it.
- **Range variable**: In Go 1.22+, each iteration gets its own variable, but `v := v` is still safe and recommended for maximum compatibility.
- **Append and stale pointers**: After `append`, the backing array may change. Any pointer into the old backing array becomes stale. Only use `&slice[i]` if you guarantee no `append` happens afterward.

## Debugging walkthrough

```go
package main

import "fmt"

type Counter struct {
	value int
}

func (c Counter) Increment() {
	c.value++
}

func main() {
	c := Counter{value: 0}
	c.Increment()
	fmt.Println(c.value)
}
```

**Symptom**: Prints `0` instead of `1`.

**Investigation**: Add a print inside `Increment`:

```go
func (c Counter) Increment() {
	c.value++
	fmt.Println("inside Increment:", c.value)
}
```

Output: `inside Increment: 1` then `0`.

**Root cause**: `Increment` has a value receiver. `c` is copied. The method modifies the copy. The original is unchanged.

**Fix**: Use a pointer receiver:

```go
func (c *Counter) Increment() {
	c.value++
}
```

And call with `c.Increment()` — Go automatically takes the address of `c` for pointer receiver calls on addressable values.

## Production notes

- **Run `go vet`** in CI. It catches many pointer mistakes, including range variable capture.
- **Use `*T` for method receivers** when the method modifies the receiver or when `T` contains a mutex (sync.Mutex must not be copied).
- **Never store `&slice[i]` after `append`**: The pointer may become invalid. Recompute after every `append`.
- **Interface nil vs pointer nil**: A `*int` variable that is nil does not make the interface holding it nil. `var p *int; var i interface{} = p; i == nil` is `false`. This is a common source of bugs in error handling.
- **Go 1.22+ fixed the range variable**: Update your go.mod `go` directive to 1.22 and remove `v := v` shims, but keep them if you support older toolchains.

## Performance implications

- Nil checks are cheap (one comparison) but missing one causes a panic. Profile: add nil checks everywhere pointers enter your code.
- Pointer receivers avoid copying the struct (O(n) in struct size). For small structs (<= 4 words), the copy cost is negligible.
- Escape to heap via returning pointers or storing pointers in interfaces costs a GC cycle. For latency-sensitive code, minimize heap-allocated pointers.
- Range variable capture (Go <1.22) doesn't just cause bugs — the shared pointer means the compiler cannot prove non-aliasing, which inhibits optimization.

## Practice task

Write a function `makeIncrementer() func() int` that returns a closure. The closure should increment and return a counter each time it's called. Do NOT use a pointer for the counter — use a value that the closure captures.

Then fix this buggy code and explain the fix:

```go
func main() {
	items := []int{10, 20, 30}
	var ptrs []*int
	for _, v := range items {
		ptrs = append(ptrs, &v)
	}
	for _, p := range ptrs {
		fmt.Println(*p)
	}
}
```

In `main()`, demonstrate both the incorrect output (pre-fix) and the correct output (post-fix).

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/23-beginner-pointer-mistakes
go test ./curriculum/modules/03-programming-fundamentals/lessons/23-beginner-pointer-mistakes
```

## Review questions

1. What happens when you dereference a nil pointer in Go?
2. Why does `v := v` inside a `for range` loop fix the aliasing bug in Go versions before 1.22?
3. What is the difference between a value receiver and a pointer receiver method?
4. After calling `append` on a slice, why might an existing `&slice[i]` pointer become invalid?
5. True or false: returning `&localVar` from a function is always a bug in Go. Explain.

## NEXT UP

Congratulations on completing Module 03: Programming Fundamentals! You now have a solid foundation in Go's type system, control flow, data structures, and memory model. Next up: Module 04 — Functions, Errors, and Data Semantics.
