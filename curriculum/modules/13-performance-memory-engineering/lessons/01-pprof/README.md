# pprof

## Learning objective

Capture CPU and heap profiles programmatically using `runtime/pprof`, serve them via `net/http/pprof`, and inspect them with `go tool pprof` to identify performance hotspots and memory patterns in Go programs.

## Why this matters

Every Go engineer will eventually debug a slow or memory-hungry service. Without profiling, you optimize based on intuition rather than evidence — and intuition is wrong more often than not. pprof is the universal language of Go performance investigation. Engineers at every major Go shop (Uber, Cloudflare, Google, Datadog) start every performance investigation by collecting a profile. Learning pprof is the first step toward data-driven optimization.

## Mental model

pprof is a statistical sampling profiler built into the Go runtime. Think of it as an x-ray machine taking snapshots at regular intervals. Each snapshot captures either "what line of code is the CPU executing right now?" (CPU profile) or "what code path allocated this memory?" (heap profile). After enough snapshots (thousands per second), the distribution of stack traces accurately reflects where time or memory is going. It is statistical, not precise — like polling 1,000 voters and extrapolating to the entire population.

## Core idea

pprof collects stack trace samples at runtime and groups them by call path. The result is a profile: a set of call stacks, each with a count of how many times that stack was observed. `go tool pprof` visualizes these counts as text, graphs, or flame graphs. The two most commonly used profile types are:

| Profile | Source | Default rate | What it measures |
|---|---|---|---|
| CPU | `runtime/pprof.StartCPUProfile` | 100 Hz | Which functions consume CPU time |
| Heap | `runtime/pprof.WriteHeapProfile` | 1 sample / 512 KB allocated | Which call sites allocate memory |
| Goroutine | `pprof.Lookup("goroutine")` | All | Stack traces of all goroutines |
| Block | `pprof.Lookup("block")` | 1 per blocking event | Where goroutines block on sync primitives |
| Mutex | `pprof.Lookup("mutex")` | 1 per mutex unlock | Where mutex contention occurs |

## Under the hood

The `runtime/pprof` package hooks into the Go runtime's sampling infrastructure. For CPU profiling, the runtime installs a SIGPROF signal handler (on Unix) that fires at the configured frequency (default 100 Hz). On each signal, the handler reads the interrupted goroutine's program counter and stack pointer from the signal context, unwinds the stack using Go's metadata table (`gopclntab`), and hashes the resulting call stack into a fixed-size hash table with an incrementing count.

For heap profiling, the sampling happens at allocation time. Each call to `runtime.mallocgc` increments a per-thread allocation counter. When the counter exceeds `MemProfileRate / 2` (default 512 KB / 2 = 256 KB), the runtime generates a random number and decides whether to sample. On a sample, it captures up to 32 stack frames and stores the record with allocated bytes, freed bytes, and object counts.

The `net/http/pprof` package registers HTTP handlers on `DefaultServeMux` at paths like `/debug/pprof/`, `/debug/pprof/heap`, and `/debug/pprof/profile`. When a request arrives, it calls the corresponding `runtime/pprof` function and streams the protobuf-encoded profile to the HTTP response.

## How Go uses it

The Go standard library uses pprof internally for:

- The `go test -bench` framework measures allocation counts and bytes per operation using `runtime.MemStats` — the same memory counters that heap profiling samples.
- The `net/http/pprof` package is the standard way to expose profiles from production HTTP services. Most Go web frameworks integrate it with a single import line.
- The `runtime` package exposes `ReadMemStats` for programmatic memory inspection and `SetCPUProfileRate` for dynamic CPU profiling control.
- `go tool pprof` is a standalone CLI that reads profile protobuf files and supports interactive exploration (`top`, `list`, `web`, `peek`, `traces`).

## Go example

```go
package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func allocateMemory(size int) []byte {
	return make([]byte, size)
}

func main() {
	cpuFile, err := os.Create("cpu.pprof")
	if err != nil {
		panic(err)
	}
	defer cpuFile.Close()

	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	result := fibonacci(40)
	fmt.Printf("fibonacci(40) = %d\n", result)

	_ = allocateMemory(10 << 20)

	heapFile, err := os.Create("heap.pprof")
	if err != nil {
		panic(err)
	}
	defer heapFile.Close()

	runtime.GC()
	if err := pprof.WriteHeapProfile(heapFile); err != nil {
		panic(err)
	}

	fmt.Println("Profiles written to cpu.pprof and heap.pprof")
	time.Sleep(100 * time.Millisecond)
}
```

Run this program to generate `cpu.pprof` and `heap.pprof` files. Then inspect them:

```bash
go tool pprof -top cpu.pprof
go tool pprof -top -alloc_space heap.pprof
```

## Step-by-step execution

