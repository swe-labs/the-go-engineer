# Module 04 Checkpoint — Functions, Errors, and Data Semantics

## Purpose

This checkpoint verifies that you understand Go's function mechanics, error handling philosophy, and data semantics before advancing to types, interfaces, and packages. You should be able to trace function calls through the stack, explain value vs pointer semantics, construct and inspect error chains, and use `defer` correctly.

## Scope

This checkpoint covers:

- Function declaration, parameters, and return values
- Multiple return values and named returns
- Variadic functions
- First-class functions and closures
- Value semantics vs pointer semantics for function arguments
- Pointer receivers and value mutation
- The `error` interface and `errors.New`
- Error wrapping with `fmt.Errorf` and `%w`
- `errors.Is` and `errors.As` for error chain inspection
- Custom error types
- `errors.Join` for multi-error aggregation
- `defer` execution order and mechanics
- `panic` and `recover` behavior
- The `init` function
- Data semantics: when copies happen and when they do not

## Time expectation

Plan for 45 to 60 minutes.

## Allowed resources

You may use:

- Module 04 lesson READMEs
- Your terminal and editor
- Your own notes

You may not copy from the answer key before attempting the questions.

## Evidence required

Submit:

1. Answers to `questions.md`
2. A code trace: given a function with `defer` and multiple return paths, predict the defer execution order before running it, then verify
3. A short paragraph explaining why Go does not have exceptions and what problems `(result, error)` solves

## Rubric

| Dimension | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------------|----------|----------------|
| **Function mechanics** | Explains stack frames, value copying on call, and closure variable capture correctly | Knows function syntax but uncertain about closures or variadic behavior | Cannot explain parameter passing |
| **Error handling** | Constructs and traverses error chains with `%w`, `errors.Is`, `errors.As`; explains why wrapping matters | Uses basic error handling but cannot trace wrapped chains | Returns strings instead of error types; ignores errors |
| **Defer and panic** | Predicts LIFO defer order correctly; explains panic stack unwinding and recover mechanics | Knows defer runs after return but confused about order with multiple defers | Cannot explain when defers execute |
| **Data semantics** | Distinguishes value vs pointer passing; explains when mutations are visible to callers | Knows pointers exist but cannot predict mutation behavior | Confuses value and pointer semantics |

**Passing threshold**: 3/5 average across all dimensions

## Retake policy

If you miss a category, return to the relevant lesson, revise your answer, and write one sentence explaining what changed in your understanding.

## NEXT UP

[Module 05 — Types, Interfaces, Packages, and Modules](../../../05-types-interfaces-packages-modules/README.md)
