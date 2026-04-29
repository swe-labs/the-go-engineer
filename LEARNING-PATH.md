# The Go Engineer Learning Path

This guide explains how to move through the stable v2.1 curriculum.

The source of truth for section ownership is [ARCHITECTURE.md](./ARCHITECTURE.md).

## Phase Overview

| Phase | Sections | Focus | Progress |
| --- | --- | --- | --- |
| 0 | s00 | machine foundation | 0% to 5% |
| 1 | s01-s04 | language foundation | 5% to 52% |
| 2 | s05-s08 | engineering core | 52% to 87% |
| 3 | s09-s10 | systems engineering | 87% to 96% |
| 4 | s11 | flagship integration | 96% to 100% |

## Recommended Paths

### Full Path

Use this path if you are new to programming or want the complete curriculum.

Rules:

- start at `s00`
- follow sections in order through `s11`
- run each lesson command
- complete section proof surfaces before moving on

### Bridge Path

Use this path if you already program but are new to Go.

Rules:

- skim `s00` if you already understand process, memory, terminal, and OS basics
- complete `s01` as a Go toolchain check
- move carefully through `s02` to `s04`; these establish Go-specific habits
- do not skip proof surfaces in `s05` through `s08`
- complete `s09` through `s11` in order

### Targeted Path

Use this path if you are returning to strengthen a specific area.

Rules:

- choose one section deliberately
- read that section README first
- check prerequisites honestly
- complete that section proof surface before claiming the skill

## Entry Points

| Goal | Start here | Prerequisite |
| --- | --- | --- |
| First Go setup and execution | `s01` Getting Started | none |
| Understand the machine model | `s00` How Computers Work | none |
| Strengthen language fundamentals | `s02` Language Basics | `s01` |
| Improve function and error design | `s03` Functions and Errors | `s02` |
| Improve type design | `s04` Types and Design | `s02`, `s03` |
| Build CLI and filesystem tools | `s05` Packages, I/O and CLI | `s03`, `s04` |
| Build backend and database code | `s06` Backend, APIs and Databases | `s05`, explicit error handling |
| Improve concurrency | `s07` Concurrency | functions, errors, packages, interfaces |
| Improve testing and profiling | `s08` Quality and Testing | `s05`, `s06`, `s07` |
| Study architecture and security | `s09` Architecture and Security | Phase 2 |
| Study production operations | `s10` Production Operations | `s06`, `s07`, `s08` |
| Integrate the whole system | `s11` Opslane Flagship | all earlier phases |

## Proof Expectations

Fast paths are allowed. Proof is not optional.

- Full Path: complete every required milestone and proof surface.
- Bridge Path: complete proof surfaces even when syntax is familiar.
- Targeted Path: complete the chosen section's proof surface before moving on.

## Navigation Rule

When a lesson mentions a concept that is introduced elsewhere, follow the inline reference:

- forward references name the later section or lesson that teaches the topic in detail
- backward references name the earlier lesson that established the idea
- sibling references point to nearby tracks that reuse the concept

## Companion Docs

| Document | Purpose |
| --- | --- |
| [README.md](./README.md) | project overview and quick start |
| [ARCHITECTURE.md](./ARCHITECTURE.md) | locked curriculum structure |
| [CURRICULUM-BLUEPRINT.md](./CURRICULUM-BLUEPRINT.md) | teaching contract |
| [docs/PROGRESSION.md](./docs/PROGRESSION.md) | phase and milestone visualization |
| [TESTING-STANDARDS.md](./TESTING-STANDARDS.md) | verification expectations |
