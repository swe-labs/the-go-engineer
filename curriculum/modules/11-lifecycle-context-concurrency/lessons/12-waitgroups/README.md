# Waitgroups

## Learning objective

Coordinate goroutine completion with `sync.WaitGroup` by correctly using `Add`, `Done`, and `Wait`, and apply the typed WaitGroup pattern for collecting results.

## Why this matters

Goroutines run independently. Without a coordination primitive, `main` or a parent function cannot know when spawned goroutines have finished. `time.Sleep` guesses are fragile — they fail under load or on slower machines. `sync.WaitGroup` provides a precise, efficient mechanism: block until a counter reaches zero. Every production Go service uses WaitGroups for graceful shutdown, batch processing, fan-out work, and ensuring background goroutines complete before the process exits.

## Mental model

A WaitGroup is a counter guarded by a condition variable. Think of it as a sign-up sheet:

- `Add(delta)` adds `delta` names to the sheet.
- `Done()` crosses off one name.
- `Wait()` blocks until the sheet is empty.

Multiple goroutines can call `Done` concurrently. The counter must reach zero exactly once. If `Add` and `Done` are mismatched, `Wait` never unblocks (counter never reaches zero) or a negative counter panics.

## Core idea

```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // work
}()
wg.Wait()
```

Contract:

- Call `wg.Add(delta)` **before** starting the goroutine, not inside it. This ensures `Wait` sees the correct count.
- Call `wg.Done()` exactly once per goroutine (usually via `defer`).
- `wg.Wait()` blocks the calling goroutine until the counter is zero.
- A `sync.WaitGroup` must not be copied after first use. Pass it by pointer.

Typed WaitGroup pattern — collect results into a slice from concurrent goroutines:

```go
results := make([]Result, n)
var wg sync.WaitGroup
for i := 0; i < n; i++ {
    wg.Add(1)
    go func(idx int) {
        defer wg.Done()
        results[idx] = fetch(idx)
    }(i)
}
wg.Wait()
// use results
```

Each goroutine writes to its own index, so no mutex is needed.

## Under the hood

`sync.WaitGroup` is a struct with three fields:

- `state` (uint64): packed fields — the counter (high 32 bits) and the number of goroutines waiting in `Wait` (low 32 bits).
- `sema` (uint32): a semaphore used to block and wake goroutines.

`Add(delta)` atomically adds `delta` to the counter. If the new counter is zero and waiters are positive, the semaphore is signaled to wake all waiters.

`Done()` calls `Add(-1)`.

`Wait()` atomically increments the waiter count, then blocks on the semaphore until the counter is zero. When `Add` from `Done` makes the counter zero and sees waiters > 0, it calls `runtime_Semrelease` which wakes all waiting goroutines.

The atomic operations on `state` ensure visibility across goroutines without additional barriers. The design avoids mutex locks in the common (non-contended) path.

## How Go uses it

- **Graceful shutdown**: `http.Server.Shutdown` waits on an internal WaitGroup that tracks in-flight requests.
- **Database connection pools**: pool maintenance goroutines use WaitGroups to coordinate cleanup.
- **Batch processing**: process a slice of items concurrently, collect results, then proceed.
- **Worker pools**: the dispatcher calls `wg.Add` per job, each worker calls `Done` when finished.
- **`testing` package**: `t.Run` subtests with `t.Parallel()` use internal WaitGroup coordination.

## Go example

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func fetch(id int) string {
	time.Sleep(time.Duration(id) * 100 * time.Millisecond)
	return fmt.Sprintf("result-%d", id)
}

