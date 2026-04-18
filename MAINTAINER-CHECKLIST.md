# Maintainer Checklist

> Use this checklist to keep the v1/v2 workflow consistent.
> All section references follow the v2.1 12-section architecture from `ARCHITECTURE.md`.

## Daily Triage

- Enforce the **GitHub Workflow**: No PRs should be reviewed unless they link to an approved Issue.
- Confirm new contributor PRs target `main` unless the work is explicitly `v1-only`.
- Add labels for `v1-only`, `v2`, `backport`, `release-blocker`, and `breaking-change` as early as possible.
- Move issues into the correct milestone and add them to the **"The Go Engineer v2"** project board.
- Watch for PRs opened against `release/v1.0.0` and redirect them to `release/v1`.

## PR Review & Merge Rules

- **Reject** PRs that try to bypass the `README`-first teaching contract or introduce "magic" early in the curriculum (see `CURRICULUM-BLUEPRINT.md`).
- Use **Squash and Merge** for PRs into `main`, `release/v1`, and later `release/v2`.
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

- Tag v2 prereleases from `main` as `v2.0.0-alpha.N`.
- Cut `release/v2` from `main` once v2 is feature complete.
- Tag beta and RC builds from `release/v2`.
- Tag final `v2.0.0` from `release/v2`.
- Keep `release/v1` for v1 patch support until you formally end support.

## Branch Hygiene

- Keep `main` as the default branch.
- Keep branch protections on `main` and `release/v1`.
- Retire `release/v1.0.0` only after all external references and protections are moved to `release/v1`.
- Auto-delete short-lived branches after merge.

## Doc Alignment Check

Before any release, verify these documents are aligned with `ARCHITECTURE.md`:

- `ROADMAP.md` — section statuses match reality
- `LEARNING-PATH.md` — phases and section boundaries correct
- `CURRICULUM-BLUEPRINT.md` — teaching contract matches README contract
- `CODE-STANDARDS.md` — NEXT UP regex and templates current
- `TESTING-STANDARDS.md` — coverage targets match section IDs
- `COMMON-MISTAKES.md` — all "Taught in" references use correct lesson IDs
- `docs/PROGRESSION.md` — milestone table matches `ARCHITECTURE.md` milestones
- `CONTRIBUTING.md` — section numbering and workflow current
