# The Go Engineer

[![CI](https://github.com/rasel9t6/the-go-engineer/actions/workflows/ci.yml/badge.svg)](https://github.com/rasel9t6/the-go-engineer/actions)
[![License: TGE License v1.0 (Non-Commercial)](https://img.shields.io/badge/License-TGE_v1.0-red.svg?style=for-the-badge&logo=data:image/svg%2Bxml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0id2hpdGUiPjxwYXRoIGQ9Ik0xNCAySDZhMiAyIDAgMCAwLTIgMnYxNmEyIDIgMCAwIDAgMiAyaDEyYTIgMiAwIDAgMCAyLTJWOGwtNi02em0tMSAxLjVMMTguNSA5SDEzVjMuNXpNNiAyMFY0aDV2Nmg2djEwSDZ6bTItOGg4djJIOHYtMnptMCA0aDh2Mkg4di0yeiIgZmlsbC1ydWxlPSJldmVub2RkIiBjbGlwLXJ1bGU9ImV2ZW5vZGQiLz48L3N2Zz4%3D)](#license)
[![GitHub Sponsors](https://img.shields.io/badge/sponsor-30363D?style=for-the-badge&logo=GitHub-Sponsors&logoColor=%23EA4AAA)](https://github.com/sponsors/rasel9t6)
[![Patreon](https://img.shields.io/badge/Patreon-F96854?style=for-the-badge&logo=patreon&logoColor=white)](https://patreon.com/rasel9t6)

Learn Go by building, testing, and operating real software.

The Go Engineer is a repo-first Go software engineering learning system. The curriculum is moving
from an alpha-era section inventory to a beta stage model that is easier for learners to navigate
from zero to production-minded engineering work.

The beta public architecture is now the main learner-facing direction. The source content still
lives in the current section folders while the beta shell is rolled out incrementally.

If you want the public explanation of that transition model, read
[docs/beta-public-architecture.md](./docs/beta-public-architecture.md).

## Start Here

Pick the release channel that matches what you want:

- `release/v1`: the stable v1 line for learners who want the current supported experience
- `release/v2`: the current public alpha checkpoint for the v2 line
- `main`: the active beta implementation branch and the fastest-moving line
- `v2.0.0-alpha.1`: the first public v2 alpha tag

If you want the safest path today, use `release/v1`.
If you want the current public v2 checkpoint, use `release/v2` or the `v2.0.0-alpha.1` tag.
If you want to follow the beta rollout as it lands, use `main`.

## Quick Start

```bash
git clone https://github.com/rasel9t6/the-go-engineer.git
cd the-go-engineer
go mod download
go version
go run ./01-foundations/01-getting-started/2-hello-world
```

## Beta Public Curriculum

The beta curriculum is organized by engineering stage.
Each stage now has a dedicated public entry page under [docs/stages](./docs/stages/README.md).

| Beta stage | Focus | Current source content |
| --- | --- | --- |
| [0 Foundation](./docs/stages/00-foundation.md) | tools, execution, terminal confidence, first-run mental models | [01-foundations/01-getting-started](./01-foundations/01-getting-started/) |
| [1 Language Fundamentals](./docs/stages/01-language-fundamentals.md) | syntax, control flow, data structures, functions, errors | [01-core-foundations/language-basics](./01-core-foundations/language-basics/), [01-foundations/03-control-flow](./01-foundations/03-control-flow/), [01-foundations/04-data-structures](./01-foundations/04-data-structures/), [01-foundations/05-functions-and-errors](./01-foundations/05-functions-and-errors/) |
| [2 Types and Design](./docs/stages/02-types-and-design.md) | structs, interfaces, composition, text and data modeling | [05-types-and-interfaces](./05-types-and-interfaces/), [06-composition](./06-composition/), [07-strings-and-text](./07-strings-and-text/) |
| [3 Modules and IO](./docs/stages/03-modules-and-io.md) | packages, modules, encoding, filesystems, CLI boundaries | [08-modules-and-packages](./08-modules-and-packages/), [09-io-and-cli](./09-io-and-cli/) |
| [4 Backend Engineering](./docs/stages/04-backend-engineering.md) | HTTP, databases, handlers, application boundaries | [10-web-and-database](./10-web-and-database/) |
| [5 Concurrency System](./docs/stages/05-concurrency-system.md) | goroutines, context, scheduling, bounded patterns | [11-concurrency](./11-concurrency/), [12-concurrency-patterns](./12-concurrency-patterns/) |
| [6 Quality and Performance](./docs/stages/06-quality-and-performance.md) | testing, benchmarking, profiling, trust in code | [13-quality-and-performance](./13-quality-and-performance/) |
| [7 Architecture](./docs/stages/07-architecture.md) | package design, service structure, system boundaries | [14-application-architecture/package-design](./14-application-architecture/package-design/), [14-application-architecture/grpc](./14-application-architecture/grpc/) |
| [8 Production Engineering](./docs/stages/08-production-engineering.md) | logging, shutdown, deployment, operating software | [14-application-architecture/structured-logging](./14-application-architecture/structured-logging/), [14-application-architecture/graceful-shutdown](./14-application-architecture/graceful-shutdown/), [14-application-architecture/docker-and-deployment](./14-application-architecture/docker-and-deployment/) |
| [9 Expert Layer](./docs/stages/09-expert-layer.md) | review, diagnosis, anti-patterns, trade-offs | beta additions planned |
| [10 Flagship Project](./docs/stages/10-flagship-project.md) | one long-running proof-of-skill product spine | [14-application-architecture/enterprise-capstone](./14-application-architecture/enterprise-capstone/) plus milestone projects across the repo |
| [11 Code Generation](./docs/stages/11-code-generation.md) | generation after understanding the code it produces | [15-code-generation](./15-code-generation/) |

## Best Entry Point By Learner Type

### Complete beginner

Start at [0 Foundation](./docs/stages/00-foundation.md).
That stage page points to the current source surface for beginner setup and first-run work.

### Experienced programmer new to Go

Start at [1 Language Fundamentals](./docs/stages/01-language-fundamentals.md).
If you want a faster ramp, skim [0 Foundation](./docs/stages/00-foundation.md) first and then use
the stage page to jump into the current source sections.

### Experienced Go developer

Jump to the stage you want to strengthen:

- concurrency: [5 Concurrency System](./docs/stages/05-concurrency-system.md)
- quality and performance: [6 Quality and Performance](./docs/stages/06-quality-and-performance.md)
- architecture: [7 Architecture](./docs/stages/07-architecture.md)
- production and operations: [8 Production Engineering](./docs/stages/08-production-engineering.md)

## What You Will Build

You will work through:

- beginner-friendly exercises and starter projects
- parsers, filesystem tools, and CLI utilities
- HTTP services and database-backed applications
- concurrency pipelines, worker pools, and timeout-aware clients
- profiling, testing, and benchmark-driven improvements
- structured logging, graceful shutdown, gRPC, and deployment-ready workflows

## Current Docs

These docs are still useful during the beta shell rollout:

| Document | Purpose |
| --- | --- |
| [LEARNING-PATH.md](./LEARNING-PATH.md) | current learning guide during the beta routing transition |
| [docs/stages/README.md](./docs/stages/README.md) | beta stage entry index and stage-by-stage public routing |
| [docs/beta-public-architecture.md](./docs/beta-public-architecture.md) | explanation of alpha source inventory versus beta public architecture |
| [docs/curriculum/README.md](./docs/curriculum/README.md) | alpha source inventory and section-by-section curriculum map |
| [COMMON-MISTAKES.md](./COMMON-MISTAKES.md) | common Go bugs and fixes |
| [ROADMAP.md](./ROADMAP.md) | public roadmap and release direction |
| [CHANGELOG.md](./CHANGELOG.md) | change history |
| [CONTRIBUTING.md](./CONTRIBUTING.md) | contribution workflow and repo rules |
| [RELEASE.md](./RELEASE.md) | release process |

Maintainers and contributors can follow the beta planning set on
[`planning/v2`](https://github.com/rasel9t6/the-go-engineer/tree/planning/v2/docs/v2).

## Current Repo Reality

During the beta rollout:

- the beta stage model is the public navigation truth
- the current section folders remain the source inventory
- some stages regroup content from multiple alpha sections
- some alpha sections split across more than one beta stage
- stage entry docs and deeper beta routing are being added incrementally

That means the README now presents the curriculum by beta stage even though the physical folder
layout is still mostly the alpha-era section structure.

For the full explanation of that relationship, see
[docs/beta-public-architecture.md](./docs/beta-public-architecture.md).

## Run Lessons And Exercises

```bash
# run a lesson
go run ./01-foundations/01-getting-started/2-hello-world

# run a starter exercise
go run ./01-foundations/05-functions-and-errors/7-order-summary/_starter

# compare with the canonical solution
go run ./01-foundations/05-functions-and-errors/7-order-summary
```

## Validate The Repo

```bash
go run ./scripts/validate_curriculum.go
go test ./...
```

## Windows Note

Some database and capstone paths use `go-sqlite3`, which requires CGO and a C compiler.
If you are on Windows, WSL2 is the smoothest setup for those paths.

## License

This project is licensed under the **The Go Engineer License (TGE License) v1.0**.

- Free for personal, educational, and non-commercial use
- Commercial use requires permission
