# Deeper distributed systems

## When you need this

Distributed systems fail in ways that single-process systems do not: network partitions, partial failures, clock skew, and cascading failures. Go developers building distributed services need to understand these failure modes and the patterns that mitigate them — consensus, consistency models, circuit breakers, and idempotency.

## When you do not need this

When simpler concrete types or standard core library primitives suffice without adding abstraction overhead.

## Prerequisites

- elective-19

## Mental model

A distributed system is a collection of independent nodes that communicate over an unreliable network. Each node can fail independently, and the network can partition (some nodes cannot reach others). The fundamental challenge is making progress despite partial failures. Patterns: timeouts detect failures, retries handle transient errors, circuit breakers prevent cascading failures, bulkheads isolate failures, and idempotency ensures safe retries.

## How it works

The network is not reliable: packets can be lost, delayed, duplicated, or reordered. Timeouts detect failures but cannot distinguish between 'node crashed' and 'node is slow but will respond eventually'. Consensus protocols (Raft) use a quorum of nodes (majority) to agree on state changes — a node cannot commit a change without approval from the majority. The CAP theorem states that in a partition, you must choose between consistency (all nodes see the same data) and availability (every request gets a response).

## Go example

```go
package main

import "fmt"

func main() {
    fmt.Println("Elective: Deeper distributed systems")
}
```

## Tradeoffs

- Consensus (Raft, etcd) requires fsync on every write — ~5-20ms per write for durability, limiting throughput to ~1000 writes/s per leader
- Eventual consistency allows local writes without coordination — throughput is limited only by the local node, ~10,000-100,000 writes/s
- Distributed transactions with 2PC: ~5-50ms per transaction, coordination overhead increases with participants
- Retry amplification: a single failed request can trigger 3 retries × 5 downstream services = 15 additional requests — exponential backoff is essential

## Common mistakes

- Assuming the network is reliable — every network call can fail; always have retry logic with exponential backoff and jitter.
- Not setting timeouts — a network call without a timeout blocks the goroutine indefinitely, leaking memory and connections.
- Using synchronous calls in a chain — a failure in any service cascades to all upstream services; use async patterns or circuit breakers.
- Ignoring clock skew — distributed timestamps are unreliable; use monotonic clocks for duration and logical clocks for ordering.
- Not testing for network failures — a system that has never been tested with network partitions will fail when a partition occurs.

## Practice task

The learner must explain the CAP tradeoff for a distributed key-value store, implement a circuit breaker that opens after 3 failures and half-opens after 30 seconds, and write a retry function with exponential backoff and jitter that calls an external service and handles transient failures.

## Production notes

Distributed systems patterns are used in every multi-service architecture: etcd provides distributed consensus for Kubernetes, Cassandra uses eventual consistency for high availability, Kafka uses distributed logs for event streaming, and Redis Sentinel uses failure detection for high availability. Go services use circuit breakers (gobreaker), retries (go-resiliency), and rate limiters (golang.org/x/time/rate) for resilience.

## Related core lessons

- Module 16.19 — Microservice splitting
- Module 17.18 — Blue/green deployments
