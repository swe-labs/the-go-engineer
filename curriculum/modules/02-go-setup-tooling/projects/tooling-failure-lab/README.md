# Tooling Failure Lab

## Learning objective

By the end of this lab, you will diagnose and fix common Go tooling failures using every tool from Module 02 — `go version`, `go run`, `go build`, `go test`, `gofmt`, `go vet`, `go doc`, `gopls`, and `go mod` commands. You will apply systematic debugging to resolve each failure and document your process.

## Why this matters

Knowing the tooling commands is not enough. Real development means encountering failures and using the tooling output to diagnose them. This lab simulates the most common Go setup and tooling problems so you can practice the diagnostic workflow before it matters on a real project.

## Prerequisites

Before starting this lab, complete all Module 02 lessons:
- 01 — Install and Verify Go
- 02 — Hello World
- 03 — go run
- 04 — go build
- 05 — go test
- 06 — gofmt
- 07 — go vet
- 08 — go doc
- 09 — Editor Setup and gopls
- 10 — Reading Compiler Errors
- 11 — Reading Runtime Errors
- 12 — Reading Test Failures
- 13 — Go Module Root Basics
- 14 — Tooling Checklist

Also complete Module 01 (Computers, Terminal, Git, and the Web) and Module 00 (Orientation).

## Available concepts

This lab draws on the following concepts. If any are unfamiliar, review the corresponding lesson:

- **Curriculum navigation** — how to find files, run commands, navigate the repo
- **code execution workflow** — `go run`, `go build`, `go test` flow
- **starter/solution pattern** — `_starter/` has broken code, `_solution/` has fixes
- **debugging methodology** — read the error, form a hypothesis, test the fix
- **process lifecycle** — how programs run and exit
- **stdin/stdout/stderr** — where output and errors go
- **shell pipelines** — chaining commands in the terminal
- **go build** — compiling Go source into a binary
- **go module cache** — where Go stores downloaded dependencies
- **compiler errors** — reading file:line:col error messages

## Lab overview

The `_starter/` directory contains a Go program with **5 deliberate failures**. Each failure is designed to be discovered by one of the Go tooling commands you learned in Module 02. Your job is to:

1. Run each tool and observe what it reports.
2. Identify which failure each tool uncovers.
3. Fix each issue using what you learned in the corresponding lesson.
4. Verify the fix by re-running the tool.

You may check `_solution/` only after you have attempted all fixes yourself.

## Tasks

### Task 1 — Discover tooling failures

Run each command below from the `_starter/` directory. For each command, record:

- What the command reported (output or error)
- Which tool from Module 02 you used
- What the failure was
- How you fixed it

```bash
go run .
go build .
go test -v .
gofmt -d main.go
go vet .
go doc Multiply
```

### Task 2 — Fix all failures

Edit the files in `_starter/` to fix every issue. The program should:

- Compile and run without errors
- Print the correct multiplication results
- Pass all tests
- Pass `gofmt -l .` with no output
- Pass `go vet .` with no output
- Show documentation when running `go doc Multiply`

### Task 3 — Document your troubleshooting

Write a troubleshooting document (`TROUBLESHOOTING.md`) that describes:

- Each failure you encountered
- The tool that helped you discover it
- The root cause of the failure
- How you fixed it

Use this format for each entry:

```markdown
## Failure: [short name]

- **Tool**: `go run .`
- **Error**: `package is not a main package`
- **Root cause**: The file was missing `package main` declaration
- **Fix**: Added `package main` at the top of the file
```

## Verification

After fixing all issues, run the following commands from `_starter/` and confirm each succeeds:

```bash
go run .
```
Expected output:
```
Multiply(3, 4) = 12
Multiply(-2, 5) = -10
Multiply(0, 7) = 0
All checks passed.
```

```bash
go test -v .
```
Expected output:
```
=== RUN   TestMultiply
--- PASS: TestMultiply (0.00s)
=== RUN   TestMultiplyZero
--- PASS: TestMultiplyZero (0.00s)
=== RUN   TestMultiplyNegative
--- PASS: TestMultiplyNegative (0.00s)
PASS
ok      tooling-failure-lab      0.xxx
```

```bash
gofmt -l .
```
Expected output: *(nothing — all files formatted)*

```bash
go vet .
```
Expected output: *(nothing — no issues)*

```bash
go doc Multiply
```
Expected output:
```
Multiply returns the product of two integers.
```

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------------|----------|----------------|
| **Correctness** | All 5 failures identified and fixed; all verification commands pass | 3-4 failures fixed; most verification passes | Fewer than 3 failures fixed |
| **Diagnostic quality** | Systematic approach: run tool, read error, identify cause, apply fix, verify | Some systematic steps but skips verification | Random trial-and-error; no record of diagnostic steps |
| **Documentation** | TROUBLESHOOTING.md covers all failures with root cause and fix for each | Covers most failures but missing root cause details | Missing or incomplete documentation |
| **Tool mastery** | Uses all 6 tools (go run, go build, go test, gofmt, go vet, go doc) correctly | Uses 4-5 tools correctly | Uses fewer than 4 tools |

## NEXT UP

[Module 03 — Programming Fundamentals with Go](../../../03-programming-fundamentals/README.md)
