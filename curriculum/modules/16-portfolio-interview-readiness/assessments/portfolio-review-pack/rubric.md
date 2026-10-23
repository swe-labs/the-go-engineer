# Rubric — Portfolio Review Pack & Interview Defense

## Assessment Overview

This assessment measures portfolio completeness, technical defense fluency, and operational depth across the 8 showcase projects of the Go Engineer curriculum.

---

## 4-Tier Scoring Rubric

| Criterion | Level 1: Inadequate | Level 2: Developing | Level 3: Proficient | Level 4: Exemplary |
|---|---|---|---|---|
| **Portfolio Audit Completeness** | Code fails to compile or fails basic unit tests; missing multiple showcase projects. | Basic tests pass; only checks file presence without verifying quality, artifacts, or runnable benchmarks. | All 8 showcase projects audited with 100% score; verifies README contracts, Go source, tests, and proof artifacts. | Thorough automated auditor with rich markdown scorecards, detailed error diagnostics, and automated verification of compile/test gates. |
| **System Architecture Defense** | Cannot explain architectural trade-offs; defaults to "that's how the tutorial wrote it." | Explains basic design patterns (MVC, handlers) but struggles to justify database consistency, isolation levels, or dependency direction. | Clearly articulates hexagonal architecture, ACID guarantees vs eventual consistency, and trade-offs of microservices vs modular monoliths. | Authoritative defense: compares storage engine internals, outbox pattern scalability, network partition behavior, and trade-offs under high write concurrency. |
| **Concurrency & Failure Recovery** | Unaware of goroutine leaks, data races, or unhandled panics crashing the process. | Implements basic mutexes or channels but leaves worker goroutines vulnerable to panics and abandoned context channels. | Demonstrates robust panic recovery, bounded worker pools, context propagation, and clean drain semantics under race detector. | Deep runtime understanding: explains scheduler M:N mechanics, lock striping cache locality, zero-allocation pooling, and backpressure telemetry. |
| **Production Deployment & Observability** | Builds huge root images; ignores CVEs; terminates abruptly without draining connections. | Uses basic Dockerfiles; logs unstructured strings; configures naive liveness probes checking database connections. | Implements multi-stage static scratch builds; non-root UID 10001; structured slog + OTel tracing; graceful drain on SIGTERM. | Implements zero-downtime rolling deploys with network drain delays, cgroup CFS quota tuning (`automaxprocs`), and automated CI vulnerability gating. |

---

## Passing Criteria

- **Minimum Score**: 80% total score (Level 3 or higher across all 4 criteria).
- **All 8 Showcase Projects Verified**: All 8 projects in Modules 08–15 must have clean, passing unit and integration tests.
