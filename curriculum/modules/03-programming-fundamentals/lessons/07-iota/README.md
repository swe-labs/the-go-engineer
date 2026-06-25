# iota

## Learning objective

Use `iota` to generate enumerated constants, reset iota per const block, skip values with `_`, create bitmask patterns with `<<` and `|`, and type iota constant blocks.

## Why this matters

Enumerated constants appear in every non-trivial codebase: HTTP status codes, log levels, permission bits, protocol tags, and state machines. Manually numbering them invites typos, merge conflicts, and gaps. `iota` tells the compiler to do the numbering automatically. Combined with bit shifting and masking, `iota` produces readable, maintainable flag definitions that are otherwise error-prone to write by hand. Senior Go engineers reach for `iota` whenever a group of related constants must be distinct.

## Mental model

`iota` is a predeclared identifier — like `true` and `false` — that represents successive integer constants within a `const` block. Think of it as a cursor that starts at 0 before the first `ConstSpec` and advances by 1 for each subsequent line. Every `const` block gets its own independent cursor; opening a new `const` block resets `iota` to 0.

```
const block start → iota = 0
  First line     → iota = 0 (used)
  Second line    → iota = 1
  Third line     → iota = 2
  ...
const block end

new const block  → iota = 0 (reset)
```

When combined with operators (`<<`, `|`, `+`), `iota` generates values that follow a pattern rather than a flat sequence.

## Core idea

```go
const (
	A = iota  // 0
	B         // 1
	C         // 2
)
```

Each `iota` increments by one for each line in the `const` block. You can skip values with `_`:

```go
const (
	_  = iota       // 0 — discarded
	KB             // 1
	MB             // 2
	GB             // 3
)
```

You can apply arithmetic:

```go
const (
	_   = 1 << (10 * iota)  // 0: 1 << 0 = 1 (ignored)
	KiB                     // 1: 1 << 10 = 1024
	MiB                     // 2: 1 << 20 = 1048576
	GiB                     // 3: 1 << 30 = 1073741824
)
```

For bitmask flags:

```go
const (
	Read  = 1 << iota  // 1
	Write              // 2
	Execute            // 4
)
```

The pattern `1 << iota` produces powers of two, which can be combined with `|` for compound permissions.

## Under the hood

`iota` is not a variable — it is a compiler intrinsic. During parsing, the compiler maintains a per-block counter. Each time it encounters a `ConstSpec` (a single constant declaration or a group within a parenthesized `const` block), it evaluates the current `iota` value, substitutes it into the expression, and then increments the counter.

The expression using `iota` is evaluated at compile time using the same constant-expression rules as any `const` declaration. The result is a constant with arbitrary precision until typed.

The Go spec defines: "Within a constant declaration, the predeclared identifier `iota` represents successive untyped integer constants. Its value is the index of the respective `ConstSpec` in that constant declaration, starting at 0."

Importantly, `iota` is scoped to the `const` block, not the file or package. Multiple `const` blocks, even adjacent, each reset `iota` to 0.

## How Go uses it

- **Byte size constants**: `KiB`, `MiB`, `GiB` as powers of 1024.
- **Permission flags**: Unix-style `Read`, `Write`, `Execute` bitmasks.
- **Protocol enumerations**: HTTP methods, gRPC status codes, DNS record types.
- **State machines**: `StateIdle`, `StateConnecting`, `StateConnected`.
- **Priority levels**: `LogDebug`, `LogInfo`, `LogWarn`, `LogError`.
- **Hardware registers**: Bit positions for device control registers.

The standard library uses `iota` extensively: `time` package constants (`Nanosecond`, `Microsecond`, etc.), `net` (`FlagUp`, `FlagBroadcast`, etc.), and `os` (`O_RDONLY`, `O_WRONLY`, `O_RDWR`).

## Go example

```go
package main

import "fmt"

type Permission uint8

const (
	Read Permission = 1 << iota
	Write
	Execute
)

const (
	_  = 1 << (10 * iota)
	KiB
	MiB
	GiB
	TiB
	PiB
)

type Weekday int

const (
	Monday Weekday = iota + 1
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func main() {
	fmt.Printf("Read=%b, Write=%b, Execute=%b\n", Read, Write, Execute)

	perm := Read | Write
	fmt.Printf("Read|Write = %b, can read? %t, can exec? %t\n",
		perm, perm&Read != 0, perm&Execute != 0)

	fmt.Printf("1 MiB = %d bytes\n", MiB)
	fmt.Printf("1 GiB = %d bytes\n", GiB)

	fmt.Printf("Wednesday = %d\n", Wednesday)
}

func hasPermission(perm, check Permission) bool {
	return perm&check != 0
}
```

