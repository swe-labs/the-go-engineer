# Scripts

This directory contains utility scripts for validating and maintaining The Go Engineer curriculum.

## Curriculum Validation

The primary tool is the curriculum validator.

### `validate_curriculum.go`

This script delegates to the canonical v3 curriculum validator (`internal/tools/curriculum`), checking:

- modular JSON manifests in `curriculum/` (`path.core.json`, `projects.json`, `assessments.json`, etc.)
- physical files, lessons, labs, project solutions, and starters
- 17 automated curriculum quality gates including graph closure, zero-magic, and artifact parity

Usage:

```bash
go run ./scripts/validate_curriculum.go
# or directly:
go run ./internal/tools/curriculum validate-all
```

## CI-equivalent Local Verification

Use this command set before final review:

```bash
go build ./...
go vet ./...
unformatted=$(gofmt -l . | grep -v '_starter' || true); test -z "$unformatted" || (echo "$unformatted" && exit 1)
go mod tidy
git diff --exit-code -- go.mod go.sum
go test ./...
go test -race ./...
go test -coverprofile coverage.out ./...
go run ./scripts/validate_curriculum.go
```

For benchmark-related changes:

```bash
go test -bench=. -benchmem -count=1 ./curriculum/modules/13-performance-memory-engineering/projects/profile-and-fix-lab/_solution
```

## Maintainer Automation Scripts

Maintainer tools for labels, issues, and repository setup live in `scripts/maintainer-scripts/`.

See [maintainer-scripts/GITHUB_AUTOMATION_GUIDE.md](./maintainer-scripts/GITHUB_AUTOMATION_GUIDE.md).
