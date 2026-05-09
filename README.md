# The Go Engineer

[![CI](https://github.com/rasel9t6/the-go-engineer/actions/workflows/ci.yml/badge.svg)](https://github.com/rasel9t6/the-go-engineer/actions)
[![License: TGE License v1.0 (Non-Commercial)](https://img.shields.io/badge/License-TGE_v1.0-red.svg?style=for-the-badge&logo=data:image/svg%2Bxml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0id2hpdGUiPjxwYXRoIGQ9Ik0xNCAySDZhMiAyIDAgMCAwLTIgMnYxNmEyIDIgMCAwIDAgMiAyaDEyYTIgMiAwIDAgMCAyLTJWOGwtNi02em0tMSAxLjVMMTguNSA5SDEzVjMuNXpNNiAyMFY0aDV2Nmg2djEwSDZ6bTItOGg4djJIOHYtMnptMCA0aDh2Mkg4di0yeiIgZmlsbC1ydWxlPSJldmVub2RkIiBjbGlwLXJ1bGU9ImV2ZW5vZGQiLz48L3N2Zz4%3D)](#license)
[![GitHub Sponsors](https://img.shields.io/badge/sponsor-30363D?style=for-the-badge&logo=GitHub-Sponsors&logoColor=%23EA4AAA)](https://github.com/sponsors/rasel9t6)
[![Patreon](https://img.shields.io/badge/Patreon-F96854?style=for-the-badge&logo=patreon&logoColor=white)](https://patreon.com/rasel9t6)

Learn Go by building, testing, and operating real software.

The Go Engineer is a repo-first Go software engineering learning system built around a progressive **5-Phase, 12-Section architecture** (v2.1).

It takes you from absolute beginner (how computers work) to building, testing, and deploying a production-grade SaaS backend ("GoScale").

For the full curriculum architecture, read [ARCHITECTURE.md](./ARCHITECTURE.md).

## Start Here

Pick the release channel that matches what you want:

- `release/v1`: the stable v1 line for learners who want the current supported experience
- `release/v2`: the current public beta checkpoint for the v2 line
- `main`: the active v2 implementation branch — fastest-moving line

## Quick Start

```bash
git clone https://github.com/rasel9t6/the-go-engineer.git
cd the-go-engineer
go mod download
go version
go run ./00-how-computers-work/1-what-is-a-program
```

## Curriculum Overview (5 Phases, 12 Sections)

See [ARCHITECTURE.md](./ARCHITECTURE.md) for the exact lesson breakdowns.

```text
Phase 0: Machine Foundation     s00  How Computers Work         (0% → 5%)
Phase 1: Language Foundation    s01  Getting Started             (5% → 12%)
                                s02  Language Basics             (12% → 28%)
                                s03  Functions & Errors          (28% → 38%)
                                s04  Types & Design              (38% → 52%)
Phase 2: Engineering Core       s05  Packages, I/O & CLI        (52% → 62%)
                                s06  Backend, APIs & Databases   (62% → 75%)
                                s07  Concurrency                 (75% → 83%)
                                s08  Quality & Testing           (83% → 87%)
Phase 3: Systems Engineering    s09  Architecture & Security     (87% → 92%)
                                s10  Production Operations       (92% → 96%)
Phase 4: Flagship Project       s11  GoScale SaaS Backend        (96% → 100%)
```

## What You Will Build

You will work through:

- beginner-friendly exercises and starter projects
- parsers, filesystem tools, and CLI utilities
- HTTP services and database-backed applications
- concurrency pipelines, worker pools, and timeout-aware clients
- profiling, testing, and benchmark-driven improvements
- structured logging, graceful shutdown, gRPC, and deployment-ready workflows
- **GoScale**: A full production-grade SaaS backend (flagship capstone)

## Current Docs

These docs guide you through the learning system:

| Document                                              | Purpose                                      |
| ----------------------------------------------------- | -------------------------------------------- |
| [ARCHITECTURE.md](./ARCHITECTURE.md)                  | Curriculum source of truth (v2.1)            |
| [ROADMAP.md](./ROADMAP.md)                            | What is built, in progress, and planned      |
| [LEARNING-PATH.md](./LEARNING-PATH.md)                | Learning guide and entry points              |
| [CURRICULUM-BLUEPRINT.md](./CURRICULUM-BLUEPRINT.md)  | Teaching contract and lesson delivery rules   |
| [CODE-STANDARDS.md](./CODE-STANDARDS.md)              | Engineering and code style standards         |
| [TESTING-STANDARDS.md](./TESTING-STANDARDS.md)        | Testing coverage and patterns                |
| [COMMON-MISTAKES.md](./COMMON-MISTAKES.md)            | 15 common Go bugs and fixes                  |
| [CONTRIBUTING.md](./CONTRIBUTING.md)                  | Contribution workflow                        |

## Run Lessons and Exercises

```bash
# Run a lesson
go run ./00-how-computers-work/1-what-is-a-program

# Run tests
go test ./...
```

## Validate the Repo

```bash
go run ./scripts/validate_curriculum.go
make test
make lint
```

## Windows Note

Some database and capstone paths use `go-sqlite3`, which requires CGO and a C compiler.
If you are on Windows, WSL2 is the smoothest setup for those paths.

## License

This project is licensed under the **The Go Engineer License (TGE License) v1.0**.

- Free for personal, educational, and non-commercial use
- Commercial use requires permission
