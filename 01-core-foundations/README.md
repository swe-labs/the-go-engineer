# Section 01: Core Foundations

## Mission

This section gets a learner from "Go is installed and I can run programs" to "I can read and
write small, typed Go programs without guessing."

## Beta Routing Rule

Section `01` is now one legacy source section that feeds two different public beta stages:

- [0 Foundation](../docs/stages/00-foundation.md) owns the canonical [01-foundations/01-getting-started](../01-foundations/01-getting-started/)
- [1 Language Fundamentals](../docs/stages/01-language-fundamentals.md) begins with `language-basics`

That means this section is still one physical source area in the repo, but it is no longer one
single learner-facing stage in the beta architecture.

By the end of Section 01, you should be comfortable with:

- running and reading small Go programs
- basic package and `main` structure
- variables and zero values
- constants and `iota`
- using simple custom types for safer code

## Who Should Start Here

### Full Path

Start with [0 Foundation](../docs/stages/00-foundation.md), complete `GT.1` to `GT.4`, then move
into `LB.1` to `LB.4` as the start of [1 Language Fundamentals](../docs/stages/01-language-fundamentals.md).

### Bridge Path

If Go is already installed and you can run both:

- `go version`
- `go run ./01-foundations/01-getting-started/2-hello-world`

you can skim `GT.1` to `GT.4` and begin the main fundamentals track at `LB.1`.

### Targeted Path

If you only want the live milestone, review these first:

- `GT.2` hello world
- `GT.4` development environment
- `LB.1` variables
- `LB.2` constants
- `LB.3` enums with `iota`

## Section Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `GT.1` | Lesson | [installation](./getting-started/1-installation) | Verifies that the learner can actually run Go before the curriculum asks for anything deeper. | entry |
| `GT.2` | Lesson | [hello world](./getting-started/2-hello-world) | Teaches the minimum runnable Go program and the shape of `main`. | `GT.1` |
| `GT.3` | Lesson | [how Go works](./getting-started/3-how-go-works) | Introduces packages, compilation, and exported names without burying the learner in theory. | `GT.2` |
| `GT.4` | Lesson | [development environment](./getting-started/4-dev-environment) | Gives the learner the basic tool loop: run, format, vet, build, and test. | `GT.3` |
| `LB.1` | Lesson | [variables](./language-basics/1-variables) | Introduces typed values, zero values, and simple declarations. | entry |
| `LB.2` | Lesson | [constants](./language-basics/2-constants) | Separates fixed values from mutable state and introduces grouped declarations. | `LB.1` |
| `LB.3` | Lesson | [enums with iota](./language-basics/3-enums) | Shows how Go models enum-like values with named types and `iota`. | `LB.2` |
| `LB.4` | Exercise | [application logger](./language-basics/4-application-logger) | Combines early fundamentals into one small, realistic logging exercise. | `LB.1`, `LB.2`, `LB.3` |

## Suggested Order

1. Work through `GT.1` to `GT.4` if you are new to Go tooling.
2. Work through `LB.1` to `LB.3` in order.
3. Complete `LB.4` without copying the finished solution line by line.

Public beta interpretation:

- `GT.1` to `GT.4` belong to `0 Foundation`
- `LB.1` to `LB.4` are the entry into `1 Language Fundamentals`

## Section Milestone

`LB.4` is the current live milestone for this section.

If you can complete it and explain:

- why named types are safer than passing raw integers everywhere
- why `iota` is useful for stable log levels
- why bounds checks matter even in simple enum-style code

then you are ready to move into control flow in Section 02.

## Pilot Role In V2

This live v2 slice keeps the current `GT.*` and `LB.*` ids and filesystem layout stable while
making the beta split explicit:

- `GT.1` to `GT.4` are the setup and first-program track for `0 Foundation`
- `LB.1` to `LB.3` are the opening language-fundamentals track for `1 Language Fundamentals`
- `LB.4` is the first real language-fundamentals milestone exercise

That keeps the repo usable for current learners while making the public beta routing honest.

## References

1. [Official Go Installation Guide](https://go.dev/doc/install)
2. [A Tour of Go](https://go.dev/tour/)
3. [Effective Go](https://go.dev/doc/effective_go)

## Next Step

After `LB.4`, continue to Section 02: Control Flow.
