# Order Summary Refactor

## Learning Objectives

- Design custom error types that implement the `error` interface
- Wrap errors with `fmt.Errorf` and `%w` to preserve the error chain
- Use `errors.Is` and `errors.As` to match and extract specific errors
- Collect multiple validation errors with `errors.Join`
- Use `defer` for resource cleanup
- Write table-driven tests for error paths

## Background

You've inherited a working-but-untidy order summary calculator. It processes orders, validates items, and computes totals, but it has no custom error types, no error wrapping, no deferred cleanup, and panics for expected failures. Your job is to refactor it into idiomatic Go that uses the error-handling patterns from this module.

## Requirements

1. **Custom error types** — Define `ValidationError` and `OrderError` that implement the `error` interface. `ValidationError` should carry a `Field` and `Rule`. `OrderError` should carry a `Code` and wrap an underlying `Err`.

2. **Error wrapping** — Use `fmt.Errorf` with `%w` to annotate every error with context (which operation failed, on what input).

3. **`errors.Is` and `errors.As`** — Define sentinel errors (`ErrEmptyItem`, `ErrNegativeQty`) and use `errors.Is` to match them. Use `errors.As` to extract `ValidationError` and `OrderError` from the chain.

4. **`errors.Join`** — In your validation function, collect all field-level failures and return them as a single multi-error with `errors.Join`.

5. **`defer` for cleanup** — Simulate a resource (e.g., a file handle or mutex lock) and use `defer` to release it.

6. **`(result, error)` return pattern** — Every function that can fail returns `(value, error)`. Do not panic for expected failures.

7. **Table-driven tests** — Cover success, validation failure, sentinel matching, custom type extraction, and multi-error scenarios.

## Starter Code Overview

The starter `main.go` provides a flat `processOrder` function that validates items inline, uses `errors.New` for every error, panics on bad input, and ignores cleanup. TODO markers identify where each refactoring should be applied.

## Project Structure

```
order-summary-refactor/
  main.go          -- starter / your working file
  main_test.go     -- your tests
  _solution/       -- reference solution (after completing)
```

## Verification

```bash
go test ./curriculum/modules/04-functions-errors-data-semantics/projects/order-summary-refactor/...
```

## Rubric

| Criteria | Excellent | Good | Needs Work |
|---|---|---|---|
| Error type design | Custom types implement `error`, `Unwrap`, and `Is` with clear fields | Types implement `error` but missing `Unwrap` or `Is` | No custom types or incorrect interface |
| Wrapping correctness | Every error uses `%w`, chain is preserved, messages are descriptive | Most errors use `%w`, some use `%v` or `errors.New` | Errors are flat strings with no wrapping |
| Error handling completeness | `errors.Is` / `errors.As` used correctly; sentinel and custom types both tested | One matching strategy used but not both | No structured error matching |
| Test coverage | Table-driven tests cover success, each failure mode, error matching, and multi-error | Tests cover success and basic failure | Tests are incomplete or missing |
| Code organization | Clean separation of types, validation, processing; `defer` used appropriately | Reasonable organization but some mixing of concerns | Monolithic, no `defer`, panics for expected failures |
