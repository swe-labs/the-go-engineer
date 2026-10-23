# Project Assessment — Bounded Worker Pool

## Purpose

Evaluate your ability to build a robust, bounded worker pool in Go. This assessment verifies that you can manage goroutine lifecycles, coordinate channel-based job distribution, isolate task panics (`recover`), implement non-blocking and deadline-driven backpressure (`select`), drain queues gracefully with `sync.WaitGroup`, and pass the Go race detector with zero data races.

## Targets

- `project-bounded-worker-pool`

## Scope

This assessment covers:

- Bounded concurrency: strictly constraining active goroutines to configured worker count `W`
- Channel lifecycle and ownership: safe channel closure guarantees with `sync.Once` and `sync.RWMutex`
- Non-blocking job dispatch: immediate rejection with `ErrQueueFull` using `select` with `default`
- Context-aware submission: honoring caller cancellation and deadlines via `SubmitWithContext`
- Panic containment boundary: isolating panicking jobs with `recover()` to prevent process termination
- Lock-free telemetry: atomic counters for active, completed, failed, and panicked jobs, plus CAS-based peak tracking
- Graceful worker drainage: connection draining with `sync.WaitGroup` and bounded shutdown timeouts (`ErrShutdownTimeout`)
- Data race elimination under concurrent execution verified by `go test -race`

## Evidence required

Submit:

1. Your completed `_starter/main.go` with worker loops, panic containment, backpressure, and shutdown implemented
2. Output of `go test -v -race .` showing all unit and concurrency tests passing with zero data races
3. A written explanation answering: Explain channel ownership: which goroutine is responsible for closing the task channel, and what happens if a sender writes to a closed channel?

## Time expectation

Plan for 90 to 120 minutes.

## Rubric

| Criterion | Weight | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------|--------------|----------|----------------|
| **Correctness** | 35% | Concurrency bound `W` never exceeded; non-blocking submit returns `ErrQueueFull`; context submit honors deadlines; all automated tests pass | Most tests pass; minor queue saturation or context edge case | Concurrency bound exceeded or submit deadlocks |
| **Panic containment & metrics** | 25% | Worker catches job panics with `recover()`, updates panicked metric, and keeps goroutine alive; CAS tracks peak workers | Panics caught but metrics inaccurate or worker terminates on panic | Unhandled job panics crash entire pool process |
| **Graceful drainage & shutdown** | 25% | `Stop()` drains all queued jobs before worker termination; `StopWithTimeout()` returns `ErrShutdownTimeout` on hung jobs; channel closed safely once | Shutdown works but pending jobs dropped or timeout unhandled | Deadlock on shutdown or send on closed channel panic |
| **Thread safety & Explanation** | 15% | Zero data races under `-race`; learner articulates channel closing ownership, panic recovery mechanics, and context propagation | Race detector clean but explanation has minor conceptual gaps | Data races detected or unable to explain channel lifecycle |

**Passing score**: 80%

## Review questions

1. Explain channel ownership: which goroutine is responsible for closing the task channel, and what happens if a sender writes to a closed channel?
2. How does your worker pool handle a task that panics: does the entire pool crash, or is the panic recovered and reported while keeping worker routines alive?
3. Walk me through how `context.WithCancel` or `context.WithTimeout` triggers worker drainage and graceful exit when terminating the application.
4. Why is unbounded goroutine spawning (`go func()`) dangerous in high-throughput backend services, and how does backpressure protect system resources?
5. What is the purpose of `atomic.CompareAndSwapInt32` in tracking peak concurrent workers without a global lock?

## Passing standard

You pass when all race detector tests pass, concurrency bounds are strictly maintained, panics are contained, and you can explain channel mechanics and shutdown lifecycles.

## Retake policy

Learner may revise code and resubmit after addressing all rubric gaps.

## NEXT UP

[Concurrent Downloader Assessment](../concurrent-downloader/README.md)
