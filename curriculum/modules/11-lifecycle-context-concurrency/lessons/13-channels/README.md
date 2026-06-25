# Channels

## Learning objective

Create, send to, and receive from unbuffered channels, and explain their blocking semantics and role as synchronization primitives.

## Why this matters

Channels are Go's primary mechanism for goroutine communication. Instead of sharing memory (a variable both goroutines can write) and protecting it with locks, channels pass ownership of data between goroutines. This "share memory by communicating" philosophy eliminates entire categories of data races and makes concurrent data flow explicit in the type system. Every concurrent Go program — HTTP servers, stream processors, pipeline workers — uses channels to move data between goroutines safely.

## Mental model

An unbuffered channel is a synchronous pipe. The sender blocks until a receiver is ready, and the receiver blocks until a sender is ready. Think of it as a handoff: sender extends a closed fist; receiver extends an open hand; at the moment they touch, the sender opens their fist and the receiver closes theirs. The data moves directly from sender's hand to receiver's hand — no buffer, no waiting area.

This "rendezvous" semantics ensures that send and receive happen at the same instant. The channel itself is just the agreement to meet.

## Core idea

```go
ch := make(chan int)   // unbuffered channel of ints
ch <- 42               // send: blocks until someone receives
v := <-ch              // receive: blocks until someone sends
close(ch)              // no more sends
```

Key properties:

- Send and receive block until both are ready (synchronous).
- The type `chan T` is a pointer — channels are reference types, passed by value means sharing the same channel.
- `chan<- T` is send-only, `<-chan T` is receive-only. Use directional channels in function signatures to enforce intent.
- Receiving from a closed channel returns the zero value immediately (with `ok=false` from comma-ok idiom).
- Sending to a closed channel panics.

## Under the hood

A channel is a runtime object of type `hchan` (defined in `runtime/chan.go`). The struct contains:

- `qcount` and `dataqsiz`: buffer state (0 for unbuffered).
- `buf`: pointer to the circular buffer (nil for unbuffered).
- `elemsize` and `elemtype`: element metadata.
- `sendx`, `recvx`: buffer indices (unused for unbuffered).
- `recvq`: linked list of goroutines blocked on receive.
- `sendq`: linked list of goroutines blocked on send.
- `lock`: mutex protecting the channel state.
- `closed`: flag.

When a goroutine sends to an unbuffered channel:

1. Lock the channel.
2. If a goroutine is waiting in `recvq`, dequeue it, copy the value directly from sender's stack to receiver's stack, unlock, and wake the receiver.
3. If no receiver is waiting, enqueue the sender goroutine in `sendq`, unlock, and park the goroutine (block).

Receive is symmetric: if a sender is in `sendq`, dequeue and copy; if not, block in `recvq`.

The direct copy from sender to receiver (step 2) avoids allocating a buffer slot. This is the "rendezvous" — the data moves from one goroutine's stack to the other without intermediate storage.

## How Go uses it

- **Request/response**: one goroutine sends a request, another receives and sends a response.
- **Coordination**: `done := make(chan struct{})` as a signal that a goroutine has completed.
- **Pipeline handoff**: each stage receives from its input channel and sends to its output channel.
- **Timeouts and cancellation**: `select` with `time.After` or `ctx.Done()`.
- **Rate limiting**: a channel as a token bucket with `select` and `default`.

## Go example

```go
package main

import (
	"fmt"
	"time"
)

func sender(ch chan<- int) {
	for i := 1; i <= 3; i++ {
		fmt.Printf("sender: sending %d\n", i)
		ch <- i
		fmt.Printf("sender: sent %d\n", i)
	}
	close(ch)
}

func main() {
	ch := make(chan int)
	go sender(ch)

	for v := range ch {
		fmt.Printf("main: received %d\n", v)
		time.Sleep(20 * time.Millisecond)
	}
	fmt.Println("main: channel closed, exiting")
}
```

Output:

```
sender: sending 1
main: received 1
sender: sent 1
sender: sending 2
main: received 2
sender: sent 2
sender: sending 3
main: received 3
sender: sent 3
main: channel closed, exiting
```

Notice the alternation: sender blocks on send until main receives, proving the rendezvous.

## Step-by-step execution

