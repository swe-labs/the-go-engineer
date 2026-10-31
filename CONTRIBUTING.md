# Contributing to The Go Engineer

This guide defines the required workflow for contributions to **The Go Engineer v3**.

## Source of Truth

Read these documents before changing curriculum, code, validation, or release surfaces:

| Document | Purpose |
|---|---|
| [ARCHITECTURE.md](./ARCHITECTURE.md) | Canonical v3 public structure |
| [`curriculum/`](./curriculum/README.md) | Modular curriculum metadata manifests |
| [`curriculum/zero-magic-authoring-spec.md`](./curriculum/zero-magic-authoring-spec.md) | Zero-magic teaching specification |
| [CODE-STANDARDS.md](./CODE-STANDARDS.md) | Go and teaching-code standards |
| [TESTING-STANDARDS.md](./TESTING-STANDARDS.md) | Verification and test expectations |
| [RELEASE.md](./RELEASE.md) | Release and branch workflow |

If documents disagree on public curriculum structure, `ARCHITECTURE.md` wins.

## Architecture Rule

The canonical v3 architecture consists of 18 core modules (`00`–`16` plus `18-flagship-opslane`) under `curriculum/modules/`, accompanied by the Companion Full-Stack Frontend Track and 20 Advanced Electives.

Allowed work includes:

- lesson corrections, depth improvements, and zero-magic explanations
- project tests, starters, and reference solutions
- curriculum metadata alignment across `curriculum/*.json`
- quality tooling improvements in `internal/tools/curriculum`
- Opslane milestone enhancements inside `curriculum/modules/18-flagship-opslane/`
- public documentation improvements

## Issue First

Every non-trivial change needs an issue before implementation.

Issue title format:

```text
[TYPE] short description
```

Allowed types:

```text
[FEAT] [FIX] [DOCS] [TEST] [CHORE] [REFACTOR] [SECURITY] [RELEASE]
```

Issue body should include:

- affected section, if applicable
- affected lesson IDs, if applicable
- expected files or folders
- validation plan
- risk notes

## Branches

| Branch | Purpose |
|---|---|
| `main` | active development and next-generation feature work |
| `release/v3` | active v3 release preparation and maintenance line |

Branch from the line that should receive the change.

Examples:

```bash
git switch main
git pull --ff-only origin main
git switch -c feat/new-feature
```

```bash
git switch release/v3
git pull --ff-only origin release/v3
git switch -c fix/patch-fix
```

## Pull Requests

Open a draft PR early and link the issue with `Closes #<issue>`.

The PR should include:

- scope
- affected sections or lessons
- validation commands run
- risk notes
- tracking metadata

Keep the PR in draft until local validation, self-review, and obvious fixes are complete.

## Commits

Use bracketed commit messages:

```text
[TYPE] short imperative description
```

Examples:

```text
[DOCS] align README with v3.0.0 architecture
[FIX] replace blocking REST API test with httptest
[TEST] add repository behavior coverage
[CHORE] tighten curriculum validator
```

Keep commits logical. Do not mix unrelated formatting, docs, and behavior changes unless the PR is explicitly a repository-wide cleanup.

## Lesson Contract & Anti-Overengineering Standards

All instructional READMEs conform to the canonical pedagogical contract defined in [`curriculum/readme.contracts.json`](./curriculum/readme.contracts.json) and enforced by the zero-magic curriculum architecture.

### 1. Core Required Headings (Mandatory for All Lessons)
Every lesson README must include these 6 core sections in order:
1. `## Mission` — Clear statement of what the learner will build, understand, and prove.
2. `## Prerequisites` — Explicit concepts or previous lesson IDs required.
3. `## Mental Model` — Zero-magic foundation of how the machine/runtime operates.
4. `## Code Walkthrough` — Line-by-line explanation of the code (`## Solution Walkthrough` for exercises/projects).
5. `## In Production` — Real-world trade-offs, operational considerations, or telemetry.
6. `## Next Steps` — Direct link and progression to the next lesson or project milestone.

### 2. Contextual Headings (Applied by Subtype Profile)
To prevent cognitive overload and avoid artificial padding, the remaining 11 headings are included only when justified by cognitive necessity and the lesson's subtype profile:
- `## Visual Model` — Mermaid diagram visualizing lifecycle, topology, or data flow.
- `## Memory Layout` — Stack vs heap allocations, struct padding, or pointer indirection.
- `## Physical Timeline` — Execution trace or concurrency sequence diagram.
- `## Machine View` — OS syscalls, assembly snippets, or runtime scheduler events.
- `## Common Failure Modes` — Root causes, compiler errors, or panic scenarios.
- `## Edge Cases` — Boundary conditions, zero values, or overflow behaviors.
- `## Benchmark & Allocation Analysis` — Benchmarks and zero-alloc profiling.
- `## Production Notes` — Hardened configuration, rate limits, or deployment caveats.
- `## Review Questions & Explanations` — Self-assessment questions with answers.
- `## Verification Surface` — Automated test commands, test cases, and assertions (required for exercises/projects).
- `## Portfolio Artifact` — Deliverables, architecture notes, and defense rubrics (required for capstones).

### 3. Subtype Profiles
- `conceptual` — Focuses on mental models, core walkthrough, and production context.
- `pattern` — Emphasizes visual models, common failure modes, and idioms.
- `deep-dive` — Adds memory layouts, machine views, and allocation analysis.
- `exercise` — Replaces code walkthrough with solution walkthrough and verification surface.
- `project` — Multi-file architecture with verification surface and test suite.
- `capstone` — Comprehensive systems capstone with portfolio artifacts and operational rubrics.

Every lesson `main.go` should include:

- copyright and license header
- module, lesson title, and level
- `WHAT YOU'LL LEARN`
- `WHY THIS MATTERS`
- exact `RUN:` command
- readable teaching comments
- `KEY TAKEAWAY`
- `NEXT UP:` footer matching canonical progression in `curriculum/path.core.json`

## Cross-References

When a lesson uses a concept from another point in the curriculum:

- forward reference later concepts briefly and name the future lesson or module
- backward reference previously established ideas when they are reused
- keep cross-references inline with the explanation rather than as detached navigation blocks

## Validation

Run the CI-equivalent bundle before final review:

```bash
go build ./...
go vet ./...
unformatted=$(gofmt -l . | grep -v '_starter' || true); test -z "$unformatted" || (echo "$unformatted" && exit 1)
go mod tidy
git diff --exit-code -- go.mod go.sum
go test ./...
go test -race ./curriculum/modules/11-lifecycle-context-concurrency/... ./curriculum/modules/18-flagship-opslane/...
go run ./internal/tools/curriculum validate-all
```

For benchmark-related changes:

```bash
go test -bench=. -benchmem -count=1 ./curriculum/modules/13-performance-memory-engineering/projects/profile-and-fix-lab/_solution
```

## Review Standard

Before final review:

- canonical v3 18-module architecture remains intact
- manifests in `curriculum/` match files and lesson metadata
- README, source, tests, and module maps agree
- `NEXT UP:` footers match curriculum metadata in `curriculum/path.core.json`
- Go code is formatted and idiomatic
- tests prove behavior where behavior should be provable
- all 17 automated curriculum quality gates pass (`go run ./internal/tools/curriculum validate-all`)
- P0/P1/P2 findings are fixed or explicitly documented

Final squash merge is maintainer-controlled.
