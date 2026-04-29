# 07 Concurrency

## Mission

This stage opens with why concurrency exists, keeps the existing goroutine, context, and time tracks, and adds sync primitives as a first-class track before moving to advanced patterns.

By the end of this stage, a learner should be able to:

- launch and coordinate concurrent work using goroutines and WaitGroups
- protect shared state using Mutexes and RWMutexes
- communicate safely between goroutines using channels
- manage timeouts, cancellation, and deadlines with `context`
- use timers, tickers, and advanced concurrency patterns (errgroup, pools)

## Stage Map

| Track | Surface | Core Job |
| --- | --- | --- |
| `GC.0-7` | Goroutines | Ground the reason concurrency exists and how goroutines communicate. |
| `SY.1-6` | Synchronization | Learn the synchronization tools (mutex, atomic) that repair shared-state bugs. |
| `CT.1-5` | Context | Learn to manage lifecycle, timeouts, and cancellation graphs. |
| `TM.1-7` | Time | Teach timers, tickers, scheduling, and cleanup. |
| `CP.1-5` | Concurrency Patterns | Master advanced patterns like `errgroup`, `sync.Pool`, and bounded pipelines. |

## Why This Stage Exists Now

The learner already knows:

- how to build backend APIs
- how to query databases
- how to handle I/O boundaries

That is enough to start asking engineering questions like:

- how do we process 10,000 items without waiting for them one by one?
- how do we prevent two requests from corrupting the same in-memory counter?
- how do we stop work early if the user cancels the request?

## Suggested Learning Flow

1. Start with `GC.0` to `GC.7` to ground the reason concurrency exists.
2. Use `SY.1` to `SY.6` to learn the tools that repair shared-state bugs.
3. Move through Context (`CT`), Time (`TM`), and Patterns (`CP`) to master runtime coordination.

## Next Step

After this section, continue to [08 Quality & Testing](../08-quality-test).
