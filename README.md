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

## Choose Your Learning Path

New to Go? → **Start at [01-core-foundations](./01-core-foundations/)**  
Know another language? → **See [LEARNING-PATH.md](./LEARNING-PATH.md) for your track**  
Experienced with Go? → **Jump to [§11 Concurrency](./11-concurrency/) or browse [curriculum map](./docs/curriculum/README.md)**

## Reference Documents

| Document | Purpose |
| -------- | ------- |
| [LEARNING-PATH.md](./LEARNING-PATH.md) | **Complete learning guide** with prerequisites and recommended tracks |
| [docs/curriculum/README.md](./docs/curriculum/README.md) | **Curriculum map** showing all 80+ lessons and dependencies |
| [COMMON-MISTAKES.md](./COMMON-MISTAKES.md) | 15 most common Go bugs with fixes and section cross-references |
| [ROADMAP.md](./ROADMAP.md) | What is built, in progress, and planned |
| [CHANGELOG.md](./CHANGELOG.md) | History of additions and bug fixes |
| [CONTRIBUTING.md](./CONTRIBUTING.md) | How to add lessons, exercises, and sections |
| [RELEASE.md](./RELEASE.md) | Release planning and process guide |

## Who is This For?

- **Complete beginners** — Never programmed before? Start at Section 01. Every line is explained.
- **Developers from other languages** — Know Python/JS/Java? Skim the basics, deep-dive into Go-specific patterns.
- **Go developers leveling up** — Already write Go? Jump to Chapters 09+ for concurrency, testing, and production patterns.

## The Elite 15-Chapter Learning Path

This repository follows a strict **Language → Runtime → IO → Quality → Architecture** progression. Every chapter cleanly decouples concepts by their primary domain.

- `01-core-foundations`: Development environment, syntax basics, constants (GS, LB).
- `02-control-flow`: Execution branches, for loops (CF).
- `03-data-structures`: Arrays, slices, maps, pointers (DS).
- `04-functions-and-errors`: First-class closures, defer, error design (FE).
- `05-types-and-interfaces`: Structs, methods, interfaces, generics (TI).
- `06-composition`: Struct embedding (CO).
- `07-strings-and-text`: UTF-8, regex, templates (ST).
- `08-modules-and-packages`: Go modules, versioning (MP).
- `09-io-and-cli`: Filesystem (FS), Encoding (EN), CLI Tools (CL).
- `10-web-and-database`: Databases (DB), Masterclass (WM), HTTP Client (HC).
- `11-concurrency`: Goroutines (GC), Context (CT), Time (TM).
- `12-concurrency-patterns`: errgroup, sync.Pool (CP).
- `13-quality-and-performance`: Testing (TE), Mocking (HM), Profiling (PR).
- `14-application-architecture`: Package Design (PD), Logging (SL), gRPC (GR), Shutdown (GS), Capstone (EC).
- `15-code-generation`: //go:generate (CG).

## Guided Learning Patterns

This curriculum is not just a linear list. It follows specific architectural decisions:

| Pattern | Behavior | Examples |
| ------- | -------- | -------- |
| **Linear** | Strictly A → B → C | Control Flow, Composition, gRPC |
| **Fork/Rejoin** | Independent paths that merge | DS.4 (Pointers) needed by DS.5 only |
| **Deep Exercise** | Synthesis of all previous prereqs | DS.6, FE.9, TI.6, GC.7 |
| **Parallel Track** | Non-blocking concurrent subjects | Race conditions run alongside channels |

## Projects & Exercises

Each module culminates in a hands-on project to test your understanding:

| Chapter | ID | Exercise | Description |
| ------- | --- | -------- | ----------- |
| **01** | LB.4 | `language-basics/4-app-logger` | Synthesise types + iota + Stringer |
| **02** | CF.4 | `4-pricing-calculator` | Map lookups · HasSuffix · switch on bool |
| **03** | DS.6 | `6-contact-manager` | Secondary index · init() · pointer returns |
| **04** | FE.9 | `8-error-handling` | Custom struct error · errors.As · defer |
| **05** | TI.6 | `6-payroll-processor` | Polymorphic slice · embedded interface |
| **06** | CO.3 | `3-bank-account` | Embedded shadowing |
| **07** | ST.6 | `6-log-parser` | Regex · scanner · strings.Builder |
| **09** | FS.7 | `filesystem/7-log-search` | WalkDir filter · Scanner per file |
| **09** | EN.6 | `encoding/6-config-parser` | Decoder · validate() · zero-value detection |
| **09** | CL.4 | `cli-tools/4-file-organizer` | --dry-run · ReadDir · Rename |
| **10** | DB.6  | `databases/6-repository` | SQLite impl · dependency injection |
| **10** | WM.11 | `enterprise-capstone` | Full production-ready server |
| **11** | GC.7  | `concurrency/7-downloader` | Semaphore pattern · WaitGroup + channel |
| **11** | TM.7  | `time-and-scheduling/7-reminder` | AfterFunc · ticker · select |
| **11** | CT.5  | `context/5-timeout-client` | DeadlineExceeded detection |
| **12** | CP.4  | `4-bounded-pipeline` | g.SetLimit · g.TryGo · pooled buffers |
| **12** | CP.5 | `5-url-checker` | Pooled client · sorting by latency |
| **13** | HM.4 | `http-client-testing/6-testify-mock` | .On/.Return · AssertNumberOfCalls |
| **13** | PR.1 | `profiling/1-cpu-profile` | flat vs cum · go tool pprof |
| **14** | SL.5 | `structured-logging/5-exercise` | PII redactor · ReplaceAttr · mapping |
| **14** | GR.2 | `grpc/1-unary` | Generated interface · interceptors |
| **14** | GS.3 | `graceful-shutdown/3-capstone` | Signal → ready=503 → drain → order |

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
