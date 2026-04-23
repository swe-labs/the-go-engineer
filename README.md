# The Go Engineer

[![CI](https://github.com/rasel9t6/the-go-engineer/actions/workflows/ci.yml/badge.svg)](https://github.com/rasel9t6/the-go-engineer/actions)
[![License: TGE License v1.0 (Non-Commercial)](https://img.shields.io/badge/License-TGE_v1.0-red.svg?style=for-the-badge&logo=data:image/svg%2Bxml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0id2hpdGUiPjxwYXRoIGQ9Ik0xNCAySDZhMiAyIDAgMCAwLTIgMnYxNmEyIDIgMCAwIDAgMiAyaDEyYTIgMiAwIDAgMCAyLTJWOGwtNi02em0tMSAxLjVMMTguNSA5SDEzVjMuNXpNNiAyMFY0aDV2Nmg2djEwSDZ6bTItOGg4djJIOHYtMnptMCA0aDh2Mkg4di0yeiIgZmlsbC1ydWxlPSJldmVub2RkIiBjbGlwLXJ1bGU9ImV2ZW5vZGQiLz48L3N2Zz4%3D)](#license)
[![GitHub Sponsors](https://img.shields.io/badge/sponsor-30363D?style=for-the-badge&logo=GitHub-Sponsors&logoColor=%23EA4AAA)](https://github.com/sponsors/rasel9t6)
[![Patreon](https://img.shields.io/badge/Patreon-F96854?style=for-the-badge&logo=patreon&logoColor=white)](https://patreon.com/rasel9t6)

Learn Go by building, testing, and operating real software.

The Go Engineer is a repo-first Go software engineering learning system built around a progressive **5-phase, 12-section architecture** (v2.1).

It takes a learner from machine fundamentals and first execution all the way to a production-shaped Go backend capstone.

For the full curriculum source of truth, read [ARCHITECTURE.md](./ARCHITECTURE.md).

## Current Status

The v2.1 stable release is shipped.

- `release/v1`: stable v1 maintenance line
- `main`: post-v2.1 implementation line
- `release/v2`: reserved for v2.1.x fixes and stable-line maintenance

## Quick Start

```bash
git clone https://github.com/rasel9t6/the-go-engineer.git
cd the-go-engineer
go mod download
go version
go run ./00-how-computers-work/1-what-is-a-program
```

## Curriculum Overview

```text
Phase 0: Machine Foundation     s00  How Computers Work          (0% -> 5%)
Phase 1: Language Foundation    s01  Getting Started             (5% -> 12%)
                                s02  Language Basics             (12% -> 28%)
                                s03  Functions & Errors          (28% -> 38%)
                                s04  Types & Design              (38% -> 52%)
Phase 2: Engineering Core       s05  Packages, I/O & CLI         (52% -> 62%)
                                s06  Backend, APIs & Databases   (62% -> 75%)
                                s07  Concurrency                 (75% -> 83%)
                                s08  Quality & Testing           (83% -> 87%)
Phase 3: Systems Engineering    s09  Architecture & Security     (87% -> 92%)
                                s10  Production Operations       (92% -> 96%)
Phase 4: Flagship Project       s11  Opslane SaaS Backend        (96% -> 100%)
```

## Source Sections

| Section | Folder | Focus |
| --- | --- | --- |
| s00 | [00-how-computers-work](./00-how-computers-work) | execution, memory, terminal, processes |
| s01 | [01-getting-started](./01-getting-started) | install, first run, toolchain basics |
| s02 | [02-language-basics](./02-language-basics) | values, control flow, data structures |
| s03 | [03-functions-errors](./03-functions-errors) | reusable logic and explicit failure handling |
| s04 | [04-types-design](./04-types-design) | structs, interfaces, composition, strings and text |
| s05 | [05-packages-io](./05-packages-io) | modules, packages, CLI, encoding, filesystem |
| s06 | [06-backend-db](./06-backend-db) | HTTP servers, API design, gRPC, databases |
| s07 | [07-concurrency](./07-concurrency) | goroutines, context, sync, pipelines |
| s08 | [08-quality-test](./08-quality-test) | tests, profiling, benchmarks |
| s09 | [09-architecture](./09-architecture) | package design, architecture patterns, security |
| s10 | [10-production](./10-production) | logging, shutdown, config, observability, deployment |
| s11 | [11-flagship](./11-flagship) | integrated capstone system |

## What You Will Build

You will work through:

- beginner-friendly exercises and starter projects
- parsers, filesystem tools, and CLI utilities
- HTTP services and database-backed applications
- concurrency pipelines, worker pools, and timeout-aware clients
- profiling, testing, and benchmark-driven improvements
- structured logging, graceful shutdown, configuration, observability, and deployment workflows
- **Opslane**: a production-shaped SaaS backend capstone

## Core Docs

| Document | Purpose |
| --- | --- |
| [ARCHITECTURE.md](./ARCHITECTURE.md) | curriculum source of truth |
| [CURRICULUM-BLUEPRINT.md](./CURRICULUM-BLUEPRINT.md) | lesson and milestone teaching contract |
| [LEARNING-PATH.md](./LEARNING-PATH.md) | entry points and path guidance |
| [ROADMAP.md](./ROADMAP.md) | beta completion status and RC focus |
| [docs/PROGRESSION.md](./docs/PROGRESSION.md) | visual progression through phases and milestones |
| [CODE-STANDARDS.md](./CODE-STANDARDS.md) | code and explanation standards |
| [TESTING-STANDARDS.md](./TESTING-STANDARDS.md) | testing patterns and verification expectations |
| [COMMON-MISTAKES.md](./COMMON-MISTAKES.md) | common Go bugs and fixes |
| [CONTRIBUTING.md](./CONTRIBUTING.md) | contribution workflow |
| [RELEASE.md](./RELEASE.md) | release and stabilization process |

## Validate the Repo

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
