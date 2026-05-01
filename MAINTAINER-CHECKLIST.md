# Maintainer Checklist

Use this checklist to keep repository workflow, branch maintenance, documentation, and releases consistent.

All section references follow the locked v2.1 12-section architecture from `ARCHITECTURE.md`.

## Daily Triage

- Enforce issue-first workflow.
- Confirm new contributor PRs target the correct branch.
- Add labels early.
- Move issues into the correct milestone and the GitHub Project: `The Go Engineer v2`.
- Redirect stable-line work to the correct release branch.
- Reject architecture changes unless they are explicitly approved.

## Project and milestone tracking

For every issue and PR:

- assign the owner when appropriate
- apply relevant labels
- attach the active milestone when one exists
- add the issue or PR to `The Go Engineer v2`
- update tracking metadata when scope changes

## Required Issue Title Format

```text
[TYPE] short description
```

Allowed types:

```text
[FEAT] [FIX] [DOCS] [TEST] [CHORE] [REFACTOR] [SECURITY] [RELEASE]
```

## Closed-Loop PR Maintenance

Before a PR is ready for final maintainer review:

- findings-first self-review must be complete
- P0/P1/P2 findings must be fixed or explicitly documented
- CI failures must be fixed or explained with evidence
- top-level PR comments must be answered
- inline review threads must be fixed, explained, or marked outdated when the diff no longer applies
- addressed threads should be resolved when permitted
- a final readiness comment should summarize changes, handled findings, validation, and remaining risks

The final squash-merge remains maintainer-only.

## PR Review Rules

- Reject PRs that bypass the README-first teaching contract.
- Reject PRs that introduce "magic" early in the curriculum.
- Reject public architecture drift unless explicitly approved.
- Require linked issues.
- Require validation evidence.
- Use **Squash and Merge**.
- Never develop directly on long-lived branches.
- Do not self-merge unless explicitly approved for that PR.

## Merge Commit Format

Final squash commits should use:

```text
[TYPE] short imperative description
```

Examples:

```text
[FEAT] add SY.1 sync.Mutex lesson
[FIX] correct DB.3 tenant-scoping example
[DOCS] align release docs with v2.1.x maintenance
```

## Branch Roles

- `main`: post-v2.1 implementation line
- `release/v1`: stable v1 maintenance line
- `release/v2`: stable v2.1.x maintenance line

## Backports

- For a stable v1 bug, fix on `release/v1` first.
- For a v2.1.x stable bug, fix on `release/v2` first.
- For post-v2.1 implementation, work on `main`.
- If a fix belongs in multiple supported lines, use `git cherry-pick -x`.

```bash
git switch <target-branch>
git pull origin <target-branch>
git cherry-pick -x <merged-commit-sha>
git push origin <target-branch>
```

## Release Flow

For v2.1.x maintenance releases:

1. Branch from `release/v2`.
2. Update `CHANGELOG.md`, release notes, and any required metadata.
3. Run local verification.
4. Open a release PR into `release/v2`.
5. Squash and merge after green CI and approval.
6. Tag from `release/v2`.

## Local Verification

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

## Documentation Alignment Check

Before release or workflow changes, verify these documents agree with `ARCHITECTURE.md` and each other:

- `README.md`
- `ARCHITECTURE.md`
- `ROADMAP.md`
- `LEARNING-PATH.md`
- `CURRICULUM-BLUEPRINT.md`
- `CODE-STANDARDS.md`
- `TESTING-STANDARDS.md`
- `COMMON-MISTAKES.md`
- `docs/PROGRESSION.md`
- `docs/ENGINEERING_ERROR_FRAMEWORK.md`
- `docs/flagship/OPSLANE_SAAS_BACKEND.md`
- `CONTRIBUTING.md`
- `RELEASE.md`
- `MAINTAINER-CHECKLIST.md`
- `AGENTS.md`
- `.github/pull_request_template.md`
- `.github/ISSUE_TEMPLATE/*`
- `.agents/skills/*/SKILL.md`

## Architecture Checklist

- [ ] Public section count remains 12.
- [ ] Section IDs remain `s00` through `s11`.
- [ ] Canonical folder names remain unchanged.
- [ ] Lesson IDs belong to the correct section.
- [ ] `curriculum.v2.json` matches filesystem paths.
- [ ] Validator passes.

## Lesson Checklist

- [ ] README follows the required section order.
- [ ] `main.go` follows the standard header and footer.
- [ ] `NEXT UP:` matches the next item ID and path in `curriculum.v2.json`.
- [ ] Section README is updated.
- [ ] Starter code compiles.
- [ ] Tests prove behavior.
- [ ] Run command works exactly as printed.
