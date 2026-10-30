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

## Canonical Contract & Anti-Overengineering Directive

### README Section Architecture

Instructional READMEs conform to the canonical 17-heading contract defined in [`curriculum/readme.contracts.json`](./curriculum/readme.contracts.json). To guarantee high pedagogical value without artificial ceremony, headings are divided into Core Required and Contextual sections:

#### Core Required Headings (Mandatory across all lessons)
1. `## Mission` — Explicit declaration of capabilities earned and verified.
2. `## Prerequisites` — Transitive dependencies and foundational concepts.
3. `## Mental Model` — Transparent mechanical foundation before code.
4. `## Code Walkthrough` — Thorough line-by-line inspection (`## Solution Walkthrough` for exercises/projects).
5. `## In Production` — Operational realities, telemetry, and failure handling.
6. `## Next Steps` — Forward progression and curriculum linkage.

#### Contextual Headings (Applied based on instructional necessity)
Applied selectively based on the lesson's subtype profile (`conceptual`, `pattern`, `deep-dive`, `exercise`, `project`, `capstone`) and real cognitive necessity:
- `## Visual Model` (Mermaid diagram of flow/lifecycle)
- `## Memory Layout` (Stack/heap byte layout and pointer indirection)
- `## Physical Timeline` (Sequence diagram of concurrent interactions)
- `## Machine View` (OS syscalls, assembly snippets, or runtime scheduling)
- `## Common Failure Modes` (Concrete traps, compiler errors, or panic scenarios)
- `## Edge Cases` (Zero values, boundary limits, or race conditions)
- `## Benchmark & Allocation Analysis` (Allocations per op and throughput)
- `## Production Notes` (Hardened configurations, security caveats)
- `## Review Questions & Explanations` (Self-check conceptual prompts)
- `## Verification Surface` (Automated test suites, assertions; required for exercises/projects)
- `## Portfolio Artifact` (Deliverables and defense rubrics; required for capstones)

## Quality Gates & Verification

All curriculum changes must pass the 17 automated quality gates:

```bash
go run ./internal/tools/curriculum validate-all
go test ./...
go vet ./...
```
