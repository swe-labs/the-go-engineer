# The Go Engineer — Architecture v3

## Canonical Architecture

This document defines the authoritative curriculum architecture for **The Go Engineer v3** on `release/v3`.

The v3 architecture establishes a zero-to-software-engineer path focused on production Go backend engineering, accompanied by an elective track and a companion full-stack track.

## Authoritative Metadata

All authoritative curriculum metadata is defined in modular JSON manifests under [`curriculum/`](./curriculum/README.md):

| Manifest | Purpose |
|---|---|
| [`curriculum/path.core.json`](./curriculum/path.core.json) | Canonical 18-module progression, instructional units, and prerequisites |
| [`curriculum/path.electives.json`](./curriculum/path.electives.json) | Advanced electives catalog, prerequisites, and learning outcomes |
| [`curriculum/projects.json`](./curriculum/projects.json) | Portfolio projects, reference solutions, and test fixtures |
| [`curriculum/assessments.json`](./curriculum/assessments.json) | Checkpoint and capstone assessment suites, rubrics, and answer keys |
| [`curriculum/concepts.json`](./curriculum/concepts.json) | Conceptual ownership, mental models, and reinforcement map |
| [`curriculum/crossrefs.json`](./curriculum/crossrefs.json) | Cross-module forward and backward references |
| [`curriculum/failures.json`](./curriculum/failures.json) | Explicit failure mode taxonomy and defensive patterns |
| [`curriculum/readme.contracts.json`](./curriculum/readme.contracts.json) | Mandatory pedagogical contract structure for instructional READMEs |
| [`curriculum/schema.v3.json`](./curriculum/schema.v3.json) | Structural JSON schema and validation definitions for all v3 manifests |
| [`curriculum/workspace.json`](./curriculum/workspace.json) | Canonical workspace and validation configuration |
| [`curriculum/legacy/curriculum.v2.json`](./curriculum/legacy/curriculum.v2.json) | Frozen legacy v2 import archive |

## Core Module Map

Instructional content is organized canonically under `curriculum/modules/`:

| Module | Title | Primary Focus | Key Deliverables |
|---|---|---|---|
| `00-orientation` | Orientation | Mental models, learning methodology, setup | Course setup verification, study plan |
| `01-computers-terminal-git-web` | Computers, Terminal, Git & Web | Shell navigation, git lifecycle, networking | Shell & Git lab, HTTP inspection |
| `02-go-setup-tooling` | Go Setup & Tooling | Toolchain, compiler, module graph, formatting | Tooling failure diagnosis lab |
| `03-programming-fundamentals` | Programming Fundamentals | Data types, memory layout, control flow | Pricing checkout, contact directory |
| `04-functions-errors-data-semantics` | Functions, Errors & Data Semantics | Stack/heap allocation, pointers, errors | Order summary refactor |
| `05-types-interfaces-packages-modules` | Types, Interfaces & Packages | Structural polymorphism, composition | Bank account, payroll processor |
| `06-testing-debugging-refactoring` | Testing, Debugging & Refactoring | Unit tests, table tests, coverage, benchmarks | Bug hunt pack, refactored library |
| `07-cli-files-json-config` | CLI, Files, JSON & Config | Flag parsing, streaming I/O, serialization | Config parser, log analyzer CLI |
| `08-http-rest-apis` | HTTP & REST APIs | Routing, middleware, JSON APIs, httptest | Book catalog REST API |
| `09-sql-postgres-persistence` | SQL, PostgreSQL & Persistence | Connection pooling, migrations, transactions | Task API with PostgreSQL |
| `10-auth-security` | Authentication & Security | Bcrypt hashing, JWT tokens, RBAC | Authenticated security service |
| `11-lifecycle-context-concurrency` | Lifecycle, Context & Concurrency | Goroutines, channels, worker pools, select | Bounded worker pool, downloader |
| `12-observability-diagnostics` | Observability & Diagnostics | Structured logging (`slog`), pprof, tracing | Observability retrofit lab |
| `13-performance-memory-engineering` | Performance & Memory Engineering | Benchmarks, memory profiles, zero-alloc | Profile & fix optimization lab |
| `14-architecture-distributed-systems` | Architecture & Distributed Systems | Hexagonal architecture, domain decoupling | Modular refactor review |
| `15-docker-cicd-deployment` | Docker, CI/CD & Deployment | Multi-stage Docker, GitHub Actions, health | Dockerized production service |
| `16-portfolio-interview-readiness` | Portfolio & Interview Readiness | Architecture defense, portfolio packaging | Portfolio review pack, rubrics |
| `18-flagship-opslane` | Flagship Opslane Capstone | Production multi-tenant incident management | Opslane SaaS platform (15 milestones) |

## Companion & Elective Tracks

- **Companion Frontend Track (`curriculum/companion-frontend/`):** 4 progressive stages exploring full-stack Go integration (Semantic HTML/CSS, Vanilla DOM & Fetch, Templ & Tailwind CSS, HTMX & Alpine.js).
- **Advanced Electives (`curriculum/electives/advanced-electives/`):** 20 advanced electives (Generics, sync.Pool, sync.Map, gRPC streaming, distributed consensus, eBPF) outside the core path.

## Quality Standards

- **Zero Magic:** Every lesson establishes explicit mental models and machine execution traces before code.
- **Runnable & Tested:** All projects include reference solutions (`_solution/`), student starter templates (`_starter/`), and comprehensive test suites (`main_test.go`).
- **Automated Quality Gates:** Enforced by `internal/tools/curriculum` across 17 automated gates.

## Validation

Run CI-equivalent curriculum validation:

```bash
go run ./internal/tools/curriculum validate-all
```
