# The Go Engineer v3 — Curriculum

> Schema: 3.1.0 | Curriculum: 3.0.0 | Status: stable | Updated: 2026-10-23

## Architecture Decision

**Decision:** Authoritative v3 modular curriculum engineered around a zero-magic phase-based path from foundational Go to production software engineering.

- **Authoritative Source of Truth:** Authoritative curriculum metadata resides in the modular JSON manifests: `curriculum/path.core.json`, `curriculum/path.electives.json`, `curriculum/projects.json`, `curriculum/assessments.json`, `curriculum/crossrefs.json`, `curriculum/concepts.json`, and `curriculum/failures.json`.
- **Physical Curriculum Delivery:** Physical lesson READMEs, runnable examples, test suites, and exercises are organized canonically under `curriculum/modules/` and `curriculum/electives/`.
- **Core Scope:** The core curriculum is dedicated to production Go backend engineering (systems, networking, persistence, concurrency, observability, architecture, and deployment). Modern server-driven full-stack frontend engineering is delivered via the companion path (semantic HTML, CSS, JavaScript, and progressive enhancement with Go templates/templ/HTMX/Tailwind/Datastar).
- **Capstone & Learning Path Flow:** Flagship Opslane (Module 18) directly requires the complete backend engineering and deployment runbook foundation of Module 15 (`core-15-19`). Module 16 (Portfolio and Interview Readiness) packages the production artifacts produced by Opslane for senior-level engineering interviews. Module 17 (Advanced Electives) provides independent specialization tracks in Generics, gRPC, and distributed systems.

## Quick Stats

| Metric                           | Count | Notes                                                             |
| -------------------------------- | ----- | ----------------------------------------------------------------- |
| Core Modules                     | 18    | Modules 00 through 16, plus Module 18 (Flagship Opslane)          |
| Elective Modules                 | 1     | Module 17 (Advanced Electives)                                    |
| Total Modules                    | 19    | Modules 00–18                                                     |
| Core Path Items (Metadata Total) | 421   | Total items array in `path.core.json` across all core modules     |
| Instructional Lessons & Labs     | 321   | 320 lessons + 1 lab across Modules 00–16                         |
| Flagship Opslane Milestones      | 15    | 14 checkpoints + 1 capstone defense in Module 18          |
| Learner Instructional Units      | 336   | 321 instructional lessons/labs + 15 Opslane milestones            |
| Elective Items                   | 20    | 20 advanced elective lessons across Generics, gRPC, and distributed patterns in Module 17 |
| Projects                         | 28    | 27 module projects + 1 flagship capstone defense project          |
| In-Path Assessment Checkpoints   | 57    | 16 module checkpoints + 15 Opslane assessments + 26 project-bound nodes |
| Assessment Registry Records      | 47    | 19 module assessments + 28 project assessments in `assessments.json`|
| Canonical Concepts               | 62    | Fully resolved ownership without leaks in `concepts.json`         |
| Operational Categories           | 13    | Mapped across 18 module coverage entries                          |

## Module Map

| #   | Module | Phase | Status | Lessons/Labs | Milestones | Projects | Assessments | Total Path Items | Path |
| --- | --- | --- | --- | ---: | ---: | ---: | ---: | ---: | --- |
| 00  | **Orientation** | orientation | stable | 10 | 0 | 1 | 2 | 13 | `curriculum/modules/00-orientation` |
| 01  | **Computers, Terminal, Git, and the Web** | foundations | stable | 15 | 0 | 1 | 1 | 17 | `curriculum/modules/01-computers-terminal-git-web` |
| 02  | **Go Setup and Tooling** | tooling | stable | 14 | 0 | 1 | 1 | 16 | `curriculum/modules/02-go-setup-tooling` |
| 03  | **Programming Fundamentals with Go** | go-core | stable | 24 | 0 | 2 | 3 | 29 | `curriculum/modules/03-programming-fundamentals` |
| 04  | **Functions, Errors, and Data Semantics** | go-core | stable | 18 | 0 | 1 | 2 | 21 | `curriculum/modules/04-functions-errors-data-semantics` |
| 05  | **Types, Interfaces, Packages, and Modules** | engineering-core | stable | 24 | 0 | 2 | 3 | 29 | `curriculum/modules/05-types-interfaces-packages-modules` |
| 06  | **Testing, Debugging, and Refactoring** | engineering-core | stable | 21 | 0 | 2 | 3 | 26 | `curriculum/modules/06-testing-debugging-refactoring` |
| 07  | **CLI, Files, JSON, Config, and Runtime Coordination** | cli-io | stable | 21 | 0 | 3 | 4 | 28 | `curriculum/modules/07-cli-files-json-config` |
| 08  | **HTTP and REST APIs** | backend | stable | 21 | 0 | 1 | 2 | 24 | `curriculum/modules/08-http-rest-apis` |
| 09  | **SQL and PostgreSQL Persistence** | data | stable | 28 | 0 | 1 | 2 | 31 | `curriculum/modules/09-sql-postgres-persistence` |
| 10  | **Authentication, Authorization, and Security** | security | stable | 25 | 0 | 1 | 2 | 28 | `curriculum/modules/10-auth-security` |
| 11  | **Lifecycle, Context, and Concurrency** | concurrency | stable | 29 | 0 | 3 | 4 | 36 | `curriculum/modules/11-lifecycle-context-concurrency` |
| 12  | **Observability and Diagnostics** | reliability | stable | 12 | 0 | 2 | 3 | 17 | `curriculum/modules/12-observability-diagnostics` |
| 13  | **Performance and Memory Engineering** | performance | stable | 8 | 0 | 2 | 3 | 13 | `curriculum/modules/13-performance-memory-engineering` |
| 14  | **Architecture and Distributed Systems** | architecture | stable | 17 | 0 | 1 | 2 | 20 | `curriculum/modules/14-architecture-distributed-systems` |
| 15  | **Docker, CI/CD, and Deployment** | delivery | stable | 19 | 0 | 1 | 2 | 22 | `curriculum/modules/15-docker-cicd-deployment` |
| 16  | **Portfolio and Interview Readiness** | career | stable | 15 | 0 | 1 | 2 | 18 | `curriculum/modules/16-portfolio-interview-readiness` |
| 18  | **Flagship Opslane** | flagship | stable | 0 | 15 | 2 | 16 | 33 | `curriculum/modules/18-flagship-opslane` |
| 17  | **Advanced Electives** | elective | stable | 20 | 0 | 0 | 0 | 20 | `curriculum/electives/advanced-electives` |
| —   | **Core Subtotal** | — | — | **321** | **15** | **28** | **57** | **421** | *(Modules 00–16, 18)* |
| —   | **Total Curriculum** | — | — | **341** | **15** | **28** | **57** | **441** | *(421 core items + 20 electives)* |

## Workload Estimates & 52-Week Pacing Guide

Curriculum completion time cannot be calculated from reading time alone. A realistic workload accounts for conceptual reading, hands-on lab exercises, independent debugging, test authoring, and capstone engineering.

### 1. Workload by Activity Layer (Core Path)

