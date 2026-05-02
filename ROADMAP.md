# Project Status and Maintenance Plan

This document describes the stable v2.1 project state and the maintenance plan for future work.

If this document and [ARCHITECTURE.md](./ARCHITECTURE.md) disagree on public curriculum structure, `ARCHITECTURE.md` is authoritative.

## Current State

The v2.1 architecture is stable and locked.

Current stable release: `v2.1.1`

The repository now has:

- 5 public learning phases
- 12 public sections, `s00` through `s11`
- 215 registered curriculum items in `curriculum.v2.json`
- strict curriculum validation through `go run ./scripts/validate_curriculum.go`
- CI-equivalent local verification, including race tests
- an integrated flagship backend path in `s11`

## Branch Model

| Branch | Role |
| --- | --- |
| `main` | active post-v2.1 implementation and integration line |
| `release/v2` | stable v2.1.x maintenance line |
| `release/v1` | stable v1 maintenance line |

Work that should affect stable v2.1 users must land on or be backported to `release/v2`. Work that expands future depth without changing the stable release line belongs on `main`.

## Stable Snapshot

| Area | Status | Standard |
| --- | --- | --- |
| Public architecture | Stable | 12 sections locked by `ARCHITECTURE.md` |
| Curriculum registry | Stable | 12 sections and 215 items validated |
| Lesson structure | Stable | README-first contract enforced by validator |
| Section READMEs | Stable | section maps align with architecture and curriculum metadata |
| Validator | Stable | strict checks fail on structural drift |
| Tests | Stable | build, vet, unit tests, race tests, and coverage generation pass locally |
| Opslane flagship | Stable baseline | integrated backend capstone remains open for future depth on `main` |
| Workflow docs | Stable | issue, PR, review, and release workflow documented |

## Version Plan

| Version | Line | Purpose |
| --- | --- | --- |
| `v2.1.1` | `release/v2` | stable curriculum completion and public documentation alignment |
| `v2.1.x` | `release/v2` | low-risk fixes, docs corrections, validator fixes, dependency/security updates |
| post-v2.1 | `main` | future lesson depth, flagship expansion, and tooling improvements |

## Maintenance Rules

Normal maintenance may:

- correct lesson explanations, tests, or examples
- improve public documentation
- strengthen validators and CI
- fix learner navigation, metadata, or README/source mismatches
- deepen Opslane without changing the public 12-section spine

Normal maintenance must not:

- add a new public root section
- rename canonical section folders
- move a lesson prefix to another section without explicit architecture approval
- weaken validator checks to hide drift
- publish a release without a clean verification record

## Release Readiness

A stable release is ready only when:

- public docs, `ARCHITECTURE.md`, and `curriculum.v2.json` agree
- `go build ./...` passes
- `go vet ./...` passes
- gofmt check passes
- `go mod tidy` produces no `go.mod` or `go.sum` diff
- `go test ./...` passes
- `go test -race ./...` passes
- `go test -coverprofile=coverage.out ./...` passes
- `go run ./scripts/validate_curriculum.go` reports success
- GitHub CI is green
- the release branch and tag are created from the verified commit

## Current Focus

After `v2.1.1`, development should stay bounded:

1. maintain the stable v2.1.x line with small, auditable fixes
2. improve Opslane depth on `main`
3. keep validator coverage strict
4. preserve public documentation as the accurate project mirror
5. avoid architecture churn unless a future major version is explicitly planned
