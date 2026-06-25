# Race detector preview

## Learning objective

Understand Go's race detector, interpret its output, add the `-race` flag to tests and builds, and identify common data race patterns.

## Why this matters

Data races are the most insidious bugs in concurrent Go programs. They cause intermittent crashes, corrupted data, and behaviour that changes between runs, architectures, and load levels. A race may pass a million tests and fail in production on the million-and-first request. Go's built-in race detector is the primary defence: it finds races deterministically by instrumenting every memory access, and it integrates with zero code changes — just add `-race`.

## Mental model

A data race occurs when two goroutines access the same memory location concurrently, and at least one access is a write. The race detector is like a traffic camera at every memory intersection: it watches every read and write, and if it sees two accesses to the same address without a synchronisation event (mutex, channel, atomic) between them, it flags a violation.

```
Goroutine 1          Goroutine 2
   read x  ─────────  write x   ← RACE! (no synchronisation)
```

## Core idea

Enable the race detector with the `-race` flag:

```bash
go test -race ./...
go build -race ./...
go run -race main.go
```

The race detector:
- Instruments **every** memory read and write in the program.
- Tracks goroutine creation and synchronisation events (mutex Lock/Unlock, channel send/receive, `sync.WaitGroup`, `sync/atomic`).
- Reports the **exact** line numbers of both racing accesses, including the goroutine stacks at the time of the access.
- Adds ~5–10x CPU and memory overhead. Never use `-race` in production — only in testing and development.

Race detector output format:

```
WARNING: DATA RACE
Read at 0x... by goroutine 7:
  main.readCounter()
      /path/main.go:25 +0x...

Previous write at 0x... by goroutine 6:
  main.incrementCounter()
      /path/main.go:18 +0x...

Goroutine 7 (running) created at:
  main.main()
      /path/main.go:12 +0x...

Goroutine 6 (finished) created at:
  main.main()
      /path/main.go:11 +0x...
```

## Under the hood

Go's race detector is based on **ThreadSanitizer (TSan)**, a LLVM-based dynamic analysis tool. When `-race` is passed, the Go compiler inserts calls to TSan's runtime library before every memory access. TSan maintains a shadow memory map that tracks, for each memory location, which goroutine last accessed it and what kind of access (read/write). On each access, TSan checks the shadow map: if another goroutine has a conflicting access without a happens-before edge (i.e., no synchronisation), TSan reports a race.

The happens-before relation is tracked through synchronisation primitives:
- `sync.Mutex.Lock` → `Unlock`
- `ch <- v` → `<-ch`
- `sync.WaitGroup.Add` → `Wait`
- `sync/atomic.*` operations
- `time.Sleep` does NOT create a happens-before edge

## How Go uses it

- **CI**: Always run `go test -race ./...` on critical packages. This is the standard practice in Go projects of any size.
- **`-race` with `go build`**: Build a race-enabled binary for stress testing in staging. Never deploy to production.
- **`-race` with fuzz testing**: `go test -fuzz=Fuzz -race` catches races with generated inputs.
- **`-race` with benchmarks**: `go test -race -bench=. -benchtime=1x` checks for races in benchmarked code.
- **`sync.Map` vs `map + sync.Mutex`**: The race detector catches unsynchronised map access regardless of which data structure you use.

## Go example

```go
package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu    sync.Mutex
	value int
}

func main() {
	var wg sync.WaitGroup
	c := &Counter{}

	// Unsynchronised write — RACE!
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.value++ // race: concurrent writes to value
		}()
	}

	wg.Wait()
	fmt.Println("Counter:", c.value)
}
```

Run with:
```
$ go run -race main.go
==================
WARNING: DATA RACE
Write at 0x... by goroutine 7:
  main.main.func1()
      main.go:16 +0x...

Previous write at 0x... by goroutine 6:
  main.main.func1()
      main.go:16 +0x...
==================
Counter: 9   // wrong! should be 10
```

Fix by adding mutex synchronisation:

```go
go func() {
    defer wg.Done()
    c.mu.Lock()
    c.value++
    c.mu.Unlock()
}()
```

## Step-by-step execution

Without the race detector:

1. 10 goroutines are spawned. Each reads `c.value`, increments, and writes back.
2. Goroutines run concurrently. Two goroutines may read `0`, both increment to `1`, and both write `1` — losing an increment.
3. The final value is non-deterministic (maybe 7, 9, or 10).
4. Without `-race`, this runs silently. The bug is invisible.

With `-race`:

