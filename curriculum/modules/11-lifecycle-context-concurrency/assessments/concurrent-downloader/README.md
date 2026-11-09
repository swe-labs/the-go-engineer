# Project Assessment — Concurrent Downloader

## Purpose

Evaluate your ability to build a resilient, high-throughput concurrent file downloader in Go. This assessment verifies that you can coordinate parallel network downloads with bounded worker pools, stream HTTP response bodies directly into `io.Writer` targets via `io.Copy`, recover from transient network failures using exponential backoff retries, track global throughput with `sync/atomic`, and propagate context cancellation.

## Targets

- `project-concurrent-downloader`

## Scope

This assessment covers:

- Bounded worker orchestration: coordinating parallel downloads using buffered channels and `sync.WaitGroup`
- Constant $O(1)$ memory streaming: piping `http.Response.Body` directly to `io.Writer` via `io.Copy`
- Transient fault recovery: exponential backoff retries (`initialBackoff * 2^attempt`) with retry limits
- Context propagation: honoring caller cancellation and deadlines (`context.WithTimeout`) mid-download
- Error handling: mapping response status codes and returning `ErrMaxRetriesExceeded` when retry limits are exhausted
- Lock-free telemetry: atomic throughput accounting of total downloaded bytes across all workers (`sync/atomic.Int64`)
- Data race elimination under concurrent execution verified by `go test -race`

## Evidence required

Submit:

1. Your completed `_starter/main.go` with streaming download, exponential backoff, and bounded worker orchestration implemented
2. Output of `go test -v -race .` showing all unit and retry tests passing with zero data races
3. A written explanation answering: How does your worker pool and backpressure mechanism handle a remote server throttling requests with 429 or dropping connections mid-stream?

## Time expectation

Plan for 90 to 120 minutes.

## Rubric

| Criterion | Weight | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------|--------------|----------|----------------|
| **Correctness** | 35% | Parallel downloads succeed across workers; retry count respected; context cancellation stops downloads; all automated tests pass | Most downloads succeed; minor retry or result indexing flaw | Non-functional downloader or goroutine deadlock |
| **Streaming & memory efficiency** | 25% | Data streamed directly with `io.Copy`; memory footprint remains constant regardless of file size; body closed properly | Streamed but body closure omitted or minor buffering overhead | Entire file loaded into unbounded memory slice (`io.ReadAll`) causing OOM |
| **Fault tolerance & backoff** | 25% | Exponential backoff cleanly recovers from transient faults; stops on context cancellation; returns `ErrMaxRetriesExceeded` | Retries present but backoff static or cancellation ignored | No retry logic on network blips |
| **Thread safety & explanation** | 15% | Total bytes tracked atomically; zero data races under `-race`; learner explains worker scaling and backpressure | Race detector clean but explanation has minor conceptual gaps | Data races detected or unable to explain backoff mechanics |

**Passing score**: 80%

## Review questions

1. How does your worker pool and backpressure mechanism handle a remote server throttling requests with 429 or dropping connections mid-stream?
2. How does your cancellation logic propagate through `context.Context` to abort in-flight chunk downloads immediately without leaking goroutines or partial files?
3. Why is `io.Copy` essential for large file downloads rather than `io.ReadAll`? What are the memory allocation and garbage collection implications?
4. How does exponential backoff prevent the "thundering herd" problem when an upstream API recovers from an outage?
5. Why is `sync/atomic.Int64` preferred over a `sync.Mutex` for tracking total downloaded bytes across all active download workers?

## Passing standard

You pass when all race detector tests pass, streaming uses constant memory, retries handle transient faults, and you can explain context propagation and backoff mechanics.

## Retake policy

Learner may revise code and resubmit after addressing all rubric gaps.

## NEXT UP

[URL Health Checker Assessment](../url-health-checker/README.md)
