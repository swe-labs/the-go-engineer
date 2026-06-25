# Backpressure

## Learning objective

Apply backpressure using bounded channels to control producer-consumer rates, implement load shedding with rejection patterns, and monitor queue depth in production.

## Why this matters

In any producer-consumer system, the producer can outpace the consumer. Without backpressure, the producer allocates unbounded memory (buffered channel grows, or values pile up in a slice), eventually causing OOM crashes or degrading latency for all consumers. Backpressure is the mechanism by which a slow consumer tells the producer "slow down." In Go, bounded channels provide backpressure by design: when the buffer is full, the producer blocks. This is the simplest, most robust backpressure mechanism available. Every production service — HTTP servers, event processors, queue consumers — must handle backpressure or fail under load.

## Mental model

Backpressure is like a pipe with a finite width. Water flows from the source (producer) through the pipe to the destination (consumer). If the destination drains slower than the source fills, the pipe fills up completely. When the pipe is full, water backs up — the source cannot add more water until the destination drains some.

The pipe's width is the channel capacity. When it is full, the producer blocks (stops generating). This blocking propagates backward through the system: the slow consumer ultimately slows the original request source, which is the correct behavior — the system self-regulates.

## Core idea

Three levels of backpressure:

1. **Blocking** (channel full → producer blocks): simplest, most common. The producer goroutine parks until a consumer drains a value. This propagates pressure upstream naturally.

2. **Load shedding** (channel full → drop the item): the producer skips non-critical work when the buffer is full. Uses `select` with `default` for non-blocking send.

3. **Rejection** (channel full → return error to caller): the API returns HTTP 429 or similar when the internal buffer is full.

```go
// Blocking backpressure (built-in)
ch := make(chan T, capacity)
ch <- item  // blocks if full

// Load shedding
select {
case ch <- item:
    // accepted
default:
    // buffer full, drop item
}

// Monitor queue depth
metrics.SetGauge("queue_depth", float64(len(ch)))
```

## Under the hood

Backpressure using bounded channels leverages the same `hchan` ring buffer mechanism. When `qcount == dataqsiz`, any send parks the goroutine in `sendq`. The goroutine is woken only when a consumer receives and a buffer slot opens. This is the runtime's built-in flow control.

The key insight: there is no special "backpressure" code path. Blocking send on a full channel is the normal buffered channel behavior. The backpressure emerges from the channel's design.

Load shedding with `select { case ch <- item: default: }` avoids parking entirely. The runtime's `select` implementation checks the `default` case when all channel cases would block. It executes `default` immediately without parking.

## How Go uses it

- **HTTP servers**: `net/http` has a fixed-size connection backlog (like a buffered channel). When the backlog is full, the OS rejects new connections.
- **Throttled access logs**: log lines are sent to a bounded channel; a single writer drains it. If the channel fills, the logging goroutine blocks, slowing the request.
- **Database connection pools**: the pool is a bounded channel of connections. Acquire blocks when all connections are in use (backpressure to the caller).
- **gRPC流**: flow control in HTTP/2 is a form of backpressure: the receiver advertises a window; the sender cannot exceed it.
- **Kubernete controllers**: work queues have bounded depth; if full, the controller refuses to enqueue more work.

## Go example

```go
package main

import (
	"fmt"
	"time"
)

func slowConsumer(ch <-chan int) {
	for v := range ch {
		fmt.Printf("consumed %d\n", v)
		time.Sleep(100 * time.Millisecond) // slow
	}
}

func fastProducer(ch chan<- int, done <-chan struct{}) {
	i := 0
	for {
		select {
		case ch <- i:
			i++
		case <-done:
			return
		}
	}
}

func main() {
	ch := make(chan int, 5)
	done := make(chan struct{})

	go slowConsumer(ch)
	go fastProducer(ch, done)

	time.Sleep(1 * time.Second)
	close(done)
	close(ch)
}
```

The buffer of 5 absorbs bursts. After 5 items, the producer blocks until the consumer drains one. The producer is rate-limited to the consumer's speed.

## Step-by-step execution

For buffer capacity 3, consumer sleeps 100ms per item:

