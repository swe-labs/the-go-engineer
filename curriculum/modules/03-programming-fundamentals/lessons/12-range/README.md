# Range

## Learning objective

Iterate over slices, arrays, maps, strings, and channels using `for range`, understand the copy semantics of range variables, use the blank identifier to discard values, and distinguish single-value from two-value range forms.

## Why this matters

`for range` is the most common loop in idiomatic Go. It appears in virtually every production codebase for reading data, computing aggregates, serializing structures, and processing streaming input. Understanding exactly what `range` copies, how it handles runes in strings, why map iteration order is non-deterministic, and how channel range terminates is essential for writing correct, efficient Go. Misunderstanding range semantics is a prolific source of bugs — especially the copy behavior, which surprises even experienced engineers.

## Mental model

`for range` is a compiler transform that expands into an ordinary `for` loop with specific iteration machinery for each data structure.

```
for i, v := range collection {
```

- `i` is the **index/key** — the position in the collection.
- `v` is the **value** — a **copy** of the element at that position.
- The loop automatically advances `i` and updates `v` each iteration.
- The loop terminates when all elements have been visited.

The blank identifier `_` discards either value:

```go
for i, _ := range slice  // discard value, keep index
for _, v := range slice  // discard index, keep value
for i := range slice     // single-value range: index only (most efficient)
```

Each collection type has unique range behaviour:

| Type | Key/index | Value | Termination |
|---|---|---|---|
| Slice/Array | `int` 0 to len-1 | copy of element | after last element |
| String | `int` byte offset | `rune` (the Unicode code point) | after last byte |
| Map | key type | copy of value | after all keys visited (non-deterministic order) |
| Channel | never assigned | value received | after channel is closed |

## Under the hood

The compiler expands `range` differently per type:

**Slice**: `for i := 0; i < len(s); i++ { v := s[i]; ... }`. The slice length is captured before iteration begins, so appending to the slice inside the loop does not extend the iteration.

**Array**: Similar to slice, but the array is copied if it is a range variable (unless you range over a pointer to the array). `for i, v := range arr` copies the entire array, then iterates over the copy. This is expensive for large arrays — use `for i, v := range &arr` to avoid the copy.

**String**: The compiler calls `utf8.DecodeRuneInString(s[i:])` on each iteration. The index is the byte offset of the start of the current rune. If the string contains invalid UTF-8, the rune is `U+FFFD` (replacement character) with width 1.

**Map**: The runtime provides `runtime.mapiterinit` and `runtime.mapiternext`. A random bucket offset is chosen at init time to prevent deterministic iteration (a security measure against hash-collision DoS attacks). The iteration produces all keys exactly once, but the order is intentionally unpredictable.

**Channel**: The compiler expands to `for { v, ok := <-ch; if !ok { break }; ... }`. The loop blocks until a value is received or the channel is closed.

## How Go uses it

- **Processing collections**: `for _, item := range items { ... }` — the canonical slice iteration.
- **Counting runes**: `for range "Hello, 世界" { n++ }` — counts Unicode code points.
- **Map lookups with iteration**: `for k, v := range m { ... }` — process all key/value pairs.
- **Channel fan-in**: `for v := range ch { results = append(results, v) }` — consume all values until close.
- **Zero-allocation read loops**: `for i := range buf { buf[i] = 0 }` — clear a buffer without copying elements.
- **Generating sequences**: `for i := range make([]struct{}, n) { ... }` — a pattern for running a loop exactly `n` times without an explicit index (uses the `range` over a slice of empty structs, which has zero allocation).

## Go example

```go
package main

import "fmt"

func main() {
	// Slice range
	nums := []int{10, 20, 30}
	fmt.Println("--- slice ---")
	for i, v := range nums {
		fmt.Printf("i=%d, v=%d\n", i, v)
	}

	// Array range (copy semantics)
	arr := [3]int{1, 2, 3}
	fmt.Println("--- array (copy) ---")
	for i, v := range arr {
		if i == 0 {
			arr[1] = 999 // modifies original, NOT the range copy
		}
		fmt.Printf("i=%d, v=%d\n", i, v)
	}
	fmt.Println("arr after:", arr) // [1, 999, 3]

	// String range (runes)
	fmt.Println("--- string ---")
	s := "Hi, 世界"
	for i, r := range s {
		fmt.Printf("byte=%d, rune=%c (U+%04X)\n", i, r, r)
	}

	// Map range (non-deterministic order)
	fmt.Println("--- map ---")
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	for k, v := range m {
		fmt.Printf("k=%q, v=%d\n", k, v)
	}

	// Single-value range
	fmt.Println("--- single-value (indices) ---")
	for i := range nums {
		fmt.Printf("index=%d\n", i)
	}

	// Discard index
	fmt.Println("--- discard index ---")
	for _, v := range nums {
		fmt.Printf("value=%d\n", v)
	}
}
```

## Step-by-step execution

Trace `for i, r := range "Hi, 世界"`:

1. Compiler computes string length: 9 bytes.
2. Iteration 0: `i=0`. Decode rune at byte 0: `'H'` (1 byte, U+0048).
3. Iteration 1: `i=1`. Decode rune at byte 1: `'i'` (1 byte, U+0069).
4. Iteration 2: `i=2`. Decode rune at byte 2: `','` (1 byte, U+002C).
5. Iteration 3: `i=3`. Decode rune at byte 3: `' '` (1 byte, U+0020).
6. Iteration 4: `i=4`. Decode rune at byte 4: `'世'` (3 bytes, U+4E16).
7. Iteration 5: `i=7`. Decode rune at byte 7: `'界'` (3 bytes, U+754C).
8. Iteration 6: `i=10`. `i >= len(s)` → exit.

