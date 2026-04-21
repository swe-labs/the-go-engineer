# Stage 07: Concurrency

## Mission

This track teaches how Go coordinates work across goroutines, channels, contexts, sync primitives, and timers without hiding the cost of background work.

By the end of this track, you should be comfortable:

- starting goroutines deliberately instead of sprinkling `go` blindly
- coordinating work with `sync.WaitGroup`, channels, and synchronization primitives
- applying cancellation and deadlines with `context.Context`
- using timers and tickers without leaking background work
- building small concurrent tools that stay readable under load

## Stage Ownership

This track belongs to [07 Concurrency](../README.md).

## Track Map

| Track | Entry | Milestone | Focus |
| --- | --- | --- | --- |
| Goroutines | [GC.1 goroutines](./goroutines) | `GC.7` | goroutines, WaitGroups, channels, and a bounded downloader |
| Context | [CT.1 background](./context) | `CT.5` | context roots, cancellation, timeouts, and timeout-aware HTTP calls |
| Sync Primitives | [SY.1 mutex and rwmutex](./sync-primitives) | `SY.6` | mutexes, `sync.Map`, atomics, race boundaries, and deadlocks |
| Time & Scheduling | [TM.1 time basics](./time-and-scheduling) | `TM.7` | time values, formatting, timers, tickers, and reminder scheduling |

## Suggested Order

1. Start with the Goroutines track.
2. Add Context once goroutine lifecycles feel concrete.
3. Learn Sync Primitives once shared-state bugs make sense.
4. Use Time & Scheduling to round out deadline and timer intuition.

## Track Milestones

The promoted outputs in this track are:

- `GC.7` concurrent downloader
- `CT.5` timeout-aware API client
- `SY.6` deadlocks
- `TM.7` console reminder

## Next Step

After the core concurrency tracks, continue to [Stage 07: Concurrency Patterns](../02-concurrency-patterns), then move on to [08 Quality & Testing](../../08-quality-test).
