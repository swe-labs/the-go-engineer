# Arrays

## Learning objective

Declare, initialize, iterate over, and compare fixed-size arrays in Go, and explain the value semantics that distinguish `[n]T` from slices.

## Why this matters

Arrays are the foundation of slices, Go's most used data structure. Understanding arrays means you understand contiguous memory layout, value semantics, and why slices exist at all. Many performance-critical Go libraries use arrays directly for stack-allocated buffers and fixed-size lookups. When you see `[32]byte` in a hash or `[16]uint32` in a network packet header, you are looking at an array.

## Mental model

Think of an array as a fixed-size block of identical slots in memory, all with the same type. The size is part of the type: `[5]int` and `[3]int` are different types. When you pass an array to a function or assign it to a new variable, Go copies the entire block — every element. This is value semantics at work.

If you picture a slice as a window into a separate backing array, the array is the building itself. The building has a fixed footprint; you cannot add more floors. The slice is just a view into some of the rooms.

## Core idea

In Go, an array type is written `[n]T` where `n` is a non-negative integer constant (the length) and `T` is the element type. The length must be a compile-time constant expression.

```go
var a [5]int              // array of 5 ints, all zero-valued
var b [3]string           // array of 3 strings, all ""
```

Arrays have two key properties:

1. **Length is part of the type**: `[5]int` and `[10]int` are incompatible types.
2. **Value semantics**: assignment and function calls copy the entire array.

You can access and modify elements with index expressions `a[i]`. Valid indices go from `0` to `n-1`; out-of-range indices cause a compile-time error or runtime panic.

## Under the hood

An array is laid out as `n` contiguous elements of type `T` in memory. The address of element `i` is `baseAddr + i * sizeof(T)`. There is no header, no length field — the length is known at compile time and baked into the type.

The compiler often allocates small arrays on the stack rather than the heap. This makes arrays extremely cheap when they do not escape to the heap.

When you assign an array to another variable, `a := b` compiles to a `memcpy` of `n * sizeof(T)` bytes. For large arrays this can be expensive, which is why Go encourages using slices (reference types) for passing data around.

## How Go uses it

Arrays appear in several standard places:

- **Fixed-size buffers**: `[32]byte` for `sha256.Sum256`, `[64]byte` for `sha512.Sum512`.
- **Network protocols**: `[4]byte` for IPv4 addresses, `[16]byte` for IPv6 or UUIDs.
- **Matrix math**: `[3][3]float64` for 3x3 transformation matrices.
- **Control flow**: `var keypress [256]bool` for debounced key states.
- **Byte slices**: `[16]byte` as a key type in maps (arrays are comparable, slices are not).

Go's `crypto` packages heavily use fixed-size arrays because they need compile-time guarantee of buffer size.

## Go example

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	var a [5]int
	a[0] = 10
	a[1] = 20
	fmt.Println("a:", a)

	// Array literal
	b := [3]string{"go", "rust", "zig"}
	fmt.Println("b:", b)

	// Ellipsis counts the initializers
	c := [...]int{2, 4, 6, 8}
	fmt.Println("c:", c, "len:", len(c))

	// Value semantics: copying
	d := c
	d[0] = 999
	fmt.Println("c after d copy:", c)
	fmt.Println("d:", d)

	// Comparison
	fmt.Println("b == [3]string{\"go\", \"rust\", \"zig\"}:", b == [3]string{"go", "rust", "zig"})

	// Type distinction
	fmt.Println("Type of a:", reflect.TypeOf(a))
	fmt.Println("Type of b:", reflect.TypeOf(b))
	fmt.Println("Type of c:", reflect.TypeOf(c))
}
```

## Step-by-step execution

For the assignment `d := c` where `c = [4]int{2, 4, 6, 8}`:

1. Compiler sees `[4]int` — length `4` is part of the type.
2. At runtime, `c` occupies `4 * 8 = 32` bytes on the stack (on 64-bit systems).
3. `d := c` allocates a new `[4]int` on the stack.
4. `memcpy` copies all 32 bytes from `c` into `d`.
5. `d[0] = 999` modifies only `d`. `c[0]` remains `2`.

For iteration `for i, v := range c`:

1. Compiler knows `len(c) = 4` at compile time — bounds check is elided.
2. Loop runs `i = 0..3`, `v` is a copy of each element.
3. Modifying `v` inside the loop does not modify `c`.

## Common mistakes

- Mistake: Passing a large array to a function expecting a slice copies the whole array.
  - Why it happens: Go passes by value. `func f(arr [1000000]int)` copies one million ints.
  - Fix: Use a slice `func f(arr []int)` or pass a pointer `func f(arr *[1000000]int)`.

- Mistake: Trying to assign `[3]int` to a `[5]int` variable.
  - Why it happens: The length is part of the type. `[3]int` and `[5]int` are different types.
  - Fix: Match the declared length or use slices for variable-sized data.

- Mistake: Using `append` on an array. `append` only works on slices.
  - Why it happens: Arrays cannot grow. You must manually create a larger array and copy.
  - Fix: Use slices if you need dynamic sizing.

- Mistake: Confusing `[n]T` with `[]T` in function signatures.
  - Why it happens: They look similar but behave completely differently (value vs reference).
  - Fix: Use slices as function parameters unless you explicitly need value semantics.

## Debugging walkthrough

Consider this code that prints incorrect results:

```go
package main

