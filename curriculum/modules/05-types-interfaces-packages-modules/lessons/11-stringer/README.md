# Stringer

## Learning objective

Implement the `fmt.Stringer` interface to control how your types are printed, distinguish `String()` from `GoString()`, and apply `Stringer` on enum-like types.

## Why this matters

`fmt.Stringer` is the most commonly implemented custom interface in Go. Whenever you `fmt.Println` or `fmt.Sprintf` with `%v` on your type, Go calls its `String()` method. Proper `String()` implementations make debugging output readable, logs informative, and user-facing messages polished. If your type doesn't have `String()`, you get `{Field1:val Field2:val}` — functional but ugly.

## Mental model

`String()` is how your type introduces itself in text. When Go needs to convert your value to a string (for printing, logging, or error messages), it checks: "does this type have a `String() string` method?" If yes, call it. If no, fall back to the default `%+v` reflection-based output. Every type can have one `String()` method that controls all default string formatting.

## Core idea

**fmt.Stringer interface**:

```go
type Stringer interface {
    String() string
}
```

Any type with a `String() string` method satisfies `fmt.Stringer`. The `fmt` package uses it for `%v`, `%s`, `%+v`, and `%#v` (unless a more specific verb like `%d` applies).

```go
type Color struct {
    Name string
    Hex  string
}

func (c Color) String() string {
    return fmt.Sprintf("%s (#%s)", c.Name, c.Hex)
}
```

**GoStringer**: The `%#v` verb (Go-syntax representation) uses `GoString() string` if defined:

```go
func (c Color) GoString() string {
    return fmt.Sprintf("Color{Name: %q, Hex: %q}", c.Name, c.Hex)
}
```

**Stringer on enum types**: Types defined with `const` and `iota` benefit greatly from `String()`:

```go
type Status int

const (
    StatusPending Status = iota
    StatusActive
    StatusInactive
)

func (s Status) String() string {
    switch s {
    case StatusPending:
        return "pending"
    case StatusActive:
        return "active"
    case StatusInactive:
        return "inactive"
    default:
        return fmt.Sprintf("Status(%d)", int(s))
    }
}
```

## Under the hood

When `fmt` needs to format a value, it checks if the value implements `Stringer` via a type assertion. If yes, it calls `String()`. This is a single interface dispatch. The result is then used directly for `%s` and `%v`, or quoted for `%q`.

## How Go uses it

- **error** values implement `Error() string`, which is similar but distinct from `String()`.
- **time.Time** implements `String()` returning `"2006-01-02 15:04:05.999999999 -0700 MST"`.
- **net.IP** implements `String()` returning `"192.0.2.1"`.
- **url.URL** implements `String()` returning the full URL.
- **big.Int** implements `String()` returning the decimal representation.

## Go example

```go
package main

import "fmt"

type Status int

const (
	StatusPending Status = iota
	StatusActive
	StatusInactive
)

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusActive:
		return "active"
	case StatusInactive:
		return "inactive"
	default:
		return fmt.Sprintf("Status(%d)", int(s))
	}
}

type Point struct {
	X, Y int
}

func (p Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func (p Point) GoString() string {
	return fmt.Sprintf("Point{X: %d, Y: %d}", p.X, p.Y)
}

func main() {
	fmt.Println(StatusPending)   // "pending"
	fmt.Println(StatusActive)    // "active"
	fmt.Println(Status(99))      // "Status(99)"

	p := Point{X: 3, Y: 4}
	fmt.Println(p)               // "(3, 4)"
	fmt.Printf("%v\n", p)        // "(3, 4)"
	fmt.Printf("%#v\n", p)       // "Point{X: 3, Y: 4}"
}
```

## Step-by-step execution

For `fmt.Println(p)` where `p = Point{X: 3, Y: 4}`:

1. `fmt.Println` receives `p` as `interface{}`.
2. Inside `fmt`, the formatter checks: does `p` implement `fmt.Stringer`?
3. Yes, `Point` has `String() string`.
4. `fmt` calls `p.String()`.
5. `p.String()` returns `"(3, 4)"`.
6. `fmt` writes `"(3, 4)\n"` to stdout.

For `fmt.Printf("%#v\n", p)`:

1. `%#v` requests Go-syntax representation.
2. Formatter checks for `GoStringer` (higher priority than `Stringer` for `%#v`).
3. `Point` has `GoString() string`.
4. Returns `"Point{X: 3, Y: 4}"`.

## Common mistakes

- Mistake: Defining `String()` with a pointer receiver but calling `fmt.Println` on a value.
  - Fix: Define `String()` with a value receiver, or always print with `&` if using pointer receiver.

- Mistake: `String()` method causing infinite recursion by calling `fmt.Sprintf` with `%v` on itself.
  - Fix: Use `%d` or convert to a primitive type explicitly:

```go
func (s Status) String() string {
    return fmt.Sprintf("Status(%d)", int(s)) // int() avoids recursion
}
```

- Mistake: Forgetting to handle all enum values in the `String()` switch, including unknown values.
  - Fix: Always include a `default` case that returns a formatted string with the raw value.

- Mistake: Defining `String()` on a type that changes the value (side effects).
  - Fix: `String()` should be a pure function — no mutation, no I/O.

## Debugging walkthrough

This code causes infinite recursion:

```go
type Status int

func (s Status) String() string {
    return fmt.Sprintf("Status(%v)", s) // BUG: calls String() again
}
```

**Symptom**: Stack overflow / infinite recursion.

**Root cause**: `fmt.Sprintf("Status(%v)", s)` calls `s.String()` because `%v` uses `Stringer`. This loops forever.

**Fix**: Convert to the underlying type first:

```go
func (s Status) String() string {
    return fmt.Sprintf("Status(%d)", int(s))
}
```

## Production notes

- **Every exported type should have a `String()` method** if it appears in logs or user-facing output.
- **`String()` must not panic**. Defensive handling of nil receivers is good practice.
- **Enum `String()` methods** are often code-generated (e.g., `stringer` tool from `golang.org/x/tools`).
- **`GoString()`** is useful for debugging; it should produce output that could be used in Go source code.
- **Consistency**: use the same style for all `String()` methods in a package.

## Performance implications

- **Interface dispatch**: Calling `String()` through `fmt.Stringer` is a single interface method call. Negligible cost.
- **Allocation**: `String()` typically allocates a new string (via `fmt.Sprintf` or string concatenation). This adds GC pressure if called frequently in hot paths.
- **Avoid String() in tight loops**: If you need to format many values in a hot path, consider formatting once and reusing.

## Practice task

Define a `Card` type representing a playing card with fields `Suit` and `Rank` (both int-based enums). Implement `String()` so a card prints as `"Ace of Spades"`, `"King of Hearts"`, etc. Also implement `GoString()` returning a Go-expression. In `main()`, print several cards with `%v` and `%#v`.

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/11-stringer
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/11-stringer
```

## Review questions

1. What method must a type implement to satisfy `fmt.Stringer`?
2. What is the difference between `String()` and `GoString()`?
3. Which `fmt` verb uses `GoStringer`?
4. Why does `fmt.Sprintf("%v", s)` inside `String()` cause infinite recursion?
5. What is the recommended way to handle unknown enum values in `String()`?

## NEXT UP

Type assertions — extracting concrete values from interface values with `x.(T)`.
