# ST.6 Config Parser Project

## Mission

Build a small config parser that turns `.env`-style text into structured data and then renders a
stable summary from that data.

This exercise is the Section 07 milestone.
It is where string helpers, regex parsing, scanner-based reading, and template-driven output come
together in one runnable artifact with tests.

## Prerequisites

Complete these first:

- `ST.1` strings
- `ST.2` formatting
- `ST.4` regex
- `ST.5` text templates

## What You Will Build

Implement a small config parser that:

1. reads multi-line config content with `bufio.Scanner`
2. uses one compiled regex to parse key-value lines
3. ignores blank lines and comments
4. stores parsed values in a `map[string]string`
5. renders a stable summary using `text/template`
6. passes the provided parsing tests

## Files

- [main.go](./main.go): complete solution with teaching comments
- [parser_test.go](./parser_test.go): tests for the parsing contract
- [_starter/main.go](./_starter/main.go): starter file with TODOs and requirements

## Run Instructions

Run the completed solution:

```bash
go run ./07-strings-and-text/6-config-parser
```

Run the tests:

```bash
go test ./07-strings-and-text/6-config-parser
```

Run the starter:

```bash
go run ./07-strings-and-text/6-config-parser/_starter
```

## Success Criteria

Your finished solution should:

- compile the regex once instead of rebuilding it inside the scan loop
- parse quoted and unquoted values into a map cleanly
- skip comments and empty lines without special-case chaos
- return malformed config lines as errors instead of printing from the parser
- render output through a template instead of scattered print statements
- pass the provided tests

## Common Failure Modes

- compiling the regex on every line
- splitting on `=` blindly and breaking quoted values
- treating the map output order as stable without sorting before rendering
- printing ad hoc output instead of using a template for the final summary

## Next Step

After you complete this exercise, continue to [Section 08](../../08-modules-and-packages) if you
are ready to move on.
