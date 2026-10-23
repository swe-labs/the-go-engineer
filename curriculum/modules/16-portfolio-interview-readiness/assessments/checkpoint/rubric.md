# Rubric — Module 16 Checkpoint: Portfolio & Interview Readiness

## Assessment Overview

This checkpoint evaluates mastery of career positioning, portfolio curation, live coding execution, system design leadership, and interview negotiation across lessons `core-16-01` through `core-16-15`.

---

## 4-Tier Scoring Rubric

| Criterion | Level 1: Inadequate | Level 2: Developing | Level 3: Proficient | Level 4: Exemplary |
|---|---|---|---|---|
| **Portfolio & Presentation** | Empty or non-functional repositories; no documentation or tests; relies on unverifiable claims. | Compiling code with minimal documentation; missing benchmarks and operational artifacts; generic resume bullets. | Curated portfolio of 8 production-grade showcase projects with passing tests, race detection, benchmarks, and Dockerfiles; impact-quantified resume. | Outstanding industry portfolio: public technical case studies, verified high-throughput benchmarks (38M ops/s), automated portfolio auditor, and open-source contributions. |
| **Live Coding & Code Quality** | Code does not compile; jumps straight to coding without asking questions; ignores errors and data races. | Solves problem with brute-force; acknowledges concurrency but uses naive locks; minimal test coverage. | Clarifies requirements before coding; implements generic, thread-safe Go code with `sync.RWMutex`; handles errors and context cancellation; writes table-driven tests. | Flawless live execution: discusses time/space complexity, cache line bouncing, escape analysis, and zero-allocation patterns; writes benchmarks and race-clean tests under time pressure. |
| **System Design Interview Leadership** | Struggles to gather requirements; suggests monolithic database for massive scale; unaware of CAP theorem or network failure modes. | Draws basic boxes (client -> load balancer -> server -> DB); suggests caching without invalidation strategy; reactive to prompts. | Drives the interview: clarifies non-functional requirements, establishes back-of-the-envelope estimates, evaluates database trade-offs, and defines graceful degradation. | Staff-level architectural mastery: designs multi-region active-active architectures, articulates consensus mechanics (Raft/Paxos), models failure blast radius, and handles network partitions. |
| **Behavioral Leadership & Negotiation** | Blames teammates for outages; unable to articulate personal contributions; hostile or passive negotiation approach. | Uses vague examples; focuses on tasks rather than measurable outcomes; accepts initial offers without benchmarking. | Frames past incidents using STAR framework with clear personal ownership and blameless RCA; conducts data-driven compensation discussions. | Transformational engineering leader: shares deep lessons from complex distributed outages; mentors junior engineers; commands top-of-market compensation through competing leverage. |

---

## Passing Criteria

- **Minimum Total Score**: 80% (Proficient or Exemplary across all categories).
- **All 10 Questions Answered**: Complete architectural designs, live coding implementations, and negotiation frameworks provided.
