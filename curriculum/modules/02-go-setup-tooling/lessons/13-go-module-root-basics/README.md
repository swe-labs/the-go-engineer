# Lesson: Go module root basics

## Mission

Understand `go.mod` as the root declaration for a Go module and learn why commands behave differently inside and outside it.

By the end, you should know what the tool or concept proves, what failure looks like, and what first debugging move to make.

## Prerequisites

You should have completed the previous lesson:

- core-02-12 (Reading test failures)

Understand how to read `go test -v` output, identify failing test names, and locate file:line failures.

## Mental Model

`go.mod` is the address label and dependency contract for a Go workspace.

This model should help you predict what the tool is doing before you memorize command flags.

## Visual Model

```text
command
  ↓
Go tool or shell
  ↓
compile / inspect / run / test / format
  ↓
output, artifact, or error
  ↓
next debugging action
```

The key habit is to ask: what stage produced this output?

## Machine View

The Go command searches for `go.mod` to determine the module root, module path, Go version, dependencies, and module cache behavior.

The Go toolchain is not one command. It is a family of commands that read files, inspect modules, compile packages, run binaries, cache artifacts, and report errors.

## Run Instructions

From the repository root:

```bash
go run ./curriculum/modules/02-go-setup-tooling/lessons/13-go-module-root-basics
go test ./curriculum/modules/02-go-setup-tooling/lessons/13-go-module-root-basics
```

Related commands to recognize:

- `go env GOMOD`
- `go list ./...`
- `go mod tidy`

## Code Walkthrough

Open `main.go`.

The lesson program uses a `toolCard` to model the tooling concept:

1. `ID` connects this file to curriculum metadata.
2. `Title` names the lesson.
3. `MentalModel` gives a beginner-safe analogy.
4. `MachineView` explains what the computer or Go tool sees.
5. `CommandPurpose` explains why the command exists.
6. `CommonMistake` names the failure pattern.
7. `Fix` gives the first recovery step.

This is not advanced Go yet. It is a small executable proof that the concept can be represented, run, and tested.

## Try It

Run `go env GOMOD` and explain why the output changes depending on your directory.

Then modify one field in `main.go`, rerun the program, and explain why the output changed.

## Common Mistakes

| Mistake | Why it happens | Correction |
|---|---|---|
| Running Go commands outside the module and blaming packages or imports. | The command hides more than one stage. | Use `go env GOMOD`, `pwd`, and the repository root before debugging imports. |
| Running from the wrong directory | Relative paths and module roots depend on location. | Run `pwd` and `go env GOMOD`. |
| Treating green editor UI as proof | Editors provide diagnostics, not full release validation. | Confirm with terminal commands. |

## Debugging Signals

Use these signals:

- `go: command not found` → shell cannot find the Go executable
- `go.mod file not found` → command is outside a module or module mode is unexpected
- `package ... is not in std` → path/import/module confusion is likely
- compile error with file/line/column → source did not type-check
- test failure with `got` and `want` → behavior disagrees with expectation

## In Production

CI, editors, dependency resolution, build cache behavior, and import paths all depend on the correct module root.

A reliable engineer treats tooling output as evidence. You do not need to memorize every flag immediately, but you do need to know what kind of proof each command creates.

## Performance Notes

The Go build and test cache can make repeated commands faster. Do not confuse cached speed with skipped correctness. If output seems surprising, inspect the exact command and package path.

## Security Notes

Do not paste secrets into commands, logs, tests, or documentation. Be careful with environment variables, shell history, and CI logs.

## Thinking Questions

1. What stage does this lesson focus on: install, compile, run, test, format, inspect, or diagnose?
2. What output proves the command worked?
3. What is one common false assumption about this tool?
4. What would be your first debugging command if it failed?
5. How will this tool appear in CI later?

## Proof of Understanding

You are complete when you can:

- explain `Go module root basics` in your own words
- run the example
- run the test
- name one failure mode
- explain the first fix
- identify the next lesson

## Next Step

Continue to:

```text
core-02-14 — ../14-tooling-checklist/README.md
```
