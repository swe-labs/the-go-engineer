# Release Guide

This document describes how to plan, prepare, and publish releases of The Go Engineer curriculum.

## Overview

The Go Engineer follows **semantic versioning** for stable releases:
- **Major**: Significant restructuring (e.g., adding entire sections)
- **Minor**: New lessons/features without breaking changes
- **Patch**: Bug fixes and documentation updates

**Release Format**: `v1.0.0`, `v1.1.0`, `v1.0.1`

## Release Process

### Step 1: Plan the Release

1. Identify issues/PRs to include in the release
2. Update [CHANGELOG.md](../CHANGELOG.md) with:
   - Date and version number
   - Major changes in sections: `Added`, `Fixed`, `Changed`, `Removed`
   - Link to detailed changes (PR numbers)

3. Update version references:
   - In documentation (README.md, ROADMAP.md)
   - In curriculum metadata if applicable

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

### Step 3: Create Release Branch

```bash
# Create release branch from main
git checkout -b release/v1.X.Y

# Commit version/changelog updates
git add CHANGELOG.md README.md ROADMAP.md
git commit -m "chore: prepare v1.X.Y release"

# Push to GitHub
git push origin release/v1.X.Y
```

### Step 4: Create Pull Request

1. Open PR from `release/v1.X.Y` to `main`
2. Title: `Release: v1.X.Y`
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

### Step 5: Merge Release

```bash
# After PR approval and all checks pass
git switch main
git pull origin main

# Merge release branch
git merge release/v1.X.Y

# Push to main
git push origin main
```

### Step 6: Create GitHub Release

1. Go to [Releases](https://github.com/rasel9t6/the-go-engineer/releases)
2. Click "Draft a new release"
3. Fill in:
   - **Tag**: `v1.X.Y`
   - **Release title**: `The Go Engineer v1.X.Y`
   - **Description**: Copy from CHANGELOG.md
   - **Attachments** (optional): Curriculum PDF, dependency graph

4. Choose "Set as latest release"
5. Publish release

### Step 7: Clean Up

```bash
# Delete release branch locally
git branch -d release/v1.X.Y

# Delete release branch on GitHub
git push origin --delete release/v1.X.Y
```

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

## Rollback Procedure

If issues are discovered after release:

```bash
# If release is not yet widely deployed
git tag -d v1.X.Y                          # Delete local tag
git push origin --delete v1.X.Y             # Delete remote tag
git push origin --delete refs/tags/v1.X.Y   # Delete tag reference

# Revert commit if already merged
git revert <commit-hash>
git push origin main

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
# If fixing critical bug in v1.0
git checkout v1.0.0
git checkout -b maintenance/v1.0

# Make fixes
git commit -m "fix: critical security issue"

# Push and create PR against v1.0 tag
git push origin maintenance/v1.0

# After merge and testing
git tag v1.0.1
git push origin v1.0.1
```

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
