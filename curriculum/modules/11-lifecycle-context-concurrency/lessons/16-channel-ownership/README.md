# Channel ownership

## Learning objective

Apply the ownership discipline: the goroutine that creates a channel owns it, sends on it, and closes it; consumers only receive until the channel is closed.

## Why this matters

Channel misuse — multiple goroutines sending without coordination, receivers closing, double closes — is the leading cause of panics in concurrent Go programs. These bugs are often intermittent and hard to reproduce. The ownership discipline eliminates entire categories of channel bugs by establishing a clear contract: one owner creates, sends, and closes; everyone else receives. This pattern is not enforced by the compiler, but it is enforced by convention, code review, and production runbooks. Mastering ownership makes your concurrent programs predictable and safe.

## Mental model

Every channel has exactly one owner: the goroutine that created it. The owner is the channel's "manager". It decides when to put things into the channel and when to close it for business. All other goroutines are "customers": they can take things out of the channel but cannot add items or close it.

If a channel has multiple producers, they must coordinate through a manager goroutine that waits for all producers to finish and then closes the channel. The channel is still owned by the coordinator.

Ownership is about responsibility, not access. Any goroutine with a reference to the channel can technically send, receive, or close it. The discipline says: don't.

## Core idea

| Role | Responsibilities |
|---|---|
| Owner | Creates the channel, sends values, closes when done |
| Consumer | Receives values, uses comma-ok or `range`, never sends or closes |

```go
// Owner
func producer() <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)          // owner closes
        for i := 0; i < n; i++ {
            ch <- i              // owner sends
        }
    }()
    return ch                    // returns read-only channel
}

// Consumer
func consumer(ch <-chan int) {   // receives read-only channel
    for v := range ch {          // only receives
        fmt.Println(v)
    }
}
```

The owner returns a `<-chan int` (receive-only) to enforce that callers cannot send. The owner retains the send direction internally.

## Under the hood

Ownership is not enforced by the runtime — it is a discipline. The `hchan` struct's `closed` flag can be set by any goroutine holding a reference. There is no "owner ID" field. This means:

- Any goroutine can call `close(ch)`.
- Any goroutine can send.
- The runtime panics on double close or send after close regardless of who does it.

The discipline is purely about human coordination. What the runtime provides is the guarantee that once `close(ch)` returns, the channel is marked closed atomically. All subsequent sends (from any goroutine) will panic. All receives will return zero values.

## How Go uses it

- **Function returning a channel**: the standard pattern — `func NewListener() <-chan Event` — the function creates the channel, launches an internal goroutine that sends events and closes on shutdown. The caller receives only.
- **Pipeline stages**: each stage owns its output channel and closes it.
- **gRPC streams**: the server owns the send side of the stream channel.
- **Database row iterator**: `sql.Rows` owns an internal channel; iteration drains it; closing the rows closes the channel.
- **`time.After`**: the `time` package creates a channel, sends one value, and the GC reclaims it.

## Go example

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func counter(limit int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 1; i <= limit; i++ {
			ch <- i
			time.Sleep(10 * time.Millisecond)
		}
	}()
	return ch
}

func main() {
	ch := counter(5)

	// main is a consumer — it only receives
	for v := range ch {
		fmt.Println(v)
	}
}

// Output:
// 1
// 2
// 3
// 4
// 5
```

## Step-by-step execution

For `counter(3)`:

1. `counter` creates `ch` (unbuffered). It owns the write-side capability.
2. Goroutine starts. It is the owner goroutine.
3. `counter` returns the channel as `<-chan int` (receive-only to the caller).
4. Owner goroutine sends `1`: blocks until `main` (consumer) receives it.
5. Main receives `1` from `range`. Unblocks owner.
6. Owner sends `2`, main receives `2`, etc.
7. After sending `3`, owner calls `close(ch)`.
8. `range` in main receives the zero value with `ok=false` and exits the loop.

The owner goroutine knows the limit (3) and closes after the last send. The consumer has no knowledge of the limit — it just receives until the channel is closed.

## Common mistakes

- **Two goroutines both trying to close** — one panics with "close of closed channel".
  - Fix: designate a single owner. Use `sync.Once` if unavoidable.
- **Consumer calling close** — the producer continues sending and panics on the next send.
  - Fix: the producer always closes. If the consumer wants to signal "stop", use a separate `done` channel.
- **Returning a bidirectional channel** — callers may accidentally send or close:

```go
func producer() chan int {  // wrong: caller can send
```

  - Fix: return `<-chan int`.
- **Accepting a send-only channel for ownership** — a function that receives `chan<- T` cannot know if the channel is already owned by someone else.
  - Fix: if the function is the owner, document it. Prefer creating the channel inside the function and returning `<-chan T`.

## Debugging walkthrough

```go
func main() {
    ch := make(chan int)
    go func() {
        for i := 0; i < 5; i++ {
            ch <- i
        }
        close(ch)
    }()
    go func() {
        for v := range ch {
            fmt.Println("consumer 1:", v)
        }
    }()
    go func() {
        for v := range ch {
            fmt.Println("consumer 2:", v)
        }
    }()
    time.Sleep(time.Second)
}
```

**Symptom**: each value is received by only one consumer. Output is interleaved.

**Root cause**: each value in an unbuffered channel is delivered to exactly one receiver. This is not a bug — it is correct channel behavior. The ownership model says the producer owns the channel and closes it; consumers share the received values.

**Fix (if you want broadcast)**: use a `chan struct{}` close-broadcast pattern or fan-out with separate channels per consumer.

## Production notes

- **Document ownership** in the function doc comment: `// Process returns a channel that the caller owns. The caller must drain it and close is handled internally.`
- **Static analysis**: `go vet` catches some channel misuses but does not enforce ownership. Use code review.
- **Channel passing**: a channel can be sent over another channel (`chan chan T`). This transfers ownership. The receiving goroutine becomes the new owner.
- **For package-level channels**, create them in an `init()` function or a constructor and clearly document who owns them.
- **Test ownership by construction**: if your tests never call `close()` on a channel created by another function, you are following the discipline.

## Performance implications

The ownership discipline itself has zero runtime cost — it is a design pattern. The performance characteristics are those of the underlying channel operations. The key benefit is not performance but correctness: ownership eliminates panic-inducing edge cases that would otherwise crash the process.

## Practice task

Implement a `merger` function:

```go
func merge(channels ...<-chan int) <-chan int
```

1. Creates one output channel.
2. Launches one goroutine per input channel that forwards values to the output.
3. Uses a `sync.WaitGroup` to know when all input channels are exhausted.
4. Closes the output channel after all inputs are done.
5. The output channel is owned by `merge` — callers only receive.

Write a test that sends 3 values on 3 separate channels and verifies the merged output contains all 9 values.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/16-channel-ownership
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/16-channel-ownership
```

## Review questions

1. Who should close a channel — the producer or the consumer?
2. Why is it important to return `<-chan T` instead of `chan T` from a constructor function?
3. How do you handle multiple producers that share one channel? Who closes it?
4. What happens if two goroutines both call `close(ch)` on the same channel?
5. What does it mean to "transfer channel ownership" and how is it implemented?

## NEXT UP

Backpressure — controlling how fast producers can send when consumers are slow, using bounded channels and rejection patterns.
