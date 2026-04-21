# Stage 05: Modules and Packages

## Mission

This track teaches how Go modules define package boundaries, dependency resolution, versioned imports, and controlled build surfaces.

By the end of this track, you should be comfortable reading and writing:

- `go.mod` and `go.sum` without treating them as magic files
- dependency workflows using `go get`, `go mod tidy`, `go list`, and `go mod why`
- semantic-versioning decisions, especially the `/v2` import-path rule
- local override workflows with `replace`
- optional platform and test surfaces through build tags

## Stage Ownership

This track belongs to [05 Packages & IO](../README.md).

## Track Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `MP.1` | Lesson | [module basics](./1-module-basics) | Explains what `go.mod` and `go.sum` actually do. | entry |
| `MP.2` | Lesson | [managing deps](./2-managing-deps) | Shows how to add, inspect, trim, and explain dependencies. | `MP.1` |
| `MP.3` | Exercise | [versioning workshop](./3-versioning) | Combines semantic versioning, `/v2` imports, and `replace`. | `MP.1`, `MP.2` |
| `MP.4` | Stretch Lesson | [build tags](./4-build-tags) | Adds conditional compilation once the module workflow is stable. | `MP.1`, `MP.2` |

## Suggested Order

1. Work through `MP.1` and `MP.2` in order.
2. Complete `MP.3` before moving on.
3. Use `MP.4` as an optional stretch lesson.

## Track Milestone

`MP.3` is the primary milestone for this track.

You are ready for the next track when you can explain:

- why `go.mod` and `go.sum` are part of the build contract
- why a module that reaches v2 must change its import path
- why `replace` is useful in local development but should be used deliberately

## Next Step

After `MP.3`, continue to [05 Stage 05: I/O and CLI](../02-io-and-cli).
