# Slice sharing and aliasing

## Learning objective

Identify when multiple slices share the same backing array, predict aliasing bugs, and write defensive code using `copy()` and capacity-limited slicing.

## Why this matters

Aliasing is the #1 source of mysterious slice bugs in production Go code. Two seemingly independent slices silently sharing the same backing array cause data corruption that is hard to reproduce and harder to debug. Every Go engineer must understand when sharing happens and how to break it.

## Mental model

Multiple slices can be windows into the same building (backing array). If one window's occupant paints a wall, every other window that looks at that wall sees the paint. Some windows are smaller (shorter length) or positioned differently (different start offset), but the building is one.

When you `append` to a slice and capacity is sufficient, you paint on the existing wall. If capacity is exceeded, the runtime builds a new building and points the window there — the old building is safe from that window's future appends, but other windows still point to it.

## Core idea

**Aliasing** occurs when two or more slice headers point to the same backing array. This happens through:

- **Assignment**: `t := s` — both headers share the array.
- **Slice expression**: `t := s[1:3]` — shares with offset.
- **Passing to functions**: receives a copy of the header, but the same backing array.
- **`append` returning the original header**: if within capacity, no new array.

**Dangers**:
- Modifying `t[i]` changes `s[i]` if they overlap.
- Appending to `t` within capacity overwrites data in `s` beyond `len(t)`.
- Slices in goroutines: unsynchronized writes through aliased slices are data races.

**Defensive techniques**:
1. `copy(dst, src)` — duplicates elements into a distinct backing array.
2. Full slice expression `s[low:high:max]` — limits capacity so `append` forces allocation.
3. Clone idiom: `append([]T(nil), src...)`.
4. `slices.Clone` (Go 1.21+) for safe copying.

## Under the hood

The slice header is just `{Data, Len, Cap}`. Two headers share an array when their `Data` pointer falls within the same allocated memory block (and their element type matches). The runtime has no tracking of how many slices reference a backing array. The garbage collector will free the backing array only when all references to it are gone.

`copy(dst, src)` calls `runtime.memmove` which handles overlapping regions correctly (it detects overlap and copies forward or backward as needed).

`slices.Clone(s)` is implemented as `append([]T(nil), s...)` — it allocates a new backing array and copies all elements.

## How Go uses it

- **`append` is copy-on-write-ish**: Within capacity, no copy; beyond capacity, new array.
- **`sort` package**: Sorts slices in place, modifying the backing array shared with the caller.
- **`io.Reader.Read`**: Populates a `[]byte` slice provided by the caller — sharing is expected.
- **`bytes.Buffer.Bytes()`**: Returns a slice sharing the buffer's backing array. Modifying it corrupts the buffer.
- **`encoding/gob`**: Uses shared slices during decode for zero-copy deserialization.

## Go example

```go
package main

import "fmt"

func main() {
	// Aliasing via slice expression
	original := []int{1, 2, 3, 4, 5, 6}
	sub := original[1:4] // shares backing array
	fmt.Println("original:", original, "sub:", sub)

	// Modify through alias
	sub[0] = 99
	fmt.Println("after sub[0]=99 — original:", original)

	// Append within capacity overwrites
	sub = append(sub, 100)
	fmt.Println("after append to sub — original:", original)

	// Defensive copy
	safe := make([]int, 3)
	copy(safe, original[1:4])
	safe[0] = -1
	fmt.Println("original after safe modification:", original)
	fmt.Println("safe:", safe)

	// Clone idiom
	clone := append([]int(nil), original...)
	clone[0] = 777
	fmt.Println("original after clone modification:", original)
	fmt.Println("clone:", clone)
}
```

## Step-by-step execution

For `original := []int{1,2,3,4,5,6}; sub := original[1:4]; sub = append(sub, 100)`:

1. `original` has backing array `[1 2 3 4 5 6]`, `len=6, cap=6`.
2. `sub := original[1:4]` — header: `Data=&arr[1], Len=3, Cap=5`. Views `[2 3 4]`.
3. `sub = append(sub, 100)` — `len=3 < cap=5`, so writes `arr[4] = 100`. `sub` now `[2 3 4 100]`, `len=4, cap=5`.
4. `original[4]` is now `100`, not `5`. The original slice sees the change because it shares the array.

For `copy(safe, original[1:4])`:

1. Allocates new `[3]int` on the heap (or stack if small enough).
2. `memmove` copies 3 integers from `&arr[1]` to `safe[0]`.
3. `safe` is now completely independent — modifying `safe` does not touch `original`.

For `clone := append([]T(nil), original...)`:

1. `[]int(nil)` is a nil slice (`Len=0, Cap=0, Data=nil`).
2. `append` on a nil slice allocates a new backing array.
3. All elements of `original` are copied into the new array.
4. `clone` is independent.

