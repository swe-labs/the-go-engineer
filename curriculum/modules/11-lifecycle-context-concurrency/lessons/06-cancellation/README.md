# Cancellation

## Learning objective

Use `context.WithCancel` to propagate cancellation signals across goroutines, implement graceful shutdown patterns, and prevent goroutine leaks by listening on `ctx.Done()`.

## Why this matters

Every goroutine you spawn must have a plan for how it stops. Without cancellation, a simple deployment rollback can leave thousands of orphaned goroutines consuming memory, holding database connections, and burning CPU. Cancellation is the primary mechanism Go provides for telling goroutines "stop what you are doing and clean up." It is essential for graceful shutdown, request deadlines, and user-initiated abort.

## Mental model

Context cancellation is a tree broadcast. The root context branches into child contexts — when any parent is cancelled, all children receive the cancellation. The `Done()` channel is like a fire alarm: once it rings, it keeps ringing, and every goroutine in the building should have a plan for what to do when it hears it.

## Core idea

`context.WithCancel(parent)` returns a derived context and a `cancel` function. Calling `cancel()` does two things:
1. Closes the derived context's `Done()` channel.
2. Recursively cancels all contexts derived from it.

Goroutines should select on `ctx.Done()` in any blocking operation:

```go
select {
case <-ctx.Done():
    return ctx.Err()
case result := <-workCh:
    return result
}
```

Cancellation propagates through the call graph by passing the context to every function that does I/O or spawns goroutines.

## Under the hood

`WithCancel` creates a `cancelCtx` struct. Internally, this struct maintains:
- A `Done` channel (created lazily on first access)
- A children map of contexts derived from it
- A reference to its parent

When `cancel()` is called:
1. The `Done` channel is created (if nil) and immediately closed.
2. `err` is set to `Canceled`.
3. Each child in the children map is cancelled recursively.
4. The `cancelCtx` is removed from its parent's children map.

The closing of a channel in Go is the only channel operation that broadcasts to all receivers — every goroutine blocked on `<-ctx.Done()` wakes up simultaneously.

## How Go uses it

- `net/http`: When a client disconnects, the server cancels the request context. All handlers, middleware, and downstream calls observe the cancellation.
- `database/sql`: `QueryContext` cancels the in-progress database query, returning the connection to the pool.
- `golang.org/x/sync/errgroup`: The first goroutine that returns a non-nil error cancels the group context, signalling all other goroutines to stop.
- Kubernetes controllers: Informer event handlers use context cancellation to stop worker goroutines during controller shutdown.

## Go example

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d shutting down: %s\n", id, ctx.Err())
			return
		default:
			fmt.Printf("Worker %d working...\n", id)
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	for i := 1; i <= 3; i++ {
		go worker(ctx, i)
	}

	time.Sleep(500 * time.Millisecond)
	fmt.Println("Cancelling workers...")
	cancel()

	time.Sleep(100 * time.Millisecond)
	fmt.Println("All workers cancelled")
}
```

## Step-by-step execution

1. `context.WithCancel(context.Background())` creates a `cancelCtx` with a nil `Done` channel.
2. Three worker goroutines start, each entering an infinite loop. In each iteration, they check `ctx.Done()` via select. Since `Done()` is nil initially, the `case <-ctx.Done()` branch blocks forever — the `default` case runs.
3. Each worker prints a message and sleeps 200ms.
4. After 500ms, `cancel()` is called. The `Done` channel is created and closed. `Err()` is set to `Canceled`.
5. All three workers' `select` statements now find the `ctx.Done()` case ready. They print their shutdown message and return.
6. The `cancel()` call returns, and after 100ms the program prints "All workers cancelled" and exits.

If `cancel()` were never called, the workers would run forever — a goroutine leak.

## Common mistakes

- Mistake: Calling `context.WithTimeout` or `context.WithCancel` in every function instead of accepting a context from the caller.
  - Why: Creates fragmented cancellation trees that are impossible to trace. A timeout set deep in the call stack overrides the parent's intent.
  - Fix: Accept `ctx context.Context` as the first parameter and derive from it only when you need a tighter deadline.

- Mistake: Ignoring the return value of `WithCancel` and using the parent context in goroutines.
  - Why: The parent context's `Done()` is never closed. Goroutines never receive cancellation.
  - Fix: Always use the derived context returned by `WithCancel`.

- Mistake: Checking `ctx.Err()` after the fact but never selecting on `ctx.Done()` in a goroutine.
  - Why: `Err()` only returns non-nil after cancellation. If the goroutine blocks on something else (e.g., a channel receive) without a select on `Done()`, it never learns about cancellation.
  - Fix: Always use `select { case <-ctx.Done(): ...; case ... }` when blocking.

## Debugging walkthrough

Consider a server that does not shut down cleanly:

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    go handleRequests(ctx)
    // wait for SIGINT
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT)
    <-sig
    fmt.Println("shutting down...")
    cancel()
    time.Sleep(3 * time.Second) // wait for goroutines
    fmt.Println("exiting")
}
```

