# The Go Engineer Learning Path

> This guide explains how to move through the v2.1 curriculum.
> The source of truth for all curriculum content is [ARCHITECTURE.md](./ARCHITECTURE.md).

---

## The 5 Phases

| Phase | Name                  | Sections | Progress  |
| ----- | --------------------- | -------- | --------- |
| 0     | Machine Foundation    | s00      | 0% → 5%   |
| 1     | Language Foundation   | s01–s04  | 5% → 52%  |
| 2     | Engineering Core      | s05–s08  | 52% → 87% |
| 3     | Systems Engineering   | s09–s10  | 87% → 96% |
| 4     | Flagship Project      | s11      | 96% → 100%|

---

## Three Ways to Move

### Full Path

Best for:

- complete beginners
- learners new to programming
- learners new to Go who want the most support

Rule:

- follow the 12 sections in order (s00 → s11)
- complete the required milestone or proof surface in each section
- do not skip repetition by default

### Bridge Path

Best for:

- experienced programmers who are new to Go
- learners who know programming already but need Go-specific instincts

Rule:

- keep the same section order
- skim setup repetition where appropriate
- do not skip proof surfaces just because the syntax feels familiar

Suggested route:

1. `00-how-computers-work` — skim if you already understand processes, memory, and the terminal
2. `01-getting-started` as a short sanity pass
3. Skim Phase 1 (s02–s04) for Go-specific patterns
4. Slow down at Phase 2 (s05–s08) where Go's error handling, concurrency, and testing discipline become critical
5. Continue through Phase 3–4 in order

### Targeted Path

Best for:

- working Go developers
- learners strengthening one weak area
- learners returning to improve a specific skill without replaying everything

Rule:

- choose one phase or section intentionally
- check its prerequisites honestly
- complete that section's proof surface before claiming mastery

---

## Recommended Entry Points

| Goal                                          | Start Here                       | Before You Start                                                 |
| --------------------------------------------- | -------------------------------- | ---------------------------------------------------------------- |
| First Go setup and execution                  | s01: Getting Started             | None                                                             |
| Understand how computers work                 | s00: How Computers Work          | None                                                             |
| Strengthen fundamentals and type design       | s04: Types & Design              | Complete or skim s01–s03                                         |
| Improve testing, profiling, and concurrency   | s07: Concurrency / s08: Quality  | Be solid on functions, errors, and interfaces (s03–s04)          |
| Build backend APIs, gRPC, and databases       | s06: Backend, APIs & Databases   | Be comfortable with I/O (s05) and error handling (s03)           |
| Architecture and security engineering         | s09: Architecture & Security     | Complete Phase 2 (s05–s08)                                       |
| Deploy and operate Go services                | s10: Production Operations       | Complete s06 (HTTP) and s07 (concurrency)                        |
| Build a portfolio-level integrated system     | s11: GoScale Flagship            | Complete all capstone exercises in s00–s10                        |

---

## Validation Floors

Fast paths are allowed. Proof is not optional.

- **Full Path**: complete every required milestone and section proof surface
- **Bridge Path**: complete the important section proof surfaces even if you skim repetition
- **Targeted Path**: complete the chosen section's proof surface before claiming mastery

---

## Companion Docs

| Document                                              | Purpose                                |
| ----------------------------------------------------- | -------------------------------------- |
| [README.md](./README.md)                              | Project overview and quick start       |
| [ARCHITECTURE.md](./ARCHITECTURE.md)                  | Curriculum source of truth (v2.1)      |
| [CURRICULUM-BLUEPRINT.md](./CURRICULUM-BLUEPRINT.md)  | Teaching contract and lesson standards  |
| [CODE-STANDARDS.md](./CODE-STANDARDS.md)              | Code style and engineering standards   |
| [COMMON-MISTAKES.md](./COMMON-MISTAKES.md)            | 15 common Go bugs and fixes            |
| [ROADMAP.md](./ROADMAP.md)                            | What is built and what is planned      |

---

## Bottom Line

Choose the phase that matches your real background, read the section's `README.md` first, then run the linked source code.
