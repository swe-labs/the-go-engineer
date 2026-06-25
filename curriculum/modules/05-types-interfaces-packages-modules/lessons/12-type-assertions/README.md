# Type assertions

## Learning objective

Extract concrete values from interface values using the type assertion syntax `x.(T)`, handle failures with the comma-ok idiom, and distinguish asserting to concrete vs interface types.

## Why this matters

When you receive a value as an interface — from a function parameter, a channel receive, or an `interface{}` — you often need the concrete type back. Type assertions are Go's mechanism for this unwrapping. They appear in JSON decoding, error inspection, type switches, and every `interface{}` consumer. Using them correctly and safely is a core Go skill.

## Mental model

Think of an interface value as a sealed box containing a concrete value of some type. A type assertion `x.(T)` tries to open the box. If the box contains a `T`, you get the value. If not, you either get a panic (single-return form) or a `false` (comma-ok form). The comma-ok form is like asking "is this a T?" and getting a yes/no answer.

## Core idea

**Syntax**:

```go
value := x.(T)        // panics if x does not hold T
value, ok := x.(T)    // ok is false if x does not hold T
```

`x` must be an interface value. `T` can be a concrete type or an interface type.

**Comma-ok form** (safe):

```go
var v interface{} = "hello"
s, ok := v.(string)
if ok {
    fmt.Println("is a string:", s)
}
```

**Panicking form** (use when you are certain):

```go
s := v.(string) // panics if v is not a string
```

**Asserting to concrete type**: `v.(string)` extracts the underlying string value.

**Asserting to interface type**: `v.(io.Reader)` checks if the concrete value implements `io.Reader`. If yes, returns the interface; if no, comma-ok returns `false`.

## Under the hood

A type assertion checks the concrete type stored in the interface's `itab`. For concrete type assertions, it compares the type pointer. For interface assertions, it checks if the concrete type's method set satisfies the target interface (may involve building a new `itab`). Both operations are O(1).

## How Go uses it

- **json.Unmarshal** returns `interface{}`; callers assert to `map[string]interface{}` or `[]interface{}`.
- **error unwrapping**: `errors.As(err, &target)` uses type assertions internally.
- **http.ResponseWriter** is often asserted to `http.Hijacker` or `http.Flusher` for specific capabilities.
- **io.Reader** wrapping: check if an `io.Reader` is also an `io.WriterTo` for optimization.
- **Type switches**: `switch v := x.(type)` uses type assertions in a structured form.

## Go example

```go
package main

import (
	"fmt"
	"io"
	"strings"
)

func printType(v interface{}) {
	s, ok := v.(string)
	if ok {
		fmt.Println("string:", s)
		return
	}
	i, ok := v.(int)
	if ok {
		fmt.Println("int:", i)
		return
	}
	b, ok := v.(bool)
	if ok {
		fmt.Println("bool:", b)
		return
	}
	fmt.Printf("unknown type: %T\n", v)
}

func main() {
	printType("hello")
	printType(42)
	printType(true)
	printType(3.14)

	var r io.Reader = strings.NewReader("hello")
	// Assert to concrete type
	sr := r.(*strings.Reader)
	fmt.Printf("concrete: %T, len: %d\n", sr, sr.Len())

	// Assert to interface type
	if w, ok := r.(io.WriterTo); ok {
		fmt.Println("also implements WriterTo")
		_ = w
	} else {
		fmt.Println("does not implement WriterTo")
	}

	// Panicking form
	v := interface{}("must be string")
	s := v.(string)
	fmt.Println("safe assertion:", s)
}
```

## Step-by-step execution

For `s, ok := v.(string)` where `v = interface{}("hello")`:

1. `v` is an interface value: `{type: string, data: &"hello"}`.
2. Type assertion `v.(string)` checks: does `v`'s concrete type equal `string`?
3. Yes. `ok` is `true`. `s` is a copy of the string `"hello"`.
4. If `v` were `42` (int), the check would fail. `ok` would be `false`. `s` would be `""` (zero value of `string`).

For panicking form `s := v.(string)`:

1. Same check. If `v` holds a `string`, returns the value.
2. If `v` does not hold a `string`, the runtime panics.

For `r.(io.WriterTo)` where `r` is `*strings.Reader`:

1. `r`'s concrete type is `*strings.Reader`.
2. Compiler checks at compile time that `*strings.Reader` could satisfy `io.WriterTo`.
3. At runtime, checks the method set of `*strings.Reader` against `io.WriterTo`.
4. `strings.Reader` does not have `WriteTo`, so `ok` is `false`.

## Common mistakes

- Mistake: Using type assertion on a non-interface value.
  - Fix: Type assertions only work on interface values. If you have a concrete type, use a type conversion instead.

- Mistake: Forgetting the comma-ok and causing a panic on unexpected types.
  - Fix: Always use the comma-ok form unless you are absolutely certain of the type.

- Mistake: Asserting to a nil interface.
  - Fix: Check `x != nil` before asserting.

- Mistake: Expecting a type assertion to convert between related types (e.g., `int` to `int64`).
  - Fix: Type assertions extract the exact concrete type. Use type conversion for related types.

- Mistake: Assuming `x.(T)` modifies `x` or the original value.
  - Fix: `x.(T)` returns a copy (for value types) or the same pointer (for pointer/reference types).

## Debugging walkthrough

This code panics:

```go
func process(v interface{}) {
    s := v.(string) // panic if v is not string
    fmt.Println(s)
}

func main() {
    process(42)
}
```

**Symptom**: `panic: interface conversion: interface {} is int, not string`.

**Root cause**: Type assertion in panicking form on a non-string value.

**Fix**: Use comma-ok form:

```go
func process(v interface{}) {
    if s, ok := v.(string); ok {
        fmt.Println(s)
    } else {
        fmt.Println("not a string")
    }
}
```

## Production notes

- **Prefer comma-ok** form in all production code. Panics from type assertions are a common source of crashes.
- **Type assertions against interface types** are useful for optional feature detection (e.g., does this `io.Writer` also support `io.ReaderFrom`?).
- **Avoid deep type assertion chains**. If you are asserting many times on the same value, consider a type switch.
- **Code generation**: `json.Unmarshal` into `map[string]interface{}` requires many type assertions. Consider using typed structs instead.

## Performance implications

- **Concrete type assertion** is a single pointer comparison (the type pointer in the `itab`). ~1 ns.
- **Interface type assertion** requires checking the method set against the concrete type. Slower but still ~5-10 ns.
- **Both are fast** — do not avoid type assertions for performance reasons.

## Practice task

Write a function `classify(items []interface{})` that iterates over a slice and prints each element's type classification: "string", "int", "bool", "float64", or "other". Use the comma-ok form. In `main()`, create a slice with values of different types and call `classify`. Verify with tests.

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/12-type-assertions
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/12-type-assertions
```

## Review questions

1. What is the difference between `x.(T)` returning one value vs two values?
2. Can you use a type assertion on a value of concrete type (not interface)? Why or why not?
3. What happens at runtime when you assert `v.(string)` and `v` holds an `int`?
4. What is the difference between asserting to a concrete type and asserting to an interface type?
5. When would you use the panicking form of type assertion in production?

## NEXT UP

Type switches — a cleaner syntax for dispatching on multiple types with `switch v := x.(type)`.
