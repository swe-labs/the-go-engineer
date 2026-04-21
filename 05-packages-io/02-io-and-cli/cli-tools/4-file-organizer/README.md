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

## Next Step

After you complete this exercise, continue to the [Encoding track](../../encoding) or back to the
[Stage 05 overview](../../README.md).