| Activity Layer | Items | On-Paper Metadata | Realistic Learner Effort | Notes |
| :--- | :---: | :---: | :---: | :--- |
| **Instructional Lessons & Labs** | 321 | 238.0 hrs | 235 – 265 hrs | Reading, mental models, running code examples, review questions (Modules 00–16) |
| **Mini-Projects** | 18 | 36.0 hrs | 36 – 72 hrs | Immediate practical reinforcement (2–4 hours per project) |
| **Portfolio Projects** | 8 | 16.0 hrs | 65 – 120 hrs | Multi-component services with tests, migrations, CI/CD (8–15 hours per project) |
| **Flagship Opslane Milestones & Projects** | 17 | 34.0 hrs | 55 – 85 hrs | Integrated production system with auth, DB, payments, worker pools, metrics (Module 18) |
| **In-Path Assessments & Defense** | 57 | 57.0 hrs | 45 – 60 hrs | Self-evaluation checklists, code reviews, and capstone oral/technical defense |
| **Core Total** | **421** | **381.0 hrs** | **435 – 600 hrs** | **Realistic full engineering workload across all 18 core modules** |

### 2. Weekly Commitment Planning Models

| Weekly Commitment | Total Weeks | Total Months | Viability & Pace Description |
| :---: | :---: | :---: | :--- |
| **12–15 hrs/week** | **40 – 52 weeks** | **10 – 12 mos** | **(Recommended)** Completes 435–600 hours of core backend engineering in 29–50 weeks of focused study, leaving deliberate slack for review, autonomous project debugging, and parallel companion frontend progression. |
| **8–10 hrs/week** | **55 – 65 weeks** | **13 – 15 mos** | Steady progress; core concepts finish in 1 year, with portfolio and Opslane defense extending into months 13–15. |
| **5 hrs/week** | **90 – 115 weeks** | **21 – 26 mos** | Casual part-time; learners reach Module 08/09 by month 12 and graduate in year 2. |

### 3. Integrated 52-Week Schedule Framework

| Weeks | Phase | Focus Modules | Key Practical Milestones |
| :---: | :--- | :--- | :--- |
| **1–3** | Foundations | Modules 00–02 | Environment verification, Git workflow, Shell + Git Lab, Tooling Failure Lab |
| **4–8** | Core Go Programming | Modules 03–04 | Syntax, memory models, functions, error handling, Pricing Checkout, Order Summary Refactor |
| **9–14** | Engineering Disciplines | Modules 05–07 | Types, interfaces, test-driven refactoring, CLI tools, Config Parser *(Companion: HTML/CSS)* |
| **15–24** | Backend & Data Persistence | Modules 08–09 | HTTP APIs, Postgres, migrations, transactions, Book Catalog REST API, Task API *(Companion: JS/DOM)* |
| **25–33** | Security, Lifecycle & Concurrency | Modules 10–11 | JWT/RBAC, goroutines, channels, worker pools, Secure Authenticated API *(Companion: templ/HTMX)* |
| **34–42** | Production Engineering | Modules 12–15 | Observability (slog/OTel), profiling (pprof), Docker multi-stage, CI/CD pipelines, deployment runbook |
| **43–50** | Flagship Implementation | Module 18 | Flagship Opslane: multi-tenant architecture, event bus, caching, graceful shutdown, production deployment |
| **51–52** | Portfolio & Graduation Review | Module 16 | Portfolio packaging, architecture defense rehearsal, resume project bullets, mock interview prep |

## Repository Structure

```
curriculum/
  modules/
    00-orientation/
    01-computers-terminal-git-web/
    02-go-setup-tooling/
    03-programming-fundamentals/
    04-functions-errors-data-semantics/
    05-types-interfaces-packages-modules/
    06-testing-debugging-refactoring/
    07-cli-files-json-config/
    08-http-rest-apis/
    09-sql-postgres-persistence/
    10-auth-security/
    11-lifecycle-context-concurrency/
    12-observability-diagnostics/
    13-performance-memory-engineering/
    14-architecture-distributed-systems/
    15-docker-cicd-deployment/
    16-portfolio-interview-readiness/
    18-flagship-opslane/
  electives/
    advanced-electives/
  companion-frontend/
    README.md
  assessments.json
  concepts.json
  crossrefs.json
  failures.json
  path.core.json
  path.electives.json
  projects.json
  readme.contracts.json
  schema.v3.json
  GRADUATION_CRITERIA.md
  PILOT_TESTING_PROTOCOL.md
internal/tools/curriculum/
```

## Companion Tracks & Quality Assurance Protocols

- **[Graduation Outcomes & Capstone Defense Rubric](GRADUATION_CRITERIA.md):** 6 observable competency outcomes (unassisted problem-solving, debugging, testing, disaster recovery, architectural defense, incident response) and formal Flagship Opslane oral defense protocol.
- **[Frontend Companion Track](companion-frontend/README.md):** 4-stage progressive server-driven hypermedia track (`templ` + `HTMX` + `Tailwind` + `Alpine.js`, with `Datastar` comparison) running parallel to backend modules and converging at Flagship Opslane.
- **[Beginner Pilot Testing Protocol](PILOT_TESTING_PROTOCOL.md):** Formal empirical zero-intervention testing methodology evaluating unassisted learner velocity and cognitive load across Modules 00–05.

---

# Core Path

## 0. Orientation

> **Phase:** orientation | **Status:** stable | **Required:** True

**Learning Goal:** Explain the learning system, repository workflow, zero-magic contract, assessments, and job-ready target.

| #   | Type    | Slug | Title | Minutes |
| --- | ------- | ---- | ----- | ------- |
| 1   | concept | how-to-use-this-repository | **How to use this repository** | 45 |
| 2   | concept | what-zero-magic-means | **What zero magic means** | 45 |
| 3   | concept | how-lessons-exercises-projects-and-checkpoints-work | **How lessons, exercises, projects, and checkpoints work** | 120 |
| 4   | concept | how-to-run-code-in-this-repository | **How to run code in this repository** | 45 |
| 5   | concept | how-starter-folders-work | **How starter folders work** | 45 |
| 6   | concept | how-assessments-work | **How assessments work** | 45 |
| 7   | concept | how-to-ask-good-debugging-questions | **How to ask good debugging questions** | 45 |
| 8   | concept | what-job-ready-means | **What job-ready means** | 45 |
| 9   | concept | how-to-use-the-roadmap | **How to use the roadmap** | 45 |
| 10   | concept | how-to-build-a-portfolio-from-this-curriculum | **How to build a portfolio from this curriculum** | 45 |

## 1. Computers, Terminal, Git, and the Web

> **Phase:** foundations | **Status:** stable | **Required:** True

**Learning Goal:** Make the machine, shell, Git, GitHub, and web basics feel safe before Go syntax.

**Prerequisites:** 0. Orientation

| #   | Type    | Slug | Title | Minutes |
| --- | ------- | ---- | ----- | ------- |
| 1   | concept | what-is-a-program | **What is a program?** | 45 |
| 2   | concept | source-code-executable-and-process | **Source code, executable, and process** | 45 |
| 3   | concept | files-bytes-directories-and-paths | **Files, bytes, directories, and paths** | 45 |
| 4   | concept | terminal-basics | **Terminal basics** | 45 |
| 5   | concept | environment-variables | **Environment variables** | 45 |
| 6   | concept | exit-codes | **Exit codes** | 45 |
| 7   | concept | how-the-os-manages-processes | **How the OS manages processes** | 45 |
| 8   | concept | memory-preview-stack-vs-heap | **Memory preview: stack vs heap** | 45 |
| 9   | concept | git-mental-model | **Git mental model** | 45 |
| 10   | concept | git-basics-status-add-commit | **Git basics: status, add, commit** | 45 |
| 11   | concept | branching-and-merging | **Branching and merging** | 45 |
| 12   | concept | github-workflow | **GitHub workflow** | 45 |
| 13   | concept | pull-requests-and-code-review | **Pull requests and code review** | 45 |
| 14   | concept | web-preview-client-server-dns-and-ports | **Web preview: client, server, DNS, and ports** | 45 |
| 15   | concept | http-request-and-response-preview | **HTTP request and response preview** | 45 |

