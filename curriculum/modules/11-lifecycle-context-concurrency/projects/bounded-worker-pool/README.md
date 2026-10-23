# Bounded Worker Pool

## Learning objective

Build a rock-solid, production-grade concurrency worker pool in Go. Master bounded goroutine lifecycle management, channel synchronization, panic containment without process crashes, backpressure handling, graceful task draining, and zero-data-race thread safety.

## Why this matters

Unbounded goroutine spawning (`go func()`) under load is a leading cause of memory exhaustion and CPU thrashing in Go backends. Production services (e.g. webhook dispatchers, payment processors, image converters) require strictly bounded worker pools with backpressure so callers know when queues are full and workers never exceed hardware thresholds.

## Prerequisites

Complete all Module 11 lessons (01–28) on goroutines, channels, synchronization primitives (`sync.Mutex`, `sync.WaitGroup`, `sync.Once`), atomic operations, `context.Context` cancellation, and race detector diagnostics.

## Available concepts

This project draws on goroutines, buffered channels, `select` statements with non-blocking defaults, `sync/atomic` counters, `sync.WaitGroup`, `panic` / `recover`, `context.WithCancel`, `context.WithTimeout`, and `go test -race`.

## Project overview

The `_starter/` directory provides type signatures, channel structures, and metric trackers. Your job is to implement the worker loop, dispatch queue, panic recovery boundary, and graceful shutdown mechanics so that:

- Concurrency is strictly bounded: At no point do active workers exceed the configured limit `W`.
- Submissions adhere to backpressure: When the job queue is saturated, `Submit()` immediately returns `ErrQueueFull`, and `SubmitWithContext()` honors cancellation deadlines.
- Panicking jobs are contained: A worker catching a panic logs the error, increments the panic counter, and continues serving subsequent jobs without terminating.
- Graceful shutdown drains all buffered jobs before terminating worker goroutines.
- The entire system passes the Go race detector with zero data races.

## Architecture

```
[ Caller ] --- Submit() ---> [ Buffered Channel (QueueSize) ]
                                   |
           +-----------------------+-----------------------+
           |                       |                       |
     [ Worker 1 ]            [ Worker 2 ]            [ Worker N ]
     (recover loop)          (recover loop)          (recover loop)
```

## Tasks

### Task 1 — Implement worker lifecycle and execution loop

In `_starter/main.go`:
- Implement `Start(ctx context.Context)`: spawn `p.workers` goroutines running `p.workerLoop`. Increment `p.wg` for each worker.
- Implement `workerLoop(ctx context.Context, workerID int)`:
  - Defer `p.wg.Done()`.
  - Select between `<-ctx.Done()` and `<-p.jobs`.
  - When a job is received, execute it via `p.executeJob(ctx, job)`.

### Task 2 — Implement panic containment & metric tracking

- Implement `executeJob(ctx context.Context, job Job)`:
  - Increment `activeWorkers` atomically; decrement on function exit.
  - Track peak concurrency using `atomic.CompareAndSwapInt32(&p.peakWorkers, ...)`.
  - Defer a `recover()` block that increments `p.panicked` counter and logs the panic with `slog.Error`.
  - Execute `job(ctx)`. Increment `p.completed` on success or `p.failed` on error.

### Task 3 — Implement backpressure dispatching

- Implement `Submit(job Job) error`:
  - Verify job is non-nil and pool is not closed under `p.mu.RLock()`.
  - Use a non-blocking `select` with `default`: push to `p.jobs` or return `ErrQueueFull`.
- Implement `SubmitWithContext(ctx context.Context, job Job) error`:
  - Select between `<-ctx.Done()` (return context error) and pushing to `p.jobs`.

### Task 4 — Implement graceful shutdown

- Implement `Stop()`:
  - Set `p.closed = true` and close `p.jobs` channel under lock (ensured only once with `sync.Once`).
  - Wait for all active and queued jobs to complete via `p.wg.Wait()`.
- Implement `StopWithTimeout(timeout time.Duration) error`:
  - Stop pool and return `ErrShutdownTimeout` if workers fail to exit within `timeout`.

## Verification

Run all unit tests with the Go race detector enabled:

```bash
go test -v -race ./curriculum/modules/11-lifecycle-context-concurrency/projects/bounded-worker-pool/_starter
```

Expected output:
```
=== RUN   TestWorkerPool_ConcurrencyBound
--- PASS: TestWorkerPool_ConcurrencyBound (0.07s)
=== RUN   TestWorkerPool_PanicContainment
--- PASS: TestWorkerPool_PanicContainment (0.00s)
=== RUN   TestWorkerPool_BackpressureAndQueueFull
--- PASS: TestWorkerPool_BackpressureAndQueueFull (0.02s)
=== RUN   TestWorkerPool_GracefulShutdown
--- PASS: TestWorkerPool_GracefulShutdown (0.05s)
=== RUN   TestWorkerPool_RaceConditionFreedom
--- PASS: TestWorkerPool_RaceConditionFreedom (0.00s)
PASS
```

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|---|---|---|---|
| **Concurrency Bound & Safety** | Active workers strictly bounded to limit; zero data races detected by `go test -race` | Concurrency mostly bounded; minor race conditions under heavy load | Unbounded goroutine creation or data races |
| **Panic Containment** | Panicking jobs recovered without crashing worker goroutine; subsequent jobs execute | Panic caught but worker dies, permanently reducing pool capacity | Unhandled panic crashes the entire process |
| **Backpressure & Deadlines** | Saturated queue returns `ErrQueueFull` or respects context timeout without blocking forever | Queue blocks indefinitely or drops jobs silently | Unchecked submission leads to out-of-memory errors |
| **Graceful Shutdown** | All queued jobs complete cleanly; closed pool rejects new submissions deterministically | Pool shuts down but in-flight jobs are abruptly killed | Deadlock or hanging on shutdown |

## NEXT UP

[Flagship Opslane Capstone](../../18-flagship-opslane/README.md)