Trace range over map with `m := map[string]int{"a": 1, "b": 2, "c": 3}`:

1. Compiler calls `runtime.mapiterinit(m)` which picks a random starting bucket.
2. Iteration: returns `("b", 2)` — unpredictable. Body executes.
3. Next: returns `("a", 1)` — unpredictable. Body executes.
4. Next: returns `("c", 3)`. Body executes.
5. Next: returns nil key (no more entries). Loop exits.

The exact order varies between runs, even on the same machine.

## Common mistakes

- **Modifying map while ranging**: Adding or deleting map entries during a range loop can cause unpredictable behaviour. Delete is safe (the deleted key may or may not appear), but insert can cause a runtime panic (concurrent map write) or data races.

- **Assuming map iteration order is deterministic**: Never rely on map iteration order. If you need stable order, collect keys into a slice and sort it.

- **Expecting range over array to not copy**: `for i, v := range arr` copies the entire array. Use `for i, v := range &arr` to avoid the copy for large arrays.

- **Appending to a slice while ranging**: `for i, v := range s { s = append(s, v) }` — the range captures the original length before iteration. Appended elements are not iterated.

- **Using range variable address outside the loop**: `for _, v := range s { go func() { fmt.Println(v) }() }` — before Go 1.22, `v` is the same variable across iterations, causing all goroutines to see the last value. In Go 1.22+, each iteration gets a new `v`.

- **Forgetting that string range yields runes, not bytes**: `for i := 0; i < len(s); i++` iterates bytes. `for i, r := range s` iterates runes. The index `i` is a byte offset, not a rune index.

## Debugging walkthrough

```go
package main

import "fmt"

func main() {
	m := map[string]int{"alice": 1, "bob": 2, "carol": 3}
	for k, v := range m {
		if k == "bob" {
			delete(m, "carol")
		}
		fmt.Println(k, v)
	}
}
```

**Symptom**: Sometimes `carol` appears, sometimes it does not. The behaviour is unpredictable.

**Root cause**: Deleting from a map during range iteration is safe but the deleted key may or may not be visited. The Go spec: "If map entries that have not yet been reached are removed during iteration, the corresponding iteration values will not be produced."

**Fix**: Collect keys first, then delete:

```go
for k := range m {
	if k == "bob" {
		delete(m, "carol")
	}
}
```

Collecting in a separate loop is fine if you only delete known keys. For bulk deletion, collect keys first.

Now consider:

```go
type Item struct {
	Value int
}

func main() {
	items := []Item{{1}, {2}, {3}}
	for _, item := range items {
		item.Value = 0
	}
	fmt.Println(items) // [{1} {2} {3}] — unchanged!
}
```

**Root cause**: `item` is a copy of the slice element. Modifying it has no effect on the original.

**Fix**: Use the index: `items[i].Value = 0`, or range over pointers: `for i := range items { items[i].Value = 0 }`.

## Production notes

- **Use single-value range over slice when you only need indices**: `for i := range slice` is more efficient and clearer than `for i, _ := range slice`.
- **Avoid map writes during range**: If you must modify, collect keys first in a separate loop.
- **For stable map iteration, sort keys**: `keys := make([]string, 0, len(m)); for k := range m { keys = append(keys, k) }; sort.Strings(keys); for _, k := range keys { v := m[k] ... }`.
- **Range over channel is blocking**: Ensure the channel is closed by the sender, not the receiver. A range over a nil channel blocks forever — a common leak.
- **String range with invalid UTF-8**: Always handle replacement character `U+FFFD` if your strings may contain invalid UTF-8. Use `utf8.ValidString` to pre-check if needed.
- **Zero-allocation n-iteration**: `for i := range n { ... }` does not work directly. Use `for i := 0; i < n; i++` or `for range make([]struct{}, n) { ... }` (the latter creates a zero-size slice, so no heap allocation occurs).

## Performance implications

- Range over a slice is **bounds-check-elimination friendly**: The compiler knows the loop variable `i` stays within `[0, len(s))` and can eliminate bounds checks on `s[i]` inside the body.
- Range over a map **allocates an iterator** on the heap. Each iteration calls into the runtime. Map iteration is 5-10x slower than slice iteration per element.
- Range over a string has **O(n) complexity** in bytes, but each rune decode involves a function call to the UTF-8 decoder. For ASCII-only strings, convert to `[]byte` first: `for _, b := range []byte(s)`.
- Range over a channel **blocks the goroutine**. The scheduler may park the goroutine until a value arrives or the channel is closed.
- Array range **copies the entire array**. For large arrays (≥1KB), range over `&arr` to avoid the copy.

## Practice task

Write a function `wordCount(s string) map[string]int` that:
- Splits `s` into words (sequences of letters).
- Returns a map of word → count.
- Uses `for range` to iterate over the string's runes.

Write a function `mergeCounts(maps ...map[string]int) map[string]int` that merges multiple word-count maps, summing values for duplicate keys.

In `main()`, call `wordCount` on `"hello world hello"` and verify the result. Then merge two counts and print the result.

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/12-range
go test ./curriculum/modules/03-programming-fundamentals/lessons/12-range
```

## Review questions

1. Does `for i, v := range arr` (where `arr` is `[3]int`) copy the array? How would you avoid the copy?
2. Why is the iteration order of a map non-deterministic in Go?
3. What is the difference between `for i := range s` and `for i := 0; i < len(s); i++` when `s` is a string?
4. What happens if a channel is never closed and you use `for v := range ch`?
5. For `for i, v := range slice`, is `v` an alias for `slice[i]` or a copy? How do you modify elements in place?

## NEXT UP

Arrays — fixed-size indexed collections with value semantics in Go.
