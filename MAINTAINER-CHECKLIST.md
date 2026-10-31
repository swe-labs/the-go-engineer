# Maintainer Checklist

Use this checklist to keep repository workflow, branch maintenance, documentation, and releases consistent.

All module references follow the canonical v3 architecture from [ARCHITECTURE.md](./ARCHITECTURE.md).

## Daily Triage

- Enforce issue-first workflow.
- Confirm new contributor PRs target the correct branch (`release/v3` or `main`).
- Add labels early.
- Move issues into the correct milestone and the GitHub Project: `The Go Engineer v3`.
- Redirect release-line work to the correct release branch.
- Reject unapproved architectural drift.

## Project and milestone tracking

For every issue and PR:

- assign the owner when appropriate
- apply relevant labels
- attach the active milestone when one exists
- add the issue or PR to `The Go Engineer v3`
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
[FEAT] add bounded worker pool lesson
[FIX] correct tenant-scoping example in Opslane
[DOCS] align release docs with v3.0.0 architecture
```

## Branch Roles

- `main`: active development and next-generation feature work
- `release/v3`: active v3 release preparation and maintenance line

## Backports

- For a v3 bug, fix on `release/v3` first and cherry-pick to `main`.
- If a fix belongs in multiple supported lines, use `git cherry-pick -x`.

```bash
git switch <target-branch>
git pull origin <target-branch>
git cherry-pick -x <merged-commit-sha>
git push origin <target-branch>
```

## Release Flow

For v3.x maintenance releases:

1. Branch from `release/v3`.
2. Update `CHANGELOG.md` and any required metadata.
3. Run local verification.
4. Open a release PR into `release/v3`.
5. Squash and merge after green CI and approval.

## Local Verification

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

## Documentation Alignment Check

Before release or workflow changes, verify these documents agree with `ARCHITECTURE.md` and each other:

- `README.md`
- `ARCHITECTURE.md`
- `ROADMAP.md`
- `LEARNING-PATH.md`
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

- [ ] Core module count remains 18 (`00` through `16`, plus `18-flagship-opslane`).
- [ ] Module directories remain canonically organized under `curriculum/modules/`.
- [ ] Item IDs belong to the declared module progression.
- [ ] Metadata manifests in `curriculum/` match filesystem paths.
- [ ] All 17 automated curriculum quality gates pass.

## Lesson Checklist

- [ ] README follows the required pedagogical section order.
- [ ] `main.go` follows the standard header and key takeaway footer.
- [ ] `NEXT UP:` matches `curriculum/path.core.json`.
- [ ] Module README is updated.
- [ ] Starter code compiles and establishes red verification state.
- [ ] Tests prove behavior.
- [ ] Run command works exactly as printed.
