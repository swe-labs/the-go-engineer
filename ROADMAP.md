# Project Status and Maintenance Plan (v3)

This document describes the project state and maintenance plan for **The Go Engineer v3**.

If this document and [ARCHITECTURE.md](./ARCHITECTURE.md) disagree on public curriculum structure, `ARCHITECTURE.md` is authoritative.

## Current State

The v3 curriculum is implemented on branch `release/v3`.

The repository features:

- 18 core modules organized canonically in `curriculum/modules/`
- 421 core path items in `curriculum/path.core.json`
- Flagship Opslane Capstone in Module 18 with 15 progressive milestones
- 20 Advanced Electives in `curriculum/electives/advanced-electives/`
- Companion Full-Stack Frontend Track in `curriculum/companion-frontend/`
- Complete reference solutions (`_solution/`) and student starter templates (`_starter/`)
- 17 automated curriculum quality gates in `internal/tools/curriculum`

## Branch Model

| Branch | Role |
|---|---|
| `main` | active development and next-generation feature work |
| `release/v3` | active v3 release preparation and maintenance line |
| `release/v2` | legacy v2.1.x maintenance line (archived) |

## Maintenance Rules

Normal maintenance may:

- refine lesson explanations, tests, and starter code
- improve public documentation and diagrams
- enhance automated quality tooling in `internal/tools/curriculum`
- address community feedback, errata, and dependency updates

Normal maintenance must not:

- introduce breaking changes to the 18-module progression without schema migration
- bypass zero-magic pedagogical standards
- weaken curriculum quality gates to mask drift
- publish unverified commits to release branches

## Release Readiness

A release is ready only when:

- public documentation and `curriculum/` metadata agree
- `go build ./...` passes
- `go vet ./...` passes
- `gofmt -l .` reports 0 unformatted files
- `go mod tidy` produces no diff in `go.mod` or `go.sum`
- `go test ./...` passes 100%
- `go test -race ./curriculum/modules/11-lifecycle-context-concurrency/... ./curriculum/modules/18-flagship-opslane/...` reports 0 races
- `go run ./internal/tools/curriculum validate-all` passes all 17 quality gates
- GitHub Actions CI is green
