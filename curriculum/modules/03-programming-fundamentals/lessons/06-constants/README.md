# Constants

## Learning objective

Declare typed and untyped constants using `const`, understand compile-time evaluation and arbitrary-precision numeric constants, and use `iota` for enumerated constant blocks.

## Why this matters

Constants make code self-documenting, prevent magic numbers from scattering through a codebase, and enable the compiler to perform optimizations unavailable to variables. Unlike most languages, Go constants are _untyped_ by default — they carry arbitrary precision until forced into a concrete type. This design allows clean arithmetic in constant expressions without overflow, then exact conversion when needed. Mastering constants is essential for API design, bitmask flags, sentinel errors, and configuration defaults.

## Mental model

A `const` declaration is a promise: this value never changes. The Go compiler treats every `const` as a literal. When you write `const secondsPerHour = 3600`, the compiler substitutes `3600` everywhere `secondsPerHour` appears. There is no runtime storage — no memory is allocated. Untyped constants are like mathematical ideals: they live in a platonic realm of infinite precision. Only when assigned to a variable or used in a typed context do they acquire a concrete type and, if necessary, range checking.

Typed constants are different: `const x int = 10` has the type `int` baked in. You cannot assign `x` to a `float64` variable without an explicit conversion.

## Core idea

```go
const Pi = 3.141592653589793238462643383279502884197  // untyped, arbitrary precision
const Pi32 float32 = 3.141592653589793238462643383279502884197  // typed, narrowed to float32
```

Key properties:

| Property | Untyped constant | Typed constant |
|---|---|---|
| Declared with | `const X = val` | `const X T = val` |
| Has a type | No (default kind: bool, rune, int, float64, complex128, string) | Yes, concrete type `T` |
| Precision | Arbitrary (limited by compiler, typically 256-bit) | Fixed by type |
| Assignment | Assignable to any compatible type | Only to matching type |
| Memory | Zero (compile-time substitution) | Zero (compile-time substitution) |

Constant expressions are evaluated at compile time using the compiler's big-number arithmetic. You can write `const huge = 1 << 100` — it compiles. But `var huge = 1 << 100` fails because `1` is an `int`, and shifting 100 bits overflows.

## Under the hood

The Go compiler represents untyped constants internally as `*big.Int`, `*big.Rat`, or `*big.Float` structures. Constant expressions are evaluated using these arbitrary-precision types, giving exact results for integer arithmetic and configurable precision (default 256 bits) for floating-point.

When a constant is used in a typed context (assigned to a variable, passed to a function), the compiler checks that the value fits in the target type. For integers: the value must be representable. For floating-point: the value is rounded to the target precision. If overflow occurs, the compiler emits a compile-time error — not a runtime wrap.

The `iota` identifier is a special compiler-managed counter that resets to `0` at each `const` block and increments by `1` for each `ConstSpec` within that block.

## How Go uses it

- **Sentinel errors**: `const ErrNotFound = errors.New("not found")` — though as of Go 1.13, convention uses `var` for sentinel errors to support `errors.Is`.
- **Bitmask flags**: `const ( Read = 1 << iota; Write; Execute )`.
- **Math constants**: `math.Pi`, `math.E` are untyped float constants.
- **Configuration defaults**: `const DefaultTimeout = 30 * time.Second`.
- **Enum-like types**: `type Color int` with a `const` block assigning `iota`.
- **Compile-time size checks**: `const _ = unsafe.Sizeof(...)` used in `init` guards.

## Go example

```go
package main

import "fmt"

const (
	UnixToGoOffset = 62135596800 // untyped int constant
	SecondsPerDay  = 86400
)

const typedHello string = "hello"

func main() {
	const localConst = 42
	fmt.Println("localConst:", localConst)

	const bigShift = 1 << 100
	fmt.Printf("1<<100 has type %T\n", bigShift) // still untyped

	// Constrain to a concrete type:
	const normalized = bigShift >> 97
	fmt.Println("normalized (1<<3):", normalized)

	// Typed constant requires explicit conversion
	var s string = typedHello
	fmt.Println("s:", s)

	// Untyped constant assigned to different types
	var i int = SecondsPerDay
	var f float64 = SecondsPerDay
	var d time.Duration = SecondsPerDay * time.Second
	fmt.Println(i, f, d)
}
```

## Step-by-step execution

Trace `const x = 1 << 100 / 1e9`:

