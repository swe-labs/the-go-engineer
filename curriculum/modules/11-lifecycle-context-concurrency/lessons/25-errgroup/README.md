# Errgroup

## Learning objective

Use `golang.org/x/sync/errgroup` to run goroutines in a fan-out pattern, propagate the first error, and cancel remaining work via context.

## Why this matters

Many Go programs need to run several independent operations in parallel and stop all of them if any one fails: fetch multiple URLs, query several shards, or run health checks against a set of endpoints. Writing this pattern by hand requires a `sync.WaitGroup`, a shared error variable with a mutex, and a cancellation mechanism. The `errgroup` package packages all three into a clean API. It is the standard tool for fan-out error handling in Go.

## Mental model

An errgroup is a collection of goroutines that share a lifecycle. When you call `g.Go(f)`, the function `f` runs in a new goroutine. If any function returns a non-nil error, the errgroup cancels the shared context, which signals all other goroutines to stop. The call to `g.Wait()` blocks until all goroutines finish and returns the first error (or nil). Think of it as a `sync.WaitGroup` with built-in error collection and context cancellation.

## Core idea

The `errgroup.Group` type has three methods:

- `Go(func() error)`: starts a goroutine that runs the function. If any prior goroutine returned an error, the function may or may not run (the group does not prevent new goroutines from starting after an error, but the shared context is cancelled).
- `Wait() error`: waits for all goroutines to finish and returns the first non-nil error, or nil if all succeeded.
- `WithContext(ctx Context) (Group, Context)`: creates a group with a derived context that is cancelled when any goroutine returns an error.

The pattern is always: create the group, call `Go` for each task, then call `Wait`.

## Under the hood

`errgroup.Group` embeds a `sync.WaitGroup` internally. Each call to `Go` increments the WaitGroup counter. When the function returns, the counter is decremented. `Wait` calls `wg.Wait()` and returns the first stored error. The error is stored exactly once using `sync.Once`: the first non-nil error from any goroutine is captured, and subsequent errors are ignored. The `WithContext` function creates a derived context via `context.WithCancel` and stores the cancel function. If a function returns an error, the cancel function is called before storing the error, propagating cancellation to all goroutines that `select` on `ctx.Done()`.

## How Go uses it

The errgroup pattern is used throughout the Go ecosystem:

- `go test` runner: runs test packages in parallel and fails fast.
- Kubernetes client-go: fans out API requests to multiple API servers.
- HashiCorp Consul: runs health checks against all nodes in parallel.
- Build systems: compile packages in parallel and stop on first error.

Go's standard library does not include errgroup because the `golang.org/x/sync` package serves as the incubation area.

## Go example

```go
package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func fetchURL(ctx context.Context, url string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(50 * time.Millisecond):
		if url == "https://bad.example" {
			return fmt.Errorf("failed to fetch %s: 500", url)
		}
		return nil
	}
}

func main() {
	urls := []string{
		"https://example.com",
		"https://api.example.com",
		"https://bad.example",
		"https://cdn.example.com",
	}

	g, ctx := errgroup.WithContext(context.Background())

	for _, url := range urls {
		url := url // capture loop variable
		g.Go(func() error {
			return fetchURL(ctx, url)
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("Fetch failed: %v\n", err)
		return
	}
	fmt.Println("All fetches succeeded")
}
```

Run:

```bash
go run .
# Output: Fetch failed: failed to fetch https://bad.example: 500
```

## Step-by-step execution

1. `errgroup.WithContext(ctx)` creates a group and a derived cancellable context.
2. For each URL, `g.Go` starts a goroutine. All four goroutines run concurrently.
3. The goroutine for `bad.example` returns an error after 50 ms.
4. Inside `errgroup`, the error is stored and the cancel function for the derived context is called.
5. Other goroutines may still be in `time.After`. When the timer fires, they check `ctx.Done()` (in this simplified example they had already passed the select). In a real fetch, the HTTP client would see the cancelled context.
6. `g.Wait()` blocks until all goroutines finish, then returns the stored error.
7. `main` prints the error.

## Common mistakes

- Mistake: Using `g.Go` with a function that never returns (infinite loop).
  - Why it happens: `g.Wait()` waits for all goroutines to finish. An infinite loop causes `Wait` to block forever.
  - Fix: Ensure every `Go` function respects context cancellation via `select`.

- Mistake: Ignoring the returned context from `WithContext`.
  - Why it happens: The developer forgets to pass the derived context to the `Go` functions, so cancellation has no effect.
  - Fix: Always pass `ctx` (the derived context) to any function that supports cancellation.

- Mistake: Capturing the loop variable incorrectly.
  - Why it happens: The closure captures the loop variable by reference, so all goroutines see the last value.
  - Fix: Use `url := url` inside the loop, or use the Go 1.22+ loop variable semantics.

- Mistake: Using `errgroup` for tasks that should not cancel others on failure.
  - Why it happens: The first error cancels all remaining work, which may not be desired (e.g., collecting partial results).
  - Fix: Use `sync.WaitGroup` + a mutex-protected error slice, or use `errgroup` with a separate context that is not derived from the group context.

## Debugging walkthrough

Buggy program:

```go
g, _ := errgroup.WithContext(context.Background())
var result int
for i := 0; i < 10; i++ {
	g.Go(func() error {
		result += i // DATA RACE
		return nil
	})
}
g.Wait()
```

Symptom: wrong result, or race detector fires.

Investigation: `g.Go` runs functions in their own goroutines. Accessing `result` without synchronization is a data race. The loop variable capture `i` is also wrong (each goroutine sees the final `i`).

Fix:

```go
g, ctx := errgroup.WithContext(context.Background())
results := make([]int, 10)
for i := 0; i < 10; i++ {
	i := i
	g.Go(func() error {
		results[i] = i * 2 // each goroutine writes its own slot, no race
		return nil
	})
}
g.Wait()
```

## Production notes

Errgroup is the recommended pattern for any "fan-out, fail-fast" workload. Use it for health checks against upstream dependencies, parallel cache warming, batch RPC fan-out, and parallel file processing. The `errgroup` package also provides `SetLimit` to bound concurrency. In production, always pass a context with an overall timeout so the group cannot outlive the request: `ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second); g, ctx := errgroup.WithContext(ctx)`.

## Performance implications

Errgroup adds negligible overhead: one `sync.WaitGroup`, one `sync.Once`, and one `context.CancelFunc` per group. The performance cost is the cost of the goroutines themselves. Use `SetLimit` to bound the number of concurrent goroutines when the fan-out is large (1000+ tasks). Without a limit, errgroup spawns one goroutine per `Go` call, which can exhaust memory or hit system limits.

## Practice task

Write a function `parallelCheck(ctx context.Context, endpoints []string) error` that uses errgroup to check each endpoint by pinging it (simulate with a sleep). If any endpoint returns an error, the function should return immediately with that error and cancel the remaining checks. Use `SetLimit(5)` to limit concurrency to 5 goroutines. Write table-driven tests that cover: all endpoints pass, one endpoint fails, and context cancellation before completion.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/25-errgroup
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/25-errgroup
```

## Review questions

1. What does `errgroup.WithContext` return, and why do you need both values?
2. What happens to the other goroutines when one goroutine in an errgroup returns an error?
3. How does `sync.Once` factor into errgroup's implementation?
4. Why is capturing the loop variable important when calling `g.Go` inside a for loop?
5. What method limits the number of concurrent goroutines in an errgroup?

## NEXT UP

Bounded worker pools
