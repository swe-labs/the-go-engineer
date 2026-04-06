# Maintainer Checklist

Use this checklist to keep the v1/v2 workflow consistent.

## Daily Triage

- Confirm new contributor PRs target `main` unless the work is explicitly `v1-only`.
- Add labels for `v1-only`, `v2`, `backport`, `release-blocker`, and `breaking-change` as early as possible.
- Move issues into the correct milestone: `v1 maintenance`, `v2 alpha`, `v2 beta`, `v2 rc`, or `v2.0.0`.
- Watch for PRs opened against `release/v1.0.0` and redirect them to `release/v1`.

## Merge Rules

- Use **Squash and Merge** for PRs into `main`, `release/v1`, and later `release/v2`.
- Never develop directly on long-lived branches.
- If a fix belongs in both supported lines, merge it once into the correct source branch and then `git cherry-pick -x` it to the other branch.
- Add the `backport` label before merge when that follow-up is required.

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
