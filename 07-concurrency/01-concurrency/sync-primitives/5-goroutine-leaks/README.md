# SY.5 Goroutine Leaks: Memory Silent Killers

## Mission

Learn to identify and prevent **Goroutine Leaks**-one of the most common causes of memory exhaustion in production Go services. Understand why goroutines get stranded and how to use `context` to ensure every concurrent task has a defined lifetime.

## Prerequisites

- `SY.4` race-conditions

## Mental Model

Think of a Goroutine Leak as **A Running Faucet in a Locked Room**.

1. **The Faucet (`go func`)**: You turned it on to get some water (work).
2. **The Lock (`blocking`)**: The drain is clogged (blocked channel), or the room is locked (return condition never met).
3. **The Flood (`OOM`)**: The water keeps rising (memory usage). Eventually, the entire house (your server) collapses under the weight of the water.

## Visual Model

```mermaid
graph TD
    M[Main Program] -- "Spawn" --> G1[Worker 1]
    G1 -- "Send to Ch" --> C[Channel]
    Note right of C: No one is receiving!
    C -- "Block" --> G1
    Note over G1: Stays in RAM forever
```

## Machine View

In Go, a goroutine is extremely cheap (starting at ~2KB of stack space), but it is **not free**.
- If a goroutine is blocked on a channel send/receive that will never complete, the Garbage Collector **cannot** clean it up.
- The stack, the local variables, and the goroutine's control structure (`g` struct) remain in the heap indefinitely.
- **Detection**: Use `runtime.NumGoroutine()` to monitor the count. In a healthy service, this number should remain stable or fluctuate within a range. If it only goes up, you have a leak.

## Run Instructions

```bash
go run ./07-concurrency/01-concurrency/sync-primitives/5-goroutine-leaks
```

## Code Walkthrough

### The Blocked Sender
In `leakGenerator`, we create an unbuffered channel and try to send to it inside a goroutine. Since no one ever reads from that channel, the goroutine is suspended forever. It will never return, and its memory will never be reclaimed.

### The Context Pattern
In `safeWorker`, we use a `select` statement to listen to a `context.Context`. This is the professional standard for goroutine management. When the parent calls `cancel()`, the worker receives the signal and returns immediately.

### Monitoring
`runtime.NumGoroutine()` is your best friend in production. Most Go monitoring tools (like Prometheus or Datadog) track this metric automatically.

## Try It

1. Launch 10,000 leaked goroutines in a loop. Watch your computer's memory usage in Task Manager / Activity Monitor.
2. Use a buffered channel in the `leakGenerator`. Does it still leak? (Hint: Yes, as soon as the buffer is full).
3. Implement a "Timeout" using `context.WithTimeout` to automatically clean up workers that take too long.

## Verification Surface

Observe the goroutine count increasing and staying high for the leaked task:

```text
=== SY.5 Goroutine Leaks ===

Initial Goroutines: 1

Scenario 1: Creating a leak...
  [Leak] Goroutine started, trying to send...
Current Goroutines: 2 (Leaked: 1)

Scenario 2: Creating a safe worker...
Final Goroutines: 2 (Leaked remains, safe exited)
```

## In Production
**Never start a goroutine without a shutdown plan.**
The #1 rule of Go concurrency is: **Know when a goroutine will stop.** If you can't point to the line of code that will cause the goroutine to return, you have a bug. Common leak sources include:
- Blocked sends on unbuffered channels.
- `for range` loops on channels that are never closed.
- Infinite `for { select {} }` loops without a cancellation case.

## Thinking Questions
1. Why can't the Go Garbage Collector detect that a goroutine is "stuck" and kill it?
2. How does a goroutine leak affect the performance of the surviving goroutines?
3. What is the difference between a goroutine leak and a memory leak?

## Next Step

We've seen goroutines that stay alive too long. Now let's look at the opposite: goroutines that get stuck waiting for each other, bringing the whole system to a halt. Continue to [SY.6 Deadlocks](../6-deadlocks/README.md).
