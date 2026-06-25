# Slices

## Learning objective

Create, manipulate, and reason about slices — Go's dynamically-sized, reference-typed view into a contiguous segment of an underlying array.

## Why this matters

Slices are the single most used data structure in Go. Every Go codebase uses them for lists, buffers, stacks, queues, CSV rows, JSON arrays, and database result sets. Unlike arrays, slices are flexible: they grow, shrink, and are cheap to pass around. Mastering slices is the key to writing idiomatic Go.

## Mental model

A slice is a window into an underlying array. The window has three fields: a pointer to the first visible element, a length (number of visible elements), and a capacity (number of elements available from that pointer to the end of the backing array).

When you `append` beyond the capacity, Go allocates a new, larger backing array, copies all elements over, and returns a new slice header pointing to the new array. The old backing array becomes unreachable if no other slice references it.

## Core idea

A slice is a descriptor of a contiguous segment of an array. The slice type is `[]T` — no length in the type. Slices are created by:

- **Slice literal**: `s := []int{1, 2, 3}` — creates a backing array of length 3.
- **`make`**: `s := make([]int, 5, 10)` — backing array length 10, slice length 5.
- **Slice expression**: `s := arr[1:4]` — views elements 1,2,3 of `arr`.
- **From another slice**: `sub := s[2:5]`.
- **`append`**: `s = append(s, 4)` — returns a new slice with the element appended.

Key properties:
- `len(s)` — number of elements in the slice.
- `cap(s)` — number of elements in the backing array from `s[0]` to the end.
- Slices are reference types: assignment or passing copies the header (ptr, len, cap), not the elements.
- A nil slice (`var s []int`) has length 0 and capacity 0, and no backing array.

## Under the hood

A slice value at runtime is a `reflect.SliceHeader`:

```go
type SliceHeader struct {
	Data uintptr // pointer to the first element of the backing array
	Len  int     // length of the slice
	Cap  int     // capacity of the slice
}
```

This header is 24 bytes on 64-bit systems (8 + 8 + 8). When you pass a slice to a function, Go copies these 24 bytes. The backing array is not copied. This is why slices are cheap to pass around even for large data.

The `append` built-in generates a call into the runtime (`runtime.growslice`) when the capacity is exceeded. `growslice` allocates a new backing array (typically 2x the old capacity for small sizes), copies elements, and returns the new slice.

## How Go uses it

- **Standard library**: `os.ReadFile` returns `[]byte`, `bytes.Buffer` reads/writes `[]byte`, `encoding/json` marshals/unmarshals `[]byte`.
- **Function signatures**: Most public Go APIs accept and return slices, not pointers to arrays.
- **Variadic functions**: `fmt.Printf(format, args...)` — `args` is `[]any`.
- **Range loops**: `for i, v := range slice {}` — iterates length elements.
- **`copy` built-in**: `copy(dst, src)` copies elements between slices with overlap handling.

## Go example

```go
package main

import "fmt"

func main() {
	// Slice literal
	fruits := []string{"apple", "banana", "cherry"}
	fmt.Println("fruits:", fruits, "len:", len(fruits), "cap:", cap(fruits))

	// make with length and capacity
	s := make([]int, 3, 6)
	fmt.Println("s:", s, "len:", len(s), "cap:", cap(s))

	// append
	s = append(s, 10, 20, 30)
	fmt.Println("after append:", s, "len:", len(s), "cap:", cap(s))

	// slice expression
	sub := fruits[1:3]
	fmt.Println("sub:", sub)

	// nil slice
	var nilSlice []int
	fmt.Println("nil slice:", nilSlice, "len:", len(nilSlice), "cap:", cap(nilSlice), "is nil:", nilSlice == nil)

	// empty slice literal
	empty := []int{}
	fmt.Println("empty slice:", empty, "len:", len(empty), "cap:", cap(empty), "is nil:", empty == nil)
}
```

## Step-by-step execution

For `s := make([]int, 2, 4)`:

1. The runtime allocates a backing array of 4 zero-valued `int`s on the heap.
2. A slice header is created: `Data` points to the array's first element, `Len = 2`, `Cap = 4`.
3. The header is assigned to `s`.

For `s = append(s, 100)`:

1. `len(s)` is `2`, `cap(s)` is `4`. `len < cap`, so space is available.
2. `s[2] = 100` writes to the backing array.
3. `Len` is incremented to `3`. No new allocation.

For `s = append(s, 200, 300, 400)`:

