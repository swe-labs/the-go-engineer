# Memory preview: stack vs heap

## Learning objective

Understand and apply Memory preview: stack vs heap in the context of professional Go software engineering.

## Why this matters

Not all data has the same lifetime. Stack memory auto-frees when functions return (fast), but only works for fixed-size, function-scoped data. Heap memory persists until explicitly freed or garbage-collected, at the cost of slower allocation and reclamation.

## Mental model

The stack is a LIFO (Last-In-First-Out) structure for function-local data — like a stack of plates, the last one put on is the first taken off. The heap is a pool of memory for data that needs to live beyond its allocating function — like a storage unit where items can be retrieved in any order.

## Core idea

Memory is finite and must be reused. Stack vs. heap is the fundamental tradeoff between allocation speed and lifetime flexibility. Understanding this tradeoff is essential for writing correct, efficient programs in any language.

## Under the hood

Go's stack grows dynamically (starting at 2 KB, growing up to 1 GB on 64-bit systems). Stack frames are copied on growth via stack copying — the GC updates all pointers in the old stack to point to the new location. Go's heap uses a page-based allocator with three tiers: page heap (large allocations), central free lists (medium), and per-P caches (small, fast).

## How Go uses it

Go uses escape analysis at compile time to decide stack vs. heap allocation. If a variable's address does not escape the function, it stays on the stack. If the compiler cannot prove the address is not used after the function returns, it allocates on the heap. Go's garbage collector handles heap reclamation automatically.

## Go example

```go
package main

import "fmt"

// StackValue creates a local variable that stays on the stack.
func StackValue() int {
	x := 42
	return x
}

// HeapPointer creates a value and returns its address (escapes to heap).
func HeapPointer() *int {
	x := 42
	return &x
}

// SumNumbers demonstrates heap allocation for dynamic-size data.
func SumNumbers(n int) int {
	nums := make([]int, n)
	for i := range nums {
		nums[i] = i
	}
	sum := 0
	for _, v := range nums {
		sum += v
	}
	return sum
}

func main() {
	s := StackValue()
	p := HeapPointer()
	sum := SumNumbers(5)
	fmt.Printf("Stack: %d, Heap: %d, Sum: %d\n", s, *p, sum)
}
```

## Step-by-step execution

1. A function is called — the runtime pushes a stack frame onto the goroutine's stack (growing if needed).
2. Local variables are allocated by decrementing the stack pointer — no malloc, no GC tracking.
3. If the variable's address is returned or stored in a globally reachable location, escape analysis marks it for heap allocation.
4. Heap allocation calls Go's memory allocator (tcmalloc-inspired), which finds a free span in the heap arena.
5. When the heap object is no longer reachable, the garbage collector marks and sweeps it in a later GC cycle.

## Common mistakes

| Mistake | Why it happens | Fix |
|---------|---------------|-----|
| Thinking 'memory' is a single uniform space | Stack and heap have different allocation patterns, lifetimes, and performance characteristics | Learn to distinguish: stack (fast, LIFO, function-scoped) vs heap (flexible, GC-managed, any lifetime) |
| Returning a pointer to a stack-allocated variable from a function | In C this causes undefined behavior, but in Go it is safe because escape analysis promotes to heap | Trust Go's escape analysis but verify with `go build -gcflags=-m` |
| Assuming heap allocation is always bad | Small, short-lived allocations on the stack are fast, but heap is necessary for dynamic lifetimes | Profile before optimizing; use stack for local data and heap only when needed |

## Debugging walkthrough

### Scenario: A Go program allocates many small objects in a tight loop, causing high GC pressure and degraded throughput.

**Cause:** Every allocation on the heap adds work for the garbage collector, which must trace and free unreachable objects.

**Fix:** Use object pools (sync.Pool), pre-allocate slices with make, or redesign to reduce allocation frequency.

### Scenario: A recursive function with no base case causes a stack overflow — the stack memory is exhausted.

**Cause:** Each recursive call pushes a stack frame, and infinite recursion fills the limited stack space (1-8 MB per goroutine in Go).

**Fix:** Add a proper base case, or convert the recursion to iteration (or tail recursion where supported).

## Production notes

Performance-critical code in Go (network proxies, databases, game servers) minimizes heap allocations by designing data structures that stay on the stack or use object pools. Understanding stack vs. heap is essential for writing efficient Go programs.

## Performance implications

- Stack allocation takes ~1-3 nanoseconds — just moving the stack pointer. Heap allocation takes ~10-100 nanoseconds plus GC overhead.
- Frequent heap allocation generates garbage that the GC must collect, causing stop-the-world pauses (typically <1ms in Go).
- Stack memory is cache-friendly (spatial locality). Heap memory can fragment and cause cache misses.

## Practice task

The learner must write a Go program with two functions — one that allocates on the stack (local variable) and one that allocates on the heap (returned pointer) — and use 'go build -gcflags=-m' to verify the compiler's escape analysis decisions.

## Tests / verification

```bash
go test ./curriculum/modules/01-computers-terminal-git-web/08-memory-preview-stack-vs-heap/
```

## Review questions

1. Learner must identify whether a given variable is stack-allocated or heap-allocated based on its lifetime and usage pattern.
2. Learner must explain why Go can safely return a pointer to a local variable while C cannot.
3. Learner must predict the relative performance of stack vs. heap allocation for a given code pattern.

## NEXT UP

Lesson 09: [Git mental model](../09-git-mental-model/README.md)
