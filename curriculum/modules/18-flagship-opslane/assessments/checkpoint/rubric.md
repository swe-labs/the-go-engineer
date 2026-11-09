# Rubric — Module 18 Checkpoint: Flagship Opslane Architecture & Capstone Defense

## Assessment Overview

This checkpoint evaluates complete end-to-end mastery across all 15 instructional nodes of Flagship Opslane (`opslane-01` through `opslane-15`), measuring architectural depth, multi-tenant security, concurrency safety, distributed tracing, and production resilience.

---

## 4-Tier Scoring Rubric

| Criterion | Level 1: Inadequate | Level 2: Developing | Level 3: Proficient | Level 4: Exemplary |
|---|---|---|---|---|
| **Multi-Tenant Security & Isolation** | Insecure design; missing tenant isolation; cross-tenant query leaks; unauthenticated API routes. | Checks tenant ID in HTTP handler but omits database enforcement; vulnerable to IDOR; flat cache keys. | Complete tenant isolation: JWT claims extraction, type-safe context propagation, parameterized SQL with tenant filters, scoped cache keys. | Defense-in-depth: combines application context scoping with PostgreSQL Row-Level Security (RLS) session variables; automated cross-tenant security test suite. |
| **Event-Driven Architecture & Outbox** | Synchronous coupling of external APIs inside DB transactions; lost events on crashes; unbounded goroutines. | Implements basic in-memory channel event bus; prone to event loss on crash or restart; naive dual-writes. | Transactional outbox pattern implemented; at-least-once delivery; bounded worker pools with backpressure; clean panic recovery. | Enterprise distributed systems design: idempotency keys with deduplication windows; dead letter queues; outbox CDC; backpressure telemetry. |
| **Caching & Performance Optimization** | No caching or naive cache with stale reads; vulnerable to cache stampedes; crashes on Redis outage. | Basic Redis GET/SET; misses cache invalidation on updates; unhandled connection failures. | Cache-aside with Redis; TTL expiration; proactive invalidation on mutations; fallback to database on cache disconnection. | Thundering herd mitigation via singleflight coalescing; multi-level caching (local sharded cache + Redis); P99 < 2ms latency SLA. |
| **Observability & Graceful Shutdown** | No logs or unstructured logs; unhandled SIGTERM drops requests; leaks DB connections on restart. | Basic log output; handles SIGINT but terminates workers abruptly; no trace propagation. | Structured slog with correlation IDs; OpenTelemetry context propagation across event bus; graceful drain on SIGTERM. | Zero-downtime rolling deploys verified; automated Prometheus alerts on P99 latency and queue depth; detailed incident postmortem analysis. |

---

## Passing Criteria

- **Minimum Total Score**: 80% (Proficient or Exemplary across all categories).
- **All 14 Milestones Tested**: Verified passing test suites across all 14 Opslane milestones.
