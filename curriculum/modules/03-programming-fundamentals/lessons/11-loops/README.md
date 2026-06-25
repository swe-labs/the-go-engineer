# Loops

## Learning objective

Write all forms of Go's `for` loop — init/condition/post, condition-only (while-style), infinite, `for range`, and use `break`, `continue`, and labeled break/continue for flow control.

## Why this matters

Loops are the fundamental mechanism for processing collections, retrying operations, waiting for events, and generating sequences. Go unifies all looping patterns under a single keyword `for`, eliminating the `while` and `do-while` variants of other languages. This simplicity means every Go engineer must be fluent in the three `for` forms and in the precise semantics of `break` and `continue`. Labeled breaks — a feature omitted from many languages — are essential for breaking out of nested loops cleanly in Go.

## Mental model

The `for` keyword is a control structure that repeats a block of code. It has three forms, all derived from the same core concept:

```
for initialization; condition; post {
	// body
}
```

Each iteration:
1. The `initialization` runs once before the first iteration (optional).
2. The `condition` is evaluated before each iteration. If `false`, the loop exits.
3. The `body` executes.
4. The `post` statement runs after each iteration (optional).
5. Back to step 2.

Omitting parts produces the other forms:

- `for condition { }` — while-style (no init, no post).
- `for { }` — infinite loop (nothing at all).
- `for i, v := range slice { }` — range loop (special form).

`break` exits the innermost enclosing `for` (or `switch`/`select`). `continue` skips to the next iteration. A label before a loop allows `break` or `continue` to target a specific enclosing loop.

## Core idea

```go
// 1. Init/condition/post (C-style)
for i := 0; i < 10; i++ {
	fmt.Println(i)
}

// 2. Condition-only (while-style)
n := 0
for n < 10 {
	n++
}

// 3. Infinite loop
for {
	// runs until break/return
}

// 4. Range loop
nums := []int{1, 2, 3}
for i, v := range nums {
	fmt.Println(i, v)
}
```

`break` and `continue`:

```go
for i := 0; i < 10; i++ {
	if i%2 == 0 {
		continue // skip even numbers
	}
	if i > 7 {
		break // exit early
	}
	fmt.Println(i) // prints: 1, 3, 5, 7
}
```

Labeled break/continue:

```go
outer:
for i := 0; i < 3; i++ {
	for j := 0; j < 3; j++ {
		if i*j > 2 {
			break outer // exits both loops
		}
		fmt.Printf("(%d,%d) ", i, j)
	}
}
// Output: (0,0) (0,1) (0,2) (1,0) (1,1) (1,2) (2,0)
```

## Under the hood

The `for` loop compiles to a sequence of machine instructions centered on a backward jump:

```
    init← executed once before loop
loop:
    check condition → jump to end if false
    body
    post
    jump to loop
end:
```

The `range` keyword is syntactic sugar the compiler expands into an ordinary `for` loop. For slices and arrays, the compiler generates:
- A pointer to the underlying array (to avoid copying the header on each iteration).
- An index variable and a copy of each element.

For maps, the compiler invokes the runtime map iterator (`runtime.mapiterinit`, `runtime.mapiternext`) which produces keys in non-deterministic order. The map iterator is designed to resist hash-flooding attacks by randomizing iteration order.

For strings, `for range` decodes UTF-8 runes on the fly using the runtime `utf8.DecodeRuneInString` function.

For channels, `for v := range ch` is equivalent to `for { v, ok := <-ch; if !ok { break } ... }`.

## How Go uses it

- **Slice iteration**: `for i, v := range items { ... }` — the most common loop in Go.
- **Map iteration**: `for k, v := range m { ... }`.
- **String rune iteration**: `for i, r := range "Hello, 世界" { ... }` yields rune-by-rune.
- **Channel consumption**: `for msg := range ch { ... }` — reads until channel is closed.
- **Retry loops**: `for retries := 0; retries < 3; retries++ { ... }` with `break` on success.
- **Polling**: `for { select { case <-ctx.Done(): return; case <-time.After(time.Second): ... } }`.
- **Nested loop exit**: Labeled break to exit outer loops from deep inside.

## Go example

```go
package main

import "fmt"

func main() {
	fmt.Println("--- init/cond/post ---")
	for i := 0; i < 5; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	fmt.Println("--- while-style ---")
	n := 0
	for n < 5 {
		fmt.Printf("%d ", n)
		n++
	}
	fmt.Println()

	fmt.Println("--- infinite with break ---")
	count := 0
	for {
		if count >= 3 {
			break
		}
		fmt.Printf("%d ", count)
		count++
	}
	fmt.Println()

	fmt.Println("--- continue ---")
	for i := 0; i < 10; i++ {
		if i%2 != 0 {
			continue
		}
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	fmt.Println("--- labeled break ---")
outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				fmt.Print("BREAK ")
				break outer
			}
			fmt.Printf("(%d,%d) ", i, j)
		}
	}
	fmt.Println()

	fmt.Println("--- labeled continue ---")
outer2:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == 1 {
				continue outer2
			}
			fmt.Printf("(%d,%d) ", i, j)
		}
	}
	fmt.Println()
}
```

## Step-by-step execution

Trace `for i := 0; i < 3; i++ { if i == 1 { continue }; fmt.Print(i) }`:

