# Opslane event bus

## Mission

Build an asynchronous, non-blocking in-memory publish-subscribe Event Bus for Opslane, implementing topic-based subscriptions, wildcard event matching, handler panic isolation, and context-aware graceful draining.

## Prerequisites

- `opslane-08`

## Mental Model

Modern domain-driven systems require loose coupling between business modules. When an order completes, downstream concerns—such as audit logging, billing updates, analytics counters, and customer notifications—should not run sequentially inside the order processing transaction. The Opslane Event Bus acts as an asynchronous pub-sub dispatcher: modules subscribe to specific topics (e.g. `"order.created"`) or wildcard topics (`"*"`). When an event is published, the bus dispatches it concurrently to all registered listener callbacks in separate goroutines. If a listener panics, the panic is caught and logged, preventing worker crashes. During application shutdown, `Close()` stops accepting new publishes and waits for all active in-flight handlers to drain.

## Visual Model

```text
Publisher (OrderService)
  │ Publish(ctx, orderEvent)
  ▼
┌─────────────────────────────────────────────────────────────┐
│ EventBus Dispatcher                                         │
│   1. Read lock (mu.RLock())                                 │
│   2. Collect topic subscribers + wildcard ("*") subscribers │
│   3. Clone subscriber slice and release lock                │
└──────────────┬───────────────────────────────┬──────────────┘
               │                               │
       go func() [Handler 1]           go func() [Handler 2]
       ├─► wg.Add(1) / defer Done()    ├─► wg.Add(1) / defer Done()
       ├─► defer recover()             ├─► defer recover()
       ▼                               ▼
┌──────────────────────────────┐ ┌────────────────────────────┐
│ Audit Logging Listener       │ │ Email Notification Listener│
└──────────────────────────────┘ └────────────────────────────┘
```

## Machine View

Under the hood:
1. `b.subscribers` maps topic strings to slices of `Subscription`.
2. To avoid holding locks during slow handler execution, `Publish` takes a read lock, shallow-copies the active subscriber slices (`append([]Subscription(nil), ...)`), and immediately unlocks `b.mu.RUnlock()`.
3. Handlers run in isolated goroutines: `b.wg.Add(1)` tracks active execution. Each goroutine establishes a panic recovery boundary with `defer recover()`, preventing unhandled errors in third-party handlers from crashing the entire process.
4. `Close(ctx)` sets `b.closed = true`, spawns a waiter on `b.wg.Wait()`, and uses `select` to return cleanly when all handlers finish or return `ctx.Err()` if shutdown timeout expires.

## Run Instructions

```bash
# Run event bus test suite
go test -v ./curriculum/modules/18-flagship-opslane/assessments/09-opslane-event-bus

# Run standalone pub-sub demo
go run ./curriculum/modules/18-flagship-opslane/assessments/09-opslane-event-bus
```

## Code Walkthrough

1. **Event Abstraction (lines 11–33)**: Defines `Event` interface and `BaseEvent` struct providing `Type`, `TenantID`, `Payload`, and `OccurredAt`.
2. **`EventBus` Struct (lines 42–55)**: Tracks subscribers map, subscriber ID generator, `sync.RWMutex`, `sync.WaitGroup`, and closed flag.
3. **`Subscribe()` & `Unsubscribe()` (lines 58–90)**: Registers handlers returning an integer subscription ID; filters and removes subscriptions cleanly under write lock.
4. **`Publish()` (lines 93–120)**:
   - Evaluates closed state.
   - Copies topic-specific and wildcard (`"*"`) subscribers.
   - Dispatches each handler in a goroutine with `wg.Add(1)` and deferred `recover()`.
5. **`Close()` (lines 123–140)**: Rejects future publishes, waits for in-flight tasks via `wg.Wait()`, and respects context timeout.

## Try It

1. Run the test suite: verify that multiple subscribers receive events concurrently.
2. Subscribe to `"*"`: verify that wildcard listeners capture all published events across all topics.
3. Simulate a panicking handler: verify that other handlers complete successfully and the bus continues operating.

Common mistakes:
- Invoking subscriber callbacks while holding a mutex lock, creating severe lock contention and potential deadlocks if a handler calls `Publish` or `Subscribe`.
- Forgetting panic recovery inside spawned goroutines, allowing any buggy listener to crash the entire application.
- Omitting `sync.WaitGroup`, causing in-flight event handlers to be killed abruptly during application shutdown.

## In Production

In multi-node microservice clusters, in-memory event buses are often augmented with distributed message brokers like Apache Kafka, RabbitMQ, or NATS. To ensure events are not lost if a node crashes before publishing, systems combine the event bus with the **Transactional Outbox Pattern**, writing events to a database table within the same transaction as business state mutations.

## Thinking Questions

1. Why should event handlers never be executed while holding the `EventBus` mutex lock?
2. What are the operational differences between an in-memory event bus and an external broker like Kafka or RabbitMQ?
3. How does the Transactional Outbox pattern prevent lost events during server crashes?

## Next Step

Proceed to [Milestone 10 — Worker Pools](../10-opslane-worker-pools/README.md) to implement bounded concurrency queues and worker pools for heavy background processing.
