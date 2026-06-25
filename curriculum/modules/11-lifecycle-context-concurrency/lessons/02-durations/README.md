# Durations

## Learning objective

Represent, compute, compare, and format spans of time using `time.Duration` and its associated constants, arithmetic, and conversion methods.

## Why this matters

Durations govern every time-bound operation in Go: HTTP timeouts, database connection pool lifetimes, cache TTLs, retry intervals, ticker periods, and context deadlines. Getting duration math wrong causes production outages — premature timeouts drop requests, and overly long timeouts exhaust goroutine pools. Go's `time.Duration` gives you nanosecond precision with type safety that prevents unit confusion.

## Mental model

A `time.Duration` is an `int64` that counts nanoseconds. The named constants `time.Nanosecond`, `time.Microsecond`, `time.Millisecond`, `time.Second`, `time.Minute`, and `time.Hour` are multipliers. Arithmetic uses ordinary integer math on the nanosecond count. The `String()` method automatically selects the most readable unit: `1h30m15s` instead of `5415000000000`.

## Core idea

`time.Duration` is defined as `type Duration int64`. The nanosecond is the base unit. All constants are derived:

```go
const (
    Nanosecond  Duration = 1
    Microsecond          = 1000 * Nanosecond
    Millisecond          = 1000 * Microsecond
    Second               = 1000 * Millisecond
    Minute               = 60 * Second
    Hour                 = 60 * Minute
)
```

Key operations on durations:

| Expression | Result |
|---|---|
| `2*time.Hour + 30*time.Minute` | `2h30m0s` |
| `5*time.Second - 1500*time.Millisecond` | `3.5s` |
| `10*time.Second / 3` | `3s` (integer division) |
| `3*time.Minute > 150*time.Second` | `true` |
| `d.Hours()`, `d.Minutes()`, `d.Seconds()` | `float64` conversions |
| `d.Round(time.Minute)`, `d.Truncate(time.Second)` | Rounding/truncation |

## Under the hood

The runtime stores `time.Duration` as a plain `int64`. Arithmetic compiles to single CPU instructions. `String()` iterates over a division table: it extracts hours (`d / Hour`), then minutes (`(d % Hour) / Minute`), then seconds, milliseconds, and nanoseconds, building the string representation. There is no heap allocation unless the duration exceeds the formatting threshold (very large values use `strconv.AppendInt`).

## How Go uses it

- `net/http`: `Timeout`, `IdleTimeout`, `ReadHeaderTimeout` on `http.Server`
- `context`: `WithTimeout(parent, dur)` wraps `WithDeadline(parent, time.Now().Add(dur))`
- `database/sql`: `SetConnMaxLifetime`, `SetConnMaxIdleTime`
- `time`: `Sleep(d)`, `NewTicker(d)`, `NewTimer(d)`, `After(d)`, `Tick(d)`
- `golang.org/x/time/rate`: `NewLimiter(rate, burst)` uses duration for rate calculation

