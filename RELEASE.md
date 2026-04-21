# Release Guide

This document describes how to plan, prepare, and publish releases of The Go Engineer curriculum.

## Overview

The Go Engineer follows **semantic versioning** for stable releases:
- **Major**: Significant restructuring (e.g., adding entire sections)
- **Minor**: New lessons/features without breaking changes
- **Patch**: Bug fixes and documentation updates

**Release Format**: `v1.0.0`, `v1.1.0`, `v1.0.1`

## Branch Roles

The project uses long-lived branches for supported major versions:

- `main`: active v2 development and prerelease integration branch
- `release/v1`: stable v1 maintenance branch for current users
- `release/v2`: created from `main` when v2 reaches feature freeze

Topic branches stay short-lived and always branch from the line they should ship to.

- `feat/...` and `fix/...` branch from `main` for v2 work
- `fix/v1-...` or `hotfix/v1-...` branch from `release/v1` for stable v1 fixes

`v1.0.0` remains an immutable release tag. It is not the permanent maintenance branch name.

## Release Process

### Step 1: Plan the Release

1. Identify issues/PRs to include in the release
2. Update [CHANGELOG.md](./CHANGELOG.md) with:
   - Date and version number
   - Major changes in sections: `Added`, `Fixed`, `Changed`, `Removed`
   - Link to detailed changes (PR numbers)

3. Update version references:
   - In documentation (README.md, ROADMAP.md)

### Step 2: Pre-Release Verification

```bash
# Verify build
make build

# Run all tests
make test

# Run with race detection
make test-race

# Generate coverage
make cover

# Validate curriculum
go run ./scripts/validate_curriculum.go

# Check for uncommitted changes
git status

# Check for dependencies that need updating
make deps-check
```

### Step 3: Choose the Correct Release Line

- For v1 patch or minor releases, work from `release/v1`
- For ongoing v2 development and alpha prereleases, work from `main`
- For v2 beta, release candidate, and final stabilization, cut `release/v2` from `main`

Do not open sync PRs from `main` into `release/v1` just to keep the branches identical. Once v2 begins, those branches are expected to diverge.

### Step 4: Create the Release Branch or Topic Branch

```bash
# One-time step: cut the long-lived v2 stabilization branch
git switch main
git pull origin main
git switch -c release/v2
git push -u origin release/v2

# Example: prepare a v1 patch release from the stable line
git switch release/v1
git pull origin release/v1
git switch -c release-prep/v1.X.Y

# Example: prepare a v2 stabilization update after release/v2 exists
git switch release/v2
git pull origin release/v2
git switch -c release-prep/v2.0.0-rc.N

# Commit version/changelog updates
git add CHANGELOG.md README.md ROADMAP.md
git commit -m "chore: prepare release metadata"

# Push to GitHub
git push origin HEAD
```

### Step 5: Create Pull Request

1. Open PR into the long-lived target branch:
   - `release-prep/v1.X.Y` -> `release/v1`
   - `release-prep/v2.0.0-rc.N` -> `release/v2`
2. Title: `Release: v1.X.Y` or `Release: v2.0.0-rc.N`
3. Description:
   ```markdown
   ## Release v1.X.Y
   
   **Date**: YYYY-MM-DD
   
   ### Summary
   [Brief description of changes]
   
   ### New Lessons
   - Lesson ID: Description
   
   ### Bug Fixes
   - Description
   
   ### Documentation
   - Description
   
   ### Breaking Changes
   - Description (if any)
   ```

4. Wait for all CI checks to pass
5. Get approval from maintainers

### Step 6: Merge Release

```bash
# After the PR is approved, merge it in GitHub using Squash and Merge.
# Then update your local long-lived branch.
git switch release/v1
git pull origin release/v1
```

Maintainers should use **Squash and Merge** for release pull requests. If the same fix also belongs in another supported branch, propagate it with `git cherry-pick -x` instead of a branch sync merge.

### Step 7: Create GitHub Release

