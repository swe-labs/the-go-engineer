# Escape analysis

## Learning objective

Explain how the Go compiler decides whether to allocate a value on the stack or heap, use `go build -gcflags=-m` to inspect escape analysis decisions, and write code that minimizes heap allocations by keeping values on the stack.

## Why this matters

Heap allocation is expensive: each allocation requires GC tracking, and each GC cycle pauses the world (though Go's concurrent GC minimizes this). Stack allocation is effectively free: the value is created and destroyed in a single instruction (adjusting the stack pointer). The difference between a heap-allocated and stack-allocated value can be 10-100x in allocation cost. Understanding escape analysis lets you write allocation-efficient code without resorting to memory pools or unsafe hacks.

## Mental model

Escape analysis is the compiler's answer to: "Does this value's address outlive the function that creates it?" If yes → heap (GC must track it). If no → stack (cheap allocation, freed when function returns). The compiler traces every pointer from creation to final use: if the pointer is returned, stored in a global, put into a heap-allocated struct, or captured by a closure — it escapes. If the pointer is only used locally or passed to functions that do not store it — it stays on the stack.

## Core idea

The Go compiler runs escape analysis on the SSA (Static Single Assignment) intermediate representation after inlining. The analysis builds a directed graph of value flows: allocations, stores, loads, and function calls. For each allocation (call to `runtime.newobject` or stack allocation), the analysis walks forward through all uses of the allocated value. If any use stores the address in a location that outlives the allocating function, the allocation is marked "heap". If all uses are within the function's lifetime, the allocation stays on the stack.

Key triggers for heap escape:

| Pattern | Why it escapes |
|---|---|
| Returning `&T{}` from a function | The address is visible after the function returns |
| Storing a pointer in a global variable | The global outlives the function |
| Passing a pointer to a function that stores it | The callee may hold the pointer after the caller returns |
| Capturing a variable by reference in a closure | The closure may outlive the enclosing function |
| Assigning a pointer to a heap-allocated struct field | The pointer follows the struct's lifetime |
| Converting a value to `interface{}` | Interface values are stored on the heap (often) |

## Under the hood

Escape analysis runs on the compiler's SSA form in `cmd/compile/internal/escape`. It starts by building an escape "forest": a set of nested locations (`Heap`, `Stack`, function parameter locations). Each allocation site is assigned a location. The analysis then performs a dataflow analysis: for each value, it tracks all operations that use the value. If a value is stored into a location that outlives its own location, the value's location is marked as escaping.

The compiler prints its reasoning with `-gcflags=-m` (one level) or `-gcflags=-m -m` (detailed). The output shows lines like:

```
./main.go:10:6: moved to heap: x
./main.go:15:16: &User{...} escapes to heap
```

The `-m -m` output includes the full escape reason chain: "does the address of x flow to a return parameter?" or "stored into heap-allocated struct field y".

## How Go uses it

- The Go compiler applies escape analysis to every function during compilation. There is no user configuration — it is always on.
- `sync.Pool` is the standard mechanism to work around unavoidable escapes: pre-allocate objects and reuse them across goroutines.
- The standard library is carefully written to minimize escapes. For example, `strings.Builder.Grow` pre-allocates the underlying buffer to prevent reallocation escapes.
- Interface method calls often force heap allocation because the interface value stores a pointer to the concrete type.
- In Go 1.22+, the compiler improved escape analysis for range variables, preventing common escape patterns in `for i, v := range slice` loops.

## Go example

```go
package main

import "fmt"

type User struct {
	ID   int
	Name string
}

// stackAllocated creates a User that stays on the stack.
func stackAllocated(id int, name string) User {
	return User{ID: id, Name: name}
}

// heapAllocated creates a User that escapes to the heap.
func heapAllocated(id int, name string) *User {
	return &User{ID: id, Name: name}
}

func main() {
	a := stackAllocated(1, "Alice")
	b := heapAllocated(2, "Bob")

	fmt.Println(a, b)
}
```

Check escape analysis:

```bash
go build -gcflags='-m' ./curriculum/modules/13-performance-memory-engineering/lessons/05-escape-analysis
```

Expected output:

```
./main.go:15:6: can inline stackAllocated
./main.go:20:6: can inline heapAllocated
./main.go:26:16: inlining call to stackAllocated
./main.go:27:16: inlining call to heapAllocated
./main.go:27:16: &User{...} escapes to heap
```

The `heapAllocated` call's `&User{...}` escapes because it is returned as a pointer. The `stackAllocated` call's `User` stays on the stack (after inlining).

## Step-by-step execution

1. The compiler parses `main.go` into an AST and performs type-checking.
2. It builds SSA form for each function. For `heapAllocated`, the SSA shows a `new` instruction (the `&User{}` allocation) and a `return` that uses its address.
3. Escape analysis traces the `new` value: it is stored in the return value, which is used by the caller (`main`). Since the caller assigns the return value to a variable (`b`) that is used after the call, the allocation escapes to heap.
4. For `stackAllocated`, the SSA creates a `User` struct value (not via `new`). The return value is a copy of the struct. The analysis finds no pointer escape — all uses of the struct are within the allocation frame.
5. `main` calls `fmt.Println(a, b)`. The interface parameter of `fmt.Println` causes both `a` and `b` to escape because interface values store pointers to the concrete type. This is visible in the `-m` output.

## Common mistakes

- **Optimizing allocation by guessing** — Changing `&T{}` to `T{}` based on superstition rather than running `-gcflags=-m` or profiling. Always verify with the compiler.
- **Assuming all function parameters are heap-allocated** — Small structs passed by value live on the stack. Only values that escape (returned, stored in interface, captured by closure) go to the heap.
- **Returning pointers to local variables unknowingly** — The compiler detects this and allocates on heap, but the developer may think it is a stack allocation because the code "looks local."
- **Using pointers everywhere to avoid copying** — Pointers cause escape to heap more often than values. Copying small structs (up to a few words) is cheaper than heap allocation plus GC.
- **Believing `make` always allocates on heap** — `make` for small slices created in a function and returned may allocate. But slices that never escape stay on the stack if their backing array is small enough (below ~64 KB).

## Debugging walkthrough

A service has high GC CPU usage. Run `go test -bench=. -benchmem` and notice `allocs/op` is unexpectedly high for a hot function:

```bash
go build -gcflags='-m -m' ./path/to/package 2>&1 | grep "escapes to heap"
```

Output:

```
./hotpath.go:42:18: &requestPool{...} escapes to heap
./hotpath.go:55:10: interface{} literal escapes to heap
```

Line 42: a `&requestPool{}` allocation escapes because the function returns `interface{}`. The fix: change the return type to the concrete type, or use `sync.Pool` for the allocation.

Line 55: a value passed to `fmt.Sprintf` is allocated on heap because `fmt.Sprintf` accepts `interface{}`. The fix: use `strconv.Itoa` instead of `fmt.Sprintf("%d", n)`.

After the fix, re-run `-gcflags=-m` to confirm the escapes are eliminated, then benchmark to measure the improvement.

## Production notes

- Escape analysis is critical for high-throughput Go services. Hot paths should have zero heap allocations per operation if possible.
- The standard library uses `sync.Pool` for allocations that cannot escape (e.g., `fmt` uses a pool for print buffers).
- Interface parameters always force heap allocation of the concrete value because the interface header stores a pointer. In hot paths, use concrete types.
- Closure variables that are modified (not just read) escape. If a closure captures a variable and modifies it, that variable is heap-allocated.
- Goroutine stack sizes start at 2 KB and grow as needed. Stack-allocated values do not cause GC pressure. Heap-allocated values do.

## Performance implications

The cost difference between stack and heap allocation is dramatic. Stack allocation is a single instruction: `SUB $N, SP`. Heap allocation calls `runtime.mallocgc`, which checks the GC state, finds a free span, potentially triggers a GC, and updates allocation counters. Stack allocation is typically 5-10 ns; heap allocation is 50-500 ns. Beyond the allocation itself, heap-allocated values must be scanned by the GC. A heap-allocated `User` struct with two fields adds GC scanning work; the same struct on the stack is invisible to the GC.

However, escape analysis is a compiler optimization, not a correctness concern. Never make code harder to read or maintain solely to avoid an escape. Profile first, then optimize the top allocation sites.

## Practice task

Write a Go program with a `Point` struct (two `float64` fields). Implement:

1. `newPoint(x, y float64) *Point` — returns a pointer (likely escapes).
2. `makePoint(x, y float64) Point` — returns a value (stays on stack after inlining).

Use `go build -gcflags='-m -m'` to verify the escape decisions. Then write benchmarks for both functions with `-benchmem` and compare `allocs/op`.

## Tests / verification

```bash
go run ./curriculum/modules/13-performance-memory-engineering/lessons/05-escape-analysis
go test ./curriculum/modules/13-performance-memory-engineering/lessons/05-escape-analysis
```

To see escape analysis decisions:

```bash
go build -gcflags='-m' ./curriculum/modules/13-performance-memory-engineering/lessons/05-escape-analysis
```

## Review questions

1. What does it mean for a value to "escape to the heap"?
2. List three code patterns that cause a value to escape.
3. How do you instruct the Go compiler to print escape analysis decisions?
4. Why do interface parameters often force heap allocation?
5. Is it always better to avoid heap allocation? When would you accept heap allocation?

## NEXT UP

Memory layout — understanding struct padding and alignment, using `unsafe.Sizeof` to measure struct sizes, reordering fields to minimize padding, and avoiding false sharing in concurrent code.
