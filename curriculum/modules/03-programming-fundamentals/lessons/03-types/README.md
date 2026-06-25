# Types

## Learning objective

Identify and use Go's basic types (int, float64, string, bool, byte, rune), distinguish typed from untyped constants, and declare custom type aliases and definitions.

## Why this matters

Every value in Go has a type. The type system catches mismatches at compile time — before your code ever runs. A function that expects a `float64` cannot silently receive an `int`; a string concatenation with a numeric type is a compile error, not a runtime surprise. Understanding the type lattice — what the basic types are, which are aliases, and how type inference works — is the foundation for writing correct Go code. In production, type mismatches are a top source of subtle bugs in cross-team codebases; mastering types eliminates an entire class of defects.

## Mental model

Think of a type as a mold or a stencil. The mold determines three things: (1) the set of values that fit, (2) the operations you can perform, and (3) the memory layout. An `int8` mold accepts values -128 to 127 and uses 1 byte. A `float64` mold accepts IEEE 754 doubles and uses 8 bytes. When you declare a variable `var x int32`, you've chosen the `int32` mold — Go will never let you pour a `string` into it without an explicit conversion.

## Core idea

Go's basic types fall into several families:

| Category | Types | Zero value |
|---|---|---|
| Boolean | `bool` | `false` |
| Signed integer | `int8`, `int16`, `int32`, `int64`, `int` | `0` |
| Unsigned integer | `uint8`, `uint16`, `uint32`, `uint64`, `uint` | `0` |
| Float | `float32`, `float64` | `0.0` |
| Complex | `complex64`, `complex128` | `0+0i` |
| Byte (alias) | `byte` (= `uint8`) | `0` |
| Rune (alias) | `rune` (= `int32`) | `0` |
| String | `string` | `""` |

**Type inference with `:=`** deduces the type from the value:

```go
x := 42       // int (platform-dependent: 32 or 64 bit)
y := 3.14     // float64
z := "hello"  // string
ok := true    // bool
```

**Typed vs untyped**: Untyped constants (like `42` or `3.14` in source code) have no fixed type until they are used. Typed constants (like `const x int = 42`) have a fixed type and cannot be used where that type is required without conversion.

**Type aliases** create an alternate name for an existing type. `type byte = uint8` means `byte` and `uint8` are interchangeable — both names refer to the same type.

**Type definitions** create a new distinct type. `type Celsius float64` creates a brand new type `Celsius` that is based on `float64` but is distinct — you cannot pass a `float64` where `Celsius` is expected without conversion.

## Under the hood

The compiler stores a type descriptor for every type in the program. For basic types, these descriptors are built into the compiler. For user-defined types (including aliases and definitions), the compiler creates new descriptors. Type aliases are represented as direct references to the aliased type's descriptor — they are literally the same type at the compiler level. Type definitions create a new descriptor with a new name but the same underlying structure.

At the machine-code level, there is no difference between `int32` and `rune` — both use 4 bytes and the same instructions. The distinction exists only at compile time for type safety.

## How Go uses it

- **`int` and `uint`** are the default numeric types. Their size is platform-dependent: 32 bits on 32-bit platforms, 64 bits on 64-bit platforms. Use `int` unless you need a specific size.
- **`byte`** is used for raw binary data (like `[]byte`). It appears in I/O, networking, and string operations.
- **`rune`** is used for Unicode code points. Every character in a Go string is a `rune`. String iteration yields `rune` values.
- **`float64`** is the default floating-point type. `float32` is used for memory-constrained or GPU-compute scenarios.
- **Type inference** (`:=`) is idiomatic. The Go team explicitly designed it so that common types are chosen sensibly: `42` → `int`, `3.14` → `float64`, `"hi"` → `string`, `true` → `bool`.
- **Type definitions** are used to create domain-specific types that prevent mixing units (e.g., `type Meter float64` and `type Feet float64`).

## Go example

```go
package main

import "fmt"

type Celsius float64
type Fahrenheit float64

func main() {
	// basic types
	var a int = 42
	var b float64 = 3.14
	var c string = "hello"
	var d bool = true

	fmt.Printf("int: %d, float64: %.2f, string: %s, bool: %t\n", a, b, c, d)

	// type inference
	x := 42
	y := 3.14
	z := "hello"
	ok := true
	fmt.Printf("inferred: %T %T %T %T\n", x, y, z, ok)

	// byte and rune
	var char byte = 'A'
	var code rune = '世'
	fmt.Printf("byte: %c (%d), rune: %c (U+%04x)\n", char, char, code, code)

	// type definition creates a distinct type
	var temp Celsius = 100.0
	fmt.Printf("Celsius: %.1f°C\n", temp)

	// type alias: byte and uint8 are the same
	var b1 byte = 255
	var b2 uint8 = b1 // no conversion needed
	fmt.Printf("byte == uint8: %d == %d\n", b1, b2)
}
```

