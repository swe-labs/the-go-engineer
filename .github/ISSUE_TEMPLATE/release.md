---
name: Release
about: Prepare release metadata or maintenance release work
title: "[RELEASE] "
labels: release
assignees: ""
---

## Release

- Version:
- Target branch:

## Scope

-

## Metadata Updates

- [ ] `CHANGELOG.md`
- [ ] `README.md` if needed
- [ ] `ROADMAP.md` if needed
- [ ] GitHub release notes

## Validation Plan

- [ ] `go build ./...`
- [ ] `go vet ./...`
- [ ] gofmt check
- [ ] `go mod tidy` no-diff check
- [ ] `go test ./...`
- [ ] `go test -race ./...`
- [ ] `go run ./scripts/validate_curriculum.go`

## Risk Notes

-