## Step-by-step execution

Trace the `Permission` block:

1. Compiler enters `const` block. `iota = 0`.
2. `Read = 1 << iota` → `1 << 0` → `1` (binary `001`).
3. `iota` increments to `1`.
4. `Write` inherits the expression `1 << iota` → `1 << 1` → `2` (binary `010`).
5. `iota` increments to `2`.
6. `Execute` → `1 << 2` → `4` (binary `100`).

Trace `perm := Read | Write`:

1. `Read` = `001`, `Write` = `010`.
2. Bitwise OR: `001 | 010 = 011` (decimal 3).
3. `perm & Read` → `011 & 001 = 001` → nonzero → `true`.
4. `perm & Execute` → `011 & 100 = 000` → zero → `false`.

## Common mistakes

- **Forgetting `iota` resets per block**: `iota` in a second `const` block starts at 0 again. Never assume `iota` continues across blocks.

- **Mixing `iota` with explicit values**: `const ( A = iota; B = 42; C )` — `C` gets `43` because `iota` still increments even when you assign an explicit value to `B`.

- **Using `iota` outside a `const` block**: This is a compile error. `iota` is only valid in `const` declarations.

- **Reordering constants**: Adding a constant in the middle of a block shifts all subsequent values. If you need stable values (e.g., for wire protocols), assign explicit numbers and document the schema.

- **Assuming `iota` starts at 1**: `iota` always starts at 0. If you want 1-based, use `iota + 1` or `_ = iota` to skip 0.

## Debugging walkthrough

```go
package main

import "fmt"

const (
	A = iota
	B
	C
)

const (
	X = iota
	Y
	Z
)

func main() {
	fmt.Println(A, B, C) // ?
	fmt.Println(X, Y, Z) // ?
}
```

**Prediction**: Both blocks start at 0 because `iota` resets per block. Output: `0 1 2` and `0 1 2`. Many newcomers expect `X` to be `3`.

Now consider:

```go
const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
)

func main() {
	fmt.Println(KB, MB, GB)
}
```

**Trace**: `iota` starts at 0. `_ = iota → 0` (discarded). iota = 1: `KB = 1 << (10 * 1) = 1024`. iota = 2: `MB = 1 << (10 * 2) = 1048576`. iota = 3: `GB = 1 << (10 * 3) = 1073741824`.

**Output**: `1024 1048576 1073741824`.

## Production notes

- **Use `iota` for internally consistent groups**: When constants only need to be distinct (not carry specific external values), `iota` is ideal.
- **Define explicit values for serialization**: Protocol buffers, JSON APIs, and databases need stable enum values. Use explicit numbers: `const StatusOK = 0`, `StatusNotFound = 1`.
- **Document the pattern**: Readers unfamiliar with `iota` may misinterpret `1 << (10 * iota)`. Add a comment explaining the generated values.
- **Combine with typed constants**: `type ByteSize uint64` provides compile-time type safety and prevents accidental mixing with raw integers.
- **Use `iota` for zero-value semantics**: When `iota` starts at 0, the zero value of a typed constant group has meaning (e.g., `Permission(0)` = no permissions, `Weekday(0)` = unknown day).

## Performance implications

- `iota` constants are resolved at compile time — zero runtime cost.
- Bitmask operations (`&`, `|`, `&^`) on `iota`-generated flags compile to single CPU instructions.
- Typed `iota` constants add no overhead compared to untyped constants; the type information is used only at compile time for type checking.
- Switch statements on `iota` constants are optimized into jump tables by the compiler.

## Practice task

Write a program that:

1. Defines a type `Flag uint8` and a `const` block with `iota` for `FlagA = 1 << iota`, `FlagB`, `FlagC`, `FlagD`.
2. Defines a function `set(f, flag Flag) Flag` that sets a flag.
3. Defines a function `clear(f, flag Flag) Flag` that clears a flag.
4. Defines a function `has(f, flag Flag) bool` that tests a flag.
5. In `main()`, create a variable with `FlagA | FlagC`, toggle `FlagB` on, toggle `FlagC` off, and print the binary representation of each step.

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/07-iota
go test ./curriculum/modules/03-programming-fundamentals/lessons/07-iota
```

## Review questions

1. What is the value of each constant in `const ( A = iota; B; C )`?
2. If you add a new constant between `B` and `C` in the block above, which constant values change?
3. What does `iota` equal at the start of a new `const` block?
4. Write a `const` block that produces the values `1, 2, 4, 8, 16` using `iota`.
5. Why should you avoid using `iota` for constants that are serialized to a database or sent over a network?

## NEXT UP

Boolean logic and how `bool` values combine with `&&`, `||`, and `!` to form predicates.
