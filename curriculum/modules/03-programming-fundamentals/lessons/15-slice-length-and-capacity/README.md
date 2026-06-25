# Slice length and capacity

## Learning objective

Distinguish between slice length and capacity, predict when `append` allocates, and use `cap`-aware slicing to control memory and growth behavior.

## Why this matters

Length and capacity govern how slices grow, when they allocate, and how much memory they waste. A junior developer treats `len` and `cap` as interchangeable. A senior engineer knows that cap controls allocation rate, that re-slicing up to `cap` is safe, and that pre-allocating capacity eliminates pointless copying. Every `append` call is a decision about capacity.

## Mental model

`len(s)` is how many elements you can read or write right now. `cap(s)` is how many total elements fit in the backing array before Go must allocate a new one.

Picture a hotel corridor: `len` is the number of occupied rooms; `cap` is the total number of rooms on the floor. You can keep opening doors (appending) until you run out of rooms, at which point the hotel builds a new floor with more rooms and moves everyone over.

## Core idea

- `len(s)` returns the number of elements in the slice — the accessible portion.
- `cap(s)` returns the capacity of the underlying backing array, measured from the first element of the slice.

Writing `s[i]` is valid only for `0 <= i < len(s)`. Reading or writing beyond `len(s)` but within `cap(s)` via a slice expression (`s[:cap(s)]`) is valid; via an index (`s[len(s)]`) is a panic.

The full slice expression `s[low:high:max]` limits the capacity of the resulting slice to `max - low` elements. This is a safety mechanism to prevent `append` from modifying unexpected parts of the backing array.

## Under the hood

`growslice` in the Go runtime handles capacity growth. The algorithm:

1. If the requested capacity fits in the current backing array, use it.
2. Otherwise, allocate a new backing array:
   - For `cap < 256`: double the capacity.
   - For `cap >= 256`: `newcap = cap + cap/4` (1.25x growth).
3. Round up to the nearest size class that the memory allocator supports.
4. Copy old elements to the new array.
5. Return new slice header (same pointer if no reallocation, new pointer otherwise).

The doubling rule at small sizes keeps amortized append cost O(1). The 1.25x rule at large sizes avoids wasting too much memory.

## How Go uses it

- **Pre-allocation**: `buf := make([]byte, 0, 1024)` for I/O buffers is standard — you know the expected size.
- **Pool recycling**: `s = s[:0]` resets length to 0 without deallocating the backing array.
- **Sub-slicing with capped capacity**: `header := raw[0:4:4]` prevents `append` on `header` from writing into `raw[4:]`.
- **`cap` checks**: Safe code checks `cap(s) > 100` before a `s[:100]` operation.
- **`hdr` decoding**: Network packet parsers use `len` and `cap` to track consumed bytes.

## Go example

```go
package main

import "fmt"

func main() {
	// Preallocated with extra capacity
	s := make([]int, 3, 8)
	fmt.Printf("initial: len=%d cap=%d s=%v\n", len(s), cap(s), s)

	// Append within capacity — no allocation
	s = append(s, 10, 20)
	fmt.Printf("after 2 appends: len=%d cap=%d s=%v\n", len(s), cap(s), s)

	// Re-slice within capacity
	t := s[1:4]
	fmt.Printf("t = s[1:4]: len=%d cap=%d t=%v\n", len(t), cap(t), t)

	// Full slice expression to limit capacity
	u := s[1:4:4]
	fmt.Printf("u = s[1:4:4]: len=%d cap=%d u=%v\n", len(u), cap(u), u)

	// append on limited-capacity slice forces allocation
	u = append(u, 99)
	fmt.Println("after u append:", u, "s:", s)

	// Reslice to full capacity
	full := s[:cap(s)]
	fmt.Println("full slice s[:cap(s)]:", full)
}
```

## Step-by-step execution

For `s := make([]int, 2, 4)` then `s = append(s, 1, 2)` then `s = append(s, 3)`:

1. `make([]int, 2, 4)` allocates backing array of size 4. `len=2`, `cap=4`.
2. `s[0]` and `s[1]` are zero. Slice `s` sees `[0, 0]`.
3. `append(s, 1)` — `len=2 < cap=4` → write `s[2]=1`, set `len=3`. No allocation.
4. `append(s, 2)` — `len=3 < cap=4` → write `s[3]=2`, set `len=4`. No allocation.
5. `append(s, 3)` — `len=4 == cap=4` → allocation needed. New cap = 8. Copy 4 elements, set `len=5`. Old backing array is no longer referenced by `s`.

For the full slice expression `t := s[1:4:4]`:

