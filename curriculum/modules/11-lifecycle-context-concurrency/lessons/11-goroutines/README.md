# Goroutines

## Learning objective

Spawn concurrent work with the `go` keyword, reason about goroutine scheduling and lifecycle, and control parallelism with `GOMAXPROCS`.

## Why this matters

Every Go program that handles multiple clients, processes streams, or performs background work uses goroutines. An HTTP server creates one goroutine per incoming request. A CLI tool that fetches data from three APIs concurrently cuts wall-clock time from the sum of latencies to the maximum latency. Goroutines are Go's answer to the question: "How do I do many things at once without the overhead of OS threads?" Mastering goroutines is prerequisite to every other concurrency topic.

## Mental model

A goroutine is a independently executing function that shares the same address space as its creator. Think of your program as a kitchen: the `main` function is the head chef following a recipe. When the head chef says `go peelPotatoes()`, a line cook peels potatoes in parallel while the head chef continues to the next step. Both chefs share the same kitchen (memory) and must coordinate access to shared tools (data structures).

Goroutines are multiplexed onto OS threads automatically. The Go scheduler decides which goroutine runs on which thread at any moment. From your perspective, every goroutine runs concurrently; from the hardware perspective, `GOMAXPROCS` goroutines run in parallel.

## Core idea

The `go` keyword before a function call creates a new goroutine:

```go
go fn()
go func() { /* closure */ }()
```

The caller returns immediately. The new goroutine executes `fn` in its own stack. The function's return value (if any) is discarded — you cannot `result := go fn()`. Communication happens through channels or shared memory with synchronization.

Goroutines are not OS threads:

| Property | OS thread | Goroutine |
|---|---|---|
| Initial stack | ~1 MB | ~2 KB (grows as needed) |
| Create/teardown | Expensive (syscall) | Cheap (userspace) |
| Scheduling | Kernel | Go runtime (M:N scheduler) |
| Context switch | ~1 µs | ~100 ns |
| Max on typical server | Thousands | Millions |

## Under the hood

A goroutine is a `g` struct in the Go runtime. It holds a small stack (initially 2 KB), an instruction pointer, and metadata. When you write `go fn()`, the compiler inserts a call to `runtime.newproc`. This:

1. Allocates a `g` struct on the heap.
2. Copies the function's arguments from the caller's stack to the new goroutine's stack.
3. Enqueues the `g` onto the current `P`'s local run queue.

The Go scheduler runs in user space. It implements an M:N model: `M` OS threads (machines) execute `P` logical processors, which schedule `G` goroutines. Each `P` has a local run queue. When a goroutine makes a blocking call (I/O, channel send/receive, `time.Sleep`, `sync.Mutex`), the `M` is released to pick up another goroutine from the `P`'s queue or from other `P`s (work stealing).

Since Go 1.14, the scheduler is **asynchronously preemptive**: a goroutine that runs too long without a function call (e.g., a tight loop) is preempted by a signal (`SIGURG`), forcing a scheduler decision. Before async preemption, a tight loop could starve other goroutines indefinitely.

## How Go uses it

- **HTTP servers**: `net/http` creates a goroutine per accepted connection via `go c.serve()`.
- **Concurrent API calls**: fan-out requests to multiple services with `go` and collect results.
- **Background workers**: periodic tasks like log rotation, cache warming, or metric aggregation.
- **`go test`**: subtests can run in parallel with `t.Parallel()`, which uses goroutines under the hood.
- **Finalizers and GC**: the runtime uses internal goroutines for `runtime.GC()` and `runtime/trace`.

## Go example

```go
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Duration(id) * 50 * time.Millisecond)
	fmt.Printf("Worker %d done\n", id)
}

func main() {
	runtime.GOMAXPROCS(2)
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()
	fmt.Println("All workers finished")
}
```

Output (order varies due to scheduling):

```
Worker 5 starting
Worker 1 starting
Worker 3 starting
Worker 2 starting
Worker 4 starting
Worker 1 done
Worker 2 done
Worker 3 done
Worker 4 done
Worker 5 done
All workers finished
```

## Step-by-step execution

For the program above with `GOMAXPROCS=2` and 5 workers:

1. `main` goroutine starts. Sets `GOMAXPROCS=2` — at most 2 OS threads run Go code simultaneously.
2. Loop `i=1..5`: each `go worker(i, &wg)` creates a new goroutine. The `wg.Add(1)` increments the counter before the goroutine starts.
3. The scheduler places goroutines onto `P0` and `P1` local queues. With 5 goroutines and 2 Ps, 2 run immediately and 3 queue.
4. Each worker prints its start message, sleeps (`time.Sleep` yields the processor), then prints its done message.
5. When a worker sleeps, the scheduler dequeues a waiting worker onto the freed P.
6. Each worker calls `wg.Done()` before exiting, decrementing the counter.
7. `main` blocks at `wg.Wait()` until the counter reaches 0.
8. After the 5th worker finishes, the counter is 0, `wg.Wait()` unblocks, and `main` prints the final message.

