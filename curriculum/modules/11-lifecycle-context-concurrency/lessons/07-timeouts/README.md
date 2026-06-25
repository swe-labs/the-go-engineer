# Timeouts

## Learning objective

Apply and propagate timeouts using `context.WithTimeout`, detect `DeadlineExceeded` errors, and implement timeout hierarchies that prevent cascading failures.

## Why this matters

Without timeouts, a slow dependency stalls your entire service. A database query that hangs for 5 minutes blocks a goroutine, holds a connection pool slot, and eventually exhausts server resources. Timeouts bound the damage: they convert an indefinite hang into a predictable error. Every I/O call in production must have a timeout. Go's `context.WithTimeout` provides the standard way to enforce them consistently across the call graph.

## Mental model

A timeout is a deadline measured from "now." `context.WithTimeout(ctx, 5*time.Second)` is equivalent to `context.WithDeadline(ctx, time.Now().Add(5*time.Second))`. Think of it as placing a bomb under the operation: if the work does not finish before the fuse burns down, the context explodes and the operation fails with `DeadlineExceeded`. The timeout propagates through every function that receives the context — you set it once at the top level and every downstream call respects it.

## Core idea

`context.WithTimeout(parent, d)` creates a context that automatically cancels after duration `d`. The returned context's `Done()` channel closes when the timeout fires. `ctx.Err()` returns `DeadlineExceeded`.

| Function | Behavior |
|---|---|
| `context.WithTimeout(ctx, d)` | Creates a child context cancelled after `d` |
| `ctx.Err() == context.DeadlineExceeded` | True after timeout fires |
| `errors.Is(err, context.DeadlineExceeded)` | Preferred check (handles wrapped errors) |

Timeout propagation: when a parent context has a timeout, derived contexts inherit it. A child context can tighten (but not loosen) the timeout via its own `WithTimeout`.

## Under the hood

`context.WithTimeout` creates a `timerCtx` that embeds `cancelCtx`. A `time.Timer` is started with `time.AfterFunc(d, cancel)`. When the timer fires, the internal `cancel()` is called, which:
1. Closes the `Done()` channel.
2. Sets `err` to `DeadlineExceeded`.
3. Cancels all child contexts.

If manual `cancel()` is called before the timer fires, it stops the timer (preventing a wasted heap entry) and sets `err` to `Canceled` instead of `DeadlineExceeded`.

## How Go uses it

- **HTTP servers**: `http.Server.Timeout` sets a handler-level timeout for each request.
- **Database queries**: `db.QueryContext(ctx, ...)` cancels the query when the context times out, returning the connection to the pool.
- **gRPC clients**: Each RPC call accepts a context; the client cancels the call if the context expires before the response arrives.
- **Kubernetes client-go**: API requests use context timeouts to fail fast when the apiserver is slow.
- **Health checks**: Probes use short timeouts (e.g., 1s) to detect unresponsive dependencies quickly.