1. Compiler reads `const x`. Expression must be a constant expression — all operands are literals or other constants.
2. `1` is an untyped integer constant.
3. `1 << 100` is evaluated using `big.Int`: result is `1267650600228229401496703205376`.
4. `1e9` (`1000000000`) is an untyped float constant. Mixed integer ÷ float promotes to float: `1267650600228229401496703205376.0 / 1000000000.0`.
5. Result: `1.2676506002282294e+30` (rounded to 256-bit precision).
6. `x` is an untyped float constant with value `~1.26765e30`.

Trace `const y int = 1 << 100`:

1. Typed constant `int`. At compile time, the value `1267650600228229401496703205376` is checked against `int` range.
2. On 64-bit platforms, `int` max is `9223372036854775807`.
3. Value exceeds range → compile error: `constant 1267650600228229401496703205376 overflows int`.

## Common mistakes

- **Confusing `const` with `var`**: `const x = math.Sqrt(4)` — function calls are not constant expressions. Only built-in functions like `len`, `cap`, `real`, `imag`, and `complex` work in constant expressions.

- **Assuming `const` prevents slice/ map mutation**: `const s = []int{1,2,3}` does not compile. Slices, maps, and functions cannot be declared `const`. Use `var` for these.

- **Forgetting `iota` resets per `const` block**: `iota` does not increment across files or across `const` blocks.

- **Overflow in typed constant**: `const x int8 = 200` compiles (constant overflow error), but `const x = 200; var y int8 = x` also fails because the conversion is checked at compile time.

- **Thinking untyped constants have dynamic type**: `const x = 42` does not have a type at compile time, but when used in `var y int = x`, `x` is assignable to `int`. It is not assignable to `*int` — there is no such thing as a constant pointer.

## Debugging walkthrough

This code does not compile:

```go
package main

import (
	"fmt"
	"math"
)

const x = math.Pow(2, 10)

func main() {
	fmt.Println(x)
}
```

**Symptom**: `math.Pow(2, 10) (value of type float64) is not constant`.

**Root cause**: `math.Pow` is a function call. Only constant expressions — literals, operators, and a few built-in functions — are allowed in `const` declarations.

**Fix**: Use a constant expression:

```go
const x = 1 << 10  // 1024
```

Or use a variable:

```go
var x = math.Pow(2, 10)
```

Now consider:

```go
const a = 1 << 100
var b int = a
```

**Symptom**: `constant 1267650600228229401496703205376 overflows int`.

**Fix**: Use a narrower shift or assign to a type that can hold the value:

```go
const a = 1 << 100
var b float64 = a  // OK: float64 can hold large values (with precision loss)
```

## Production notes

- **Prefer `const` over `var` for unchanging values**: Constants aid readability, enable compiler optimizations, and prevent accidental mutation. The compiler can fold constant expressions at compile time.
- **Use typed constants for exported API values**: `const DefaultPort uint16 = 8080` documents the expected type in godoc and prevents misuse.
- **Use untyped constants for internal arithmetic**: Allows the constant to flow into any compatible numeric type without explicit conversion.
- **Constants in switch statements**: The compiler can optimize constant-only switch statements into jump tables.
- **`iota` with bitmask patterns**: Document the pattern explicitly. Readers unfamiliar with `iota` often misinterpret `1 << iota` sequences.

## Performance implications

- Constants have zero runtime cost — they are inlined at compile time. There is no memory load, no register spill.
- Constant expressions are evaluated once at compile time, not repeated at runtime.
- A switch on a constant value is compiled into a jump table (O(1)), while a switch on a variable uses a binary search or comparison chain.
- Using `const` instead of `var` for `string` values saves a global string allocation — the constant is emitted directly in the data section.
- Array sizes declared with `const` are known at compile time; the compiler can perform bounds-check elimination.

## Practice task

Write a program that:

1. Defines typed constants for `DaysInWeek = 7`, `HoursInDay = 24`, `MinutesInHour = 60`, `SecondsInMinute = 60`.
2. Defines an untyped constant `SecondsPerWeek` using a constant expression from the above constants.
3. Prints `SecondsPerWeek`.
4. Defines a `const` block with `iota` for weekdays: `Monday = iota` through `Sunday`.
5. Prints the numeric value of `Saturday`.

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/06-constants
go test ./curriculum/modules/03-programming-fundamentals/lessons/06-constants
```

## Review questions

1. Why does `const x = math.Sqrt(16)` fail to compile, but `const x = 4` works?
2. What is the difference between `const x = 42` and `const x int = 42`?
3. Can you declare a `const` slice? Why or why not?
4. What is the result of `const x = 1 << 100; var y int64 = x` and why?
5. How does `iota` behave across multiple `const` blocks in the same package?

## NEXT UP

iota — the constant generator for enumerated values and bitmask patterns in Go.
