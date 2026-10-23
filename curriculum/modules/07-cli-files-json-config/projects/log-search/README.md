# Log Search

## Learning objective

Build a high-performance streaming log search CLI tool in Go. Master buffered I/O with `bufio.Scanner`, memory-efficient stream processing with bounded memory footprints, regular expression compilation, log-level hierarchy filtering, and contextual window buffering (`-B` and `-A` lines).

## Why this matters

Production log files frequently reach gigabytes in size. Inexperienced developers often load entire files into memory using `os.ReadFile`, triggering out-of-memory (OOM) crashes in production environments. Writing robust infrastructure tools requires streaming data line-by-line using `bufio.Scanner`, maintaining constant O(1) memory regardless of input file size.

## Prerequisites

Complete all Module 07 lessons (01–21) on streaming I/O with `io.Reader`, `bufio.Scanner`, regular expressions with `regexp`, command-line flags with `flag`, and string parsing.

## Available concepts

This project draws on `bufio`, `io`, `regexp`, `strings`, `flag`, ring buffers for preceding context, and streaming error handling.

## Project overview

The `_starter/` directory provides type signatures, log-level enums, and CLI flag bindings. Your job is to implement the stream searching engine so that:

- Input streams (`io.Reader` from file or `os.Stdin`) are parsed line-by-line without accumulating unbounded memory.
- Substring searches support case-insensitive and case-sensitive modes.
- Regular expression queries match complex patterns with syntax validation.
- Log level filtering respects severity hierarchies (`DEBUG` < `INFO` < `WARN` < `ERROR`).
- Context lines around matches (`-B <num>` lines before, `-A <num>` lines after) are emitted with match markers (`>` for hit, ` ` for context).

## CLI Flags

| Flag | Default | Description |
|---|---|---|
| `-file` | `(stdin)` | Log file path to search |
| `-query` | `""` | Search string or regex pattern |
| `-regex` | `false` | Treat query as regular expression |
| `-case-sensitive` | `false` | Enable case-sensitive matching |
| `-level` | `ALL` | Minimum log level filter (`DEBUG`, `INFO`, `WARN`, `ERROR`) |
| `-B` | `0` | Number of context lines before match |
| `-A` | `0` | Number of context lines after match |

## Tasks

### Task 1 — Implement stream scanner & query matching

In `_starter/main.go`:
- In `SearchStream(ctx context.Context, r io.Reader, cfg SearchConfig) (SearchResult, error)`:
  - Compile query string or regex pattern. Validate regex and return `ErrInvalidRegex` on syntax error.
  - Initialize `bufio.Scanner(r)` with an expanded 1MB buffer capacity.
  - Process lines sequentially, counting `TotalScanned`.

### Task 2 — Implement log level filtering

- Implement `extractLogLevel(line string) string`:
  - Inspect line for common level markers like `[INFO]`, `[ERROR]`, `LEVEL=WARN`.
- Compare extracted level against `cfg.MinLevel` priority.

### Task 3 — Implement context window buffering

- Maintain a ring buffer of capacity `cfg.ContextBefore` for preceding context lines.
- When a line matches:
  - Flush all lines in the before-buffer marked with `IsMatch = false`.
  - Append the matching line marked with `IsMatch = true`.
  - Set a countdown counter for `cfg.ContextAfter` lines.
- For non-matching lines within the after-countdown, emit them with `IsMatch = false`.

## Verification

Run all automated unit tests:

```bash
go test -v ./curriculum/modules/07-cli-files-json-config/projects/log-search/_starter
```

Expected output:
```
=== RUN   TestSearchStream_Substring
--- PASS: TestSearchStream_Substring (0.00s)
=== RUN   TestSearchStream_Regex
--- PASS: TestSearchStream_Regex (0.00s)
=== RUN   TestSearchStream_MinLevelFiltering
--- PASS: TestSearchStream_MinLevelFiltering (0.00s)
=== RUN   TestSearchStream_ContextLines
--- PASS: TestSearchStream_ContextLines (0.00s)
PASS
```

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|---|---|---|---|
| **Streaming & Memory Efficiency** | Line-by-line processing via `bufio.Scanner`; memory remains constant O(1) regardless of input size | Scanner used with reasonable memory | Entire file loaded into memory with `os.ReadFile` |
| **Pattern & Regex Handling** | Substring and regex searches handled correctly; case sensitivity toggle works; regex errors caught | Searches work but case sensitivity buggy | Regex compilation errors panic |
| **Context Window Handling** | Before and after context lines accurately collected using a ring buffer; match markers clean | Context mostly works with off-by-one line counts | Context lines missed or duplicated |
| **Severity Hierarchy** | Log levels correctly ranked (`ERROR` > `WARN` > `INFO` > `DEBUG`); filters accurately | Levels checked with simple string matching | Misclassifies log severity |

## NEXT UP

[Config Parser](../config-parser/README.md)
