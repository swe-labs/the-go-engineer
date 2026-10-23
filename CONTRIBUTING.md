# Contributing to The Go Engineer

This guide defines the required workflow for contributions to **The Go Engineer v3**.

## Source of Truth

Read these documents before changing curriculum, code, validation, or release surfaces:

| Document | Purpose |
|---|---|
| [ARCHITECTURE.md](./ARCHITECTURE.md) | Canonical v3 public structure |
| [`curriculum/`](./curriculum/README.md) | Modular curriculum metadata manifests |
| [CURRICULUM-BLUEPRINT.md](./CURRICULUM-BLUEPRINT.md) | Zero-magic teaching contract |
| [CODE-STANDARDS.md](./CODE-STANDARDS.md) | Go and teaching-code standards |
| [TESTING-STANDARDS.md](./TESTING-STANDARDS.md) | Verification and test expectations |
| [RELEASE.md](./RELEASE.md) | Release and branch workflow |

If documents disagree on public curriculum structure, `ARCHITECTURE.md` wins.

## Architecture Rule

The canonical v3 architecture consists of 18 core modules (`00`–`16` plus `18-flagship-opslane`) under `curriculum/modules/`, accompanied by the Companion Full-Stack Frontend Track and 20 Advanced Electives.

Allowed work includes:

- lesson corrections, depth improvements, and zero-magic explanations
- project tests, starters, and reference solutions
- curriculum metadata alignment across `curriculum/*.json`
- quality tooling improvements in `internal/tools/curriculum`
- Opslane milestone enhancements inside `curriculum/modules/18-flagship-opslane/`
- public documentation improvements

## Issue First

Every non-trivial change needs an issue before implementation.

Issue title format:

```text
[TYPE] short description
```

Allowed types:

```text
[FEAT] [FIX] [DOCS] [TEST] [CHORE] [REFACTOR] [SECURITY] [RELEASE]
```

Issue body should include:

- affected section, if applicable
- affected lesson IDs, if applicable
- expected files or folders
- validation plan
- risk notes

## Branches

| Branch | Purpose |
|---|---|
| `main` | active development and next-generation feature work |
| `release/v3` | active v3 release preparation and maintenance line |

Branch from the line that should receive the change.

Examples:

```bash
git switch main
git pull --ff-only origin main
git switch -c feat/new-feature
```

```bash
git switch release/v3
git pull --ff-only origin release/v3
git switch -c fix/patch-fix
```

## Pull Requests

Open a draft PR early and link the issue with `Closes #<issue>`.

The PR should include:

- scope
- affected sections or lessons
- validation commands run
- risk notes
- tracking metadata

Keep the PR in draft until local validation, self-review, and obvious fixes are complete.

## Commits

Use bracketed commit messages:

```text
[TYPE] short imperative description
```

Examples:

```text
[DOCS] align README with v3.0.0 architecture
[FIX] replace blocking REST API test with httptest
[TEST] add repository behavior coverage
[CHORE] tighten curriculum validator
```

Keep commits logical. Do not mix unrelated formatting, docs, and behavior changes unless the PR is explicitly a repository-wide cleanup.

## Lesson Contract

Every learner-facing lesson README must include these sections in order:

1. `## Mission`
2. `## Prerequisites`
3. `## Mental Model`
4. `## Visual Model`
5. `## Machine View`
6. `## Run Instructions`
7. `## Code Walkthrough`
8. `## Try It`
9. `## In Production`
10. `## Thinking Questions`
11. `## Next Step`

Exercises replace `## Code Walkthrough` with `## Solution Walkthrough` and include `## Verification Surface`.

Every lesson `main.go` should include:

- copyright and license header
- module, lesson title, and level
- `WHAT YOU'LL LEARN`
- `WHY THIS MATTERS`
- exact `RUN:` command
- readable teaching comments
- `KEY TAKEAWAY`
- `NEXT UP:` footer matching canonical progression in `curriculum/path.core.json`

## Cross-References

When a lesson uses a concept from another point in the curriculum:

- forward reference later concepts briefly and name the future lesson or module
- backward reference previously established ideas when they are reused
- keep cross-references inline with the explanation rather than as detached navigation blocks

## Validation

Run the CI-equivalent bundle before final review:

```bash
go build ./...
go vet ./...
unformatted=$(gofmt -l . | grep -v '_starter' || true); test -z "$unformatted" || (echo "$unformatted" && exit 1)
go mod tidy
git diff --exit-code -- go.mod go.sum
go test ./...
go test -race ./curriculum/modules/11-lifecycle-context-concurrency/... ./curriculum/modules/18-flagship-opslane/...
go run ./internal/tools/curriculum validate-all
```

For benchmark-related changes:

```bash
go test -bench=. -benchmem -count=1 ./curriculum/modules/13-performance-memory-engineering/projects/profile-and-fix-lab/_solution
```

## Review Standard

Before final review:

- canonical v3 18-module architecture remains intact
- manifests in `curriculum/` match files and lesson metadata
- README, source, tests, and module maps agree
- `NEXT UP:` footers match curriculum metadata in `curriculum/path.core.json`
- Go code is formatted and idiomatic
- tests prove behavior where behavior should be provable
- all 17 automated curriculum quality gates pass (`go run ./internal/tools/curriculum validate-all`)
- P0/P1/P2 findings are fixed or explicitly documented

Final squash merge is maintainer-controlled.
