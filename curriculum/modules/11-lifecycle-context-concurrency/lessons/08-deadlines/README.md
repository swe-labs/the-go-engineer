# Deadlines

## Learning objective

Apply absolute deadlines using `context.WithDeadline`, distinguish deadlines from relative timeouts, and use deadlines to enforce latency SLOs in HTTP servers and downstream calls.

## Why this matters

Relative timeouts (`WithTimeout`) are convenient, but some scenarios demand absolute deadlines: "this request must complete by 12:00:00 UTC" or "the database call must finish before the HTTP handler's remaining time expires." Absolute deadlines compose correctly across service boundaries — a downstream service can compute its remaining time from the deadline rather than guessing a timeout value. This is how gRPC and distributed tracing propagate latency budgets.

## Mental model

A deadline is an absolute point in time. `context.WithDeadline(ctx, t)` is like a fixed appointment: when the clock reaches `t`, the context cancels. A timeout is a deadline computed as `now + d`. `WithTimeout` is a convenience wrapper around `WithDeadline`. The key difference: deadlines compose. A downstream service can inspect the deadline, compute the remaining time, and decide if it can complete the work before the deadline expires.

## Core idea

`context.WithDeadline(parent, deadline time.Time)` cancels the context when the system clock reaches `deadline`. The context's `Deadline()` method returns the absolute time.

| Function | Behavior |
|---|---|
| `context.WithDeadline(ctx, t)` | Cancels when wall clock reaches `t` |
| `ctx.Deadline()` | Returns `(t, true)` if a deadline is set |
| `time.Until(deadline)` | Remaining time before deadline |

Comparing `WithTimeout` vs `WithDeadline`:

| WithTimeout | WithDeadline |
|---|---|
| `WithTimeout(ctx, 5*time.Second)` | `WithDeadline(ctx, time.Now().Add(5*time.Second))` |
| Relative to now | Absolute point in time |
| Cannot inspect remaining budget | `ctx.Deadline()` gives the absolute time |
| Composes poorly across services | Composes naturally across service boundaries |

## Under the hood

`context.WithDeadline` works identically to `WithTimeout` internally. It creates a `timerCtx` that wraps `cancelCtx` with a timer. The timer is set to fire at the specified absolute time. If the parent context already has an earlier deadline, the new context's deadline is capped to the parent's. This ensures that derived contexts never outlive their parent.

## How Go uses it

- **gRPC**: Each RPC has a deadline computed by the client. The server reads `ctx.Deadline()` to compute its remaining processing time.
- **HTTP servers**: `http.Server.BaseContext` sets a per-request deadline. Handlers read the remaining time to decide how long to wait for downstream calls.
- **Distributed tracing**: OpenTelemetry and Jaeger propagate deadlines in trace context. Each span checks the remaining time before making outgoing calls.
- **Kubernetes admission webhooks**: The API server sets a deadline on the webhook call. If the webhook does not respond in time, the API server fails the admission request.

