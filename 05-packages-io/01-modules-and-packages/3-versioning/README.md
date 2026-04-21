# MP.3 Versioning Workshop

## Mission

Build a small semantic-versioning and module-policy demo that makes Go's version rules concrete.

This exercise is the Stage 05 milestone.
It is where module identity, compatibility boundaries, and `replace` reasoning come together in
one runnable artifact with tests.

## Prerequisites

Complete these first:

- `MP.1` module basics
- `MP.2` managing dependencies

## What You Will Build

Implement a small versioning helper that:

1. models semantic versions with a `Version` struct
2. formats versions consistently as `vMAJOR.MINOR.PATCH`
3. detects compatibility based on the major version boundary
4. compares versions to decide which one is newer
5. prints a clear explanation of the `/v2` import-path rule
6. passes the provided unit tests

## Files

- [main.go](./main.go): complete solution with teaching comments
- [version_test.go](./version_test.go): tests for the version helper behavior
- [_starter/main.go](./_starter/main.go): starter file with TODOs and requirements

## Run Instructions

Run the completed solution:

```bash
go run ./05-packages-io/01-modules-and-packages/3-versioning
```

Run the tests:

```bash
go test ./05-packages-io/01-modules-and-packages/3-versioning
```

Run the starter:

```bash
go run ./05-packages-io/01-modules-and-packages/3-versioning/_starter
```

## Success Criteria

Your finished solution should:

- return versions in the `vX.Y.Z` format
- treat matching major versions as compatible
- compare major, then minor, then patch when deciding which version is newer
- explain why v2+ modules require a new import path
- explain when `replace` helps and when it should be used carefully
- pass the provided tests

## Common Failure Modes

- treating minor or patch bumps as breaking changes
- comparing versions lexicographically instead of numerically
- forgetting that a v2 module needs `/v2` in the import path
- using `replace` as a permanent escape hatch instead of a deliberate local-development tool

## Next Step

After you complete this exercise, continue to [Stage 05](../../02-io-and-cli) if you are ready
to move on.
If you want one more stretch lesson first, visit [`MP.4` build tags](../4-build-tags).
