# Project Assessment — Order Summary Refactor

## Purpose

Evaluate your ability to design custom error types, wrap errors with `fmt.Errorf` and `%w`, use `errors.Is` and `errors.As` for structured error matching, aggregate validation errors with `errors.Join`, and apply `defer` for resource cleanup. This assessment verifies that you can refactor a panic-heavy codebase into idiomatic Go error handling.

## Targets

- `project-order-summary-refactor`

## Scope

This assessment covers:

- Custom error types implementing the `error` interface
- `fmt.Errorf` with `%w` for error wrapping and chain preservation
- Sentinel errors (`ErrEmptyItem`, `ErrNegativeQty`) and `errors.Is` matching
- `errors.As` for extracting custom error types from wrapped chains
- `errors.Join` for multi-error aggregation
- `defer` for guaranteed resource cleanup
- The `(result, error)` return pattern replacing panics
- Table-driven tests covering error paths

## Evidence required

Submit:

1. Your refactored `_starter/main.go` with custom error types, wrapping, and deferred cleanup
2. Your test file with table-driven tests covering success, validation failure, sentinel matching, custom type extraction, and multi-error scenarios
3. Output of `go test -v .` showing all tests passing
4. A brief explanation of one error chain you constructed: show the wrapping layers and how `errors.Is` or `errors.As` traverses the chain

## Time expectation

Plan for 90 to 120 minutes.

## Rubric

| Criterion | Weight | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------|--------------|----------|----------------|
| **Error type design** | 25% | Custom types implement `error`, `Unwrap`, and carry descriptive fields (`Field`, `Rule`, `Code`); sentinel errors are package-level variables | Types implement `error` but missing `Unwrap` or fields are generic | No custom error types or incorrect interface implementation |
| **Wrapping correctness** | 25% | Every error uses `%w` with context; error chain is preserved end-to-end; `errors.Is` and `errors.As` both work on wrapped chains | Most errors use `%w` but some use `%v` or `errors.New` losing the chain | Errors are flat strings with no wrapping; chain is broken |
| **Multi-error and cleanup** | 20% | `errors.Join` collects all validation failures; `defer` guarantees cleanup even on error paths; no panics for expected failures | `errors.Join` used but cleanup missing or panics remain | Neither multi-error nor defer used; panics on expected failures |
| **Test coverage** | 20% | Table-driven tests cover success, each validation failure, sentinel matching via `errors.Is`, custom type extraction via `errors.As`, and multi-error unwrapping | Tests cover success and basic failure but not error matching | Tests incomplete or missing error-path coverage |
| **Code organization** | 10% | Clean separation: types in logical groups, validation distinct from processing, `defer` used idiomatically | Reasonable but some mixing of validation and business logic | Monolithic function with inline error strings |

**Passing score**: 75%

## Review questions

Be prepared to answer these questions during a review session:

1. What is the difference between `fmt.Errorf("failed: %v", err)` and `fmt.Errorf("failed: %w", err)`? Why does this distinction matter for `errors.Is`?
2. Show me your `ValidationError` type. How does `errors.As` find it when it is wrapped inside another error?
3. Walk through your `errors.Join` usage: what happens when you call `errors.Is` on a joined error — does it check each sub-error?
4. Why did you use `defer` for cleanup instead of putting the cleanup call at the end of the function? What happens if a return statement exits early?
5. Your original code used `panic` for bad input. Why is `(result, error)` better for expected failures? When is `panic` still appropriate in Go?

## Passing standard

You pass when all tests pass, no panics exist for expected failures, your error chains are traversable with `errors.Is`/`errors.As`, and you can explain the wrapping mechanics without guessing.

## Retake policy

If you score below 75%, review the Module 04 lessons on error wrapping, custom error types, and defer, then refactor and resubmit with a correction note.

## NEXT UP

[Module 04 Checkpoint](../checkpoint/README.md)