1. Init: `i = 0`. Check condition: `0 < 3` → true.
2. Body: `i == 1` → false. Print `0`.
3. Post: `i++` → `i = 1`.
4. Check condition: `1 < 3` → true.
5. Body: `i == 1` → true. `continue` → skip rest of body, go to post.
6. Post: `i++` → `i = 2`.
7. Check condition: `2 < 3` → true.
8. Body: `i == 1` → false. Print `2`.
9. Post: `i++` → `i = 3`.
10. Check condition: `3 < 3` → false. Exit loop.
Output: `02`.

Trace labeled break in nested loops:

1. Outer `i=0`. Inner `j=0`: print `(0,0)`. Inner `j=1`: print `(0,1)`. Inner `j=2`: print `(0,2)`.
2. Outer `i=1`. Inner `j=0`: print `(1,0)`. Inner `j=1`: `i==1 && j==1` → `break outer`. Jump to after outer loop.
3. Inner `j=2` and outer `i=2` never execute.

## Common mistakes

- **Modifying the range variable**: `for i, v := range slice { v = 0 }` — `v` is a copy, so the slice is not modified. Use `slice[i] = 0`.
- **Infinite loop from missing increment**: `for i := 0; i < 10; { fmt.Println(i) }` — `i` never changes. Add the post statement.
- **Assuming `break` in a switch exits the loop**: `break` inside a `switch` that is inside a `for` breaks the switch, not the loop. Use a labeled break or a boolean flag.
- **Forgetting that `range` over a map copies values**: `for k, v := range m { v.Field = 1 }` modifies the copy, not the map value. Access the map directly: `m[k] = ...`.
- **Closing a channel inside a range loop**: `for v := range ch { close(ch) }` — panics: close of closed channel. Use `close(ch)` only in the sender goroutine.
- **Using `range` on nil slice/map**: Ranging over a nil slice or nil map is fine — it iterates zero times. Ranging over a nil channel blocks forever.

## Debugging walkthrough

```go
package main

import "fmt"

func main() {
	nums := []int{1, 2, 3, 4, 5}
	for _, v := range nums {
		if v%2 == 0 {
			v = 0 // Tries to zero out evens?
		}
	}
	fmt.Println(nums) // [1 2 3 4 5] — unchanged!
}
```

**Symptom**: Slice is not modified.

**Root cause**: `v` is a copy of each element. Assigning to `v` does not affect the original slice.

**Fix**: Use the index:

```go
for i, v := range nums {
	if v%2 == 0 {
		nums[i] = 0
	}
}
```

Now consider:

```go
package main

import "fmt"

func main() {
	for i := 0; i < 5; i++ {
		switch i {
		case 3:
			break
		}
		fmt.Print(i)
	}
}
```

**Symptom**: Loop does not stop at 3. Output: `01234`.

**Root cause**: `break` in a `switch` exits the switch, not the for loop. The for loop continues to the next increment.

**Fix**: Use a labeled break or restructure:

```go
loop:
	for i := 0; i < 5; i++ {
		switch i {
		case 3:
			break loop
		}
		fmt.Print(i)
	}
// Output: 012
```

## Production notes

- **Prefer `for range` over C-style for slices/maps/strings**: It is less error prone and automatically handles bounds.
- **Avoid `for range` channel loops that never close**: A range over a channel that is never closed leaks the goroutine and blocks forever. Ensure channels are closed appropriately.
- **Use labeled break sparingly**: It is powerful but can confuse readers. Prefer extracting the inner loop into a function that returns early.
- **Watch for performance**: `for i := range slice` avoids copying the element. Use `for i := range slice` instead of `for i, v := range slice` when you only need indices.
- **Loop variable capture in closures**: Before Go 1.22, loop variables were captured by reference, causing bugs with goroutines. In Go 1.22+, each iteration gets a new variable. If supporting older versions, copy loop variables explicitly.

## Performance implications

- `for range` over a slice is as fast as a C-style `for i := 0; i < len(s); i++`. The compiler eliminates bounds checks when it can prove the index is within bounds.
- `for range` over a map accesses the hash table on each iteration. The cost is O(1) per key but with a high constant factor (hashing, equality checks, overflow bucket traversal).
- Empty `for {}` without any blocking operation consumes 100% CPU. Use `runtime.Gosched()` or a synchronization primitive inside the loop.
- Loop-invariant code is hoisted out of the loop by the compiler's optimization passes (but manually hoisting improves readability).
- Range over a string (`for i, r := range s`) calls `utf8.DecodeRuneInString` per iteration. For ASCII-only strings, convert to `[]byte` and iterate that for better performance.

## Practice task

Write a function `sumUntil(nums []int, target int) (int, int)` that:
- Iterates over `nums` using a C-style `for` loop.
- Returns the sum of elements and the count of elements used.
- Stops when the cumulative sum reaches or exceeds `target`.
- Uses `break` when the target is reached.
- Uses `continue` to skip negative values (they do not count toward the sum or count).

In `main()`, test with `[]int{3, -1, 5, 2, -3, 4}` and `target = 10`. Print the result.

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/11-loops
go test ./curriculum/modules/03-programming-fundamentals/lessons/11-loops
```

## Review questions

1. What are the three forms of `for` in Go?
2. What is the output of `for i := 0; i < 3; i++ { if i == 1 { continue }; fmt.Print(i) }`?
3. How do you break out of an outer loop from inside a nested loop?
4. Why does `for i, v := range slice { v = 0 }` not modify the slice? How would you fix it?
5. What happens when you range over a nil slice, nil map, or nil channel?

## NEXT UP

Range — deeper coverage of `for range` iteration over slices, arrays, maps, strings, and channels.
