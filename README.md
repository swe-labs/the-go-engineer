# The Go Engineer

[![CI](https://github.com/swe-labs/the-go-engineer/actions/workflows/ci.yml/badge.svg)](https://github.com/swe-labs/the-go-engineer/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](./LICENSE)

> **Release status:** `v3.0.0` on the `release/v3` line. Authoritative curriculum metadata lives in [`curriculum/`](./curriculum/README.md).

## What is The Go Engineer?

**The Go Engineer** is an open, repository-first curriculum that takes software engineers from foundational programming to production-grade Go systems engineering.

Every lesson is built on a **zero-magic** standard: mental models of the operating system and the Go runtime, runnable self-contained examples, automated test suites, and explicit failure-mode engineering — with no unexplained framework abstractions. Nothing is hand-waved. If a lesson introduces a mechanism, it explains what the machine actually does: how memory is laid out, how goroutines are scheduled, how the database connection pool drains, and what happens when the process receives `SIGTERM`.

The curriculum is not a collection of snippets. It is a single governed repository: 19 modules, hundreds of lessons and labs, module project capstones, an assessment registry, and an automated quality suite that keeps the whole thing honest.

## Who this is for

- **Career transitioners** moving into backend software engineering who want a single, self-contained path instead of a scattered tutorial trail.
- **Computer science graduates** who know the theory but have never built, operated, or defended a production-shaped system.
- **Working backend engineers** leveling up in Go who want depth on concurrency, persistence, observability, performance, and deployment rather than surface-level framework tours.
- **Self-directed learners** who prefer runnable code, automated proof, and explicit reasoning over video courses and copy-paste walkthroughs.

You should be comfortable using a terminal. Prior Go experience is helpful but not required — Module 00 assumes nothing.

## What you will become able to do

By graduation you will be able to, without hand-holding:

- **Read and explain Go from first principles** — value vs. reference semantics, stack and heap allocation, interfaces, and the runtime's scheduling model.
- **Design and build backend services** — HTTP APIs, relational persistence with SQL and PostgreSQL, authentication, authorization, and tenant isolation.
- **Write concurrent code that is correct** — goroutines, channels, worker pools, `context` cancellation, and graceful shutdown under real load.
- **Debug unassisted** — reproduce a failure, read a stack trace and a heap profile, and form a hypothesis from evidence rather than guessing.
- **Engineer for performance** — benchmark, profile, and optimize with `pprof` and memory-layout reasoning instead of superstition.
- **Operate and ship software** — structured logging, metrics, tracing, containers, CI/CD, and reproducible deployment runbooks.
- **Defend your work** — explain architecture tradeoffs, respond to code review, and present a portfolio-grade system to a hiring panel.

## Core learning path

The core path runs through 18 modules in order, from orientation to the flagship capstone. Module 17 (Advanced Electives) is an optional companion track.

| # | Module | Phase | What it delivers |
|---|---|---|---|
| 00 | [Orientation](./curriculum/modules/00-orientation/README.md) | orientation | The learning system, workflow, and zero-magic contract |
| 01 | [Computers, Terminal, Git, and the Web](./curriculum/modules/01-computers-terminal-git-web/README.md) | foundations | The machine, shell, Git, and web basics |
| 02 | [Go Setup and Tooling](./curriculum/modules/02-go-setup-tooling/README.md) | tooling | A normal, repeatable Go toolchain |
| 03 | [Programming Fundamentals with Go](./curriculum/modules/03-programming-fundamentals/README.md) | go-core | Core programming without hidden assumptions |
| 04 | [Functions, Errors, and Data Semantics](./curriculum/modules/04-functions-errors-data-semantics/README.md) | go-core | Boundaries, failure, cleanup, and mutation |
| 05 | [Types, Interfaces, Packages, and Modules](./curriculum/modules/05-types-interfaces-packages-modules/README.md) | engineering-core | Maintainable, modular Go programs |
| 06 | [Testing, Debugging, and Refactoring](./curriculum/modules/06-testing-debugging-refactoring/README.md) | engineering-core | Proof, diagnosis, and safe change |
| 07 | [CLI, Files, JSON, Config, and Runtime Coordination](./curriculum/modules/07-cli-files-json-config/README.md) | cli-io | Useful local programs and process lifecycle |
| 08 | [HTTP and REST APIs](./curriculum/modules/08-http-rest-apis/README.md) | backend | The first serious backend service |
| 09 | [SQL and PostgreSQL Persistence](./curriculum/modules/09-sql-postgres-persistence/README.md) | data | Real relational backend work |
| 10 | [Authentication, Authorization, and Security](./curriculum/modules/10-auth-security/README.md) | security | Explicit, secure backend behavior |
| 11 | [Lifecycle, Context, and Concurrency](./curriculum/modules/11-lifecycle-context-concurrency/README.md) | concurrency | Safe concurrent systems and request lifecycle |
| 12 | [Observability and Diagnostics](./curriculum/modules/12-observability-diagnostics/README.md) | reliability | Services you can see into |
| 13 | [Performance and Memory Engineering](./curriculum/modules/13-performance-memory-engineering/README.md) | performance | Profiling and benchmark-driven optimization |
| 14 | [Architecture and Distributed Systems](./curriculum/modules/14-architecture-distributed-systems/README.md) | architecture | Architecture as tradeoff reasoning |
| 15 | [Docker, CI/CD, and Deployment](./curriculum/modules/15-docker-cicd-deployment/README.md) | delivery | Shippable, repeatable software |
| 16 | [Portfolio and Interview Readiness](./curriculum/modules/16-portfolio-interview-readiness/README.md) | career | Competence converted into employability |
| 18 | [Flagship Opslane](./curriculum/modules/18-flagship-opslane/README.md) | flagship | Final graduation proof |

The full catalogue, pacing guide, and workload estimates are in [`curriculum/README.md`](./curriculum/README.md). Visual progression charts are in [`docs/PROGRESSION.md`](./docs/PROGRESSION.md).

## Projects and proof system

Learning is proven, not assumed. Each module ends in a **project capstone** with a starter and a solution, a grading rubric whose weights sum to 100%, and an automatically verified artifact. Alongside the projects sits an **assessment registry** of checkpoints, quizzes, and defenses, each with a README, questions, rubric, and answer key.

Proof is verified mechanically. The curriculum ships a structural validator that checks every item's README, files, links, prerequisite closure, concept ownership, and assessment artifacts, so a lesson cannot quietly drift out of compliance. You can run the same gates a maintainer runs (see below).

The **dual-state starter/solution model** means every project begins as a runnable starter with failing tests and ends as a passing reference solution — so "done" always means "the tests pass", not "it looks finished".

## How to start

### Prerequisites

- **Go:** 1.27+ (declared in [`go.mod`](./go.mod))
- **Git:** standard git CLI
- **Operating system:** Linux, macOS, or Windows (WSL2 or MinGW-w64 recommended for CGO race detection)

### Setup

```bash
# 1. Clone the repository
git clone https://github.com/swe-labs/the-go-engineer.git
cd the-go-engineer

# 2. Download dependencies
go mod download
```

Then open the first lesson and follow it in order:

- [`curriculum/modules/00-orientation/lessons/01-how-to-use-this-repository/README.md`](./curriculum/modules/00-orientation/lessons/01-how-to-use-this-repository/README.md)

## How to validate your progress

The repository includes both a curriculum structural validator and standard Go test suites:

```bash
# 1. Run the full automated curriculum quality gate
go run ./scripts/validate_curriculum.go validate-all

# 2. Run the full unit and integration test suite
go test ./...

# 3. Run static analysis
go vet ./...

# 4. Verify code formatting
gofmt -l .

# 5. Run the race detector on the concurrency and flagship modules
go test -race ./curriculum/modules/11-lifecycle-context-concurrency/... ./curriculum/modules/18-flagship-opslane/...
```

Add `--strict` to the first command to promote every advisory to a hard failure; this is the release gate. Each module's own README lists its project and assessment verification commands.

## Advanced electives

[Module 17: Advanced Electives](./curriculum/electives/advanced-electives/README.md) is an optional, non-blocking specialization track of 26 lessons across five streams — advanced Go typing and metaprogramming, modern RPC and serialization (Protobuf, gRPC), high-throughput concurrency primitives, distributed-systems patterns, and Kubernetes operations. Take it alongside or after the capstone to extend your stack without diluting the core path.

## Flagship Opslane

[Module 18: Flagship Opslane](./curriculum/modules/18-flagship-opslane/README.md) is the graduation capstone: a multi-tenant incident-management SaaS backend that integrates everything from Modules 00–15 — API design, persistence, authentication, tenant isolation, an event bus, bounded worker pools, caching, observability, and graceful shutdown — into one system you run, operate, and defend. The full architecture is specified in [`docs/flagship/OPSLANE_SAAS_BACKEND.md`](./docs/flagship/OPSLANE_SAAS_BACKEND.md).

## Contributing and maintenance

Contributions are welcome. Before opening a pull request, read:

- **[Contribution Guidelines (`CONTRIBUTING.md`)](./CONTRIBUTING.md):** workflow, review expectations, and the authoring standard.
- **[Maintainer Checklist (`MAINTAINER-CHECKLIST.md`)](./MAINTAINER-CHECKLIST.md):** the release gate and verification steps.
- **[Zero-Magic Authoring Specification](./curriculum/zero-magic-authoring-spec.md):** the principles every lesson and project must follow.
- **[System Architecture (`ARCHITECTURE.md`)](./ARCHITECTURE.md):** the canonical v3 architecture and metadata manifests.
- **[Quality Rubrics](./curriculum/quality-rubric.md):** the evaluation criteria for code, proof, and documentation.
- **[Release History (`CHANGELOG.md`)](./CHANGELOG.md):** version history and migration notes.

## License

This project is licensed under the [MIT License](./LICENSE).

- Free for personal learning, forked study, modification, and redistribution.
- Commercial reuse, including paid courses and training programs, is permitted under the same terms.
- Contributions are welcome: submitting a pull request licenses your contribution to the project under MIT.
