# Type conversions

## Learning objective

Convert between Go types explicitly using the `T(x)` conversion syntax, numeric conversions between integer and float types, string/byte slice interconversion, the `strconv` package for parsing and formatting, and recognize when unsafe conversions are necessary or dangerous.

## Why this matters

Go has no implicit type coercion. Every type change must be spelled out. This is a feature, not a bug: it makes data flow visible at a glance and prevents entire classes of silent data corruption bugs. In practice, reading user input, serializing to JSON, working with protocol buffers, and performing arithmetic between different numeric types all demand type conversions. Knowing the rules — and the pitfalls — separates reliable code from fragile code.

## Mental model

Think of a Go value as a box labelled with a type. The box contains bits. A type conversion `T(x)` is a machine that takes box `x`, reads its bits, and produces a new box labelled `T` containing the closest representable value. The conversion may:

- Change the bit pattern (e.g., `float64` to `int` truncates).
- Keep the same bits but change the label (e.g., `string` to `[]byte` reinterprets the memory).
- Fail at compile time if the types are incompatible (e.g., `string` to `int`).

Conversions never silently lose information unless you ask for a narrower type. When the source value cannot be represented in the target type, integer conversions wrap around (two's complement truncation); float-to-int conversions truncate toward zero.

## Core idea

The syntax `T(x)` converts expression `x` to type `T`. Both `T` and `x` must have the same _underlying type_ or be numeric kinds, strings, or slices with compatible element types. The conversion is a compile-time checked operation that may produce a runtime value.

```go
var i int = 42
var f float64 = float64(i)  // explicit conversion
var b byte = byte(i)         // narrower type: truncates if i > 255
```

Conversions fall into families:

| Family | Example | Behaviour |
|---|---|---|
| Numeric | `int(3.14)`, `float64(7)` | Arithmetic representation change |
| String ↔ byte slice | `string([]byte{65})` → `"A"` | Reinterpret memory as UTF-8 |
| String ↔ rune slice | `string([]rune{8226})` → `"•"` | Encode/decode runes |
| `strconv` | `strconv.Atoi("42")` | Parse or format to/from string |
| `unsafe` | `(*int)(unsafe.Pointer(&f))` | Reinterpret bit pattern (dangerous) |

## Under the hood

The Go compiler represents every value as a sequence of bytes in memory. A conversion from `int32` to `int64` is a sign-extending move instruction at the assembly level. A conversion from `float64` to `int` calls a `cvttsd2si` instruction on x86-64 that truncates toward zero. A conversion from `string` to `[]byte` allocates a new byte slice and copies the string's backing memory — the string's immutability guarantee is preserved.

For `strconv` functions, the runtime parses ASCII digits character by character, handling sign, base prefixes, and overflow detection. `strconv.Itoa` builds digits from least significant to most, then reverses.

The `unsafe` package bypasses the type system entirely. `unsafe.Pointer` is the Go equivalent of `void*` in C. Conversions through `unsafe` are not portable, not guaranteed by the Go spec, and can crash the runtime if misused.

## How Go uses it

Every non-trivial Go codebase contains conversions:

- **Reading numbers from strings**: `strconv.Atoi`, `strconv.ParseFloat` appear in CLI flags, config files, HTTP query parameters, and CSV parsing.
- **JSON encoding**: `json.Marshal` internally converts Go types to their JSON representations; fields typed as `float64` in decoded JSON often need conversion to `int`.
- **Protocol buffers**: Generated `proto` types use `int32`, `int64`, `uint32` etc.; business logic converts to platform `int`.
- **Byte-level I/O**: `[]byte` is the currency of network and file I/O; `string(buf)` converts received bytes to strings for processing.
- **Unsafe optimizations**: High-performance parsers sometimes use `unsafe` to cast between `[]byte` and `string` without allocation.

## Go example

```go
package main

import (
	"fmt"
	"math"
	"strconv"
)

func main() {
	var i int = 255
	var f float64 = float64(i)
	fmt.Printf("int → float64: %d → %f\n", i, f)

	var trunc int = int(3.99)
	fmt.Printf("float64 → int (truncate): %d\n", trunc)

	s := strconv.Itoa(42)
	n, _ := strconv.Atoi("42")
	fmt.Printf("strconv.Itoa(42) = %q, strconv.Atoi(\"42\") = %d\n", s, n)

	gpa, _ := strconv.ParseFloat("3.75", 64)
	fmt.Printf("ParseFloat: %f\n", gpa)

	formatted := strconv.FormatFloat(math.Pi, 'f', 4, 64)
	fmt.Printf("FormatFloat(Pi, 4 digits): %s\n", formatted)

	b := []byte("Go")
	s2 := string(b)
	fmt.Printf("[]byte → string: %q\n", s2)

	large := math.MaxUint16
	narrow := byte(large)
	fmt.Printf("uint16(%d) → byte wraps to: %d\n", large, narrow)
}
```

## Step-by-step execution

Trace `strconv.Atoi(" -42")`:

1. `Atoi` calls `ParseInt(s, 10, 0)` with base 10 and bit size 0 (maps to `int`).
2. `ParseInt` skips leading whitespace.
3. Reads the `-` sign, sets negative flag.
4. Accumulates digits: `'4'` → `4`, `'2'` → `4*10 + 2 = 42`.
5. Applies sign: `-42`.
6. Checks overflow against `math.MinInt`/`math.MaxInt` for 64-bit.
7. Returns `-42, nil`.

Trace `byte(256)`:

1. Constant `256` has type `int` with binary `00000001 00000000`.
2. Target type `byte` (uint8) can hold only the lower 8 bits.
3. Compiler truncates: `00000000` → `0`.
4. This is compile-time constant conversion; the compiler warns or treats overflow as a compile error if the source is an untyped constant.

## Common mistakes

- **Expecting `int(3.14)` to round**: It truncates toward zero. `int(3.99)` is `3`, `int(-3.99)` is `-3`. Use `math.Round`, `math.Floor`, or `math.Ceil` for explicit rounding.

- **Confusing `string(65)` with `strconv.Itoa(65)`**: `string(65)` yields `"A"` (the rune with code point 65), not `"65"`. Always use `strconv.Itoa` for decimal string representation.

- **Forgetting `strconv.Atoi` returns two values**: The error must be handled. Ignoring it with `_` silently produces zero on bad input.

- **Assuming `float64(int)` is lossless for all ints**: `int64` values with magnitude > 2^53 lose precision when converted to `float64` because `float64` has only 52 mantissa bits. `float64(9_007_199_254_740_993)` rounds.

- **Using `unsafe` conversions carelessly**: `unsafe.Pointer` conversions that violate memory alignment or aliasing rules can cause crashes that are nearly impossible to debug.

- **Integer overflow during conversion**: `int8(128)` compiles, but `var x int = 128; int8(x)` silently wraps to `-128`. There is no runtime overflow check.

## Debugging walkthrough

Consider this bug:

```go
package main

import (
	"fmt"
	"strconv"
)

func main() {
	price := "19.99"
	cents := int(strconv.Atoi(price)) // compile error!
	fmt.Println(cents)
}
```

**Symptom**: Does not compile: "cannot convert string to type int" and "multiple-value strconv.Atoi in single-value context".

**Fix step by step**:
1. `strconv.Atoi` returns `(int, error)` and only works on integer strings. For `"19.99"`, use `strconv.ParseFloat`.
2. `strconv.ParseFloat` returns `(float64, error)`.
3. Convert `float64` to cents: multiply by 100, then convert to `int`.

```go
f, err := strconv.ParseFloat(price, 64)
if err != nil {
	panic(err)
}
cents := int(math.Round(f * 100))
fmt.Println(cents) // 1999
```

Now consider a silent truncation bug:

```go
var big uint64 = 1_000_000_000_000
small := int32(big) // silent: wraps!
fmt.Println(small)
```

**Investigation**: Print the hex values. `fmt.Printf("big: %X\n", big)` shows `E8D4A51000`. `int32` keeps only the low 32 bits: `0x4A51000 = 1215752192`.

**Fix**: Check the range explicitly or use `math/bits` utilities.

## Production notes

- **Validate before converting**: Always check `strconv` errors. Never use `_` to discard the error.
- **Prefer `strconv` over `fmt.Sprint`**: `fmt.Sprintf("%d", n)` allocates and is ~3x slower than `strconv.Itoa(n)`.
- **Use named types for safety**: Instead of passing raw `int` IDs everywhere, define `type UserID int` and convert explicitly at boundaries.
- **Watch for 32-bit platforms**: `int` is 32 bits on 32-bit architectures. Code that assumes `int` is 64 bits will break. Use `int64` explicitly for portable 64-bit storage.
- **`unsafe` in production**: Only use `unsafe` when profiling proves a conversion bottleneck and you have a benchmark to validate the improvement. Document every unsafe block with the safety invariant it relies on.

## Performance implications

- Numeric conversions between same-width types (e.g., `int32` → `uint32`) are free: they change only the type label, not the bits.
- Widening conversions (`int32` → `int64`) are a single register move with sign extension.
- Narrowing conversions (`int64` → `int32`) are a single register move; the high bits are discarded.
- `string` ↔ `[]byte` conversions allocate and copy. For zero-copy, use `unsafe` (string ⇒ `[]byte`) or `reflect` (slice ⇒ string), but only after profiling.
- `strconv.Atoi` on a 10-digit string takes ~50-100ns. `strconv.ParseFloat` takes ~100-200ns. Both are allocation-free for most inputs.
- `strconv.FormatFloat` allocates — the result is a new string. Pre-allocate buffers with `appendFloat` when formatting many floats in a loop.

## Practice task

Write a function `parseCelsius(input string) (float64, error)` that:
- Accepts a temperature string like `"23.5"`, `"-10"`, or `"37.2C"` (optional trailing `C`).
- Parses the numeric portion using `strconv.ParseFloat`.
- Returns the temperature as a `float64`.
- Returns an error if parsing fails or if any suffix is present that is not `"C"`.

Then write a function `fahrenheitToCelsius(f float64) float64` that converts Fahrenheit to Celsius.

In `main()`, read a command-line argument, parse it, convert Fahrenheit to Celsius if the input ends with `"F"`, and print the result.

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/05-type-conversions
go test ./curriculum/modules/03-programming-fundamentals/lessons/05-type-conversions
```

## Review questions

1. What is the difference between `string(65)` and `strconv.Itoa(65)`?
2. Why does `int(3.99)` produce `3` and not `4`? How would you get `4`?
3. What happens when you convert `math.MaxInt64` to `float64` and back to `int64`? Is the value preserved?
4. Explain the output of `var x int8 = 127; fmt.Println(int8(x + 1))`.
5. When would you use `unsafe.Pointer` for type conversion, and what risks does it carry?

## NEXT UP

Constants and how the `const` keyword creates compile-time-evaluated named values with arbitrary precision.