1. Go to [Releases](https://github.com/rasel9t6/the-go-engineer/releases)
2. Click "Draft a new release"
3. Fill in:
   - **Tag**: `v1.X.Y`
   - **Release title**: `The Go Engineer v1.X.Y`
   - **Description**: Copy from CHANGELOG.md
   - **Attachments** (optional): Curriculum PDF

4. Choose "Set as latest release"
5. Publish release

Use prerelease tags during the v2 rollout:

- `v2.0.0-alpha.N` from `main`
- `v2.0.0-beta.N` from `release/v2`
- `v2.0.0-rc.N` from `release/v2`

Mark alpha, beta, and RC builds as prereleases on GitHub so stable v1 users are not silently moved early.

### Step 8: Clean Up

```bash
# Delete the short-lived topic branch locally
git branch -d release-prep/v1.X.Y

# Delete the short-lived topic branch on GitHub
git push origin --delete release-prep/v1.X.Y
```

Keep `release/v1` and later `release/v2` as permanent branches while those lines are supported.

## Release Checklist

Before releasing, verify:

- [ ] All issues in the milestone are closed
- [ ] All PRs in the release are merged
- [ ] CHANGELOG.md is updated with all changes
- [ ] ROADMAP.md status indicators are accurate
- [ ] All CI checks pass (`build`, `test`, `vet`, `fmt`)
- [ ] Test coverage is maintained (> 75%)
- [ ] No security vulnerabilities in dependencies
- [ ] Documentation is up to date
- [ ] No breaking changes without major version bump
- [ ] Version numbers updated in appropriate files
- [ ] PR target matches the release line (`release/v1` or `release/v2`)
- [ ] Any cross-line fix has a planned `cherry-pick -x` follow-up

## Rollback Procedure

If issues are discovered after release:

```bash
# If release is not yet widely deployed
git tag -d v1.X.Y                          # Delete local tag
git push origin --delete v1.X.Y             # Delete remote tag
git push origin --delete refs/tags/v1.X.Y   # Delete tag reference

# Revert commit if already merged
git revert <commit-hash>
git push origin <release-branch>

# Create patch release when ready
# e.g., v1.X.1 for security fix on v1.X.0
```

## Version Numbering Guidelines

### When to bump MAJOR (v2.0.0)

- Significant curriculum restructuring
- Changing core learning path
- Removing entire sections (breaking change for learners)

### When to bump MINOR (v1.1.0)

- Adding new sections
- Adding new lessons
- Adding exercises
- Significant documentation improvements

### When to bump PATCH (v1.0.1)

- Bug fixes (typos, incorrect code)
- Documentation corrections
- CI/build improvements
- Dependency updates

## Maintenance Releases

For maintaining older releases (security fixes, critical bugs):

```bash
# If fixing a critical bug in the stable v1 line
git switch release/v1
git pull origin release/v1
git switch -c hotfix/v1.X.Z

# Make fixes
git commit -m "fix: critical security issue"

# Push and create PR against release/v1
git push origin hotfix/v1.X.Z

# After merge and testing
git tag v1.X.(Z+1)
git push origin v1.X.(Z+1)
```

If the same fix is also needed on `main`, cherry-pick it forward with `git cherry-pick -x`.

## Dependency Updates

### Regular Updates

```bash
# Check for newer versions
make deps-check

# Update to latest patch versions
make deps-update

# Test thoroughly
make test-race
make cover

# Commit if all tests pass
git commit -m "chore: update dependencies"
```

### Security Updates

If a security vulnerability is discovered:

1. **Immediate patch release**:
   - Apply the security fix
   - Create v1.X.(Z+1) release
   - Document in CHANGELOG with `[SECURITY]` prefix
   - Notify users/contributors

2. **Example**:
   ```markdown
   ## [1.0.1] - 2024-04-15 [SECURITY]
   
   ### Security
   - Fixed SQL injection vulnerability in database lesson (CVE-2024-XXXXX)
   - Updated all dependencies to patched versions
   
   ### Recommendation
   Users are urged to update immediately.
   ```

## Release Announcements

For major releases, announce via:

1. **GitHub Releases page** (automatic)
2. **README.md** (update featured changes)
3. **CHANGELOG.md** (already done)
4. **Contributors** - Thank those who contributed
5. **Social media** (optional): Share on Twitter/X, LinkedIn, etc.

### Example Announcement

```
🎉 The Go Engineer v1.2.0 is released!

New in this release:
✨ 8 new lessons on Concurrency Patterns
✨ Enhanced Docker & deployment section
🐛 Fixed curriculum mapping issues
📚 Comprehensive testing standards guide

📖 Full changelog: [link]
🚀 Get started: [link]

Thank you to all contributors! @mention @mention @mention
```

## Post-Release

### Monitoring

- Monitor issues for any problems reported with new release
- Check CI status of dependent projects
- Respond quickly to security reports

### Next Steps

1. **Create milestone** for next version
2. **Triage** all open issues and PRs
3. **Plan** next release content based on community feedback
4. **Update** contributor discussions/roadmap as needed

## GitHub Repository Settings

Recommended repository settings for this workflow:

- protect `main`, `release/v1`, and later `release/v2`
- require pull requests, status checks, and at least one maintainer review
- disable direct pushes to protected branches
- prefer **Squash and Merge** for pull requests
- enable automatic deletion of head branches after merge

Ahead/behind counts between `main` and `release/v1` are normal once v2 work starts. The goal is not identical histories across major versions; the goal is intentional propagation of the fixes that should exist in both lines.

## Release Schedule

The Go Engineer follows an irregular release schedule based on:
- Completion of curriculum sections
- Community feedback and requests
- Bug/security fixes requiring immediate release

**No minimum interval** between releases, but typically:
- Major releases: Once per year (after 2-3 new sections)
- Minor releases: 2-4 times per year
- Patch releases: As needed (1-2 weeks when issues discovered)

## Tools & Resources

- **GitHub Release Notes**: https://github.com/rasel9t6/the-go-engineer/releases
- **Changelog Format**: https://keepachangelog.com/
- **Semantic Versioning**: https://semver.org/
- **Commit Conventions**: Conventional Commits (feat:, fix:, docs:, etc.)
