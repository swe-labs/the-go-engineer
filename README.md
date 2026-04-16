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
go run ./01-getting-started/2-hello-world
```

## Public Curriculum (11 Stages)

The curriculum is organized into 11 stages. Each stage has a dedicated entry page under [docs/stages](./docs/stages/README.md).

| Stage | Focus | Sections | Source |
| --- | --- | --- | --- |
<| [01 Getting Started](./docs/stages/01-getting-started.md) | install, hello world, dev setup | s01 | [01-getting-started](./01-getting-started/) |
| [02 Language Basics](./docs/stages/02-language-basics.md) | variables, types, control flow, data structures | s02-s03 | [02-language-basics](./02-language-basics/) |
| [03 Functions & Errors](./docs/stages/03-functions-errors.md) | functions, params, error handling | s04 | [03-functions-errors](./03-functions-errors/) |
| [04 Types & Design](./docs/stages/04-types-design.md) | structs, interfaces, composition, strings | s05-s07 | [04-types-design](./04-types-design/), [06-composition](./06-composition/), [07-strings-and-text](./07-strings-and-text/) |
| [05 Packages & IO](./docs/stages/05-packages-io.md) | modules, packages, CLI, files | s08-s09 | [08-modules-and-packages](./08-modules-and-packages/), [09-io-and-cli](./09-io-and-cli/) |
| [06 Backend & DB](./docs/stages/06-backend-db.md) | HTTP, databases, handlers | s10 | [10-web-and-database](./10-web-and-database/) |
| [07 Concurrency](./docs/stages/07-concurrency.md) | goroutines, channels, patterns | s11-s12 | [11-concurrency](./11-concurrency/), [12-concurrency-patterns](./12-concurrency-patterns/) |
| [08 Quality & Test](./docs/stages/08-quality-test.md) | testing, benchmarking, profiling | s13 | [13-quality-and-performance](./13-quality-and-performance/) |
| [09 Architecture](./docs/stages/09-architecture.md) | package design, services | s14 | [14-application-architecture/package-design](./14-application-architecture/package-design/), [14-application-architecture/grpc](./14-application-architecture/grpc/) |
| [10 Production](./docs/stages/10-production.md) | logging, shutdown, deployment | s14b | [14-application-architecture/structured-logging](./14-application-architecture/structured-logging/), [14-application-architecture/graceful-shutdown](./14-application-architecture/graceful-shutdown/), [14-application-architecture/docker-and-deployment](./14-application-architecture/docker-and-deployment/) |
| [11 Flagship](./docs/stages/11-flagship.md) | full GoScale project | s15-s17 | [14-application-architecture/enterprise-capstone](./14-application-architecture/enterprise-capstone/), [15-code-generation](./15-code-generation/) |

## Best Entry Point By Learner Type

### Complete beginner

Start at [Stage 01: Getting Started](./docs/stages/01-getting-started.md).
That stage page points to the current source surface for beginner setup and first-run work.

### Experienced programmer new to Go

Start at [Stage 02: Language Basics](./docs/stages/02-language-basics.md).
If you want a faster ramp, skim [Stage 01](./docs/stages/01-getting-started.md) first and then use
the stage page to jump into the current source sections.

### Experienced Go developer

Jump to the stage you want to strengthen:

- concurrency: [Stage 07: Concurrency](./docs/stages/07-concurrency.md)
- quality and test: [Stage 08: Quality & Test](./docs/stages/08-quality-test.md)
- architecture: [Stage 09: Architecture](./docs/stages/09-architecture.md)
- production: [Stage 10: Production](./docs/stages/10-production.md)

## What You Will Build

You will work through:

- beginner-friendly exercises and starter projects
- parsers, filesystem tools, and CLI utilities
- HTTP services and database-backed applications
- concurrency pipelines, worker pools, and timeout-aware clients
- profiling, testing, and benchmark-driven improvements
- structured logging, graceful shutdown, gRPC, and deployment-ready workflows

## Current Docs

These docs are still useful:

| Document | Purpose |
| --- | --- |
| [LEARNING-PATH.md](./LEARNING-PATH.md) | current learning guide |
| [docs/stages/README.md](./docs/stages/README.md) | stage entry index |
| [docs/beta-public-architecture.md](./docs/beta-public-architecture.md) | explanation of source vs public architecture |
| [docs/curriculum/README.md](./docs/curriculum/README.md) | section-by-section curriculum map |
| [COMMON-MISTAKES.md](./COMMON-MISTAKES.md) | common Go bugs and fixes |
| [ROADMAP.md](./ROADMAP.md) | public roadmap |
| [CHANGELOG.md](./CHANGELOG.md) | change history |
| [CONTRIBUTING.md](./CONTRIBUTING.md) | contribution workflow |
| [RELEASE.md](./RELEASE.md) | release process |

Maintainers and contributors can follow planning on
[`planning/v2`](https://github.com/rasel9t6/the-go-engineer/tree/planning/v2/docs/v2).

## Current Repo Reality

- the 11-stage model is the public navigation truth
- the current section folders remain the source inventory
- content is being migrated to `01-foundations/` incrementally

That means the README now presents the curriculum by stage even though the physical folder
layout is still migrating to the new format.

For the full explanation of that relationship, see
[docs/beta-public-architecture.md](./docs/beta-public-architecture.md).

## Run Lessons And Exercises

```bash
# run a lesson
go run ./01-getting-started/2-hello-world

# run a starter exercise
go run ./03-functions-errors/7-order-summary/_starter

# compare with the canonical solution
go run ./03-functions-errors/7-order-summary
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