import "fmt"

func zero(arr [5]int) {
	for i := range arr {
		arr[i] = 0
	}
}

func main() {
	nums := [5]int{1, 2, 3, 4, 5}
	zero(nums)
	fmt.Println(nums)
}
```

**Symptom**: Prints `[1 2 3 4 5]` even though `zero` sets every element to `0`.

**Investigation**: Add print statements inside `zero`:

```go
func zero(arr [5]int) {
	fmt.Printf("inside zero, arr before: %v, pointer: %p\n", arr, &arr)
	for i := range arr {
		arr[i] = 0
	}
	fmt.Printf("inside zero, arr after: %v\n", arr)
}
```

Output:
```
inside zero, arr before: [1 2 3 4 5], pointer: 0xc0000100a0
inside zero, arr after: [0 0 0 0 0]
[1 2 3 4 5]
```

**Root cause**: Arrays use value semantics. `zero` receives a copy of `nums`. The modifications affect the copy, not the original.

**Fix**: Pass a pointer to the array, or change the function to return a new array, or use a slice.

```go
func zero(arr *[5]int) {
	for i := range arr {
		arr[i] = 0
	}
}
// call: zero(&nums)
```

## Production notes

- **Arrays as map keys**: `[16]byte` is comparable and can be a map key. This is common for caches keyed by hash or UUID.
- **Avoid large arrays on the stack**: Arrays larger than the stack frame may cause a stack overflow. The threshold varies but arrays above a few KB often escape to the heap.
- **Fixed-size lookup tables**: Small arrays like `[256]byte` for ASCII tables are idiomatic and fast — they compile to immediate loads.
- **Never use arrays for dynamic collections**: That is what slices are for. Arrays in public APIs are rare and should signal "fixed size is semantically important".

## Performance implications

- Arrays are stack-allocated when they don't escape, making access as fast as local variables — a single CPU instruction with fixed offset.
- Copying an array of `n` elements is `O(n)` in time and space. A `[1000]int` copy is ~8 KB of `memcpy`. For small arrays (<= 4 machine words), the copy is inlined as register moves.
- Bounds checking is elided when the compiler can prove the index is within `[0, n)`. Iterating `for i := 0; i < len(a); i++` typically has zero overhead.
- Accessing elements through a pointer to an array (`*[n]T`) adds one level of indirection but avoids the copy cost.

## Practice task

Write a function `sumAndAverage(nums [8]float64) (sum float64, avg float64)` that computes the sum and average of the elements. Then write a function `reverse(arr *[6]int)` that reverses the array in place.

In `main()`, demonstrate:
1. Creating an `[8]float64` literal and printing sum and average.
2. Creating a `[6]int` literal, reversing it, and printing the result.
3. Passing a `[6]int` to another function that tries to modify it (to show value semantics).

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/13-arrays
go test ./curriculum/modules/03-programming-fundamentals/lessons/13-arrays
```

## Review questions

1. Why is `[5]int` a different type than `[3]int`? What practical problems does this solve?
2. What happens when you assign one array variable to another? Does the original change?
3. How does `[...]int{1,2,3}` determine its length? Is `...` an operator?
4. Can two arrays of different lengths be compared with `==`? What does the compiler say?
5. Why does `append` not work on arrays? What would you use instead for dynamic growth?

## NEXT UP

Slices — Go's dynamic, flexible view into arrays.