## 2. Go Setup and Tooling

> **Phase:** tooling | **Status:** stable | **Required:** True

**Learning Goal:** Make the Go toolchain normal from day one.

**Prerequisites:** 1. Computers, Terminal, Git, and the Web

| #   | Type    | Slug | Title | Minutes |
| --- | ------- | ---- | ----- | ------- |
| 1   | concept | install-and-verify-go | **Install and verify Go** | 45 |
| 2   | concept | hello-world | **Hello World** | 30 |
| 3   | concept | go-run | **go run** | 30 |
| 4   | concept | go-build | **go build** | 45 |
| 5   | testing | go-test | **go test** | 45 |
| 6   | concept | gofmt | **gofmt** | 30 |
| 7   | concept | go-vet | **go vet** | 30 |
| 8   | concept | go-doc | **go doc** | 30 |
| 9   | concept | editor-setup-and-gopls | **Editor setup and gopls** | 45 |
| 10   | concept | reading-compiler-errors | **Reading compiler errors** | 45 |
| 11   | concept | reading-runtime-errors | **Reading runtime errors** | 45 |
| 12   | testing | reading-test-failures | **Reading test failures** | 45 |
| 13   | concept | go-module-root-basics | **Go module root basics** | 45 |
| 14   | concept | tooling-checklist | **Tooling checklist** | 30 |

## 3. Programming Fundamentals with Go

> **Phase:** go-core | **Status:** stable | **Required:** True

**Learning Goal:** Teach core programming through Go without hidden assumptions.

**Prerequisites:** 2. Go Setup and Tooling

| #   | Type    | Slug | Title | Minutes |
| --- | ------- | ---- | ----- | ------- |
| 1   | concept | values-and-expressions | **Values and expressions** | 45 |
| 2   | concept | variables | **Variables** | 45 |
| 3   | concept | types | **Types** | 45 |
| 4   | concept | zero-values | **Zero values** | 45 |
| 5   | concept | type-conversions | **Type conversions** | 45 |
| 6   | concept | constants | **Constants** | 45 |
| 7   | concept | iota | **iota** | 45 |
| 8   | concept | boolean-logic | **Boolean logic** | 45 |
| 9   | concept | if-and-else | **If and else** | 45 |
| 10   | concept | switch | **Switch** | 45 |
| 11   | concept | loops | **Loops** | 45 |
| 12   | concept | range | **Range** | 45 |
| 13   | concept | arrays | **Arrays** | 45 |
| 14   | concept | slices | **Slices** | 45 |
| 15   | concept | slice-length-and-capacity | **Slice length and capacity** | 45 |
| 16   | concept | slice-sharing-and-aliasing | **Slice sharing and aliasing** | 45 |
| 17   | concept | maps | **Maps** | 45 |
| 18   | concept | comma-ok-idiom | **Comma-ok idiom** | 45 |
| 19   | concept | strings | **Strings** | 45 |
| 20   | concept | bytes | **Bytes** | 45 |
| 21   | concept | runes-and-unicode | **Runes and Unicode** | 45 |
| 22   | concept | pointers-as-addresses | **Pointers as addresses** | 45 |
| 23   | concept | beginner-pointer-mistakes | **Beginner pointer mistakes** | 45 |
| 24   | concept | algorithmic-problem-solving | **Algorithmic problem-solving and complexity** | 45 |

## 4. Functions, Errors, and Data Semantics

> **Phase:** go-core | **Status:** stable | **Required:** True

**Learning Goal:** Teach boundaries, failure, cleanup, mutation, and data behavior.

**Prerequisites:** 3. Programming Fundamentals with Go

| #   | Type    | Slug | Title | Minutes |
| --- | ------- | ---- | ----- | ------- |
| 1   | concept | function-basics | **Function basics** | 45 |
| 2   | concept | parameters | **Parameters** | 45 |
| 3   | concept | return-values | **Return values** | 45 |
| 4   | concept | multiple-return-values | **Multiple return values** | 45 |
| 5   | concept | scope | **Scope** | 45 |
| 6   | concept | call-stack-mental-model | **Call stack mental model** | 45 |
| 7   | concept | passing-by-value | **Passing by value** | 45 |
| 8   | concept | pointer-and-value-mutation-behavior | **Pointer and value mutation behavior** | 45 |
| 9   | concept | errors-as-values | **Errors as values** | 45 |
| 10   | concept | error-wrapping | **Error wrapping** | 45 |
| 11   | concept | errors-is | **errors.Is** | 45 |
| 12   | concept | errors-as | **errors.As** | 45 |
| 13   | concept | validation | **Validation** | 45 |
| 14   | concept | orchestration | **Orchestration** | 45 |
| 15   | concept | defer-mechanics | **defer mechanics** *(Workshop: Deterministic Cleanup — Part 1)* | 45 |
| 16   | concept | defer-for-cleanup | **defer for cleanup** *(Workshop: Deterministic Cleanup — Part 2)* | 45 |
| 17   | concept | panic-and-recover | **panic and recover** | 45 |
| 18   | concept | custom-error-types | **Custom error types** | 45 |

## 5. Types, Interfaces, Packages, and Modules

> **Phase:** engineering-core | **Status:** stable | **Required:** True

**Learning Goal:** Turn scripts into maintainable Go programs.

**Prerequisites:** 4. Functions, Errors, and Data Semantics

| #   | Type    | Slug | Title | Minutes |
| --- | ------- | ---- | ----- | ------- |
| 1   | concept | structs | **Structs** | 45 |
| 2   | concept | methods | **Methods** | 45 |
| 3   | concept | pointer-receivers | **Pointer receivers** *(Workshop: Method Receivers & Mutation — Part 1)* | 45 |
| 4   | concept | value-receivers | **Value receivers** *(Workshop: Method Receivers & Mutation — Part 2)* | 45 |
| 5   | concept | receiver-sets | **Receiver sets** *(Workshop: Method Receivers & Mutation — Part 3)* | 45 |
| 6   | concept | composition | **Composition** | 45 |
| 7   | concept | embedding | **Embedding** | 45 |
| 8   | concept | interfaces | **Interfaces** | 45 |
| 9   | concept | small-interfaces | **Small interfaces** | 45 |
| 10   | concept | interface-embedding | **Interface embedding** | 45 |
| 11   | concept | stringer | **Stringer** | 45 |
| 12   | concept | type-assertions | **Type assertions** *(Workshop: Dynamic Type Inspection — Part 1)* | 45 |
| 13   | concept | type-switches | **Type switches** *(Workshop: Dynamic Type Inspection — Part 2)* | 45 |
| 14   | concept | nil-interfaces | **Nil interfaces** | 45 |
| 15   | concept | package-names | **Package names** | 45 |
| 16   | concept | export-rules | **Export rules** | 45 |
| 17   | concept | project-layout | **Project layout** | 120 |
| 18   | concept | modules | **Modules** | 45 |
| 19   | concept | go-mod | **go.mod** | 45 |
| 20   | concept | go-sum | **go.sum** | 45 |
| 21   | concept | dependency-management | **Dependency management** | 45 |
| 22   | concept | workspaces | **Workspaces** | 45 |
| 23   | concept | documentation-comments | **Documentation comments** | 45 |
| 24   | concept | generics-fundamentals | **Generics fundamentals** | 45 |