## Go example

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	d := 2*time.Hour + 30*time.Minute + 15*time.Second
	fmt.Println("Duration:", d)
	fmt.Println("Hours:", d.Hours())
	fmt.Println("Minutes:", d.Minutes())
	fmt.Println("Seconds:", d.Seconds())
	fmt.Println("Milliseconds:", d.Milliseconds())

	a := 3 * time.Second
	b := 1500 * time.Millisecond
	fmt.Println("a + b:", a+b)
	fmt.Println("a - b:", a-b)
	fmt.Println("a > b:", a > b)
	fmt.Println("a / b:", a/b)

	const tick = 500 * time.Millisecond
	start := time.Now()
	time.Sleep(tick)
	fmt.Println("Slept for:", time.Since(start).Round(time.Millisecond))
}
```

## Step-by-step execution

1. `d := 2*time.Hour + 30*time.Minute + 15*time.Second` — each constant is `time.Duration`, multiplication and addition are int64 operations: `7200000000000 + 180000000000 + 15000000000 = 7395000000000`.
2. `d.Hours()` divides the int64 count by `int64(time.Hour)` and returns the float64 `2.5041666666666667`.
3. `a + b` computes `3000000000 + 1500000000 = 4500000000` (4.5s).
4. `a / b` performs integer division: `3000000000 / 1500000000 = 2`.
5. `time.Sleep(tick)` tells the scheduler to park the current goroutine for at least 500ms.
6. `time.Since(start).Round(time.Millisecond)` subtracts monotonic clocks, then rounds to the nearest millisecond.

## Common mistakes

- Mistake: Writing `time.Duration(5000)` expecting 5 seconds.
  - Why: The argument is in nanoseconds, so `5000` = 5 microseconds.
  - Fix: Write `5000 * time.Millisecond` or `5 * time.Second`.

- Mistake: Using `int64` directly for duration math.
  - Why: Loses type safety — adding `int64` to `time.Time` compiles without error but produces garbage.
  - Fix: Always use `time.Duration` constants and convert with `time.Duration(n) * time.Second`.

- Mistake: Dividing durations for ratios.
  - Why: Integer division truncates: `10*time.Second / 3` = `3s`, losing `1/3` of a second.
  - Fix: Use `d.Seconds()` to get float64 before division.

- Mistake: Passing a negative duration to `Sleep`.
  - Why: `time.Sleep` returns immediately for negative durations, which can cause busy loops.
  - Fix: Guard with `if d > 0 { time.Sleep(d) }`.

## Debugging walkthrough

Consider a rate limiter that goes silent:

```go
func rateLimit(events int, per time.Duration) {
    interval := per / time.Duration(events)
    for i := 0; i < events; i++ {
        fmt.Println("event", i)
        time.Sleep(interval)
    }
}
```

**Symptom**: When called with `rateLimit(100, 1*time.Second)`, the function completes instantly — no events are printed with delays.

**Investigation**: Add a print of `interval`:

```go
fmt.Println("interval:", interval)  // interval: 10ms
```

That looks correct. But then check the caller — maybe it passes `int64` without conversion:

```go
rateLimit(100, 1000) // 1000 * Nanosecond = 1 microsecond!
```

**Root cause**: The caller supplied a bare integer `1000` instead of `1 * time.Second`. Go treats `1000` as `time.Duration(1000)` = 1000 nanoseconds = 1 microsecond.

**Fix**: Always use named duration constants at call sites: `rateLimit(100, 1*time.Second)`.

## Production notes

- Store durations as `time.Duration` (not `int64` for milliseconds) in config structs; use `time.ParseDuration` for env-var parsing.
- For JSON serialization, use a custom type that marshals as `"1s"` (string) rather than raw nanoseconds.
- Round durations before logging to avoid overwhelming readers with nanosecond precision: `d.Round(time.Millisecond)`.
- Beware of `time.Duration` overflow: `math.MaxInt64` nanoseconds is ~292 years. For epoch-offset times, prefer `time.Time`.
- Use `d / time.Millisecond` to convert duration to millisecond int64 for legacy systems.

## Performance implications

- Duration arithmetic is integer math — sub-nanosecond cost.
- `d.String()` allocates a string; cache the result if logging the same duration repeatedly.
- `d.Round()` and `d.Truncate()` do not allocate (they return a new `time.Duration`).
- `time.Since(start)` is ~30ns total (two monotonic reads + subtraction).
- `time.Sleep(d)` parks the goroutine and yields the OS thread — no CPU burn.

## Practice task

Write a function `humanDuration(d time.Duration) string` that returns a human-readable string like `"2h 30m 15s"`, omitting zero components. For example, `1h0m0s` should produce `"1h"` and `0s` should produce `"0s"`. Handle edge cases: negative durations return the absolute value with a `"-"` prefix. Use integer division and modulo on the nanosecond count.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/02-durations
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/02-durations
```

The existing tests verify duration arithmetic, comparison, and conversion. After the practice task, add table-driven tests covering `humanDuration` with cases for full duration, zero components, and negative values.

## Review questions

1. What is the underlying type of `time.Duration`, and what unit does it count in?
2. Why does `5*time.Second / 2` produce `2s` instead of `2.5s`?
3. If `d = 90 * time.Second`, what do `d.Round(time.Minute)` and `d.Truncate(time.Minute)` produce?
4. Write the expression that converts 2500 milliseconds to a `time.Duration`.
5. What happens if you pass `-1 * time.Second` to `time.Sleep`?

## NEXT UP

Timers — one-shot events that fire after a specified duration.