The order of start/done messages is non-deterministic because the scheduler interleaves goroutines.

## Common mistakes

- **No synchronization**: starting a goroutine and exiting `main` before it runs. The program exits when `main` returns; pending goroutines are terminated.
  - Fix: use `sync.WaitGroup` or channels to coordinate.
- **`GOMAXPROCS` misunderstanding**: `GOMAXPROCS` limits parallel execution of Go code, not total goroutines. 100 goroutines with `GOMAXPROCS=1` still run concurrently (interleaved on one thread).
- **Assuming goroutines start instantly**: the `go` statement enqueues the goroutine; it may not run until the next scheduler point.
- **Starving other goroutines**: a tight loop without function calls blocks the P. If `GOMAXPROCS=1`, no other goroutine runs until the loop yields.
  - Fix: insert `runtime.Gosched()` in tight loops to voluntarily yield.
- **Capturing loop variables**: `for i := 0; i < 5; i++ { go func() { fmt.Println(i) }() }` — all goroutines see the final `i` value.
  - Fix: `go func(n int) { fmt.Println(n) }(i)`.

## Debugging walkthrough

Consider this program that exits prematurely:

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		fmt.Println("background work")
	}()
	fmt.Println("main done")
}
```

**Symptom**: sometimes "background work" does not print.

**Investigation**: add a `time.Sleep` to confirm the goroutine runs. Then check whether the issue is timing or never started.

**Root cause**: `main` returns before the goroutine is scheduled. The program exits.

**Fix 1 — WaitGroup**:

```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
	defer wg.Done()
	fmt.Println("background work")
}()
wg.Wait()
```

**Fix 2 — channel**:

```go
done := make(chan struct{})
go func() {
	fmt.Println("background work")
	close(done)
}()
<-done
```

## Production notes

- **Goroutine panics**: an unhandled panic in a goroutine crashes the entire process, not just the goroutine. Defer `recover` inside every long-lived goroutine to isolate failures.
- **Connection handling**: every `net/http` connection creates a goroutine. Under load, this means tens of thousands of goroutines. The runtime handles this well, but each goroutine consumes ~2 KB stack + heap overhead for its closure.
- **Graceful shutdown**: use `context.Context` + `sync.WaitGroup` to drain in-flight goroutines before `os.Signal`.
- **Logging**: always include a goroutine identifier (or use `runtime.GoID()` via `runtime.Stack`) in logs during debugging, but never use `runtime.GoID()` in production — it is considered an anti-pattern.
- **Limit creation**: for CPU-bound work, cap active goroutines to `runtime.GOMAXPROCS(0)` to avoid oversubscription.

## Performance implications

| Operation | Approximate cost |
|---|---|
| Goroutine creation | ~200 ns (on modern CPU) |
| Goroutine context switch | ~100 ns (userspace) |
| OS thread context switch | ~1 µs (kernel) |
| Goroutine stack (idle) | ~2 KB |
| OS thread stack (idle) | ~1 MB |

Creating a goroutine is ~5000x cheaper than creating an OS thread. However, goroutines are not free. 1,000,000 goroutines consuming 2 KB each = 2 GB of stack memory before any work. Use bounded worker pools for massive concurrency.

`GOMAXPROCS` greater than `runtime.NumCPU()` can degrade performance due to cache contention and thread switching. The default is `NumCPU`.

## Practice task

Write a program that:

1. Launches 10 goroutines, each printing its ID and sleeping for `ID * 100` milliseconds.
2. Uses `sync.WaitGroup` to wait for all goroutines to finish before `main` exits.
3. Reports the total elapsed time.
4. Prints the current `GOMAXPROCS` value and the number of CPU cores.

Run it and observe the non-deterministic print order. Then set `GOMAXPROCS=1` and compare the output sequence and total time.

## Tests / verification

```go
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/11-goroutines
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/11-goroutines
```

## Review questions

1. What is the initial stack size of a goroutine, and how does it compare to an OS thread?
2. What does `runtime.GOMAXPROCS(0)` return? What happens if you set it to 1 on an 8-core machine?
3. Why can't you capture the return value of `go fn()`? How do you get results from a goroutine?
4. What happens when `main` returns while 3 goroutines are still running?
5. Describe a scenario where creating too many goroutines can harm performance despite each one being "cheap".

## NEXT UP

Waitgroups — coordinate the completion of multiple goroutines with `sync.WaitGroup`.
