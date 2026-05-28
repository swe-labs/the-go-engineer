# Module 02 — Go Setup and Tooling

## Mission

This module makes the Go toolchain predictable from day one.

Before programming fundamentals begin, you need to trust the basic loop:

```text
edit → format → run/build → test → inspect errors → repeat
```

If the toolchain feels magical, every later failure becomes harder. This module removes that magic by explaining what each command does, what output means, and how to recover from common setup problems.

## Who this module is for

This module is for:

- beginners installing Go for the first time
- learners who can run commands but do not know what they prove
- developers new to Go tooling
- contributors who need the repository's run/test/validation workflow

## What you will build

You will build a working Go tooling workflow:

- verify Go installation
- run a Hello World program
- use `go run`
- use `go build`
- use `go test`
- format with `gofmt`
- inspect suspicious code with `go vet`
- read documentation with `go doc`
- connect editor tooling to terminal proof
- read compiler, runtime, and test failures
- understand module roots and `go.mod`
- use a repeatable tooling checklist

## Prerequisites

You should have completed:

```text
curriculum/modules/01-computers-terminal-git-web/
```

You should be able to:

- open a terminal
- use `pwd`, `ls`, and `cd`
- explain source code vs executable vs process
- explain basic exit codes
- understand why current working directory matters

## Concept map

```text
install Go
  ↓
verify go command
  ↓
hello world
  ↓
go run → go build → go test
  ↓
gofmt → go vet → go doc
  ↓
editor + gopls
  ↓
compiler errors → runtime errors → test failures
  ↓
module root + go.mod + module cache
  ↓
repeatable tooling checklist
```

## Lessons

| # | ID | Lesson | Outcome |
|---:|---|---|---|
| 1 | `core-02-01` | [Install and verify Go](./lessons/01-install-and-verify-go/README.md) | Install Go and verify that the `go` command, version, environment, and workspace assumptions are visible. |
| 2 | `core-02-02` | [Hello World](./lessons/02-hello-world/README.md) | Write and run the smallest useful Go program while understanding package, import, function, and output at a beginner level. |
| 3 | `core-02-03` | [go run](./lessons/03-go-run/README.md) | Understand `go run` as compile-then-execute for quick feedback, not as an interpreter. |
| 4 | `core-02-04` | [go build](./lessons/04-go-build/README.md) | Understand `go build` as the command that produces reusable executable artifacts. |
| 5 | `core-02-05` | [go test](./lessons/05-go-test/README.md) | Learn how `go test` compiles packages with test files and runs `Test...` functions as proof. |
| 6 | `core-02-06` | [gofmt](./lessons/06-gofmt/README.md) | Use `gofmt` to make Go formatting automatic and non-negotiable. |
| 7 | `core-02-07` | [go vet](./lessons/07-go-vet/README.md) | Use `go vet` as a static analyzer that catches suspicious code patterns tests may miss. |
| 8 | `core-02-08` | [go doc](./lessons/08-go-doc/README.md) | Learn to inspect package documentation from the terminal instead of guessing API behavior. |
| 9 | `core-02-09` | [Editor setup and gopls](./lessons/09-editor-setup-and-gopls/README.md) | Understand what your editor and Go language server do, and how to verify they are helping rather than hiding problems. |
| 10 | `core-02-10` | [Reading compiler errors](./lessons/10-reading-compiler-errors/README.md) | Learn to read compiler errors as structured location-and-cause reports, not as scary walls of text. |
| 11 | `core-02-11` | [Reading runtime errors](./lessons/11-reading-runtime-errors/README.md) | Understand runtime errors as failures that happen after a program starts executing. |
| 12 | `core-02-12` | [Reading test failures](./lessons/12-reading-test-failures/README.md) | Learn to use test failure output to identify expected behavior, actual behavior, and the smallest broken assumption. |
| 13 | `core-02-13` | [Go module root basics](./lessons/13-go-module-root-basics/README.md) | Understand `go.mod` as the root declaration for a Go module and learn why commands behave differently inside and outside it. |
| 14 | `core-02-14` | [Tooling checklist](./lessons/14-tooling-checklist/README.md) | Build a repeatable checklist for verifying Go tooling before writing larger programs. |

## Labs

This module has no separate lab. Each lesson includes a runnable proof task and starter/solution practice surface.

## Projects

This module has no portfolio project. Portfolio artifacts begin after enough programming fundamentals exist.

## Assessments

Complete the checkpoint:

```text
curriculum/modules/02-go-setup-tooling/assessments/checkpoint/
```

## Common failure modes

| Failure | Symptom | Fix |
|---|---|---|
| Go not on PATH | `go: command not found` | Reopen terminal, inspect PATH, verify install path. |
| Wrong module root | imports or `go list` behave unexpectedly | Run `pwd` and `go env GOMOD`. |
| Treating `go run` as interpretation | compile errors feel surprising | Remember `go run` compiles before executing. |
| Ignoring formatting | noisy diffs and inconsistent style | Run `gofmt` or enable format-on-save. |
| Misreading errors | fixing random files | Read file, line, column, first error, then inspect source. |
| Trusting editor only | CI fails after local green editor UI | Confirm with terminal commands. |

## Completion checklist

You are ready for Module 03 when you can:

- [ ] verify Go installation with `go version`
- [ ] inspect Go environment with `go env`
- [ ] run a small Go program
- [ ] build a Go command
- [ ] run package tests
- [ ] format files with `gofmt`
- [ ] run `go vet`
- [ ] use `go doc`
- [ ] explain compiler vs runtime vs test failures
- [ ] find the module root with `go env GOMOD`
- [ ] complete the tooling checklist
- [ ] pass the Module 02 checkpoint

## Next module

Continue to:

```text
curriculum/modules/03-programming-fundamentals/
```

Module 03 starts programming fundamentals with a stable Go toolchain underneath you.
