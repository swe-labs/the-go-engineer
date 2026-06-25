# CPU profiling

## Learning objective

Collect CPU profiles from Go programs, interpret the output of `go tool pprof`, identify hotspots using `top`, `list`, and `web` commands, and read flame graphs to prioritize optimization targets.

## Why this matters

CPU time is the most constrained resource in compute-bound services. Every microsecond of unnecessary CPU work compounds under load: a function that takes 10 µs per call becomes 10 ms of CPU per 1,000 requests. CPU profiling reveals exactly which functions consume CPU time, replacing guesswork with evidence. Professional Go optimizers never optimize without a CPU profile — the profile tells you where to invest effort for the greatest return.

## Mental model

CPU profiling is like a photographer taking a flash photo of the CPU every 10 ms. Each photo captures "what function is the CPU executing right now?" After 3,000 photos (30 seconds), the stack trace that appears in 37% of photos likely consumes approximately 37% of CPU. It is not exact — sampling error exists — but with enough samples, the error is small enough to guide optimization decisions. Functions that never appear in any photo consume negligible CPU and should not be optimized first.

## Core idea

A CPU profile is a collection of stack trace samples, each representing the program's call stack at an instant when the CPU was executing user code. The `go tool pprof` interactive shell provides four main views:

| Command | What it shows |
|---|---|
| `top` | Flat and cumulative CPU consumption per function, sorted by flat time |
| `list <func>` | Per-line CPU consumption within a specific function |
| `web` | An SVG call graph where node size = CPU consumption, edge width = call frequency |
| `peek <func>` | Shows callers and callees of a function with their CPU contributions |

The key distinction in every view: **flat** time is samples where the function was at the top of the call stack (executing its own code). **Cumulative** time includes samples where the function was anywhere on the call stack (it called something that consumed CPU). A high flat + low cumulative time means the function is CPU-intensive. Low flat + high cumulative means the function is a caller of CPU-intensive code.

## Under the hood

When `pprof.StartCPUProfile` is called, the Go runtime programs the OS interval timer (via `setitimer` on Linux, `ITIMER_PROF`) to deliver `SIGPROF` at 100 Hz. The signal handler (`sigprof` in `runtime/signal_unix.go`) reads the interrupted thread's `m` (OS thread) and `gp` (goroutine pointer) from thread-local storage. It calls `gentraceback` to unwind the goroutine's stack: starting from the PC/SP saved in the signal context, it walks frame-by-frame using Go's stack map metadata. Each frame's PC is mapped to a function name and line number via the `gopclntab` (Go Program Counter Line Number Table). The resulting stack trace is hashed with FNV-1a into a fixed-size hash table. On `StopCPUProfile`, the runtime serializes all buckets into protobuf format.

## How Go uses it

- The `runtime/pprof` package is the programmatic API for CPU profiling.
- The `net/http/pprof` package exposes CPU profiling via HTTP: `GET /debug/pprof/profile?seconds=30`.
- `go test -cpuprofile=cpu.out` collects a CPU profile of test execution.
- `go tool pprof` is the analysis CLI. It reads protobuf profiles and supports interactive and non-interactive modes.
- Third-party tools like `Graphviz` (for `web` command) and `FlameGraph` scripts enable richer visualizations.

## Go example

```go
package main

import (
	"fmt"
	"math"
	"os"
	"runtime/pprof"
)

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func countPrimes(limit int) int {
	count := 0
	for i := 0; i < limit; i++ {
		if isPrime(i) {
			count++
		}
	}
	return count
}

func main() {
	f, err := os.Create("cpu_profile.pprof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	result := countPrimes(100000)
	fmt.Printf("Primes up to 100000: %d\n", result)
}
```

After running, inspect with:

```bash
go tool pprof -top cpu_profile.pprof
go tool pprof -web cpu_profile.pprof     # requires Graphviz
```

## Step-by-step execution

1. `pprof.StartCPUProfile(f)` enables sampling. The runtime begins receiving SIGPROF signals.
2. `countPrimes(100000)` starts iterating from 0 to 100,000.
3. For each `i`, `isPrime(i)` is called. Inside `isPrime`, `math.Sqrt` computes the square root and the loop checks divisibility.
4. The signal handler fires approximately every 10 ms. If the CPU is inside `isPrime` when the signal arrives, that call stack (countPrimes → isPrime) is recorded.
5. After processing 100,000 numbers, `main` proceeds to the defer.
6. `pprof.StopCPUProfile()` serializes the hash table of samples to `cpu_profile.pprof`.
7. Running `go tool pprof -top cpu_profile.pprof` shows the distribution. `isPrime` dominates because the inner loop executes most iterations.

