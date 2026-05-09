# Scripts

This directory contains utility scripts for validating and maintaining The Go Engineer curriculum.

## Curriculum Validation

The primary tool in this directory is the curriculum validator, which ensures that all content adheres to the `curriculum.v2.json` schema and architecture guidelines.

### `validate_curriculum.go`

This Go script performs a comprehensive check of the curriculum state:

- Validates all paths and references in `curriculum.v2.json`
- Checks that all referenced files exist
- Verifies that `NEXT UP:` footers in lesson files correctly match the curriculum progression
- Validates the `README.md` structure for foundation lessons and exercises (`s00`-`s04`)
- Verifies that foundation `Visual Model` sections contain Mermaid diagrams
- Ensures run-mode foundation lessons include a runnable `main.go`
- Confirms text encoding and section labels

**Usage:**

```bash
go run ./scripts/validate_curriculum.go
```

## Maintainer Automation Scripts

Automation tools for managing GitHub issues, labels, and workflow automation have been moved to the `maintainer-scripts/` subdirectory to separate maintainer tools from curriculum scripts.

See [maintainer-scripts/GITHUB_AUTOMATION_GUIDE.md](./maintainer-scripts/GITHUB_AUTOMATION_GUIDE.md) for full documentation on how to use those tools.
