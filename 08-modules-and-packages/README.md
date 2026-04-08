# Section 08: Modules and Packages

## Mission

This section teaches you how Go modules define package boundaries, dependency resolution, versioned
imports, and controlled build surfaces.

By the end of Section 08, you should be comfortable reading and writing:

- `go.mod` and `go.sum` with confidence instead of treating them as magic files
- dependency-management workflows using `go get`, `go mod tidy`, `go list`, and `go mod why`
- semantic-versioning decisions, especially the `/v2` import-path rule
- local override workflows with `replace`
- optional platform and test surfaces through build tags

## Who Should Start Here

### Full Path

Start here after completing Section 07 in order.

### Bridge Path

You can move faster if you already understand:

- how Go import paths map to package directories
- basic command-line use of the `go` tool
- semantic-versioning vocabulary like major, minor, and patch

Even on the bridge path, do not skip `MP.1`.
It removes a lot of mystery from the rest of the section.

### Targeted Path

If you are here mainly for dependency hygiene and versioned module behavior, review these first:

- `MP.1` module basics
- `MP.2` managing dependencies

## Section Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `MP.1` | Lesson | [module basics](./1-module-basics) | Explains what `go.mod` and `go.sum` actually do so dependency work stops feeling magical. | entry |
| `MP.2` | Lesson | [managing deps](./2-managing-deps) | Shows how to add, inspect, trim, and explain dependencies using the Go toolchain. | `MP.1` |
| `MP.3` | Exercise | [versioning workshop](./3-versioning) | Combines semantic versioning, `/v2` imports, and `replace` reasoning into one milestone. | `MP.1`, `MP.2` |
| `MP.4` | Stretch Lesson | [build tags](./4-build-tags) | Adds conditional compilation and tagged tests once the core module workflow is solid. | `MP.1`, `MP.2` |

## Suggested Order

1. Work through `MP.1` and `MP.2` in order.
2. Complete `MP.3` without copying the finished solution line by line.
3. Use `MP.4` as an optional stretch lesson before moving to the next section.

## Section Milestone

`MP.3` is the current live milestone for this pilot section.

If you can complete it and explain:

- why `go.mod` and `go.sum` are part of the build contract, not random generated clutter
- why a module that reaches v2 must change its import path
- why `replace` is useful in local development but should be used deliberately

then you are ready to move into file I/O and CLI tooling in Section 09.

## Pilot Role In V2

This live v2 slice keeps the current Section 08 paths and `MP.*` ids stable while upgrading the
learner-facing structure:

- `MP.1` and `MP.2` are the core lessons
- `MP.3` is the milestone exercise
- `MP.4` is an optional stretch lesson

That keeps the section useful for current learners while the broader v2 migration continues.

## Legacy To Pilot Mapping

This pilot intentionally avoids breaking the current Section 08 filesystem layout.

- `MP.1` stays at `08-modules-and-packages/1-module-basics`
- `MP.2` stays at `08-modules-and-packages/2-managing-deps`
- `MP.3` stays at `08-modules-and-packages/3-versioning`, now treated as the section milestone
  surface with a starter and tests
- `MP.4` stays at `08-modules-and-packages/4-build-tags` as an optional stretch lesson

## References

1. [Using Go Modules](https://go.dev/blog/using-go-modules)
2. [Go Modules Reference](https://go.dev/ref/mod)
3. [Module version numbering](https://go.dev/doc/modules/version-numbers)

## Next Step

After `MP.3`, continue to [Section 09: I/O and CLI Tools](../09-io-and-cli).
If you want one more stretch lesson before moving on, finish `MP.4` first.