## Go example

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func operationWithDeadline(ctx context.Context) error {
	deadline, ok := ctx.Deadline()
	if ok {
		fmt.Println("Deadline set:", deadline.Format(time.RFC3339))
		fmt.Println("Time until deadline:", time.Until(deadline).Round(time.Millisecond))
	} else {
		fmt.Println("No deadline set")
	}

	select {
	case <-time.After(1 * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	deadline := time.Now().Add(200 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	err := operationWithDeadline(ctx)
	if errors.Is(err, context.DeadlineExceeded) {
		fmt.Println("Deadline exceeded")
	} else if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Completed before deadline")
	}
}
```

## Step-by-step execution

1. `time.Now().Add(200 * time.Millisecond)` computes an absolute deadline 200ms from now.
2. `context.WithDeadline(context.Background(), deadline)` creates a `timerCtx`. The timer is set to fire at the deadline.
3. `operationWithDeadline` calls `ctx.Deadline()`, which returns the absolute deadline and `true`.
4. The select waits for either a 1-second sleep or context cancellation.
5. After 200ms, the timer fires. The context is cancelled with `DeadlineExceeded`.
6. The `ctx.Done()` case becomes ready. The function returns `DeadlineExceeded`.
7. `main` prints "Deadline exceeded".

If the operation completed before the deadline, `defer cancel()` would stop the timer.

## Common mistakes

- Mistake: Passing a deadline in the past.
  - Why: `WithDeadline` with a past time creates an immediately-cancelled context. The operation never runs.
  - Fix: Check `deadline.After(time.Now())` before creating the context, or design the code to handle immediate cancellation.

- Mistake: Using `WithDeadline` when you should use `WithTimeout`.
  - Why: Both work, but `WithTimeout` is clearer for relative durations. Use `WithDeadline` only when you have an absolute time (e.g., from a parent context or a configuration value).
  - Fix: Use `WithTimeout` for "wait 5 seconds" and `WithDeadline` for "finish by 12:00:00".

- Mistake: Ignoring the parent context's deadline.
  - Why: If the parent has an earlier deadline, your derived deadline has no effect — the parent cancels first.
  - Fix: Check the remaining time: `remaining := time.Until(parentDeadline)` and set a timeout shorter than `remaining`.

- Mistake: Comparing `deadline` with `time.Now()` manually to check if time is up.
  - Why: This duplicates what the context does. The context's `Done()` channel closes when the deadline is reached.
  - Fix: Always select on `ctx.Done()` and let the context handle timing.

## Debugging walkthrough

Consider a batch processor that consistently fails with `DeadlineExceeded`:

```go
func batchProcess(items []string) error {
    deadline := time.Now().Add(5 * time.Second)
    ctx, cancel := context.WithDeadline(context.Background(), deadline)
    defer cancel()

    for _, item := range items {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }
        if err := processItem(ctx, item); err != nil {
            return err
        }
    }
    return nil
}
```

**Symptom**: Even with a few items, `DeadlineExceeded` is returned after exactly 5 seconds.

**Investigation**: Add logging inside `processItem`:

```go
func processItem(ctx context.Context, item string) error {
    if deadline, ok := ctx.Deadline(); ok {
        log.Printf("Processing %s: %v remaining", item, time.Until(deadline))
    }
    // ...
}
```

The log shows that the first item consumes 4.5 seconds of the deadline. Each subsequent item has less than 500ms remaining.

**Root cause**: The deadline is set once at the start of `batchProcess`. Each `processItem` call reduces the remaining time. A slow first item starves the rest.

**Fix**: Either increase the deadline, process items concurrently, or set per-item sub-deadlines:

```go
itemCtx, itemCancel := context.WithTimeout(ctx, 2*time.Second)
defer itemCancel()
return processItem(itemCtx, item)
```

## Production notes

- Propagate deadlines from incoming requests to outgoing calls: `ctx.Deadline()` gives the downstream service the ability to compute its remaining time.
- gRPC automatically propagates deadlines. Set the client-side deadline with `grpc.WithTimeout` and read it on the server with `ctx.Deadline()`.
- In HTTP services, set `http.Server.WriteTimeout` and `ReadTimeout` — these act as hard deadlines for the entire request lifecycle.
- For batch jobs, use absolute deadlines: "this job must finish by 02:00:00" rather than "this job may run for 1 hour." Absolute deadlines survive restarts and rescheduling.
- Monitor the gap between deadline and actual completion time. A consistently small gap indicates the deadline is too tight.

## Performance implications

- `WithDeadline` has the same cost as `WithTimeout`: one `timerCtx` allocation and one timer heap insertion.
- `ctx.Deadline()` returns a stored `time.Time` — O(1), no allocation.
- `time.Until(deadline)` is a monotonic subtraction — ~30ns.
- Creating many short-lived deadline contexts in a tight loop adds GC pressure. Reuse contexts where possible.

## Practice task

Write a function `remaining(ctx context.Context) (time.Duration, bool)` that returns the time remaining before the context's deadline and whether a deadline is set. If no deadline is set, return `(0, false)`. If the deadline has passed, return `(0, true)`. Write a main function that creates a context with a 500ms deadline, sleeps 200ms, calls `remaining`, and verifies the remaining time is approximately 300ms.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/08-deadlines
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/08-deadlines
```

The existing tests verify deadline exceeded, not-exceeded, equivalence with timeout, past deadlines, and correct deadline return. After the practice task, add tests for `remaining` covering each case.

## Review questions

1. How does `context.WithDeadline(ctx, t)` differ from `context.WithTimeout(ctx, d)`?
2. What happens when you pass a deadline in the past to `WithDeadline`?
3. Why do absolute deadlines compose better across service boundaries than relative timeouts?
4. How do you retrieve the absolute deadline from a context?
5. If a parent context has a deadline of 5 seconds and you derive a child with a deadline of 10 seconds, when does the child actually expire?

## NEXT UP

Context values with caveats — carrying request-scoped data through the context tree with type-safe keys and understanding the limitations.