## Step-by-step execution

For `var temp Celsius = 100.0`:

1. Compiler sees the type definition `type Celsius float64`. It registers `Celsius` as a new named type with underlying type `float64`.
2. Compiler allocates 8 bytes for `temp` on the stack.
3. The value `100.0` is converted from untyped float to `Celsius` (which is represented as `float64` bits).
4. At runtime, `temp` holds the IEEE 754 representation of `100.0`.

For `x := 42`:

1. Compiler sees the untyped integer literal `42`.
2. It applies default type rules: untyped integer → `int`.
3. `x` is declared as `int` with value `42`.

## Common mistakes

- Mistake: `var x int = 3.14` expecting implicit truncation.
  - Why: Go never implicitly converts between numeric types. This is a compile error.
  - Fix: `var x int = int(3.14)` (truncates to 3), or use a float type.

- Mistake: Confusing type aliases with type definitions.
  - Why: `type myInt = int` means `myInt` and `int` are interchangeable. `type myInt int` creates a distinct type.
  - Fix: Understand that `=` after the type name means alias, not new type.

- Mistake: Passing a typed constant where another numeric type is expected.
  - Why: `const x int = 5` is typed `int`. Passing to a `float64` parameter fails at compile time.
  - Fix: Use untyped constants when flexibility is needed: `const x = 5`.

## Debugging walkthrough

Consider this broken code:

```go
package main

import "fmt"

type Score int

func main() {
	highScore := 100
	printScore(highScore)
}

func printScore(s Score) {
	fmt.Println("Score:", s)
}
```

**Symptom**: Compile error: `cannot use highScore (variable of type int) as type Score in argument to printScore`.

**Investigation**: `type Score int` is a type definition, not an alias. `Score` and `int` are distinct types. `highScore` is `int` (inferred from `100`), but `printScore` expects `Score`.

**Root cause**: `Score` is a distinct type. Go requires explicit conversion between distinct types even if they have the same underlying type.

**Fix**: Convert explicitly: `printScore(Score(highScore))`.

Alternatively, change the definition to an alias: `type Score = int` (then no conversion needed, but you lose type safety).

## Production notes

- **Use type definitions for domain modeling** — `type UserID int64` and `type OrderID int64` prevent passing a user ID where an order ID is expected. This catches bugs at compile time.
- **Prefer `int` and `float64`** for general use. Only use specific-size types when interfacing with binary protocols, fixed-size structs, or memory-constrained environments.
- **`byte` vs `uint8`** — use `byte` when the data represents bytes (binary I/O, strings), not for general math.
- **`rune` vs `int32`** — use `rune` when the value is a Unicode code point. This communicates intent.
- **Type aliases** are rare. They are mostly used for gradual code migration (e.g., when moving a type between packages).

## Performance implications

- All basic types map directly to hardware registers. There is no boxing or wrapping overhead.
- `int` is the most efficient integer type because it matches the native word size. Smaller types (`int8`, `int16`) may require extra instructions for sign/zero extension.
- `float64` is faster than `float32` on modern x86 CPUs (SSE instructions natively operate on 64-bit floats; converting to/from 32-bit adds overhead).
- Type definitions have zero runtime cost — they are erased by the compiler, leaving only the underlying type's machine code.

## Practice task

Write a function `convertUnits(value float64, from, to string) (float64, error)` that converts between temperature scales: Celsius (`"C"`), Fahrenheit (`"F"`), and Kelvin (`"K"`). Formulas:
- F = C × 9/5 + 32
- K = C + 273.15
- C = (F − 32) × 5/9
- C = K − 273.15

Return an error for unknown unit strings. In `main()`, test at least 5 conversions and print results. Use type inference with `:=` throughout.

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/03-types
go test ./curriculum/modules/03-programming-fundamentals/lessons/03-types
```

## Review questions

1. Explain the difference between `type Score int` and `type Score = int`. When would you use each?
2. Given `const x = 100` and `const y int = 100`, can you pass `x` and `y` to a function expecting `float64`? Explain.
3. Debug: `var r rune = 'A'; var b byte = r` — does this compile? Why or why not? How would you fix it?
4. Tradeoff: Why does Go use `int` (platform-dependent size) instead of fixing it to `int64`? What are the tradeoffs on a 32-bit platform?
5. What type does `x := 'A'` infer? (Hint: character literals in Go have a default type.)

## NEXT UP

Zero values — the guaranteed default value for every Go type and why this design choice prevents an entire category of bugs.
