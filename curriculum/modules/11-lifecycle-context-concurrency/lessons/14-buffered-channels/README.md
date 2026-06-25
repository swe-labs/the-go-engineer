# Buffered channels

## Learning objective

Create and use buffered channels with explicit capacity, explain when send blocks and when receive blocks, and apply the channel-as-semaphore pattern.

## Why this matters

Unbuffered channels force a rendezvous: every send must meet a receive. This guarantees synchronization but limits throughput. Buffered channels decouple sender and receiver: the sender can enqueue up to `capacity` items without waiting. This enables burst handling, rate limiting, and pipeline stages where producers are faster than consumers. Understanding buffered channel blocking semantics is essential for building efficient concurrent systems that absorb load spikes without dropping data.

## Mental model

A buffered channel is a mailbox with `capacity` slots. The sender drops a letter into a slot and leaves — no need to wait for the recipient to open it. If all slots are full, the sender must wait until a slot opens. The recipient checks the mailbox: if there is a letter, they take it; if the mailbox is empty, they wait.

The buffer decouples the timing of send and receive. Unlike the unbuffered handoff, sender and receiver operate at their own pace, within the buffer's capacity.

## Core idea

```go
ch := make(chan T, capacity)
```

- Send blocks **only when the buffer is full**.
- Receive blocks **only when the buffer is empty**.
- `cap(ch)` returns the buffer capacity.
- `len(ch)` returns the current number of buffered elements.

| State | Send | Receive |
|---|---|---|
| Empty (len=0) | Proceeds (non-blocking) | Blocks |
| Partial (0 < len < cap) | Proceeds | Proceeds |
| Full (len=cap) | Blocks | Proceeds |

As a semaphore: a buffered channel with capacity `N` and `struct{}` elements acts as a counting semaphore. Send acquires a token (blocks if none available), receive releases it.

```go
sem := make(chan struct{}, maxConcurrent)
sem <- struct{}{}   // acquire (block if at capacity)
// ... do work ...
<-sem               // release
```

## Under the hood

A buffered channel's `hchan` struct (same as unbuffered) sets `dataqsiz` to the capacity. The `buf` field points to a circular buffer of `dataqsiz * elemsize` bytes. The `sendx` and `recvx` indices track the next send and receive positions in the circular buffer.

When a goroutine sends to a buffered channel with space:

1. Lock the channel.
2. Copy the value into `buf[sendx]`.
3. Increment `sendx` (wrap around the ring).
4. Increment `qcount`.
5. If a receiver is waiting on `recvq`, wake one receiver (it will take from `buf[recvx]`).
6. Unlock, return.

When the buffer is full, the sender enqueues in `sendq` and parks, just like unbuffered.

The key difference: with buffer space, the send completes without a matching receiver. The value sits in the ring buffer until a receiver arrives.

## How Go uses it

- **Worker pools**: a buffered job channel with capacity `N` holds pending jobs. Workers receive when idle.
- **Rate limiting**: a buffered channel of tokens released at a fixed interval (ticker + bucket).
- **Request batching**: buffer incoming requests up to a limit before processing a batch.
- **Log aggregation**: multiple goroutines send log lines to a buffered channel; a single writer drains and flushes to disk.
- **Connection pooling**: a buffered channel of `net.Conn` values acts as a reusable pool.

## Go example

```go
package main

import (
	"fmt"
	"sync"
)

const maxWorkers = 3

func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Printf("worker %d processing job %d\n", id, j)
	}
}

func main() {
	jobs := make(chan int, 10)
	var wg sync.WaitGroup

	for i := 1; i <= maxWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}

	for j := 1; j <= 10; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
}
```

Output (order may vary):

```
worker 3 processing job 1
worker 1 processing job 2
worker 2 processing job 3
worker 3 processing job 4
...
```

The buffer of 10 allows the producer to enqueue all jobs without waiting for workers to be ready.

## Step-by-step execution

For `ch := make(chan int, 3)` with sends `ch <- 1; ch <- 2; ch <- 3; ch <- 4`:

1. Send `1`: buffer at 0/3 → copy `1` to `buf[0]`. `qcount=1`, `sendx=1`. Returns immediately.
2. Send `2`: `qcount=2`, `sendx=2`. Returns.
3. Send `3`: `qcount=3`, `sendx=0` (wrap). Returns.
4. Send `4`: `qcount == dataqsiz` → buffer full. Enqueue goroutine in `sendq`, park. Blocks until a receive frees a slot.

