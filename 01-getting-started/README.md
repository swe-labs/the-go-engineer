# Section 01: Getting Started

## Mission

This section helps a complete beginner cross the first real threshold:

- Go is installed and verified
- the terminal is no longer mysterious
- a tiny program can run successfully
- the learner can explain the basic shape of a Go program

The goal is not to race into syntax. The goal is to replace fear with a reliable first-run loop.

## Why This Section Exists Now

Before a learner can think about control flow, data structures, or functions, they need one calm, repeatable truth:

1. open the repo
2. run a program
3. read the output
4. change something small
5. run it again

That loop is the real foundation for everything that follows.

## Section Map

| ID | Lesson | What It Unlocks |
| --- | --- | --- |
| `GT.1` | [installation](./01-installation/) | confirms that Go is installed and runnable |
| `GT.2` | [hello world](./02-hello-world/) | teaches the minimum executable Go program |
| `GT.3` | [how Go works](./03-how-go-works/) | explains packages, imports, and compilation at a beginner-safe level |
| `GT.4` | [development environment](./04-dev-environment/) | establishes the everyday command loop used across the repo |
| `GT.5` | [go tools](./05-go-tools/) | mastering `fmt`, `vet`, and `doc` for professional code |
| `GT.6` | [compiler errors](./06-reading-compiler-errors/) | decoding the machine's feedback and treating it as a partner |

## Zero-Magic Boundary

This section stays intentionally small. It does not try to teach:

- variables in depth
- branching and loops
- functions as a design tool
- data structures as a problem-solving system

It only teaches enough to make the learner operational and calm.

## How To Use This Section

For each lesson:

1. read the `README.md`
2. open `main.go`
3. run the lesson
4. make one small change
5. run it again

That order matters. The README explains what to look for before the learner stares at code.

## Finish This Section When

- `go run` feels familiar instead of scary
- the learner can explain what `package main` and `func main()` are doing
- the learner can move between folders and run lessons without panic
- the learner can use the basic command loop: `go run`, `go fmt`, `go build`, and `go test`

## Next Step

Continue to [Section 02: Language Basics](../02-language-basics/README.md).
