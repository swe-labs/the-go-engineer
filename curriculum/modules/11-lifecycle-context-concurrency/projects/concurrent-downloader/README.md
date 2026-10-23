# Concurrent Downloader

## Learning objective

Build a resilient, high-throughput concurrent file downloader in Go. Master bounded worker pool coordination, HTTP streaming with `io.Copy`, transient fault recovery with exponential backoff retries, aggregate atomic byte accounting, and cancellation propagation.

## Why this matters

Network I/O is inherently prone to transient failures: flaky connections, rate limits (HTTP 429), and temporary upstream outages (HTTP 503). Writing production-grade download clients requires bounded worker concurrency to prevent socket exhaustion, automatic retry logic with backoff, and robust streaming directly into destinations without unbounded memory allocation.

## Prerequisites

Complete all Module 11 lessons (01–28) on goroutines, channels, synchronization with `sync.WaitGroup`, atomic operations with `sync/atomic`, and context timeouts with `context.WithTimeout`.

## Available concepts

This project draws on worker pools, `io.Writer`, `io.Copy`, `http.Client` contexts, exponential backoff with `time.After`, `sync/atomic.Int64`, and `go test -race`.

## Project overview

The `_starter/` directory provides type definitions, job models, and stub methods. Your job is to implement the concurrent download coordinator and retry mechanism so that:

- Multiple download jobs are processed concurrently without exceeding the configured worker bound.
- Streams are piped directly to destinations (`io.Writer`) via `io.Copy` to maintain constant O(1) memory usage.
- Transient errors (network blips, non-2xx status) trigger retries with exponential backoff up to `maxRetries`.
- Cancellation via `context.Context` aborts in-flight downloads promptly.
- Total downloaded bytes are tracked atomically across all workers.
- The test suite executes with `go test -race` with zero data races.

## Architecture

```
[ Download Jobs ] ---> [ Buffered Queue ] ---> [ Bounded Workers (1..W) ]
                                                       |
                                           (downloadWithRetry)
                                           - Attempt 1: Fail (503)
                                           - Wait initialBackoff
                                           - Attempt 2: Success (200)
                                                       |
                                            (io.Copy to Destination)
                                                       |
                                           (atomic.AddInt64 TotalBytes)
```

## Tasks

### Task 1 — Implement streaming download

In `_starter/main.go`:
- Implement `downloadSingle(ctx context.Context, job DownloadJob) (int64, error)`:
  - Create `http.NewRequestWithContext` using `http.MethodGet`.
  - Execute with `d.client.Do(req)`.
  - Validate response status is 2xx; return `ErrDownloadFailed` on 4xx/5xx.
  - Stream `resp.Body` into `job.Destination` using `io.Copy`.
  - Accumulate written bytes into `d.totalBytes` atomically.

### Task 2 — Implement exponential backoff retries

- Implement `downloadWithRetry(ctx context.Context, job DownloadJob) (int64, int, error)`:
  - Attempt `downloadSingle` in a loop up to `d.maxRetries`.
  - If successful, return written bytes and retry count.
  - Check for context cancellation before sleeping.
  - Wait `backoff` duration (`backoff *= 2`), then retry.
  - Return `ErrMaxRetriesExceeded` if all retries fail.

### Task 3 — Implement bounded concurrent coordinator

- Implement `DownloadAll(ctx context.Context, jobs []DownloadJob) []DownloadResult`:
  - Launch `min(d.workers, len(jobs))` worker goroutines.
  - Synchronize workers using `sync.WaitGroup`.
  - Record download results into output slice preserving original job index.

## Verification

Run all unit tests with the Go race detector enabled:

```bash
go test -v -race ./curriculum/modules/11-lifecycle-context-concurrency/projects/concurrent-downloader/_starter
```

Expected output:
```
=== RUN   TestConcurrentDownloader_Success
--- PASS: TestConcurrentDownloader_Success (0.00s)
=== RUN   TestConcurrentDownloader_RetryOnTransientFailure
--- PASS: TestConcurrentDownloader_RetryOnTransientFailure (0.01s)
=== RUN   TestConcurrentDownloader_MaxRetriesExceeded
--- PASS: TestConcurrentDownloader_MaxRetriesExceeded (0.02s)
=== RUN   TestConcurrentDownloader_ContextCancellation
--- PASS: TestConcurrentDownloader_ContextCancellation (0.20s)
PASS
```

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|---|---|---|---|
| **Streaming & Memory Efficiency** | Data streamed directly with `io.Copy`; memory footprint remains constant regardless of file size | Buffer used with reasonable memory | Entire file loaded into unbounded memory slice (`io.ReadAll`) |
| **Fault Tolerance & Backoff** | Exponential backoff cleanly recovers from transient faults; stops on context cancellation | Retries without backoff or fails to track attempts | No retry logic on network blips |
| **Concurrency Bound & Safety** | Worker count strictly bounded; zero data races detected by race detector | Concurrency bounded but slight lock contention | Spawns unbounded goroutines |
| **Coordination & Metrics** | Atomic byte counting; clean result aggregation preserving input ordering | Byte counting correct but without atomic primitives | Race conditions in metric accumulation |

## NEXT UP

[Module 12: Observability & Profiling](../../12-observability-profiling/README.md)