## 6. Testing, Debugging, and Refactoring

> **Phase:** engineering-core | **Status:** stable | **Required:** True

**Learning Goal:** Move proof, diagnosis, and safe change before serious backend work.

**Prerequisites:** 5. Types, Interfaces, Packages, and Modules

| #   | Type    | Slug | Title | Minutes |
| --- | ------- | ---- | ----- | ------- |
| 1   | testing | why-testing-exists | **Why testing exists** | 45 |
| 2   | testing | unit-testing | **Unit testing** | 45 |
| 3   | testing | table-driven-tests | **Table-driven tests** | 45 |
| 4   | testing | subtests | **Subtests** | 45 |
| 5   | testing | t-cleanup | **t.Cleanup** | 45 |
| 6   | testing | test-fixtures | **Test fixtures** | 45 |
| 7   | testing | golden-files | **Golden files** | 45 |
| 8   | testing | testable-design-with-io-writer | **Testable design with io.Writer** | 45 |
| 9   | testing | fakes-before-mocks | **Fakes before mocks** | 45 |
| 10   | testing | mocking-tradeoffs | **Mocking tradeoffs** | 45 |
| 11   | debugging | debugging-mindset | **Debugging mindset** | 45 |
| 12   | testing | reading-stack-traces | **Reading stack traces** | 45 |
| 13   | debugging | debugging-panics | **Debugging panics** | 45 |
| 14   | debugging | delve-basics | **Delve basics** *(Workshop: Delve Debugging Workflow — Step 1)* | 45 |
| 15   | debugging | breakpoints | **Breakpoints** *(Workshop: Delve Debugging Workflow — Step 2)* | 45 |
| 16   | testing | stepping-through-code | **Stepping through code** *(Workshop: Delve Debugging Workflow — Step 3)* | 45 |
| 17   | testing | inspecting-variables | **Inspecting variables** *(Workshop: Delve Debugging Workflow — Step 4)* | 45 |
| 18   | debugging | debugging-tests | **Debugging tests** | 45 |
| 19   | testing | refactoring-safely | **Refactoring safely** | 45 |
| 20   | testing | fuzz-testing | **Fuzz testing** | 45 |
| 21   | testing | race-detector-preview | **Race detector preview** | 45 |

## 7. CLI, Files, JSON, Config, and Runtime Coordination

> **Phase:** cli-io | **Status:** stable | **Required:** True

**Learning Goal:** Build useful local programs and understand process lifecycle before servers.

**Prerequisites:** 6. Testing, Debugging, and Refactoring

| #   | Type    | Slug | Title | Minutes |
| --- | ------- | ---- | ----- | ------- |
| 1   | concept | cli-args | **CLI args** | 45 |
| 2   | concept | flags | **Flags** | 45 |
| 3   | concept | subcommands | **Subcommands** | 45 |
| 4   | concept | standard-input-and-standard-output | **Standard input and standard output** | 45 |
| 5   | concept | io-reader | **io.Reader** | 45 |
| 6   | concept | io-writer | **io.Writer** | 45 |
| 7   | concept | files | **Files** | 45 |
| 8   | concept | paths | **Paths** | 45 |
| 9   | concept | directories | **Directories** | 45 |
| 10   | concept | temp-files | **Temp files** | 45 |
| 11   | concept | fs-fs | **fs.FS** | 45 |
| 12   | concept | json-marshal | **JSON marshal** *(Workshop: In-Memory JSON — Part 1)* | 45 |
| 13   | concept | json-unmarshal | **JSON unmarshal** *(Workshop: In-Memory JSON — Part 2)* | 45 |
| 14   | concept | json-encoder | **JSON encoder** *(Workshop: Streaming JSON — Part 1)* | 45 |
| 15   | concept | json-decoder | **JSON decoder** *(Workshop: Streaming JSON — Part 2)* | 45 |
| 16   | concept | config-files | **Config files** | 45 |
| 17   | concept | environment-variables-for-configuration | **Environment variables for configuration** | 45 |
| 18   | concept | config-validation | **Config validation** | 45 |
| 19   | concept | text-templates | **Text templates** | 45 |
| 20   | concept | base64-as-encoding-not-encryption | **Base64 as encoding, not encryption** | 45 |
| 21   | testing | cli-testability | **CLI testability** | 45 |

## 8. HTTP and REST APIs

> **Phase:** backend | **Status:** stable | **Required:** True

**Learning Goal:** Build the first serious backend service.

**Prerequisites:** 7. CLI, Files, JSON, Config, and Runtime Coordination

| #   | Type    | Slug | Title | Minutes |
| --- | ------- | ---- | ----- | ------- |
| 1   | concept | http-as-a-protocol | **HTTP as a protocol** | 45 |
| 2   | concept | net-http-basics | **net/http basics** | 45 |
| 3   | concept | handler-lifecycle | **Handler lifecycle** | 45 |
| 4   | concept | routing | **Routing** | 45 |
| 5   | concept | request-parsing | **Request parsing** | 45 |
| 6   | concept | input-validation | **Input validation** | 45 |
| 7   | concept | response-writing | **Response writing** | 45 |
| 8   | concept | api-error-shape | **API error shape** | 45 |
| 9   | concept | status-codes | **Status codes** | 45 |
| 10   | concept | middleware-pattern | **Middleware pattern** | 45 |
| 11   | concept | error-middleware | **Error middleware** | 45 |
| 12   | testing | handler-testing-with-httptest | **Handler testing with httptest** | 45 |
| 13   | concept | body-size-limits | **Body size limits** | 45 |
| 14   | concept | server-timeouts | **Server timeouts** | 45 |
| 15   | concept | graceful-shutdown | **Graceful shutdown** | 45 |
| 16   | concept | health-and-readiness-probes | **Health and readiness probes** | 45 |
| 17   | concept | rest-design-principles | **REST design principles** | 45 |
| 18   | concept | api-versioning | **API versioning** | 45 |
| 19   | concept | pagination-and-filtering | **Pagination and filtering** | 45 |
| 20   | concept | openapi-basics | **OpenAPI basics** | 45 |
| 21   | concept | http-client | **HTTP Client** | 60 |

## 9. SQL and PostgreSQL Persistence

> **Phase:** data | **Status:** stable | **Required:** True

**Learning Goal:** Teach real relational backend work.

**Prerequisites:** 8. HTTP and REST APIs

