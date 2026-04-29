# Scripts

This directory contains utility scripts for validating and maintaining The Go Engineer curriculum.

## Curriculum Validation

The primary tool is the curriculum validator.

### `validate_curriculum.go`

This script checks curriculum state, including:

- paths and references in `curriculum.v2.json`
- referenced files
- `NEXT UP:` footers in lesson files
- README structure for covered lesson categories
- required runnable lesson surfaces
- text encoding and section labels

Usage:

```bash
go run ./scripts/validate_curriculum.go
```

## CI-equivalent Local Verification

Use this command set before final review:

```bash
go build ./...
go vet ./...
unformatted=$(gofmt -l .); test -z "$unformatted" || (echo "$unformatted" && exit 1)
go mod tidy
git diff --exit-code -- go.mod go.sum
go test ./...
go test -race ./...
go test -coverprofile coverage.out ./...
go run ./scripts/validate_curriculum.go
```

For benchmark-related changes:

```bash
go test -bench=. -benchmem -count=1 ./08-quality-test/01-quality-and-performance/testing/benchmarks/
```

## Maintainer Automation Scripts

Maintainer tools for labels, issues, and repository setup live in `scripts/maintainer-scripts/`.

See [maintainer-scripts/GITHUB_AUTOMATION_GUIDE.md](./maintainer-scripts/GITHUB_AUTOMATION_GUIDE.md).
