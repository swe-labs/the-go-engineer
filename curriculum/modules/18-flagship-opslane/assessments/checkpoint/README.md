# Module 18 Assessment — Flagship Opslane Architecture & Capstone Defense

## Purpose

Validate mastery for the target curriculum items.

## Targets

- `opslane-01`
- `opslane-02`
- `opslane-03`
- `opslane-04`
- `opslane-05`
- `opslane-06`
- `opslane-07`
- `opslane-08`
- `opslane-09`
- `opslane-10`
- `opslane-11`
- `opslane-12`
- `opslane-13`
- `opslane-14`
- `opslane-15`

## Evidence

Submit answers, code evidence, test output, and explanation where required.

---

## 1. Visual Model: Complete Opslane Flagship Distributed Architecture

```
+-----------------------------------------------------------------------------------------+
|                                OPSLANE DISTRIBUTED SYSTEM                               |
+-----------------------------------------------------------------------------------------+
|                                                                                         |
|      Ingress Traffic                                                                    |
|             |                                                                           |
|             v                                                                           |
|   +---------------------------------------------------------------------------------+   |
|   |                       HTTP API Layer (Go 1.22 Router)                           |   |
|   |  - Middleware: Correlation ID, Structured Slog, Prometheus Metrics, OTel Spans  |   |
|   |  - Tenant Context Extraction -> Multi-Tenant Isolation Enforcement             |   |
|   |  - Auth: JWT Claims Verification & Role-Based Access Control (RBAC)             |   |
|   +---------------------------------------------------------------------------------+   |
|             |                                                  |                        |
|             v                                                  v                        |
|   +------------------------------------+    +---------------------------------------+   |
|   |       Order & Payment Engine       |    |         Caching & State Layer         |   |
|   |  - State Machine Validation        |    |  - Redis Cache Aside (P99 < 2ms)      |   |
|   |  - Database ACID Transactions      |    |  - Sharded Local Fallback             |   |
|   |  - Idempotency Key De-duplication  |    +---------------------------------------+   |
|   +------------------------------------+                       |                        |
|             |                                                  |                        |
|             v                                                  v                        |
|   +---------------------------------------------------------------------------------+   |
|   |                    PostgreSQL Relational Persistence Engine                     |   |
|   |  - Multi-tenant foreign keys, pgxpool connection limits, schema migrations       |   |
|   |  - Transactional Outbox Pattern for atomic event persistence                    |   |
|   +---------------------------------------------------------------------------------+   |
|             |                                                                           |
|             v                                                                           |
|   +---------------------------------------------------------------------------------+   |
|   |                       Internal Asynchronous Event Bus                           |   |
|   |  - Publisher/Subscriber fan-out, non-blocking delivery, consumer backpressure   |   |
|   +---------------------------------------------------------------------------------+   |
|             |                                                                           |
|             v                                                                           |
|   +---------------------------------------------------------------------------------+   |
|   |                         Bounded Concurrent Worker Pool                          |   |
|   |  - Panic isolation, fixed worker goroutines, graceful drain on SIGTERM           |   |
|   +---------------------------------------------------------------------------------+   |
|                                                                                         |
+-----------------------------------------------------------------------------------------+
```

---

## 2. Core Architectural Scenarios

### Scenario 1: Multi-Tenant Data Isolation & Query Security (`opslane-01` to `opslane-05`)
- How does the Opslane HTTP middleware inject tenant identity into the request `context.Context`?
- What prevents SQL queries from accidentally omitting the `WHERE tenant_id = $1` clause? Explain row-level security (RLS) vs repository layer contract enforcement.
- How do cryptographic token claims guarantee that User A in Tenant X cannot view or modify orders in Tenant Y?

### Scenario 2: Asynchronous Event Bus & Bounded Worker Pools (`opslane-07` to `opslane-10`)
- Walk through the lifecycle of an order: HTTP Request -> Database Insert -> Event Emission -> Asynchronous Notification.
- Why is emitting an event directly inside a database transaction an anti-pattern? (Explain the dual-write problem and transactional outbox pattern).
- How does the bounded worker pool throttle background processing when event publishing spikes by 50x?

### Scenario 3: Caching Invalidation & Cache-Aside Resilience (`opslane-11`)
- Explain the Opslane cache-aside pattern with Redis.
- How does Opslane handle cache stampedes (thundering herd problem) when a hot product record expires?
- If the Redis cluster suffers a network partition, does Opslane fail requests or fall back to PostgreSQL? How is database pool exhaustion prevented?

### Scenario 4: Telemetry, Observability, and Zero-Downtime Deployment (`opslane-12` to `opslane-14`)
- How does a single correlation ID propagate from the HTTP boundary down to SQL queries and background worker goroutines?
- Detail the exact graceful shutdown sequence when Kubernetes sends `SIGTERM`: listener shutdown, drain wait delay, worker queue drain, database pool closure.

---

## 3. Invariant Production Rules

1. **Strict Multi-Tenant Boundary**: All database reads, writes, and cache keys MUST incorporate the tenant identifier.
2. **Never Publish Events Before Commit**: Asynchronous events must only trigger after the local database transaction is committed to disk.
3. **Bounded Memory and Goroutines**: No task creation may spawn an unbounded `go func()`; all background processing must route through the bounded worker pool.
4. **Zero-Downtime Signal Handling**: Process termination MUST drain existing in-flight connections with context timeout before calling `os.Exit(0)`.

---

## 4. 4-Tier Assessment Rubric

| Level | Criteria & Observable Evidence |
|---|---|
| **Level 1: Inadequate** | Missing multi-tenant isolation; event bus causes data races or leaks goroutines; crashes on SIGTERM; no tests. |
| **Level 2: Developing** | Implements basic CRUD; lacks tenant isolation enforcement; emits events before DB commit; unbuffered worker channels. |
| **Level 3: Proficient** | All 14 milestones pass automated tests; multi-tenant RLS/context isolation verified; transactional outbox pattern used; structured slog + OTel tracing; graceful drain on SIGTERM. |
| **Level 4: Exemplary** | Exceptional distributed systems architecture. Explains cache-line bouncing, database connection pool queuing theory, thundering herd mitigation, and distributed tracing context propagation; demonstrates oral capstone defense. |

---

## 5. Submission Requirements

Submit your answers to the 10 questions in `questions.md` demonstrating comprehensive mastery of the entire Opslane distributed platform.