## Common mistakes

- **Profiling a cold process** — A CPU profile collected when the service just started and has no request traffic shows runtime startup code (GC initialization, map zeroing) that disappears after warmup. Always profile under realistic load.
- **Optimizing functions that consume <5% of CPU** — Optimizing a function at 3% improves overall performance by at most 3%, but the optimization may introduce bugs. Focus on the top 3-5 functions that cumulatively account for >50% of samples.
- **Misreading flat vs cumulative time** — Flat time is samples where the function was on top of the stack. Cumulative includes samples where it was anywhere on the stack. A high flat means the function is CPU-intensive. A low flat but high cumulative means it calls something expensive.
- **Using `-sample_index` without understanding the units** — CPU profiles have one sample type (CPU samples). Memory profiles have four (alloc_objects, alloc_space, inuse_objects, inuse_space). Mixing these up gives misleading output.

## Debugging walkthrough

A web service is slower after a deployment. Collect a 30-second CPU profile:

```bash
curl -o cpu.pprof http://localhost:6060/debug/pprof/profile?seconds=30
go tool pprof cpu.pprof
```

```
(pprof) top
Showing nodes accounting for 45.6s, 87.5% of 52.1s total
      flat  flat%   sum%        cum   cum%
    18.2s 34.9% 34.9%     18.2s 34.9%  encoding/json.(*Decoder).reflectValue
    12.1s 23.2% 58.1%     12.1s 23.2%  runtime.mallocgc
     8.3s 15.9% 74.0%     30.5s 58.5%  main.handleRequest
```

`encoding/json.(*Decoder).reflectValue` has 34.9% flat CPU. This is JSON decoding using reflection. The 12.1s in `mallocgc` indicates heavy allocation. The fix: use `json.Decoder` with `DisallowUnknownFields` or switch to a faster JSON library (`github.com/goccy/go-json`, `github.com/valyala/fastjson`). After the fix, collect another profile and compare: the `encoding/json` node should drop below 5%.

## Production notes

- Collect CPU profiles during peak traffic, not during idle windows. A 30-second profile under load is statistically significant.
- In Kubernetes, use an init container or sidecar to collect profiles on-demand rather than running the profiler continuously.
- CPU profiling is safe in production — the 5% overhead is acceptable for short investigation windows (30-60 seconds). Do not leave CPU profiling enabled 24/7 without understanding the capacity impact.
- Store profiles with metadata: commit hash, timestamp, and a brief description of the workload. This enables historical comparisons.

## Performance implications

The signal-based sampling mechanism adds deterministic overhead: each SIGPROF signal triggers a context switch to the signal handler, stack unwinding, and hash table insertion. At 100 Hz, this is approximately 5% CPU overhead. The overhead is proportional to the sample rate: at 1000 Hz, overhead approaches 30%. Go's default of 100 Hz balances statistical significance with overhead. For most investigations, 30 seconds at 100 Hz (3,000 samples) provides sufficient accuracy.

## Practice task

Write a Go program that:

1. Defines a function `slowFunction(n int) int` that computes the sum of all primes below `n` by trial division.
2. Defines a function `fastFunction(n int) int` that computes the same sum using a Sieve of Eratosthenes.
3. Calls both functions with `n = 50000` while a CPU profile is active.
4. Saves the profile to `compare.pprof`.
5. Runs `go tool pprof -top compare.pprof` and notes which function appears more in the profile.

## Tests / verification

```bash
go run ./curriculum/modules/13-performance-memory-engineering/lessons/02-cpu-profiling
go test ./curriculum/modules/13-performance-memory-engineering/lessons/02-cpu-profiling
```

## Review questions

1. What is the difference between flat and cumulative CPU time in `go tool pprof` output?
2. Why is 100 Hz the default CPU profiling rate? What tradeoff does it represent?
3. How would you collect a CPU profile from a running HTTP service without modifying its source code?
4. If a function has 2% flat time but 40% cumulative time, what does that indicate about the function's role?
5. What command generates an SVG call graph from a CPU profile?

## NEXT UP

Memory profiling — understanding heap profiles, distinguishing alloc_objects from alloc_space and inuse_objects from inuse_space, and detecting memory leaks through heap dump analysis.
