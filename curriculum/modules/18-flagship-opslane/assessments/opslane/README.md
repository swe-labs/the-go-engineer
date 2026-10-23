# Project Assessment — Opslane Flagship Capstone Defense

## Purpose

Validate mastery for the target curriculum items.

## Targets

- `project-opslane`

## Evidence

Submit answers, code evidence, test output, and explanation where required.

---

## 1. Visual Model: Capstone Defense Evaluation Pipeline

```
+-----------------------------------------------------------------------------------------+
|                         OPSLANE CAPSTONE DEFENSE WORKFLOW                               |
+-----------------------------------------------------------------------------------------+
|                                                                                         |
|   Candidate Submission                                                                  |
|   +-----------------------+     Evaluation Gates          Defense Panel Examination     |
|   | 14 Milestones Code    |    +----------------------+   +-----------------------+     |
|   | 100% Passing Tests    |--->| - Static AST Audit   |-->| Oral System Design    |     |
|   | Multi-Stage Dockerfile|    | - Race Detection     |   | Live Concurrency Walk |     |
|   | Prometheus / OTel     |    | - RLS Tenant Checks  |   | Chaos/Failure Defense |     |
|   +-----------------------+    +----------------------+   +-----------------------+     |
|                                                                       |                 |
|                                                                       v                 |
|                                                            Formal Graduation Gate       |
|                                                            (Senior Backend Engineer)    |
|                                                                                         |
+-----------------------------------------------------------------------------------------+
```

---

## 2. Capstone Submission Requirements

1. **Repository Completeness**:
   - Complete working implementation in `projects/opslane-flagship/_solution/`.
   - All 14 milestone unit and integration test suites passing cleanly with zero data races (`go test -race`).
2. **Production Artifacts**:
   - Multi-stage Dockerfile producing a static, non-root `scratch` container image.
   - `docker-compose.yml` declaring PostgreSQL, Redis, Prometheus, and the Opslane backend.
   - Database migrations for multi-tenant schema and outbox tables.
3. **Oral Capstone Defense**:
   - Defend the architecture across the 10 questions in `questions.md` covering distributed transactions, tenant isolation, cache coherency, backpressure, and zero-downtime rolling deployments.

---

## 3. 4-Tier Assessment Rubric

| Level | Criteria & Observable Evidence |
|---|---|
| **Level 1: Inadequate** | Code fails to compile or fails automated tests; cross-tenant query leaks; missing graceful shutdown. |
| **Level 2: Developing** | Basic features pass; lacks transactional outbox; unhandled cache stampedes; naive mutexes. |
| **Level 3: Proficient** | All 14 milestones pass with race detection; multi-tenant RLS verified; cache-aside with Redis; OTel tracing; graceful drain on SIGTERM. |
| **Level 4: Exemplary** | Masterful distributed systems engineer. Defends transaction isolation levels, outbox event pipelines, singleflight cache stampede mitigation, and zero-downtime rolling deploys before a panel of senior engineers. |

---

## 4. Verification

```bash
# Verify all 14 milestone test suites pass cleanly
go test -v ./curriculum/modules/18-flagship-opslane/assessments/...

# Verify flagship solution passes all integration tests
go test -v -race ./curriculum/modules/18-flagship-opslane/projects/opslane-flagship/_solution
```
