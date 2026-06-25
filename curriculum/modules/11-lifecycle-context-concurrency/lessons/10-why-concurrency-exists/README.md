# Why concurrency exists

## Learning objective

Distinguish concurrency from parallelism, explain Amdahl's law and its impact on scaling, classify workloads as CPU-bound or IO-bound, and describe Go's concurrency model based on goroutines and channels.

## Why this matters

Concurrency is not about making programs faster — it is about making programs structured for the world they run in. Modern systems wait: for network responses, disk I/O, database queries, and user input. Sequential code wastes CPU cycles waiting. Concurrency lets one part of the program make progress while another part waits. Go was designed from the ground up for this world, and understanding why concurrency exists shapes every design decision — from how you structure a handler to how you architect a microservice.

## Mental model

Concurrency is the composition of independently executing tasks. Parallelism is the simultaneous execution of multiple tasks. Think of concurrency as the structure of your program (how you _write_ the code) and parallelism as the execution (how the hardware _runs_ the code). Go gives you concurrency primitives (goroutines, channels) that can be parallelized by the runtime if multiple CPU cores are available — but the code is the same either way.

## Core idea

**Concurrency vs parallelism**: Concurrency is about dealing with many things at once (design). Parallelism is about doing many things at once (execution). A concurrent program may or may not run in parallel.

**Amdahl's law**: The speedup from parallelization is limited by the sequential portion of the workload. If 10% of a task must be sequential, the maximum speedup is 10x, no matter how many CPUs you add.

**CPU-bound vs IO-bound**: A CPU-bound task (e.g., video encoding) is limited by processor speed. Adding more goroutines does not help — it adds context-switch overhead. An IO-bound task (e.g., HTTP requests) is limited by waiting. Concurrency helps by allowing other goroutines to run while one waits.

**Go's model**: Goroutines are lightweight (starting stack ~4KB) and multiplexed onto OS threads by the Go scheduler (M:N scheduling). Channels provide safe communication between goroutines. The slogan: "Do not communicate by sharing memory; instead, share memory by communicating."

## Under the hood

Go's runtime implements M:N scheduling: M goroutines are multiplexed onto N OS threads. The scheduler (G-M-P model) has three components:
- **G** (goroutine): A lightweight thread with its own stack, registers, and state.
- **M** (machine): An OS thread.
- **P** (processor): A logical processor that holds a queue of runnable goroutines.

When a goroutine blocks on I/O or a channel operation, the scheduler detaches the M from the P, starts a new M (or wakes a sleeping one) to keep the P busy, and the blocked goroutine eventually resumes on another M. This allows Go to handle hundreds of thousands of concurrent goroutines efficiently.

## How Go uses it

- **HTTP servers**: Each incoming request is typically handled by its own goroutine. The server handles thousands of concurrent connections with a small number of OS threads.
- **Database connection pools**: Queries run in goroutines. When a query blocks waiting for a connection, the goroutine yields, and another goroutine runs.
- **Pipeline pattern**: Data flows through stages connected by channels. Each stage runs in its own goroutine. This is the idiomatic Go way to compose concurrent operations.
- **Worker pools**: A fixed number of goroutines process work from a shared channel. This bounds concurrency and prevents resource exhaustion.

## Go example

```go
package main

import (
	"fmt"
	"time"
)

type task func(int) int

func sequential(tasks []task, input int) int {
	result := input
	for _, t := range tasks {
		result = t(result)
	}
	return result
}

func slowSquare(n int) int {
	time.Sleep(10 * time.Millisecond)
	return n * n
}

func slowDouble(n int) int {
	time.Sleep(10 * time.Millisecond)
	return n * 2
}

func main() {
	tasks := []task{slowSquare, slowDouble, slowSquare}

	start := time.Now()
	result := sequential(tasks, 3)
	elapsed := time.Since(start)
	fmt.Printf("Sequential: %d (took %v)\n", result, elapsed)
}
```

## Step-by-step execution

1. `tasks` is a slice of three functions, each simulating 10ms of work.
2. `sequential` calls them one at a time: `slowSquare(3) = 9` (10ms), `slowDouble(9) = 18` (10ms), `slowSquare(18) = 324` (10ms).
3. Total elapsed time is ~30ms (3 x 10ms).
4. If these tasks were IO-bound and independent, they could run concurrently: launch three goroutines, each running one task, and combine results in ~10ms.
5. Go's concurrency model makes this restructuring straightforward: replace the sequential loop with goroutines and channels.

The key insight: concurrency exists because the alternative (sequential execution of independent work) wastes time waiting. Even a single-CPU machine benefits from concurrency — while one goroutine waits for I/O, another makes progress.

## Common mistakes

