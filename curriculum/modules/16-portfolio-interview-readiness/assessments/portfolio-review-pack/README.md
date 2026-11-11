# Project Assessment — Portfolio Review Pack

## Purpose

Mastery is validated against the following target items from Portfolio and Interview Readiness:

## Targets

- `project-portfolio-review-pack`

## Evidence required

Provide written answers to each scenario below, any required Go code or patch, captured test output proving it passes, and a short explanation of your design decisions.

---

## Visual Model: Portfolio Audit & Interview Defense Pipeline

```
+-----------------------------------------------------------------------------------------+
|                                PORTFOLIO AUDIT ENGINE                                   |
+-----------------------------------------------------------------------------------------+
|                                                                                         |
|   Showcase Projects (08-15)                                                             |
|   +-----------------------+     Inspection Pipeline       Scorecard Generator           |
|   | 08: REST API          |    +----------------------+   +-----------------------+     |
|   | 09: Postgres & SQL    |--->| - Verify README.md   |-->| Markdown Table        |     |
|   | 10: Auth & RBAC       |    | - Verify *.go source |   | Pass/Fail Status      |     |
|   | 11: Worker Pool       |    | - Verify *_test.go   |   | Gaps & Deficiencies   |     |
|   | 12: Observability     |    | - Verify Artifacts   |   | Readiness Calculation |     |
|   | 13: Rate Limiter      |    +----------------------+   +-----------------------+     |
|   | 14: Hexagonal Arch    |                                           |                 |
|   | 15: Docker & CI/CD    |                                           v                 |
|   +-----------------------+                             Automated Defense Briefs        |
|                                                                                         |
+-----------------------------------------------------------------------------------------+
                                             |
                                             v
+-----------------------------------------------------------------------------------------+
|                         SENIOR TECHNICAL INTERVIEW DEFENSE                              |
+-----------------------------------------------------------------------------------------+
|                                                                                         |
|   [Interviewer Challenge]                                                               |
|   "Why use lock striping with ring buffers instead of sync.Map in your rate limiter?"   |
|                                                                                         |
|   [Senior Engineer Response Matrix]                                                     |
|   1. Machine Mechanics: Cache line ping-pong, atomic contention under heavy writes.     |
|   2. Architectural Trade-offs: Memory footprint vs 38M ops/s zero-allocation latency.   |
|   3. Observability & Failure Modes: Window skew, memory bounds under high cardinality.  |
|                                                                                         |
+-----------------------------------------------------------------------------------------+
```

---

## Mental Model: From Code Author to Engineering Authority

Junior engineers focus on whether code runs. Senior engineers focus on **trade-offs, operational boundaries, and system behavior under failure**:
1. **Verifiable Quality Over Claims**: An automated auditor inspects code, tests, and documentation, ensuring zero hollow claims on your resume.
2. **Defensible Trade-Offs**: Every decision in your codebase (Postgres transactions vs eventual consistency, worker pools vs raw goroutines, hexagonal architecture vs MVC) is backed by concrete operational rationale.
3. **Failure Invariant Awareness**: You articulate what breaks when scale increases 10x, when the database partitions, or when downstream APIs slow down.

---

## Machine View: Static & Filesystem Inspection Architecture

```
Filesystem Workspace Root
 └── curriculum/modules/
      ├── 08-web-services-apis/projects/book-catalog
      │    ├── README.md           (Documented contracts, endpoints, error schema)
      │    ├── main.go             (Clean Go 1.22 routing, recovery middleware)
      │    └── main_test.go        (Unit and route integration assertions)
      ├── ...
      └── 15-docker-cicd-deployment/projects/dockerized-service-cicd
           ├── README.md           (Multi-stage build specs, probe behavior)
           ├── Dockerfile          (Scratch base, UID 10001 appuser)
           ├── docker-compose.yml  (Resource limits, healthcheck config)
           ├── main.go             (Graceful drain shutdown, probe endpoints)
           └── main_test.go        (Health check and lifecycle tests)

Audit Engine Execution Flow:
  1. Stat(path) -> Directory exists
  2. ReadFile(README.md) -> Valid content > 100 bytes
  3. Glob(*.go) -> Source files found
  4. Glob(*_test.go) -> Test files found
  5. CheckArtifacts() -> Dockerfile/Benchmarks/Migrations
  6. Calculate Weighted Score:
     Score = (README*0.30) + (Code*0.30) + (Tests*0.25) + (Proof*0.15)
  7. Overall Portfolio Status = Score >= 80% && PassedProjects >= 7
```

---

## Deliverables & Tasks

1. **Portfolio Audit Engine**:
   - Complete `_starter/main.go` so all unit tests in `_starter/main_test.go` pass.
   - Run the audit against all 8 showcase projects (Modules 08–15).
2. **Scorecard Generation**:
   - Emit a formatted markdown report matching the showcase matrix.
3. **Defense Brief Mastery**:
   - Answer the 10 portfolio and technical defense questions in `questions.md` demonstrating mastery of backend engineering trade-offs.

---

## Common Pitfalls & Anti-Patterns

1. **Superficial Code Review**: Merely checking if files exist without validating test coverage, error handling, or memory efficiency.
2. **Passive Defense**: Responding to interview challenges with "that's just how standard libraries do it" rather than explaining the runtime mechanics.
3. **Ignoring Red Phase**: Skipping starter validation or writing tests that pass trivially without asserting business invariants.

---

## Verification & Testing

```bash
# Verify starter fails on missing implementation
go test ./curriculum/modules/16-portfolio-interview-readiness/projects/portfolio-review-pack/_starter

# Verify solution passes all portfolio audit tests
go test -v ./curriculum/modules/16-portfolio-interview-readiness/projects/portfolio-review-pack/_solution

# Run the live portfolio audit tool
go run ./curriculum/modules/16-portfolio-interview-readiness/projects/portfolio-review-pack/_solution/main.go
```

## Rubric

| Criterion | Weight | Passing standard |
| --- | --- | --- |
| Audit engine correctness | 35% | `_starter/main.go` passes all `_starter/main_test.go` tests and computes the documented weighted score. |
| Scorecard generation | 25% | Emits a formatted markdown report with per-project pass/fail status and readiness calculation. |
| Defense brief mastery | 30% | Answers in `questions.md` defend architectural trade-offs with runtime mechanics and failure-mode awareness. |
| Explanation quality | 10% | Written reasoning is precise and links each decision to operational consequences. |

A passing score is 80%, with no single criterion below 60%.

## Retake policy

If you score below the passing standard, review the corresponding lesson and project deliverables against the rubric, fix the identified gaps, and resubmit with a short note explaining what changed and why.

## NEXT UP

[Module 16 Checkpoint](../checkpoint/README.md)
