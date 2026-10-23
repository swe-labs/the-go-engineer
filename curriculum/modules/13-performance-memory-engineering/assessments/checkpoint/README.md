# Module 13 Assessment — Performance and Memory Engineering

## Purpose

Validate mastery for the target curriculum items.

## Targets

- `core-13-01`
- `core-13-02`
- `core-13-03`
- `core-13-04`
- `core-13-05`
- `core-13-06`
- `core-13-07`
- `core-13-08`

## Evidence

Submit answers, code evidence, test output, and explanation where required.

---

## 1. Deep-Dive Memory & Performance Diagnostic Scenarios

Analyze and answer each technical scenario, referencing Go runtime internals (mcache, mcentral, mheap, GC pacer, write barrier) and profiling methodologies.

### Scenario 1: Stack vs. Heap Allocation & Escape Analysis (`core-13-01`, `core-13-02`)
Consider the two functions below:
```go
func createItemByPointer(name string, value int) *Item {
    return &Item{Name: name, Value: value}
}

func createItemByValue(name string, value int) Item {
    return Item{Name: name, Value: value}
}
```
- Explain step-by-step how the Go compiler (`cmd/compile`) decides whether an allocated variable resides on the goroutine execution stack or escapes to the shared runtime heap.
- What compiler flags reveal escape analysis decisions? Explain the diagnostic output `moved to heap: item` and `leaking param: name`.
- Under what circumstances can a pointer parameter *not* escape to the heap? Give a concrete code example.

### Scenario 2: Go Garbage Collector Mechanics & Pacer Tuning (`core-13-03`, `core-13-04`)
A real-time financial order-matching engine in Go experiences periodic 15ms latency spikes that correlate directly with GC cycles, despite running on 32-core nodes with 64GB RAM.
- Explain the phases of Go's tri-color concurrent mark-and-sweep garbage collector (Sweep Termination, Concurrent Mark, Mark Termination). Which phases require STW (Stop-The-World)?
- What is `GOGC` (default `100`), and how does the GC Pacer decide when to trigger the next GC cycle?
- How does `GOMEMLIMIT` (introduced in Go 1.19) prevent container OOM killer evictions while stabilizing GC frequency under fluctuating memory pressure?

### Scenario 3: `sync.Pool` Mechanics & GC Interactions (`core-13-05`)
An engineer implements a global `sync.Pool` to reuse 1MB byte buffers across workers. During testing under erratic traffic bursts, memory spikes dramatically and the application experiences sudden latency hiccups.
- Explain how `sync.Pool` stores objects per-P (processor) using private and shared pools to eliminate mutex contention.
- What happens to pooled objects during a Garbage Collection cycle? How does pool victim cache (introduced in Go 1.13) smooth out cache clearing across multiple GC cycles?
- Why should pooled buffers that exceed a maximum threshold not be returned to `sync.Pool`?

### Scenario 4: CPU Cache Hierarchy, Padding, and False Sharing (`core-13-06`)
A high-throughput counter array `[8]uint64` updated concurrently by 8 goroutines on separate CPU cores runs significantly *slower* than 8 independent local variables.
- Explain the concept of a CPU cache line (typically 64 bytes) and the phenomenon of **False Sharing**.
- How do you use cache-line padding (e.g. `_ [56]byte` or `cpu.CacheLinePad`) in Go structs to prevent cache bouncing across CPU cores?
- When does cache-line padding become counterproductive (e.g. memory cache footprint inflation)?

### Scenario 5: Contention Profiling & Micro-Benchmarking (`core-13-07`, `core-13-08`)
An engineer writes a benchmark:
```go
func BenchmarkMap(b *testing.B) {
    m := make(map[int]int)
    for i := 0; i < b.N; i++ {
        m[i] = i
    }
}
```
- Identify at least three flaws in this benchmark methodology (setup overhead included in timing, map reallocation amortized over `b.N`, lack of `b.ResetTimer()`, compiler dead-code elimination).
- How do you profile mutex contention using `runtime.SetMutexProfileFraction` and `go tool pprof -mutex`?
- How does `testing.B.RunParallel` evaluate multi-threaded lock scaling compared to standard `b.N` iteration loops?

---

## 2. Practical Memory Optimization Challenge

Given an unoptimized JSON log parser that decodes records, normalizes fields, and filters errors, write an optimized Go implementation that:
1. Reuses scratch memory buffers using `sync.Pool`.
2. Uses `json.NewDecoder` or zero-allocation byte scanning instead of repeated `json.Unmarshal([]byte)` on ephemeral slices.
3. Preallocates output slices using slice header capacity hints.
4. Validates zero heap allocation in the hot loop using `testing.AllocsPerRun`.

---

## 3. Invariant Performance Guidelines

1. **Profile Before Optimizing**: Never guess where bottlenecks exist; prove them using `pprof` flame graphs or benchmark allocation reports (`-benchmem`).
2. **Value Semantics First**: Default to passing and returning values on the stack unless ownership transfer, mutation, or struct size (>128 bytes) explicitly requires pointer semantics.
3. **Guard Pooled Capacity**: Never return arbitrarily large buffers to a `sync.Pool`; enforce a strict capacity threshold.
4. **Isolate Diagnostic Ports**: Never mount `net/http/pprof` handlers on unauthenticated public HTTP endpoints.

---

## 4. 4-Tier Assessment Rubric

| Level | Criteria & Observable Evidence |
|---|---|
| **Level 1: Novice** | Believes pointers are always faster than values; writes microbenchmarks without `b.ResetTimer`; cannot read escape analysis compiler output; unaware of STW GC pauses. |
| **Level 2: Advanced Beginner** | Can capture pprof profiles but struggles to interpret flame graphs; applies `sync.Pool` without resetting buffers; understands GC basics but cannot tune `GOGC` or `GOMEMLIMIT`. |
| **Level 3: Competent** | Explains tri-color GC phases accurately; eliminates heap escapes systematically; implements `sync.Pool` with capacity guards; profiles mutex contention using `pprof -mutex`. |
| **Level 4: Proficient** | Masters systems-level memory engineering; articulates cache line mechanics and resolves false sharing; tunes GC pacer parameters for low-latency SLAs; proves zero-allocation guarantees with compiler tooling. |

---

## 5. Submission Requirements

Submit your detailed answers for Scenarios 1–5 alongside your verified code and benchmark output for the practical optimization challenge.
