# Track A: Goroutines

## Mission

This track teaches you how Go schedules concurrent work and how to coordinate that work safely
with WaitGroups and channels before moving into more advanced synchronization topics.

## Track Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `GC.1` | Lesson | [goroutines](./1-goroutine) | Introduces goroutines, scheduler basics, and closure-capture pitfalls. | entry |
| `GC.2` | Lesson | [WaitGroups](./2-wait-group) | Teaches the standard barrier pattern for waiting on concurrent work. | `GC.1` |
| `GC.3` | Lesson | [channels](./3-channels) | Introduces typed communication between goroutines. | `GC.1`, `GC.2` |
| `GC.4` | Lesson | [buffered channels](./4-channels-buffered) | Explains queue-like channel behavior and bounded async flow. | `GC.3` |
| `GC.5` | Lesson | [closing channels](./5-channels-closing) | Shows how receivers learn a producer is done. | `GC.3`, `GC.4` |
| `GC.6` | Lesson | [pipeline project](./6-project-1) | Combines goroutines, channels, and cancellation into one mini flow. | `GC.3`, `GC.4`, `GC.5` |
| `GC.7` | Exercise | [concurrent downloader](./7-downloader) | Applies coordination, fan-out, and bounded concurrency in one milestone. | `GC.1`, `GC.2`, `GC.3`, `GC.4`, `GC.5`, `GC.6` |

## Suggested Order

1. Work through `GC.1` to `GC.6` in order.
2. Complete `GC.7` as the live goroutines milestone.
3. Use the later goroutines lessons as legacy reference after the milestone if you want deeper coverage.

## Track Milestone

`GC.7` is the current goroutines track milestone.

If you can complete it and explain:

- why goroutines still need explicit coordination instead of "fire and forget"
- why bounded concurrency is safer than launching unlimited work
- why channel-based result aggregation is often cleaner than shared mutable state

then the core goroutines part of Stage 07 is doing its job.

## Legacy Reference Surfaces

These lessons remain available, but they are not part of the live v2 track map yet:

- `GC.8` race conditions
- `GC.9` select deep dive
- `GC.10` sync primitives

## Next Step

After `GC.7`, continue to the [Stage 07 overview](../README.md) or move into the
[Context track](../context).
