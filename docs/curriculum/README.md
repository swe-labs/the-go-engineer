# The Go Engineer Curriculum Map

Welcome to the complete learning path for **The Go Engineer**. This directory contains the public
curriculum map plus visual references that help learners and contributors understand how the repo
fits together.

## Curriculum Overview

The Go Engineer teaches **15 major sections** across **80+ lessons** through a deliberate
progression:

**Language → Runtime → I/O → Quality → Architecture**

### Section Breakdown

#### Core Foundations (Sections 01-03)

Master the language fundamentals in 15 lessons.

- **§01 — Core Foundations** (8 lessons)
  - Getting Started: installation, hello world, compilation model, dev tools (`GS.1–GS.4`)
  - Language Basics: variables, constants, iota, application-logger exercise (`LB.1–LB.4`)
- **§02 — Control Flow** (4 lessons)
  - for-loop, if/else, switch, pricing-calculator exercise (`CF.1–CF.4`)
- **§03 — Data Structures** (6 lessons)
  - arrays, slices, maps, pointers, advanced slicing, contact-manager exercise (`DS.1–DS.6`)

#### Functions and Errors (Section 04)

Error-handling design patterns in 9 lessons.

- **§04 — Functions & Errors** (9 lessons)
  - functions, closures, variadic functions, multiple returns
  - custom errors, error wrapping, defer, panic/recover, error-handling exercise (`FE.1–FE.9`)

#### Types and Design (Sections 05-07)

Object-oriented and data-modeling patterns in 11 lessons.

- **§05 — Types & Interfaces** (6 lessons)
  - structs, methods, interfaces, Stringer, generics, payroll-processor exercise (`TI.1–TI.6`)
- **§06 — Composition** (2 lessons)
  - composition vs inheritance, embedding (`CO.1–CO.2`)
- **§07 — Strings & Text** (5 core lessons plus the live v2 milestone exercise)
  - string internals, formatting, unicode/runes, regex, text templates (`ST.1–ST.5`)

#### Modules and I/O (Sections 08-09)

Real-world input/output and module management in 18 lessons.

- **§08 — Modules & Packages** (3 lessons)
  - module basics, managing dependencies, versioning (`MP.1–MP.3`)
- **§09 — I/O & CLI** (15 lessons)
  - Filesystem (8): files, paths, directories, temp files, embed, I/O patterns, `fs.FS` testing (`FS.1–FS.8`)
  - Encoding (5): JSON marshalling/unmarshalling, streaming, base64 (`EN.1–EN.5`)
  - CLI Tools (3): args, flags, subcommands (`CL.1–CL.3`)

#### Web and Database (Section 10)

Full-stack development in 19 lessons.

- **§10 — Web & Database** (19 lessons)
  - Databases (5): connecting, INSERT, SELECT, prepared statements, transactions (`DB.1–DB.5`)
  - Web Masterclass (11): routing, DI, templates, middleware, sessions, auth, forms, CRUD, pagination, comments (`WM.1–WM.11`)
  - HTTP Client (2): basic GET, refactoring for testability (`HC.1–HC.2`)

#### Concurrency (Sections 11-12)

Parallel programming patterns in 23 lessons.

- **§11 — Concurrency** (17 lessons)
  - Goroutines & Channels (10): goroutines, WaitGroups, channels, buffered channels, closing, pipelines, races, select, sync primitives (`GC.1–GC.10`)
  - Context (4): Background/TODO, WithCancel, WithTimeout, WithValue (`CT.1–CT.4`)
  - Time & Scheduling (7): time basics, formatting, timers/tickers, random, scheduler, timezones (`TM.1–TM.7`)
- **§12 — Concurrency Patterns** (6 lessons)
  - errgroup basics, errgroup + context, sync.Pool, bounded pipeline, URL checker, worker pool (`CP.1–CP.6`)

#### Quality and Performance (Section 13)

Testing and optimization in 14 lessons.

- **§13 — Quality & Performance** (14 lessons)
  - Testing (4): unit tests, table-driven, HTTP handler tests, benchmarking (`TE.1–TE.4`)
  - Mocking (4): manual mocks, function injection, table-driven mocks, testify/mock (`HM.1–HM.4`)
  - Profiling (2): CPU profiling, live `pprof` endpoint (`PR.1–PR.2`)

#### Architecture (Section 14)

Production-grade design patterns in 15 lessons.

