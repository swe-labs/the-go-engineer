# Validation Backbone

Validation is the quality backbone of the project.

The validator should protect the curriculum from broken paths, weak metadata, missing learner-facing files, vague explanations, and incomplete release artifacts.

## Validation profiles

### Metadata validation

Command:

```bash
go run ./tools/validate/curriculum validate-metadata
```

Checks:

- JSON parsing
- schema compatibility
- module graph
- item graph
- prerequisite resolution
- next item resolution
- concept ownership
- concept reinforcement
- cross-reference targets
- cross-reference reason quality
- project targets
- assessment targets
- README contract metadata
- failure coverage metadata
- migration coverage
- no placeholder/scaffold statuses
- canonical path policy

### Repository validation

Command:

```bash
go run ./tools/validate/curriculum validate-repository --strict-repository
```

Checks actual learner-facing files:

- module READMEs
- lesson READMEs
- lab READMEs
- project READMEs
- assessment READMEs
- Go source files
- tests
- starter directories
- solution directories
- assets and diagrams
- required README headings
- `NEXT UP` correctness
- code formatting
- test presence
- rubric files
- answer keys

### Full validation

Command:

```bash
go run ./tools/validate/curriculum validate-all --strict-repository
```

Runs all gates.

## Strict repository mode

Strict mode is intentionally unforgiving.

It should fail when metadata says a file exists but the file is missing.

This is correct. Metadata completion and full repository completion are separate stages.

## CI workflows

The repository uses four workflow files:

```text
.github/workflows/validate-metadata.yml
.github/workflows/validate-repository.yml
.github/workflows/validate-release.yml
.github/workflows/release.yml
```

### `validate-metadata.yml`

Runs on metadata and validator changes.

Must pass before metadata PRs merge.

### `validate-repository.yml`

Runs on metadata, curriculum, tools, and docs changes.

Must pass before learner-facing content PRs merge.

### `validate-release.yml`

Runs full release gates.

Must pass before release candidates.

### `release.yml`

Builds release artifacts after validation.

## Go validator responsibilities

The Go validator under `tools/validate/curriculum/` is the stable CI-critical gate.

It should prioritize:

- deterministic checks
- strict failure messages
- stable output
- no network dependency
- no hidden state
- no reliance on generated files unless explicitly validating `dist/`

## Python validator responsibilities

Python validators under `tools/validate/repository/` may provide richer file inspections and faster iteration.

They can check:

- README heading depth
- prose placeholder patterns
- diagram references
- markdown links
- code block languages
- exercise folder shape
- content word-count thresholds
- local asset references

Python can be more flexible. Go should be the release gate.

## Failure severity

Use these levels:

| Severity | Meaning | Merge allowed? |
|---|---|---|
| P0 | Broken graph, missing required file, invalid JSON, failing release validation | No |
| P1 | Missing proof surface, weak assessment, missing project rubric, broken code/test | No |
| P2 | Quality issue that weakens learning but does not break execution | Usually no |
| P3 | Editorial improvement, naming polish, optional enhancement | Yes with tracking |

## Validator design rules

- Error messages must name the exact file or ID.
- Errors must be actionable.
- Do not hide errors behind warnings when they block completion.
- Do not weaken validation to pass an incomplete repo.
- Add allowlists only with explicit reason and expiration.
- Keep generated reports machine-readable.

## Required release output

A release should include:

```text
dist/curriculum.v3.full.generated.json
dist/validation-report.json or validation-report.txt
dist/completion-report.json
dist/migration-report.json
dist/release-notes.md
```

`dist/` files are generated artifacts. Do not hand-edit.
