# The Go Engineer v3 — Graduation Outcomes & Capstone Defense Rubric

> **Target Standard:** Produce an autonomous, capable software engineer who can build, explain, test, deploy, monitor, and recover Go backend software.
> **Status:** Active Curriculum Standard (`v3.0.0`)
> **Scope:** Evaluates mastery across all 18 mandatory core modules (Modules 00–16 and Module 18 Flagship Opslane). Module 17 (Advanced Electives) provides optional specialization tracks and is not required for core graduation.

---

## 1. The 6 Core Graduation Outcomes

Graduation from *The Go Engineer* is determined by observable engineering behavior, not passive lesson completion. A graduating engineer must independently prove competence across six core dimensions:

### Outcome 1: Unassisted Problem Decomposition
- **Capability:** Given an unfamiliar requirement or ambiguous specification, decompose it into clear input/output contracts, boundary constraints, and loop invariants.
- **Evidence:** Completing algorithmic tasks (custom multi-key sorting via `slices.SortFunc`, binary search tier lookups, order-preserving deduplication, two-pointer partitioning) and business domain logic without copying prior solutions or relying on boilerplate generators.

### Outcome 2: Autonomous Debugging & Execution Tracing
- **Capability:** Diagnose panics, nil pointer dereferences, data races, deadlocks, and goroutine leaks using Delve, stack traces, and the Go race detector (`-race`).
- **Evidence:** Tracing through execution state tables, setting conditional breakpoints, analyzing compiler escape decisions, and identifying the exact hardware or runtime root cause before applying a minimal, verified regression fix.

### Outcome 3: Automated Verification & Test Rigor
- **Capability:** Write robust, deterministic, table-driven unit and integration test suites that prove correctness under both success and failure conditions.
- **Evidence:** Testing boundary values ($N=0, 1, \infty$), validating error returns with `errors.Is`/`errors.As`, structuring subtests with `t.Run()`, avoiding non-deterministic time-dependent assertions, and asserting transactional rollback on simulated database failures.

### Outcome 4: End-to-End Delivery & Disaster Recovery Proof
- **Capability:** Package services into minimal, secure multi-stage Docker containers, configure CI/CD pipelines, and execute verified database disaster recovery drills.
- **Evidence:** Performing live container destruction and proving reproducible recovery: backup artifact generation (`pg_dump`) $\to$ volume and container purge $\to$ automated restore (`pg_restore`) $\to$ automated SQL row count and sequence integrity verification. The drill must satisfy defined **Recovery Point Objective (RPO: zero loss of committed transactions prior to backup timestamp)** and **Recovery Time Objective (RTO: complete restore and service health check passed in $\le 5$ minutes)**.

### Outcome 5: Code Review & Architectural Tradeoff Defense
- **Capability:** Articulate *why* a particular design was chosen, analyze time/space Big-O and memory allocation costs, and defend Architecture Decision Records (ADRs).
- **Evidence:** Providing actionable, constructive peer reviews, choosing between generics and interfaces using explicit engineering criteria, and conducting a formal oral and technical defense of Flagship Opslane.

### Outcome 6: Live Operational Incident Response
- **Capability:** Instrument services with structured logging (`slog`), Prometheus metrics, and OpenTelemetry spans, and respond deterministically to production failures.
- **Evidence:** Triage simulated operational failures: deadlocks, connection pool starvation, memory exhaustion, token tampering, and cross-tenant leakage attempts with deterministic 401/403/404 isolation.

---

## 2. Flagship Opslane Capstone Defense Rubric

The capstone defense evaluates the learner's complete implementation of **Flagship Opslane** (Module 18) across 5 observable competency levels:

| Score | Rating | Meaning | Defense Criteria |
| :---: | :--- | :--- | :--- |
| **5** | **Exemplary** | Autonomous Technical Mastery | Demonstrates complete technical mastery; explains runtime mechanics when relevant to debugging; zero data races under `-race`; deterministic failure handling; reproducible disaster recovery drill satisfying defined RPO/RTO. |
| **4** | **Proficient** | Core Engineering Competence | Meets all mandatory specification requirements; clean idiomatic Go; deterministic unit and integration tests; robust error handling; clearly defends architectural decisions and operational trade-offs. |
| **3** | **Competent** | Foundational Competence | Implements core behavior; tests cover major happy paths; requires minor guidance during technical defense of distributed tradeoffs or failure modes. |
| **2** | **Novice** | Developing Competence | Incomplete specification requirements; brittle error handling; struggles to explain race conditions, context lifecycle, or tenant isolation boundaries. |
| **1** | **Unprepared**| Incomplete | Code non-functional or copied without comprehension; missing tests; service panics or leaks resources under failure scenarios. |

> **Evaluation Disclaimer:** Meeting the curriculum graduation standard certifies demonstrated technical capability against the observable engineering benchmarks above. It confirms readiness to contribute autonomously to production Go codebases within defined operational bounds; it does not claim universal hiring guarantees across unrelated tech stacks or unverified operational environments.

---

## 3. Capstone Defense Protocol (Oral & Technical Defense)

1. **Architecture Walkthrough (15 min):** Explain the multi-tenant architecture, database schema, transaction boundaries, and event pipeline.
2. **Failure Simulation (15 min):** The reviewer injects a simulated fault (e.g., terminate database container, simulate network timeout, submit unauthorized cross-tenant requests). The learner must diagnose the symptom, explain runtime behavior, and demonstrate recovery.
3. **Code Quality & Concurrency Audit (15 min):** Review worker pool implementation under the Go race detector (`go test -race ./...`), inspecting channel ownership, mutex scopes, and context cancellation.
4. **Graduation Sign-Off:** Requires a minimum score of **4/5** across all core dimensions.
