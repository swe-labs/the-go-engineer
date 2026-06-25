# Event-driven basics

## Learning objective

Design and implement an in-memory event bus in Go, use publish/subscribe to decouple producers from consumers, and compose multiple handlers for the same event without coordination.

## Why this matters

When a customer places an order, multiple things must happen: send a confirmation email, reserve inventory, update analytics, notify the fulfillment team. In a tightly coupled system, the order service calls email, inventory, analytics, and fulfillment directly. Adding a new downstream action means changing the order service. Event-driven architecture flips this: the order service publishes an `order.placed` event. Any number of subscribers react to it. The order service does not know about the subscribers, and subscribers do not know about each other.

## Mental model

An event bus is a mailing list. The producer posts a message to the list (publish). Everyone subscribed to the list receives the message (subscribe). The producer does not know who is subscribed. Subscribers do not know about each other. New subscribers can join without the producer's permission. A subscriber can leave without affecting anyone.

In Go, the event bus is typically an in-memory struct with a `Publish` method and a `Subscribe` method. Handlers are functions that match a signature. The bus calls each handler synchronously or asynchronously.

## Core idea

An event is a record that something happened. It has a type (what happened) and data (the relevant payload).

```go
type Event struct {
    Type string
    Data interface{}
}
```

An event bus manages subscriptions and dispatches events:

```go
type EventBus struct {
    handlers map[string][]EventHandler
}

func (b *EventBus) Subscribe(eventType string, handler EventHandler) {
    b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *EventBus) Publish(event Event) {
    for _, h := range b.handlers[event.Type] {
        h(event)
    }
}
```

Publish/subscribe provides three forms of decoupling:

1. **Space decoupling**: the producer and consumer do not know each other's identity or location.
2. **Time decoupling**: the producer and consumer do not need to be active at the same time (if the bus persists events).
3. **Synchronization decoupling**: the producer is not blocked by the consumer's processing (if the bus is asynchronous).

In this lesson, we build an in-memory, synchronous event bus. Synchronous means `Publish` blocks until all handlers complete. This is simpler to reason about and ensures handlers are processed in publish order. Asynchronous buses (using goroutines) add complexity for ordering, error handling, and backpressure.

## Under the hood

The in-memory event bus uses a `sync.RWMutex` to protect the handlers map during concurrent subscribe and publish operations. On publish, the bus reads the handlers slice under a read lock (allowing concurrent publishes), then iterates the slice and calls each handler.

Because handlers are called synchronously in the same goroutine as the publisher, the publisher blocks until all handlers complete. This means:
- Event ordering is preserved: handlers see events in publish order.
- If a handler panics, the publisher sees the panic.
- If a handler blocks indefinitely, the publisher blocks indefinitely.

For these reasons, production event buses add timeouts, retries, and asynchronous dispatch.

## How Go uses it

Go's standard library does not provide an event bus, but the pattern is everywhere in the Go ecosystem:

- `net/http` server sends events like `httptest.Server` close notifications.
- `context.Context` cancellation is a kind of event: when a context is cancelled, all listeners receive the cancellation.
- `database/sql` driver notifications.
- Many Go projects implement in-memory event buses for cross-module communication within a single process.

## Go example

```go
package main

import (
	"fmt"
	"sync"
)

type Event struct {
	Type string
	Data interface{}
}

type EventHandler func(Event)

type EventBus struct {
	mu       sync.RWMutex
	handlers map[string][]EventHandler
}

func NewEventBus() *EventBus {
	return &EventBus{handlers: make(map[string][]EventHandler)}
}

func (b *EventBus) Subscribe(eventType string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *EventBus) Publish(event Event) {
	b.mu.RLock()
	handlers := b.handlers[event.Type]
	b.mu.RUnlock()
	for _, h := range handlers {
		h(event)
	}
}

type OrderPlacedData struct {
	OrderID string
	Total   float64
}

type NotificationService struct {
	notifications []string
	mu            sync.Mutex
}

func (n *NotificationService) OnOrderPlaced(event Event) {
	data := event.Data.(OrderPlacedData)
	n.mu.Lock()
	n.notifications = append(n.notifications,
		fmt.Sprintf("Order %s for $%.2f placed", data.OrderID, data.Total))
	n.mu.Unlock()
}

func (n *NotificationService) Notifications() []string {
	n.mu.Lock()
	defer n.mu.Unlock()
	result := make([]string, len(n.notifications))
	copy(result, n.notifications)
	return result
}

type InventoryService struct {
	reserved []string
	mu       sync.Mutex
}

func (i *InventoryService) OnOrderPlaced(event Event) {
	data := event.Data.(OrderPlacedData)
	i.mu.Lock()
	i.reserved = append(i.reserved,
		fmt.Sprintf("Inventory reserved for %s", data.OrderID))
	i.mu.Unlock()
}

func (i *InventoryService) Reserved() []string {
	i.mu.Lock()
	defer i.mu.Unlock()
	result := make([]string, len(i.reserved))
	copy(result, i.reserved)
	return result
}

func main() {
	bus := NewEventBus()
	notifier := &NotificationService{}
	inventory := &InventoryService{}

	bus.Subscribe("order.placed", notifier.OnOrderPlaced)
	bus.Subscribe("order.placed", inventory.OnOrderPlaced)

	bus.Publish(Event{
		Type: "order.placed",
		Data: OrderPlacedData{OrderID: "ORD-001", Total: 99.99},
	})

	fmt.Println("Notifications:", notifier.Notifications())
	fmt.Println("Inventory:", inventory.Reserved())
}
```

## Step-by-step execution