When a receiver does `<-ch`:
1. Copy `buf[recvx]` (`recvx=0`, value `1`) to receiver.
2. `recvx=1`, `qcount=2`.
3. If goroutines are waiting in `sendq`, dequeue one, copy its value to `buf[sendx]`, and wake it.

## Common mistakes

- **Using buffered channels when unbuffered is correct**: a buffered channel with capacity 1 is not the same as an unbuffered channel. Capacity 1 lets one send proceed without a receiver; unbuffered always requires the rendezvous. For synchronization signals (e.g., `done`), use unbuffered `chan struct{}`.
- **Assuming buffered channel sends never block**: they block when full. In high-throughput systems, a full buffer causes backpressure. Monitor `len(ch)` and `cap(ch)` in production.
- **Closing too early**: if you close a buffered channel while values remain in the buffer, receivers can still read them. The channel delivers remaining values, then returns zero values with `ok=false`. This is correct behavior, but the sender must not send after closing.
- **Mixing send-only and receive-only channel directions incorrectly**: a `chan<- T` (send-only) can be created from a bidirectional `chan T` with implicit conversion. But you cannot convert a send-only back to bidirectional.

## Debugging walkthrough

```go
func main() {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	ch <- 3  // blocks here
	fmt.Println(<-ch)
}
```

**Symptom**: program deadlocks. The runtime prints `fatal error: all goroutines are asleep - deadlock!`.

**Investigation**: the main goroutine is the only goroutine. After sending 3 values to a buffer of capacity 2, the third send blocks. But there is no concurrent goroutine to receive.

**Root cause**: buffer capacity exceeded with no concurrent consumer.

**Fix**: launch a consumer goroutine before sending, or increase buffer capacity, or both.

Alternatively, use `select` with `default` for non-blocking send:

```go
select {
case ch <- 3:
default:
	fmt.Println("buffer full, dropping value")
}
```

## Production notes

- **Choose buffer size based on measurement, not guesswork**. Profile the producer-consumer rate difference. A buffer that is too large wastes memory; too small causes excessive blocking.
- **Monitor queue depth**: `len(ch)` tells you the current backlog. Export this as a Prometheus gauge in production:

```go
queueDepth := float64(len(jobCh))
metrics.SetGauge("job_queue_depth", queueDepth)
```

- **Buffered channels are not for persistent queues**. If the process crashes, buffered values in channels are lost. Use proper message queues (Kafka, RabbitMQ) for durability.
- **The channel-as-semaphore pattern** is elegant but limited: you cannot time out on acquire without `select`. For complex semaphore needs, consider `golang.org/x/sync/semaphore`.

## Performance implications

| Operation | Approximate cost |
|---|---|
| Buffered send (space available) | ~20-40 ns (lock + memcpy into ring) |
| Buffered receive (data available) | ~20-40 ns |
| Park (buffer full/empty) | Same as unbuffered (~100-200 ns with wake) |

Buffered operations are cheaper than unbuffered in the non-blocking case because no goroutine scheduling happens — just a lock, a memcpy into the ring buffer, and unlock. This makes buffered channels suitable for high-throughput data transfer between goroutines.

The ring buffer is pre-allocated at `make` time. There is no per-operation heap allocation (the value is copied into the pre-allocated buffer).

## Practice task

Implement a simple rate limiter using a buffered channel as a token bucket:

1. Create a buffered channel of `struct{}` with capacity 5.
2. Launch a goroutine that fills the channel with one token every 200ms (up to capacity).
3. In `main`, loop 10 times attempting to acquire a token (receive from channel) with a timeout of 100ms using `select`.
4. Print "processed" on success, "rate limited" on timeout.
5. Measure how many succeed and how many are rate limited.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/14-buffered-channels
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/14-buffered-channels
```

## Review questions

1. When does a send to a buffered channel block?
2. What does `cap(ch)` vs `len(ch)` return for a buffered channel?
3. Can you range over a buffered channel without closing it? What happens?
4. How does a buffered channel with capacity 1 differ from an unbuffered channel?
5. What is the channel-as-semaphore pattern and when would you use it?

## NEXT UP

Closing channels — the `close` function, the comma-ok idiom, range over channel, and the "who closes?" contract.
