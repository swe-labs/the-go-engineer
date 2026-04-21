# FS.7 Log Search Tool

## Mission

Build a small log-search tool that walks directories, scans files line by line, and reports
matching results without loading everything into memory.

This exercise is the Filesystem track milestone for Stage 05.

## Prerequisites

Complete these first:

- `FS.1` files
- `FS.2` paths
- `FS.3` directories
- `FS.6` I/O patterns

## What You Will Build

Implement a search tool that:

1. walks a root directory recursively
2. filters for `.log` and `.txt` files
3. scans files line by line with `bufio.Scanner`
4. performs case-insensitive matching
5. reports filename, line number, and line text for each match

## Files

- [main.go](./main.go): complete solution with teaching comments
- [search_test.go](./search_test.go): tests for the search helpers
- [_starter/main.go](./_starter/main.go): starter file with TODOs and requirements

## Run Instructions

Run the completed solution:

```bash
go run ./05-packages-io/02-io-and-cli/filesystem/7-log-search
```

Run the tests:

```bash
go test ./05-packages-io/02-io-and-cli/filesystem/7-log-search
```

Run the starter:

```bash
go run ./05-packages-io/02-io-and-cli/filesystem/7-log-search/_starter
```

## Success Criteria

Your finished solution should:

- traverse directories with `filepath.WalkDir`
- stream file contents line by line instead of reading whole files
- skip unsupported file types cleanly
- keep unreadable files from crashing the full search
- pass the provided tests

## Next Step

After you complete this exercise, continue to [`FS.8` fs.FS testing seam](../8-fs-testing-seam)
or back to the [Stage 05 overview](../../README.md).