| #   | Type    | Slug | Title | Minutes |
| --- | ------- | ---- | ----- | ------- |
| 1   | concept | why-databases-exist | **Why databases exist** | 45 |
| 2   | concept | relational-modeling | **Relational modeling** | 45 |
| 3   | concept | tables-rows-and-columns | **Tables, rows, and columns** | 45 |
| 4   | concept | primary-keys | **Primary keys** | 45 |
| 5   | concept | foreign-keys | **Foreign keys** | 45 |
| 6   | concept | constraints | **Constraints** | 45 |
| 7   | concept | sqlite-as-a-local-learning-tool | **SQLite as a local learning tool** | 45 |
| 8   | concept | postgresql-as-production-default | **PostgreSQL as production default** | 45 |
| 9   | concept | docker-compose-for-postgres | **Docker Compose for Postgres** | 45 |
| 10   | concept | database-sql | **database/sql** | 45 |
| 11   | concept | sql-db-as-a-pool | **sql.DB as a pool** | 45 |
| 12   | concept | insert | **INSERT** | 45 |
| 13   | concept | select | **SELECT** | 45 |
| 14   | concept | scanning-rows | **Scanning rows** | 45 |
| 15   | concept | null-handling | **Null handling** | 45 |
| 16   | concept | prepared-statements | **Prepared statements** | 45 |
| 17   | concept | transactions | **Transactions** | 45 |
| 18   | concept | rollback-discipline | **Rollback discipline** | 45 |
| 19   | concept | migrations | **Migrations** | 45 |
| 20   | concept | joins | **Joins** | 45 |
| 21   | concept | indexes | **Indexes** | 45 |
| 22   | concept | explain | **EXPLAIN** | 45 |
| 23   | concept | query-timeouts-with-context | **Query timeouts with context** | 45 |
| 24   | concept | repository-and-service-seam | **Repository and service seam** | 45 |
| 25   | testing | integration-tests | **Integration tests** | 45 |
| 26   | concept | sqlc | **sqlc** | 45 |
| 27   | concept | update-operations | **UPDATE operations** | 45 |
| 28   | concept | delete-operations | **DELETE operations** | 45 |

## 10. Authentication, Authorization, and Security

> **Phase:** security | **Status:** stable | **Required:** True

**Learning Goal:** Make secure backend behavior explicit.

**Prerequisites:** 9. SQL and PostgreSQL Persistence

| #   | Type    | Slug | Title | Minutes |
| --- | ------- | ---- | ----- | ------- |
| 1   | security | threat-modeling | **Threat modeling** | 45 |
| 2   | security | trust-boundaries | **Trust boundaries** | 45 |
| 3   | security | input-validation-for-security | **Input validation for security** | 45 |
| 4   | security | authentication-vs-authorization | **Authentication vs authorization** | 45 |
| 5   | security | password-hashing | **Password hashing** | 45 |
| 6   | security | sessions-and-cookies | **Sessions and cookies** | 45 |
| 7   | security | jwt-implementation | **JWT implementation** | 45 |
| 8   | security | jwt-risks | **JWT risks** | 45 |
| 9   | security | oauth2-and-oidc-overview | **OAuth2 and OIDC overview** | 45 |
| 10   | security | api-keys | **API keys** | 45 |
| 11   | security | rbac | **RBAC** | 45 |
| 12   | security | abac-and-policy-checks | **ABAC and policy checks** | 45 |
| 13   | security | tenant-isolation-basics | **Tenant isolation basics** | 45 |
| 14   | security | sql-injection-prevention | **SQL injection prevention** | 45 |
| 15   | security | xss | **XSS** | 45 |
| 16   | security | csrf | **CSRF** | 45 |
| 17   | security | cors | **CORS** | 45 |
| 18   | security | rate-limiting | **Rate limiting** | 45 |
| 19   | security | tls-and-https | **TLS and HTTPS** | 45 |
| 20   | security | secrets-management | **Secrets management** | 45 |
| 21   | security | dependency-security | **Dependency security** | 45 |
| 22   | security | govulncheck | **govulncheck** | 45 |
| 23   | security | owasp-for-go-apis | **OWASP for Go APIs** | 45 |
| 24   | security | secure-logging | **Secure logging** | 45 |
| 25   | concept | refresh-tokens | **Refresh tokens** | 45 |

## 11. Lifecycle, Context, and Concurrency

> **Phase:** concurrency | **Status:** stable | **Required:** True

**Learning Goal:** Teach safe concurrent systems and request lifecycle after learners understand services.

**Prerequisites:** 10. Authentication, Authorization, and Security

| #   | Type    | Slug | Title | Minutes |
| --- | ------- | ---- | ----- | ------- |
| 1   | concept | time-basics | **Time basics** *(Workshop: Time Mechanics & Scheduling — Part 1)* | 45 |
| 2   | concept | durations | **Durations** *(Workshop: Time Mechanics & Scheduling — Part 2)* | 45 |
| 3   | concept | timers | **Timers** *(Workshop: Time Mechanics & Scheduling — Part 3)* | 45 |
| 4   | concept | tickers | **Tickers** *(Workshop: Time Mechanics & Scheduling — Part 4)* | 45 |
| 5   | concept | context-basics | **Context basics** | 45 |
| 6   | concept | cancellation | **Cancellation** | 45 |
| 7   | concept | timeouts | **Timeouts** | 45 |
| 8   | concept | deadlines | **Deadlines** | 45 |
| 9   | concept | context-values-with-caveats | **Context values with caveats** | 45 |
| 10   | concept | why-concurrency-exists | **Why concurrency exists** | 45 |
| 11   | concept | goroutines | **Goroutines** | 45 |
| 12   | concept | waitgroups | **WaitGroups** | 45 |
| 13   | concept | channels | **Channels** | 45 |
| 14   | concept | buffered-channels | **Buffered channels** | 45 |
| 15   | concept | closing-channels | **Closing channels** | 45 |
| 16   | concept | channel-ownership | **Channel ownership** | 45 |
| 17   | concept | backpressure | **Backpressure** | 45 |
| 18   | concept | pipelines | **Pipelines** | 45 |
| 19   | concept | mutex-and-rwmutex | **Mutex and RWMutex** | 45 |
| 20   | concept | atomics-with-caution | **Atomics with caution** | 45 |
| 21   | concept | race-conditions | **Race conditions** | 45 |
| 22   | concept | race-detector | **Race detector** | 45 |
| 23   | concept | goroutine-leaks | **Goroutine leaks** | 45 |
| 24   | concept | deadlocks | **Deadlocks** | 45 |
| 25   | concept | errgroup | **errgroup** | 45 |
| 26   | concept | bounded-worker-pools | **Bounded worker pools** | 45 |
| 27   | concept | retries-and-backoff | **Retries and backoff** | 45 |
| 28   | concept | idempotency-keys | **Idempotency keys** | 45 |
| 29   | concept | sync-once | **sync.Once** | 30 |

## 12. Observability and Diagnostics

> **Phase:** reliability | **Status:** stable | **Required:** True

**Learning Goal:** Make services observable before optimizing them.

**Prerequisites:** 11. Lifecycle, Context, and Concurrency

| #   | Type    | Slug | Title | Minutes |
| --- | ------- | ---- | ----- | ------- |
| 1   | operations | why-observability-exists | **Why observability exists** | 45 |
| 2   | operations | structured-logging-with-slog | **Structured logging with slog** | 45 |
| 3   | operations | request-scoped-logging | **Request-scoped logging** | 45 |
| 4   | operations | correlation-ids | **Correlation IDs** | 45 |
| 5   | operations | pii-redaction | **PII redaction** | 45 |
| 6   | operations | metrics | **Metrics** | 45 |
| 7   | operations | prometheus | **Prometheus** | 45 |
| 8   | operations | tracing | **Tracing** | 45 |
| 9   | operations | opentelemetry | **OpenTelemetry** | 45 |
| 10   | operations | alerting-mindset | **Alerting mindset** | 45 |
| 11   | operations | incident-debugging | **Incident debugging** | 45 |
| 12   | operations | reliability-review | **Reliability review** | 45 |

## 13. Performance and Memory Engineering

> **Phase:** performance | **Status:** stable | **Required:** True

**Learning Goal:** Make services efficient through profiling, benchmark-driven optimization, and memory understanding.

**Prerequisites:** 12. Observability and Diagnostics

