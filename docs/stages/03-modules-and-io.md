# 3 Modules and IO

## Purpose

`3 Modules and IO` teaches how Go code crosses package, process, encoding, and filesystem
boundaries.

## Who This Is For

- learners who understand core language and type design
- developers who want more confidence with real file, CLI, and serialization work

## Mental Model

Software becomes useful when it can cross boundaries safely.
This stage teaches how code interacts with modules, command-line input, encoded data, and the local
filesystem without becoming messy.

## Why This Stage Exists

This is where learners move from in-memory modeling into real program boundaries.

The goal is to stop treating modules, files, flags, and encoded data as mysterious plumbing and
start treating them as ordinary parts of software engineering.

## What You Should Learn Here

- module and package boundaries
- dependency management basics
- command-line tool patterns
- JSON and other encoding workflows
- filesystem operations, temporary files, and file-search style tasks

## Stage Shape

This stage has two connected parts:

1. `modules-and-packages`
   - module boundaries, dependency management, versioning, and build-surface reasoning
2. `io-and-cli`
   - CLI tools, encoding flows, and filesystem workflows that interact with the outside world

The stage is intentionally ordered from "how is this code organized and versioned?" to
"how does this program interact with real inputs and outputs?"

## Current Source Content

- [08-modules-and-packages](../../08-modules-and-packages/)
- [09-io-and-cli/cli-tools](../../09-io-and-cli/cli-tools/)
- [09-io-and-cli/encoding](../../09-io-and-cli/encoding/)
- [09-io-and-cli/filesystem](../../09-io-and-cli/filesystem/)

## Where This Stage Starts

This stage starts at [08-modules-and-packages](../../08-modules-and-packages/).

If `2 Types and Design` taught you how to shape programs well, this stage teaches you how those
programs behave at package, tool, file, and serialization boundaries.

## Recommended Order

Use this order across the source surfaces:

1. [08-modules-and-packages](../../08-modules-and-packages/)
2. [09-io-and-cli/cli-tools](../../09-io-and-cli/cli-tools/)
3. [09-io-and-cli/encoding](../../09-io-and-cli/encoding/)
4. [09-io-and-cli/filesystem](../../09-io-and-cli/filesystem/)

Do not jump straight into filesystem-heavy tooling before the module and package rules feel clear.
The package boundary model makes the later I/O patterns easier to reason about.

## Path Guidance

### Full Path

Complete Section `08` first, then work through the three Section `09` tracks and their milestone
surfaces.

### Bridge Path

You can move faster if modules, packages, and standard-library tooling already feel somewhat
familiar, but do not skip:

- `MP.3`
- `CL.4`
- `EN.6`
- `FS.7`

Those are the main proof surfaces that show you can work across package, CLI, encoding, and
filesystem boundaries.

### Targeted Path

If you enter this stage late, choose the track that matches your real boundary gap and then check
the earlier milestone surfaces honestly.

Examples:

- weak on modules and dependency behavior: start with `modules-and-packages`
- weak on CLI surfaces and flags: start with `cli-tools`
- weak on JSON and config flows: start with `encoding`
- weak on files and directories: start with `filesystem`

## Stage Milestones

The current milestone backbone is:

- `MP.3` versioning workshop
- `CL.4` file organizer
- `EN.6` config parser
- `FS.7` log search tool

If you can complete those four milestone surfaces honestly, you have the practical foundation for
the next beta stage.

## Finish This Stage When

- you can explain what modules and packages are doing in the repo
- you can build small CLI and file-processing tools without getting lost
- you can serialize and deserialize data confidently
- you understand practical input/output boundaries in Go programs

More concretely:

- you can explain why `go.mod` and `go.sum` are part of the build contract
- you can build small tools that use args, flags, or subcommands without hidden magic
- you can choose between encoding, decoding, and streaming workflows intentionally
- you can write filesystem utilities without tangling traversal, reading, and matching logic

## Next Stage

Move to [4 Backend Engineering](./04-backend-engineering.md).
