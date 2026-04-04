# The Go Engineer Curriculum Map

Welcome to the complete learning path for **The Go Engineer**. This directory contains comprehensive curriculum documentation and visual references.

## 📚 Curriculum Overview

The Go Engineer teaches **15 major sections** across **80+ lessons** through a carefully designed **Language → Runtime → IO → Quality → Architecture** progression.

### Section Breakdown

#### **Core Foundation (Sections 01-03)**
Master the language fundamentals in 15 lessons.

- **§01 — Core Foundations** (8 lessons)
  - Getting Started: Go installation, hello-world, compilation model, dev tools (GS.1–4)
  - Language Basics: variables, constants, iota, application-logger exercise (LB.1–4)

- **§02 — Control Flow** (4 lessons)
  - for-loop, if/else, switch, pricing-calculator exercise (CF.1–4)

- **§03 — Data Structures** (6 lessons)
  - arrays, slices, maps, pointers, advanced-slicing, contact-manager exercise (DS.1–6)

#### **Functions & Errors (Section 04)**
Error handling design patterns in 9 lessons.

- **§04 — Functions & Errors** (9 lessons)
  - functions, closures, variadic functions, multiple returns
  - custom errors, error-wrapping, defer, panic/recover, error-handling exercise (FE.1–9)

#### **Types & Design (Sections 05-07)**
Object-oriented and functional patterns in 11 lessons.

- **§05 — Types & Interfaces** (6 lessons)
  - structs, methods, interfaces, Stringer, generics, payroll-processor exercise (TI.1–6)

- **§06 — Composition** (2 lessons)
  - composition vs inheritance, embedding (CO.1–2)

- **§07 — Strings & Text** (5 lessons)
  - string internals, formatting, unicode/runes, regex, text templates (ST.1–5)

#### **Modules & IO (Sections 08-09)**
Real-world input/output and module management in 18 lessons.

- **§08 — Modules & Packages** (3 lessons)
  - module basics, managing dependencies, versioning (MP.1–3)

- **§09 — IO & CLI** (15 lessons)
  - Filesystem (8): files, paths, directories, temp files, embed, io.Reader/Writer patterns, fs.FS testing (FS.1–8)
  - Encoding (5): JSON marshalling/unmarshalling, streaming, base64 (EN.1–5)
  - CLI Tools (3): args, flags, subcommands (CL.1–3)

#### **Web & Database (Section 10)**
Full-stack development in 19 lessons.

- **§10 — Web & Database** (19 lessons)
  - Databases (5): connecting, INSERT, SELECT, prepared statements, transactions (DB.1–5)
  - Web Masterclass (11): routing, DI, templates, middleware, sessions, auth, forms, CRUD, pagination, comments (WM.1–11)
  - HTTP Client (2): basic GET, refactoring for testability (HC.1–2)

> **NOTE**: This section is the most complex; learners **must** understand middleware and DI before advanced lessons.

#### **Concurrency (Sections 11-12)**
Parallel programming patterns in 21 lessons.

- **§11 — Concurrency** (17 lessons)
  - Goroutines & Channels (10): goroutines, WaitGroups, channels, buffered channels, closing, pipelines, races, select, sync primitives (GC.1–10)
  - Context (4): Background/TODO, WithCancel, WithTimeout, WithValue (CT.1–4)
  - Time & Scheduling (7): time basics, formatting, timers/tickers, random, scheduler, timezones (TM.1–7)

- **§12 — Concurrency Patterns** (3 lessons)
  - errgroup basics, errgroup + context, sync.Pool (CP.1–3)

#### **Quality & Performance (Section 13)**
Testing and optimization in 14 lessons.

- **§13 — Quality & Performance** (14 lessons)
  - Testing (4): unit tests, table-driven, HTTP handler tests, benchmarking (TE.1–4)
  - Mocking (4): manual mocks, function-injection, table-driven, testify/mock (HM.1–4)
  - Profiling (2): CPU profiling, live pprof HTTP endpoint (PR.1–2)

#### **Architecture (Section 14)**
Production-grade design patterns in 15 lessons.

