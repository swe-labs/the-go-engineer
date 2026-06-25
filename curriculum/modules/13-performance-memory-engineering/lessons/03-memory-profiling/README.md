# Memory profiling

## Learning objective

Collect and interpret Go heap profiles, distinguish between alloc_objects, alloc_space, inuse_objects, and inuse_space views, detect memory leaks by comparing heap profiles over time, and identify the call sites responsible for excessive allocation.

## Why this matters

Memory leaks in Go are rare but devastating: a service that grows by 1 MB per minute will OOM in 16 hours on a 1 GB container. Unlike CPU hotspots (which cause latency), memory leaks cause crashes. Heap profiling is the only reliable way to find the allocation call site behind a leak. Even without leaks, understanding allocation patterns helps reduce GC pressure, which directly impacts latency — every GC pause steals CPU from request processing.

## Mental model

A heap profile is a photograph of the heap at a single moment. Each sample represents approximately 512 KB of live allocation. If a call site appears in 10 samples, it is responsible for roughly 5 MB of live heap. The profile is a probability-weighted estimate: allocation sites that allocate rarely may not appear in any sample, and heavy allocators appear in proportion to their allocation volume. Profiles collected at different times can be compared: if a call site's inuse_space grows between profile A and profile B, that site is leaking.

## Core idea

The Go runtime tracks allocation at every `mallocgc` call using sampling. The heap profile can be viewed through four lenses:

| View | `-sample_index` | What it measures | Use case |
|---|---|---|---|
| `inuse_space` | 0 | Bytes of live (unfreed) memory at profile time | Memory leak detection |
| `inuse_objects` | 1 | Count of live objects at profile time | Identifying many small allocations |
| `alloc_space` | 2 | Cumulative bytes allocated since process start | Understanding GC pressure |
| `alloc_objects` | 3 | Cumulative count of allocations since start | Identifying allocation frequency |

The default view for `go tool pprof` on a heap profile is `inuse_space`. For leak detection, always use `inuse_space`. For understanding allocation rate, use `alloc_space` with `-diff_base` to compare two profiles.

## Under the hood

Heap profiling lives in `runtime/mprof.go`. Each M (OS thread) maintains a per-thread allocation counter. When `mallocgc` is called, it increments the counter by the allocation size. When the counter exceeds `MemProfileRate / 2` (default 512 KB / 2 = 256 KB), the runtime generates a random number between 0 and `MemProfileRate - 1`. If the counter exceeds this random number, a sample is taken. This produces an average sampling rate of one sample per `MemProfileRate` bytes allocated.

On a sample, `recordstacktrace` captures up to 32 PC values from the current goroutine's stack. The record is stored in a global slice of `MemProfileRecord` entries. Each record contains: `AllocBytes` (total bytes allocated at this site), `FreeBytes` (bytes freed), `AllocObjects` (count allocated), `FreeObjects` (count freed), and the stack trace. On `WriteHeapProfile`, the runtime iterates all records, computes `inuse = AllocBytes - FreeBytes`, and serializes as a pprof profile.

The key sampling invariance: the probability of sampling a particular allocation is proportional to its size, not its frequency. A 1 KB allocation is 1000x more likely to be sampled than a 1-byte allocation.

## How Go uses it

- `runtime/pprof.WriteHeapProfile` writes the standard heap profile.
- `net/http/pprof` exposes `/debug/pprof/heap` for HTTP-accessible profiling.
- `runtime.ReadMemStats` provides programmatic access to the same allocation counters without profile serialization.
- The `go test -memprofile=mem.out` flag collects a heap profile of test execution.
- `go tool pprof` with `-base` flag compares two profiles (useful for detecting leaks between two points in time).

## Go example

```go
package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

type LeakHolder struct {
	data [][]byte
}

func (lh *LeakHolder) leak(size int) {
	buf := make([]byte, size)
	lh.data = append(lh.data, buf)
}

func (lh *LeakHolder) noLeak(size int) {
	_ = make([]byte, size)
}

func main() {
	var holder LeakHolder

	for i := 0; i < 100; i++ {
		holder.leak(1 << 20)
	}

	runtime.GC()

	f, err := os.Create("heap_profile.pprof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := pprof.WriteHeapProfile(f); err != nil {
		panic(err)
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("HeapAlloc: %d MB\n", m.HeapAlloc/1024/1024)
	fmt.Printf("Heap profile written to heap_profile.pprof\n")

	_ = holder
}
```

Analyze with:

```bash
go tool pprof -inuse_space heap_profile.pprof
(pprof) top10
(pprof) list main.LeakHolder.leak
```

## Step-by-step execution

