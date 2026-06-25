# Bounded worker pools

## Learning objective

Implement a bounded worker pool using a buffered channel as a semaphore, dispatch jobs to a fixed number of workers, and collect results.

## Why this matters

Starting an unbounded number of goroutines for each incoming task crashes a service under load. Each goroutine consumes stack memory, scheduler overhead, and potentially file descriptors or network connections. A bounded worker pool caps resource usage: at most N goroutines are active at any time, and excess tasks queue up in a channel. This is the fundamental pattern for controlling concurrency in production Go services.

## Mental model

A worker pool is like a checkout counter at a supermarket. There are N cashiers (workers) who process customers (jobs). If all cashiers are busy, new customers line up in a queue (buffered channel). When a cashier finishes, they call the next customer from the queue. The number of cashiers stays constant, so the store never has more than N customers being served at once, regardless of how many customers enter.

## Core idea

The bounded worker pool pattern uses three components:

1. **Job channel**: a buffered channel that holds pending jobs. The buffer size limits how many jobs can queue.
2. **Worker goroutines**: a fixed number of goroutines that read from the job channel, process each job, and write results to the result channel.
3. **Result channel**: a channel for collecting outputs (optional if workers have side effects).

The key insight: a buffered channel with capacity N acts as a semaphore. Sending to the channel acquires a permit; receiving releases it. The number of workers equals the number of goroutines, not the channel capacity.

## Under the hood

Each worker runs a `for` loop that reads from the jobs channel using a `select` with `ctx.Done()` for cancellation. When the jobs channel is closed, the workers receive the zero value with `ok == false` and exit. A separate goroutine closes the result channel after all workers have exited (using `sync.WaitGroup`). The Go scheduler multiplexes the N worker goroutines across available OS threads. If workers block on I/O, the scheduler automatically runs other workers.

## How Go uses it

The `net/http` server uses a bounded goroutine pool (the `Server` has a `MaxConcurrentRequests` limit). The `database/sql` connection pool is a bounded resource pool where each connection is obtained via a channel-based semaphore. The `go test` runner uses a bounded pool to run tests in parallel within a package. Many web frameworks use worker pools to handle request processing with a fixed goroutine limit.

## Go example

```go
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID    int
	Delay time.Duration
}

type Result struct {
	JobID  int
	Output string
}

func worker(ctx context.Context, id int, jobs <-chan Job, results chan<- Result) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			time.Sleep(job.Delay)
			results <- Result{
				JobID:  job.ID,
				Output: fmt.Sprintf("worker %d processed job %d", id, job.ID),
			}
		}
	}
}

func main() {
	numJobs := 10
	numWorkers := 3

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(ctx, id, jobs, results)
		}(w)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- Job{ID: j, Delay: time.Duration(j%3) * 10 * time.Millisecond}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Println(r.Output)
	}
}
```

## Step-by-step execution

1. `main` creates buffered channels for jobs (10) and results (10).
2. `main` starts 3 worker goroutines. Each blocks on `<-jobs`.
3. `main` sends 10 jobs to `jobs`. Each worker picks up a job as it becomes ready.
4. `main` closes `jobs`. Workers will exit after processing their current job when they next read from `jobs` and see `ok == false`.
5. A separate goroutine waits for all workers (via `wg.Wait()`), then closes `results`.
6. `main` ranges over `results`, printing each result. After `results` is closed, the loop ends.
7. At most 3 workers are active at any time. The goroutine count is bounded.

## Common mistakes

- Mistake: Using a buffered channel as the worker pool limit but starting unbounded goroutines.
  - Why it happens: The channel buffer limits queued jobs, not goroutines. You must start exactly N goroutines.
  - Fix: The limit is the number of worker goroutines you start, not the channel buffer size.

- Mistake: Closing the results channel before all workers have finished sending.
  - Why it happens: If `close(results)` executes while a worker is still sending, it panics.
  - Fix: Close results only after all workers have exited (use `sync.WaitGroup`).

- Mistake: Not cancelling workers on shutdown.
  - Why it happens: Workers wait on `<-jobs` with no way to unblock if the system is shutting down.
  - Fix: Pass a `context.Context` and `select` on `ctx.Done()` as shown above.

- Mistake: Creating a new channel for each job instead of using a shared job channel.
  - Why it happens: The developer sends jobs by creating a channel per worker.
  - Fix: Use one shared buffered `chan Job` that all workers read from. This load-balances automatically.

## Debugging walkthrough

Buggy program:

```go
jobs := make(chan Job, 10)
for w := 1; w <= 3; w++ {
	go worker(context.Background(), w, jobs, nil)
}
for j := 1; j <= 10; j++ {
	jobs <- Job{ID: j}
}
// no close(jobs) -- workers wait forever
```

Symptom: program finishes, but 3 goroutines are leaked (blocked on `<-jobs`).

Investigation: `runtime.NumGoroutine()` returns 5 (1 main + 1 for the runtime + 3 leaked workers).

Fix: `close(jobs)` after sending all jobs. Workers will exit when they see `ok == false`.

## Production notes

Choose the worker count based on the bottleneck: CPU-bound workloads use `runtime.GOMAXPROCS(0)` workers. I/O-bound workloads use a higher count, but never more than the downstream system can handle. Set the job channel buffer size to absorb bursts: a common heuristic is `numWorkers * 2`. Always add a backlog limit: if the job channel is full, either block (applying backpressure) or reject the job. Monitor the channel depth as a Prometheus gauge to detect when workers are falling behind.

## Performance implications

Worker pools reduce goroutine contention by limiting the number of active goroutines. With N workers, at most N goroutines contend for scheduler time. The job channel send/receive adds ~50 ns per operation. The result channel adds another ~50 ns. The trade-off is between goroutine overhead (stack + scheduling) and channel synchronization cost. For workloads where each job takes >1 ms, the channel overhead is negligible. For sub-microsecond tasks, a single goroutine with batching outperforms a pool.

## Practice task

Implement a `Pool` struct with `Start(ctx, numWorkers int, jobs <-chan Job) <-chan Result` method. Internally, it spawns `numWorkers` goroutines, reads from `jobs`, processes, and sends to `result` channel. When `jobs` is closed and all workers finish, it closes `result`. Write a table-driven test that verifies: correct number of results, workers do not exceed the bound, and context cancellation stops all workers.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/26-bounded-worker-pools
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/26-bounded-worker-pools
```

## Review questions

1. What is the difference between a buffered channel used as a semaphore versus a buffered channel used as a job queue?
2. How do you ensure all workers exit cleanly when there are no more jobs?
3. Why should `close(results)` happen after all workers have finished, not immediately after `close(jobs)`?
4. How do you choose the number of workers for a CPU-bound workload vs. an I/O-bound workload?
5. What happens if the job channel is full and a producer tries to send?

## NEXT UP

Retries and backoff
