# Tickers

## Learning objective

Create, use, and stop periodic tickers with `time.NewTicker` and `time.Tick` for recurring tasks such as heartbeats, polling, and periodic cleanup.

## Why this matters

Production services need recurring actions: health-check loops, metrics emission, log rotation, cache eviction, and session cleanup. Tickers provide a clean channel-based mechanism for periodic work. Using raw `time.Sleep` in a loop drifts over time because the sleep duration does not account for the work duration. Tickers maintain a steady interval regardless of how long the receiver takes to process each tick.

## Mental model

A ticker is a recurring timer. Instead of firing once, it fires repeatedly at a fixed interval. Each tick sends the current time on its channel. Think of a metronome: it ticks at a steady beat regardless of how the musician plays between beats. If the receiver is slow and misses a tick, the ticker drops the missed tick — it does not queue them.

## Core idea

`time.NewTicker(d)` creates a `*time.Ticker` that sends the current time on `C` every `d`. The ticker runs until explicitly stopped with `ticker.Stop()`.

| Operation | Behavior |
|---|---|
| `time.NewTicker(d)` | Creates a ticker that ticks every `d` |
| `<-ticker.C` | Receives the next tick time |
| `ticker.Stop()` | Stops the ticker; closes no channels |
| `time.Tick(d)` | Returns `ticker.C` directly (no Stop ability — use with caution) |

## Under the hood

A `time.Ticker` uses the same runtime timer heap as `time.Timer`. The difference is the `period` field: when a ticker fires, the runtime does not remove it from the heap. Instead, it advances the deadline by `period` and re-inserts it. If the channel send would block (receiver is slow), the runtime skips the tick and continues to the next deadline. This means slow consumers miss ticks rather than queuing them.

## How Go uses it

- `net/http`: HTTP/2 ping frames use tickers to measure connection idle time.
- `runtime/metrics`: The runtime emits GC and goroutine metrics on a periodic ticker.
- Kubernetes client-go: Informer reflectors use tickers for resync periods.
- `database/sql`: Connection pool health checks use tickers to evict stale connections.
- Monitoring agents: Metrics exporters (Prometheus, OpenTelemetry) use tickers for collection intervals.

## Go example

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(200 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println("Tick at:", time.Now().Format("15:04:05.000"))
			<-ticker.C
		}
		ticker.Stop()
		done <- true
	}()

	<-done
	fmt.Println("Ticker stopped")
}
```

## Step-by-step execution

1. `time.NewTicker(200 * time.Millisecond)` creates a ticker. The runtime inserts it into the timer heap with a deadline 200ms from now and `period = 200ms`.
2. The goroutine loops: it prints a message, then blocks on `<-ticker.C`.
3. After 200ms, the runtime fires the ticker: it sends `time.Now()` on `ticker.C` and advances the deadline by 200ms.
4. The goroutine receives, prints, and loops back to block on the channel.
5. After 5 iterations, `ticker.Stop()` removes the ticker from the runtime heap entirely.
6. The goroutine signals `done`, and main exits.

If any iteration's processing took >200ms, the runtime would skip the next tick — the consumer would never receive a queued backlog. This prevents unbounded memory growth.

## Common mistakes

- Mistake: Using `time.Tick(d)` without a way to stop the ticker.
  - Why: `time.Tick` returns a channel but hides the `*time.Ticker`. The ticker leaks and runs forever, consuming heap space.
  - Fix: Use `time.NewTicker` and call `Stop` when done.

- Mistake: Assuming ticks are evenly spaced regardless of receiver speed.
  - Why: If the receiver takes longer than the tick interval, the ticker drops ticks. The interval between received ticks may vary.
  - Fix: Measure the actual interval with `time.Since(lastTick)` if you need to detect drift.

- Mistake: Reading from `ticker.C` after `Stop`.
  - Why: `Stop` does not close the channel. A goroutine blocked on `<-ticker.C` remains blocked indefinitely.
  - Fix: Use a separate done channel or `context.Context` with select.

- Mistake: Creating a ticker with a zero or negative duration.
  - Why: Panics with "non-positive interval for NewTicker".
  - Fix: Validate the interval before creating the ticker.

## Debugging walkthrough

Consider a metrics exporter that emits duplicate values:

```go
func emitMetrics() {
    ticker := time.NewTicker(1 * time.Second)
    for range ticker.C {
        fmt.Println("emit metric", time.Now().Unix())
    }
}
```

**Symptom**: Metrics are emitted twice per second, not once.

**Investigation**: Check if `emitMetrics` is launched more than once:

```go
// In main
go emitMetrics()
go emitMetrics() // duplicate!
```

There is no guard against multiple callers.

**Root cause**: The function is started in two goroutines, each creating its own ticker. Both tickers tick at 1s intervals, producing two emissions per second.

**Fix**: Use a singleton pattern or a sync mechanism:

```go
var once sync.Once
once.Do(func() { go emitMetrics() })
```

## Production notes

- Always pair `time.NewTicker` with `defer ticker.Stop()` to prevent leaks.
- Use `time.NewTicker` with `context.Context` for cancellation-aware periodic loops:
  ```go
  select {
  case <-ticker.C:
      // do work
  case <-ctx.Done():
      ticker.Stop()
      return
  }
  ```
- For heartbeats, use a short tick interval (e.g., 5s) and monitor the time since the last received tick to detect producer stalls.
- Ticker accuracy depends on system load; the runtime guarantees only that ticks are not delivered early. On an overloaded system, ticks may arrive late.
- In tests, use a small tick interval (1ms–10ms) to avoid slow tests. Never depend on real-time timing in unit tests.

## Performance implications

- Each active ticker occupies a slot in the runtime timer heap. Thousands of simultaneous tickers increase heap maintenance overhead.
- `ticker.Stop` is O(log n) on the timer heap.
- Receiving from `ticker.C` does not allocate (the tick time is already in the channel buffer).
- Tickers with sub-millisecond intervals may be limited by OS timer resolution (typically 1ms on Linux, ~15ms on Windows).
- A stopped ticker can be garbage collected immediately.

## Practice task

Write a function `heartbeat(ctx context.Context, interval time.Duration) <-chan time.Time` that returns a channel. Every `interval`, the channel receives the current time. When `ctx` is cancelled, the ticker is stopped and the channel is closed. Use `time.NewTicker` and a goroutine with a select over `ticker.C` and `ctx.Done()`. Write a main function that receives 3 beats, cancels the context, and verifies the channel closes.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/04-tickers
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/04-tickers
```

The existing tests verify ticker firing, multiple ticks, stopping, and zero-duration panic. After completing the practice task, add tests for `heartbeat` covering normal operation and cancellation.

## Review questions

1. How does `time.NewTicker` differ from `time.NewTimer`?
2. What happens if the consumer of `ticker.C` takes longer than the tick interval?
3. Why is `time.Tick(d)` considered unsafe for long-lived code?
4. Does `ticker.Stop()` close the `ticker.C` channel?
5. What happens if you pass a negative duration to `time.NewTicker`?

## NEXT UP

Context basics — the standard way to carry deadlines, cancellation signals, and request-scoped values across API boundaries.
