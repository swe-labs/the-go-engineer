# Benchmarks

## Learning objective

Write and run Go benchmark functions using `testing.B`, interpret benchmark output, use `-benchmem` to measure allocations, compare benchmark results with `benchstat`, and apply benchmarks to prevent performance regressions.

## Why this matters

Benchmarks are the only way to objectively measure code performance. Without benchmarks, engineers rely on intuition about what is fast vs slow — and intuition is systematically wrong. Benchmarks catch regressions before they reach production, validate optimization efforts, and provide the data needed to make evidence-based engineering decisions. In Go, benchmarks are a first-class citizen: every `go test` invocation supports benchmarks with no additional tools.

## Mental model

A benchmark is a stopwatch that runs a function many times and reports the average time per run. The framework finds the right number of repetitions to get a stable average — just like measuring a runner's lap time by timing 100 laps, not 1 lap. The key insight: benchmarks run the same code path repeatedly to cancel out measurement noise from GC, OS scheduling, and CPU frequency scaling. The result is a statistically stable measurement of execution time per operation.

## Core idea

A benchmark function in Go has the signature `func BenchmarkXxx(b *testing.B)`. The `testing.B` type provides:

| Field/Method | Purpose |
|---|---|
| `b.N` | The number of iterations to run. Set by the framework — do not modify. |
| `b.ResetTimer()` | Zero the benchmark timer, discarding setup time. |
| `b.ReportAllocs()` | Enable allocation reporting for this benchmark. |
| `b.Run(name, fn)` | Run a sub-benchmark, enabling hierarchical organization. |
| `b.SetBytes(n)` | Report throughput in bytes/sec (for I/O-like benchmarks). |
| `b.StartTimer()` / `b.StopTimer()` | Manually control the timer interval. |

The framework runs the benchmark with increasing `b.N` until the total duration exceeds the benchmark time (default 1 second). The reported `ns/op` is the total wall-clock time divided by `b.N`.

## Under the hood

When `go test -bench=.` runs, the testing framework parses benchmark function names, creates a `testing.B` for each, and calls the benchmark function. The framework starts with `b.N = 1` and runs the function once. It measures the wall-clock duration with `time.Now()`. If the duration is below the minimum benchmark time (default 1 second), it computes a scaling factor: `target_duration / measured_duration * 1.1`. It then runs again with the estimated N. It repeats this calibration loop until either the runtime is above the minimum time or N stops growing.

Internally, `testing.B` calls `runtime.ReadMemStats` before and after the benchmark loop to compute `B/op` and `allocs/op`. The reported value is the delta divided by `b.N`. The framework also tracks the number of iterations completed and reports the minimum time across runs (not the average) to reduce noise from GC or OS scheduling.

## How Go uses it

- The standard library maintains benchmarks for virtually every package. Changes to `strings`, `fmt`, `encoding/json`, and `net/http` are benchmarked on every commit.
- `benchstat` (part of `golang.org/x/perf`) computes statistical summaries and detects significant performance changes.
- `go test -bench=. -benchmem` runs all benchmarks with allocation reporting.
- `go test -bench=. -count=10` runs each benchmark 10 times, enabling statistical analysis with benchstat.
- `go test -bench=. -benchtime=10s` extends the minimum benchmark duration for more stable measurements.

## Go example

```go
package main

import (
	"fmt"
	"strings"
)

func concatPlus(parts []string) string {
	result := ""
	for _, p := range parts {
		result += p
	}
	return result
}

func concatBuilder(parts []string) string {
	var b strings.Builder
	for _, p := range parts {
		b.WriteString(p)
	}
	return b.String()
}

func main() {
	parts := []string{"hello", " ", "world", " ", "from", " ", "Go"}
	fmt.Println("concatPlus:", concatPlus(parts))
	fmt.Println("concatBuilder:", concatBuilder(parts))
}
```

Benchmarks:

```go
func BenchmarkConcatPlus(b *testing.B) {
	parts := []string{"hello", " ", "world", " ", "from", " ", "Go"}
	for i := 0; i < b.N; i++ {
		concatPlus(parts)
	}
}

func BenchmarkConcatBuilder(b *testing.B) {
	parts := []string{"hello", " ", "world", " ", "from", " ", "Go"}
	for i := 0; i < b.N; i++ {
		concatBuilder(parts)
	}
}
```

Run with:

```bash
go test -bench=. -benchmem
```

Expected output pattern:

```
BenchmarkConcatPlus-8      5000000   320 ns/op   96 B/op   7 allocs/op
BenchmarkConcatBuilder-8  20000000    85 ns/op   64 B/op   1 allocs/op
```

`strings.Builder` is ~4x faster and makes 1 allocation instead of 7.

## Step-by-step execution

