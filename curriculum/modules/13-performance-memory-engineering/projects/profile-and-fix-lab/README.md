# Profile and Fix Lab

## Learning objective

Profile a CPU and memory bottlenecked Go batch-processing application using Go's built-in `pprof` diagnostic tooling. Identify heap allocation hotspots, escape analysis decisions, and memory churn using flame graphs, heap profiles, and memory benchmarks. Apply targeted mechanical optimizations—including object pooling with `sync.Pool`, preallocation of slice headers, and buffered string building—to achieve measurable speedup and allocation reductions.

## Why this matters

Writing code that compiles and passes unit tests is only table stakes. Under high production concurrency or high-throughput event streaming, naive allocations and unpooled byte buffers trigger severe garbage collection (GC) latency spikes and heap exhaustion (OOM). Learning to diagnose bottlenecks scientifically with `pprof` flame graphs and benchmark profiling before modifying a single line of code is an essential senior Go engineering capability.

## Prerequisites

Complete Module 13 lessons on Go memory model, heap vs stack allocation, escape analysis (`go build -gcflags="-m"`), runtime garbage collection mechanics, and `pprof` profiling.

## Available concepts

This project draws on `runtime/pprof`, `net/http/pprof`, `sync.Pool`, `strings.Builder`, slice capacity preallocation (`make([]T, 0, cap)`), `testing.B` benchmarks with `b.ReportAllocs()`, and `go tool pprof` CLI/web flame graphs.

## Project overview

The `_starter/` directory contains an unoptimized log event batch processor (`EventProcessor`). It successfully aggregates JSON log streams into summary reports, but exhibits severe performance and allocation anti-patterns:
- Naive string concatenation in inner loops (`digest += s.Service + ...`) causing repeated heap copies.
- Unpooled slice buffers for JSON serialization, triggering heap churn on every batch.
- Slice `append` without capacity hints, causing repeated geometric reallocations.
- Missing live `net/http/pprof` endpoints for runtime inspection.

Your mission is to profile the application, identify the root causes of allocation churn, and apply targeted Go optimizations without altering the functional output.

## Profiling Workflow

### 1. Collect Benchmark Profiles

Run the benchmark suite with CPU and memory profiling enabled:

```bash
go test -bench="." -benchmem -cpuprofile=cpu.pprof -memprofile=mem.pprof ./curriculum/modules/13-performance-memory-engineering/projects/profile-and-fix-lab/_starter
```

### 2. Inspect Hotspots with pprof CLI

Examine the top memory allocation functions:

```bash
go tool pprof -top mem.pprof
```

Or open the interactive flame graph in your browser:

```bash
go tool pprof -http=:8081 cpu.pprof
```

### 3. Check Compiler Escape Analysis

Inspect which variables escape from stack to heap:

```bash
go build -gcflags="-m" ./curriculum/modules/13-performance-memory-engineering/projects/profile-and-fix-lab/_starter
```

## Tasks

### Task 1 — Eliminate Heap Churn with `sync.Pool`

In `_starter/main.go`:
- Define a package-level or processor-level `sync.Pool` that yields reusable `*bytes.Buffer` instances initialized with `make([]byte, 0, 4096)`.
- In `FormatReportJSON`:
  - Acquire a buffer from the pool via `.Get().(*bytes.Buffer)`.
  - Reset the buffer with `buf.Reset()`.
  - Ensure the buffer is returned to the pool with `defer pool.Put(buf)`.
  - Encode directly into the pooled buffer and stream to the `io.Writer`.

### Task 2 — Eliminate String Concat and Preallocate Slices

In `ProcessBatch`:
- Replace dynamic `var services []ServiceStat` with preallocated capacity:
  ```go
  services := make([]ServiceStat, 0, len(serviceMap))
  ```
- Replace naive string concatenation (`+` and `+=`) with `strings.Builder`:
  - Call `b.Grow(...)` with an estimate of the required capacity.
  - Use `b.WriteString` and `b.WriteByte` to construct `digest_summary`.

### Task 3 — Expose Live `net/http/pprof` Endpoints

- Import `_ "net/http/pprof"`.
- Implement `StartPprofServer(addr string) *http.Server` to launch a non-blocking background HTTP server on the designated address (e.g., `:6060`).
- Ensure `/debug/pprof/` responds to HTTP GET with status 200 OK.

## Verification

Run the automated test and benchmark suite against the solution:

```bash
go test -v -bench="." ./curriculum/modules/13-performance-memory-engineering/projects/profile-and-fix-lab/_solution
```

Expected output:
```
=== RUN   TestProcessBatch_Correctness
--- PASS: TestProcessBatch_Correctness (0.00s)
=== RUN   TestProcessBatch_EdgeCases
--- PASS: TestProcessBatch_EdgeCases (0.00s)
=== RUN   TestPprofServer_Endpoints
--- PASS: TestPprofServer_Endpoints (0.00s)
goos: windows
goarch: amd64
BenchmarkProcessBatch-12         17779         67328 ns/op        9526 B/op        118 allocs/op
PASS
```

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|---|---|---|---|
| **Profiling Evidence** | Identifies exact allocation hotspots using `pprof` top and flame graph analysis | Identifies general area of slowness without profile metrics | No profiling performed |
| **Allocation Reductions** | Uses `sync.Pool`, `strings.Builder.Grow`, and slice preallocation to slash heap allocations | Uses `strings.Builder` but ignores buffer pooling and slice preallocation | Remains heavily bottlenecked with string concat |
| **pprof Integration** | Exposes `/debug/pprof/*` endpoints on a dedicated administration port | Mounts pprof but blocks the main execution goroutine | No pprof endpoints |
| **Output Equivalence** | Produces byte-for-byte identical JSON reports and digest strings | Output correct but formatting or sorting unstable | Breaks functional correctness |

## NEXT UP

[Performance Analysis Lab](../performance-analysis-lab/README.md)