| #   | Type    | Slug | Title | Minutes |
| --- | ------- | ---- | ----- | ------- |
| 1   | operations | pprof | **pprof** | 45 |
| 2   | operations | cpu-profiling | **CPU profiling** | 45 |
| 3   | operations | memory-profiling | **Memory profiling** | 45 |
| 4   | operations | benchmarks | **Benchmarks** | 45 |
| 5   | operations | escape-analysis | **Escape analysis** | 45 |
| 6   | operations | memory-layout | **Memory layout** | 45 |
| 7   | operations | caching-basics | **Caching basics** | 45 |
| 8   | operations | cache-invalidation | **Cache invalidation** | 45 |

## 14. Architecture and Distributed Systems

> **Phase:** architecture | **Status:** stable | **Required:** True

**Learning Goal:** Teach architecture as tradeoff reasoning with distributed systems fundamentals.

**Prerequisites:** 13. Performance and Memory Engineering

| #   | Type    | Slug | Title | Minutes |
| --- | ------- | ---- | ----- | ------- |
| 1   | architecture | why-architecture-exists | **Why architecture exists** | 45 |
| 2   | architecture | package-boundaries | **Package boundaries** | 45 |
| 3   | architecture | service-layer | **Service layer** | 45 |
| 4   | architecture | repository-pattern-deep-dive | **Repository pattern deep dive** | 45 |
| 5   | architecture | modular-monolith | **Modular monolith** | 45 |
| 6   | architecture | hexagonal-architecture | **Hexagonal architecture** | 45 |
| 7   | architecture | domain-modeling | **Domain modeling** | 45 |
| 8   | architecture | invariants | **Invariants** | 45 |
| 9   | architecture | event-driven-basics | **Event-driven basics** | 45 |
| 10   | architecture | queues | **Queues** | 45 |
| 11   | architecture | retries-and-idempotent-consumers | **Retries and idempotent consumers** | 45 |
| 12   | architecture | payment-workflow-design | **Payment workflow design** | 45 |
| 13   | architecture | multi-tenancy-architecture | **Multi-tenancy architecture** | 45 |
| 14   | architecture | caching-architecture | **Caching architecture** | 45 |
| 15   | architecture | when-to-split-services | **When to split services** | 45 |
| 16   | architecture | microservices-as-a-tradeoff-not-a-default | **Microservices as a tradeoff, not a default** | 45 |
| 17   | architecture | architecture-decision-records | **Architecture decision records** | 45 |

## 15. Docker, CI/CD, and Deployment

> **Phase:** delivery | **Status:** stable | **Required:** True

**Learning Goal:** Make software shippable and repeatable.

**Prerequisites:** 14. Architecture and Distributed Systems

| #   | Type    | Slug | Title | Minutes |
| --- | ------- | ---- | ----- | ------- |
| 1   | operations | linux-processes | **Linux processes** | 45 |
| 2   | operations | signals | **Signals** | 45 |
| 3   | operations | logs | **Logs** | 45 |
| 4   | operations | file-permissions | **File permissions** | 45 |
| 5   | operations | networking-basics | **Networking basics** | 45 |
| 6   | operations | docker-basics | **Docker basics** | 45 |
| 7   | operations | docker-images-and-layers | **Docker images and layers** | 45 |
| 8   | operations | multi-stage-builds | **Multi-stage builds** | 45 |
| 9   | operations | docker-compose | **Docker Compose** | 45 |
| 10   | operations | container-health-checks | **Container health checks** | 45 |
| 11   | operations | github-actions | **GitHub Actions** | 45 |
| 12   | operations | ci-pipeline | **CI pipeline** | 45 |
| 13   | operations | release-artifacts | **Release artifacts** | 45 |
| 14   | operations | config-in-deployment | **Config in deployment** | 45 |
| 15   | operations | secrets-in-deployment | **Secrets in deployment** | 45 |
| 16   | operations | vulnerability-scanning | **Vulnerability scanning** | 45 |
| 17   | operations | one-deployment-target | **One deployment target** | 45 |
| 18   | operations | rollback-basics | **Rollback basics** | 45 |
| 19   | operations | deployment-runbook | **Deployment runbook** | 45 |

## 16. Portfolio and Interview Readiness

> **Phase:** career | **Status:** stable | **Required:** True

**Learning Goal:** Convert competence into employability.

**Prerequisites:** 15. Docker, CI/CD, and Deployment

| #   | Type    | Slug | Title | Minutes |
| --- | ------- | ---- | ----- | ------- |
| 1   | career  | portfolio-strategy | **Portfolio strategy** | 30 |
| 2   | career  | project-selection-and-scoping | **Project selection and scoping** | 30 |
| 3   | career  | code-review-readiness | **Code review readiness** | 30 |
| 4   | career  | technical-writing-for-engineers | **Technical writing for engineers** | 30 |
| 5   | career  | resume-and-linkedin-optimization | **Resume and LinkedIn optimization** | 30 |
| 6   | career  | behavioral-interview-preparation | **Behavioral interview preparation** | 30 |
| 7   | career  | system-design-interview-basics | **System design interview basics** | 30 |
| 8   | career  | live-coding-interview-practice | **Live coding interview practice** | 30 |
| 9   | career  | take-home-project-strategy | **Take-home project strategy** | 30 |
| 10   | career  | negotiation-fundamentals | **Negotiation fundamentals** | 30 |
| 11   | career  | open-source-contribution | **Open source contribution** | 30 |
| 12   | career  | technical-blog-writing | **Technical blog writing** | 30 |
| 13   | career  | conference-talk-preparation | **Conference talk preparation** | 30 |
| 14   | career  | mentorship-and-code-review | **Mentorship and code review** | 30 |
| 15   | career  | career-roadmap | **Career roadmap** | 30 |

## 18. Flagship Opslane

> **Phase:** flagship | **Status:** stable | **Required:** True

**Learning Goal:** Final graduation proof through a portfolio-grade backend system.

**Prerequisites:** 15. Docker, CI/CD, and Deployment (Technical milestone: `core-15-19`)

> **Architecture Note:** Module 16 (Portfolio and Interview Readiness) is intentionally decoupled from technical implementation so learners can construct Flagship Opslane first, producing the production artifact that Module 16 packages for interviews.

| #   | Type    | Slug | Title | Minutes |
| --- | ------- | ---- | ----- | ------- |
| 1   | capstone | opslane-foundation-and-configuration | **Opslane foundation and configuration** | 120 |
| 2   | capstone | opslane-database-and-models | **Opslane database and models** | 120 |
| 3   | capstone | opslane-authentication | **Opslane authentication** | 120 |
| 4   | capstone | opslane-authorization | **Opslane authorization** | 120 |
| 5   | capstone | opslane-tenant-isolation | **Opslane tenant isolation** | 120 |
| 6   | capstone | opslane-http-api-layer | **Opslane HTTP API layer** | 120 |
| 7   | capstone | opslane-order-processing | **Opslane order processing** | 120 |
| 8   | capstone | opslane-payment-pipeline | **Opslane payment pipeline** | 120 |
| 9   | capstone | opslane-event-bus | **Opslane event bus** | 120 |
| 10   | capstone | opslane-worker-pools | **Opslane worker pools** | 120 |
| 11   | capstone | opslane-caching-layer | **Opslane caching layer** | 120 |
| 12   | capstone | opslane-observability | **Opslane observability** | 120 |
| 13   | capstone | opslane-graceful-shutdown | **Opslane graceful shutdown** | 120 |
| 14   | capstone | opslane-deployment | **Opslane deployment** | 120 |
| 15   | capstone | opslane-portfolio-defense | **Opslane portfolio defense** | 120 |

