# Getting Started Track

## Mission

This track removes the setup friction that usually blocks beginners before they ever reach real Go
fundamentals.

## Beta Stage Ownership

This track belongs to the public beta stage [0 Foundation](../../docs/stages/00-foundation.md).

It is the setup and first-run half of Section `01`, not the beginning of `1 Language Fundamentals`.

By the end of this track, a learner should be able to:

- verify a working Go installation
- run a simple Go program
- explain the purpose of `package main` and `func main()`
- use the basic command loop: `go run`, `go fmt`, `go vet`, `go build`, `go test`

## Track Map

| ID | Surface | Why It Matters | Requires |
| --- | --- | --- | --- |
| `GT.1` | [installation](./1-installation) | Confirms Go is installed and the environment is real before deeper lessons begin. | entry |
| `GT.2` | [hello world](./2-hello-world) | Teaches the minimum runnable Go program and the shape of `main`. | `GT.1` |
| `GT.3` | [how Go works](./3-how-go-works) | Explains packages, compilation, and exported names in plain terms. | `GT.2` |
| `GT.4` | [development environment](./4-dev-environment) | Establishes the everyday tool loop learners will use for the rest of the repo. | `GT.3` |

## Suggested Order

1. Run `GT.1` and confirm the environment works.
2. Read and run `GT.2` carefully.
3. Use `GT.3` to understand what Go is doing under the hood.
4. Finish `GT.4` before moving into `LB.1` if you are new to Go tooling.

## Where This Leads

After this track, move into [Language Basics](../language-basics/) and start with `LB.1`.

In beta routing terms:

- this track is the end of `0 Foundation`
- `LB.1` is the start of [1 Language Fundamentals](../../docs/stages/01-language-fundamentals.md)

If you already know how to install Go and run simple programs, you can use this track as a quick
tooling refresher instead of a hard gate.
