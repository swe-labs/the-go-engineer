# EN.6 Config Parser

## Mission

Build a small JSON config loader that reads from disk, decodes from a stream, and validates
required fields after parsing.

This exercise is the Encoding track milestone for Section 09.

## Prerequisites

Complete these first:

- `EN.1` JSON marshalling
- `EN.2` JSON unmarshalling
- `EN.3` JSON encoder
- `EN.4` JSON decoder

## What You Will Build

Implement a config loader that:

1. opens a JSON file from disk
2. decodes it with `json.NewDecoder`
3. stores the data in a typed `AppConfig` struct
4. validates required fields after decoding
5. returns useful wrapped errors on failure

## Files

- [main.go](./main.go): complete solution with teaching comments
- [config_test.go](./config_test.go): tests for parsing and validation behavior
- [_starter/main.go](./_starter/main.go): starter file with TODOs and requirements

## Run Instructions

Run the completed solution:

```bash
go run ./09-io-and-cli/encoding/6-config-parser
```

Run the tests:

```bash
go test ./09-io-and-cli/encoding/6-config-parser
```

Run the starter:

```bash
go run ./09-io-and-cli/encoding/6-config-parser/_starter
```

## Success Criteria

Your finished solution should:

- decode directly from an open file with `json.NewDecoder`
- validate required fields after parsing
- distinguish missing files from invalid JSON and invalid config data
- use `%w` when wrapping lower-level errors
- pass the provided tests

## Next Step

After you complete this exercise, continue to the [Filesystem track](../../filesystem) or back to
the [Section 09 overview](../../README.md).