- **§14 — Application Architecture** (15 lessons)
  - Package Design (3): naming, visibility, project layout (`PD.1–PD.3`)
  - Docker (3): single-stage Dockerfile, multi-stage builds, layer caching (`DO.1–DO.3`)
  - Logging (4): slog basics, context-keyed logger, custom handlers, zerolog comparison (`SL.1–SL.4`)
  - gRPC (3): proto definition, unary server, unary client (`GR.1–GR.3`)
  - Graceful Shutdown (2): signal.NotifyContext, HTTP graceful drain (`GS.1–GS.2`)
  - Enterprise Capstone (1): full REST API with PostgreSQL and Docker Compose (`EC.1`)

#### Code Generation (Section 15)

Build-time automation in 1 lesson.

- **§15 — Code Generation** (1 lesson)
  - `go:generate` plus canonical tools such as mockery, stringer, and sqlc (`CG.1`)

## Curriculum Statistics

| Metric | Count |
| --- | --- |
| Total Sections | 15 |
| Total Lessons | 80+ |
| Entry Points | 15 |
| Exercises | 25+ |
| Prerequisite Links | 200+ |
| Total Concepts | 250+ |

## Visual Dependency Graph

See [dependency-graphs.html](./dependency-graphs.html) for an interactive visual map of lesson
dependencies.

How to use it:

1. Open `dependency-graphs.html` in a browser.
2. Click a lesson to inspect its prerequisites and dependents.
3. Follow the graph edges to understand the learning path.

## How to Navigate

### If you're a complete beginner

Start at **§01 Getting Started** (`GS.1`) and move in order. The curriculum is designed so each
lesson builds on previous concepts.

### If you know another language

Skim **§01-03**, go deep on **§04-07**, then move into **§09+** for real applications.

### If you're an experienced Go developer

Jump directly to:

- **§11-12** for advanced concurrency
- **§13** for testing patterns
- **§14** for production architecture

### If you're learning a specific topic

Use the dependency graph to find related lessons across topics like databases, web, testing,
concurrency, logging, and Docker.

## Curriculum Validation

Run the validator to ensure the curriculum still matches the filesystem:

```bash
go run ./scripts/validate_curriculum.go
```

This verifies:

- all lessons in `curriculum.json` have corresponding directories
- no orphaned directories lack curriculum entries
- lesson run/test command paths stay consistent

## Lesson Structure

Each lesson directory follows a consistent shape:

```text
NN-lesson-name/
├── main.go
├── README.md
└── _starter/
```

- `main.go`: primary lesson code with full explanation
- `README.md`: optional additional context for complex lessons
- `_starter/`: optional starting point for follow-along exercises

### `main.go` Header Template

```go
// ============================================================================
// Section N: Section Name — Lesson Title
// Level: Beginner | Intermediate | Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Key concept 1
//   - Key concept 2
//
// ENGINEERING DEPTH:
//   Explanation of why this matters in production Go
//
// RUN: go run ./path/to/lesson
// ============================================================================
```

## Prerequisites and Dependencies

The curriculum uses a directed acyclic graph (DAG) of dependencies:

- **Entry points**: 15 lessons with no prerequisites
- **Chains**: linear progressions like `GS.1 → GS.2 → GS.3 → GS.4`
- **Fans**: lessons depending on multiple prerequisites
- **Exercises**: capstone-style lessons combining several concepts

## Reading `curriculum.json`

The master curriculum lives in [`/curriculum.json`](../../curriculum.json):

```json
{
  "sections": [
    {
      "id": "s01",
      "num": "§01",
      "title": "Core Foundations",
      "lessons": [
        {
          "id": "GS.1",
          "name": "installation",
          "concept": "Verify Go binary; inspect GOROOT, GOPATH",
          "prerequisites": [],
          "is_entry": true,
          "is_exercise": false,
          "path": "01-core-foundations/getting-started/1-installation"
        }
      ]
    }
  ]
}
```

Key fields:

- `is_entry`: true when the lesson has no prerequisites
- `is_exercise`: true when the lesson synthesizes multiple concepts
- `prerequisites`: array of lesson IDs required before the lesson
- `path`: directory relative to the repository root

## Next Steps

- **Start learning**: [01-core-foundations/getting-started](../../01-core-foundations/getting-started)
- **See dependencies**: [dependency-graphs.html](./dependency-graphs.html)
- **Contribute**: [CONTRIBUTING.md](../../CONTRIBUTING.md)

---

**Last Updated**: April 2026  
**Curriculum Version**: 1.0  
**Total Lines of Go Code**: 8000+
