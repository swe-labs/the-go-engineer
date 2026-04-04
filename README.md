# The Go Engineer: Learn Go by Building Real Projects

[![CI](https://github.com/rasel9t6/the-go-engineer/actions/workflows/ci.yml/badge.svg)](https://github.com/rasel9t6/the-go-engineer/actions)
[![License: TGE License v1.0 (Non-Commercial)](https://img.shields.io/badge/License-TGE_v1.0-red.svg?style=for-the-badge&logo=data:image/svg%2Bxml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0id2hpdGUiPjxwYXRoIGQ9Ik0xNCAySDZhMiAyIDAgMCAwLTIgMnYxNmEyIDIgMCAwIDAgMiAyaDEyYTIgMiAwIDAgMCAyLTJWOGwtNi02em0tMSAxLjVMMTguNSA5SDEzVjMuNXpNNiAyMFY0aDV2Nmg2djEwSDZ6bTItOGg4djJIOHYtMnptMCA0aDh2Mkg4di0yeiIgZmlsbC1ydWxlPSJldmVub2RkIiBjbGlwLXJ1bGU9ImV2ZW5vZGQiLz48L3N2Zz4%3D)](#-license)
[![GitHub Sponsors](https://img.shields.io/badge/sponsor-30363D?style=for-the-badge&logo=GitHub-Sponsors&logoColor=#EA4AAA)](https://github.com/sponsors/rasel9t6)
[![Patreon](https://img.shields.io/badge/Patreon-F96854?style=for-the-badge&logo=patreon&logoColor=white)](https://patreon.com/rasel9t6)

Welcome to **The Go Engineer** — the definitive open-source Go curriculum. Every section teaches through **practical examples, real-world projects, and hands-on exercises** — not just syntax. You'll build servers, CLI tools, concurrent pipelines, REST APIs, and production-grade applications while learning the engineering depth behind every concept.

## Quick Start

```bash
# 1. Install Go: https://go.dev/dl/  (see 01-core-foundations/getting-started)
# 2. Clone this repository
git clone https://github.com/rasel9t6/the-go-engineer.git
cd the-go-engineer

# 3. Download dependencies
go mod download

# 4. Verify Go is working
go version

# 5. Run your first program
go run ./01-core-foundations/getting-started/2-hello-world
```

## Reference Documents

| Document | Purpose |
| -------- | ------- |
| [COMMON-MISTAKES.md](./COMMON-MISTAKES.md) | 15 most common Go bugs with fixes and section cross-references |
| [ROADMAP.md](./ROADMAP.md) | What is built, in progress, and planned |
| [CHANGELOG.md](./CHANGELOG.md) | History of additions and bug fixes |
| [CONTRIBUTING.md](./CONTRIBUTING.md) | How to add lessons, exercises, and sections |

## Who is This For?

- **Complete beginners** — Never programmed before? Start at Section 01. Every line is explained.
- **Developers from other languages** — Know Python/JS/Java? Skim the basics, deep-dive into Go-specific patterns.
- **Go developers leveling up** — Already write Go? Jump to Chapters 09+ for concurrency, testing, and production patterns.

## The Elite 15-Chapter Learning Path

This repository follows a strict **Language → Runtime → IO → Quality → Architecture** progression. Every chapter cleanly decouples concepts by their primary domain.

- `01-core-foundations`: Development environment, syntax basics, constants, variable types.
- `02-control-flow`: Execution branches, for loops, line-of-sight principle.
- `03-data-structures`: Deep dives into arrays, slices, maps, pointers, and heap allocations.
- `04-functions-and-errors`: First-class closures, defer mechanics, and idiomatic error design.
- `05-types-and-interfaces`: Structs, methods, and duck-typing abstractions.
- `06-composition`: Achieving flexible architecture through struct embedding.
- `07-strings-and-text`: UTF-8 internals, regex buffers, and templates.
- `08-modules-and-packages`: Go modules, versioning, dependency tracking.
- `09-io-and-cli`: Standard inputs, filesystem traversal, encoding formats (JSON/XML), and flag-driven CLI tools.
- `10-web-and-database`: SQL migrations, HTTP REST Web Servers, Routing, and raw HTTP Client executions.
- `11-concurrency`: Goroutine mechanics, waitgroups, channel syncing, time/scheduling logic, context lifetimes.
- `12-concurrency-patterns`: High efficiency fan-out, errgroup bounds, and zero-allocation sync.Pool.
- `13-quality-and-performance`: Unit testing arrays, mocking external interfaces, profiling flamegraphs, benchmarking.
- `14-application-architecture`: Full scale project layers, Docker orchestration, Structured Logging, gRPC protocols.
- `15-code-generation`: Abstracted tooling via `//go:generate`.

## Projects & Exercises

Each module culminates in a hands-on project to test your understanding:

| Chapter | Exercise | Description |
| ------- | -------- | ----------- |
| **01** Core Foundations | `language-basics/4-application-logger` | Application Logger with severity levels |
| **02** Control Flow | `4-pricing-calculator` | Pricing Calculator engine |
| **03** Data Structures | `6-contact-manager` | Slice-based Contact Manager System |
| **04** Functions & Errors | `8-error-handling` | Custom mathematical error handling |
| **05** Types & Interfaces | `6-payroll-processor` | Polymorphic User Payroll Processor |
| **06** Composition | `3-bank-account` | Bank Account System with deposits/withdrawals |
| **07** Strings & Text | `6-log-parser` | Log File Parsing System |
| **09** IO and CLI | `cli-tools/4-file-organizer` | CLI file organizer by extension |
| **09** IO and CLI | `filesystem/7-log-search` | Directory traversal log search tool |
| **09** IO and CLI | `encoding/6-config-parser` | JSON config file parser with validation |
| **10** Web & Database | `databases/6-repository` | CRUD SQLite App using Repository Pattern |
| **10** Web & Database | `web-masterclass/1-routing/exercise` | Multi-route Bookstore Web API |
| **11** Concurrency | `concurrency/7-downloader` | Concurrent Multi-File Downloader |
| **11** Concurrency | `time-and-scheduling/7-reminder` | Console reminder with countdown timer |
| **11** Concurrency | `context/5-timeout-client` | Timeout-aware HTTP API client |
| **12** Concurrency Patterns | `4-bounded-pipeline-exercise` | Image resizer with bounded concurrency via `errgroup` |
| **12** Concurrency Patterns | `5-url-checker-exercise` | URL health checker with zero-alloc pooled clients |
| **13** Quality & Performance | `http-client-testing/6-testify-mock` | Mocking an external REST API Data Fetcher |
| **13** Quality & Performance | `profiling/1-cpu-profile` | Profile slow vs fast log processor |
| **14** Application Architecture | `enterprise-capstone/cmd/api` | **The Multi-Package Docker Enterprise Backend** |
| **14** Application Architecture | `structured-logging/2-context-logger` | HTTP middleware request-scoped logger extraction |
| **14** Application Architecture | `grpc/1-unary` | Type-safe OrderService client/server via Interceptors |
| **14** Application Architecture | `graceful-shutdown/3-capstone` | Complete readiness → HTTP drain → generic shutdown |

## How to Use This Repository

The best way to learn is by **reading the inline comments** and **running the code**.

```bash
# Run any lesson
go run ./CHAPTER/LESSON

# Examples:
go run ./01-core-foundations/getting-started/2-hello-world
go run ./01-core-foundations/language-basics/1-variables
go run ./11-concurrency/concurrency/3-channels
go run ./10-web-and-database/web-masterclass/1-routing

# Chapter 13 — Quality and Testing
go run ./13-quality-and-performance/profiling/1-cpu-profile
```

### Self-Challenge Mode

Most exercises include a `_starter/` directory with TODO stubs:

```bash
# Try the exercise yourself first:
go run ./02-control-flow/4-pricing-calculator/_starter

# Then compare with the solution:
go run ./02-control-flow/4-pricing-calculator
```

For the grand finale, boot the entire Enterprise Backend cluster (Database + Migrations + API) using Docker:

```bash
# Run the massive Chapter 14 Capstone project
cd 14-application-architecture/enterprise-capstone
docker-compose up -d --build
```

## Running the Tests

To verify your environment is set up correctly, run the full test suite over the entire domain structure:

```bash
# Run all tests
go test ./...

# Run tests with race detection
go test -race ./...
```

## Windows Users — CGO Note

Chapter 10 (`databases`) and Chapter 14 (`enterprise-capstone`) use `go-sqlite3`, which requires CGO and a C compiler. On Windows without WSL:

1. Install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/)
2. Set environment: `$env:CGO_ENABLED = "1"` (PowerShell)

We recommend [WSL2](https://docs.microsoft.com/en-us/windows/wsl/) for the best experience.

## 📜 License

This project is licensed under the **The Go Engineer License (TGE License) v1.0**.

- ✅ Free for personal, educational, and non-commercial use
- ❌ Commercial use is strictly prohibited without permission
