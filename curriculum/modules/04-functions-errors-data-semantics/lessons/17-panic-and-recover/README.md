# panic and recover

## Learning objective

Distinguish between errors (expected failures) and panics (unexpected invariants), use `recover` in a deferred function to catch a panic, and build a panic-recovery middleware pattern.

## Why this matters

Go's philosophy is explicit error handling: return errors for expected failures. But some failures are not expected -- nil pointer dereferences, out-of-bounds access, or invariant violations. These are panics. A panic in a production HTTP server should not crash the entire process; it should return a 500 response to the affected request. The `recover` mechanism, combined with `defer`, lets you build this safety net. Understanding when to panic and when to return an error is a hallmark of mature Go engineering.

## Mental model

Think of `panic` as pulling a fire alarm. Normal operation (returning errors) is like a door closing normally. A panic is a building evacuation: every floor (function frame) must be evacuated in order, and each floor runs its cleanup checklist (deferred functions) before leaving. `recover` is the fire chief who can call off the evacuation at a specific floor -- the stack stops unwinding and normal execution resumes.

## Core idea

- **`panic(v interface{})`**: Immediately stops execution of the current function, begins unwinding the stack, running deferred functions in each frame. If the unwinding reaches the top of the goroutine's stack, the program crashes with the panic value and stack trace.
- **`recover() interface{}`**: Called inside a deferred function, it stops the unwinding and returns the panic value. If there is no panic, `recover()` returns `nil`.
- **`recover` only works in a deferred function**: Calling `recover` directly (not in a defer) returns `nil`. Calling it in a deferred function in a callee frame (not the panicking frame) returns `nil`.

The contract: `panic` is for programmer bugs and invariants that should never happen. `error` is for everything else (invalid input, network failures, disk full).

## Under the hood

The Go runtime maintains a per-goroutine stack of `_panic` structs. When `panic` is called:

1. `runtime.gopanic` creates a `_panic` struct with the panic value.
2. It enters a loop that walks the goroutine's defer list (LIFO).
3. For each `_defer`, it executes the deferred function.
4. If the deferred function calls `runtime.gorecover`, the `_panic` is marked recovered, the panic value is returned, and the loop terminates.
5. The goroutine continues executing from after the `recover()` call in the deferred function.
6. If no `recover` is found, `gopanic` prints the panic value and stack trace, then exits the process (or propagates to the goroutine's parent).

A double panic (panic inside a deferred function during unwinding) is detected by checking if a `_panic` struct already exists. The runtime then prints both panic values and exits immediately.

## How Go uses it

```go
// HTTP recovery middleware (net/http server already does this)
func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                log.Printf("panic recovered: %v", err)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

// Per-item recovery in batch processing
func processBatch(items []Item) {
    for _, item := range items {
        func() {
            defer func() {
                if err := recover(); err != nil {
                    log.Printf("skipping item %v due to panic: %v", item, err)
                }
            }()
            processItem(item) // may panic
        }()
    }
}
```

## Go example

```go
package main

import (
	"errors"
	"fmt"
	"log"
)

var ErrInvalidInput = errors.New("invalid input")

func riskyOperation(input string) (result string, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic: %v", r)
			err = fmt.Errorf("panic recovered: %v", r)
		}
	}()
	if input == "" {
		panic("input must not be empty")
	}
	if input == "error" {
		return "", ErrInvalidInput
	}
	return "processed: " + input, nil
}

func main() {
	for _, input := range []string{"hello", "error", ""} {
		result, err := riskyOperation(input)
		if err != nil {
			fmt.Printf("riskyOperation(%q): error: %v\n", input, err)
		} else {
			fmt.Printf("riskyOperation(%q): %s\n", input, result)
		}
	}
}
```

## Step-by-step execution

For `riskyOperation("")`:

1. `defer func() { ... }()` is registered.
2. `input == ""` is true → `panic("input must not be empty")` is called.
3. The deferred function runs (stack unwinding).
4. Inside the deferred function, `recover()` returns `"input must not be empty"`.
5. The panic is stopped. The deferred function assigns the error to the named return `err`.
6. The function returns `("", error(...))` normally.
7. `main` prints the error.

Without the deferred `recover`, step 4-6 would be: the program prints the panic value and stack trace, then exits with a non-zero status.

## Common mistakes

- **Calling `recover` outside a deferred function**: `recover()` returns `nil` if called directly in a function body, not inside a `defer`. It must be inside a `defer` that runs during the unwinding.

- **Calling `recover` in the wrong goroutine**: A panic only unwinds the stack of the goroutine that panicked. `recover()` in another goroutine cannot catch it.

- **Using panic for routine error handling**: Panic should not be used like an exception. Never do `panic(err)` and `recover()` as a substitute for `return err`. This is non-idiomatic and makes control flow impossible to follow.

- **Swallowing panics silently**: `recover()` that does nothing is dangerous. At minimum, log the panic value. Silent recovery hides bugs.

- **Believing `recover` catches all panics**: A panic that happens in a new goroutine cannot be recovered by the parent goroutine's deferred function.

## Debugging walkthrough

```go
package main

import "fmt"

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered:", r)
		}
	}()
	fmt.Println("before")
	panickingFunc()
	fmt.Println("after") // does this run?
}

func panickingFunc() {
	panic("oh no!")
}
```

**Prediction**: Many beginners expect `"after"` to print.

**Actual output**:
```
before
recovered: oh no!
```

**Explanation**: When `panickingFunc` panics, the deferred function in `main` runs during unwinding. After recovery, execution continues after the deferred function, which is after the `defer` statement in `main`. But `"after"` never prints because the `defer` is registered in `main`, and after recovery, `main` returns immediately -- control does NOT resume at the point of the panic. It resumes after the `defer` block, which is at the end of `main`.

**Root cause**: Recovery does not resume at the panic site. It resumes after the deferred function that called `recover`. In `main`, this is at the end of the function.

## Production notes

- **Always recover in HTTP handlers**: The `net/http` server already recovers from panics in handler goroutines. If you write your own server, add recovery middleware. Never let a panic crash the process.

- **Use panic only for truly exceptional conditions**: Nil pointer dereference, index out of range, type assertion failure, closed channel send. These indicate programmer bugs, not runtime conditions.

- **Library code should never panic**: Library functions should always return errors. If a library panics, the caller cannot recover (they might not expect it). The only exception is `regexp.MustCompile` and similar "must" variants, which are explicitly designed for initialization.

- **Log stack traces on recovery**: Use `debug.Stack()` or `log.Printf("panic: %v\n%s", r, debug.Stack())` to capture the stack trace for debugging. The panic's own stack trace is lost after recovery.

- **Recover in batch processing per-item**: If processing a batch of items, wrap each item in a closure with deferred recover so one bad item doesn't fail the entire batch.

## Performance implications

- `panic` and `recover` are not designed for performance. They unwind the stack, which involves iterating the defer list and calling functions.
- The cost of a panic-recover pair is dominated by stack unwinding. For deep call stacks (100+ frames), this can be milliseconds.
- Never use panic-recover for control flow. It is 100x-1000x slower than returning an error.
- Deferred functions with recover use the slower runtime defer path (not inlined). This adds a small overhead to the function even when no panic occurs.

## Practice task

Write a function `safeDivide(a, b int) (result int, err error)` that:
1. Panics if `b == 0`.
2. Recovers from the panic in a deferred function and returns an error `"division by zero"`.
3. Returns `a / b` normally.

Then write a batch processor `safeBatch(dividends, divisors []int) (results []int, err error)` that processes each pair, recovering individually from each division-by-zero panic so a single bad pair doesn't stop the batch.

## Tests / verification

```bash
go run ./curriculum/modules/04-functions-errors-data-semantics/lessons/17-panic-and-recover
go test ./curriculum/modules/04-functions-errors-data-semantics/lessons/17-panic-and-recover
```

## Review questions

1. What is the difference between a panic and an error in Go?
2. Where must `recover()` be called to have an effect?
3. After a successful `recover()`, where does execution continue?
4. Can a `recover()` in one goroutine catch a panic from another goroutine?
5. Why should library code avoid calling `panic`?

## NEXT UP

Custom error types
