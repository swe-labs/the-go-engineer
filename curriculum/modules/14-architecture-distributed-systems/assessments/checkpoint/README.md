# Module 14 Assessment — Architecture and Distributed Systems

## Purpose

Validate mastery for the target curriculum items.

## Targets

- `core-14-01`
- `core-14-02`
- `core-14-03`
- `core-14-04`
- `core-14-05`
- `core-14-06`
- `core-14-07`
- `core-14-08`
- `core-14-09`
- `core-14-10`
- `core-14-11`
- `core-14-12`
- `core-14-13`
- `core-14-14`
- `core-14-15`
- `core-14-16`
- `core-14-17`

## Evidence

Submit answers, code evidence, test output, and explanation where required.

---

## 1. Deep-Dive Distributed Systems & Architecture Scenarios

Answer each scenario in detail, providing architectural diagrams, Go interface definitions, and trade-off analyses where relevant.

### Scenario 1: Clean Architecture vs. Monolithic Sprawl (`core-14-01` through `core-14-04`)
A team maintains a 50,000-line Go API service. Over time, HTTP handlers have begun executing direct SQL queries via global `sql.DB` variables and calling third-party payment gateways inline. Adding a new payment provider or testing business rules now takes weeks.
- Explain the principles of Hexagonal Architecture (Ports and Adapters). How does it invert infrastructure dependencies?
- Why does idiomatic Go favor small, consumer-defined interfaces (`io.Reader`, single-method ports) over broad producer-defined interface hierarchies?
- Contrast Constructor Dependency Injection with Global Service Locators. How does constructor injection impact race detection and test parallelization (`t.Parallel()`)?

### Scenario 2: Synchronous RPC vs. Asynchronous Event-Driven Architecture (`core-14-05` through `core-14-08`)
In an e-commerce platform, the Checkout Service calls the Inventory Service, Fraud Service, Payment Service, and Notification Service sequentially over synchronous HTTP REST calls. During Black Friday sales, p99 checkout latency skyrockets to 18 seconds and multiple requests fail with 504 Gateway Timeout.
- Explain the concept of **Temporal Coupling** in synchronous microservice architectures.
- How would you refactor this checkout workflow using event-driven choreography or orchestration with an event broker (e.g., Kafka or RabbitMQ)?
- How do you guarantee the **Transactional Outbox Pattern** in Go to ensure database state changes and outbound event publishing are atomic without distributed 2PC transactions?

### Scenario 3: Consensus, State Machines, and Raft (`core-14-09` through `core-14-12`)
A distributed key-value store relies on three Go nodes to manage cluster state and leadership elections.
- Explain the three roles in the Raft consensus algorithm: Leader, Follower, and Candidate.
- How does Raft handle split votes and network partitions (split-brain)?
- What is a quorum ($Q = \lfloor N/2 \rfloor + 1$), and why does a 3-node cluster tolerate only 1 failure, whereas a 5-node cluster tolerates 2 failures?
- How does log compaction (snapshotting) prevent Raft write-ahead logs from consuming unbounded disk space?

### Scenario 4: Resiliency Patterns — Circuit Breakers, Retries, and Backoff (`core-14-13` through `core-14-17`)
A Go downstream service experiences intermittent network packet loss, causing upstream clients to retry requests aggressively. This avalanche of retries overwhelms the struggling service, causing complete outage (Retry Storm).
- Explain the three states of a **Circuit Breaker** (Closed, Open, Half-Open). How does it protect downstream dependencies?
- Explain the difference between naive retries, exponential backoff, and **Full Jitter**. Why is randomized jitter essential to prevent synchronized retry spikes?
- What are **Bulkheads** and how do bounded worker channels in Go prevent one failing downstream integration from consuming all server worker goroutines?

---

## 2. Practical Architecture Challenge

Design and implement the Go interface contracts and service orchestration for a distributed resilient order processor:
1. Define pure domain interfaces: `PaymentPort`, `InventoryPort`, `OutboxRepository`.
2. Implement an `OrderOrchestrator` that:
   - Evaluates payment with a Circuit Breaker wrapper.
   - If payment succeeds, writes the order and an outbox event within a database transaction.
   - If payment fails with a transient error, applies exponential backoff with jitter before failing.
   - Propagates `context.Context` with timeout across all operations.

---

## 3. Invariant Architectural Principles

1. **Dependencies Point Inward**: Core business domain logic must never import outer infrastructure packages (database drivers, HTTP routers, cloud SDKs).
2. **Never Retry Non-Idempotent Mutations**: Never automatically retry non-idempotent operations without an idempotency key header (`Idempotency-Key`).
3. **Fail Fast with Circuit Breakers**: Downstream dependencies known to be offline must fail immediately in the calling service without exhausting connection pools.
4. **Zero Shared Mutable State Across Goroutines**: Isolate mutable state inside single-goroutine actors or guard via explicit port interfaces.

---

## 4. 4-Tier Assessment Rubric

| Level | Criteria & Observable Evidence |
|---|---|
| **Level 1: Novice** | Confuses clean architecture with directory organization; uses global DB singletons; implements unbounded retry loops causing retry storms; cannot explain Raft split-brain. |
| **Level 2: Advanced Beginner** | Understands ports and adapters conceptually but defines interfaces in the infrastructure layer; uses retries without backoff or jitter; explains circuit breakers superficially without state transitions. |
| **Level 3: Competent** | Designs proper hexagonal boundaries; defines consumer-oriented port interfaces; implements exponential backoff with full jitter and circuit breakers; explains Raft consensus, quorums, and outbox patterns accurately. |
| **Level 4: Proficient** | Exceptional distributed systems architect. Articulates consensus trade-offs (CAP theorem, PACELC); designs fault-tolerant systems using transactional outbox and event sourcing; implements bulkhead isolation in Go concurrency pipelines. |

---

## 5. Submission Requirements

Submit your architectural breakdown addressing Scenarios 1–4 and the interface and orchestration code for the practical challenge.