**Symptom**: The program exits immediately after `cancel()` without waiting for goroutines to finish.

**Investigation**: Add a print inside `handleRequests` to see if it observes cancellation:

```go
func handleRequests(ctx context.Context) {
    <-ctx.Done()
    fmt.Println("handleRequests: context cancelled")
    time.Sleep(2 * time.Second) // cleanup
    fmt.Println("handleRequests: cleanup done")
}
```

The second print never appears.

**Root cause**: `cancel()` returns immediately. It closes the channel and sets the error, but it does not wait for goroutines to acknowledge. The `time.Sleep` in main is insufficient.

**Fix**: Use a `sync.WaitGroup` or a dedicated channel to wait for goroutine completion:

```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    handleRequests(ctx)
}()
<-sig
cancel()
wg.Wait() // waits for handleRequests to return
```

## Production notes

- Always defer the cancel function: `defer cancel()` — otherwise, derived contexts remain tracked in the parent and leak memory.
- Use `errgroup` for groups of goroutines where the first error should cancel all others.
- In HTTP servers, call `server.Shutdown(ctx)` which waits for active connections to complete within a deadline.
- Set a hard deadline for graceful shutdown: if cleanup takes too long, exit anyway. Use `context.WithTimeout` for the shutdown context.
- Monitor `context.Canceled` errors in your telemetry — they indicate cancelled operations, not bugs. Differentiate from `DeadlineExceeded`.

## Performance implications

- `WithCancel` allocates a `cancelCtx` (~64 bytes) and, when `Done()` is first accessed, a channel.
- Cancelling a context with many children iterates over a map — O(n) where n is the number of children.
- Selecting on `ctx.Done()` in a hot loop (without blocking) adds a branch per iteration. Use a non-blocking check with `select { case <-ctx.Done(): ...; default: }` to keep the overhead minimal.
- Context cancellation is designed for correctness, not throughput. Derive once, check in blocking calls.

## Practice task

Write a function `fanOut(ctx context.Context, n int, fn func(int))` that spawns `n` goroutines, each calling `fn(id)` in a loop until `ctx` is cancelled. Use `context.WithCancel` in `main()` to create a context, run `fanOut` for 3 workers, sleep 100ms, then cancel. Verify that all goroutines stop. Use `sync.WaitGroup` to ensure all goroutines complete before `main` exits.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/06-cancellation
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/06-cancellation
```

The existing tests verify that cancellation stops goroutines, the error is `Canceled`, cancellation broadcasts to multiple goroutines, and graceful shutdown completes. After the practice task, add tests for `fanOut` covering normal cancellation and immediate cancellation.

## Review questions

1. What happens to the `Done()` channel when `cancel()` is called?
2. Does cancelling a child context propagate to the parent?
3. What is the difference between `ctx.Err()` returning `Canceled` vs `DeadlineExceeded`?
4. Why must you always use the derived context from `WithCancel` rather than the parent?
5. How do you ensure all goroutines complete before `main` exits after cancellation?

## NEXT UP

Timeouts — using `context.WithTimeout` to limit operation duration and detect stalled operations.
