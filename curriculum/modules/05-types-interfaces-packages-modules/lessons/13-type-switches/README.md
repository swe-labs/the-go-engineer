# Type switches

## Learning objective

Write and use type switches to dispatch behavior based on the dynamic type of an interface value, and distinguish type switches from type assertions.

## Why this matters

Real Go code receives values through interfaces — `error`, `io.Reader`, `any`, or custom interfaces. A type assertion answers "is this specific type?" but a type switch answers "which of many types is this?" Handling every possible type with separate assertions is repetitive, error-prone, and obscures intent. Type switches are the idiomatic, readable way to branch on dynamic type.

## Mental model

A type switch is a sorting machine. Packages arrive on a conveyor belt (the interface value). Each package has a label (its dynamic type). The machine checks the label and routes the package to the correct bin:

- `int` packages go to the blue bin
- `string` packages go to the red bin
- Anything unrecognised goes to the default bin

The variable `v := x.(type)` gives you the package contents already unboxed for the bin it lands in — you do not need to assert again inside the case body.

## Core idea

A type switch is a switch statement where the switch expression is `v := x.(type)`. Each `case` clause lists one or more types. When a match is found, the variable `v` has that concrete type inside the case body.

```go
switch v := x.(type) {
case int:
    // v is int here
case string:
    // v is string here
default:
    // v has the same type as x (interface type)
}
```

Key rules:
- The `.(type)` syntax is only valid inside a `switch` statement, not in any other context.
- `case` ordering matters: the first matching case wins.
- `case nil` matches a nil interface value.
- `default` matches if no other case matches.
- No `fallthrough` in type switches — each case is exclusive.
- A type switch does **not** require `interface{}`; it works with any interface type.

## Under the hood

The compiler transforms a type switch into a sequence of type assertions, but optimises the dispatch. For non-empty interfaces, the compiler compares the type pointer from the interface's `itab` against the target type's `*rtype`. For `any` (empty interface), the compiler may use a hash-based jump table for O(1) dispatch when there are many cases.

The variable `v` declared in `switch v := x.(type)` is a new variable scoped to the switch. In each case arm, the compiler re-binds `v` to the concrete type, which may require an implicit copy of the underlying value.

## How Go uses it

- **`fmt.Sprintf` and friends**: the `fmt` package uses a type switch on the `any` argument to determine formatting (`int`, `string`, `error`, etc.).
- **`encoding/json`**: unmarshaling uses type switches on `any` to decode JSON objects, arrays, strings, numbers, booleans, and null.
- **Error classification**: inspecting an `error` value to check if it is `*os.PathError`, `*net.DNSError`, or a custom type.
- **`database/sql`**: scanning result values into Go types based on the column type.
- **`log/slog`**: handling different argument types in structured logging.

## Go example

```go
package main

import (
	"fmt"
	"time"
)

// describe returns a human-readable description of any value.
func describe(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil"
	case int:
		return fmt.Sprintf("integer %d", x)
	case string:
		return fmt.Sprintf("string %q (len %d)", x, len(x))
	case bool:
		if x {
			return "true"
		}
		return "false"
	case time.Duration:
		return fmt.Sprintf("duration %v (%d ns)", x, int64(x))
	case fmt.Stringer:
		return fmt.Sprintf("Stringer: %s", x.String())
	default:
		return fmt.Sprintf("unknown type %T", v)
	}
}

func main() {
	inputs := []any{nil, 42, "hello", true, time.Second * 3, 3.14}
	for _, v := range inputs {
		fmt.Println(describe(v))
	}
}
```

## Step-by-step execution

For `describe(42)` with `v = 42`:

1. The `switch x := v.(type)` statement evaluates the dynamic type of `v`: `int`.
2. Go compares the dynamic type against each case in source order: `nil` → no, `int` → yes.
3. The variable `x` is bound as `int` with value `42`.
4. The body `fmt.Sprintf("integer %d", x)` executes, producing `"integer 42"`.
5. The switch exits (no fallthrough).

For `describe("hello")`:

1. Dynamic type: `string`.
2. Falls past `nil` and `int` cases.
3. Matches `case string:`.
4. `x` is bound as `string` with value `"hello"`.
5. Produces `string "hello" (len 5)`.

For `describe(3.14)`:

1. Dynamic type: `float64`.
2. Falls past `nil`, `int`, `string`, `bool`, `time.Duration`, `fmt.Stringer` (does not implement it).
3. Reaches `default`.
4. `x` retains the original `any` type. Produces `"unknown type float64"`.

## Common mistakes

- **Using `x` (not the switch variable) inside a case body**: The variable declared in `switch v := x.(type)` is the typed one. Using `x` directly still has the interface type and needs another assertion.
- **Forgetting `default`**: A type switch compiles without `default`, but values of unexpected types pass silently. Always include `default` unless you are certain all types are handled.
- **Expecting fallthrough**: Type switches do not support `fallthrough`. Each case is independent.
- **Case ordering for interface types**: If `case fmt.Stringer:` appears before `case string:`, a `string` value matches `fmt.Stringer` (since `string` implements `fmt.Stringer` via `fmt.Sprintf` conventions). Order more specific (concrete) types before more general (interface) types.
- **Assuming the switch variable is a pointer**: The variable `v` in `case *int:` is `*int`, but the variable in `case int:` is `int`. The extraction unpacks the interface value — if the interface holds `*int`, only the `*int` case matches, not `int`.

## Debugging walkthrough

Consider this broken code:

```go
func inspect(v any) string {
	switch v := v.(type) {
	case string:
		return v  // compile error: cannot use v (type string) as type string in return
	}
	return ""
}
```

**Symptom**: Compile error: `cannot use v (type string) as type string in return argument` — wait, that seems wrong. Let's clarify. The real bug is subtler.

Actual broken code:

```go
func classify(v any) string {
	var result string
	switch v := v.(type) {
	case int:
		result = fmt.Sprintf("int %d", v)
	case string:
		result = fmt.Sprintf("str %s", v)
	}
	return result
}

func main() {
	fmt.Println(classify(42))
	fmt.Println(classify(3.14))  // prints empty string!
}
```

**Symptom**: `classify(3.14)` returns an empty string.

**Root cause**: No `default` case. When `v` is `float64`, no case matches, `result` stays at its zero value `""`.

**Fix**: Add a `default` case:

```go
default:
    result = fmt.Sprintf("unknown %T", v)
```

## Production notes

- **Exhaustive matching**: When adding a new concrete type to a system, trace all type switches that dispatch on that interface. The compiler does not warn about missing cases — use `default` and a panic or log for unhandled types in critical paths.
- **Error classification pattern**: Type switches are the idiomatic way to inspect errors:
  ```go
  var perr *os.PathError
  if errors.As(err, &perr) { /* ... */ }
  ```
  But when you have many possible error types, a type switch on a custom error interface is simpler.
- **gRPC interceptor pattern**: Type switches on `any` are used in middleware to inspect and modify context values.
- **JSON decoding flexibility**: Use type switches to handle "weakly typed" JSON where a field might be a string or number.

## Performance implications

- A type switch with concrete types compiles to an efficient jump table or binary search over type pointers — O(1) or O(log n) dispatch.
- Type switches are significantly faster than a chain of `if`/`else` type assertions because the compiler can optimise the entire dispatch in one pass.
- When a case matches an interface type (e.g., `case fmt.Stringer:`), the runtime must check interface satisfaction, which is slightly slower than a concrete type match.
- For hot paths, place the most common type first to minimise comparisons.

## Practice task

Write a function `mathSwitch(v any) (float64, error)` that:

- Returns the value unchanged if `v` is `float64`.
- Returns `float64(x)` if `v` is `int`.
- Returns `float64(x)` if `v` is `int64`.
- Returns `x * 2` if `v` is `float32`.
- Returns `0, fmt.Errorf("unsupported type: %T", v)` for any other type.

Then write a `main()` that calls it with `int(10)`, `int64(20)`, `float32(3.5)`, `float64(2.5)`, and `"hello"` and prints each result. Run and verify.

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/13-type-switches
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/13-type-switches
```

## Review questions

1. What syntax distinguishes a type switch from a regular switch?
2. Why does case ordering matter when mixing concrete types and interface types?
3. What does `case nil:` match in a type switch?
4. Does a type switch support `fallthrough`? What happens if you try?
5. How does the type of `v` differ in a `case int:` arm versus the `default` arm?

## NEXT UP

Nil interfaces — what happens when an interface variable is nil, and why `interface != nil` often fails to detect it.
