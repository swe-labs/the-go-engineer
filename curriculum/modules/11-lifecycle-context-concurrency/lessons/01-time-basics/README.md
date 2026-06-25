# Time basics

## Learning objective

Use `time.Time`, `time.Now`, `time.Format`, `time.Parse`, `time.Duration`, `time.Since`, and `time.Until` to read, display, and manipulate timestamps in Go programs.

## Why this matters

Every production system tracks time: request logs, cache TTLs, retry backoff, session expiry, metrics, and scheduled jobs. Go's `time` package is the foundation for all of these. A professional Go engineer must be fluent in time representation, parsing, formatting, and duration math — bugs here silently corrupt data and break SLAs.

## Mental model

Think of `time.Time` as a struct with two internal clocks: a wall clock for calendar operations (year/month/day, timezone) and a monotonic clock for measuring elapsed time. The monotonic component is subtracted when you call `time.Since` or `time.Until`, so duration measurements remain accurate even if the system clock jumps due to NTP adjustments.

## Core idea

`time.Time` is Go's canonical type for an instant in time. Every `time.Time` value carries:
- Wall clock data (seconds since 1970-01-01 + nanosecond offset + location)
- Monotonic clock reading (for accurate duration measurement when created by `time.Now`)

Key functions:

| Function | Purpose |
|---|---|
| `time.Now()` | Current local time (wall + monotonic) |
| `t.Format(layout)` | Format as string using reference layout |
| `time.Parse(layout, s)` | Parse string into `time.Time` |
| `t.Unix()` | Seconds since Unix epoch |
| `time.Since(t)` | Shortcut for `time.Now().Sub(t)` |
| `time.Until(t)` | Shortcut for `t.Sub(time.Now())` |

## Under the hood

`time.Time` is a 24-byte struct: three `uint64` fields — wall (wall clock seconds + flag + nsec), ext (monotonic offset or extended seconds), and loc (pointer to time.Location). The first bit of wall distinguishes wall clock from monotonic. When `time.Now()` is called, the runtime reads `CLOCK_REALTIME` and `CLOCK_MONOTONIC` via platform syscalls (vdso on Linux). Parsing and formatting follow the reference time `Mon Jan 2 15:04:05 MST 2006` — each component maps to a specific number.

## How Go uses it

Go's standard library embeds `time.Time` everywhere:

- `net/http`: `Request.Header` timestamps, `Server.WriteTimeout`
- `database/sql`: `Scan` into `time.Time`, prepared statements with `datetime` columns
- `context`: `WithDeadline` and `WithTimeout` accept `time.Time` / `time.Duration`
- `log`: standard logger prefixes each line with `t.Format(time.DateTime)`
- `crypto/tls`: certificate validity windows using `time.Now()` comparisons

