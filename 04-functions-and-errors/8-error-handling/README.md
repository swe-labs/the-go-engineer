# FE.9 Error Handling Project

## Mission

Build a small "safe math" program that handles failure explicitly and exposes useful error details
without fragile string matching.

This exercise is the Section 04 milestone.
It is where functions, multiple return values, custom errors, wrapping habits, and defer-based
cleanup come together in one runnable artifact.

## Prerequisites

Complete these first:

- `FE.1` functions
- `FE.4` multiple return values
- `FE.5` custom errors
- `FE.6` error wrapping
- `FE.7` defer
- `FE.8` panic and recover

## What You Will Build

Implement a small safe math library and a runnable demo that:

1. defines a `MathError` type with structured fields
2. implements `safeDivide(a, b int) (float64, error)`
3. implements `safeModulo(a, b int) (int, error)`
4. implements `safeSqrt(n float64) (float64, error)`
5. prints both success results and inspectable failure details
6. uses defer for small, intentional completion logging rather than as a substitute for errors

## Files

- [main.go](./main.go): complete solution with teaching comments
- [_starter/main.go](./_starter/main.go): starter file with TODOs and requirements

## Run Instructions

Run the completed solution:

```bash
go run ./04-functions-and-errors/8-error-handling
```

Run the starter:

```bash
go run ./04-functions-and-errors/8-error-handling/_starter
```

## Success Criteria

Your finished solution should:

- return custom errors instead of panicking for ordinary math failures
- expose error details through structured fields instead of string matching
- handle both success and failure paths in `main()`
- show at least one small use of defer that keeps the flow clearer
- remain easy to read without overbuilding package structure

## Common Failure Modes

- using `panic` for divide-by-zero or modulo-by-zero
- matching on `err.Error()` instead of inspecting the concrete error value
- hiding too much behavior inside defer
- building one giant `main()` function with no small helpers

## Next Step

After you complete this exercise, move to [FE.10 functional options pattern](../9-functional-options)
for a stretch lesson, or continue to [Section 05](../../05-types-and-interfaces) if you are ready
to move on.
