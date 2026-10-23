# Module 11 Assessment — Lifecycle, Context, and Concurrency

## Purpose

Verify that the learner can independently apply and explain goroutines, channels, select, context cancellation, sync primitives, deadlock prevention, pipeline patterns, and thread-safe systems engineering.

## Targets

- `core-11-01`
- `core-11-02`
- `core-11-03`
- `core-11-04`
- `core-11-05`
- `core-11-06`
- `core-11-07`
- `core-11-08`
- `core-11-09`
- `core-11-10`
- `core-11-11`
- `core-11-12`
- `core-11-13`
- `core-11-14`
- `core-11-15`
- `core-11-16`
- `core-11-17`
- `core-11-18`
- `core-11-19`
- `core-11-20`
- `core-11-21`
- `core-11-22`
- `core-11-23`
- `core-11-24`
- `core-11-25`
- `core-11-26`
- `core-11-27`
- `core-11-28`
- `core-11-29`

## Scope

This checkpoint covers:

- Goroutine scheduling mechanics: M:N runtime scheduler model (Goroutines, OS Threads, Logical Processors), work-stealing, and preemption
- Channel semantics: unbuffered handshakes, buffered queues, directional types (`<-chan`, `chan<-`), and nil channel behavior
- Multiplexing with `select`: non-blocking channel polling with `default`, timeout channels (`time.After`), and fair random selection
- Synchronization primitives: `sync.Mutex`, `sync.RWMutex`, `sync.WaitGroup`, `sync.Once`, and `golang.org/x/sync/errgroup`
- Lock-free synchronization: atomic counters, atomic booleans, and CAS (`CompareAndSwap`) via `sync/atomic`
- Request lifecycle management: deadline propagation, graceful cancellation chains, and values with `context.Context`
- High-throughput concurrency patterns: bounded worker pools, fan-out/fan-in topologies, streaming pipelines, and leaky bucket / token bucket rate limiting
- Concurrency defect diagnosis: data races, race conditions, deadlocks, livelocks, starvation, and goroutine leaks
- Tooling and verification: race detector (`go test -race`), runtime stack traces, and pprof goroutine profiling

## Evidence required

Submit:

1. **Concurrent Go program**: Bounded concurrency architecture coordinating work distribution via channels and synchronizing worker lifecycles
2. **Context-aware cancellation**: Implementation demonstrating rapid deadline abortion and resource cleanup across sub-goroutines
3. **Race detector verification**: Test suite execution showing zero data races under `go test -v -race .`
4. **Written explanations**: Answers to the five review questions below

## Time expectation

Plan for 2 to 3 hours.

## Rubric

| Criterion | Weight | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------|--------------|----------|----------------|
| **Conceptual clarity** | 30% | Clearly explains Go's M:N scheduler, goroutine leak mechanisms, channel ownership rules, and difference between data race and race condition | Explains basics of goroutines, channels, and context with minor gaps | Cannot explain M:N scheduling, channel states, or goroutine leaks |
| **Implementation** | 35% | Complete, clean concurrent implementation with bounded workers, non-blocking channels, context propagation, and graceful shutdown | Working implementation with minor race or timeout flaw | Unbounded goroutine creation or deadlock-prone channel logic |
| **Proof** | 25% | Strong test suite passing with `go test -race` with zero race warnings, verifying cancellation and queue bounds | Basic tests without race detector or missing edge cases | Tests fail or trigger data races under `-race` |
| **Communication** | 10% | Professional explanation connecting concurrency patterns to backend throughput, resource bounds, and system resilience | Basic explanation with some clarity | Unable to articulate concurrency trade-offs or race conditions |

**Passing score**: 80%

## Review questions

1. Explain Go's goroutine scheduler. How does M:N threading work, and what causes goroutine starvation?
2. Walk me through a goroutine leak. Write the code that causes it, explain why it leaks, and show the fix.
3. Explain the difference between a data race and a race condition. Can you have a race condition without a data race?
4. What are the rules of channel ownership? Who should close a channel, and why is sending to a closed channel a panic?
5. What is the difference between `context.WithCancel` and `context.WithTimeout`, and why should `cancel()` always be deferred immediately?

## Passing standard

Score of 80% or higher across the four rubric criteria, with all required evidence submitted.

## Retake policy

Learner may retry after fixing rubric gaps and writing a short correction note.

## NEXT UP

[Module 12: Observability and Diagnostics](../../12-observability-diagnostics/README.md)