## Go example

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	fmt.Println("Now (Unix):", now.Unix())
	fmt.Println("Now (RFC3339):", now.Format(time.RFC3339))
	fmt.Println("Now (custom):", now.Format("2006-01-02 15:04:05"))

	layout := "2006-01-02"
	parsed, err := time.Parse(layout, "2026-06-01")
	if err != nil {
		panic(err)
	}
	fmt.Println("Parsed:", parsed.Format(time.RFC3339))

	then := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	fmt.Println("Since then:", time.Since(then))
	fmt.Println("Until next year:", time.Until(then.AddDate(1, 0, 0)))
}
```

## Step-by-step execution

1. `time.Now()` queries the kernel for the current wall clock and monotonic time.
2. `now.Unix()` extracts the integer seconds since epoch (no allocation).
3. `now.Format(time.RFC3339)` walks through the format string, mapping each token (`2006`, `01`, `02`, `T`, `15`, `04`, `05`, `Z07:00`) to the corresponding field of now.
4. `time.Parse(layout, "2026-06-01")` does the reverse: reads the string left to right, extracts year/month/day from positions matching the layout, and returns a `time.Time` with zero hours/minutes/seconds in UTC.
5. `time.Since(then)` subtracts the monotonic components of `then` from `now`, yielding a `time.Duration`. This is accurate even if the wall clock was adjusted between `then` and `now`.

## Common mistakes

- Mistake: Using `time.Now()` to measure elapsed time and comparing wall clock values.
  - Why: The wall clock can jump forward or backward due to NTP, daylight saving, or leap seconds.
  - Fix: Use `time.Since(start)` — it uses the monotonic clock.

- Mistake: Storing timestamps as Unix `int64`.
  - Why: Loses timezone information, human readability, and the monotonic component.
  - Fix: Store `time.Time` values; serialize with `time.RFC3339Nano` for interchange.

- Mistake: Comparing `time.Time` with `==`.
  - Why: `time.Time` has unexported fields including the monotonic reading and location pointer — two values representing the same instant may differ on these.
  - Fix: Use `t.Equal(other)` which compares only the wall clock instant.

- Mistake: Passing a timezone-unaware string to `time.Parse`.
  - Why: Without a timezone token in the layout, Parse assumes UTC.
  - Fix: Parse with `time.RFC3339` or specify `time.ParseInLocation`.

## Debugging walkthrough

Consider this buggy retry loop:

```go
func retry(attempts int) {
	start := time.Now()
	for i := 0; i < attempts; i++ {
		if time.Since(start) > 5*time.Second {
			fmt.Println("timeout")
			return
		}
		// do work
		time.Sleep(2 * time.Second)
	}
}
```

**Symptom**: The retry loop sometimes exits instantly on the first iteration with "timeout" even though no real time has passed.

**Investigation**: Print `start` monotonic component:

```go
fmt.Println(start.String()) // includes "m=+0.000000000"
```

If `start` was obtained via `time.Date(...)` (no monotonic clock), then `time.Since(start)` does not use monotonic subtraction — it falls back to wall-clock subtraction. If the system clock was adjusted backward after `retry` was called, `time.Since` can produce a large negative value, and the comparison `> 5*time.Second` evaluates true.

**Root cause**: `time.Date` does not attach a monotonic reading. `time.Since` on a non-monotonic time falls back to unreliable wall-clock arithmetic.

**Fix**: Capture `start` with `time.Now()` (which includes the monotonic clock) and use `time.Since` on it.

## Production notes

- In HTTP servers, set `ReadTimeout` and `WriteTimeout` using `time.Duration` to prevent slow-client attacks.
- Always pass `context.Context` with deadlines rather than storing `time.Time` in request-scoped state.
- Use `time.Time` as the canonical type; convert to Unix timestamps only at API boundaries (JSON, protobuf, DB).
- For high-precision logging, use `time.RFC3339Nano` to capture nanosecond resolution.
- The zero value of `time.Time` is year 1, month 1, day 1 — not the Unix epoch. Check `t.IsZero()` when handling optional timestamps.

## Performance implications

- `time.Now()` is fast (~30ns on modern hardware via vdso).
- `time.Format` allocates — avoid in hot paths. Cache formatted strings or format at the boundary.
- `time.Parse` is ~500ns and allocates multiple strings internally.
- Comparing `time.Time` with `Equal()` is slightly slower than `==` but necessary for correctness.
- `time.Since` is nearly free (monotonic subtraction of two uint64 values).

## Practice task

Write a function `parseTimestamp(s string) (time.Time, error)` that accepts timestamps in either RFC3339 (`"2026-06-01T14:30:00Z"`) or Unix integer (`"1717200000"`) format and returns the parsed `time.Time`. For Unix strings, parse the integer and use `time.Unix(sec, 0)`. Then write a main function that prints the parsed time formatted as `"2006-01-02 15:04:05"`.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/01-time-basics
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/01-time-basics
```

The existing tests verify `Parse`, `Format`, `Since`, and `Until`. After completing the practice task, add tests for `parseTimestamp` covering both formats and error handling for invalid input.

## Review questions

1. What are the two internal clock components of `time.Time`, and which one does `time.Since` use?
2. Why does `time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)` not carry a monotonic clock reading?
3. How do you correctly compare two `time.Time` values for equality?
4. Given `now := time.Now()`, write the expression that produces a formatted string like `"2026-Jun-01"`.
5. What happens if you pass the wrong layout string to `time.Parse`?

## NEXT UP

Durations — how to represent, compute, and compare spans of time.
