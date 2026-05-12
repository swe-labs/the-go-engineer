# Track B: Profiling

## Mission

Learn to investigate the runtime cost of your code. This track introduces Go's powerful profiling tools (`pprof`), escape analysis, and memory layout optimization, turning performance tuning into a data-driven process.

## Track Map

| ID | Topic | Surface | Why It Matters |
| --- | --- | --- | --- |
| `PR.1` | **CPU Profiling** | [`./01-cpu-profile`](./01-cpu-profile) | Finding "hot" functions in offline profiles. |
| `PR.2` | **Live pprof** | [`./02-http-pprof`](./02-http-pprof) | Inspecting production-like traffic in real-time. |
| `PR.3` | **Memory Profiling** | [`./03-memory-profiling`](./03-memory-profiling) | Identifying leak sources and high-allocation paths. |
| `PR.4` | **Escape Analysis** | [`./04-escape-analysis`](./04-escape-analysis) | Understanding Stack vs Heap allocation decisions. |
| `PR.5` | **Benchmark-Driven Dev**| [`./05-benchmark-driven-development`](./05-benchmark-driven-development) | Using benchmarks to drive optimization loops. |
| `PR.6` | **Memory Layout** | [`./06-memory-layout`](./06-memory-layout) | Impact of struct padding and cache locality. |

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