1. The compiler instruments each `c.value++` to call TSan before the read and after the write.
2. Goroutine 6 reads `c.value`. TSan notes: "G6 read at address X".
3. Goroutine 7 also reads `c.value` without any synchronisation between G6 and G7. TSan checks the shadow map: "G6 accessed this, no happens-before edge between G6 and G7 → RACE".
4. TSan records both goroutine stacks and the exact lines.
5. At program exit (or when the race is detected), TSan prints the warning.

## Common mistakes

- **Assuming `-race` detects all races**: It only detects races that actually happen during execution. A race on a rarely-taken code path may not trigger. Run under realistic load.
- **Running `-race` in production**: The overhead (5–10x) makes it unsuitable for production. Some teams run a single canary instance with `-race` to catch races under live traffic.
- **Confusing race with deadlock**: The race detector does not report deadlocks. A deadlock-free program can still have races, and a race-free program can still deadlock.
- **Believing `time.Sleep` synchronises**: `time.Sleep` does not create a happens-before edge. `time.Sleep(time.Second)` does not guarantee a concurrent write has completed.
- **Ignoring the stack traces in the race report**: The report includes two goroutine stacks — one for each racing access. Read both. The "Previous" access is the one that established the conflicting state.

## Debugging walkthrough

You have a web server with a shared cache:

```go
var cache = make(map[string]string)

func handler(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Path
    if val, ok := cache[key]; ok {
        fmt.Fprint(w, val)
        return
    }
    val := fetchFromDB(key)
    cache[key] = val  // race: concurrent map writes
    fmt.Fprint(w, val)
}
```

Run `go test -race ./handlers/`. The race detector reports:

```
WARNING: DATA RACE
Write at 0x... by goroutine 12:
  main.handler()
      handlers/cache.go:15 +0x...

Previous write at 0x... by goroutine 8:
  main.handler()
      handlers/cache.go:15 +0x...
```

**Investigation**: Line 15 is `cache[key] = val`. Multiple goroutines (`handler` is called per request) write to the same map without synchronisation. Go maps are not safe for concurrent writes — the runtime panics on concurrent map access, but the race detector catches it first.

**Fix**: Add a `sync.RWMutex`:

```go
var (
    cache = make(map[string]string)
    mu    sync.RWMutex
)

func handler(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Path
    mu.RLock()
    val, ok := cache[key]
    mu.RUnlock()
    if ok {
        fmt.Fprint(w, val)
        return
    }
    val = fetchFromDB(key)
    mu.Lock()
    cache[key] = val
    mu.Unlock()
    fmt.Fprint(w, val)
}
```

## Production notes

- **CI integration**: Add `go test -race ./...` to every CI pipeline. Set `GORACE="halt_on_error=1"` to make the test binary exit with a non-zero code on the first race.
- **Cost**: The race detector adds ~5–10x CPU and memory. A test suite that takes 10 seconds without `-race` may take 60 seconds with it. Run `-race` on a subset of packages on every commit and the full suite nightly.
- **False positives**: Rare but possible. TSan may report a race on benign concurrent access (e.g. two goroutines reading the same package-level variable that is never written after init). Use `sync/atomic` to suppress false positives.
- **`GORACE` environment variable**: Controls TSan behaviour: `log_path`, `exitcode`, `halt_on_error`, `history_size`. See `go doc runtime/race`.

## Performance implications

- CPU overhead: 5–10x slowdown. Memory overhead: 5–10x increase (shadow memory).
- Binary size: The instrumented binary is ~2x larger.
- Runtime: Even without races, the instrumented code runs slower because every memory access calls TSan.
- Never use `-race` in production. Build a separate staging binary with `-race` for load testing.

## Practice task

Fix the data race in the provided `Counter` program. First, run `go run -race main.go` to confirm the race. Then add a `sync.Mutex` to protect `c.value++`. After the fix, run `go run -race main.go` again — the race warning should disappear and the counter should always print `10`.

## Tests / verification

```bash
go run -race ./curriculum/modules/06-testing-debugging-refactoring/lessons/21-race-detector-preview
go test -race ./curriculum/modules/06-testing-debugging-refactoring/lessons/21-race-detector-preview
```

## Review questions

1. What does a data race mean in terms of goroutine memory access?
2. How do you enable the race detector in Go?
3. Why should you not run the race detector in production?
4. What is a "happens-before" relation, and which Go operations create one?
5. What does the race detector output tell you about each racing access?

## NEXT UP

Congratulations on completing Module 06! You now understand Go's testing philosophy, debugging workflows, and safe refactoring techniques. Next up: Module 07 — CLI, Files, JSON, and Configuration.
