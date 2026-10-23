# The Go Engineer v3 — Graduation Outcomes & Capstone Defense Rubric

> **Target Standard:** Produce an autonomous, capable software engineer who can build, explain, test, deploy, monitor, and recover production Go software.  
> **Status:** Active Release Candidate Standard (`v3.1.0-rc.1`)

---

## 1. The 6 Core Graduation Outcomes

Graduation from *The Go Engineer* is determined by observable engineering behavior, not passive lesson completion. A graduating engineer must independently prove competence across six core dimensions:

### Outcome 1: Unassisted Problem Decomposition
- **Capability:** Given an unfamiliar requirement or ambiguous specification, decompose it into input/output contracts, boundary constraints, and loop invariants.
- **Evidence:** Completing algorithmic tasks (custom multi-key sorting via `slices.SortFunc`, binary search tier lookups, order-preserving deduplication, two-pointer partitioning) and business logic without copying prior solutions or relying on boilerplate generators.

### Outcome 2: Autonomous Debugging & Execution Tracing
- **Capability:** Diagnose panics, nil pointer dereferences, data races, deadlocks, and goroutine leaks using Delve, stack traces, and the Go race detector (`-race`).
- **Evidence:** Tracing through execution state tables, setting conditional breakpoints, inspecting compiler memory escape, and explaining the exact hardware/runtime root cause before applying a minimal, verified fix.

### Outcome 3: Automated Verification & Test Rigor
- **Capability:** Write robust, deterministic, table-driven unit and integration test suites that prove correctness under both success and failure conditions.
- **Evidence:** Testing boundary values ($N=0, 1, \infty$), validating error returns with `errors.Is`/`errors.As`, using subtests with `t.Run()`, avoiding flaky time-dependent assertions, and asserting transactional rollback on database errors.

### Outcome 4: End-to-End Delivery & Disaster Recovery Proof
- **Capability:** Package services into minimal, secure multi-stage Docker containers, configure CI/CD pipelines, and execute verified database disaster recovery drills.
- **Evidence:** Performing live container destruction and proving zero data loss through `pg_dump` $\to$ container purge $\to$ `pg_restore` $\to$ automated SQL row count and sequence integrity verification.

### Outcome 5: Code Review & Architectural Tradeoff Defense
- **Capability:** Articulate *why* a particular design was chosen, analyze time/space Big-O and memory allocation costs, and defend Architecture Decision Records (ADRs).
- **Evidence:** Providing actionable, constructive peer reviews, choosing between generics and interfaces with explicit criteria, and conducting a formal oral/technical defense of Flagship Opslane.

### Outcome 6: Live Operational Incident Response
- **Capability:** Instrument services with structured logging (`slog`), Prometheus metrics, and OpenTelemetry spans, and respond deterministically to production failures.
- **Evidence:** Triage simulated operational failures: deadlocks, connection pool starvation, memory exhaustion, token tampering, and cross-tenant leakage attempts with deterministic 401/403/404 isolation.

---

## 2. Flagship Opslane Capstone Defense Rubric

The capstone defense evaluates the learner's complete implementation of **Flagship Opslane** (Module 18) across 5 rigorous levels:

| Score | Rating | Meaning | Defense Criteria |
| :---: | :--- | :--- | :--- |
| **5** | **Exemplary** | Production Ready | Demonstrates deep mastery; explains runtime mechanics (scheduler, GC, memory layout); handles edge cases; zero race conditions; full DR drill verified. |
| **4** | **Proficient** | Industry Hireable | Meets all requirements; clean idiomatic Go; comprehensive tests; robust error handling; defends design decisions clearly. |
| **3** | **Competent** | Junior Standard | Implements core behavior; tests cover major paths; requires minor guidance during technical defense of distributed tradeoffs. |
| **2** | **Novice** | Developing | Incomplete requirements; toy error handling; struggles to explain race conditions or tenant isolation boundaries. |
| **1** | **Unprepared**| Failing | Code copied without comprehension; tests missing; service panics under failure scenarios. |

---

## 3. Capstone Defense Protocol (Oral & Technical Defense)

1. **Architecture Walkthrough (15 min):** Explain the multi-tenant architecture, database schema, transaction boundaries, and event pipeline.
2. **Failure Simulation (15 min):** The reviewer injects a simulated fault (e.g., kill Postgres container, simulate network partition, submit unauthorized cross-tenant requests). The learner must diagnose the symptom, explain runtime behavior, and demonstrate recovery.
3. **Code Quality & Concurrency Audit (15 min):** Review worker pool implementation under the Go race detector (`go test -race ./...`), inspecting channel ownership and context cancellation.
4. **Graduation Sign-Off:** Requires a minimum score of **4/5** across all core dimensions.
