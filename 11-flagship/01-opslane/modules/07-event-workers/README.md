# OPSL.7 Event Bus and Worker Pools

## Mission

Add bounded asynchronous work without turning background processing into invisible chaos.

## What This Module Builds

- event bus primitives
- worker pool boundaries
- queueing and backpressure behavior
- safe drain rules for background work

## You Are Here If

- `OPSL.6` is complete
- payment and order workflows are ready to emit asynchronous follow-up work
- you are prepared to reason about bounded concurrency

## Proof Surface

This module is implemented in the current tree.

Run:

```bash
go test ./11-flagship/01-opslane/internal/events/...
go test ./11-flagship/01-opslane/internal/workers/...
go run ./11-flagship/01-opslane/scripts/progress.go
```

The proof surface now covers:

- bounded event publishing
- queue-full backpressure
- context-aware publishing and submission
- worker pool drain behavior
- handler error reporting
- order, payment, and notification worker adapters

Implemented files:

- `internal/events/bus.go`
- `internal/events/types.go`
- `internal/workers/pool.go`
- `internal/workers/order_processor.go`
- `internal/workers/payment_processor.go`
- `internal/workers/notification_worker.go`

Implementation map: [SURFACE.md](./SURFACE.md)

## Required Files and Boundaries

Worker coordination should stay explicit and bounded.
Avoid untracked goroutine spawning inside handlers or repositories.

## Mental Model

Async work is not "fire and forget".

It is:

```text
event -> bounded queue -> fixed workers -> explicit handler -> observable error path
```

If the queue is full, the system must decide what to do.
OPSL.7 keeps that decision visible by returning `ErrQueueFull` instead of hiding unlimited buffering.

## Machine View

When code publishes an event, the machine does not create infinite work.
It tries to put one small event into one bounded channel.

```text
Publisher
    |
    v
Event Bus channel with capacity N
    |
    v
Worker Pool channel with capacity M
    |
    v
Fixed number of worker goroutines
    |
    v
OrderProcessor / PaymentProcessor / NotificationWorker
```

The important constraint:

```text
queue capacity + worker count = the upper bound of in-memory async pressure
```

That is what prevents "one request creates unlimited goroutines" from becoming a production outage.

## Diagram

```text
HTTP/service code
      |
      v
 events.Bus.TryPublish
      |
      | queue full -> ErrQueueFull
      v
 events.Event
      |
      v
 workers.Pool.TrySubmit
      |
      | queue full -> ErrQueueFull
      v
 fixed worker goroutine
      |
      +--> OrderProcessor
      +--> PaymentProcessor
      +--> NotificationWorker
```

## Try It

Change a worker pool test from `QueueSize: 1` to `QueueSize: 2`.
Then observe that the second `TrySubmit` no longer returns `ErrQueueFull` until the queue reaches the new capacity.

Next, change `Workers: 2` to `Workers: 1` in the drain test.
The test should still pass, but work drains with less parallel capacity.

## Engineering Questions

- What happens when the queue is full?
- When should the system reject work instead of buffering more?
- How do you inspect or drain a stuck worker safely?

## Next Step

Next: `OPSL.8` -> `11-flagship/01-opslane/modules/08-caching`

Open `11-flagship/01-opslane/modules/08-caching/README.md` to continue.
