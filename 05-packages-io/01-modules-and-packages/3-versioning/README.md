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


## 



## 



## 



## 



## 



## 




## Mental Model

Think of this as the conceptual blueprint. The components interact by exchanging state, defining clear boundaries between what is requested and what is provided.

## Visual Model

Visualizing this process involves tracing the execution path from the input entry point, through the processing layers, and out to the final output or side effect.

## Machine View

At the hardware level, this translates into specific memory allocations, CPU instruction cycles, and OS-level system calls to manage resources efficiently.

## Solution Walkthrough

The solution demonstrates a complete implementation, proving the concept by bridging the individual requirements into a single, cohesive executable.

## Try It

Run the code locally. Modify the inputs, toggle the conditions, and observe how the output shifts. Experimentation is the fastest way to cement your understanding.

## Verification Surface

The correctness of this component is proven by its associated test suite. We verify boundaries, handle edge cases, and ensure performance constraints are met.

## In Production

Version management is one of the most common sources of production incidents in Go services. When a team bumps a dependency's major version without updating import paths, the build breaks at the worst possible time — usually right before a release. In large organizations with hundreds of internal modules, a single breaking change in a shared library can cascade through dozens of services. Go's module system forces these decisions to be explicit: a `/v2` import path is a loud signal that the API contract changed, and `go.sum` ensures that no dependency changes silently between builds. Teams that treat `replace` directives as permanent fixtures instead of local-development tools often discover that CI builds diverge from local builds, because `replace` only works in the main module. Understanding semantic versioning deeply — not just as a naming convention but as a compatibility contract — prevents entire categories of deployment failures.

## Thinking Questions

1. Why does Go require a new import path (`/v2`) for major version bumps instead of just changing the version number in `go.mod`?
2. If you maintain a library used by 50 internal services, what strategy would you use to roll out a breaking API change without forcing all consumers to update simultaneously?
3. When is a `replace` directive appropriate in a production `go.mod`, and what risks does it introduce if left permanently?
4. How would you detect that a dependency update introduced a subtle behavioral change that does not break the API signature but changes the output?

## Next Step

After you complete this exercise, continue to [Stage 05](../../02-io-and-cli) if you are ready
to move on.
If you want one more stretch lesson first, visit [`MP.4` build tags](../4-build-tags).


