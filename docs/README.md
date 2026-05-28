# Go Engineer Documentation

This directory contains the maintainer documentation for the Go Engineer curriculum repository.

The active repository architecture is:

```text
metadata/     Source of truth for graph, concepts, projects, assessments, contracts, migration.
curriculum/   Learner-facing content: READMEs, code, tests, labs, projects, assessments, diagrams, assets.
tools/        Validation, generation, audit, migration, and authoring automation.
docs/         Maintainer documentation and quality standards.
dist/         Generated release artifacts only. Never hand-edit.
```

## Documentation map

| Document | Purpose |
|---|---|
| [`architecture.md`](./architecture.md) | Repository architecture, directory responsibilities, path conventions, and naming rules. |
| [`metadata-contract.md`](./metadata-contract.md) | Metadata file ownership, required fields, graph rules, path policy, and schema expectations. |
| [`content-quality-standard.md`](./content-quality-standard.md) | World-class lesson, lab, project, assessment, and module quality requirements. |
| [`validation-backbone.md`](./validation-backbone.md) | Validation profiles, strict gates, CI usage, and what must pass before release. |
| [`authoring-guide.md`](./authoring-guide.md) | How to create or update modules, lessons, labs, projects, and assessments. |
| [`migration-guide.md`](./migration-guide.md) | v2 to v3 migration policy and how to preserve legacy coverage. |
| [`release-process.md`](./release-process.md) | Release checklist, artifact generation, tagging, and rollback rules. |
| [`contributor-guide.md`](./contributor-guide.md) | Branch, issue, PR, commit, and local workflow rules. |
| [`review-process.md`](./review-process.md) | Findings-first review method and quality bar for approving changes. |

## Source-of-truth order

When files disagree, use this order:

1. `metadata/schema.v3.json`
2. `metadata/workspace.json`
3. `metadata/path.core.json` and `metadata/path.electives.json`
4. `metadata/concepts.json`
5. `metadata/projects.json` and `metadata/assessments.json`
6. `metadata/crossrefs.json`
7. `metadata/readme.contracts.json`
8. `docs/*.md`
9. learner-facing `curriculum/**/*.md`

The docs explain policy. The metadata defines the executable curriculum contract.

## Required validation

Run these before a change is considered complete:

```bash
make validate-metadata
make validate-repository
make validate-release
```

If `make` is unavailable, use the direct commands:

```bash
go run ./tools/validate/curriculum validate-metadata
go run ./tools/validate/curriculum validate-repository --strict-repository
go run ./tools/validate/curriculum validate-all --strict-repository
```

## Maintainer principle

Do not accept "looks good" as proof.

A curriculum item is complete only when metadata, content, code, tests, assets, assessments, and validation all agree.