For `bus.Publish(Event{Type: "order.placed", Data: ...})`:

1. `Publish` acquires a read lock on the handlers map.
2. It looks up the handler slice for key `"order.placed"`. Two handlers are registered: notification and inventory.
3. The read lock is released.
4. The first handler (`NotificationService.OnOrderPlaced`) is called with the event.
5. Inside the handler, the event data is type-asserted to `OrderPlacedData`. The handler formats a notification string and appends it to the notifications slice under a mutex.
6. The second handler (`InventoryService.OnOrderPlaced`) is called with the same event.
7. Inside the handler, the event data is type-asserted. A reservation string is appended to the reserved slice under a mutex.
8. `Publish` returns. Both handlers have processed the same event.

The key property: the `NotificationService` and `InventoryService` have no knowledge of each other. If you add a third subscriber (e.g., `AnalyticsService`), you register it with the bus and the existing code does not change.

## Common mistakes

- **Sharing state between handlers**: handler A sets a field on the event data and handler B reads it. This creates implicit coupling between handlers. Event data should be immutable after publish. If handlers need to share state, use a separate channel or repository.
- **Blocking the publisher**: if a handler makes a slow HTTP call, the publisher blocks until it completes. Use goroutines for slow handlers, but be aware of ordering and error handling.
- **Panic in a handler**: if a handler panics, the publisher sees the panic and may crash. Wrap handler calls in a recover to isolate failures.
- **No event schema**: using `interface{}` for event data means runtime type assertions can panic. Define a typed struct for each event type and use a type switch in handlers.
- **Subscription leaks**: subscribing in a goroutine that terminates without unsubscribing. The handler stays registered forever, creating a memory leak and possibly calling back into a dead goroutine. Provide an `Unsubscribe` method or use subscription tokens.

## Debugging walkthrough

```go
bus.Subscribe("order.placed", func(event Event) {
    data := event.Data.(OrderPlacedData)
    fmt.Println(data.OrderID)
})
bus.Publish(Event{Type: "order.placed", Data: "wrong type"})
```

**Symptom**: Panic: `interface conversion: interface {} is string, not OrderPlacedData`.

**Root cause**: The publisher used a string as event data, but the handler expected `OrderPlacedData`. The type assertion panics on mismatch.

**Fix (publisher side)**: Publish correctly typed data:
```go
bus.Publish(Event{Type: "order.placed", Data: OrderPlacedData{OrderID: "ORD-001"}})
```

**Fix (handler side)**: Use the comma-ok idiom to handle unexpected types gracefully:
```go
data, ok := event.Data.(OrderPlacedData)
if !ok {
    log.Printf("unexpected event data type: %T", event.Data)
    return
}
```

## Production notes

- **Event versioning**: event schemas evolve. An `order.placed` event might add a `discount` field. Old subscribers that expect the old schema may break. Use protobuf or Avro for schema evolution, or version the event type (`order.placed.v2`).
- **At-least-once delivery**: the in-memory bus does not provide delivery guarantees. If the process crashes after publish but before a handler runs, the event is lost. For durable events, use a message broker (RabbitMQ, Kafka, AWS SQS).
- **Dead-letter handling**: if a handler repeatedly fails (e.g., database is down), the event should be moved to a dead-letter queue after N retries. The in-memory bus does not provide this; a production event bus should.
- **Observability**: instrument the bus with metrics (events published, handler duration, handler errors) and structured logging. This lets you detect slow handlers, failing handlers, and event throughput.
- **Synchronous vs asynchronous**: the synchronous bus is ideal for in-process events where the publisher needs the side effects to complete before responding (e.g., sending a confirmation email before returning the HTTP response). For fire-and-forget events, use async dispatch with a goroutine pool.

## Performance implications

- **Handler call overhead**: each handler is a function call. With 10 subscribers per event type and 1000 events/second, that is 10,000 handler calls/second. Go handles this easily.
- **Type assertion overhead**: asserting `event.Data.(OrderPlacedData)` is fast (a few nanoseconds). Use the comma-ok form to avoid panics.
- **Lock contention**: the mutex protects the handler map on subscribe and publish. Subscriptions are rare (setup time). Publishes are common (runtime). Use `sync.RWMutex` so multiple publishes do not block each other.
- **Goroutine overhead**: if you add goroutines for async handlers, each goroutine costs ~4KB of stack. For thousands of concurrent events, use a worker pool.

## Practice task

Add a third subscriber to the event bus: a `ReportingService` that counts the total value of all placed orders. It should maintain a running total and expose a `TotalRevenue() float64` method. Then add an `order.refunded` event type with a `RefundOrder` handler in the `ReportingService` that subtracts the refunded amount. Write table-driven tests for both event types with multiple subscribers.

## Tests / verification

```bash
go run ./curriculum/modules/14-architecture-distributed-systems/lessons/09-event-driven-basics
go test ./curriculum/modules/14-architecture-distributed-systems/lessons/09-event-driven-basics
```

The tests verify that events are delivered to all subscribers, that unrelated event types do not trigger handlers, and that each subscriber processes events independently. After completing the practice task, add tests for the reporting subscriber and the refund event flow.

## Review questions

1. What are the three forms of decoupling that event-driven architecture provides?
2. In the synchronous in-memory bus, what happens if a handler blocks for 10 seconds?
3. How would you add error handling to the event bus so that a failing handler does not prevent other handlers from running?
4. What is the risk of using `interface{}` for event data, and how do you mitigate it?
5. When would you choose an in-process event bus over a message broker like RabbitMQ or Kafka?

## NEXT UP

Queues -- durable, asynchronous work distribution with message brokers.
