# Closing channels

## Learning objective

Use `close(ch)` to signal completion, detect closed channels with the comma-ok idiom, iterate with `range`, and apply the "who closes?" discipline.

## Why this matters

An unclosed channel causes goroutine leaks: a receiver ranging over the channel blocks forever. A double-closed channel panics. A send on a closed channel panics. These are not edge cases — they are the most common channel bugs in production Go systems. Mastering `close` semantics — what happens when, who should call it, and how receivers detect it — is essential for writing correct concurrent programs.

## Mental model

Closing a channel is like putting a "no more deliveries" sign on a mailbox. The mail carrier stops delivering new letters. Recipients can still take any letters already in the box. After the last letter is taken, the mailbox reports "no more letters" (the zero value with `ok=false`). Once closed, no one can add new letters to the box.

You can close a channel only from the sender side. The sender knows when all values have been sent. The receiver only knows when the channel is closed. This asymmetry is the foundation of the ownership discipline.

## Core idea

```go
close(ch)              // signal: no more sends
v, ok := <-ch          // ok == false when ch is closed
for v := range ch { }  // exits when ch is closed
```

Rules:

- Only the sending goroutine should close a channel. Never close from the receiver.
- Closing a channel is not required. Channels are garbage collected when unreferenced. But `range` loops and receivers that need to know "done" require `close`.
- Sending on a closed channel panics: `panic: send on closed channel`.
- Closing a closed channel panics: `panic: close of closed channel`.
- Closing a nil channel panics: `panic: close of nil channel`.
- Receiving from a closed channel always succeeds, returning the zero value with `ok=false`, immediately (no blocking).

## Under the hood

`close(ch)` calls `runtime.closechan`. The function:

1. Locks the channel.
2. Checks `closed` flag — if already set, unlocks and panics.
3. Sets `closed = 1`.
4. Iterates over `recvq`: for each goroutine in the receive wait queue, dequeues it, copies the zero value to its receive variable with `ok=false`, and adds it to a "run next" list (not immediately scheduled, but marked runnable).
5. Iterates over `sendq`: for each goroutine in the send wait queue, dequeues it and schedules a panic (`"send on closed channel"`) when that goroutine runs.
6. Unlocks and wakes all the goroutines from step 4.

Key detail: the zero value is only delivered to goroutines that are already blocked on receive when `close` is called. A later receive that happens after close sees the buffer (if any) drained, then zero values with `ok=false`.

## How Go uses it

- **Pipeline termination**: the final stage closes its output channel when all input has been processed. Upstream stages detect this via `range`.
- **Broadcast signal**: `close(done)` in a `select` broadcasts to all waiting goroutines that work should stop.
- **Worker pool shutdown**: close the jobs channel to signal workers that no more work is coming. Workers exit their `range` loop cleanly.
- **Result collection**: the producer closes the results channel; the consumer ranges until done.

## Go example

```go
package main

import (
	"fmt"
	"time"
)

func produce(ch chan<- int) {
	for i := 1; i <= 5; i++ {
		ch <- i
		time.Sleep(10 * time.Millisecond)
	}
	close(ch)
}

func main() {
	ch := make(chan int, 3)
	go produce(ch)

	for v := range ch {
		fmt.Println("received:", v)
	}
	fmt.Println("producer done, channel closed")
}
```

Output:

```
received: 1
received: 2
received: 3
received: 4
received: 5
producer done, channel closed
```

## Step-by-step execution

For `produce` sending 1..5 into a buffered channel of capacity 3:

1. `produce` sends 1, 2, 3 without blocking (buffer 3).
2. Send 4 blocks: buffer full, producer parks in `sendq`.
3. Main receives 1: copies from buffer, wakes producer. Producer sends 4, blocks again (buffer was 2/3 before send).
4. Main receives 2: wakes producer. Producer sends 5, buffer 2/3, returns.
5. Producer calls `close(ch)`. Buffer contains `[3, 4, 5]`. No goroutines in `recvq`. `closed=1`.
6. Main receives 3: `v=3, ok=true`. Buffer `[4, 5]`.
7. Main receives 4, 5. Buffer empty.
8. Main tries next `range` iteration: buffer empty, `closed=true`, so returns zero value with `ok=false`. `range` exits.
9. Main prints "producer done, channel closed".

## Common mistakes

- **Closing from the receiver**: the receiver does not know if all sends are done. If the receiver closes and the sender tries to send, the sender panics.
  - Fix: the producer always closes.
- **Double close**: calling `close(ch)` twice panics.
  - Fix: use `sync.Once` if multiple code paths could close, or restructure to have a single closer.
- **Sending after close**: the sender must know the channel is closed. The ownership discipline (one producer) prevents this. For multiple senders, use a separate `done` channel or `sync.WaitGroup`.
- **Not closing a channel**: a `range` loop never exits, leaking the receiving goroutine.
- **Closing a nil channel**: `var ch chan int; close(ch)` panics. Always `make` a channel before closing.

## Debugging walkthrough

```go
func main() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 3; i++ {
			ch <- i
		}
	}()
	for v := range ch {
		fmt.Println(v)
	}
}
```

**Symptom**: program prints 0, 1, 2 and then hangs forever.

**Investigation**: the `range` loop expects the channel to be closed when no more values arrive. The producer never calls `close(ch)`.

**Root cause**: missing `close(ch)` in the producer goroutine.

**Fix**: add `close(ch)` after the loop in the producer:

```go
go func() {
    for i := 0; i < 3; i++ {
        ch <- i
    }
    close(ch)
}()
```

## Production notes

- **Use `defer close(ch)`** in the producer to ensure the channel is closed even if the producer panics or returns early.
- **For multiple producers**, use a separate `sync.WaitGroup` to know when all producers are done, then close the channel in a dedicated goroutine:

```go
go func() {
    wg.Wait()
    close(ch)
}()
```

- **Checking `closed` is not exported** — you cannot call a function `isClosed(ch)` (it would race anyway). Use the comma-ok idiom on receive.
- **`close(ch)` for broadcast**: closing a `chan struct{}` unblocks all goroutines waiting on `<-ch` or `case <-ch:`. This is the idiomatic way to signal shutdown to multiple workers.
- **Panic recovery**: if untrusted code might send on a closed channel, wrap sends in a `recover`. But prefer to fix the design.

## Performance implications

`close(ch)` is O(N) where N is the number of goroutines blocked on receive (or send). Each blocked goroutine is woken individually. In practice, N is usually small (single digits) because channels are designed for communication between a few goroutines, not broadcast to thousands. For broadcast to many goroutines, use `close(done)` pattern with `select` in each worker rather than relying on a single channel with thousands of waiters.

## Practice task

Write a function `fanOut(ch <-chan int, n int) []<-chan int` that:

1. Takes an input channel and a count `n`.
2. Creates `n` output channels.
3. Launches `n` goroutines, each reading from `ch` and forwarding to its own output channel.
4. Closes all output channels when `ch` is closed.
5. Verify that draining all output channels yields the same values as reading from `ch`.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/15-closing-channels
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/15-closing-channels
```

## Review questions

1. What does `v, ok := <-ch` return when `ch` is closed and empty?
2. What happens if you call `close(ch)` twice?
3. Who should close a channel — the sender or the receiver? Why?
4. How do you make a `range` loop over a channel exit cleanly?
5. What happens to goroutines blocked on `<-ch` when the channel is closed?

## NEXT UP

Channel ownership — the discipline of who creates, sends, closes, and receives for safe channel usage.
