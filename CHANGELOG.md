# Changelog

All notable stable changes to The Go Engineer are documented here.

Format: `Version - Date`.

## v3.0.0 - 2026-10-23

### Added

- **Canonical 18-Module Core Curriculum:** Rebuilt the curriculum around an 18-module progressive systems engineering path (421 core path items in `curriculum/path.core.json`).
- **Complete Reference Solutions and Starter Scaffolds:** Added reference implementations (`_solution/`), student starter templates (`_starter/`), and automated test suites (`main_test.go`) across all core portfolio projects (Modules 02–16).
- **Flagship Capstone Project (Opslane):** Integrated 15-milestone SaaS backend capstone (Module 18) demonstrating multi-tenant data models, PostgreSQL migrations, JWT authentication, RBAC authorization, event bus, worker pools, caching, observability, and graceful shutdown.
- **Companion Frontend Track:** Added 4-stage companion full-stack track under `curriculum/companion-frontend` covering Semantic HTML/CSS, Vanilla DOM & Fetch, Templ & Tailwind CSS, and HTMX & Alpine.js.
- **Zero-Magic Pedagogical Standards:** Added 92 golden lessons with verified machine-level mental models, stack/heap allocation explanations, and explicit failure-mode engineering.
- **Automated Curriculum Quality Gates:** Added 17 automated quality gates in `internal/tools/curriculum` enforcing schema integrity, cycle-free graph closure, zero-magic standards, tag uniqueness, concept ownership, failure engineering coverage, README pedagogical contracts, and workspace synchronization.
- **Portfolio Auditor Tool:** Added automated CLI audit engine in Module 16 (`portfolio-review-pack`) to score project completeness and provide technical interview defense guides.

### Changed

- Transitioned canonical curriculum metadata from legacy monolithic structure to modular JSON files (`path.core.json`, `path.electives.json`, `projects.json`, `assessments.json`, `crossrefs.json`, `concepts.json`, `failures.json`, `readme.contracts.json`, `workspace.json`, `schema.v3.json`).
- Updated active schema format to `3.1.0` and promoted core curriculum version to `3.0.0` stable.
- Rewrote root `README.md` and `curriculum/README.md` with transparent 52-week pacing schedules and workload estimates (435–600 learner hours).
- Modernized Go dependencies and tooling across the repository for Go 1.24+ compatibility.

### Fixed

- Fixed unformatted Go files across project solutions and assessment milestones to conform with `gofmt` CI rules.
- Fixed dependency declarations in `go.mod` and `go.sum` via `go mod tidy`.
- Corrected physical path mappings in `path.core.json` for Module 05 lab 17, orientation lesson 03, and Opslane milestone paths.

### Validation

- `go run ./internal/tools/curriculum validate-all` (17/17 quality gates passed)
- `go test ./...` (all unit and integration tests passed)
- `go vet ./...` (0 warnings/errors)
- `gofmt -l .` (0 unformatted files)
- `go test -race ./curriculum/modules/11-lifecycle-context-concurrency/... ./curriculum/modules/18-flagship-opslane/...` (0 data races)

## v2.1.1 - 2026-04-29

### Added

- Added stricter curriculum validation for locked section outputs, stable section status, canonical section README track labels, and engineering README heading order.
- Added deterministic `httptest` coverage for the REST API exercise so `go test ./...` no longer launches a blocking demo server.
- Added public release notes for the v2.1.1 stable completion release.

### Changed

- Reworked public documentation so README, roadmap, release guide, learning path, contribution guide, and progression docs reflect the current stable v2.1 architecture.
- Updated public status language from migration-oriented wording to stable v2.1.x maintenance wording.
- Aligned curriculum metadata and architecture outputs for s05 and s10.
- Clarified that the Section 09 gRPC material is supporting reference content, not an extra public architecture track.
- Promoted all v2 section metadata to stable status.
- Replaced heading normalization in the validator with strict contract checks.

### Fixed

- Fixed the REST API exercise test hang by replacing `go run main.go` output testing with in-process handler tests.
- Fixed REST API invalid ID behavior and JSON response error handling.
- Fixed API.9 output and context-cancellation behavior.
- Fixed public documentation links that pointed to stale release surfaces.
- Fixed learner README heading drift and exercise verification-section ordering.

### Validation

- `go build ./...`
- `go vet ./...`
- `gofmt -l .`
- `go mod tidy`
- `git diff --exit-code -- go.mod go.sum`
- `go test ./...`
- `go test -race ./...`
- `go test -coverprofile coverage.out ./...`
- `go run ./scripts/validate_curriculum.go`

The curriculum validator reports:

```text
Success! 597 files with run commands validated, and 12 v2 sections plus 215 v2 items checked.
```

## v2.1.0 - 2026-04-22

### Added

- Published the stable v2.1 curriculum line.
- Established the 5-phase, 12-section public architecture.
- Added the `release/v2` stable maintenance line.
- Published the first stable v2.1 release artifact.

### Changed

- Aligned root documentation, roadmap, release workflow, and maintainer workflow with the v2.1 architecture.
- Standardized branch, issue, PR, and commit workflow around bracketed change types.
- Made `curriculum.v2.json` the active machine-readable registry for the v2 line.

## Earlier History

Earlier repository history includes the original v1 curriculum, the v2 planning period, and the v2.1 architecture consolidation. Those changes are preserved in git history. Current public documentation should be treated as authoritative for the stable v2.1 line.
