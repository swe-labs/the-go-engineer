# CQRS

## When you need this

CRUD-based architectures use the same data model for reads and writes, causing contention when read-heavy and write-heavy workloads compete for the same tables. Complex queries require joins across multiple tables, slowing down read performance. CQRS solves this by separating read and write models into independent paths, each optimized for its specific workload.

## When you do not need this

When simpler concrete types or standard core library primitives suffice without adding abstraction overhead.

## Prerequisites

- elective-15

## Mental model

CQRS splits the system into two sides. The command side receives write operations, validates them, and appends events to an event store. The query side reads from denormalized read models that are optimized for specific queries. Events propagate from the command side to the query side via an event bus. The read model is eventually consistent with the write model — there is a propagation delay.

## How it works

Commands are validated, applied to the domain model, and persisted as events in an event store (append-only log). The event store is the source of truth. A projection (event handler) subscribes to events and updates a read model (typically a denormalized database table or a separate read-optimized store). The read model is rebuilt by replaying all events from the beginning. Read models can be specialized: one for order listing (denormalized with customer name and total), another for order statistics (aggregated counts).

## Go example

```go
package main

import "fmt"

func main() {
    fmt.Println("Elective: CQRS")
}
```

## Tradeoffs

- CQRS improves read performance: read models are optimized for specific queries (denormalized, indexed) — single-table lookups instead of joins
- CQRS improves write scalability: command side can be scaled independently from the query side
- CQRS adds latency: events must propagate from command store to read model — propagation delay is typically 100ms-1s
- Storage cost: data is stored multiple times (command store + read models + event store) — 2-5x more storage than normalized single-database approach

## Common mistakes

- Applying CQRS to CRUD apps — CQRS adds complexity (event bus, separate stores, eventual consistency) that is not justified for simple CRUD.
- Sharing the same database between command and query sides — the database is a coupling point; use separate databases or schemas.
- Ignoring eventual consistency — users see stale data after a write; communicate the propagation delay and use UI patterns (optimistic updates).
- Using synchronous event propagation — the command handler waits for the read model update before returning, eliminating the benefits of CQRS.
- Not having a rebuild mechanism — read models can become corrupt or out of sync; provide a way to rebuild them from the event store.

## Practice task

The learner must design a CQRS system for an order management service with separate command and query models, implement an event store that persists commands as events, create a read model projection that consumes events and updates a denormalized orders_read table, and demonstrate that a write is visible on the read side only after event propagation completes (eventual consistency).

## Production notes

CQRS is used in high-scale systems where read and write workloads diverge: e-commerce (order writes are transactional, product catalog reads are denormalized and cached), banking (transaction writes require strict validation, balance queries need fast lookups), and analytics (write events to a log, read aggregated projections). Go services use CQRS with NATS for event propagation and PostgreSQL for read models.

## Related core lessons

- Module 16.15 — sync.Pool
- Module 17.14 — sync.Map
