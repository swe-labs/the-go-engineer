# Release Guide

This document defines the release process for The Go Engineer.

## Release Lines

| Line | Branch | Purpose |
| --- | --- | --- |
| v2.1.x | `release/v2` | stable v2.1 maintenance |
| v1.x | `release/v1` | stable v1 maintenance |
| post-v2.1 | `main` | future implementation and integration work |

The current stable release is `v2.1.1`.

## Versioning

The project uses semantic versioning for stable releases:

- major: public architecture changes or incompatible curriculum restructuring
- minor: new stable content inside the locked architecture
- patch: bug fixes, documentation corrections, validator fixes, CI changes, and low-risk curriculum repairs

Do not reuse an existing tag. If `v2.1.0` already exists, the next stable patch is `v2.1.1`.

## Release Preconditions

Before preparing a release:

- the release issue is linked
- the target branch is identified
- public docs match `ARCHITECTURE.md` and `curriculum.v2.json`
- all P0/P1/P2 review findings are fixed or explicitly accepted by the maintainer
- GitHub CI is passing
- local verification evidence is available

## Release Preparation

For a v2.1.x release:

```bash
git switch release/v2
git pull --ff-only origin release/v2
git switch -c release/v2.1.x-prep
```

Update:

- [CHANGELOG.md](./CHANGELOG.md)
- [README.md](./README.md)
- [ROADMAP.md](./ROADMAP.md)
- release notes for the target version
- any public workflow document that changed

Commit release metadata with:

```bash
git commit -m "[RELEASE] prepare v2.1.x"
```

## Required Verification

Run the CI-equivalent bundle:

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

For benchmark-related releases:

```bash
go test -bench=. -benchmem -count=1 ./08-quality-test/01-quality-and-performance/testing/benchmarks/
```

Remove generated local artifacts such as `coverage.out` unless they are intentionally tracked.

## Release PR

Open a PR into the release branch.

PR title:

```text
[RELEASE] prepare v2.1.x
```

PR body:

```markdown
Closes #<issue>

## Release

- Version: v2.1.x
- Target branch: release/v2

## Scope

-

## Validation

- [ ] go build ./...
- [ ] go vet ./...
- [ ] gofmt check
- [ ] go mod tidy no-diff check
- [ ] go test ./...
- [ ] go test -race ./...
- [ ] go test -coverprofile coverage.out ./...
- [ ] go run ./scripts/validate_curriculum.go

## Risk

-
```

Keep the PR open until CI is green and review findings are handled.

## Merge and Tag

After approval and green CI:

```bash
git switch release/v2
git pull --ff-only origin release/v2
git tag v2.1.x
git push origin v2.1.x
```

Do not force-update stable tags.

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
