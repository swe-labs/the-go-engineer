# Tools

This directory contains all automation for the Go Engineer curriculum.

It is intentionally neutral: tools are named by responsibility, not by the human or assistant that uses them.

```text
validate/   strict validation gates for metadata and learner-facing repository files
generate/   deterministic generators from metadata to curriculum files and dist artifacts
audit/      reporting tools for gaps, quality, consistency, and completion
migrate/    v2-to-v3 migration mapping and reporting utilities
authoring/  reusable authoring workflows for humans and assistants
```

## Required layout

The tools assume the final repository structure:

```text
metadata/     source of truth
curriculum/   learner-facing curriculum files
tools/        automation
docs/         maintainer documentation
dist/         generated release artifacts
```

## Common commands

```bash
go run ./tools/validate/curriculum validate-metadata
go run ./tools/validate/curriculum validate-repository --strict-repository
go run ./tools/validate/curriculum validate-all --strict-repository

python3 tools/audit/audit_completion.py --metadata-dir metadata --curriculum-dir curriculum --output dist/completion-report.json
python3 tools/generate/generate_full_snapshot.py --metadata-dir metadata --curriculum-dir curriculum --output dist/curriculum.v3.full.generated.json
python3 tools/migrate/migration_report.py --metadata-dir metadata --output dist/migration-report.json
```

## Rule

Python tools may generate, migrate, and audit. The Go validator is the release gate.
