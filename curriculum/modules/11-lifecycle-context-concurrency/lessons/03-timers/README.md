# Timers

## Learning objective

Create, stop, and reset one-shot timed events using `time.NewTimer`, `time.After`, `timer.Stop`, and `timer.Reset` in concurrent Go programs.

## Why this matters

Timers are the building block for timeouts, retries, heartbeat intervals, and delayed task execution. Every production Go service uses timers — in HTTP client timeouts, database connection retries, graceful shutdown windows, and rate-limit backoff. Misusing timers leaks goroutines or causes premature timeouts that cascade into system failures.

## Mental model

A timer is a one-shot alarm that fires exactly once after a specified duration. Think of it as a countdown that sends a pulse on a channel when it reaches zero. You can cancel the alarm before it rings (`Stop`), or rewind it to count down from a new duration (`Reset`). Once it rings, the channel holds exactly one tick — be sure to drain it before resetting.

## Core idea

`time.NewTimer(d)` creates a `*time.Timer` with a `C` channel (`<-chan time.Time`) that receives the current time when `d` elapses.

Key timer operations:

| Operation | Behavior |
|---|---|
| `time.NewTimer(d)` | Creates a timer that fires after `d` |
| `<-timer.C` | Blocks until timer fires (or returns immediately if already fired) |
| `timer.Stop()` | Prevents the timer from firing; returns `true` if stopped before firing |
| `timer.Reset(d)` | Changes the timer to fire after `d`; must call on stopped/expired timers |
| `time.After(d)` | Returns a channel directly, equivalent to `NewTimer(d).C` |

## Under the hood

Each `time.Timer` is managed by the Go runtime's timer heap (a min-heap keyed by absolute wake time). Internally, a `timer` struct contains the fire time, period (for tickers), the channel to send on, and a callback function. When `NewTimer` is called, the runtime inserts the timer into the heap. The scheduler's `sysmon` thread periodically checks the heap and fires timers whose deadline has passed by sending on their channel. `Stop` removes the timer from the heap. `Reset` removes and re-inserts with the new duration.

## How Go uses it

- `net/http`: `Transport.DialTimeout` uses `time.After` under the hood to limit connection establishment.
- `context`: `WithTimeout` creates a timer via `time.AfterFunc` that calls the cancel function.
- `database/sql`: Connection retry loops use timers for backoff intervals.
- `os/signal`: Signal notification uses a channel pattern analogous to timers.

