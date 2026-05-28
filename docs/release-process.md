# Release Process

This document defines how to release the curriculum.

A release is a validated snapshot of metadata, learner-facing content, tools, docs, and generated artifacts.

## Release requirements

A release is allowed only when:

- metadata validation passes
- strict repository validation passes
- release validation passes
- migration coverage has zero unmapped required legacy items
- completion report shows required core completion
- generated artifacts are reproducible
- no secrets are committed
- no forbidden folder names exist
- root is clean
- release notes are generated or reviewed

## Release commands

```bash
make validate-metadata
make validate-repository
make validate-release
make release-artifacts
```

Direct equivalents:

```bash
go run ./tools/validate/curriculum validate-metadata
go run ./tools/validate/curriculum validate-repository --strict-repository
go run ./tools/validate/curriculum validate-all --strict-repository
```

## Release artifacts

Generated artifacts belong in `dist/`:

```text
dist/curriculum.v3.full.generated.json
dist/validation-report.json
dist/completion-report.json
dist/migration-report.json
dist/release-notes.md
```

Do not hand-edit `dist/` files.

Regenerate them from metadata and curriculum sources.

## Versioning

Use semantic version-style tags:

```text
v3.0.0
v3.0.1
v3.1.0
```

Suggested meaning:

| Version type | Use |
|---|---|
| Patch | docs fixes, validator bugfixes, typo fixes, non-structural content corrections |
| Minor | new lessons, stronger projects, new validation checks, new electives |
| Major | architecture changes, path contract changes, graph restructuring |

## Release workflow

1. Create release branch.
2. Confirm metadata is stable.
3. Confirm repository validation passes.
4. Generate artifacts.
5. Review completion report.
6. Review migration report.
7. Review release notes.
8. Tag release.
9. Run GitHub release workflow.
10. Verify downloadable artifacts.
11. Announce release.

## Pre-release checklist

- [ ] `metadata/VALIDATION.metadata.json` is current.
- [ ] `metadata/legacy/unmapped-v2-report.json` shows zero unmapped items.
- [ ] All required curriculum files exist.
- [ ] All Go code is formatted.
- [ ] All tests pass.
- [ ] All project rubrics exist.
- [ ] All assessments have questions, answer keys, and rubrics.
- [ ] All referenced assets exist.
- [ ] All docs match current architecture.
- [ ] `dist/` regenerated.
- [ ] No `.env` files tracked.
- [ ] No `AGENTS.md` root file tracked.
- [ ] No forbidden folder names.

## Rollback

Rollback if:

- release artifacts are incomplete
- validation failed after tag
- serious learner-facing path errors are discovered
- secret files were included
- license or attribution issue is discovered

Rollback steps:

1. Remove or mark release as withdrawn.
2. Revoke affected artifacts.
3. Open a fix issue.
4. Patch source files.
5. Regenerate artifacts.
6. Re-tag only when clean.

## Release notes standard

Release notes should include:

```text
## Highlights
## Added
## Changed
## Fixed
## Validation
## Migration
## Known limitations
```

Known limitations must not hide release blockers. If a limitation breaks required core learning, the release should not ship.
