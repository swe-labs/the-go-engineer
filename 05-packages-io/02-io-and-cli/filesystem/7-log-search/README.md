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

Log searching is a daily operational activity. When a production incident occurs at 3 AM, the first response is almost always "search the logs." Tools like `grep`, `ripgrep`, and centralized log platforms (Elasticsearch, Loki, CloudWatch) all implement the same core pattern this exercise teaches: walk a directory, filter by file type, scan line by line, and report matches with context. The line-by-line streaming approach matters because production log files can be gigabytes in size — loading an entire file into memory with `os.ReadFile` would crash the search tool or the host machine. The `bufio.Scanner` pattern keeps memory usage constant regardless of file size. In real systems, you would also need to handle rotated logs (`.log.1`, `.log.gz`), binary files that should be skipped, and encoding issues in logs written by different services. The ability to build a reliable file-scanning tool from standard library primitives — without external dependencies — is a skill that transfers directly to production debugging.

## Thinking Questions

1. Why does `bufio.Scanner` use a fixed-size buffer, and what happens when a log line exceeds that buffer size?
2. How would you extend this tool to search compressed `.gz` log files without extracting them to disk first?
3. What are the trade-offs between searching logs locally with a tool like this versus shipping them to a centralized logging platform?
4. If you need to search 100 GB of log files quickly, what concurrency strategy would you use and what bottleneck would you hit first — CPU or disk I/O?

## Next Step

After you complete this exercise, continue to [`FS.8` fs.FS testing seam](../8-fs-testing-seam)
or back to the [Stage 05 overview](../../README.md).