1. `go test -bench=.` discovers `BenchmarkConcatPlus` and `BenchmarkConcatBuilder`.
2. For `BenchmarkConcatPlus`, the framework sets `b.N = 1` and calls the function.
3. The loop runs once: `concatPlus` builds the string using `+=` (which allocates a new string on each concatenation).
4. The framework measures ~300 ns. Since this is below 1 second, it computes `N = 1s / 300ns * 1.1 ≈ 3,600,000`.
5. The framework runs with `b.N = 3,600,000`, calls `runtime.ReadMemStats`, runs the loop, calls `ReadMemStats` again, and measures the delta.
6. Results: `ns/op = total_time / N`, `B/op = (alloc_bytes_after - alloc_bytes_before) / N`, `allocs/op = (alloc_count_after - alloc_count_before) / N`.
7. For `BenchmarkConcatBuilder`, the same process yields fewer allocations (1 vs 7) and lower `ns/op`.

## Common mistakes

- **Benchmarking I/O-bound work** — The benchmark framework measures wall-clock time. If the function makes HTTP calls or reads from disk, the reported `ns/op` includes I/O wait time, which is non-deterministic. Use `-benchtime` and `-count` for I/O benchmarks, or mock the I/O layer.
- **Not resetting the timer** — If the benchmark function has setup code (file reads, allocation of test data), that setup time is included unless `b.ResetTimer()` is called before the loop.
- **Benchmarking with compiler optimizations disabled** — Running with `-gcflags='-N -l'` measures unoptimized code. Always benchmark with default compiler flags.
- **The compiler eliminating the benchmark** — If the benchmarked function has no side effects and returns a value that is never used, the compiler may eliminate the call. Assign the result to a package-level variable: `var result = compute()`.
- **Insufficient iterations** — A single run of a fast function (nanoseconds) has high relative noise. Let the framework determine `b.N` rather than setting it manually.

## Debugging walkthrough

A microservice is slower after adding a new feature. Before debugging the business logic, write a benchmark for the new code path:

```go
func BenchmarkNewFeature(b *testing.B) {
	payload := generateLargePayload()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processPayload(payload)
	}
}
```

Run with benchstat:

```bash
go test -bench=BenchmarkNewFeature -count=10 -benchmem > old.txt
# apply code fix
go test -bench=BenchmarkNewFeature -count=10 -benchmem > new.txt
benchstat old.txt new.txt
```

Output:

```
name             old time/op    new time/op    delta
NewFeature-8      12.5µs ± 3%    8.2µs ± 2%   -34.4%

name             old alloc/op   new alloc/op   delta
NewFeature-8      4.2kB ± 0%    2.1kB ± 0%   -50.0%
```

`benchstat` reports the mean, the variation (±), and the delta with statistical significance. A delta of -34.4% with low variance confirms the optimization worked.

## Production notes

- Add benchmarks to every performance-sensitive package. Run them in CI on every pull request.
- Use `benchstat` to compare benchmark results across commits. A regression of >5% should block the PR.
- Store benchmark results in a database (e.g., `golang.org/x/perf/storage`) to track performance over time.
- Benchmark on hardware that matches production. CPU frequency scaling, hyperthreading, and cache size affect results. Use `-count=10` to account for variance.
- For allocation-sensitive code, `-benchmem` is essential. A 2x reduction in allocs/op often yields a larger latency improvement than the ns/op numbers suggest because of reduced GC pressure.

## Performance implications

Benchmarks themselves have negligible overhead — the framework is highly optimized. However, the code being benchmarked must be production-representative. Common pitfalls: benchmarking only the happy path (fast), benchmarking with tiny inputs (unrealistic), or benchmarking without GC pressure (production has concurrent allocations). The `-benchtime` flag on `testing.B` controls the minimum sampling duration; longer times reduce noise but increase CI runtime. For most projects, a `-benchtime=2s` and `-count=5` provides sufficient accuracy.

## Practice task

Write a Go program with three functions:

1. `sumRange(n int) int` — sums integers 1 to n using a `for` loop.
2. `sumFormula(n int) int` — computes `n * (n + 1) / 2` directly.
3. `sumRecursive(n int) int` — sums recursively (no tail-call optimization).

Write benchmarks for all three with `benchmem`. Run `benchstat` to compare. Which is fastest? Which allocates the most?

## Tests / verification

```bash
go run ./curriculum/modules/13-performance-memory-engineering/lessons/04-benchmarks
go test ./curriculum/modules/13-performance-memory-engineering/lessons/04-benchmarks
go test -bench=. ./curriculum/modules/13-performance-memory-engineering/lessons/04-benchmarks
```

## Review questions

1. What is `b.N` and how does the testing framework determine its value?
2. Why should you call `b.ResetTimer()` after setup code in a benchmark?
3. What does `-benchmem` report, and why is it important?
4. How does `benchstat` differ from manually comparing benchmark output?
5. What happens if a benchmark function's result is discarded and the compiler eliminates the call? How do you prevent this?

## NEXT UP

Escape analysis — understanding stack vs heap allocation decisions made by the Go compiler, using `-gcflags=-m` to inspect escape decisions, and writing allocation-efficient code by keeping values on the stack.
