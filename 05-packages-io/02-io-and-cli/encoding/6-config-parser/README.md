# EN.6 Config Parser

## Mission

Build a small JSON config loader that reads from disk, decodes from a stream, and validates
required fields after parsing.

This exercise is the Encoding track milestone for Stage 05.

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
go run ./05-packages-io/02-io-and-cli/encoding/6-config-parser
```

Run the tests:

```bash
go test ./05-packages-io/02-io-and-cli/encoding/6-config-parser
```

Run the starter:

```bash
go run ./05-packages-io/02-io-and-cli/encoding/6-config-parser/_starter
```

## Success Criteria

Your finished solution should:

- decode directly from an open file with `json.NewDecoder`
- validate required fields after parsing
- distinguish missing files from invalid JSON and invalid config data
- use `%w` when wrapping lower-level errors
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

Configuration parsing is one of the first things that runs when a production service starts, and it is one of the most common sources of startup failures. A misconfigured JSON file — a missing comma, a wrong type, or an absent required field — can prevent an entire fleet of servers from booting after a deploy. Production config loaders validate aggressively on startup so that bad configuration fails fast and loudly, rather than silently producing wrong behavior at runtime. The distinction between "file not found," "invalid JSON syntax," and "valid JSON but missing required fields" matters because each failure points to a different operational fix: a missing file means a deployment packaging error, invalid JSON means a human editing mistake, and missing fields mean the config schema changed without updating all environments. Teams that wrap errors properly at each layer — as this exercise teaches with `%w` — can trace the root cause from a single log line instead of guessing.

## Thinking Questions

1. Why is it important to validate config fields after parsing rather than relying on zero values to signal "not set"?
2. If your service reads config from both a file and environment variables, which source should take precedence and why?
3. How would you handle config changes that need to take effect without restarting the service?
4. What risks does `json.NewDecoder` introduce compared to `json.Unmarshal` when reading from an untrusted source?

## Next Step

After you complete this exercise, continue to the [Filesystem track](../../filesystem) or back to
the [Stage 05 overview](../../README.md).


