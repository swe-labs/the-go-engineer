# The Go Engineer Learning Path

> This guide explains how to move through the v2.1 curriculum.
> The source of truth for structure and section ownership is [ARCHITECTURE.md](./ARCHITECTURE.md).

## The 5 Phases

| Phase | Name | Sections | Progress |
| --- | --- | --- | --- |
| 0 | Machine Foundation | s00 | 0% -> 5% |
| 1 | Language Foundation | s01-s04 | 5% -> 52% |
| 2 | Engineering Core | s05-s08 | 52% -> 87% |
| 3 | Systems Engineering | s09-s10 | 87% -> 96% |
| 4 | Flagship Project | s11 | 96% -> 100% |

## Three Ways to Move

### Full Path

Best for:

- complete beginners
- learners new to programming
- learners new to Go who want the most support

Rule:

- follow the 12 sections in order (`s00` -> `s11`)
- complete each section's proof surface
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

1. Skim [`00-how-computers-work`](./00-how-computers-work) if you already understand processes, memory, and the terminal.
2. Use [`01-getting-started`](./01-getting-started) as a short sanity pass.
3. Move carefully through Phase 1 (`s02` to `s04`) to absorb Go-specific patterns.
4. Slow down at Phase 2 (`s05` to `s08`) where packages, I/O, concurrency, and testing habits matter more.
5. Continue through Phase 3 and Phase 4 in order.

### Targeted Path

Best for:

- working Go developers
- learners strengthening one weak area
- learners returning to improve a specific skill without replaying everything

Rule:

- choose one phase or section intentionally
- check its prerequisites honestly
- complete that section's proof surface before claiming mastery

## Recommended Entry Points

| Goal | Start Here | Before You Start |
| --- | --- | --- |
| First Go setup and execution | `s01` Getting Started | None |
| Understand how computers work | `s00` How Computers Work | None |
| Strengthen fundamentals and type design | `s04` Types & Design | Complete or skim `s01` to `s03` |
| Improve backend, API, and database fluency | `s06` Backend, APIs & Databases | Be comfortable with `s05` and explicit errors from `s03` |
| Improve concurrency and testing | `s07` Concurrency / `s08` Quality & Testing | Be solid on functions, errors, packages, and interfaces |
| Focus on architecture and security | `s09` Architecture & Security | Complete Phase 2 |
| Focus on deployment and operations | `s10` Production Operations | Complete `s06` and `s07` first |
| Build an integrated system | `s11` GoScale Flagship | Complete all earlier section milestones |

## Validation Floors

Fast paths are allowed. Proof is not optional.

- **Full Path**: complete every required milestone and section proof surface
- **Bridge Path**: complete the important proof surfaces even if you skim repetition
- **Targeted Path**: complete the chosen section's proof surface before claiming mastery

## Companion Docs

| Document | Purpose |
| --- | --- |
| [README.md](./README.md) | project overview and quick start |
| [ARCHITECTURE.md](./ARCHITECTURE.md) | curriculum source of truth |
| [CURRICULUM-BLUEPRINT.md](./CURRICULUM-BLUEPRINT.md) | teaching contract and README-first standard |
| [ROADMAP.md](./ROADMAP.md) | beta completion and RC priorities |
| [docs/PROGRESSION.md](./docs/PROGRESSION.md) | milestone progression and stage overview |

## Bottom Line

Choose the phase that matches your real background, read the section `README.md` first, then run the linked code and complete the proof surfaces.
