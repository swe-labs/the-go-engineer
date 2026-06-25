# Queues

## Learning objective

Implement an in-memory FIFO queue, build a producer-consumer pattern using buffered Go channels, and deploy a concurrent worker pool that drains work from a shared queue.

## Why this matters

Queues are the backbone of asynchronous processing in every production Go service. HTTP request handlers must respond in milliseconds, but many operations -- sending email, resizing images, processing payments -- take seconds. Without a queue, those slow operations block the response, causing timeouts and poor user experience. With a queue, the handler enqueues work in microseconds and returns immediately. Background workers drain the queue and process work at their own pace. Every major Go deployment -- Kubernetes controllers, Kafka consumers, task queues -- relies on this pattern.

## Mental model

A queue is a line of people waiting for a service window. The first person in line is served first. New arrivals join at the back. This is First-In-First-Out (FIFO) order. In software, the "people" are messages or tasks, the "line" is a data structure (or a channel), and the "service windows" are concurrent worker goroutines. The queue decouples producers (who add items) from consumers (who remove and process them). Producers never wait for consumers, and consumers never wait for producers -- they meet only through the queue.

## Core idea

A FIFO queue exposes three fundamental operations:

- **Enqueue**: add an item to the tail.
- **Dequeue**: remove and return the item at the head.
- **Peek/Len**: inspect the head or the queue depth without modifying it.

Go provides two ways to build queues:

| Approach | Characteristics |
|---|---|
| Slice-backed (`[]T`) | Simple, resizable, needs mutex for concurrency |
| `chan T` | Built-in concurrent safe, blocking or buffered, language primitive |

A **buffered channel** IS a queue: it holds items in FIFO order up to its capacity. A **worker pool** is a set of goroutines that all read from the same channel, competing for work. Go's `chan` makes this pattern trivial compared to languages that require external libraries.

## Under the hood

A slice-backed queue stores items in a contiguous block of memory. Enqueue appends; dequeue re-slices `q.items[1:]`. The underlying array does not shrink, so repeated enqueue/dequeue cycles leak memory unless the slice is periodically re-allocated. A ring buffer (circular array) solves this but is more complex.

A `chan T` queue is implemented by the Go runtime as a circular buffer inside the hchan struct. When a goroutine sends on a buffered channel and the buffer is not full, the runtime copies the value into the buffer and wakes a waiting receiver (if any). When it receives and the buffer is not empty, the runtime copies from the buffer. This entire mechanism is lock-free in the fast path (no mutex per send/receive), making it extremely efficient.

## How Go uses it

Go's standard library and runtime use queues pervasively:

- **`net/http`**: the HTTP server uses a goroutine-per-connection model. Incoming connections are queued in the OS listen backlog and dispatched by the runtime's network poller.
- **`log`**: the standard logger writes synchronously, but many production log libraries use an internal queue and a dedicated writer goroutine to avoid blocking the caller.
- **`database/sql`**: the connection pool is a queue of `*sql.DB` connections. Acquire dequeues, release enqueues.
- **`time.After`**: internally uses a heap-based timer queue.
- **Kubernetes informers**: use a work queue (a priority queue) to deduplicate and order object reconciliation.

## Go example

```go
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Queue struct {
	mu    sync.Mutex
	items []int
}

func (q *Queue) Enqueue(item int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
}

func (q *Queue) Dequeue() (int, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return 0, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

type WorkerPool struct {
	queue     *Queue
	workers   int
	processed atomic.Int64
	wg        sync.WaitGroup
}

func NewWorkerPool(q *Queue, n int) *WorkerPool {
	return &WorkerPool{queue: q, workers: n}
}

func (wp *WorkerPool) Start(fn func(int)) {
	for range wp.workers {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for {
				item, ok := wp.queue.Dequeue()
				if !ok {
					return
				}
				fn(item)
				wp.processed.Add(1)
			}
		}()
	}
}

func (wp *WorkerPool) Wait()  { wp.wg.Wait() }
func (wp *WorkerPool) Count() int64 { return wp.processed.Load() }

func main() {
	q := &Queue{}
	for i := range 20 {
		q.Enqueue(i)
	}

	pool := NewWorkerPool(q, 4)
	pool.Start(func(n int) {
		fmt.Printf("worker consumed: %d\n", n)
	})
	pool.Wait()
	fmt.Printf("total processed: %d\n", pool.Count())
}
```

## Step-by-step execution

For the worker pool example with 20 items and 4 workers:

