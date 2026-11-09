# Rubric — Project Assessment: Opslane Flagship Capstone Defense

## Assessment Overview

This rubric evaluates the flagship capstone implementation and oral defense for graduation from The Go Engineer curriculum as a production-grade Senior Backend Engineer.

---

## 4-Tier Scoring Rubric

| Criterion | Level 1: Inadequate | Level 2: Developing | Level 3: Proficient | Level 4: Exemplary |
|---|---|---|---|---|
| **Core Distributed Architecture** | System crashes under basic load; missing ACID transactional safety; unbounded goroutines. | Implements basic CRUD; lacks outbox pattern; dual-write defects; unable to articulate trade-offs. | Full 14 milestones implemented; passes race detection; outbox pattern guaranteed; bounded worker pool. | Masterful system design: explains WAL flushing, transaction isolation levels, CQRS-lite read optimization, and distributed tracing context extraction. |
| **Multi-Tenant Security** | Cross-tenant data leaks possible; missing authentication or authorization checks. | Basic JWT checks; missing database parameterized tenant enforcement; flat cache keys. | Multi-tenant isolation verified across HTTP, database queries, and Redis key namespacing; passes security test suite. | Defense-in-depth: PostgreSQL RLS session variables combined with application context keys; formal security proof presented during defense. |
| **Resilience & Fault Tolerance** | Cascading failures take down entire service; deadlocks on channel operations; drops traffic on shutdown. | Handles errors but lacks backpressure; worker panics crash the process; naive shutdown. | Panic recovery in worker loops; backpressure on saturated queues; graceful drain sequence on SIGTERM. | Zero-downtime rolling deploys verified; circuit breaker fault containment; automated Prometheus alerting on P99 latency and queue saturation. |
| **Oral Defense & Technical Fluency** | Unable to explain architectural decisions; defers to tutorials; cannot walk through code. | Explains what the code does, but struggles to justify why specific architectural trade-offs were chosen. | Confidently walks panel through end-to-end request lifecycle; justifies database, caching, and concurrency choices. | Authoritative systems leadership: articulates failure modes, capacity planning calculations, and multi-region scaling paths under cross-examination. |

---

## Passing Criteria

- **Minimum Total Score**: 80% (Proficient or Exemplary across all 4 categories).
- **All 14 Milestones Tested**: Verified passing test suites across all 14 Opslane milestones.
