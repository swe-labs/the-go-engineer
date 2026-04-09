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

## What You Should Learn Here

- module and package boundaries
- dependency management basics
- command-line tool patterns
- JSON and other encoding workflows
- filesystem operations, temporary files, and file-search style tasks

## Current Source Content

- [08-modules-and-packages](../../08-modules-and-packages/)
- [09-io-and-cli/cli-tools](../../09-io-and-cli/cli-tools/)
- [09-io-and-cli/encoding](../../09-io-and-cli/encoding/)
- [09-io-and-cli/filesystem](../../09-io-and-cli/filesystem/)

## Finish This Stage When

- you can explain what modules and packages are doing in the repo
- you can build small CLI and file-processing tools without getting lost
- you can serialize and deserialize data confidently
- you understand practical input/output boundaries in Go programs

## Next Stage

Move to [4 Backend Engineering](./04-backend-engineering.md).
