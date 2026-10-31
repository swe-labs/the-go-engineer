# The Go Engineer

[![CI](https://github.com/rasel9t6/the-go-engineer/actions/workflows/ci.yml/badge.svg)](https://github.com/rasel9t6/the-go-engineer/actions)
[![License: TGE License v1.0](https://img.shields.io/badge/License-TGE_v1.0-red.svg)](#license)

**The Go Engineer** is an open, repository-first curriculum designed to take software engineers from foundational programming to production-grade Go systems engineering.

Every lesson is built on a **zero-magic** standard: mental models of the operating system and Go runtime, runnable self-contained examples, automated test suites, and explicit failure-mode engineering—without unexplained framework abstractions.

---

## Release Status

- **Current Release:** `v3.0.0`
- **Active Release Line:** `release/v3`
- **Previous Release Line:** `v2.1.1` (archived in git history)
- **Authoritative Curriculum Metadata:** [`curriculum/`](./curriculum/README.md)

---

## Curriculum Structure

The canonical v3 curriculum consists of **18 core modules** organized into progressive engineering phases, culminating in the **Opslane Flagship Capstone** and **Portfolio Defense**:

| Phase | Modules | Primary Focus | Key Projects & Milestones |
|---|---|---|---|
| **0. Orientation & Foundations** | `00`–`02` | Terminal, Git, operating system concepts, and Go toolchain mechanics | Course Setup Verification, Shell & Git Lab, Tooling Failure Lab |
| **1. Go Core & Semantics** | `03`–`04` | Types, memory layout, control flow, functions, and error handling | Pricing Checkout, Contact Directory, Order Summary Refactor |
| **2. Engineering Core** | `05`–`07` | Interfaces, polymorphism, test-driven refactoring, CLI tools, JSON, files | Bank Account, Payroll Processor, Bug Hunt Pack, Config Parser |
| **3. Backend & Persistence** | `08`–`09` | HTTP/REST APIs, routing, middleware, SQL, PostgreSQL migrations, connection pools | Book Catalog REST API, Task API with PostgreSQL |
| **4. Security & Concurrency** | `10`–`11` | Bcrypt, JWT auth, RBAC, goroutines, channels, bounded worker pools, context | Secure Authenticated API, Bounded Worker Pool, Concurrent Downloader |
| **5. Production Systems** | `12`–`15` | Structured logging (`slog`), Prometheus, OpenTelemetry, profiling (`pprof`), hexagonal architecture, Docker, CI/CD | PII Redactor, Observability Retrofit, Performance Lab, Dockerized Service |
| **6. Career & Portfolio** | `16` | Portfolio review engine, technical documentation, architectural interview defense | Portfolio Review Pack & Defense Guides |
| **7. Flagship Capstone** | `18` | Multi-tenant incident management, event bus, worker pools, caching, and graceful shutdown | Flagship Opslane Platform (15 milestones) |

### Companion Tracks

- **[Companion Frontend Track](./curriculum/companion-frontend/README.md):** 4 progressive stages exploring full-stack Go integration (Semantic HTML/CSS, Vanilla DOM & Fetch, Templ & Tailwind CSS, and HTMX & Alpine.js).
- **[Advanced Electives (Module 17)](./curriculum/path.electives.json):** Independent advanced topics (distributed consensus, gRPC streaming, eBPF) outside the core capstone path.

---

## Quick Start

### Prerequisites

- **Go:** 1.24+ (declared in [`go.mod`](./go.mod))
- **Git:** Standard git CLI
- **Operating System:** Linux, macOS, or Windows (WSL2 or MinGW-w64 recommended for CGO race detection)

### Setup & First Lesson

```bash
# 1. Clone the repository
git clone https://github.com/rasel9t6/the-go-engineer.git
cd the-go-engineer

# 2. Download dependencies
go mod download

# 3. Run your first lesson
go run ./curriculum/modules/00-orientation/lessons/01-how-to-use-this-repository/main.go
```

---

## Verification & Testing

The repository includes both curriculum structural validators and standard Go test suites:

```bash
# 1. Run all 17 automated curriculum quality gates
go run ./internal/tools/curriculum validate-all

# 2. Run the full unit and integration test suite
go test ./...

# 3. Run static analysis
go vet ./...

# 4. Verify code formatting
gofmt -l .

# 5. Run concurrency race detector on core modules
go test -race ./curriculum/modules/11-lifecycle-context-concurrency/... ./curriculum/modules/18-flagship-opslane/...
```

---

## Authoritative Documentation

- **[System Architecture (`ARCHITECTURE.md`)](./ARCHITECTURE.md):** Canonical v3 architecture, metadata manifests, and zero-magic engineering standards.
- **[Curriculum Map & Pacing Guide](./curriculum/README.md):** Comprehensive 52-week roadmap, workload estimates, and module catalog.
- **[Visual Progression Charts (`docs/PROGRESSION.md`)](./docs/PROGRESSION.md):** Interactive flowcharts detailing phase transitions and module dependency flow.
- **[Flagship Opslane Capstone Architecture (`docs/flagship/OPSLANE_SAAS_BACKEND.md`)](./docs/flagship/OPSLANE_SAAS_BACKEND.md):** Production multi-tenant SaaS backend architecture and milestone rubrics.
- **[Quality Rubrics](./curriculum/quality-rubric.md):** 5-point evaluation criteria for code quality, proof, and documentation.
- **[Zero-Magic Authoring Specification](./curriculum/zero-magic-authoring-spec.md):** Principles for mental models and transparent systems engineering.
- **[Release History & Changelog](./CHANGELOG.md):** Version changelog and migration notes.
- **[Contribution Guidelines](./CONTRIBUTING.md):** Workflow for contributors.

---

## License

This project is licensed under The Go Engineer License (TGE License) v1.0.

- Free for personal, educational, and non-commercial learning.
- Commercial redistribution, corporate training programs, or commercial platform hosting requires written permission.
