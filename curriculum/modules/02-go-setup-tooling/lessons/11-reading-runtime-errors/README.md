# Lesson: Reading runtime errors

## Mission

Understand runtime errors as failures that happen after a program starts executing.

By the end, you should know what the tool or concept proves, what failure looks like, and what first debugging move to make.

## Prerequisites

You should have completed the previous lesson:

- core-02-10 (Reading compiler errors)

Understand the file:line:col format of Go compiler errors and how to read them.

## Mental Model

A runtime error is a machine saying: the instructions were valid enough to start, but something went wrong while doing them.

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

Runtime failures can come from panics, nil pointer use, missing files, invalid inputs, network failures, permissions, or explicit error returns.

The Go toolchain is not one command. It is a family of commands that read files, inspect modules, compile packages, run binaries, cache artifacts, and report errors.

## Run Instructions

From the repository root:

```bash
go run ./curriculum/modules/02-go-setup-tooling/lessons/11-reading-runtime-errors
go test ./curriculum/modules/02-go-setup-tooling/lessons/11-reading-runtime-errors
```

Related commands to recognize:

- `go run .`
- `panic`
- `stack trace`

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

Compare compile failure, test failure, and runtime failure using their timing and evidence.

Then modify one field in `main.go`, rerun the program, and explain why the output changed.

## Common Mistakes

| Mistake | Why it happens | Correction |
|---|---|---|
| Treating every error as a compiler problem. | The command hides more than one stage. | Ask when the failure happened: before build, during test, or while running. |
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

Production debugging often starts with runtime logs, stack traces, status codes, metrics, and recent changes.

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

- explain `Reading runtime errors` in your own words
- run the example
- run the test
- name one failure mode
- explain the first fix
- identify the next lesson

## Next Step

Continue to:

```text
core-02-12 — ../12-reading-test-failures/README.md
```