1. Producer sends 0, 1, 2 — buffer fills to 3.
2. Producer tries to send 3: buffer full → goroutine parks in `sendq`.
3. Consumer receives 0 (after 100ms). Buffer now 2/3. Scheduler wakes producer.
4. Producer sends 3. Buffer is 3/3 again. Producer parks.
5. Consumer receives 1. Wakes producer. Producer sends 4.
6. This cycle continues. Producer is throttled to the consumer's rate of 10 items/second.

The average send rate of the producer equals the consumer's receive rate, regardless of how fast the producer could run.

## Common mistakes

- **Unbounded buffered channel**: `make(chan T, 1000000)` defeats backpressure. The producer can enqueue a million items before blocking, consuming 1M × sizeof(T) memory. If T is a pointer or struct, this can be gigabytes.
  - Fix: choose a small, measured buffer size. Let the producer block.
- **Confusing backpressure with load shedding**: blocking backpressure preserves every item but slows the producer. Load shedding drops items. Choose based on whether data loss is acceptable.
- **No monitoring**: a full buffer with blocking backpressure is invisible unless you monitor `len(ch)`. Add metrics.
- **Forgetting to drain**: if the consumer stops reading (e.g., due to a bug), the producer blocks forever. Combine with context cancellation for bounded blocking.

## Debugging walkthrough

```go
func main() {
    ch := make(chan int, 100)
    go func() {
        for i := 0; ; i++ {
            ch <- i
        }
    }()
    time.Sleep(time.Second)
    fmt.Println(len(ch))
    // Output: 100
}
```

**Symptom**: `len(ch)` is exactly the buffer capacity, never more.

**Investigation**: the producer blocks at capacity 100. There is no consumer, so the buffer stays at 100.

**Root cause**: no consumer goroutine is reading from the channel.

**Fix**: add a consumer. Or use a context with timeout to bound blocking:

```go
select {
case ch <- i:
case <-ctx.Done():
    return
}
```

## Production notes

- **Buffer size**: start with a small buffer (e.g., 10-100) and monitor `len(ch)/cap(ch)` in production. Increase only if you see excessive producer blocking with available consumer capacity.
- **Load shedding for non-critical work**: logging, analytics, and metrics can be dropped when the buffer is full. Use `select` with `default` and increment a counter.
- **Rejection for APIs**: when an internal buffer is full, the HTTP handler should return `429 Too Many Requests` or `503 Service Unavailable` instead of blocking the request goroutine.
- **Context timeouts**: combine backpressure with `context.WithTimeout` or `context.WithDeadline` so that blocked producers can abort if the consumer does not recover.
- **Backpressure propagates**: in a pipeline, a slow final stage backs up all previous stages. This is good — it means the entire system slows down together rather than crashing.

## Performance implications

Backpressure via blocking send has the same cost as any buffered channel send (lock + memcpy). The blocking/unblocking of the producer goroutine incurs scheduling overhead (~100 ns park + ~100 ns wake). This is negligible compared to the work being throttled.

Load shedding avoids scheduling overhead entirely (no park). It is the cheapest option when dropping items is acceptable.

Monitoring `len(ch)` is a single atomic read + subtraction (O(1)). Export it as a gauge for dashboards.

## Practice task

Build a simple HTTP server with backpressure:

1. Create a `jobs` buffered channel with capacity 10.
2. Launch 3 worker goroutines that receive from `jobs` and "process" (sleep 200ms).
3. The HTTP handler (listen on `:8080`) does a non-blocking send to `jobs`:
   - If accepted, return `202 Accepted`.
   - If buffer full, return `503 Service Unavailable`.
4. Test with 20 concurrent requests. Observe how many are accepted vs rejected.
5. Add a `/metrics` endpoint that returns `queue_depth` as a gauge.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/17-backpressure
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/17-backpressure
```

## Review questions

1. How does a bounded buffered channel provide backpressure?
2. What is the difference between blocking backpressure and load shedding?
3. How can you monitor queue depth in production?
4. What happens to the producer when the buffer is full and no consumer is reading?
5. Why is an unbounded buffered channel dangerous in a production system?

## NEXT UP

Pipelines — compose multiple stages connected by channels, with fan-out, fan-in, and cancellation propagation.