## Go example

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	timer := time.NewTimer(1 * time.Second)
	start := time.Now()
	<-timer.C
	fmt.Println("Timer fired after:", time.Since(start).Round(time.Millisecond))

	timer2 := time.NewTimer(500 * time.Millisecond)
	stopped := timer2.Stop()
	if stopped {
		fmt.Println("Timer stopped before firing")
	}

	timer3 := time.NewTimer(1 * time.Second)
	go func() {
		<-timer3.C
		fmt.Println("Timer 3 fired")
	}()
	time.Sleep(100 * time.Millisecond)
	reset := timer3.Reset(2 * time.Second)
	fmt.Println("Timer 3 reset, was active:", reset)
	time.Sleep(2500 * time.Millisecond)

	select {
	case <-time.After(100 * time.Millisecond):
		fmt.Println("time.After fired")
	}
}
```

## Step-by-step execution

1. `time.NewTimer(1 * time.Second)` allocates a `Timer` struct and inserts it into the runtime heap with a deadline 1s from now.
2. `<-timer.C` blocks the main goroutine. The scheduler parks it. After 1s, `sysmon` or the P's timer goroutine sends `time.Now()` on `timer.C`. The main goroutine is made runnable again.
3. `timer2.Stop()` removes the timer from the heap before it fires, preventing the channel send. Returns `true` because the timer was active.
4. `timer3.Reset(2 * time.Second)` removes and re-inserts `timer3` with a new deadline 2s from now. Returns `true` because the timer had not yet fired.
5. `time.After(100 * time.Millisecond)` is shorthand: `NewTimer(d).C`. It creates a timer and immediately returns its channel.

## Common mistakes

- Mistake: Calling `Reset` on an active timer without draining the channel.
  - Why: If the timer already fired, its channel holds a stale value. `Reset` followed by `<-timer.C` returns immediately with the old value.
  - Fix: Drain the channel before resetting: `if !timer.Stop() { <-timer.C }`.

- Mistake: Forgetting to stop a timer after use.
  - Why: The timer remains in the runtime heap until it fires, keeping resources alive and preventing GC.
  - Fix: Always defer `timer.Stop()` when the timer lifetime is tied to a function scope.

- Mistake: Using `time.After` in a loop that recreates it every iteration.
  - Why: Each iteration allocates a new timer that remains in the heap until it fires.
  - Fix: Use `time.NewTimer` once and call `Reset` after draining.

- Mistake: Reading from `timer.C` after `Stop` returns `false`.
  - Why: `Stop` returning `false` means the timer already fired. The channel contains a value that must be drained.
  - Fix: Always drain when `Stop` returns `false`: `if !timer.Stop() { <-timer.C }`.

## Debugging walkthrough

Consider a HTTP client timeout that is not respected:

```go
func fetchWithTimeout(url string, timeout time.Duration) {
    timer := time.NewTimer(timeout)
    defer timer.Stop()
    go func() {
        // simulate HTTP call
        time.Sleep(2 * time.Second)
        fmt.Println("response received")
    }()
    <-timer.C
    fmt.Println("timeout")
}
```

**Symptom**: The function always prints "timeout" immediately, even before the timeout duration.

**Investigation**: Add a timestamp at the timer creation and at the channel read.

```go
fmt.Println("timer created at", time.Now().Format(time.StampMilli))
<-timer.C
fmt.Println("timer fired at", time.Now().Format(time.StampMilli))
```

The timer fires almost instantly.

**Root cause**: The goroutine is launched but the main goroutine immediately blocks on `<-timer.C`. The timer is created correctly, but nothing is printed before it — the problem is that the caller may pass a very small or negative timeout. With a negative timeout, `NewTimer` creates a timer that fires immediately.

**Fix**: Validate the timeout before creating the timer:

```go
if timeout <= 0 {
    timeout = 5 * time.Second // default
}
```

## Production notes

- Always use `defer timer.Stop()` to ensure cleanup, even on the happy path.
- Use `time.After` for simple `select` cases where the timer is used once and discarded.
- Use `time.NewTimer` with `Reset` when the same timer is reused (e.g., retry loops with backoff).
- Never reuse a timer without draining its channel first; use the helper pattern:
  ```go
  if !timer.Stop() {
      select {
      case <-timer.C:
      default:
      }
  }
  timer.Reset(newDuration)
  ```
- Timers fire from the runtime's system goroutine — keep the channel receiver responsive.

## Performance implications

- `time.NewTimer` allocates a `Timer` struct on the heap.
- Active timers consume runtime heap space — many concurrent timers increase scheduler overhead.
- `timer.Stop` is an O(log n) operation on the timer heap.
- `timer.Reset` is O(log n) remove + O(log n) insert.
- `time.After` leaks if the timer fires after the `select` chooses another case — the timer goroutine survives until the duration elapses. Use `NewTimer` + `Stop` in long-lived selects.

## Practice task

Write a function `delayedPrint(msg string, delay time.Duration, cancel <-chan struct{}) bool` that:
- Waits for `delay`, then prints `msg` and returns `true`.
- If `cancel` is closed before `delay` expires, returns `false` without printing.
- Uses `time.NewTimer` and a `select` over both the timer and `cancel`.

Then write a main function that tests both the success and cancellation paths.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/03-timers
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/03-timers
```

The existing tests verify timer fire, stop, reset, and `time.After`. After the practice task, add table-driven tests for `delayedPrint` covering the success case, cancellation case, and zero-delay edge case.

## Review questions

1. What does `timer.Stop()` return if the timer has already fired?
2. Why must you drain `timer.C` before calling `timer.Reset()`?
3. What is the difference between `time.After(d)` and `time.NewTimer(d).C`?
4. What happens if you call `time.NewTimer(-1 * time.Second)`?
5. How does the Go runtime track pending timers internally?

## NEXT UP

Tickers — periodic events that fire repeatedly at a fixed interval.
