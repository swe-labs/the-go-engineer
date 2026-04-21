# Track C: Time and Scheduling

## Mission

This track teaches you how Go models time values and timed events so you can build reminders,
timeouts, and interval-based background work without leaking timers or tickers.

## Track Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `TM.1` | Lesson | [time basics](./1-time) | Introduces `time.Time`, `time.Duration`, and elapsed-time calculations. | entry |
| `TM.2` | Lesson | [formatting](./2-formatting) | Teaches parsing and formatting with Go's reference-time model. | `TM.1` |
| `TM.3` | Lesson | [timers and tickers](./3-timer-and-ticker) | Explains one-shot and repeating timed events plus cleanup rules. | `TM.1`, `TM.2` |
| `TM.7` | Exercise | [console reminder](./7-reminder) | Combines timers, tickers, and command-line input in one runnable milestone. | `TM.1`, `TM.2`, `TM.3` |

## Suggested Order

1. Work through `TM.1` to `TM.3` in order.
2. Complete `TM.7` as the live time-track milestone.
3. Use the legacy reference lessons later if you want deeper scheduling coverage.

## Track Milestone

`TM.7` is the current time-and-scheduling track milestone.

If you can complete it and explain:

- why `time.NewTicker` needs `Stop()` cleanup
- why `time.AfterFunc` is useful for one-shot delayed work
- why timer-driven programs still need explicit exit signaling

then the time part of Stage 07 is doing its job.

## Legacy Reference Surfaces

These lessons remain available, but they are not part of the live v2 track map yet:

- `TM.4` random numbers
- `TM.5` scheduler
- `TM.6` timezones

## Next Step

After `TM.7`, continue to the [Stage 07 overview](../README.md) or move on to
[Stage 07: Concurrency Patterns](../../02-concurrency-patterns).

