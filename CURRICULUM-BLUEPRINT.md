# The Go Engineer Curriculum Blueprint (v3)

This document defines how the canonical v3 architecture behaves as an instructional engineering system.

If this file and [ARCHITECTURE.md](./ARCHITECTURE.md) disagree on public structure, `ARCHITECTURE.md` wins.

## Core Promise

The curriculum guides a learner from fundamentals to production Go backend systems engineering with deep understanding, rigorous proof, and production reliability.

The expected progression is:

- understand what a line of Go does at the CPU, operating system, and runtime level
- explain why that design exists and evaluate architectural trade-offs
- predict what can break under concurrency, network failure, or load
- prove behavior with automated tests, race detection, benchmarks, and profiles
- integrate isolated lessons into production-grade systems culminating in the Opslane Flagship Platform

## Architecture & Structure

The canonical v3 curriculum consists of **18 core modules** organized into progressive engineering phases:

- Phase 0: Orientation & Foundations (Modules 00–02)
- Phase 1: Go Core & Semantics (Modules 03–04)
- Phase 2: Engineering Core (Modules 05–07)
- Phase 3: Backend & Persistence (Modules 08–09)
- Phase 4: Security & Concurrency (Modules 10–11)
- Phase 5: Production Systems (Modules 12–15)
- Phase 6: Career & Portfolio (Module 16)
- Phase 7: Flagship Capstone (Module 18)

Companion tracks include the Companion Full-Stack Frontend Track (`curriculum/companion-frontend/`) and 20 Advanced Electives (`curriculum/electives/`).

## Teaching Rules

### README first, code second

Every learner-facing lesson teaches through `README.md` before source code.
The learner should understand the mission, prerequisites, mental model, and machine execution view before reading `main.go`.

### Code is never skipped

Do not replace code with prose. Explain the code, run the code, then modify the code.

### Zero magic

Each section teaches only what has been earned. If a concept depends on later ideas, either defer it or add an explicit forward reference.

### Explanations answer how, why, and what changes

Good teaching surfaces explain:

- what the line or block does
- why it exists
- what would change if the design changed
- what failure mode or trap a learner is likely to encounter

### Engineering depth is stage-aware

Foundational modules establish concrete mental models (stack vs heap, pointer semantics, goroutine mechanics). Production modules focus on systems architecture, observability, fault tolerance, and multi-tenant security.

## Canonical Contract

### Required README Sections

Instructional READMEs conform to the contract defined in `curriculum/readme.contracts.json`:

1. `## Mission`
2. `## Prerequisites`
3. `## Mental Model`
4. `## Visual Model`
5. `## Machine View`
6. `## Run Instructions`
7. `## Code Walkthrough`
8. `## Try It`
9. `## In Production`
10. `## Thinking Questions`
11. `## Next Step`

For exercises and projects:
- `## Solution Walkthrough`
- `## Verification Surface`

## Quality Gates & Verification

All curriculum changes must pass the 17 automated quality gates:

```bash
go run ./internal/tools/curriculum validate-all
go test ./...
go vet ./...
```