1. `os.Create("cpu.pprof")` opens a new file for writing the CPU profile.
2. `pprof.StartCPUProfile(cpuFile)` starts the sampling timer (SIGPROF at 100 Hz on Unix). The runtime begins recording stack traces on each signal.
3. `fibonacci(40)` executes for roughly 0.5–2 seconds (depending on CPU). During this time, the runtime collects ~50–200 stack samples.
4. `pprof.StopCPUProfile()` stops the timer, serializes all collected samples into protobuf format, and writes to `cpu.pprof`.
5. `allocateMemory(10 << 20)` allocates a 10 MB slice, triggering a heap sampling event.
6. `os.Create("heap.pprof")` opens a new file for the heap profile.
7. `runtime.GC()` forces a garbage collection so the heap profile reflects only live allocations.
8. `pprof.WriteHeapProfile(heapFile)` iterates all `MemProfileRecord` entries, computes `inuse = AllocBytes - FreeBytes` for each, and serializes the profile.

## Common mistakes

- **Forgetting to stop the CPU profile** — If the program exits without calling `pprof.StopCPUProfile()`, the profile file is empty or truncated. Always use `defer pprof.StopCPUProfile()` right after `StartCPUProfile`.
- **Profiling without realistic load** — A CPU profile of an idle process shows only GC and runtime overhead. Always profile under load that resembles production traffic.
- **Serving pprof on a public port** — The `net/http/pprof` handlers have no authentication. In production, serve them on a separate admin port (e.g., `:6060`) that is not exposed externally.
- **Reading alloc_space instead of inuse_space** — `alloc_space` shows cumulative allocations since process start. `inuse_space` shows currently live memory. For memory leak investigation, use `inuse_space`.

## Debugging walkthrough

Consider a service that is consuming 2 GB of RAM unexpectedly:

```bash
go tool pprof -inuse_space http://localhost:6060/debug/pprof/heap
```

Inside the pprof interactive shell:

```
(pprof) top10
Showing nodes accounting for 1.8GB, 90% of 2GB
      flat  flat%   sum%        cum   cum%
     1.2GB 60.00% 60.00%      1.2GB 60.00%  main.cacheEntries
     0.4GB 20.00% 80.00%      0.4GB 20.00%  main.parseResponse
     0.2GB 10.00% 90.00%      0.2GB 10.00%  main.newRequest
```

The `top10` output shows `main.cacheEntries` is responsible for 60% of live heap. Next:

```
(pprof) list cacheEntries
```

This shows the exact lines within `cacheEntries` that allocate. If the function caches entries without eviction, the fix is to add an LRU eviction policy or TTL.

## Production notes

- Import `net/http/pprof` with a blank identifier: `import _ "net/http/pprof"`. This registers debug handlers on the default mux. For security, use a separate HTTP server on an internal port.
- Profile overhead: CPU profiling adds ~5% CPU overhead at 100 Hz. Heap profiling adds ~1-3% allocation overhead. Enable profiling only during investigation windows, or use a continuous profiler (Datadog, Google Cloud Profiler) that respects rate limits.
- In Kubernetes, expose pprof on a separate container port and configure the liveness probe to collect a heap profile on OOM risk.
- Profile file sizes: a 30-second CPU profile is typically 50-200 KB. Heap profiles are 10-50 KB. These are safe to store and compare.

## Performance implications

Profiling itself has a cost. The signal handler for CPU profiling runs at 100 Hz and performs stack unwinding, hashing, and hash table insertion. This adds approximately 5% CPU overhead. Heap profiling samples 1 in every ~512 KB of allocation, adding roughly 1-3% overhead to allocation-heavy code. In latency-sensitive paths, consider disabling profiling during normal operation and enabling it only when investigating a specific issue.

## Practice task

Write a Go program that:

1. Starts a CPU profile and writes it to `practice.pprof`.
2. Computes `fibonacci(35)` 10 times in a loop.
3. Allocates 100 slices of 1 MB each, keeping references alive in a slice.
4. Stops the CPU profile.
5. Writes a heap profile to `practice_heap.pprof`.
6. Prints the file sizes of both profile files.

Then run `go tool pprof -top practice.pprof` and `go tool pprof -top -inuse_space practice_heap.pprof` and note the top entry in each.

## Tests / verification

```bash
go run ./curriculum/modules/13-performance-memory-engineering/lessons/01-pprof
go test ./curriculum/modules/13-performance-memory-engineering/lessons/01-pprof
```

## Review questions

1. What is the default CPU profiling sample rate in Go? How would you change it?
2. Why does heap profiling use sampling instead of recording every allocation?
3. What is the difference between `-inuse_space` and `-alloc_space` in `go tool pprof`?
4. How would you expose pprof endpoints on a production HTTP service safely?
5. What happens if you forget to call `pprof.StopCPUProfile()` before the program exits?

## NEXT UP

CPU profiling — how to interpret CPU profiles, identify hotspots, use flame graphs, and apply the `top`/`list`/`web` commands to drill into performance bottlenecks.
