# Contributor Guide

This guide defines how to contribute to the Go Engineer curriculum.

## Before you start

Read:

1. `docs/architecture.md`
2. `docs/metadata-contract.md`
3. `docs/content-quality-standard.md`
4. `docs/validation-backbone.md`
5. this guide

## Branch naming

Use short, descriptive branch names:

```text
feat/module-08-handler-lifecycle
fix/assessment-module-09-rubric
docs/metadata-contract
test/repository-validator-assets
chore/root-cleanup
migration/v2-api-track
```

Allowed prefixes:

```text
feat/
fix/
docs/
test/
chore/
refactor/
security/
release/
migration/
```

## Commit messages

Use bracketed type prefixes:

```text
[FEAT] add handler lifecycle lesson
[FIX] correct SQL assessment target
[DOCS] document metadata path policy
[TEST] add README contract validation
[CHORE] clean root documentation
[MIGRATION] map v2 API lessons to electives
```

Allowed types:

```text
[FEAT]
[FIX]
[DOCS]
[TEST]
[CHORE]
[REFACTOR]
[SECURITY]
[RELEASE]
[MIGRATION]
```

## PR expectations

A PR must include:

- summary
- affected modules/lessons/projects/assessments
- metadata changes
- learner-facing changes
- validation commands run
- screenshots or excerpts when docs/assets changed
- known risks
- follow-up work, if any

## Required local validation

For metadata-only changes:

```bash
make validate-metadata
```

For curriculum changes:

```bash
make validate-metadata
make validate-repository
```

For release or architecture changes:

```bash
make validate-release
```

## Contribution boundaries

Do not:

- add root-level architecture docs for old versions
- add `.env` files
- add `AGENTS.md`
- add tool-branded folders
- edit `dist/` by hand
- bypass metadata when adding learner-facing files
- weaken validation to pass incomplete content

Do:

- update metadata first
- keep paths canonical
- write README-first lessons
- add tests when behavior matters
- add rubrics for projects
- add assessment evidence requirements
- preserve migration traceability

## Issue template

Use this structure for non-trivial work:

```text
## Scope
What will change?

## Affected files
Which metadata/content/tool/docs files?

## Learner impact
How does this improve the curriculum?

## Validation plan
Which commands will run?

## Risks
What might break?
```

## PR review readiness

A PR is ready when:

- validation passes
- content quality checklist is complete
- no placeholder text remains
- all linked files exist
- review comments are addressed
- generated artifacts are reproducible when relevant
