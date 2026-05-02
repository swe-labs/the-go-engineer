# Contributing to The Go Engineer

This guide defines the required workflow for contributions to the stable v2.1 repository.

## Source of Truth

Read these documents before changing curriculum, code, validation, or release surfaces:

| Document | Purpose |
| --- | --- |
| [ARCHITECTURE.md](./ARCHITECTURE.md) | locked v2.1 public structure |
| [curriculum.v2.json](./curriculum.v2.json) | machine-readable curriculum registry |
| [CURRICULUM-BLUEPRINT.md](./CURRICULUM-BLUEPRINT.md) | README-first teaching contract |
| [CODE-STANDARDS.md](./CODE-STANDARDS.md) | Go and teaching-code standards |
| [TESTING-STANDARDS.md](./TESTING-STANDARDS.md) | verification expectations |
| [RELEASE.md](./RELEASE.md) | release and branch process |

If documents disagree on public curriculum structure, `ARCHITECTURE.md` wins.

## Architecture Rule

The public architecture is locked at 12 sections, `s00` through `s11`.

Do not add, remove, rename, or reorder public root sections unless a maintainer explicitly approves architecture work.

Allowed work includes:

- lesson corrections and depth improvements
- tests and proof surfaces
- README/source/curriculum metadata alignment
- validator and CI improvements
- Opslane flagship depth inside `s11`
- public documentation cleanup

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
| --- | --- |
| `main` | active post-v2.1 implementation and integration line |
| `release/v2` | stable v2.1.x maintenance line |
| `release/v1` | stable v1 maintenance line |

Branch from the line that should receive the change.

Examples:

```bash
git switch main
git pull --ff-only origin main
git switch -c docs/public-release-readiness
```

```bash
git switch release/v2
git pull --ff-only origin release/v2
git switch -c release/v2.1.x-prep
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
[DOCS] align README with v2.1.1 stable release
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
- section, lesson title, and level matching `curriculum.v2.json`
- `WHAT YOU'LL LEARN`
- `WHY THIS MATTERS`
- exact `RUN:` command matching `curriculum.v2.json`
- readable teaching comments
- Machine Role comments for major constructs that explain role, boundary, invariant, or failure mode
- `KEY TAKEAWAY`
- `NEXT UP:` footer matching `curriculum.v2.json`

## Cross-References

When a lesson uses a concept from another point in the curriculum:

- use `[!NOTE]` for prerequisite context, backward references, and gentle forward references
- use `[!TIP]` for actionable navigation, rerun suggestions, or learner practice advice
- keep cross-references inline rather than as detached navigation blocks
- include the lesson ID and a clickable local `README.md` link when referencing a specific lesson
- avoid legacy `Forward Reference` or `Backward Reference` labels

## Validation

Run the CI-equivalent bundle before final review:

```bash
go build ./...
go vet ./...
unformatted=$(gofmt -l .); test -z "$unformatted" || (echo "$unformatted" && exit 1)
go mod tidy
git diff --exit-code -- go.mod go.sum
go test ./...
go test -race ./...
go test -coverprofile=coverage.out ./...
go run ./scripts/validate_curriculum.go
```

On PowerShell, quote the coverage flag if needed:

```powershell
go test "-coverprofile=coverage.out" ./...
```

For benchmark-related changes:

```bash
go test -bench=. -benchmem -count=1 ./08-quality-test/01-quality-and-performance/testing/benchmarks/
```

## Review Standard

Before final review:

- architecture v2.1 remains intact
- `curriculum.v2.json` matches files and lesson metadata
- README, source, tests, and section maps agree
- `NEXT UP:` footers match curriculum metadata
- Go code is formatted and idiomatic
- tests prove behavior where behavior should be provable
- CI-equivalent validation passes locally
- P0/P1/P2 findings are fixed or explicitly documented

Final squash merge is maintainer-controlled.