- Mistake: Believing concurrency always improves performance.
  - Why: Concurrency adds overhead (goroutine creation, scheduling, synchronization). For CPU-bound tasks on a single core, a sequential loop is faster.
  - Fix: Profile first. Concurrency helps IO-bound or latency-sensitive workloads. For CPU-bound work, consider parallelism with `runtime.GOMAXPROCS`.

- Mistake: Confusing concurrency with parallelism.
  - Why: A concurrent program may still run on a single OS thread. True parallelism requires multiple CPU cores and `GOMAXPROCS > 1`.
  - Fix: Design for concurrency (structure). The runtime handles parallelism (execution).

- Mistake: Spawning unbounded goroutines.
  - Why: Goroutines are cheap, but not free. Each goroutine consumes stack space (~4KB minimum) and scheduler overhead.
  - Fix: Use bounded worker pools or errgroup with a semaphore to limit concurrency.

- Mistake: Ignoring Amdahl's law when designing concurrent systems.
  - Why: If one bottleneck serializes all work (e.g., a shared mutex, a single database write connection), adding goroutines does not improve throughput.
  - Fix: Identify and minimize serial portions. Use sharding, partitioning, or lock-free data structures.

## Debugging walkthrough

Consider a concurrent web scraper that is slower than the sequential version:

```go
func scrape(urls []string) []string {
    results := make([]string, len(urls))
    var wg sync.WaitGroup
    for i, url := range urls {
        wg.Add(1)
        go func(i int, url string) {
            defer wg.Done()
            resp, _ := http.Get(url)
            body, _ := io.ReadAll(resp.Body)
            resp.Body.Close()
            results[i] = string(body)
        }(i, url)
    }
    wg.Wait()
    return results
}
```

**Symptom**: Scraping 100 URLs concurrently is slower than scraping them sequentially.

**Investigation**: Add timing per request:

```go
go func(i int, url string) {
    start := time.Now()
    resp, err := http.Get(url)
    log.Printf("url=%s took=%v err=%v", url, time.Since(start), err)
    // ...
}(i, url)
```

The logs show that most requests spend 95% of their time waiting for the `http.Get` to establish a connection. The system's file descriptor limit (ulimit) caps the number of concurrent connections.

**Root cause**: The default `http.Transport` has a limit on the number of concurrent connections per host. When 100 goroutines try to connect to the same host simultaneously, most wait in the connection pool queue. The sequential version completes each request without queueing.

**Fix**: Increase the transport's `MaxIdleConnsPerHost` and `MaxConnsPerHost`, or throttle the goroutines with a semaphore:

```go
sem := make(chan struct{}, 10) // max 10 concurrent
go func(i int, url string) {
    sem <- struct{}{}
    defer func() { <-sem }()
    // ... http.Get
}(i, url)
```

## Production notes

- Use `GOMAXPROCS` to control CPU parallelism. The runtime sets it to the number of cores by default.
- For IO-bound workloads, concurrency improves throughput even on a single core. For CPU-bound workloads, match concurrency to available cores.
- Monitor goroutine count in production. A growing goroutine count indicates a leak or unbounded concurrency.
- Use `errgroup` for orchestrating concurrent operations where any error should cancel all others.
- Respect rate limits and connection pool limits when designing concurrent systems — unbounded concurrency degrades performance.

## Performance implications

- Goroutine creation is fast (~1KB stack + ~500ns startup).
- Channel sends/receives are ~50ns (unbuffered) when uncontended.
- The Go scheduler makes ~10ns decisions to switch goroutines.
- Context switching between goroutines is far cheaper than OS thread context switching (~1us vs ~10us).
- For CPU-bound work, match goroutine count to `runtime.GOMAXPROCS` to avoid oversubscription.
- For IO-bound work, goroutine count should match the number of concurrent I/O operations, not CPU cores.

## Practice task

Revise the `sequential` example to run the three tasks concurrently using goroutines and channels. Create a function `concurrent(tasks []task, input int) int` that runs each task in its own goroutine, chaining results through channels. Measure and print the elapsed time for both sequential and concurrent execution. The concurrent version should be ~3x faster than sequential because the tasks simulate IO-bound work.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/10-why-concurrency-exists
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/10-why-concurrency-exists
```

The existing tests verify the sequential function's correctness and timing characteristics. After completing the practice task, add tests for `concurrent` that verify the result is correct and the execution time is less than the sum of individual task times.

## Review questions

1. What is the difference between concurrency and parallelism?
2. According to Amdahl's law, if 5% of a workload is sequential, what is the maximum speedup from parallelization?
3. Why are goroutines lighter than OS threads?
4. When would adding more goroutines hurt, not help, performance?
5. What is the M:N scheduling model used by Go's runtime?

## NEXT UP

Goroutines — spawning lightweight concurrent threads of execution with the `go` keyword.