For `ch <- 42` in sender goroutine, `v := <-ch` in main:

1. Sender evaluates `ch <- 42`. Locks `ch`. Checks `recvq`.
2. No receiver waiting → allocates a `sudog` (internal runtime struct representing a blocked goroutine), enqueues it in `sendq`, unlocks, parks sender.
3. Main evaluates `<-ch`. Locks `ch`. Checks `sendq`.
4. Finds the sender's `sudog` in `sendq`. Copies the value (42) directly from sender's stack to main's stack.
5. Unlocks `ch`, wakes the sender goroutine.
6. Sender resumes — the send operation completes. Main continues with `v == 42`.

The entire handoff is synchronous: sender does not proceed until main has received.

## Common mistakes

- **Sending on an unbuffered channel without a receiver** — the goroutine blocks forever. If it is the only goroutine, the program deadlocks. Always ensure a matching receive exists.
- **Forgetting to close a channel** — a `range` loop over a channel never exits, leaking the receiving goroutine. Close the channel when all sends are done.
- **Sending after closing** — panics: `send on closed channel`. Use the ownership discipline (one producer, closes when done).
- **Using channels for simple locking** — a channel with a single-element buffer can act as a mutex, but `sync.Mutex` is clearer and cheaper.
- **Nil channels always block** — `var ch chan int` is nil. Sending to or receiving from nil blocks forever. Use `make` to create a usable channel.

## Debugging walkthrough

```go
package main

func main() {
	ch := make(chan int)
	ch <- 1  // blocks forever
}
```

**Symptom**: program hangs immediately with `fatal error: all goroutines are asleep - deadlock!`.

**Investigation**: the runtime detects that the only goroutine (main) is blocked trying to send with no receiver.

**Root cause**: unbuffered channel send without a concurrent receive.

**Fix**: launch a receiver goroutine before sending:

```go
ch := make(chan int)
go func() {
    v := <-ch
    fmt.Println(v)
}()
ch <- 1
```

## Production notes

- **Use directional channel types** in function signatures: `func produce(ch chan<- int)` and `func consume(ch <-chan int)`. The compiler enforces that producers only send and consumers only receive, catching misuse at compile time.
- **Unbuffered channels provide backpressure by design**: a producer cannot outrun a consumer because each send blocks until the consumer receives. This is a feature, not a bug.
- **Channel of `struct{}`** is the idiomatic signal type. It conveys zero information (zero-size value) and indicates synchronization-only intent.
- **Avoid very large values on channels**: the value is copied from sender to receiver. For large structs, pass a pointer (`chan *LargeStruct`), but be mindful of ownership — the pointer must not be mutated after sending.
- **Beware goroutine leaks**: a goroutine blocked on a send with no receiver is a leak. Always ensure a matching receiver will eventually run.

## Performance implications

| Operation | Approximate cost |
|---|---|
| Unbuffered channel send/receive (rendezvous) | ~50-100 ns |
| Channel lock contention | Depends on goroutine count |
| Sudog allocation | Heap allocation each time a goroutine blocks |

The rendezvous is not free: the runtime must lock the channel, search the wait queue, copy the value, and schedule the woken goroutine. In high-throughput systems, batching (buffered channels) or alternative patterns (batch processing) reduce channel operations per unit of work.

Channel operations are still orders of magnitude cheaper than OS-level IPC. On modern hardware, millions of channel operations per second are achievable.

## Practice task

Write a function `pingPong(count int)` that:

1. Creates an unbuffered `chan string`.
2. Launches two goroutines: "ping" sends `"ping"` and waits for `"pong"`; "pong" sends `"pong"` and waits for `"ping"`.
3. Alternates `count` times.
4. Uses a second channel to signal completion to `main`.

Run and observe the strict alternation.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/13-channels
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/13-channels
```

## Review questions

1. What happens when you send to an unbuffered channel with no receiver ready?
2. What is the zero value of a channel variable (`var ch chan int`)? What happens when you send to it?
3. How does an unbuffered channel provide backpressure?
4. What does `close(ch)` do to a blocked receive operation?
5. Why does `for v := range ch` block forever if the channel is never closed?

## NEXT UP

Buffered channels — decouple sender and receiver with a bounded buffer, enabling asynchronous communication.
