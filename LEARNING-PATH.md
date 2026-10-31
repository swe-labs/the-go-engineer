# The Go Engineer Learning Path (v3)

This guide explains how to navigate through the canonical v3 curriculum on `release/v3`.

The source of truth for module progression is [`curriculum/path.core.json`](./curriculum/path.core.json) and [ARCHITECTURE.md](./ARCHITECTURE.md).

> **Visual Diagrams:** For interactive Mermaid flowcharts mapping both phase-level and module-level dependency graphs, see [`docs/PROGRESSION.md`](./docs/PROGRESSION.md).

## Phase Overview

| Phase | Modules | Focus | Key Projects & Milestones |
|---|---|---|---|
| **0. Orientation & Foundations** | `00`–`02` | Terminal, Git, operating system concepts, and Go toolchain mechanics | Course Setup Verification, Shell & Git Lab, Tooling Failure Lab |
| **1. Go Core & Semantics** | `03`–`04` | Types, memory layout, control flow, functions, and error handling | Pricing Checkout, Contact Directory, Order Summary Refactor |
| **2. Engineering Core** | `05`–`07` | Interfaces, polymorphism, test-driven refactoring, CLI tools, JSON, files | Bank Account, Payroll Processor, Bug Hunt Pack, Config Parser |
| **3. Backend & Persistence** | `08`–`09` | HTTP/REST APIs, routing, middleware, SQL, PostgreSQL migrations, connection pools | Book Catalog REST API, Task API with PostgreSQL |
| **4. Security & Concurrency** | `10`–`11` | Bcrypt, JWT auth, RBAC, goroutines, channels, bounded worker pools, context | Secure Authenticated API, Bounded Worker Pool, Concurrent Downloader |
| **5. Production Systems** | `12`–`15` | Structured logging (`slog`), Prometheus, OpenTelemetry, profiling (`pprof`), hexagonal architecture, Docker, CI/CD | PII Redactor, Observability Retrofit, Performance Lab, Dockerized Service |
| **6. Career & Portfolio** | `16` | Portfolio review engine, technical documentation, architectural interview defense | Portfolio Review Pack & Defense Guides |
| **7. Flagship Capstone** | `18` | Multi-tenant incident management, event bus, worker pools, caching, and graceful shutdown | Flagship Opslane Platform (15 milestones) |

## Companion Tracks

- **[Companion Full-Stack Frontend Track](./curriculum/companion-frontend/README.md):** 4 progressive stages (HTML/CSS, DOM/Fetch, Templ/Tailwind, HTMX/Alpine).
- **[Advanced Electives](./curriculum/path.electives.json):** 20 deep-dive electives in `curriculum/electives/advanced-electives/` (Generics, sync.Pool, sync.Map, gRPC streaming, consensus, eBPF).

## Recommended Learning Paths

### 1. Full Path (Zero to Production Engineer)
Follow Modules 00 through 18 sequentially. Complete all lessons, labs, project starter refactors, and milestone assessments.

### 2. Experienced Engineer Bridge Path
- Skim Modules 00–01 if already comfortable with Unix/Git.
- Complete Module 02 (tooling) and Modules 03–04 to cement Go pointer/value semantics and memory layout.
- Complete projects in Modules 05–15.
- Build and defend the Flagship Opslane capstone (Module 18) and prepare portfolio artifacts (Module 16).

## Verification & Proof

Every module provides automated verification:

- **Lessons:** Runnable code with verified machine execution traces (`go run main.go`).
- **Projects:** Reference solutions (`_solution/`), starter templates (`_starter/`), and test suites (`go test ./...`).
- **Milestones:** Comprehensive assessments, rubrics, and answer keys.
- **Repository Quality Gates:** 17 automated gates checked via `go run ./internal/tools/curriculum validate-all`.