1. `LeakHolder.leak` allocates a 1 MB `[]byte` slice and appends it to `LeakHolder.data`. The slice's backing array is heap-allocated because it outlives the `leak` call (stored in the holder's field).
2. After 100 iterations, `holder.data` holds 100 MB of live allocations.
3. `runtime.GC()` runs a garbage collection, freeing any unreachable allocations. All 100 MB in `holder.data` is still reachable, so it survives.
4. `pprof.WriteHeapProfile(f)` iterates the global `MemProfileRecord` slice, filters to records with non-zero `inuse`, and writes them to the file.
5. `runtime.ReadMemStats(&m)` reads the same counters directly — `HeapAlloc` reports ~100 MB.
6. `go tool pprof -inuse_space heap_profile.pprof` shows `main.LeakHolder.leak` as the top entry with approximately 100 MB of `inuse_space`.

## Common mistakes

- **Reading alloc_space instead of inuse_space** — `alloc_space` shows cumulative bytes allocated since process start. For leak detection, `inuse_space` is the relevant metric. Use `-sample_index=0` or `-inuse_space` explicitly.
- **Profiling at the wrong point in the request lifecycle** — A heap profile collected during peak request processing shows inflated `inuse_space` because many requests are in-flight. Collect profiles during idle periods to see the baseline, then compare with peak profiles.
- **Forgetting that heap profiles show allocation call sites, not current owners** — A profile shows which function allocated the memory, but not which code path currently holds the reference. A buffer allocated by `conn.setup` and held by a leaked goroutine appears as "allocated by conn.setup" even though the goroutine is the root cause.
- **Not calling `runtime.GC()` before profiling** — Without a GC, the heap profile includes unreachable-but-uncollected objects as `inuse`, producing false positives. Always call `runtime.GC()` before `WriteHeapProfile` for accurate `inuse` data.

## Debugging walkthrough

A service OOMs after 6 hours of uptime. Create two heap profiles — one at startup (baseline), one just before OOM:

```bash
curl -o baseline.pprof http://localhost:6060/debug/pprof/heap
# wait 6 hours
curl -o before_oom.pprof http://localhost:6060/debug/pprof/heap
go tool pprof -base baseline.pprof before_oom.pprof
```

```
(pprof) top
  flat  flat%   sum%        cum   cum%
  450MB 90.00% 90.00%     450MB 90.00%  main.cacheEntries
```

The `-base` diff shows which call sites grew between the two profiles. Here, `main.cacheEntries` grew by 450 MB. The function likely populates a global map without eviction. The fix: add TTL eviction or LRU limits to the cache.

## Production notes

- Heap profiling sampling rate is controlled by `runtime.MemProfileRate`. The default of 512 KB balances accuracy with overhead. Lowering it increases accuracy but adds allocation overhead.
- Heap profiles include only live allocations after `runtime.GC()`. For understanding total allocation rate (including freed memory), use `go tool pprof -alloc_space` or query `runtime.MemStats.TotalAlloc`.
- In Kubernetes, collect heap profiles periodically (every 5 minutes) during incident investigation. Compare consecutive profiles to detect growing allocations.
- Continuous profiling tools (Datadog, Google Cloud Profiler) collect heap profiles from production services and surface trends over time.

## Performance implications

Heap profiling adds overhead proportional to the sampling rate. At the default `MemProfileRate` of 512 KB, the overhead is approximately 1-3% for allocation-heavy workloads. Each sample captures a stack trace (up to 32 frames), stores it in a hash table, and eventually serializes it. The serialization itself (calling `WriteHeapProfile`) is CPU-intensive and can take hundreds of milliseconds for large heaps. Do not call `WriteHeapProfile` on every request.

## Practice task

Write a Go program that:

1. Creates a `sync.Map` and stores 100,000 entries with 128-byte values.
2. Collects a heap profile to `before.pprof`.
3. Adds another 100,000 entries to the same map.
4. Collects a heap profile to `after.pprof`.
5. Runs `go tool pprof -base before.pprof after.pprof` and identifies which function grew.

## Tests / verification

```bash
go run ./curriculum/modules/13-performance-memory-engineering/lessons/03-memory-profiling
go test ./curriculum/modules/13-performance-memory-engineering/lessons/03-memory-profiling
```

## Review questions

1. What is the default heap sampling rate in Go? How is each sample triggered?
2. Why should you call `runtime.GC()` before `WriteHeapProfile` when looking for memory leaks?
3. What is the difference between `inuse_space` and `alloc_space` in a heap profile?
4. How does `go tool pprof -base` help detect memory leaks?
5. A function allocates many small objects frequently but uses little total memory. Which heap profile view best reveals this function?

## NEXT UP

Benchmarks — writing benchmark functions with `testing.B`, using the `-bench` flag, analyzing results with `benchstat`, measuring allocation with `-benchmem`, and comparing performance across code changes.