func main() {
	ids := []int{1, 2, 3, 4, 5}
	results := make([]string, len(ids))

	var wg sync.WaitGroup
	for i, id := range ids {
		wg.Add(1)
		go func(idx, id int) {
			defer wg.Done()
			results[idx] = fetch(id)
		}(i, id)
	}

	wg.Wait()

	for _, r := range results {
		fmt.Println(r)
	}
}
```

Output (after ~500 ms total):

```
result-1
result-2
result-3
result-4
result-5
```

## Step-by-step execution

For 5 items with the typed WaitGroup pattern:

1. `main` creates `results` as a `[]string` of length 5, all zero values (`""`).
2. Loop `i=0..4`: `wg.Add(1)` increments counter to 1, 2, 3, 4, 5.
3. Each `go func()` captures its index and ID. The goroutine is enqueued on the run queue.
4. After the loop, `wg.Wait()` blocks `main`.
5. Each goroutine runs (order non-deterministic), calls `fetch(id)` which sleeps, then assigns to `results[idx]`.
6. After assignment, `defer wg.Done()` decrements the counter.
7. When the 5th goroutine calls `Done`, counter reaches 0. `Add` (via `Done`) signals the semaphore.
8. `main` unblocks from `Wait`, prints results, exits.

The key insight: the loop must finish calling `wg.Add` before any goroutine completes, otherwise a `Done` could make the counter negative. Starting goroutines inside the loop while `Add` happens before each `go` is safe because the goroutine's first chance to call `Done` is at its return, which is after the current iteration's `go` returns control to the loop.

## Common mistakes

- **Copying a WaitGroup**: `func doStuff(wg sync.WaitGroup)` copies the struct (including the internal counter). The copy's `Wait` may return immediately or race.
  - Fix: pass `*sync.WaitGroup`.
- **`Add` after `Wait` has started**: once `Wait` observes a positive counter and blocks, subsequent `Add` calls are not seen. `Wait` only unblocks when the counter it first observed reaches zero.
  - Fix: ensure all `Add` calls happen before calling `Wait`.
- **Wrong `Add` placement**: calling `wg.Add(1)` inside the goroutine risks a race where `Wait` starts before any `Add` executes, waits on zero, and returns too early.
  - Fix: call `Add` before `go`.
- **Negative counter**: calling `Done` more times than `Add` panics: `sync: negative WaitGroup counter`.
- **Reusing a WaitGroup for multiple batches**: forgetting to reset the counter. Either use a new `sync.WaitGroup` per batch, or ensure counter is zero before reusing.

## Debugging walkthrough

Consider this code that never finishes:

```go
var wg sync.WaitGroup
for i := 0; i < 5; i++ {
    go func() {
        wg.Add(1)
        defer wg.Done()
        fmt.Println(i)
    }()
}
wg.Wait()
```

**Symptom**: program hangs at `wg.Wait()`. It never prints anything or prints garbage values.

**Investigation**: the goroutines call `wg.Add(1)` after they start, but `wg.Wait()` in `main` runs immediately after the loop (before any goroutine runs) and sees counter = 0, so it returns instantly. The program exits before goroutines run.

**Root cause**: `Add` placed inside the goroutine instead of before `go`.

**Fix**:

```go
for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(n int) {
        defer wg.Done()
        fmt.Println(n)
    }(i)
}
wg.Wait()
```

## Production notes

- **Always use `defer wg.Done()`** immediately after the `go` line. If the goroutine panics, `defer` still runs and decrements the counter. Without `defer`, a panic causes a permanent hang in any caller waiting on `Wait`.
- **Combine with context cancellation**: in production, `wg.Wait()` blocks forever if a goroutine hangs. Wrap `Wait` with a `select` on `ctx.Done()`:

```go
c := make(chan struct{})
go func() {
    wg.Wait()
    close(c)
}()
select {
case <-c:
case <-ctx.Done():
}
```

- **Prefer typed WaitGroup for result collection** over using a mutex-protected slice append. Each goroutine writes to a reserved index; no locking needed.
- **For dynamic work counts**, use `sync.WaitGroup` with periodic `Add` calls as new work is discovered, but ensure `Add` is never called after `Wait` has returned.

## Performance implications

| Operation | Cost |
|---|---|
| `wg.Add(1)` | One atomic add (fast, ~5 ns) |
| `wg.Done()` | One atomic add (fast, ~5 ns) |
| `wg.Wait()` (contended) | One atomic add + semaphore park (~100 ns if already zero, ~µs if blocking) |
| Semaphore wake (N waiters) | OS-dependent (µs range) |

WaitGroup is one of the cheapest synchronization primitives in Go. In the uncontended case, `Add`/`Done` are single atomic instructions. The semaphore is only invoked when `Wait` actually blocks.

Avoid calling `wg.Add` with large deltas in a hot loop where the counter is never zero between iterations — this is fine in practice, but be aware that each `Add` is an atomic operation with store barrier.

## Practice task

Write a function `concurrentSum(nums []int, workers int) int` that:

1. Splits `nums` into `workers` equal chunks.
2. Launches `workers` goroutines, each summing its chunk into a shared results slice.
3. Uses `sync.WaitGroup` to wait for all workers.
4. Returns the total sum.
5. Tests with `nums = []int{1,2,3,4,5,6,7,8,9,10}` and `workers = 3`.

Verify the result is 55 regardless of the number of workers.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/12-waitgroups
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/12-waitgroups
```

## Review questions

1. Where must `wg.Add(1)` be called — before or inside the goroutine? Why?
2. What happens if you call `wg.Done()` more times than `wg.Add()`?
3. Can you pass a `sync.WaitGroup` by value to a function? Why or why not?
4. In the typed WaitGroup pattern, why is a mutex not needed when writing results from concurrent goroutines?
5. How would you add a timeout to `wg.Wait()`?

## NEXT UP

Channels — communicate between goroutines using typed message passing with `chan T`.
