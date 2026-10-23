# Performance Analysis Lab

## Learning objective

Diagnose and resolve severe synchronization contention and heap memory allocation escapes in a high-throughput Go rate-limiting engine. Learn how to profile mutex contention with `pprof -mutex`, perform compiler escape analysis with `go build -gcflags="-m"`, and refactor naive synchronization into lock-striped shards (`[64]sync.Mutex`) with circular ring buffers and stack-allocated value returns.

## Why this matters

Under modest load, naive Go concurrency patterns—such as a single global mutex protecting a shared map—often pass functional tests without raising alarms. However, under high parallel concurrency (thousands of requests/sec across multiple CPU cores), goroutines pile up waiting on the single mutex lock convoy, causing throughput collapse and severe latency degradation. Knowing how to identify lock contention and escape analysis inefficiencies turns an engineer into an indispensable systems architect.

## Prerequisites

Complete Module 13 lessons on Go memory model, goroutine synchronization, mutex contention profiling, escape analysis, and benchmarking with `b.RunParallel()`.

## Available concepts

This project draws on `sync.Mutex`, sharded lock striping, FNV-1a hashing, circular ring buffers (`[Capacity]int64`), escape analysis (`go build -gcflags="-m"`), stack vs heap value semantics, and `testing.B.RunParallel`.

## Project overview

The `_starter/` directory implements a rate-limiting service (`ShardedRateLimiter`). While functionally valid for single-threaded usage, under high concurrent load it suffers from three critical bottlenecks:
1. **Global Lock Contention**: A single `sync.Mutex` serializes all requests across all client keys.
2. **Heap Memory Escape**: `CheckLimit` constructs and returns a heap-allocated pointer `*TokenResult`, producing millions of ephemeral objects for the garbage collector.
3. **Memory Leaks & Slice Reallocation**: Client buckets store unbounded slices (`[]time.Time`) that continuously reallocate and retain memory.

Your mission is to analyze these bottlenecks and implement a lock-striped, zero-allocation rate limiter.

## Profiling Workflow

### 1. Identify Lock Contention

Run the benchmark under parallel load:

```bash
go test -bench="BenchmarkRateLimiter_Parallel" -benchmem ./curriculum/modules/13-performance-memory-engineering/projects/performance-analysis-lab/_starter
```

Notice the high ns/op latency caused by goroutines waiting for the single mutex.

### 2. Inspect Heap Escapes

Analyze why objects escape to the heap:

```bash
go build -gcflags="-m" ./curriculum/modules/13-performance-memory-engineering/projects/performance-analysis-lab/_starter
```

Identify the line where `&TokenResult{...} escapes to heap`.

## Tasks

### Task 1 — Lock Striping with 64 Shards

In `_starter/main.go`:
- Refactor `ShardedRateLimiter` to maintain an array of 64 independent shards:
  ```go
  type clientShard struct {
      mu      sync.Mutex
      buckets map[string]*CircularRingWindow
  }
  ```
- Partition client keys across the 64 shards using an FNV-1a hash function:
  ```go
  func (l *ShardedRateLimiter) getShard(key string) *clientShard
  ```
- Acquire and release only the specific shard's lock, allowing 64 concurrent goroutines to operate simultaneously without blocking each other.

### Task 2 — Fixed Circular Ring Buffer

- Replace unbounded `[]time.Time` slices with a `CircularRingWindow` using a fixed-capacity slice of int64 Unix nano timestamps.
- Implement circular index wrap-around (`(head + 1) % capacity`) to prune expired timestamps and record new events in $O(1)$ time with zero slice allocations.

### Task 3 — Zero-Allocation Stack Value Returns

- Update `CheckLimit(key string) TokenResult` to return the struct **by value** rather than by pointer (`*TokenResult`).
- Verify via `go build -gcflags="-m"` that `CheckLimit` results remain on the goroutine stack, eliminating heap allocations during rate limit evaluations.

## Verification

Run the automated tests and parallel benchmark:

```bash
go test -v -bench="BenchmarkRateLimiter_Parallel" ./curriculum/modules/13-performance-memory-engineering/projects/performance-analysis-lab/_solution
```

Expected output:
```
=== RUN   TestRateLimiter_BurstAllowanceAndThrottle
--- PASS: TestRateLimiter_BurstAllowanceAndThrottle (0.25s)
=== RUN   TestRateLimiter_DistinctClients
--- PASS: TestRateLimiter_DistinctClients (0.00s)
=== RUN   TestRateLimiter_ConcurrentContention
--- PASS: TestRateLimiter_ConcurrentContention (0.00s)
goos: windows
goarch: amd64
BenchmarkRateLimiter_Parallel-12    	38696997	        27.08 ns/op	      20 B/op	       1 allocs/op
PASS
```

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|---|---|---|---|
| **Lock Striping** | Partitions keyspace across 64 independent mutex shards with clean hashing; zero cross-key contention | Uses multiple shards but hash distribution is skewed or has fewer than 16 shards | Uses single global mutex |
| **Memory Architecture** | Circular ring buffer bounds memory to fixed capacity with zero unbounded slice growth | Ring buffer works but re-allocates slices on purge | Unbounded slice leak |
| **Escape Reduction** | Result returned by value on stack; benchmark shows sub-30ns latency | Returned by pointer or escapes via unnecessary interface cast | High allocation count per op |
| **Correctness** | Enforces burst limits, correctly resets after expiration window, and isolates client tenants | Rate limiting mostly works but window reset calculation drifts | Fails concurrency or throttle tests |

## NEXT UP

[Modular Refactor Review](../../14-architecture-distributed-systems/projects/modular-refactor-review/README.md)
