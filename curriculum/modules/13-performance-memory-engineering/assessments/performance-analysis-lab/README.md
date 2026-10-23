# Project Assessment — Performance Analysis Lab

## Purpose

Validate mastery for the target curriculum items.

## Targets

- `project-performance-analysis-lab`

## Evidence

Submit answers, code evidence, test output, and explanation where required.

---

## 1. Assessment Overview & Industrial Context

Concurrency bottlenecks in Go are rarely exposed during basic unit testing. A single shared mutex guarding a map or state structure will pass hundreds of sequential tests effortlessly. However, when deployed under high parallel concurrency (tens of thousands of requests per second across 16+ CPU cores), goroutines pile up in a lock convoy waiting to acquire the single mutex lock. This causes CPU cache bouncing, lock contention stalls, and throughput collapse.

Compounding this issue, naive implementations that return pointer allocations (`*TokenResult`) from hot loop methods force the Go compiler to allocate results on the heap rather than the stack, generating millions of short-lived objects that trigger frequent garbage collection cycles.

In this assessment, you will analyze a high-throughput rate-limiting engine (`ShardedRateLimiter`). You will use parallel benchmarks (`b.RunParallel`) and escape analysis (`-gcflags="-m"`) to diagnose lock contention and heap escapes, and refactor the engine into a lock-striped architecture with 64 independent shards, fixed circular ring buffers, and zero-allocation value returns.

---

## 2. Core Competencies Tested

- **Parallel Concurrency Benchmarking**: Evaluating lock contention using `testing.B.RunParallel` across multiple CPU threads.
- **Lock Striping & Sharding**: Partitioning global lock contention into an array of 64 independent shards using FNV-1a hashing.
- **Compiler Escape Analysis**: Diagnosing heap escapes using `go build -gcflags="-m"` and converting pointer return types into stack-allocated values (`TokenResult` vs `*TokenResult`).
- **Fixed-Capacity Circular Ring Buffers**: Replacing unbounded, continually allocating slices (`[]time.Time`) with fixed-size circular ring buffers (`[Capacity]int64`) with zero heap churn.
- **Cache-Line Padding & False Sharing**: Understanding CPU cache architecture and avoiding false sharing across concurrent mutexes.

---

## 3. Zero-Magic Architectural Breakdown

```
[ Unsharded Bottleneck ]
Requests (10,000 req/s across 16 cores)
    │   │   │   │   │   │
    └───┴───┴───┼───┴───┘
                ▼
        [ Single sync.Mutex ]  <--- CRITICAL LOCK CONVOY BOTTLENECK
                │
                ▼
        map[string]*ClientBucket (unbounded []time.Time slice, heap escapes)

[ Lock-Striped 64-Shard Architecture ]
Requests
    │
    ▼ FNV-1a Hash(clientKey) % 64
┌───────────┬───────────┬───────────┬─── ... ───┬───────────┐
│ Shard 0   │ Shard 1   │ Shard 2   │           │ Shard 63  │
│ [Mutex 0] │ [Mutex 1] │ [Mutex 2] │           │ [Mutex 63]│
└─────┬─────┴─────┬─────┴─────┬─────┴───────────┴─────┬─────┘
      │           │           │                       │
      ▼           ▼           ▼                       ▼
   RingBuf     RingBuf     RingBuf                 RingBuf
```

### Why FNV-1a Hash Striping Eliminates Contention
Rather than forcing all cores to synchronize on a single memory address, FNV-1a hashing uniformly distributes client keys across 64 shards. Goroutines accessing different client keys acquire independent mutexes, allowing true parallel execution across all CPU cores without lock convoys.

### Stack Allocation vs. Heap Escapes
When a method returns a pointer `*TokenResult`, the compiler cannot prove the object's lifetime ends with the function call, so it forces the allocation to the heap. Returning `TokenResult` by value allows the entire struct to reside directly on the goroutine's execution stack, eliminating heap allocation and GC tracking entirely.

---

## 4. Implementation Requirements & Task Breakdown

### Task 1: Lock Striping with 64 Independent Shards
In `_starter/main.go`:
- Define `clientShard`:
  ```go
  type clientShard struct {
      mu      sync.Mutex
      buckets map[string]*CircularRingWindow
  }
  ```
