# Opslane worker pools

## Mission

Design and implement a robust, bounded worker pool in Go, managing fixed concurrency ceilings, buffered backpressure channels, atomic performance telemetry, panic isolation, and graceful drain-and-stop shutdown.

## Prerequisites

- `opslane-09`

## Mental Model

Unbounded goroutine creation (`go func()`) under high incoming traffic is a common cause of production outages. Spawning 100,000 goroutines to process bulk orders exhausts memory, triggers severe garbage collection pressure, and degrades CPU cache efficiency. A bounded worker pool establishes a rigid concurrency ceiling: a fixed number of long-lived worker goroutines (e.g. 4 or 16) pull jobs from a bounded buffered channel (`chan Job`). If the buffer fills, callers experience backpressure rather than silently exhausting server memory. Telemetry metrics track active workers and peak concurrency using atomic primitives (`sync/atomic`), and `Shutdown()` ensures all in-flight jobs are completed before exiting.

## Visual Model

```text
Job Submissions (Submit(job))
  │
  ▼
┌─────────────────────────────────────────────────────────────┐
│ Bounded Queue (chan Job, capacity = 50)                     │
│   [Job 1] [Job 2] [Job 3] ... (Backpressure on full)        │
└──────────────┬───────────────┬───────────────┬──────────────┘
               │               │               │
               ▼               ▼               ▼
        ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
        │  Worker 1   │ │  Worker 2   │ │  Worker 3   │ (Fixed Pool)
        │ (goroutine) │ │ (goroutine) │ │ (goroutine) │
        └──────┬──────┘ └──────┬──────┘ └──────┬──────┘
               │               │               │
               └───────────────┼───────────────┘
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ Telemetry & Panic Isolation                                 │
│   - atomic.AddInt32(&activeWorkers, 1)                      │
│   - CAS maxConcurrencySeen                                  │
│   - defer recover()                                         │
│   - atomic.AddInt64(&completedJobs, 1)                      │
└─────────────────────────────────────────────────────────────┘
```

## Machine View

Under the hood:
1. `p.jobs = make(chan Job, queueCap)` allocates a bounded ring buffer in memory. `Submit` attempts a non-blocking channel send unless the channel is full.
2. `executeJob` tracks peak concurrency using a lock-free Compare-And-Swap (CAS) loop:
   ```go
   currentMax := p.maxConcurrencySeen.Load()
   if active <= currentMax { break }
   if p.maxConcurrencySeen.CompareAndSwap(currentMax, active) { break }
   ```
3. Panic recovery wraps `j.Execute(p.ctx)`. If a job panics, `p.activeWorkers.Add(-1)` and `p.completedJobs.Add(1)` are guaranteed to execute via `defer`, keeping telemetry accurate and keeping the worker goroutine alive.
4. In `Shutdown(ctx)`: `close(p.jobs)` signals to workers that no more jobs will be enqueued. Workers drain remaining buffered jobs from the channel and exit when the channel closes. `p.wg.Wait()` blocks until all workers call `p.wg.Done()`.

## Run Instructions

```bash
# Run worker pool test suite
go test -v ./curriculum/modules/18-flagship-opslane/assessments/10-opslane-worker-pools

# Run standalone pool execution
go run ./curriculum/modules/18-flagship-opslane/assessments/10-opslane-worker-pools
```

## Code Walkthrough

1. **Job Interface & SimpleJob (lines 16–28)**: Declares `Job` (`ID()`, `Execute(ctx)`) and standard closure-based `SimpleJob`.
2. **`BoundedWorkerPool` Struct (lines 31–42)**: Holds worker count, buffered channel, waitgroup, cancellation context, and atomic telemetry counters.
3. **`NewBoundedWorkerPool` (lines 44–60)**: Spawns $N$ `workerLoop` goroutines, registering each with `pool.wg.Add(1)`.
4. **`workerLoop` & `executeJob` (lines 62–113)**: Selects between context cancellation and job receipt; executes jobs with atomic tracking and panic recovery.
5. **`Submit()` (lines 116–127)**: Verifies closed state and pushes job into `jobs` channel.
6. **`Shutdown()` (lines 130–150)**: Atomic swap `closed.Swap(true)`, closes `jobs` channel, waits for workers to drain buffered work, and returns cleanly or with timeout error.

## Try It

1. Run the test suite: verify that 20 jobs submitted to a 4-worker pool complete with `maxSeen <= 4`.
2. Test panic isolation: submit a job that panics and verify that subsequent jobs complete normally.
3. Test shutdown: submit jobs, immediately initiate `Shutdown`, and verify all queued jobs drain to completion.

Common mistakes:
- Closing a channel from multiple goroutines or while producers are still sending, causing fatal panics.
- Unbounded concurrency: spawning an unconstrained goroutine per request instead of pooling.
- Forgetting to drain buffered channel jobs during shutdown, discarding in-flight customer work.

## In Production

In high-throughput services, worker pool capacities should be configured based on available CPU cores (`runtime.NumCPU()`) for CPU-bound tasks, or based on downstream connection limits (database connection pool size) for I/O-bound tasks.

## Thinking Questions

1. Why is a bounded worker pool preferred over launching an unconstrained `go func()` per incoming request?
2. How does `close(p.jobs)` allow worker goroutines to drain buffered jobs and terminate cleanly?
3. How does Compare-And-Swap (CAS) update `maxConcurrencySeen` safely without mutex locks?

## Next Step

Proceed to [Milestone 11 — Caching Layer](../11-opslane-caching-layer/README.md) to implement in-memory LRU / TTL caching with singleflight cache stampede protection.
