# Nil interfaces

## Learning objective

Distinguish between a nil interface value and an interface holding a nil pointer, explain why `interface != nil` can fail to detect a nil concrete value, and apply safe patterns to avoid typed nil traps.

## Why this matters

The typed nil trap is one of the most common bugs in production Go code. A function returns a nil `*os.PathError` wrapped in an `error` interface; the caller checks `if err != nil` — and it evaluates to `true`, even though the concrete pointer is nil. This causes nil pointer dereference panics, misrouted error handling, and subtle logic errors that are extremely hard to debug. Every professional Go engineer must understand this distinction.

## Mental model

An interface value is a box with two compartments:
- The **type label** (what kind of value is inside)
- The **data pointer** (the actual value)

A truly nil interface is an empty box — both compartments are empty (nil, nil).

A typed nil is a box with a label but nothing useful inside — e.g. `(*os.PathError, nil)`. The box is not empty (it has a label), so `if box != nil` says "not empty".

Checking the box for nil only tells you whether the box itself is empty. It does **not** tell you whether the contents are a nil pointer.

## Core idea

```go
var a any = nil          // true nil: type=nil, data=nil → a == nil is true

var p *int = nil
var b any = p            // typed nil: type=*int, data=nil → b == nil is false
```

An interface value is nil **only when both its type and value are nil**. Assigning a nil pointer of concrete type to an interface sets the type field, making the interface non-nil.

```go
func returnsNilPtr() *int { return nil }
var err error = returnsNilPtr()  // err is (*int, nil), not nil!
err == nil                       // false — the trap
```

## Under the hood

In the Go runtime, an interface value is represented as a `runtime.eface` (for `any`) or `runtime.iface` (for non-empty interfaces):

```go
type eface struct {
    _type *_type
    data  unsafe.Pointer
}

type iface struct {
    tab  *itab
    data unsafe.Pointer
}
```

The `== nil` comparison checks whether `_type == nil` (for eface) or `tab == nil` (for iface). When you assign a nil `*int` to `any`, `_type` points to the `*int` type metadata and `data` is nil. Since `_type != nil`, the comparison returns false.

## How Go uses it

- **`error` interface**: The most common victim of the typed nil trap. A function with `func() error` returning a nil `*MyError` pointer creates a non-nil `error` value.
- **`io.Reader` / `io.Writer`**: Wrapper types that hold a nil `*bytes.Buffer` or nil `*os.File` produce non-nil interfaces.
- **gRPC interceptors**: Returning a nil `*status.Status` as `error` triggers the trap.
- **JSON marshaling**: `json.Marshal` on a typed nil interface produces `"null"`, not an error — surprising when the value "should be empty".

## Go example

```go
package main

import "fmt"

// MyError is a custom error type.
type MyError struct {
	Msg string
}

func (e *MyError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return e.Msg
}

// doWork returns a nil *MyError wrapped in error.
func doWork() error {
	var err *MyError // nil pointer of concrete type
	return err       // wraps as (*MyError, nil) — non-nil interface!
}

// doWorkSafe returns a nil error when there is no error.
func doWorkSafe() error {
	return nil // true nil interface
}

func main() {
	err := doWork()
	fmt.Printf("doWork: err == nil = %v, err = %v\n", err == nil, err)

	err2 := doWorkSafe()
	fmt.Printf("doWorkSafe: err == nil = %v, err = %v\n", err2 == nil, err2)

	// Safe pattern: use a helper to unwrap typed nils.
	fmt.Println("Safe nil check:", safeError(err))
}

func safeError(err error) string {
	if err == nil {
		return "no error"
	}
	return fmt.Sprintf("error: %v", err)
}
```

## Step-by-step execution

For `doWork()`:

1. Declare `var err *MyError` — `err` is a nil pointer of type `*MyError`.
2. `return err` — Go wraps the value into an `error` interface. The interface gets type `*MyError` and data pointing to nil.
3. Caller receives the `error` value `e` with `(*MyError, nil)`.
4. `e == nil` compares the interface's type field: `*MyError != nil`, so result is `false`.

For `doWorkSafe()`:

1. `return nil` — returns a literal nil. Go creates an interface with type=nil, data=nil.
2. Caller receives a truly nil `error`.
3. `e == nil` is `true`.

## Common mistakes

- **Returning named error variables of concrete type**: `func f() error { var err *MyError; return err }` creates a typed nil. Use `return nil` for the success path.
- **Checking `err != nil` after a function that returns `(*MyError, error)`**: The caller should always compare the `error` return, not the concrete pointer.
- **Using reflection**: `reflect.ValueOf(err).IsNil()` panics if `err` is a nil interface. You must check `err == nil` first.
- **Assuming `json.Unmarshal` produces nil errors for empty input**: It does — but wrapping it in a custom error type introduces the trap.

## Debugging walkthrough

Consider this seemingly correct code:

```go
type Result struct {
    Value int
    Err   *MyError
}

func (r *Result) Error() string {
    if r.Err != nil {
        return r.Err.Msg
    }
    return ""
}

func process() *Result {
    // ... returns nil Result on success
    return nil
}

func main() {
    r := process()
    if r != nil {
        fmt.Println("has result")
    }
}
```

**Symptom**: Works fine — `r` is a true nil `*Result`.

But change `process` to return `error`:

```go
func process() error {
    var r *Result
    return r  // typed nil!
}

func main() {
    if err := process(); err != nil {
        fmt.Println("error occurred") // THIS PRINTS!
    }
}
```

**Root cause**: `return r` wraps nil `*Result` into `error` as `(*Result, nil)`. The interface is non-nil.

**Fix**: Always return a bare `nil` for the success case:

```go
func process() error {
    var r *Result
    if r == nil {
        return nil  // true nil interface
    }
    return r
}
```

## Production notes

- **Linter rule**: The Go vet command (`go vet`) catches some typed nil assignments. Run `go vet ./...` in CI.
- **Interface constructors**: A function returning an interface should return `nil` explicitly on success, never a nil concrete pointer.
- **`errors.Is` and `errors.As`**: These functions handle typed nils correctly because they inspect the unwrapped error chain, but they still need a non-nil entry point.
- **Protocol buffers / gRPC**: Protobuf-generated types are pointers. Returning `nil *pb.MyMessage` as `error` is a common variant of the trap.

## Performance implications

- The typed nil check itself has negligible cost — it is two pointer comparisons.
- The real cost is debugging time. Typed nil bugs have caused production outages in major Go services. The runtime cost is zero; the human cost is high.
- Reflection-based nil checks (`reflect.ValueOf(x).IsNil()`) are expensive (~50-100ns) and should not be used in hot paths.

## Practice task

Write a function `getError(shouldFail bool) error` that:

- If `shouldFail` is true, returns a non-nil `*MyError` with message "something went wrong".
- If `shouldFail` is false, returns `nil` (a true nil interface, not a nil `*MyError`).

Then write `main()` that calls `getError(false)` and prints whether the result `== nil`. Also call `getError(true)` and print its message. Run and verify both cases work correctly.

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/14-nil-interfaces
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/14-nil-interfaces
```

## Review questions

1. What two components make up an interface value at runtime?
2. Why does `var e error = (*MyError)(nil)` result in `e != nil`?
3. What happens when you compare a typed nil interface with `== nil`?
4. How can `go vet` help catch typed nil traps?
5. What is the correct way to return a nil error from a function whose return type is an interface?

## NEXT UP

Package names — how Go organises code into packages, naming conventions, and the relationship between directories and packages.
