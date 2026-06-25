# Context basics

## Learning objective

Create, derive, and use Go's `context.Context` to carry cancellation signals, deadlines, and request-scoped values across goroutine boundaries.

## Why this matters

Every production Go service passes `context.Context` as the first parameter to virtually every function: HTTP handlers, database queries, RPC calls, and background workers. The context package provides the standard mechanism for propagating deadlines, cancellation, and request-scoped metadata through the call graph. Without it, goroutines leak, timeouts become inconsistent, and request tracing is impossible.

## Mental model

A context is a tree. The root is `context.Background()`, which never cancels and carries nothing. Derived contexts branch off with additional properties: `WithCancel` adds a cancel trigger, `WithTimeout` adds a time limit, and `WithValue` adds a key-value pair. Cancellation flows downward: when a parent cancels, all its children cancel. A cancelled child does not cancel its parent.

## Core idea

`context.Context` is an interface with four methods:

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key any) any
}
```

| Method | Returns |
|---|---|
| `Deadline()` | The time at which this context will be cancelled, if one is set |
| `Done()` | A channel that is closed when the context is cancelled or times out |
| `Err()` | `nil` while active, `Canceled` or `DeadlineExceeded` after cancellation |
| `Value(key)` | The value associated with key, or `nil` |

Root contexts:

| Function | Purpose |
|---|---|
| `context.Background()` | Root for all contexts; never cancelled, no values, no deadline |

## Under the hood

`context.Background()` returns a singleton `emptyCtx` struct. `context.WithCancel(parent)` creates a `cancelCtx` that embeds the parent and adds a `Done` channel (lazily created on first call to `Done()`). When `cancel()` is called, the channel is closed and the `cancelCtx` iterates over its children, cancelling each one. `context.WithTimeout` creates a `timerCtx` that wraps `cancelCtx` with a `time.Timer` — when the timer fires, it calls `cancel()`. Context trees are stored as linked lists; `Value(key)` walks up the tree until it finds a match.

## How Go uses it

- `net/http`: Every `http.Handler` receives a `context.Context` via `r.Context()`. The server cancels it when the connection closes.
- `database/sql`: `QueryContext`, `ExecContext`, `PrepareContext` all accept a context for cancellation and timeouts.
- `net`: Dialers accept contexts for connection timeouts.
- `golang.org/x/sync/errgroup`: Uses `WithCancel` to propagate the first error as cancellation.
- `google.golang.org/grpc`: Every RPC call receives a context for deadlines and metadata.

## Go example

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	fmt.Println("Background ctx err:", ctx.Err())
	fmt.Println("Background ctx done:", ctx.Done())

	ctxCancel, cancel := context.WithCancel(ctx)
	cancel()
	fmt.Println("Canceled ctx err:", ctxCancel.Err())

	ctxTimeout, timeoutCancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer timeoutCancel()
	<-ctxTimeout.Done()
	fmt.Println("Timeout ctx err:", ctxTimeout.Err())
}
```

## Step-by-step execution

1. `context.Background()` returns a singleton `emptyCtx`. `Err()` returns `nil` because it is never cancelled.
2. `context.WithCancel(ctx)` creates a `cancelCtx` with a lazily-initialized `Done` channel.
3. Calling `cancel()` closes the `Done` channel and sets `Err()` to `Canceled`. Any goroutine selecting on `<-ctxCancel.Done()` wakes up immediately.
4. `context.WithTimeout(ctx, 100*time.Millisecond)` creates a `timerCtx`. The runtime starts a background timer.
5. After 100ms, the timer fires, calling `cancel()` internally. The `Done` channel closes. `Err()` returns `DeadlineExceeded`.
6. `defer timeoutCancel()` ensures the underlying timer is stopped even if the timeout fires naturally, preventing the timer from lingering in the runtime heap.

## Common mistakes

