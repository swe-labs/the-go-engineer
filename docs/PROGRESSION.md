# The Go Engineer — Visual Progression

> This document visualises the learning journey through the v2.1 curriculum.
> Section IDs and milestones must match `ARCHITECTURE.md`. If they conflict, `ARCHITECTURE.md` wins.

---

## Overall Progress Model

```text
Phase 0       Phase 1              Phase 2                Phase 3            Phase 4
(0–5%)        (5–52%)              (52–87%)               (87–96%)           (96–100%)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
MACHINE       LANGUAGE             ENGINEERING            SYSTEMS            FLAGSHIP
FOUNDATION    FOUNDATION           CORE                   ENGINEERING        PROJECT
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
"I understand "I write Go          "I build systems       "I design and      "I build and
 the machine"  fluently"            that work at scale"    operate systems"   ship it"
```

---

## Learning Path: Zero to Senior

### Path A: Complete Beginner

```text
START
  │
  ▼
[s00] How Computers Work ──→ [s01] Getting Started ──→ [s02] Language Basics ──→ [s03] Functions & Errors
  (0%)                        (5%)                      (12%)                     (28%)
  │
  ▼
[s04] Types & Design ──→ [s05] Packages, I/O & CLI ──→ [s06] Backend, APIs & DB ──→ [s07] Concurrency
  (38%)                    (52%)                          (62%)                       (75%)
  │
  ▼
[s08] Quality & Testing ──→ [s09] Architecture & Security ──→ [s10] Production Ops ──→ [s11] GoScale
  (83%)                       (87%)                            (92%)                    (96% → 100%)
```

---

## Engineering Context Growth

```text
Engineering Context %

100% ┤
     │                                                          ████ GoScale ████
 90% ┤                                                    ████
     │                                              ████
 80% ┤                                        ████
     │                                  █████
 70% ┤                           █████
     │                     █████
 60% ┤               █████                    ●●●● Phase 2: Engineering Core
     │          █████                         ●●●●
 50% ┤     ●●●●                          ●●●●
     │     ●●●●                       ●●●●
 40% ┤     ●●●●                 ●●●●
     │     ●●●●           ●●●●
 30% ┤     ●●●●      ●●●●
     │     ●●●●●●●●●●●●
 20% ┤     ●●●●●●●●●●●●●●  ●●●● Phase 0–1: Foundation
     │                      ●●●●
 10% ┤
     │
  0% └──────────────────────────────────────────────────────────────
         s00–s02       s03–s04       s05–s08        s09–s10       s11

     ● = Syntax + Basic Why
     █ = Full Engineering Context
```

---

## Key Milestones

These milestones match the Milestone Progression table in `ARCHITECTURE.md`:

| %    | Milestone                 | Lesson   | Engineering Checkpoint                          |
| ---- | ------------------------- | -------- | ----------------------------------------------- |
| 5%   | Machine model checkpoint  | HC.5     | Explain CPU, memory, processes, signals          |
| 10%  | First program             | GT.2     | Run and modify Hello World                       |
| 13%  | Toolchain confident       | GT.6     | Read and fix compiler errors                     |
| 18%  | Pricing Checkout          | CF.7     | Handle all control flow cases cleanly            |
| 24%  | Contact Directory         | DS.6     | Use slices + maps + pointers together            |
| 30%  | Order Summary             | FE.7     | Handle validation, propagation, orchestration    |
| 44%  | Payroll Processor         | TI.10    | Implement interfaces, type switches, generics    |
| 51%  | Config Parser (strings)   | ST.6     | Parse text with templates and regex              |
| 58%  | Log Search CLI            | FS.7     | Build a CLI tool with filesystem I/O             |
| 66%  | REST API with timeouts    | HS.10    | Full HTTP middleware stack with timeouts          |
| 70%  | gRPC Service              | API.9    | Build a gRPC service with interceptors           |
| 74%  | Repository Pattern        | DB.6     | Database CRUD with transactions                  |
| 77%  | Concurrent Downloader     | GC.7     | Handle race conditions and goroutine lifecycle    |
| 81%  | URL Health Checker        | CP.5     | Debug concurrent failures                        |
| 85%  | Benchmark optimisation    | PR.5     | Profile and optimise with before/after evidence   |
| 88%  | Modular Refactor Capstone | ARCH.9   | Refactor to hexagonal architecture               |
| 91%  | Secure API                | SEC.11   | Apply OWASP Top 10 checklist                     |
| 94%  | Dockerised Service        | DEPLOY.3 | Deploy to container with CI/CD                   |
| 96%  | Shutdown Capstone         | GS.3     | Graceful shutdown with zero-downtime              |
| 100% | GoScale Complete          | s11      | Full production-grade SaaS backend               |

---

## Core Engineering Principles Per Phase

### Phase 0: Machine Foundation

- **Core Idea**: Code is instructions for a computer
- **Key Insight**: Programs, memory, processes, and signals
- **Engineering Mind**: "What does the computer actually do?"

### Phase 1: Language Foundation

- **Core Idea**: Variables, types, control flow, functions, errors
- **Key Insight**: Fail fast, validate early, handle errors explicitly
- **Engineering Mind**: "What happens when things go wrong?"

### Phase 2: Engineering Core

- **Core Idea**: Concurrency is about coordination, not parallelism
- **Key Insight**: Race conditions, deadlocks, goroutine leaks, context cancellation
- **Engineering Mind**: "What breaks when 1000 users use this?"

### Phase 3: Systems Engineering

- **Core Idea**: Systems are observed, not guessed
- **Key Insight**: Architecture trade-offs, security hardening, deployment
- **Engineering Mind**: "How do I debug at 3 AM?"

### Phase 4: Flagship Project

- **Core Idea**: Everything comes together
- **Key Insight**: Real-world complexity, trade-offs, decisions under pressure
- **Engineering Mind**: "What would a senior engineer do?"

---

## The Promise

```text
By completing this curriculum, you will be able to:

  ✓ Explain how a computer executes your code
  ✓ Write Go code from scratch
  ✓ Structure code for maintainability
  ✓ Handle errors properly with the three-tier framework
  ✓ Write concurrent code safely
  ✓ Test code comprehensively
  ✓ Profile and optimise performance
  ✓ Build production HTTP and gRPC services
  ✓ Work with databases reliably
  ✓ Secure your applications against common vulnerabilities
  ✓ Deploy and operate systems
  ✓ Think like a senior engineer
```

---

_This document is maintained alongside ARCHITECTURE.md. The curriculum will evolve as sections are completed._
