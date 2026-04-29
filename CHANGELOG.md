# Changelog

All notable stable changes to The Go Engineer are documented here.

Format: `Version - Date`.

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
- `go test -coverprofile=coverage.out ./...`
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
