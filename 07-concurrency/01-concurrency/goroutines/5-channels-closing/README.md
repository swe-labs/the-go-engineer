# GC.5 Closing Channels: Signaling the End

## Mission

Learn how to gracefully signal the end of a data stream or broadcast a system-wide shutdown using `close()`. Master the patterns for safely consuming data until a channel is empty and avoid the "panic on send" disaster.

## Prerequisites

- `GC.4` buffered-channels

## Mental Model

Think of Closing a Channel as **Turning Off the Lights at a Diner**.

1. **The Signal**: When the owner (the sender) flips the switch (`close`), they are saying: "We aren't taking any more orders."
2. **The Buffer**: If there are already orders being prepared (items in the buffer), the customers (receivers) can still finish them.
3. **The Zero Value**: Once the food is gone, anyone else who walks in gets an immediate "Sorry, we're closed" (Zero Value + `ok=false`).

## Visual Model

```mermaid
graph TD
    S[Sender] -- "Send val 1" --> C[Channel]
    S -- "Send val 2" --> C
    S -- "close(ch)" --> C
    C -- "Recv val 1" --> R[Receiver]
    C -- "Recv val 2" --> R
    C -- "Recv Zero Value" --> R
    Note over R: "ok" becomes false
```

## Machine View

In Go's `hchan` struct, `close()` sets a `closed` flag to 1.
- **Receivers**: Any goroutines currently blocked on `recvq` are immediately unparked and receive the zero value.
- **Senders**: Any goroutine trying to send to a closed channel will **panic** immediately. This is because a send to a closed channel is considered a logical bug in your concurrent design.
- **Garbage Collection**: Channels are eventually cleaned up by the GC even if you don't close them. Closing is for **communication**, not for memory management.

## Run Instructions

```bash
go run ./07-concurrency/01-concurrency/goroutines/5-channels-closing
```

## Code Walkthrough

### `range` Loop
The `for val := range ch` loop is the idiomatic way to consume a channel. It automatically stops as soon as the channel is closed and emptied.

### The Comma-OK Pattern
`value, ok := <-ch` allows you to manually check if a channel is still open. If `ok` is false, the channel is closed and `value` is the zero value of the type.

### Broadcast Shutdown
A closed channel unblocks **all** receivers. This is a powerful pattern for graceful shutdown: one close signal can tell 1,000 workers to stop immediately.

## Try It

1. Comment out `close(taskQueue)` in the first example. Run the code. Why does it deadlock at the end? (Hint: The `range` loop is waiting for a 6th item that will never come).
2. Try sending a value to `signals` **after** it has been closed. Observe the panic.
3. Replace `<-shutdown` with `time.Sleep(1 * time.Second)` in the worker. Notice how the broadcast signal is much more precise and immediate.

## Verification Surface

Observe the two patterns for detection and the broadcast signal at the end:

```text
=== Closing Channels ===

1) range over channel:
   Step 1: Compile source code -> done
   ...
   Pipeline complete!

2) comma-ok pattern:
   Received: 42 (ok=true)
   Channel closed - no more data

3) close() as broadcast signal:
   Worker 1: waiting for shutdown signal...
   Main: sending shutdown signal...
   Worker 1: received shutdown, cleaning up
```

## In Production
**Only close from the sender side.**
If you have multiple senders, closing is more complex (you might need a `sync.Once` or a coordinator). If you have one sender and multiple receivers, closing is the perfect way to say "I'm done sending."

## Thinking Questions
1. Why does Go return the zero value of the type after a channel is closed?
2. What happens if you try to close a channel that is already closed?
3. In what scenario would you *not* want to close a channel?

## Next Step

Next: `GC.6` -> `07-concurrency/01-concurrency/goroutines/6-project-1`

Open `07-concurrency/01-concurrency/goroutines/6-project-1/README.md` to continue.