1. `low=1, high=4, max=4`.
2. `len(t) = 4-1 = 3`. `cap(t) = 4-1 = 3`.
3. `t` shares the backing array with `s` at offset 1.
4. `append(t, x)` creates capacity 3 < 3+1=4 → allocates new array. The data at `s[4]` is protected.

## Common mistakes

- Mistake: Re-slicing without realizing `len` vs `cap`.
  - Why it happens: `s[1:3]` creates a slice whose capacity starts at index 1, so `cap` is `originalCap - 1`. Appending to the sub-slice can overwrite `s[3]`.
  - Fix: Use `s[1:3:3]` full slice expression when you want `append` to allocate.

- Mistake: Iterating `for i := 0; i < cap(s); i++` expecting the same result as `range`.
  - Why it happens: `cap(s)` may be larger than `len(s)`. Indexing `s[i]` for `i >= len(s)` panics.
  - Fix: Always use `len(s)` or `for range` for iteration.

- Mistake: Not preallocating and paying for many reallocations.
  - Why it happens: `append` one element at a time with starting capacity 0 triggers O(log n) allocations.
  - Fix: `make([]T, 0, expectedN)` or estimate capacity.

- Mistake: Confusing `make([]int, n)` with `make([]int, 0, n)`.
  - Why it happens: `make([]int, n)` allocates and zeroes n elements (`len = n, cap = n`). `make([]int, 0, n)` allocates size n but `len = 0`. Appending to the first starts at position n.
  - Fix: Use `make([]int, 0, n)` if you want to start empty and append.

## Debugging walkthrough

```go
package main

import "fmt"

func main() {
	original := []int{1, 2, 3, 4, 5}
	sub := original[0:3]
	fmt.Println("sub before:", sub, "len:", len(sub), "cap:", cap(sub))
	sub = append(sub, 99)
	fmt.Println("sub after:", sub)
	fmt.Println("original:", original)
}
```

**Symptom**: `original` becomes `[1 2 3 99 5]` unexpectedly.

**Investigation**: Print `cap(sub)`:

```
sub before: [1 2 3] len: 3 cap: 5
```

**Root cause**: `sub` has `cap=5` because it was sliced from `original` which had `cap=5`. The `append` of `99` fits within capacity, so it writes to `original[3]`.

**Fix**: Use the full slice expression to limit capacity:

```go
sub := original[0:3:3]
```

Now `cap(sub) = 3`, and `append` allocates a new backing array.

## Production notes

- **Pre-allocation rule of thumb**: If you know n within 20%, `make([]T, 0, n)` saves ~50% of allocation time.
- **Capacity check before large operation**: `if cap(s) < needed { s = append(make([]T, 0, needed), s...) }` can be faster than a loop of appends.
- **Buffer reuse**: In servers, `buf = buf[:0]` is preferred over allocating a new `make([]byte, size)` for each request.
- **`cap(s) == 0`** is true for both nil and empty slices with capacity 0. Use `len(s) == 0` to check emptiness.

## Performance implications

- Reallocation cost: `growslice` is O(n) in the current length because it copies all elements. Frequent reallocation in hot loops can dominate CPU time.
- Capacity rounding: Go's memory allocator rounds up to size classes, so `make([]byte, 0, 10)` may actually allocate 16 bytes. This is not a bug; it's alignment.
- Small capacity overhead: The slice header itself is always 24 bytes on 64-bit. The backing array overhead is the allocated capacity times element size.
- `cap` checking is free: `len` and `cap` are fields in the header, read with a single MOV instruction.

## Practice task

Write a function `growSteps(initialCap, numAppends int) []int` that starts with a slice of length 0 and capacity `initialCap`, then appends one element at a time `numAppends` times. After each append, print the old and new capacity when a reallocation occurs. Return the final slice.

In `main()`, call `growSteps(2, 10)` to see when Go doubles capacity and when it switches to 1.25x growth.

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/15-slice-length-and-capacity
go test ./curriculum/modules/03-programming-fundamentals/lessons/15-slice-length-and-capacity
```

## Review questions

1. What is the result of `make([]int, 5, 10)[:10]`? Does it panic?
2. After `s := []int{1,2,3}; t := s[:2]`, what are `len(t)` and `cap(t)`?
3. True or false: `append` always allocates a new backing array. Explain.
4. What does the full slice expression `a[low:high:max]` prevent?
5. Why does Go double capacity for small slices but only grow by 25% for large slices?

## NEXT UP

Slice sharing and aliasing — how multiple slices can share the same backing array, and how to avoid data corruption bugs.
