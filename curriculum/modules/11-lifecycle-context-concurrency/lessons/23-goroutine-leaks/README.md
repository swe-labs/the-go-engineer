# Goroutine leaks

## Learning objective

Identify why goroutines leak, detect leaks using runtime metrics and profiling, and apply prevention patterns with context cancellation and channel ownership.

## Why this matters

A goroutine leak is a permanent memory leak: the goroutine's stack (initially 2 KB, growing to megabytes after usage) is never freed, and any objects the goroutine references cannot be garbage collected. Leaked goroutines accumulate under load and cause out-of-memory crashes. According to postmortems from Uber, Docker, and Lyft, goroutine leaks are the most common cause of memory exhaustion in Go microservices after plain heap allocation.

## Mental model

Every goroutine you spawn must have a guaranteed exit path. If you cannot prove that a goroutine will eventually reach a `return` statement, it is a leak waiting to happen. The goroutine itself is an allocation (stack + metadata) that persists until the goroutine exits. A blocked goroutine is not freed — it sits in the scheduler's run queue in `_Gwaiting` state forever. The Go runtime tracks every goroutine in an internal slice (`allgs`), and there is no timeout, no GC for goroutines, and no automatic reaper.

## Core idea

A goroutine leaks when it enters a state from which it can never return. The three most common causes:

1. **Blocked send**: sending on a channel that no goroutine will ever receive from.
2. **Blocked receive**: receiving from a channel that no goroutine will ever send to (or close).
3. **Infinite loop**: a `for` or `select` that never receives a stop signal.

All three share the same root cause: the goroutine was started without a mechanism to tell it when to stop.

Detection techniques:

- `runtime.NumGoroutine()` in tests: assert the count returns to baseline.
- `net/http/pprof`: `go tool pprof http://localhost:6060/debug/pprof/goroutine`
- `/debug/pprof/goroutine?debug=2`: full stack dump of every goroutine.
- Prometheus `go_goroutines` metric.

## Under the hood

The Go runtime scheduler uses an M:P:G model. A leaked goroutine in `_Gwaiting` state is removed from its P's run queue and stored in a global list. The GC must still scan the goroutine's stack during every GC cycle to find roots, so leaked goroutines increase GC CPU usage even when they are "idle". The goroutine's stack starts at 2 KB and grows as needed; a leaked goroutine that has processed work may hold a multi-megabyte stack. The `allgs` slice is never compacted — it grows monotonically.

## How Go uses it

The standard library uses context cancellation and channel closing as the two primary mechanisms to stop goroutines. `net/http` server cancels the request context when the client disconnects. `database/sql` uses context cancellation to abort queries. The `os/signal` package uses a channel that the runtime closes on signal delivery.

## Go example

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func leakyWorker(done chan bool) {
	<-done // blocks forever if nothing sends
}

func safeWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func main() {
	// This goroutine leaks: nothing ever sends on done
	done := make(chan bool)
	go leakyWorker(done)
	fmt.Println("Leaky worker started (will leak)")

	// This goroutine exits cleanly via context
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	go safeWorker(ctx)

	time.Sleep(100 * time.Millisecond)
	fmt.Println("Safe worker exited via context cancellation")
	fmt.Println("Tip: use runtime.NumGoroutine() to detect leaks")
}
```

## Step-by-step execution

For the leaky worker:

1. `main` creates an unbuffered `done` channel.
2. `main` starts `leakyWorker` in a new goroutine. The goroutine executes `<-done` and blocks.
3. `main` continues and sleeps for 100 ms.
4. `main` exits (after the main function returns, the program terminates).
5. If the program ran as a long-lived service, the `leakyWorker` goroutine would block forever, its stack never freed.

For the safe worker:

1. `main` creates a cancellable context with a 50 ms timeout.
2. `main` starts `safeWorker`. The goroutine enters a `for/select`.
3. Every 10 ms, `safeWorker` takes the `default` case and sleeps.
4. After 50 ms, the context's deadline expires. `ctx.Done()` channel is closed.
5. `safeWorker` receives from `ctx.Done()` and returns. The goroutine exits cleanly.

## Common mistakes

- Mistake: Starting a goroutine in an HTTP handler without a shutdown signal.
  - Why it happens: The handler returns, but the goroutine continues running. If it blocks on a channel send, it leaks.
  - Fix: Use the request context: `select { case <-ctx.Done(): return; case ch <- data: }`.

- Mistake: Using `time.Sleep` as a cleanup mechanism.
  - Why it happens: The developer adds `time.Sleep(1 * time.Second)` before returning, assuming the goroutine will finish.
  - Fix: Use a `sync.WaitGroup` or a context to wait for the goroutine to complete.

- Mistake: Sending to a channel in a goroutine when the receiver may have stopped.
  - Why it happens: The sender does not know when the receiver stops listening.
  - Fix: Make the sender responsible for closing the channel, or use a done channel that both sides select on.

- Mistake: Starting goroutines inside a loop without limiting concurrency.
  - Why it happens: Each iteration spawns a goroutine, and if they block on an external resource, thousands pile up.
  - Fix: Use a bounded worker pool (see Lesson 26).

## Debugging walkthrough

Symptom: a Go service's memory grows until OOM. Prometheus `go_goroutines` shows a steady increase over time.

Investigation:

1. Fetch the goroutine profile: `wget http://localhost:6060/debug/pprof/goroutine?debug=2 -O stacks.txt`.
2. Count goroutines by state: `grep "goroutine" stacks.txt | wc -l` shows 5000+.
3. Search for `_Gwaiting` blocks: look for stack traces ending in `chan send` or `chan receive`.
4. Find the blocking line: a stack trace shows `main.sendMetric` blocking on `ch <- m`.
5. Root cause: the metric sender goroutine has no timeout on the send and the receiver fell behind. Under load, senders pile up.

Fix:

```go
select {
case ch <- m:
case <-ctx.Done():
}
```

## Production notes

Every exported function that starts a goroutine must document how that goroutine stops. In production, set a goroutine count alert in your monitoring system. For Go services, a sharp increase in `go_goroutines` often precedes an OOM kill by minutes. Use `net/http/pprof` in every binary (behind an authenticated endpoint) so you can inspect goroutine dumps live. The `runtime.NumGoroutine()` function is cheap and safe to export as a Prometheus gauge.

## Performance implications

A leaked goroutine's stack consumes memory that is never reclaimed. A goroutine blocked on a channel or mutex uses zero CPU but still occupies scheduler data structures. The `allgs` slice grows linearly with leaked goroutines, and GC scans each stack during every cycle. For a service with 10,000 leaked goroutines each holding a 64 KB stack, that is 640 MB of unreclaimable memory plus GC overhead for scanning all those stacks on every cycle.

## Practice task

Write a function `startLeakySender(ch chan<- int)` that leaks if no receiver reads from `ch`. Then write `startSafeSender(ctx context.Context, ch chan<- int)` that uses `select` with `ctx.Done()` to avoid the leak. Write a test that starts both, waits, and asserts that `runtime.NumGoroutine()` returns to baseline only for the safe version.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/23-goroutine-leaks
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/23-goroutine-leaks
```

## Review questions

1. What are the three most common causes of goroutine leaks?
2. How does `runtime.NumGoroutine()` help detect leaks?
3. What is the difference between a goroutine that is sleeping via `time.Sleep` and one that is blocked on a channel receive?
4. Why does a goroutine leak cause GC overhead even when the goroutine is blocked?
5. What two mechanisms does Go provide to stop goroutines safely?

## NEXT UP

Deadlocks
