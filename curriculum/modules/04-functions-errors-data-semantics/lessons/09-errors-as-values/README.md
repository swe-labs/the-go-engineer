# Errors as values

## Learning objective

Use Go's `error` interface to represent and propagate errors, apply the `if err != nil` pattern consistently, and distinguish `nil` errors from non-nil errors containing a `nil` pointer.

## Why this matters

Go does not have exceptions. Instead, errors are ordinary values that functions return. This design forces you to handle errors explicitly, making control flow visible and predictable. The `if err != nil` pattern appears more than any other in production Go code. Mastering it is required to write robust Go programs.

## Mental model

An error is a value that implements the `error` interface:

```go
type error interface {
    Error() string
}
```

Think of a function returning a result as carrying a second box alongside the result. If the second box is `nil`, everything worked. If it is not `nil`, something went wrong — open the box, read the message, and decide what to do.

Because `error` is an interface, any type that has an `Error() string` method qualifies as an error. This includes `errors.New`, `fmt.Errorf`, sentinel errors, and custom error types.

## Core idea

Functions that can fail return `(T, error)`:

```go
func parseInt(s string) (int, error) {
    n, err := strconv.Atoi(s)
    if err != nil {
        return 0, err
    }
    return n, nil
}
```

The caller checks the error immediately:

```go
n, err := parseInt("42")
if err != nil {
    log.Fatal(err)
}
fmt.Println(n)
```

**Nil error** means success. A non-nil error means failure — the other return values should be ignored (they are typically zero values).

**Sentinel errors** are named package-level error values that callers can compare with `==`:

```go
var ErrNotFound = errors.New("item not found")

func lookup(id int) (string, error) {
    if id < 1 {
        return "", ErrNotFound
    }
    return "item-" + itoa(id), nil
}
```

**Error is a value** — you can store errors in variables, return them from functions, put them in slices, and inspect them programmatically.

## Under the hood

`errors.New` returns a pointer to a struct with a single string field. Two `errors.New` calls with the same message produce different pointer values — they are not `==`-equal.

`fmt.Errorf` uses `errors.New` internally for the simple case, but with `%w` (Go 1.13+), it wraps an existing error, enabling `errors.Is` and `errors.As` to unwrap the chain.

The `error` interface is two words wide (type pointer + data pointer). A `nil` error has both pointers set to `nil`. A non-nil error stored in an `error` variable has a non-nil type pointer. This is why `if err != nil` can be true even when the underlying concrete error value is nil — the interface variable is not nil.

## How Go uses it

- **`if err != nil`** is the most common pattern. Every Go engineer writes it dozens of times per day.
- **Sentinel errors**: `io.EOF`, `sql.ErrNoRows`, `os.ErrNotExist` — defined as package-level variables for `==` comparison.
- **Error wrapping**: `fmt.Errorf("context: %w", err)` preserves the original error for unwrapping.
- **Error inspection**: `errors.Is(err, target)` and `errors.As(err, target)` unwrap the chain.
- **Custom error types**: Structs implementing `Error() string` can carry additional context (codes, fields, stack traces).

## Go example

```go
package main

import (
	"errors"
	"fmt"
)

var ErrNegative = errors.New("value is negative")
var ErrZero = errors.New("value is zero")

// Returns (T, error) — idiomatic Go.
func positiveSquare(n int) (int, error) {
	if n < 0 {
		return 0, ErrNegative
	}
	if n == 0 {
		return 0, ErrZero
	}
	return n * n, nil
}

func main() {
	inputs := []int{4, -1, 0, 9}
	for _, n := range inputs {
		result, err := positiveSquare(n)
		if err != nil {
			fmt.Printf("positiveSquare(%d) error: %v\n", n, err)
		} else {
			fmt.Printf("positiveSquare(%d) = %d\n", n, result)
		}
	}

	// Sentinel error comparison.
	_, err := positiveSquare(-5)
	if err == ErrNegative {
		fmt.Println("Got expected negative error")
	}

	// Error as a value: store in a slice.
	errs := []error{ErrNegative, ErrZero, errors.New("custom error")}
	fmt.Println("Stored errors:")
	for _, e := range errs {
		fmt.Println(" -", e)
	}
}
```

## Step-by-step execution

For `positiveSquare(-1)`:

1. `n < 0` is true.
2. Returns `0, ErrNegative`.
3. Caller receives `result = 0, err = ErrNegative`.
4. `if err != nil` → true.
5. Prints `"positiveSquare(-1) error: value is negative"`.

For `positiveSquare(4)`:

1. `n < 0` false, `n == 0` false.
2. Computes `n * n = 16`.
3. Returns `16, nil`.
4. `if err != nil` → false.
5. Prints `"positiveSquare(4) = 16"`.

## Common mistakes

- **Checking `err != nil` after ignoring it**: `result, _ := positiveSquare(-1)` discards the error. The caller gets `0` and thinks it is a valid result.
- **Comparing sentinel errors with `==` on wrapped errors**: `err == ErrNegative` is false if `err` was wrapped with `%w`. Use `errors.Is(err, ErrNegative)` instead.
- **Returning a nil pointer in an `error` interface**: `return nil` with type `*MyError` in a function returning `error` — the returned `error` is non-nil even though the pointer is nil. Always use `return nil` (typed as `error`), not `return *MyError(nil)`.
- **Swallowing errors in a defer**: `defer file.Close()` ignores the close error. In critical paths, capture and check it.

## Debugging walkthrough

This code has a subtle nil-interface bug:

```go
package main

import "fmt"

type MyError struct {
	Code int
}

func (e *MyError) Error() string {
	return fmt.Sprintf("code %d", e.Code)
}

func doSomething(flag bool) error {
	var me *MyError
	if flag {
		me = &MyError{Code: 42}
	}
	return me // me is nil, but the returned error is NOT nil!
}

func main() {
	err := doSomething(false)
	fmt.Println(err == nil) // prints false!
}
```

**Symptom**: `err == nil` is `false` even though `doSomething(false)` should succeed.

**Root cause**: `me` is a `*MyError` typed nil. When returned as `error`, the interface has a non-nil type descriptor (`*MyError`) but a nil data pointer. The interface is not nil.

**Fix**: Always return a bare `nil` when you have no error:

```go
func doSomething(flag bool) error {
	if flag {
		return &MyError{Code: 42}
	}
	return nil
}
```

## Production notes

- **Always handle errors**: Every `err` should be checked, logged, wrapped, or explicitly discarded with a comment explaining why.
- **Sentinel error naming**: `ErrXxx` for package-level exported errors, `errXxx` for unexported. This is Go convention.
- **Error messages**: Lowercase, no trailing punctuation. `fmt.Errorf("reading config: %w", err)`.
- **Don't panic**: Reserve `panic` for truly unrecoverable states (programmer bugs, corruption). Use errors for everything else.
- **Wrap errors with context**: `fmt.Errorf("connecting to %s: %w", addr, err)` preserves the original and adds context.

## Performance implications

- **`errors.New`** allocates a struct on the heap. The same sentinel error created once at package init costs one allocation.
- **`fmt.Errorf` without `%w`** allocates a new error each call. On hot paths, consider caching errors or using sentinel errors.
- **`fmt.Errorf with %w`** is slightly more expensive because it builds a wrapping structure.
- **Custom error types** can avoid allocations by implementing `Error()` on a value receiver and using pre-allocated error instances.
- **Error checking** (`if err != nil`) is a single interface comparison — effectively free.

## Practice task

Write a function `parseAge(s string) (int, error)` that:
- Uses `strconv.Atoi` to convert.
- Returns `ErrInvalidAge` if the string does not parse.
- Returns `ErrNegativeAge` if the age is negative.
- Returns `ErrUnreasonableAge` if the age is > 150.
- Otherwise returns `(age, nil)`.

Define the three sentinel errors as package-level variables. In `main()`, test with `"25"`, `"-1"`, `"abc"`, `"200"`.

## Tests / verification

```bash
go run ./curriculum/modules/04-functions-errors-data-semantics/lessons/09-errors-as-values
go test ./curriculum/modules/04-functions-errors-data-semantics/lessons/09-errors-as-values
```

## Review questions

1. What is the `error` interface definition in Go?
2. What does a `nil` error mean in a function that returns `(T, error)`?
3. How do you define a sentinel error that callers can compare with `==`?
4. Why can `if err != nil` be true when the underlying error pointer is nil? How do you fix it?
5. What is the difference between `errors.New("msg")` and `fmt.Errorf("msg")`?

## NEXT UP

Error wrapping — wrapping errors with `fmt.Errorf` and `%w`, and the `errors.Is` / `errors.As` functions.
