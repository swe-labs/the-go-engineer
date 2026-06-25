# Zero values

## Learning objective

Predict the zero value of any Go type and explain how Go's zero-value guarantee prevents undefined-behavior bugs common in other languages.

## Why this matters

In C and C++, an uninitialized variable contains whatever garbage was left in memory. This causes Heisenbugs — intermittent failures that change behavior when you add a print statement or run under a debugger. Go eliminates this entire class of bugs by guaranteeing that every declared variable starts at a well-defined zero. In production, this means no mysterious crashes from uninitialized pointers, no garbage string lengths, and no indeterminate boolean branches. The zero-value design is a safety feature that every Go engineer relies on, often without thinking about it.

## Mental model

Think of each variable declaration as Go writing a default value into the box before you get to use it. The default depends on what the box is designed to hold:
- Numeric boxes get `0` (including `0.0` for floats).
- Boolean boxes get `false`.
- String boxes get `""` (empty string).
- Container boxes (slices, maps, channels) and pointer boxes get `nil`.
- Interface boxes get `nil`.
- Function boxes get `nil`.

The key insight: `nil` is a valid value. It is not "nothing" — it is a specific zero value that you can check for and that some built-in functions handle gracefully.

## Core idea

Every type in Go has a zero value:

| Type | Zero value |
|---|---|
| `bool` | `false` |
| `int`, `int8`, ..., `uint64` | `0` |
| `float32`, `float64` | `0.0` |
| `complex64`, `complex128` | `0+0i` |
| `string` | `""` |
| `pointer` | `nil` |
| `slice` | `nil` (but `len` and `cap` are `0`) |
| `map` | `nil` |
| `channel` | `nil` |
| `function` | `nil` |
| `interface` | `nil` |
| `struct` | Each field gets its own zero value |
| `array` | Each element gets its own zero value |

The zero value for `nil` types is not useless. The Go runtime defines safe behavior for many operations on nil values:
- Calling `len()` and `cap()` on a nil slice returns `0`.
- Ranging over a nil slice or map iterates zero times (no panic).
- Receiving from a nil channel blocks forever (useful in select patterns).
- Calling methods on a nil receiver is allowed (the method can handle it).

However, writing to a nil map or indexing a nil slice causes a panic.

## Under the hood

The Go compiler implements zero-value initialization by emitting an instruction sequence that zeros the allocated memory. For local variables, the compiler may use `REP STOSQ` (on x86) or equivalent to fill the stack frame with zeros. For heap-allocated variables, the runtime's `newobject` calls `memclrNoHeapPointers` to zero the memory. The compiler also performs escape analysis: if a variable's address does not escape the function, it stays on the stack and is zeroed there; otherwise it is heap-allocated and zeroed by the allocator.

For composite types (structs, arrays), zeroing is recursive but is implemented as a single bulk zero operation on the entire memory block rather than field-by-field — the compiler knows the total size and emits a single memset.

## How Go uses it

- **`var buf bytes.Buffer`** — the zero value of `bytes.Buffer` is an empty buffer ready to use. No constructor call needed.
- **`var mu sync.Mutex`** — the zero value of `sync.Mutex` is an unlocked mutex. This is why Go structs often document "zero value is ready to use."
- **`http.HandlerFunc`** can be `nil` before being assigned. Code checks `if h != nil` before calling.
- **Maps** are `nil` by default. A nil map is safe to read (returns zero value for missing keys) but panics on write. This is why `make(map[string]int)` is needed before insertion.
- **Slices** are `nil` by default but `append` works on nil slices — it allocates the underlying array on first append.

## Go example

```go
package main

import "fmt"

type Config struct {
	Host string
	Port int
	Debug bool
}

func main() {
	// zero values for basic types
	var i int
	var f float64
	var s string
	var b bool
	fmt.Printf("int: %d, float64: %.1f, string: %q, bool: %t\n", i, f, s, b)

	// zero value for pointer
	var p *int
	fmt.Printf("pointer: %v, p == nil: %t\n", p, p == nil)

	// nil slice is safe to range over and append
	var nums []int
	fmt.Printf("nil slice: len=%d, cap=%d, isNil=%t\n", len(nums), cap(nums), nums == nil)
	for range nums {
		// never executes
	}
	nums = append(nums, 1) // append works on nil slices
	fmt.Println("after append:", nums)

	// nil map — safe to read, panics on write
	var m map[string]int
	fmt.Printf("nil map: len=%d, m == nil: %t\n", len(m), m == nil)
	_ = m["key"] // safe — returns zero value (0)
	// m["key"] = 1 // would panic: assignment to entry in nil map

	// nil channel — safe to range over (never yields)
	var ch chan int
	fmt.Printf("nil channel: %v, ch == nil: %t\n", ch, ch == nil)

	// zero value struct
	var cfg Config
	fmt.Printf("zero struct: Host=%q, Port=%d, Debug=%t\n", cfg.Host, cfg.Port, cfg.Debug)
}
```

## Step-by-step execution

For `var nums []int` followed by `nums = append(nums, 1)`:

