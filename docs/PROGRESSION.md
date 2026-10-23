# The Go Engineer Progression (v3)

This document visualizes the canonical v3 learner journey from foundational concepts to production systems engineering on `release/v3`.

Module progression and item metadata are formally specified in [`curriculum/path.core.json`](../curriculum/path.core.json) and [ARCHITECTURE.md](../ARCHITECTURE.md).

## Phase Progression

```mermaid
flowchart TD
    P0["Phase 0: Orientation & Foundations\nModules 00–02\nTerminal, Git, OS concepts, Go toolchain"] --> P1["Phase 1: Go Core & Semantics\nModules 03–04\nTypes, memory layout, control flow, functions, errors"]
    P1 --> P2["Phase 2: Engineering Core\nModules 05–07\nInterfaces, polymorphism, testing, CLI tools, files"]
    P2 --> P3["Phase 3: Backend & Persistence\nModules 08–09\nHTTP/REST APIs, routing, middleware, SQL, PostgreSQL"]
    P3 --> P4["Phase 4: Security & Concurrency\nModules 10–11\nJWT auth, RBAC, goroutines, channels, worker pools, context"]
    P4 --> P5["Phase 5: Production Systems\nModules 12–15\nslog, pprof, tracing, architecture, Docker, CI/CD"]
    P5 --> P7["Phase 7: Flagship Capstone\nModule 18\nOpslane Multi-tenant SaaS Platform (15 milestones)"]
    P7 --> P6["Phase 6: Career & Portfolio\nModule 16\nPortfolio Packaging, Technical Writing & Architecture Defense"]
```

## Module Flow

```mermaid
flowchart TD
    M00["00: Orientation"] --> M01["01: Computers, Terminal, Git & Web"]
    M01 --> M02["02: Go Setup & Tooling"]
    M02 --> M03["03: Programming Fundamentals"]
    M03 --> M04["04: Functions, Errors & Data Semantics"]
    M04 --> M05["05: Types, Interfaces & Packages"]
    M05 --> M06["06: Testing, Debugging & Refactoring"]
    M06 --> M07["07: CLI, Files, JSON & Config"]
    M07 --> M08["08: HTTP & REST APIs"]
    M08 --> M09["09: SQL, PostgreSQL & Persistence"]
    M09 --> M10["10: Auth & Security"]
    M10 --> M11["11: Lifecycle, Context & Concurrency"]
    M11 --> M12["12: Observability & Diagnostics"]
    M12 --> M13["13: Performance & Memory Engineering"]
    M13 --> M14["14: Architecture & Distributed Systems"]
    M14 --> M15["15: Docker, CI/CD & Deployment"]
    M15 --> M18["18: Flagship Opslane Capstone"]
    M18 --> M16["16: Portfolio & Interview Readiness"]
```

## Engineering Context Growth

| Phase | Learner Shift | Engineering Weight |
|---|---|---|
| **Phase 0: Foundations** (00–02) | Mental model of OS, process memory, and toolchain mechanics | Concrete, execution-level |
| **Phase 1: Go Core** (03–04) | Understand stack vs heap allocation, pointer semantics, errors | Foundational systems |
| **Phase 2: Engineering Core** (05–07) | Design structural polymorphism, test table suites, streaming I/O | Standard engineering |
| **Phase 3: Persistence** (08–09) | Build robust REST handlers, connection pools, and DB transactions | High, production-facing |
| **Phase 4: Concurrency** (10–11) | Coordinate bounded goroutine pools, cancellation, and JWT/RBAC | High, concurrency-critical |
| **Phase 5: Production Systems** (12–15) | Diagnose memory bottlenecks (`pprof`), trace requests, deploy via CI/CD | Advanced systems |
| **Phase 7: Flagship Capstone** (18) | Build multi-tenant SaaS backend with 15 progressive milestones | Full system integration |
| **Phase 6: Portfolio & Defense** (16) | Defend architecture, package production deliverables, interview prep | Professional readiness |

## Key Projects & Milestones

| Module | Project / Milestone | Core Deliverable | Verification Mode |
|---|---|---|---|
| `00` | Course Setup Verification | Verify runtime, editor, Git workflow, study plan | Automated script |
| `01` | Shell & Git Navigation Lab | Shell execution trace, commit hygiene | Terminal inspection |
| `02` | Tooling Failure Lab | Diagnose PATH drift, broken build tags, module conflicts | Red-to-green repair |
| `03` | Pricing Checkout & Contact Directory | Multi-branch pricing engine, pointer-backed address book | Unit test suite |
| `04` | Order Summary Refactor | Sentinel errors, custom error types, stack allocations | Table test suite |
| `05` | Bank Account & Payroll Processor | Interface abstraction, mockable polymorphism | Test suite with mocks |
| `06` | Bug Hunt Pack & Refactored Library | TDD bug fixes, benchmark assertions, coverage audit | Coverage + race tests |
| `07` | Config Parser & Log Analyzer | Streaming JSON decode, flag parsing, file I/O | CLI integration test |
| `08` | Book Catalog REST API | HTTP router, middleware chain, httptest fixtures | HTTP integration tests |
| `09` | Task API with PostgreSQL | Connection pooling, SQL migrations, ACID transactions | Database test suite |
| `10` | Secure Authenticated API | Bcrypt password hashing, JWT claims, RBAC middleware | Auth security tests |
| `11` | Bounded Worker Pool & Downloader | Worker pool, semaphore pattern, graceful context cancel | Race-detector test suite |
| `12` | Observability Retrofit | Structured `slog`, OpenTelemetry tracing, Prometheus | Metrics inspection |
| `13` | Performance & Memory Lab | Benchmark optimization, pprof heap/CPU analysis | Zero-alloc benchmarks |
| `14` | Modular Refactor Review | Hexagonal decoupling, domain invariants, ADR documentation | Domain boundary audit |
| `15` | Dockerized Production Service | Multi-stage Dockerfile, container health checks, CI/CD runbook | Docker health check |
| `18` | Flagship Opslane Platform | 15 milestones: multi-tenancy, workers, caching, graceful shutdown | Complete capstone suite |
| `16` | Portfolio Review Pack & Defense | Technical README, architecture diagrams, live defense answers | Portfolio rubric audit |

## Completion Standard

A learner who completes the curriculum demonstrates:

1. **Zero Magic Understanding:** Explains memory layouts, pointer dereferencing, goroutine scheduler mechanics, and CPU cache friendliness.
2. **Production Code Quality:** Idiomatic Go adhering to standard conventions, table-driven tests, race detection, and structured logging.
3. **Resilient Systems Engineering:** Defensive error classification, bounded concurrency, graceful shutdown, and resource leak prevention.
4. **Architectural Mastery:** Capable of designing, building, profiling, and defending multi-tenant backend systems at scale.
