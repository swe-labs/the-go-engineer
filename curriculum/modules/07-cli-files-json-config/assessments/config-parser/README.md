# Project Assessment — Config Parser

## Purpose

Evaluate your ability to build a layered configuration system following the Twelve-Factor App methodology. This assessment verifies that you can implement precedence hierarchies (defaults < file < environment), parse JSON and key-value formats, validate schema invariants, and design testable code with dependency injection.

## Targets

- `project-config-parser`

## Scope

This assessment covers:

- JSON parsing with `encoding/json` and struct tags
- Key-value file parsing (line-by-line `KEY=VALUE`)
- File extension detection for format routing
- Environment variable overrides with `os.LookupEnv`
- Dependency injection for mockable environment lookups
- `strconv.Atoi` and `strings.EqualFold` for type conversions
- Schema validation with error aggregation
- Sentinel errors (`ErrUnsupportedFormat`, `ErrValidation`)
- Configuration precedence: Defaults < File < Environment

## Evidence required

Submit:

1. Your completed `_starter/main.go` with file loading, environment overrides, and validation implemented
2. Output of `go test -v .` showing all tests passing
3. A brief explanation of why `envLookup` is injected as a function parameter instead of calling `os.LookupEnv` directly — what problem does this solve for testing?

## Time expectation

Plan for 60 to 90 minutes.

## Rubric

| Criterion | Weight | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------|--------------|----------|----------------|
| **Precedence ordering** | 30% | Environment variables strictly override file configs, which strictly override defaults; precedence verified for all field types (string, int, bool) | Overrides work but order partially reversed or incomplete for nested struct fields | Overrides fail; final config uses defaults regardless of env |
| **Schema validation** | 25% | Port range (1–65535), required strings, and positive bounds enforced; clean aggregated error message wrapping `ErrValidation` | Basic validation present with partial checks | Invalid configurations accepted silently |
| **Format flexibility** | 20% | JSON and KV formats both supported; format detection based on file extension is clean; `ErrUnsupportedFormat` returned for unknown extensions | Only JSON supported or KV parsing fragile | Parsing crashes on valid JSON or panics on malformed files |
| **Testability and design** | 25% | `envLookup` dependency injected via constructor; all tests use mock environment without touching `os.Setenv`; deterministic and parallel-safe | Tests work but mutate global process environment with `os.Setenv` | Untested override logic or no dependency injection |

**Passing score**: 75%

## Review questions

1. Explain the Twelve-Factor App configuration principle. Why should config come from the environment rather than hardcoded values?
2. Walk through the precedence: if `config.json` sets `server.port = 3000` and `APP_SERVER_PORT=8080` is set in the environment, what is the final port? Why?
3. Your `envLookup` is a function parameter. How does this make testing easier compared to calling `os.LookupEnv` directly inside the function?
4. What does your validation return when both `server.port = 0` and `database.name = ""`? Is it a single error or an aggregated error? How does the caller inspect it?
5. What happens when a user provides a file with extension `.yaml`? Show the error path from detection to `ErrUnsupportedFormat`.

## Passing standard

You pass when all tests pass, precedence is correctly enforced, validation catches all invariants, and you can explain dependency injection and error aggregation mechanically.

## Retake policy

If you score below 75%, review Module 07 lessons on JSON parsing, environment variables, and configuration patterns, then resubmit.

## NEXT UP

[Module 07 Checkpoint](../checkpoint/README.md)