- Refactor `ShardedRateLimiter` to contain `shards [64]clientShard`.
- Implement `getShard(key string) *clientShard`:
  - Calculate 64-bit FNV-1a hash of `key`.
  - Return `&l.shards[hash % 64]`.
- Lock and unlock only the target shard's mutex during `Allow(key string)`.

### Task 2: Fixed-Size Circular Ring Buffer
- Define `CircularRingWindow`:
  - Fixed-size int64 slice storing Unix nanosecond timestamps.
  - Pointers: `head int`, `count int`, `capacity int`.
- In `Allow(key string) (TokenResult, bool)`:
  - Prune timestamps older than `now - windowDuration`.
  - If `count < capacity`, record `now` at `head = (head + 1) % capacity` and increment `count`.
  - Zero heap allocations in steady state.

### Task 3: Zero-Allocation Stack Value Returns
- Update `Allow` signature to return `TokenResult` by value rather than pointer:
  ```go
  func (l *ShardedRateLimiter) Allow(key string) (TokenResult, bool)
  ```
- Verify with `go build -gcflags="-m"` that `Allow` produces zero heap escapes.

---

## 5. Verification & Acceptance Criteria

Run parallel benchmarks and tests:

```bash
# Verify functional correctness
go test -v ./curriculum/modules/13-performance-memory-engineering/projects/performance-analysis-lab/_starter

# Benchmark parallel performance under heavy concurrency
go test -bench="BenchmarkRateLimiter_Parallel" -benchmem ./curriculum/modules/13-performance-memory-engineering/projects/performance-analysis-lab/_starter
```

### Acceptance Checklist
- [ ] `TestShardedRateLimiter_Accuracy`: verifies rate limiting logic under single-threaded and concurrent load.
- [ ] `BenchmarkRateLimiter_Parallel`: achieves at least an 8x throughput improvement over single-mutex baseline on multi-core systems.
- [ ] Benchmark verifies 0 `B/op` and 0 `allocs/op` on the steady-state `Allow` execution path.
- [ ] `go build -gcflags="-m"` confirms `TokenResult` does not escape to heap.

---

## 6. Common Pitfalls & Anti-Patterns

1. **Deadlocks in Multi-Shard Operations**:
   *Anti-pattern*: Attempting to acquire locks on multiple shards without consistent ordering, triggering circular deadlocks.
   *Correction*: If an operation touches multiple shards, always sort shard indices and acquire locks in ascending index order.

2. **False Sharing across Shards**:
   *Anti-pattern*: Declaring `shards [64]clientShard` where each mutex is adjacent in memory, causing CPU cache lines (64 bytes) to bounce between cores.
   *Correction*: Ensure each shard struct is padded to 64 bytes (`_ [CacheLineSize]byte`) or accessed via pointer indirection to prevent false sharing.

3. **Unbounded Ring Buffers**:
   *Anti-pattern*: Reallocating the ring buffer slice when full.
   *Correction*: Ring buffers must overwrite old entries or reject based on fixed capacity; never resize dynamically during rate limit checks.

---

## 7. 4-Tier Assessment Rubric

| Level | Criteria & Observable Evidence |
|---|---|
| **Level 1: Novice** | Code fails to compile or uses a global lock. Sharding logic panics on negative hash indices. Tests produce race conditions under `go test -race`. |
| **Level 2: Advanced Beginner** | Sharding implemented but uses fewer shards or poor hash distribution. Returns pointer allocations that escape to the heap. Memory continues to grow due to unbounded slices. |
| **Level 3: Competent** | 64 shards partitioned via FNV-1a. Fixed circular ring buffer implemented. Stack value returns achieve 0 allocs/op in benchmarks. Multi-threaded benchmark throughput scales cleanly. |
| **Level 4: Proficient** | Exceptional concurrency engineering. Addresses false sharing with cache-line padding. Proves zero-allocation guarantee via compiler escape analysis. Articulates throughput limits based on CPU memory bus architecture. |

---

## 8. Submission & Reflection

1. Submit benchmark output for both sequential and parallel rate limiter executions.
2. Submit the compiler escape analysis log snippet for `Allow`.
3. Answer: *Why is FNV-1a preferred over cryptographic hash functions (such as SHA-256) for internal shard key distribution?*
