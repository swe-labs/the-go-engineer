# Module 16 Assessment — Portfolio and Interview Readiness

## Purpose

Validate mastery for the target curriculum items.

## Targets

- `core-16-01`
- `core-16-02`
- `core-16-03`
- `core-16-04`
- `core-16-05`
- `core-16-06`
- `core-16-07`
- `core-16-08`
- `core-16-09`
- `core-16-10`
- `core-16-11`
- `core-16-12`
- `core-16-13`
- `core-16-14`
- `core-16-15`

## Evidence

Submit answers, code evidence, test output, and explanation where required.

---

## 1. Visual Model: Career Engineering & Interview Execution Lifecycle

```
+-----------------------------------------------------------------------------------------+
|                        ENGINEERING PORTFOLIO & REPUTATION ENGINE                        |
+-----------------------------------------------------------------------------------------+
|                                                                                         |
|   Technical Artifacts          Documentation & Outreach       Hiring Funnel             |
|   +-----------------------+    +-----------------------+    +-----------------------+   |
|   | 8 Showcase Projects   |--->| Deep Architectural    |--->| Inbound Senior Lead   |   |
|   | Open Source Prs       |    | Technical Blog Posts  |    | Target Company List   |   |
|   | Verified Benchmarks   |    | Public Case Studies   |    | Optimized Resume      |   |
|   +-----------------------+    +-----------------------+    +-----------------------+   |
|                                                                         |               |
+-------------------------------------------------------------------------|---------------+
                                                                          v
+-----------------------------------------------------------------------------------------+
|                       MULTI-STAGE SENIOR INTERVIEW DEFENSE                              |
+-----------------------------------------------------------------------------------------+
|                                                                                         |
|   Stage 1: Live Coding & Concurrency                                                    |
|   - Clean idiomatic Go, table-driven tests, data race detection, error wrapping.        |
|                                                                                         |
|   Stage 2: Distributed System Design                                                    |
|   - Back-of-the-envelope estimation, storage trade-offs, isolation levels, outbox event |
|     pipelines, caching invalidation, failure blast radius containment.                  |
|                                                                                         |
|   Stage 3: Behavioral & Staff-Plus Leadership                                           |
|   - STAR framework, production incident postmortems, mentorship, cross-team alignment.  |
|                                                                                         |
|   Stage 4: Total Compensation Negotiation                                               |
|   - Market benchmarking, competing offer leverage, equity vesting mechanics, leveling. |
|                                                                                         |
+-----------------------------------------------------------------------------------------+
```

---

## 2. Core Scenarios

### Scenario 1: Portfolio Strategy & Architectural Case Studies (`core-16-01` to `core-16-04`)
You are curating your GitHub profile for Senior Backend Engineer roles at top-tier distributed systems and cloud infrastructure companies.
- Why are "to-do list apps" or generic tutorial clones detrimental to senior candidacy?
- What are the four essential artifacts every showcase repository must provide?
- Structure a README for a distributed rate limiter that highlights throughput benchmarks (38M ops/s), concurrency safety under `go test -race`, and zero-allocation memory profiling.

### Scenario 2: System Design & Live Coding Mechanics (`core-16-07`, `core-16-08`, `core-16-09`)
In a 45-minute live coding screen, you are asked to design a concurrent in-memory cache with TTL expiration.
- How do you structure your communication before typing code?
- Why is an active background cleaner goroutine with `time.Ticker` preferable to per-key timers (`time.AfterFunc`)?
- How do you write a table-driven test and benchmark in Go to prove concurrent correctness and low latency to the interviewer?

### Scenario 3: Behavioral Leadership & Incident Postmortems (`core-16-06`, `core-16-14`)
An interviewer asks: *"Tell me about a time a production release caused a severe outage and how you handled the resolution and post-incident process."*
- Frame your response using the STAR method (Situation, Task, Action, Result).
- Emphasize blameless postmortem culture, root cause analysis (RCA), corrective action tracking, and automated regression testing.

### Scenario 4: Offer Negotiation & Career Growth (`core-16-10`, `core-16-15`)
You receive an initial offer of $160,000 base + equity, but market benchmarks for your level indicate $185,000 base + equity.
- How do you craft a professional, non-adversarial counter-proposal leveraging your engineering qualifications and competing pipeline opportunities?
- Explain the key components of equity compensation: RSUs vs ISOs, vesting schedules, 83(b) elections, and liquidity events.

---

## 3. Invariant Career Rules

1. **Never Make Unverifiable Claims**: If a technology or architectural pattern is listed on your resume, you must be prepared to write it, profile it, and debug it live.
2. **Code Speaks Through Tests**: In live coding and take-home exercises, untested code is considered broken code. Always write tests first or concurrently.
3. **Always Quantify Engineering Impact**: Frame accomplishments in terms of measurable operational metrics: latency reduction (P99), cost savings, throughput increases, or incident MTTR.
4. **Negotiate from Value, Never Entitlement**: Ground compensation discussions in verified market data and specific technical contributions.

---

## 4. 4-Tier Assessment Rubric

| Level | Criteria & Observable Evidence |
|---|---|
| **Level 1: Inadequate** | Portfolio contains incomplete tutorial clones; unable to write idiomatic Go live; no test coverage; reactive interview communication. |
| **Level 2: Developing** | Projects compile but lack benchmarks and architecture briefs; writes functional code but ignores data races and error handling; generic behavioral answers. |
| **Level 3: Proficient** | 8 showcase projects with 100% test coverage; leads system design interviews with clear trade-offs and capacity planning; uses STAR framework effectively. |
| **Level 4: Exemplary** | Authoritative systems engineer. Publishes technical case studies and open-source packages; effortlessly navigates runtime internals, pprof profiles, and distributed consensus in interviews; negotiates senior/staff compensation with professional market leverage. |

---

## 5. Submission Requirements

Submit your answers to the 10 questions in `questions.md` demonstrating mastery of portfolio curation, live technical communication, and career engineering.