## Common mistakes

- Mistake: Assuming `append` never modifies the original slice.
  - Why it happens: Within capacity, `append` modifies the backing array. The original slice may see changes in the region beyond its length.
  - Fix: Use full slice expression or `copy` to isolate.

- Mistake: Storing a slice returned from a method that documents "the slice shares memory".
  - Why it happens: `bytes.Buffer.Bytes()` returns a slice sharing the buffer's internal array. Appending to the returned slice corrupts the buffer.
  - Fix: Copy the data immediately: `buf := append([]byte(nil), b.Bytes()...)`.

- Mistake: Concurrent goroutines appending to aliased slices.
  - Why it happens: Two goroutines holding slices with the same backing array call `append`. Even if `append` allocates, the read and write to the old backing array race.
  - Fix: Give each goroutine its own copy, or use explicit synchronization.

- Mistake: Using `append` on a sub-slice in a loop and accumulating results.
  - Why it happens: Each iteration may share or not share the backing array depending on capacity, leading to intermittent corruption.
  - Fix: `copy` the elements you need at each iteration, or preallocate.

## Debugging walkthrough

```go
package main

import "fmt"

func main() {
	data := []int{0, 1, 2, 3, 4, 5}
	results := [][]int{}
	for i := 0; i < len(data); i += 2 {
		chunk := data[i : i+2]
		results = append(results, chunk)
	}
	data[1] = 99
	data[3] = 88
	fmt.Println("results:", results)
}
```

**Symptom**: Prints `results: [[0 99] [2 88] [4 5]]` — all chunks share the same backing array, so mutating `data` changes all chunks.

**Investigation**: After the loop, inspect `results`:

```go
for i, ch := range results {
	fmt.Printf("chunk %d: %v, ptr=%p\n", i, ch, ch)
}
```

Output:
```
chunk 0: [0 99], ptr=0xc000010020
chunk 1: [2 88], ptr=0xc000010024
chunk 2: [4 5], ptr=0xc000010028
```

All three are contiguous pointers into `data`.

**Root cause**: `chunk := data[i:i+2]` does not copy — it aliases. `append(results, chunk)` appends the alias, not a copy.

**Fix**: Copy each chunk:

```go
chunk := make([]int, 2)
copy(chunk, data[i:i+2])
```

Or use full slice expression and immediate append copy:

```go
chunk := append([]int(nil), data[i:i+2]...)
```

## Production notes

- **API contracts**: Document whether a returned slice shares memory. For example, `os.ReadFile` returns a new `[]byte` that you own. `bytes.Buffer.Bytes()` explicitly warns about sharing.
- **Zero-copy optimization**: In hot paths, intentional aliasing avoids allocation. This is valid when ownership is clear (e.g., `sort.Slice` sorts in place).
- **`append` inside a loop**: Always assume `append` may reallocate, but if it doesn't, all loop iterations share the same backing array. Copy to a target slice to be safe.
- **`strings.Split` returns substrings sharing the input backing array**: If you keep substrings after modifying the input, the input's backing array may be garbage-collected, but the substrings still reference it — no issue because Go's GC tracks via `SliceHeader.Data`, but the split strings pin the entire backing array, increasing memory usage.

## Performance implications

- `copy` is O(n) and uses the fastest memory copy the CPU supports (AVX, SSE, etc. on modern hardware).
- Intentional aliasing (not copying) is free — no allocation, no `memmove`. Use it when the data flow is unambiguous.
- Full slice expressions cost nothing at runtime — they just set the `Cap` field in the header.
- Cloning large slices unnecessarily is the #1 memory waste in Go programs. Only clone when you need isolation.
- `slices.Clone` from Go 1.21+ is the idiomatic way to clone; it's as fast as `append([]T(nil), s...)`.

## Practice task

Write a function `splitAndCopy(data []int, chunkSize int) [][]int` that splits `data` into non-overlapping chunks of `chunkSize` without aliasing — each chunk must have its own backing array. Then write a function `unsafeSplit(data []int, chunkSize int) [][]int` that does the same but with aliasing. In `main()`, demonstrate that modifying `data[i]` after calling `splitAndCopy` does not affect the chunks, but doing the same after `unsafeSplit` does.

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/16-slice-sharing-and-aliasing
go test ./curriculum/modules/03-programming-fundamentals/lessons/16-slice-sharing-and-aliasing
```

## Review questions

1. After `t := s[2:5]`, do `t` and `s` share the same backing array?
2. How can you append to a sub-slice without risking overwriting the parent's data?
3. What does `copy()` do when the source and destination overlap? Is it safe?
4. Why does `bytes.Buffer.Bytes()` document that the returned slice shares memory?
5. How would you safely clone a slice in Go 1.21+? In earlier versions?

## NEXT UP

Maps — Go's built-in hash table for key-value associations.
