# Track B: Profiling

## Mission

Learn to investigate the runtime cost of your code. This track introduces Go's powerful profiling tools (`pprof`), escape analysis, and memory layout optimization, turning performance tuning into a data-driven process.

## Track Map

| ID | Topic | Surface | Why It Matters |
| --- | --- | --- | --- |
| `PR.1` | **CPU Profiling** | [`./1-cpu-profile`](./1-cpu-profile) | Finding "hot" functions in offline profiles. |
| `PR.2` | **Live pprof** | [`./3-http-pprof`](./3-http-pprof) | Inspecting production-like traffic in real-time. |
| `PR.3` | **Memory Profiling** | [`./3-memory-profiling`](./3-memory-profiling) | Identifying leak sources and high-allocation paths. |
| `PR.4` | **Escape Analysis** | [`./4-escape-analysis`](./4-escape-analysis) | Understanding Stack vs Heap allocation decisions. |
| `PR.5` | **Benchmark-Driven Dev**| [`./5-benchmark-driven-development`](./5-benchmark-driven-development) | Using benchmarks to drive optimization loops. |
| `PR.6` | **Memory Layout** | [`./6-memory-layout`](./6-memory-layout) | Impact of struct padding and cache locality. |

## Suggested Order

1. **Measurement**: `PR.1` -> `PR.3` (Capturing data).
2. **Analysis**: `PR.4` -> `PR.5` (Understanding the 'Why' and Iterating).
3. **Deep Dive**: `PR.6` (Machine-level optimization).

## Track Milestone

You have mastered this track when you can:
- Generate a CPU profile and identify the most expensive function call.
- Explain how to check if a variable escapes to the heap using `go build -gcflags="-m"`.
- Optimize a struct for better memory usage (e.g., reordering fields).
- Use `go tool pprof` to visualize the call graph.

## Next Up

After mastering performance, move to [09 Architecture & Security](../../../09-architecture) to learn how to structure these efficient components into large systems.