---

# Advanced Electives

## 17. Advanced Electives

> **Phase:** elective | **Status:** stable | **Required:** False

**Learning Goal:** Keep advanced Go and backend topics without blocking the core learner path.

**Prerequisites:** 16. Portfolio and Interview Readiness

| #   | Type      | Slug | Title | Minutes |
| --- | --------- | ---- | ----- | ------- |
| 1   | reference | generics | **Generics** | 45 |
| 2   | reference | complex-generic-constraints | **Complex generic constraints** | 45 |
| 3   | reference | generic-data-structures | **Generic data structures** | 45 |
| 4   | reference | functional-options | **Functional options** | 45 |
| 5   | reference | method-values | **Method values** | 45 |
| 6   | reference | build-tags | **Build tags** | 45 |
| 7   | reference | go-generate | **go generate** | 45 |
| 8   | reference | mockery | **Mockery** | 45 |
| 9   | reference | protobuf | **Protobuf** | 45 |
| 10   | reference | grpc-fundamentals | **gRPC fundamentals** | 45 |
| 11   | reference | grpc-streaming | **gRPC streaming** | 45 |
| 12   | reference | grpc-interceptors | **gRPC interceptors** | 45 |
| 13   | reference | rest-vs-grpc-tradeoffs | **REST vs gRPC tradeoffs** | 45 |
| 14   | reference | sync-map | **sync.Map** | 45 |
| 15   | reference | sync-pool | **sync.Pool** | 45 |
| 16   | reference | cqrs | **CQRS** | 45 |
| 17   | reference | feature-flags | **Feature flags** | 45 |
| 18   | reference | blue-green-deployments | **Blue/green deployments** | 45 |
| 19   | reference | microservice-splitting | **Microservice splitting** | 45 |
| 20   | reference | deeper-distributed-systems | **Deeper distributed systems** | 45 |

---

# Projects

| #   | Module | Title | Type | Stage | Difficulty |
| --- | ------ | ----- | ---- | ----- | ---------- |
| 1   | Mod 00 | **Course Setup Verification** | mini-project | foundations | foundation |
| 2   | Mod 01 | **Shell + Git Lab** | mini-project | foundations | foundation |
| 3   | Mod 02 | **Tooling Failure Lab** | mini-project | tooling | foundation |
| 4   | Mod 03 | **Contact Directory** | mini-project | go-core | core |
| 5   | Mod 03 | **Pricing Checkout** | mini-project | go-core | core |
| 6   | Mod 04 | **Order Summary Refactor** | mini-project | go-core | core |
| 7   | Mod 05 | **Bank Account** | mini-project | engineering-core | core |
| 8   | Mod 05 | **Payroll Processor** | mini-project | engineering-core | core |
| 9   | Mod 06 | **Bug Hunt Pack** | mini-project | engineering-core | core |
| 10   | Mod 06 | **Refactor Under Tests** | mini-project | engineering-core | core |
| 11   | Mod 07 | **Config Parser** | mini-project | cli-io | core |
| 12   | Mod 07 | **File Organizer** | mini-project | cli-io | core |
| 13   | Mod 07 | **Log Search** | mini-project | cli-io | core |
| 14   | Mod 08 | **Book Catalog REST API** | portfolio-project | backend | intermediate |
| 15   | Mod 09 | **Task API with Postgres** | portfolio-project | data | intermediate |
| 16   | Mod 10 | **Secure Authenticated API** | portfolio-project | security | intermediate |
| 17   | Mod 11 | **Bounded Worker Pool** | mini-project | concurrency | intermediate |
| 18   | Mod 11 | **Concurrent Downloader** | mini-project | concurrency | intermediate |
| 19   | Mod 11 | **URL Health Checker** | mini-project | concurrency | intermediate |
| 20   | Mod 12 | **Observability Retrofit** | portfolio-project | reliability | advanced |
| 21   | Mod 12 | **PII Redactor** | portfolio-project | reliability | intermediate |
| 22   | Mod 13 | **Performance Analysis Lab** | mini-project | core | advanced |
| 23   | Mod 13 | **Profile and Fix Lab** | portfolio-project | core | advanced |
| 24   | Mod 14 | **Modular Refactor Review** | mini-project | architecture | advanced |
| 25   | Mod 15 | **Dockerized Service with CI/CD** | portfolio-project | delivery | advanced |
| 26   | Mod 16 | **Portfolio Review Pack** | portfolio-project | career | advanced |
| 27   | Mod 17 | **Advanced Electives Capstone** | mini-project | elective | advanced |
| 28   | Mod 18 | **Opslane Flagship** | flagship | flagship | flagship |

---

# Assessments

The assessment framework evaluates mastery across end-of-module milestones and practical project implementations.

### Module Assessments (19 Records)

| ID | Title | Type | Passing Score |
| -- | ----- | ---- | ------------- |
| `assessment-module-00` | **Module 00 Assessment — Orientation** | project-review | 80% |
| `assessment-module-01` | **Module 01 Assessment — Computers, Terminal, Git, and the Web** | project-review | 80% |
| `assessment-module-02` | **Module 02 Assessment — Go Setup and Tooling** | project-review | 80% |
| `assessment-module-03` | **Module 03 Assessment — Programming Fundamentals with Go** | project-review | 80% |
| `assessment-module-04` | **Module 04 Assessment — Functions, Errors, and Data Semantics** | project-review | 80% |
| `assessment-module-05` | **Module 05 Assessment — Types, Interfaces, Packages, and Modules** | project-review | 80% |
| `assessment-module-06` | **Module 06 Assessment — Testing, Debugging, and Refactoring** | project-review | 80% |
| `assessment-module-07` | **Module 07 Assessment — CLI, Files, JSON, and Configuration** | project-review | 80% |
| `assessment-module-08` | **Module 08 Assessment — HTTP and REST APIs** | project-review | 80% |
| `assessment-module-09` | **Module 09 Assessment — SQL and PostgreSQL Persistence** | project-review | 80% |
| `assessment-module-10` | **Module 10 Assessment — Authentication, Authorization, and Security** | project-review | 80% |
| `assessment-module-11` | **Module 11 Assessment — Lifecycle, Context, and Concurrency** | project-review | 80% |
| `assessment-module-12` | **Module 12 Assessment — Observability and Diagnostics** | project-review | 80% |
| `assessment-module-13` | **Module 13 Assessment — Performance and Memory Engineering** | project-review | 80% |
| `assessment-module-14` | **Module 14 Assessment — Architecture and Distributed Systems** | project-review | 80% |
| `assessment-module-15` | **Module 15 Assessment — Docker, CI/CD, and Deployment** | project-review | 80% |
| `assessment-module-16` | **Module 16 Assessment — Portfolio and Interview Readiness** | portfolio-review | 80% |
| `assessment-module-17` | **Module 17 Assessment — Electives and Advanced Engineering** | project-review | 80% |
| `assessment-module-18` | **Module 18 Assessment — Flagship Opslane** | capstone-defense | 80% |

### Project Assessments (28 Records)

