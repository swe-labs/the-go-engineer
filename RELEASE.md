# Release Guide

This document defines the release process for The Go Engineer.

## Release Lines

| Line | Branch | Purpose |
|---|---|---|
| v3.x | `release/v3` | active v3 release and maintenance line |
| v2.1.x | `release/v2` | legacy v2.1 maintenance line (archived) |
| v1.x | `release/v1` | legacy v1 maintenance line (archived) |

## Versioning

The project uses semantic versioning for stable releases:

- major: public architecture changes or incompatible curriculum restructuring
- minor: new stable content inside the locked architecture
- patch: bug fixes, documentation corrections, validator fixes, CI changes, and low-risk curriculum repairs

## Release Preconditions

Before preparing a release:

- the release issue is linked
- the target branch is identified (`release/v3`)
- public docs match [ARCHITECTURE.md](./ARCHITECTURE.md) and [`curriculum/`](./curriculum/README.md) manifests
- all P0/P1/P2 review findings are fixed
- GitHub CI is passing
- local verification evidence is verified

## Required Verification

Run the CI-equivalent bundle:

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

For benchmark-related releases:

```bash
go test -bench=. -benchmem -count=1 ./curriculum/modules/13-performance-memory-engineering/projects/profile-and-fix-lab/_solution
```

## GitHub Release

Create the GitHub release from the stable tag:

```bash
gh release create v2.1.x --target release/v2 --title "The Go Engineer v2.1.x" --notes-file RELEASE-NOTES-v2.1.x.md
```

The release notes should include:

- release purpose
- major fixes and documentation changes
- validation evidence
- branch and tag references
- known risks, if any

## Backports

If a fix must ship to more than one supported line:

1. merge the fix into the source branch
2. cherry-pick with provenance
3. open a follow-up PR if branch protection requires it

```bash
git switch <target-branch>
git pull --ff-only origin <target-branch>
git cherry-pick -x <merged-commit-sha>
git push origin HEAD
```

## Security Releases

For security fixes:

1. use `[SECURITY]` in the issue, PR, and commit title
2. make the smallest safe fix
3. add tests where possible
4. avoid exposing secrets or exploit details in public discussion
5. publish a patch release
6. document user-facing impact in the changelog

## Rollback

If a tag was pushed incorrectly before release publication:

```bash
git tag -d v2.1.x
git push origin --delete v2.1.x
```

If the release was already published, prefer a revert followed by a patch release.

## Checklist

- [ ] Release issue is linked.
- [ ] Target branch is correct.
- [ ] Public docs are current.
- [ ] `CHANGELOG.md` is updated.
- [ ] Release notes are ready.
- [ ] Local verification passes.
- [ ] GitHub CI passes.
- [ ] No architecture drift.
- [ ] Tag is created from the verified release branch.
- [ ] GitHub release is published.
