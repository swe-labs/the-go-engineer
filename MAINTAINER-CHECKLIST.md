# Maintainer Checklist

> Use this checklist to keep the v1/v2 workflow consistent.
> All section references follow the v2.1 12-section architecture from `ARCHITECTURE.md`.

## Daily Triage

- Enforce the **GitHub Workflow**: no PRs should be reviewed unless they link to an approved issue.
- Confirm new contributor PRs target `main` unless the work is explicitly `v1-only`.
- Add labels for `v1-only`, `v2`, `backport`, `release-blocker`, and `breaking-change` as early as possible.
- Move issues into the correct milestone and add them to the **"The Go Engineer v2"** project board.
- Watch for PRs opened against `release/v1.0.0` and redirect them to `release/v1`.

## PR Review & Merge Rules

- **Reject** PRs that try to bypass the `README`-first teaching contract or introduce "magic" early in the curriculum (see `CURRICULUM-BLUEPRINT.md`).
- Use **Squash and Merge** for PRs into `main`, `release/v1`, and `release/v2`.
- Never develop directly on long-lived branches.
- If a fix belongs in both supported lines, merge it once into the correct source branch and then `git cherry-pick -x` it to the other branch.
- Add the `backport` label before merge when that follow-up is required.
- Verify that new lessons are registered in `curriculum.v2.json` and pass `go run ./scripts/validate_curriculum.go`.

## Backports

- For a stable-user bug, fix it on `release/v1` first.
- For a v2-only bug, fix it on `main`.
- After merge, cherry-pick with:

```bash
git switch <target-branch>
git pull origin <target-branch>
git cherry-pick -x <merged-commit-sha>
git push origin <target-branch>
```

- Open a follow-up PR if the destination branch is protected.
- Remove the `backport` label only after the second branch has the fix.

## Release Flow

- Tag v2 prereleases from `main` as `v2.1.0-alpha.N` while the stabilization line is still open.
- Keep `release/v2` as the active RC stabilization branch once beta-complete work is cut there.
- Tag beta and RC builds from `release/v2`.
- Tag final `v2.1.0` from `release/v2`.
- Keep `release/v1` for v1 patch support until you formally end support.

### RC.1 Gate

Before tagging `v2.1.0-rc.1`, verify all of the following on `release/v2`:

- `make build`
- `make test`
- `make test-race`
- `make bench`
- `make run-hello`
- `make run-env`
- `go run ./scripts/validate_curriculum.go`
- the release-prep PR is green and approved
- all `release-blocker` issues in the `v2 rc` milestone are closed or explicitly accepted for deferment
- the release branch is clean immediately before tagging

If GNU Make is not available in the maintainer environment, run the documented direct Go-command equivalents from `RELEASE.md` instead of skipping the gate.

If `go test -race ./...` cannot run locally because the environment lacks a supported CGO toolchain or C compiler, do not silently waive it. Use the CI race check on the release-prep PR as the release gate for that item.

### After `v2.1.0-rc.1` Is Published

- Keep `release/v2` focused on release blockers, validation findings, and release-facing polish only.
- Open all RC findings in the `v2 rc` milestone and keep them on the **The Go Engineer v2** project board.
- Do not resume beta migration or broad architecture work on `release/v2`.
- Tag final `v2.1.0` only after RC blockers are closed or explicitly deferred.

### Stable `v2.1.0` Gate

Before tagging `v2.1.0`, verify all of the following on `release/v2`:

- `make build`
- `make test`
- `make test-race`
- `make bench`
- `make run-hello`
- `make run-env`
- `go run ./scripts/validate_curriculum.go`
- the `release-prep/v2.1.0` PR is green and approved
- `v2.1.0-rc.1` feedback has been reviewed and true blockers are fixed or explicitly deferred
- all remaining `release-blocker` issues in the `v2 rc` milestone are closed or intentionally deferred
- the release branch is clean immediately before tagging

If GNU Make is not available in the maintainer environment, run the documented direct Go-command equivalents from `RELEASE.md`.

If local `go test -race ./...` is still blocked by the maintainer environment, require the CI race check on the stable release-prep PR before tagging.

## Branch Hygiene

- Keep `main` as the default branch.
- Keep branch protections on `main`, `release/v1`, and `release/v2`.
- Retire `release/v1.0.0` only after all external references and protections are moved to `release/v1`.
- Auto-delete short-lived branches after merge.

## Doc Alignment Check

Before any release, verify these documents are aligned with `ARCHITECTURE.md`:

- `ROADMAP.md` - section statuses match reality
- `LEARNING-PATH.md` - phases and section boundaries correct
- `CURRICULUM-BLUEPRINT.md` - teaching contract matches README contract
- `CODE-STANDARDS.md` - NEXT UP regex and templates current
- `TESTING-STANDARDS.md` - coverage targets match section IDs
- `COMMON-MISTAKES.md` - all "Taught in" references use correct lesson IDs
- `docs/PROGRESSION.md` - milestone table matches `ARCHITECTURE.md` milestones
- `CONTRIBUTING.md` - section numbering and workflow current