| ID | Title | Type | Passing Score |
| -- | ----- | ---- | ------------- |
| `assessment-project-shell-git-lab` | **Assessment — Shell + Git Lab** | project-review | 80% |
| `assessment-project-tooling-failure-lab` | **Assessment — Tooling Failure Lab** | project-review | 80% |
| `assessment-project-course-setup` | **Assessment — Course Setup Verification** | project-review | 80% |
| `assessment-project-performance-analysis-lab` | **Assessment — Performance Analysis Lab** | project-review | 80% |
| `assessment-project-advanced-electives` | **Assessment — Advanced Electives Capstone** | project-review | 80% |
| `assessment-project-pricing-checkout` | **Assessment — Pricing Checkout** | project-review | 80% |
| `assessment-project-contact-directory` | **Assessment — Contact Directory** | project-review | 80% |
| `assessment-project-order-summary-refactor` | **Assessment — Order Summary Refactor** | project-review | 80% |
| `assessment-project-bank-account` | **Assessment — Bank Account** | project-review | 80% |
| `assessment-project-payroll-processor` | **Assessment — Payroll Processor** | project-review | 80% |
| `assessment-project-bug-hunt-pack` | **Assessment — Bug Hunt Pack** | project-review | 80% |
| `assessment-project-refactor-under-tests` | **Assessment — Refactor Under Tests** | project-review | 80% |
| `assessment-project-file-organizer` | **Assessment — File Organizer** | project-review | 80% |
| `assessment-project-log-search` | **Assessment — Log Search** | project-review | 80% |
| `assessment-project-config-parser` | **Assessment — Config Parser** | project-review | 80% |
| `assessment-project-book-catalog-rest-api` | **Assessment — Book Catalog REST API** | project-review | 80% |
| `assessment-project-task-api-postgres` | **Assessment — Task API with Postgres** | project-review | 80% |
| `assessment-project-secure-authenticated-api` | **Assessment — Secure Authenticated API** | project-review | 80% |
| `assessment-project-concurrent-downloader` | **Assessment — Concurrent Downloader** | project-review | 80% |
| `assessment-project-url-health-checker` | **Assessment — URL Health Checker** | project-review | 80% |
| `assessment-project-bounded-worker-pool` | **Assessment — Bounded Worker Pool** | project-review | 80% |
| `assessment-project-pii-redactor` | **Assessment — PII Redactor** | project-review | 80% |
| `assessment-project-observability-retrofit` | **Assessment — Observability Retrofit** | project-review | 80% |
| `assessment-project-profile-and-fix-lab` | **Assessment — Profile and Fix Lab** | project-review | 80% |
| `assessment-project-modular-refactor-review` | **Assessment — Modular Refactor Review** | project-review | 80% |
| `assessment-project-dockerized-service-cicd` | **Assessment — Dockerized Service with CI/CD** | project-review | 80% |
| `assessment-project-portfolio-review-pack` | **Assessment — Portfolio Review Pack** | project-review | 80% |
| `assessment-project-opslane` | **Assessment — Opslane Flagship** | project-review | 80% |

---

# Cross-References

Total cross-references: **2089**

### Style Rules

| Relation      | Template                                                     |
| ------------- | ------------------------------------------------------------ |
| builds_on     | Builds on: Module X.Y — Lesson Title                         |
| preview_only  | Preview only: Explained fully in Module X.Y — Lesson Title   |
| reinforced_in | Reinforced in: Project X — Project Title                     |
| related       | Related: Module X.Y — Lesson Title                           |

### Relation Distribution

| Relation      | Count | Description |
| ------------- | ----- | ----------- |
| related       | 944   | Bidirectional conceptual links between related lessons |
| builds_on     | 587   | Strict directional pedagogical prerequisites |
| reinforced_in | 512   | Practical application links to projects and labs |
| preview_only  | 46    | Early mental-model previews developed fully later |

---

# Migration: v2 → v3

**Source Schema:** 2

**Target Schema:** 3.1.0

**Default Policy:** Every legacy item must be mapped to keep, rewrite, split, merge, move-to-elective, replace, or archive. Every migrated README must also receive a legacy_readme_action: reuse, rewrite, split, merge, or archive.

**Frozen Source:** `curriculum/legacy/curriculum.v2.json`

### High-Confidence Moves

| Legacy Pattern        | Outcome          | Target    | Reason                                                                    |
| --------------------- | ---------------- | --------- | ------------------------------------------------------------------------- |
| API.4-API.9           | move-to-elective | module-17 | Protobuf/gRPC is valuable but should not block the core beginner path.    |
| TI.14-TI.15           | move-to-elective | module-17 | Complex generics and generic data structures are advanced.                |
| TI.12                 | move-to-elective | module-17 | Functional options are useful but not core beginner material.             |
| CP.3                  | move-to-elective | module-17 | sync.Pool is performance-specialized.                                     |
| SY.2 sync.Map portion | move-to-elective | module-17 | sync.Map is specialized and should follow normal map+mutex reasoning.     |
| ARCH.7                | move-to-elective | module-17 | CQRS is advanced architecture.                                            |
| DEPLOY.2              | move-to-elective | module-17 | Blue/green deployment is advanced deployment strategy.                    |
| TE.1-TE.11            | rewrite          | module-06 | Testing and debugging must move before backend work.                      |
| DB.*                 | rewrite          | module-09 | Database path must become SQL/Postgres-first with migrations and EXPLAIN. |
| SEC.*                | rewrite          | module-10 | Security needs authn/authz strategy and secure API practice.              |

### Authoritative Quality Gates (17/17 Passing)

- `validate-schema`: 10 canonical JSON manifests valid against Schema 3.1.0
- `validate-graph`: Strict DAG acyclicity across all item dependencies
- `validate-lessons`: Zero machine role comments or unmotivated markers
- `validate-prerequisite-closure`: Strict pedagogical progression across module boundaries
- `validate-project-binding`: All 28 projects bound to valid curriculum placement anchors
- `validate-no-placeholder-zm`: Zero mechanical placeholders across zero-magic blocks
- `validate-assessment-completeness`: All 47 assessments equipped with weighted criteria
- `validate-unique-tags`: Tag taxonomy consistency across core and electives
- `validate-crossrefs`: Complete bidirectional graph integrity across all entities
- `validate-concept-ownership`: All 62 concepts mapped to unambiguous canonical owners
- `validate-cognitive-load`: Module pacing and cognitive density within verified bounds
- `validate-operational-failures`: 13 operational failure categories mapped across 18 modules
- `validate-failure-engineering`: Failure modes and detection strategies verified
- `validate-zero-magic`: 92 golden lessons certified against pedagogical spec
- `validate-physical-artifacts`: 441/441 physical lesson and elective READMEs verified on disk
- `validate-assessment-artifacts`: 47 assessments + 14 Opslane milestones verified with complete physical suites
- `validate-starter-parity`: 27 project starters (RED) and solutions (GREEN) verified with valid Go ASTs

### Validation Commands

| Command | Purpose | Execution |
| --- | --- | --- |
| **Run All Gates** | Execute full 17-gate curriculum test suite | `go run ./scripts/validate_curriculum.go` |
| **Validate Schema** | Check JSON manifests and key ordering | `go run ./internal/tools/curriculum validate-schema` |
| **Validate Graph** | Check DAG dependency structure | `go run ./internal/tools/curriculum validate-graph` |
| **Validate Physical** | Verify 441 physical READMEs | `go run ./internal/tools/curriculum validate-physical-artifacts` |
| **Validate Assessments** | Verify 4-file assessment suites | `go run ./internal/tools/curriculum validate-assessment-artifacts` |
| **Validate Code Parity** | Verify Go code and AST correctness | `go run ./internal/tools/curriculum validate-starter-parity` |

---

_Generated from curriculum v3.0.0 | Schema 3.1.0 | 2026-10-23_
