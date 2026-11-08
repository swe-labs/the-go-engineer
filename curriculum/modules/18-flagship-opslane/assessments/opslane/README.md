# Project Assessment — Opslane Flagship Capstone Defense

## Purpose

Mastery is validated against the following target items from Flagship Opslane:

## Targets

- `project-opslane`

## Evidence required

Provide written answers to each scenario below, any required Go code or patch, captured test output proving it passes, and a short explanation of your design decisions.

---

## Visual Model: Capstone Defense Evaluation Pipeline

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

## Capstone Submission Requirements

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

## 4-Tier Assessment Rubric

| Level | Criteria & Observable Evidence |
|---|---|
| **Level 1: Inadequate** | Code fails to compile or fails automated tests; cross-tenant query leaks; missing graceful shutdown. |
| **Level 2: Developing** | Basic features pass; lacks transactional outbox; unhandled cache stampedes; naive mutexes. |
| **Level 3: Proficient** | All 14 milestones pass with race detection; multi-tenant RLS verified; cache-aside with Redis; OTel tracing; graceful drain on SIGTERM. |
| **Level 4: Exemplary** | Masterful distributed systems engineer. Defends transaction isolation levels, outbox event pipelines, singleflight cache stampede mitigation, and zero-downtime rolling deploys before a panel of senior engineers. |

---

## Verification

```bash
# Verify all 14 milestone test suites pass cleanly
go test -v ./curriculum/modules/18-flagship-opslane/assessments/...

# Verify flagship solution passes all integration tests
go test -v -race ./curriculum/modules/18-flagship-opslane/projects/opslane-flagship/_solution
```

## Retake policy

If you score below the passing standard, review the corresponding lesson and project deliverables against the rubric, fix the identified gaps, and resubmit with a short note explaining what changed and why.