1. `nums` is declared as a slice of int. Go sets nums to zero: `nil` (pointer = nil, len = 0, cap = 0).
2. The `for range nums` loop checks len(nums) which is 0 — zero iterations.
3. `append(nums, 1)` is called. The runtime checks: cap is 0, so a new underlying array is allocated (size 1).
4. `nums` is updated to point to the new array with len=1, cap=1, element [0]=1.

For `var m map[string]int` then `_ = m["key"]`:

1. `m` is declared as a map. Go sets it to `nil`.
2. Reading `m["key"]` looks up in a nil map. The runtime handles this: returns the zero value for the value type (`0` for int), and the boolean `false`.
3. Writing `m["key"] = 1` would call the map assign runtime function, which detects the nil map and panics.

## Common mistakes

- Mistake: Writing to a nil map.
  - Why: `var m map[string]int` creates a nil map. `m["key"] = 1` panics.
  - Fix: Initialize with `make(map[string]int)` or a composite literal `map[string]int{}`.

- Mistake: Calling methods on a nil pointer receiver that dereferences the receiver.
  - Why: Go allows calling methods on nil receivers, but if the method accesses fields, it panics.
  - Fix: Check for nil at the start of the method.

- Mistake: Assuming a nil slice and an empty slice are the same.
  - Why: `var s []int` is nil. `s := []int{}` is non-nil but empty. JSON serialization treats them differently: nil → `null`, empty → `[]`.
  - Fix: Be explicit about which you want. Use `var` for nil, `make` or `{}` for empty.

## Debugging walkthrough

Consider this broken code:

```go
package main

import "fmt"

func main() {
	var scores map[string]int
	scores["alice"] = 95
	scores["bob"] = 87
	fmt.Println(scores)
}
```

**Symptom**: Runtime panic: `assignment to entry in nil map`.

**Investigation**: The panic traceback points to the first write `scores["alice"] = 95`. The variable `scores` was declared with `var scores map[string]int` which initializes it to `nil`.

**Root cause**: A nil map cannot accept writes. The zero value of a map is `nil`, and the Go runtime panics when you try to write to it.

**Fix**: Initialize the map before writing:

```go
scores := make(map[string]int)
scores["alice"] = 95
scores["bob"] = 87
```

Or use a composite literal:

```go
scores := map[string]int{
	"alice": 95,
	"bob":   87,
}
```

## Production notes

- **Prefer zero-value initialization** for structs that document "zero value is ready to use." The standard library does this extensively (`bytes.Buffer`, `sync.Mutex`, `log.Logger`).
- **Use nil checks for optional fields**: `if config.Logger != nil { config.Logger.Println(...) }` instead of requiring a noop logger.
- **JSON marshaling** distinguishes nil and empty slices. If your API should return `[]` instead of `null`, initialize the slice: `items := make([]string, 0)`.
- **Nil channels in select** are an advanced pattern: a nil channel is never selected, so you can dynamically enable/disable cases by assigning nil to a channel.
- **Factory functions** are unnecessary when the zero value is useful. Only write `NewSomething()` when initialization is required beyond zeroing.

## Performance implications

- Zero-value initialization is a single `memclr` or bulk zero instruction. For stack variables, there is no runtime cost at all on most architectures — the stack frame is wiped at function entry.
- For heap-allocated variables, the allocator zeroes memory during allocation. This is O(n) in the size of the allocation, but the memory bandwidth is usually the bottleneck.
- Comparing `nil` to a pointer is a single register comparison — essentially free.
- Nil slices and nil maps use minimal memory (a few machine words for the descriptor). An empty slice (`[]int{}`) allocates a zero-length underlying array which still has a non-nil pointer — slightly more memory.

## Practice task

Write a function `safeSet(m map[string]int, key string, value int) map[string]int` that:
- If `m` is nil, creates a new map with `make`.
- Sets `m[key] = value`.
- Returns the map.

Write another function `safeGet(m map[string]int, key string) (int, bool)` that returns the value and whether the key exists, handling nil maps gracefully.

In `main()`, demonstrate calling `safeSet` with a nil map, then with a non-nil map, and use `safeGet` to retrieve values. Also demonstrate a nil slice that is safe to range over and append to.

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/04-zero-values
go test ./curriculum/modules/03-programming-fundamentals/lessons/04-zero-values
```

## Review questions

1. What is the zero value of a pointer? What happens if you dereference a nil pointer in Go?
2. Explain the difference between `var s []int` and `s := []int{}`. When would you choose one over the other for JSON serialization?
3. Debug: `var m map[string]int; m["x"] = 1; fmt.Println(m["x"])` — why does `m["x"]` print `0` even though the write panics? (Trick question: the write panics first.)
4. Tradeoff: Go guarantees zero values for all variables. What are the tradeoffs compared to a language that requires explicit initialization? Consider safety, performance, and developer experience.
5. A nil slice has `len` and `cap` of 0. Is a nil slice equal to `nil`? Is an empty slice (`[]int{}`) equal to `nil`? What does `reflect.DeepEqual` say?

## NEXT UP

Type conversions — how to convert between Go's basic types with explicit syntax, the strconv package, and when to use unsafe conversions.