1. `len(s)` is `3`, need 3 more elements. New length would be `6`, exceeding `cap(s)=4`.
2. `growslice` allocates a new backing array of capacity `8` (double the old capacity for small slices).
3. Elements `{0, 0, 100, 200, 300, 400}` are copied to the new array.
4. A new header with `Len=6, Cap=8` is returned.

## Common mistakes

- Mistake: Expecting `append` to modify the slice in place.
  - Why it happens: `append` may return a new slice header. If you don't assign the result, the original slice stays unchanged.
  - Fix: Always write `s = append(s, val)`.

- Mistake: Confusing nil and empty slices.
  - Why it happens: Both have `len == 0`, but `var s []int` is nil while `s := []int{}` is non-nil. They behave the same with `append` and `range`, but `json.Marshal(nilSlice)` produces `null` and `json.Marshal(emptySlice)` produces `[]`.
  - Fix: Use `var s []T` for "no slice yet", `make([]T, 0)` for "empty but ready to use".

- Mistake: Reading beyond `len` but within `cap`.
  - Why it happens: `cap` tells you the backing array size, but `len` bounds the accessible elements. `s[:cap(s)]` works because slice literals go up to cap, but `s[3]` panics if `len(s) = 2`.
  - Fix: Only index up to `len(s)-1`.

- Mistake: Using slices as map keys.
  - Why it happens: Slices are not comparable (no `==` operator). This is a compile error.
  - Fix: Convert to string or use arrays as keys.

## Debugging walkthrough

```go
package main

import "fmt"

func main() {
	data := make([]int, 0, 5)
	data = append(data, 1, 2, 3)
	process(data)
	fmt.Println("after process:", data)
}

func process(items []int) {
	items = append(items, 4)
	fmt.Println("inside process:", items)
}
```

**Symptom**: Prints `inside process: [1 2 3 4]` but `after process: [1 2 3]`.

**Investigation**: Add prints of `len` and `cap`:

```go
func process(items []int) {
	fmt.Printf("before append: ptr=%p len=%d cap=%d\n", items, len(items), cap(items))
	items = append(items, 4)
	fmt.Printf("after append: ptr=%p len=%d cap=%d\n", items, len(items), cap(items))
}
```

Output:
```
before append: ptr=0xc0000123a0 len=3 cap=5
after append: ptr=0xc0000123a0 len=4 cap=5
after process: [1 2 3]
```

**Root cause**: The slice header is copied when passed to `process`. The backing array is shared, so `append` within capacity modifies the same backing array. But the new `len=4` exists only in `process`'s copy of the header. The caller's header still has `len=3`.

**Fix**: Return the new slice: `data = process(data)`, or pass a pointer to the slice.

## Production notes

- **Preallocate with `make`**: If you know the final size, `make([]T, 0, expectedN)` and `append` avoids reallocation.
- **Slice reuse**: Avoid repeated `append` in hot loops by resetting `s = s[:0]` on a pre-allocated slice.
- **JSON**: `null` vs `[]` matters. For JSON APIs, prefer `make([]T, 0)` over `var s []T` to get `[]` instead of `null`.
- **Don't store slices in structs unless necessary**: A slice in a struct ties the struct to heap-allocated memory and complicates ownership. Prefer owning the backing array with an array when appropriate.

## Performance implications

- `append` is amortized O(1) per element. The growth factor (2x for small slices, 1.25x for large) ensures that the average cost per append is constant.
- Access is O(1) with a single bounds check. The compiler may elide the check when it can prove safety.
- Slicing (`s[low:high]`) is O(1) — it just copies the header. No data is moved.
- `copy()` is O(n) and compiles to `memmove`.
- Nil slices cost nothing: the zero value of `[]T` is nil, ready to use with `append`.

## Practice task

Write a function `merge(a, b []int) []int` that returns a new slice containing all elements of `a` followed by all elements of `b`. Then write a function `dedup(sorted []int) []int` that returns a new slice with consecutive duplicates removed. In `main()`, demonstrate both with sample data and print results.

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/14-slices
go test ./curriculum/modules/03-programming-fundamentals/lessons/14-slices
```

## Review questions

1. What is the difference between `var s []int` and `s := []int{}`? When does it matter?
2. After `s = append(s, x)`, why must you assign the result back to `s`?
3. A slice header is 24 bytes. What three fields does it contain?
4. Can you compare two slices with `==`? Why or why not?
5. What does `make([]int, 5)` allocate? What is `len` and `cap` of the result?

## NEXT UP

Slice length and capacity — understanding the difference and how they affect growth, bounds, and performance.