1. `main` enqueues integers 0 through 19 into the slice-backed queue.
2. `NewWorkerPool` creates a pool with 4 workers, each tracking processed via `atomic.Int64`.
3. `pool.Start` launches 4 goroutines. Each goroutine loops, calling `Dequeue`.
4. `Dequeue` acquires the mutex, checks `len(q.items) == 0`, reads `q.items[0]`, re-slices to `q.items[1:]`, releases the mutex.
5. Each goroutine calls the handler `fn(item)` then increments `processed`.
6. When all 20 items are dequeued, `Dequeue` returns `false` and each goroutine exits its loop.
7. `pool.Wait()` blocks until all 4 goroutines call `wg.Done()`.
8. `pool.Count()` returns 20, confirming every item was processed exactly once.

The key property: multiple goroutines compete for the same queue, but the mutex guarantees each item is dequeued by exactly one worker.

## Common mistakes

- **Unbounded queue growth**: if producers are faster than consumers, an in-memory queue grows without bound and consumes all RAM. Always bound queue size with a buffered channel capacity or implement back-pressure.
- **Nil channel blocks forever**: sending on a nil channel blocks permanently. Ensure channels are initialized before use.
- **Closing the channel while producers still send**: closing a channel causes panics on subsequent sends. Use a `sync.WaitGroup` or a coordination pattern to ensure all producers finish before close.
- **Reading from a closed channel**: after a channel is closed and drained, receives return the zero value immediately. Use the comma-ok idiom: `v, ok := <-ch`.
- **Mutex copy**: a `sync.Mutex` must not be copied after first use. Always use a pointer receiver for the queue type.

## Debugging walkthrough

Consider this broken worker pool:

```go
type Pool struct {
	ch chan int
	wg sync.WaitGroup
}

func (p *Pool) Start() {
	for i := 0; i < 4; i++ {
		go func() {
			for v := range p.ch {
				fmt.Println(v)
			}
		}()
	}
}
```

**Symptom**: Workers process some items but hang on exit; `main` never terminates.

**Investigation**: The `range p.ch` loop only exits when the channel is closed, but nowhere in the code does `close(p.ch)` get called. The workers block forever reading from an open, empty channel.

**Fix**: Close the channel after all producers finish:

```go
func main() {
	p := &Pool{ch: make(chan int, 100)}
	p.Start()
	for i := 0; i < 100; i++ {
		p.ch <- i
	}
	close(p.ch)  // signal workers to stop
	p.wg.Wait()
}
```

**Root cause**: The worker goroutines never receive their termination signal. The `range` loop over a channel only exits when the channel is closed. Without `close()`, the workers are orphaned goroutines.

## Production notes

- **Back-pressure**: when the queue reaches its capacity, producers should either block or fail fast, not pile memory. Use `select` with a default case: `select { case ch <- item: default: return ErrQueueFull }`.
- **Graceful shutdown**: in production, you need to drain in-flight work before shutdown. Track outstanding work with a `sync.WaitGroup` and signal workers via context cancellation.
- **Observability**: expose queue depth, enqueue rate, dequeue rate, and worker count as Prometheus metrics. A growing queue depth is the leading indicator of consumer failure.
- **Queue as a service boundary**: when a microservice provides a queue, document capacity, retention, and delivery guarantees. Clients depend on these semantics.

## Performance implications

- **Slice-backed queue**: enqueue is amortized O(1) (append grows capacity exponentially). Dequeue is O(1) but leaves a nil slot at the front. After many cycles, the underlying array is full of nil slots, wasting memory. Periodically re-slice: `q.items = append([]T(nil), q.items...)`.
- **Channel queue**: send/receive on a buffered channel is lock-free in the fast path. Cost is roughly 10-30ns per operation. A channel of capacity N uses a fixed-size circular buffer -- zero allocations after creation.
- **Contention**: when many goroutines contend on a mutex-backed queue, throughput collapses. Channel-based queues scale better because the runtime uses a lock-free internal buffer and a per-M (machine) run queue for goroutine scheduling.
- **Ring buffer**: for maximum throughput, use a pre-allocated ring buffer. It avoids slice reallocation and GC pressure. The Go runtime's own channel uses a ring buffer internally.

## Practice task

Write an in-memory queue that stores strings, then build a worker pool of 3 goroutines that reads from it. Each worker should convert the string to uppercase and print it. Enqueue at least 10 strings. Verify that all strings are processed and that the total matches the enqueued count.

## Tests / verification

```bash
go run ./curriculum/modules/14-architecture-distributed-systems/lessons/10-queues
go test ./curriculum/modules/14-architecture-distributed-systems/lessons/10-queues
```

The tests verify FIFO ordering, empty-queue error handling, closed-queue rejection, and the worker pool processing all items.

## Review questions

1. Why does a slice-backed queue need a mutex, but a channel-backed queue does not?
2. What happens to a buffered channel send when the buffer is full? What happens to a receive when the buffer is empty?
3. How does `range` over a channel behave differently from `range` over a slice?
4. If a worker pool's queue depth grows indefinitely, what is the most likely root cause?
5. What is the purpose of `close(ch)` in a producer-consumer pattern?

## NEXT UP

Retries and idempotent consumers.
