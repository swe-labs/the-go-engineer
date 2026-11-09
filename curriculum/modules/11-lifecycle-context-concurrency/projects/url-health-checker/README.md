# URL Health Checker

## Learning objective

Build a concurrent HTTP health check probe engine in Go. Master bounded fan-out / fan-in worker patterns, per-request deadline enforcement with `context.WithTimeout`, latency instrumentation, thread-safe status change detection, and zero-race concurrency.

## Why this matters

Distributed microservice architectures require continuous synthetic monitoring to detect degraded dependencies, network partitions, and flap cycles before users are impacted. Writing resilient monitoring probes requires bounded worker pools, deadline propagation to HTTP clients, and atomic state tracking.

## Prerequisites

Complete all Module 11 lessons (01–28) on goroutines, channel fan-out/fan-in patterns, `context.Context` cancellation/timeouts, and thread-safe data structures with `sync.RWMutex`.

## Available concepts

This project draws on goroutines, worker queues, `sync.WaitGroup`, `sync.RWMutex`, `http.Client` request contexts, `time.Duration` latency calculation, state machine transitions, and `go test -race`.

## Project overview

The `_starter/` directory provides models, configuration structs, and stub methods. Your job is to implement the concurrent probing engine and transition tracker so that:

- A batch of target URLs is probed concurrently without exceeding the configured worker count.
- Every HTTP request strictly honors a timeout context; hanging endpoints do not block workers indefinitely.
- Probe results (`UP` or `DOWN`, HTTP status code, latency) preserve original URL slice ordering.
- State transitions (e.g. `UP` → `DOWN` or `DOWN` → `UP`) are detected and recorded thread-safely.
- The entire engine passes `go test -race` with zero data races.

## Architecture

```
[ URLs Slice ] ---> [ Work Channel ] ---> [ Workers (1..W) ] ---> [ Output Results ]
                                                   |
                                     (ProbeTarget with context timeout)
                                                   |
                                     (Thread-safe Transition Tracker)
```

## Tasks

### Task 1 — Implement single probe with timeout

In `_starter/main.go`:
- Implement `ProbeTarget(ctx context.Context, targetURL string) CheckResult`:
  - Create a child context with `h.timeout`.
  - Issue `http.NewRequestWithContext` using `http.MethodGet`.
  - Execute request with `h.client.Do(req)`.
  - Measure elapsed latency and classify status (`UP` for 2xx/3xx, `DOWN` for errors or 5xx).
  - Invoke `recordTransition(url, status)`.

### Task 2 — Implement transition detector

- Implement `recordTransition(url string, currentStatus EndpointStatus)`:
  - Protect `h.lastStatus` map access with `h.mu.Lock()`.
  - If a previous status exists and differs from `currentStatus`, record a `TransitionEvent`.
  - Update `h.lastStatus[url]`.

### Task 3 — Implement bounded concurrent batch execution

- Implement `CheckAll(ctx context.Context, urls []string) []CheckResult`:
  - Distribute tasks across a bounded number of worker goroutines (`min(h.workers, len(urls))`).
  - Coordinate synchronization with `sync.WaitGroup`.
  - Store probe results into the output slice preserving the original URL input index.

## Verification

Run all unit tests with the Go race detector enabled:

```bash
go test -v -race ./curriculum/modules/11-lifecycle-context-concurrency/projects/url-health-checker/_starter
```

Expected output:
```
=== RUN   TestProbeTarget_SuccessAndFailure
--- PASS: TestProbeTarget_SuccessAndFailure (0.00s)
=== RUN   TestProbeTarget_TimeoutHandling
--- PASS: TestProbeTarget_TimeoutHandling (0.20s)
=== RUN   TestCheckAll_ConcurrentExecution
--- PASS: TestCheckAll_ConcurrentExecution (0.04s)
=== RUN   TestStatusTransitions
--- PASS: TestStatusTransitions (0.00s)
PASS
```

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|---|---|---|---|
| **Concurrency & Worker Bound** | Bounded concurrency strictly maintained; slice indices preserved without race conditions | Concurrency bounded but results unordered or lock contention | Unbounded goroutine spawning per target |
| **Timeout & Context Control** | Every probe governed by timeout context; slow/hung endpoints aborted promptly | Timeout set globally on client without per-request context | Hangs indefinitely on slow endpoints |
| **State & Transition Tracking** | Thread-safe transition history recording; detects flapping state changes | State recorded but race conditions possible on concurrent checks | Misses status changes or crashes on concurrent access |
| **Code Quality & Idioms** | Idiomatic standard library usage; clean separation of probe and coordination logic | Functional but tightly coupled probe logic | Non-idiomatic channel or mutex handling |

## NEXT UP

[Concurrent Downloader](../concurrent-downloader/README.md)