- **§14 — Application Architecture** (15 lessons)
  - Package Design (3): naming, visibility, project layout (PD.1–3)
  - Docker (3): single-stage Dockerfile, multi-stage builds, layer caching (DO.1–3)
  - Logging (4): slog basics, context-keyed logger, custom handlers, zerolog comparison (SL.1–4)
  - gRPC (3): proto definition, unary server, unary client (GR.1–3)
  - Graceful Shutdown (2): signal.NotifyContext, HTTP graceful drain (GS.1–2)
  - Enterprise Capstone (1): Full REST API with PostgreSQL & Docker Compose (EC.1)

#### **Code Generation (Section 15)**
Build-time automation in 1 lesson.

- **§15 — Code Generation** (1 lesson)
  - go:generate directive and canonical tools (mockery, stringer, sqlc) (CG.1)

---

## 📊 Curriculum Statistics

| Metric | Count |
|--------|-------|
| Total Sections | 15 |
| Total Lessons | 80+ |
| Entry Points | 15 (one per section) |
| Exercises | 25+ |
| Prerequisite Links | 200+ |
| Total Concepts | 250+ |

---

## 🗺️ Visual Dependency Graph

See `dependency-graphs.html` in this directory for an interactive visual map of all lesson dependencies.

**How to use**:
1. Open `dependency-graphs.html` in a browser
2. Click on any lesson to see its prerequisites and dependents
3. Follow the arrows to understand the learning path

---

## 📖 How to Navigate

### If you're a **Complete Beginner**
Start at **§01 Getting Started** (GS.1) and follow sequentially. The curriculum is designed so that each lesson builds on previous concepts.

### If you know another language
Skim §01-03 (basics), deep-dive into §04-07 (Go idioms), then §09+ (real applications).

### If you're an experienced Go developer
Jump directly to:
- §11-12 for advanced concurrency
- §13 for testing patterns
- §14 for production architecture

### If you're learning a specific topic
Use the dependency graph to find all related lessons:
```
Topics: Database, Web, Testing, Concurrency, Logging, Docker, etc.
```

---

## ✅ Curriculum Validation

Run the curriculum validation tool to ensure all lessons are correctly mapped to the filesystem:

```bash
go run ./scripts/validate_curriculum.go
```

This verifies:
- All lessons in `curriculum.json` have corresponding directories
- No orphaned directories lack curriculum entries
- All paths are consistent

---

## 📝 Lesson Structure

Each lesson directory follows a consistent structure:

```
NN-lesson-name/
├── main.go              (Primary lesson code with full explanation)
├── README.md            (Optional: Additional context for complex lessons)
└── _starter/            (Optional: Starting point for follow-along)
```

### main.go Header Template
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

---

## 🎯 Prerequisites & Dependencies

The curriculum uses a **directed acyclic graph (DAG)** of dependencies:

- **Entry Points**: 15 lessons with no prerequisites (one per section)
- **Chains**: Linear progressions (e.g., GS.1 → GS.2 → GS.3 → GS.4)
- **Fans**: Lessons depending on multiple prerequisites (e.g., WM.8 requires WM.1-7 + DB.1-3)
- **Exercises**: Capstone lessons combining 5+ prerequisites

---

## 📚 Reading the Curriculum JSON

The master curriculum is defined in `/curriculum.json`:

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
        // ... more lessons
      ]
    }
  ]
}
```

**Key Fields**:
- `is_entry`: true if this lesson has no prerequisites
- `is_exercise`: true if this is a capstone combining multiple concepts
- `prerequisites`: Array of lesson IDs required before this one
- `path`: Directory relative to repository root

---

## 🚀 Next Steps

- **Start Learning**: Begin at [01-core-foundations/getting-started](../../01-core-foundations/getting-started)
- **See Dependencies**: Open `dependency-graphs.html`
- **Contribute**: Read [CONTRIBUTING.md](../../CONTRIBUTING.md) for adding lessons

---

**Last Updated**: April 2026  
**Curriculum Version**: 1.0  
**Total Lines of Go Code**: 8000+
