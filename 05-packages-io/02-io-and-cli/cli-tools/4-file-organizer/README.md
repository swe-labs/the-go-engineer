# CL.4 File Organizer

## Mission

Build a small CLI that groups files by extension and can preview its work safely before moving
anything.

This exercise is the CLI track milestone for Stage 05.

## Prerequisites

Complete these first:

- `CL.1` args
- `CL.2` flags
- `CL.3` subcommands

## What You Will Build

Implement a CLI tool that:

1. accepts `--dir` and `--dry-run` flags
2. reads a directory with `os.ReadDir`
3. groups files into subdirectories by extension
4. skips files without extensions cleanly
5. previews moves without mutating the filesystem in dry-run mode

## Files

- [main.go](./main.go): complete solution with teaching comments
- [_starter/main.go](./_starter/main.go): starter file with TODOs and requirements

## Run Instructions

Run the completed solution:

```bash
go run ./05-packages-io/02-io-and-cli/cli-tools/4-file-organizer --dir=./my-folder
```

Run the starter:

```bash
go run ./05-packages-io/02-io-and-cli/cli-tools/4-file-organizer/_starter --dir=./my-folder
```

## Success Criteria

Your finished solution should:

- parse flags with the standard `flag` package
- fail safely when `--dir` is missing
- keep `--dry-run` non-destructive
- create extension directories only when needed
- move or preview each file clearly

## Common Failure Modes

- using `os.Args` manually when typed flags would be clearer
- moving files immediately without a dry-run path
- forgetting to skip directories
- forgetting that `filepath.Ext` returns the leading dot


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

Dry-run modes are a critical safety pattern in production CLI tools. Tools like `terraform plan`, `kubectl diff`, and database migration runners all implement preview-before-mutate behavior because filesystem and infrastructure mutations are difficult or impossible to reverse. In production environments, a file-organizing tool would need to handle concurrent access (another process writing to the same directory), permission errors on restricted directories, symlink loops that cause infinite recursion, and atomic moves across filesystem boundaries where `os.Rename` silently fails. The pattern of separating the "decide what to do" logic from the "do it" logic — which this exercise teaches through the dry-run flag — is the foundation of safe operational tooling. Teams that skip dry-run behavior in internal tools learn the hard way when a script reorganizes a production asset directory and breaks serving paths.

## Thinking Questions

1. Why is `os.Rename` not guaranteed to work when source and destination are on different filesystems, and what would you do instead?
2. If two instances of this tool run concurrently on the same directory, what race conditions could occur and how would you prevent them?
3. How would you extend this tool to support an "undo" operation that reverses the last organize run?
4. Why does `filepath.Ext` include the leading dot, and what edge cases does that create for files like `.gitignore` or `archive.tar.gz`?

## Next Step

After you complete this exercise, continue to the [Encoding track](../../encoding) or back to the
[Stage 05 overview](../../README.md).