## Go example

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func slowOperation(ctx context.Context) error {
	select {
	case <-time.After(2 * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	err := slowOperation(ctx)
	if errors.Is(err, context.DeadlineExceeded) {
		fmt.Println("Operation timed out:", err)
	} else if err != nil {
		fmt.Println("Operation failed:", err)
	} else {
		fmt.Println("Operation succeeded")
	}
}
```

## Step-by-step execution

1. `context.WithTimeout(context.Background(), 300*time.Millisecond)` creates a `timerCtx` with a 300ms timer.
2. `slowOperation` is called with the timeout context.
3. Inside `slowOperation`, `select` waits for either `time.After(2*time.Second)` or `<-ctx.Done()`.
4. After 300ms, the timer fires. The internal `cancel()` closes the `Done()` channel and sets `Err()` to `DeadlineExceeded`.
5. The `ctx.Done()` case in `select` becomes ready. `slowOperation` returns `ctx.Err()` = `DeadlineExceeded`.
6. `main` checks `errors.Is(err, context.DeadlineExceeded)` and prints "Operation timed out".
7. `defer cancel()` runs. Since the timer already fired, calling `cancel()` again is a no-op.

If `slowOperation` had completed before 300ms, the `defer cancel()` would stop the timer, preventing it from firing unnecessarily.

## Common mistakes

- Mistake: Setting a timeout but never checking `ctx.Done()` in blocking operations.
  - Why: The timeout fires and closes the channel, but if no goroutine is selecting on it, the operation continues indefinitely.
  - Fix: Always select on `ctx.Done()` in any operation that blocks.

- Mistake: Creating a new timeout context per function call instead of propagating the parent timeout.
  - Why: Each `WithTimeout` in the call stack creates a new timer, fragmenting the timeout. A child may have a longer timeout than the parent, which is meaningless because the parent cancels first.
  - Fix: Accept a context from the caller and derive from it only when you need a stricter deadline.

- Mistake: Using `ctx.Err()` instead of `errors.Is` to check for `DeadlineExceeded`.
  - Why: `ctx.Err()` returns the raw sentinel. But a function might wrap it: `fmt.Errorf("query failed: %w", ctx.Err())`. Then `ctx.Err() == DeadlineExceeded` fails.
  - Fix: Always use `errors.Is(err, context.DeadlineExceeded)`.

- Mistake: Forgetting to call the cancel function returned by `WithTimeout`.
  - Why: The timer remains in the runtime heap until it fires, even if the operation completes early. This leaks memory and CPU.
  - Fix: Always `defer cancel()` immediately after creating the timeout context.

## Debugging walkthrough

Consider a service that occasionally returns "context deadline exceeded" for fast queries:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    result, err := queryDB(ctx, "SELECT ...")
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    fmt.Fprint(w, result)
}

func queryDB(ctx context.Context, q string) (string, error) {
    ctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
    defer cancel()
    // actual query...
    time.Sleep(20 * time.Millisecond)
    return "ok", nil
}
```

**Symptom**: Under load, `queryDB` returns `DeadlineExceeded` even though the query takes only 20ms.

**Investigation**: Check the parent context's deadline. Add logging:

```go
if d, ok := ctx.Deadline(); ok {
    log.Printf("parent deadline: %v, time until: %v", d, time.Until(d))
}
```

Under load, the parent context (from `http.Request`) already has a short deadline because the HTTP server's `Timeout` is set to 100ms. By the time `queryDB` runs, only 40ms remain. `queryDB` derives a 50ms timeout from this — but the parent context cancels first (after 40ms), causing `DeadlineExceeded`.

**Root cause**: `WithTimeout` on an already-short-lived parent context does not extend the parent's deadline. The parent's deadline is the hard upper bound.

**Fix**: Use a fixed timeout independent of the parent, or check how much time remains before deriving:

```go
deadline := time.Now().Add(50 * time.Millisecond)
if parentDeadline, ok := ctx.Deadline(); ok && parentDeadline.Before(deadline) {
    deadline = parentDeadline
}
childCtx, cancel := context.WithDeadline(ctx, deadline)
```

## Production notes

- Set a global default timeout for all outgoing HTTP calls using `http.Client.Timeout`.
- Use `context.WithTimeout` at the top of HTTP handlers to ensure all downstream calls share the same deadline.
- Log `DeadlineExceeded` at a lower severity than unexpected errors — they indicate slow dependencies, not bugs.
- In microservices, propagate client-set timeouts via gRPC deadline propagation or HTTP `Grpc-Timeout` headers.
- Monitor timeout rates per endpoint to detect degrading dependencies.

## Performance implications

- `WithTimeout` creates a `timerCtx` and starts a timer. Timer allocation and heap insertion cost roughly 1–2 microseconds.
- Cancelling a timeout context before the timer fires stops the timer, avoiding unnecessary heap wakeups.
- Checking `errors.Is(err, context.DeadlineExceeded)` is fast (pointer comparison after unwrapping).
- Each pending timer consumes runtime heap memory — do not create millions of short-lived timeout contexts per second.

## Practice task

Write a function `fetchWithTimeout(url string, timeout time.Duration) (string, error)` that:
- Creates a context with the given timeout using `context.WithTimeout`.
- Simulates an HTTP request by selecting on `time.After` (random duration between 100-500ms) and `ctx.Done()`.
- If the timeout fires first, returns `"", context.DeadlineExceeded`.
- If the simulated request completes first, returns `"response", nil`.

Write a main function that calls `fetchWithTimeout` with a 150ms timeout and verifies the error is `DeadlineExceeded`.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/07-timeouts
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/07-timeouts
```

The existing tests verify timeout exceeded, timeout respected, cancel-before-expiry, and negative-timeout behavior. After the practice task, add tests for `fetchWithTimeout` covering both the success and timeout paths.

## Review questions

1. What is the relationship between `context.WithTimeout` and `context.WithDeadline`?
2. What does `ctx.Err()` return after a timeout fires?
3. Why should you use `errors.Is(err, context.DeadlineExceeded)` instead of `err == context.DeadlineExceeded`?
4. What happens if you derive a 10-second timeout from a context that already expires in 5 seconds?
5. Why must you always `defer cancel()` after calling `context.WithTimeout`?

## NEXT UP

Deadlines — absolute time limits using `context.WithDeadline` and how they differ from relative timeouts.
