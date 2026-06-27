# GitHub workflow

## Learning objective

Understand and apply the GitHub workflow in the context of professional Go software engineering.

## Why this matters

Direct push to shared branches causes chaos. The GitHub workflow formalizes code review, automated testing, and controlled integration — making collaboration scalable from solo developers to thousands of contributors.

## Mental model

GitHub extends Git with a collaborative layer: remotes host branches for sharing, pull requests are proposals for merging branches with discussion attached, and actions automate verification. The workflow is propose-verify-merge.

## Core idea

The GitHub workflow is a standardized collaboration pattern: fork or clone a repository, create a feature branch, push changes, open a pull request, discuss and review, then merge. This process ensures every change is reviewed, tested, and traceable.

## Under the hood

A pull request is essentially a request to merge two Git refs. GitHub stores metadata (comments, reviews, status checks) in its database, linked to the PR number. When a merge button is clicked, GitHub performs `git merge` on the server, or for squash merge, uses `git merge --squash`. Branch protection rules are enforced server-side before the merge is allowed.

## How Go uses it

The Go project itself uses GitHub with pull requests, code review, and CI. The standard library, community packages, and most Go tools follow the same fork-branch-PR-review workflow. GitHub Actions runs Go tests on every PR for most open-source Go projects.

## Go example

The example program simulates the GitHub workflow as a state machine: a repository with forks, feature branches, pull requests with review status, and merge validation through branch protection rules.

## Step-by-step execution

1. Fork the upstream repository on GitHub to create a personal copy under your account.
2. Clone your fork locally: `git clone https://github.com/your-username/repo.git`.
3. Create a feature branch, make changes, commit logically, and push: `git push origin feature-branch`.
4. On GitHub, open a pull request from your feature branch to the upstream repository's target branch.
5. Address review feedback with additional commits, then the maintainer merges the PR.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Making changes directly on main | Not familiar with branching workflow | Always create a feature branch for changes |
| Pushing sensitive data (API keys) to a public repo | Committing secrets without realizing Git history retains them | Use `.gitignore` and `git secrets` or similar pre-commit hooks |
| Opening a PR with unrelated changes in a single commit | Not committing logically after each change | Use `git add -p` and make focused, single-purpose commits |

## Debugging walkthrough

**Scenario: PR shows "This branch has conflicts that must be resolved."**
- **Cause:** The target branch has changed since the feature branch was created.
- **Diagnosis:** Click "Resolve conflicts" on GitHub or check locally with `git merge main`.
- **Resolution:** Resolve conflicts locally: `git merge main` on the feature branch, resolve conflicts, commit, and push.

**Scenario: New commits pushed after PR approval bypass the review.**
- **Cause:** Branch protection rules are not configured to dismiss stale reviews.
- **Diagnosis:** Check the repository's branch protection settings.
- **Resolution:** Enable "Dismiss stale pull request approvals when new commits are pushed" in branch protection rules.

## Production notes

GitHub is the largest code hosting platform in the world. Millions of developers use pull requests daily for code review, CI verification, and collaborative development. The PR workflow is the industry standard for team contribution.

## Performance implications

- Pushing to GitHub is network-bound — a typical push of a few commits takes 1-5 seconds.
- Large PRs (1000+ lines changed) are slower to review and more likely to have merge conflicts.
- GitHub Actions CI adds 1-10 minutes per PR depending on test suite size.

## Practice task

Create a fork of a repository, clone it, create a feature branch, push it, and create a pull request with a descriptive title and body explaining the change.

## Tests / verification

```bash
go test ./curriculum/modules/01-computers-terminal-git-web/12-github-workflow/
```

## Review questions

1. What is the correct sequence of steps in the GitHub workflow?
2. What do the remote names `upstream` and `origin` mean?
3. What happens when a PR is merged via "Create a merge commit" vs "Squash and merge" vs "Rebase and merge"?
4. Why should you never force-push to a shared branch?
5. What is the purpose of branch protection rules?

## NEXT UP

[Lesson 13: Pull requests and code review](../13-pull-requests-and-code-review/README.md)
