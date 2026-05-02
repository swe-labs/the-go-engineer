# The Go Engineer v2.1.1

Release date: 2026-04-29

Target branch: `release/v2`

## Summary

`v2.1.1` is the stable curriculum completion release for the v2.1 line. It aligns public documentation, curriculum metadata, section maps, validator behavior, and verification evidence around the locked 5-phase, 12-section architecture.

The release keeps the public architecture unchanged and focuses on correctness, learner navigation, validation strictness, and release readiness.

## Included Changes

- Confirms 12 v2 sections and 215 registered curriculum items.
- Aligns `ARCHITECTURE.md`, `curriculum.v2.json`, section READMEs, and public documentation.
- Replaces the blocking REST API test with deterministic `httptest` coverage.
- Tightens REST API invalid ID and response handling.
- Fixes API.9 output and context-cancellation behavior.
- Makes curriculum validation stricter instead of normalizing drift.
- Enforces stable section status, expected section outputs, canonical section README labels, and engineering README heading order.
- Updates public README, roadmap, learning path, release guide, contribution guide, changelog, and progression documentation for the stable v2.1.x line.

## Validation

Local verification passed:

```text
go build ./...
go vet ./...
gofmt -l .
go mod tidy
git diff --exit-code -- go.mod go.sum
go test ./...
go test -race ./...
go test -coverprofile=coverage.out ./...
go run ./scripts/validate_curriculum.go
git diff --check
```

Curriculum validator:

```text
Success! 597 files with run commands validated, and 12 v2 sections plus 215 v2 items checked.
```

## Compatibility

- Public architecture remains v2.1.
- No public root sections are added, removed, renamed, or reordered.
- Stable maintenance continues on `release/v2`.
- Future implementation depth continues on `main`.

## Known Notes

- Windows users need a CGO-capable C compiler for `go test -race ./...` and SQLite-backed paths.
- Generated local files such as `coverage.out` are verification artifacts and are not release assets.