- Mistake: Passing `nil` instead of a context.
  - Why: `nil` contexts cause panics when the code calls `ctx.Done()`. Always pass `context.Background()` as the default.
  - Fix: Use `context.Background()` for root contexts and always accept a context parameter from callers.

- Mistake: Passing `context.WithValue` results as if they were lightweight.
  - Why: `WithValue` creates a new context node. Deep chains of `WithValue` increase lookup cost and memory per request.
  - Fix: Keep value keys to a minimum (request IDs, auth tokens) and push configuration via explicit parameters.

- Mistake: Storing a context in a struct.
  - Why: The Go documentation explicitly warns against this. Contexts are request-scoped and should flow through function calls, not be stored.
  - Fix: Pass context as the first parameter to every function that needs it.

- Mistake: Ignoring the return value of `context.WithTimeout` / `context.WithCancel`.
  - Why: The returned context is the derived one; the parent context is unchanged. Using the parent context in goroutines means cancellation never reaches them.
  - Fix: Always use the derived context returned by these functions.

## Debugging walkthrough

Consider a database query that never times out:

```go
func queryDB(ctx context.Context) {
    ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
    defer cancel()
    // simulates a long DB call
    time.Sleep(1 * time.Second)
    fmt.Println("query done")
}

func main() {
    queryDB(context.Background())
    fmt.Println("main done")
}
```

**Symptom**: The function takes 1 second and does not abort after 100ms.

**Investigation**: Add a select over `ctx.Done()` inside `queryDB`:

```go
select {
case <-time.After(1 * time.Second):
    fmt.Println("query done")
case <-ctx.Done():
    fmt.Println("context cancelled:", ctx.Err())
}
```

**Root cause**: The timeout context is created but never checked. The `time.Sleep` does not observe context cancellation. The timeout fires and closes `ctx.Done()`, but nobody is listening.

**Fix**: Always select on `ctx.Done()` in any blocking operation, or pass the context to a library that supports it (like `sql.DB`).

## Production notes

- Accept `context.Context` as the first parameter of every function that makes I/O calls or spawns goroutines.
- Always defer the cancel function returned by `WithCancel` / `WithTimeout` to release resources.
- Use `context.Background()` in `main()`, `init()`, and tests as the root context.
- Propagate context through middleware, handler chains, and service boundaries.
- Never pass a `nil` context — use `context.Background()` instead.

## Performance implications

- `context.Background()` returns a singleton — zero allocation per call.
- `context.WithCancel` allocates a `cancelCtx` struct (~64 bytes) and, when `Done()` is first called, a channel.
- `context.WithTimeout` adds a timer allocation to the timer heap.
- `Value(key)` walks up the context tree — O(depth) lookup cost. Keep the depth small.
- Creating a new context per request is negligible; creating thousands per second in a hot path adds GC pressure.

## Practice task

Write a function `withTimeoutAndCancel(parent context.Context, d time.Duration) (ctx context.Context, cancel context.CancelFunc)` that returns a context that times out after `d` but also supports manual cancellation. Then write a main function demonstrating that calling `cancel()` before the timeout returns `context.Canceled` instead of `DeadlineExceeded`. Use `context.WithTimeout` and check the error after cancelling.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/05-context-basics
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/05-context-basics
```

The existing tests verify `Background`, `WithCancel`, `WithTimeout`, and cancellation behavior. After completing the practice task, add tests confirming that `cancel()` before timeout produces `Canceled` and not `DeadlineExceeded`.

## Review questions

1. What is the difference between `context.Background()` and a context derived with `WithCancel`?
2. What happens to child contexts when a parent context is cancelled?
3. Why should you not store a context in a struct?
4. What does `ctx.Done()` return and when is it closed?
5. If you call `cancel()` before a `context.WithTimeout` fires, what does `ctx.Err()` return?

## NEXT UP

Cancellation — using `context.WithCancel` to signal goroutines to stop and implement graceful shutdown.
