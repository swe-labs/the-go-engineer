# Module 03 Checkpoint — Programming Fundamentals with Go

## Purpose

This checkpoint verifies that you have internalized Go's core programming constructs before advancing to functions, errors, and data semantics. The goal is not memorization — the goal is mechanical reasoning: you can predict what Go code does before running it, explain why zero values behave as they do, and trace control flow through loops and conditionals.

## Scope

This checkpoint covers:

- Values, expressions, and operator precedence
- Variable declaration styles (`:=`, `var`, `var` with type)
- Go's type system: numeric types, strings, bools, and type conversions
- Zero values and their practical implications
- Constants, `iota`, and typed enumerations
- Boolean logic and short-circuit evaluation
- Control flow: `if/else`, `switch`, `for` loops, and `range`
- Arrays vs slices: length, capacity, backing arrays, and `append` growth
- Slice sharing, aliasing, and the copy-on-grow boundary
- Maps, the comma-ok idiom, and iteration order
- Strings, bytes, runes, and Unicode handling
- Pointers: addressing, dereferencing, and nil pointer safety
- Common beginner pointer mistakes

## Time expectation

Plan for 45 to 60 minutes.

## Allowed resources

You may use:

- Module 03 lesson READMEs
- Your terminal and editor
- Your own notes

You may not copy from the answer key before attempting the questions.

## Evidence required

Submit:

1. Answers to `questions.md`
2. A short code trace: given a provided slice operation, predict the output before running it, then verify by running it and explaining any difference
3. A brief paragraph explaining why Go distinguishes between arrays and slices

## Rubric

| Dimension | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------------|----------|----------------|
| **Type system understanding** | Correctly explains zero values, type conversions, and why `int` and `int64` are distinct types | Knows basic types but uncertain about conversions or zero values | Cannot explain type distinctions |
| **Control flow reasoning** | Traces loops, range, and switch statements accurately; predicts output before running | Mostly correct with minor loop-counting errors | Cannot trace a basic for loop |
| **Slice mechanics** | Explains len/cap, backing array sharing, and append growth behavior | Understands basic slice usage but not sharing or capacity | Confuses arrays and slices |
| **Map and string handling** | Uses comma-ok idiom correctly; explains rune vs byte distinction | Basic map usage without edge-case awareness | Cannot explain map iteration order or rune concept |
| **Pointer fundamentals** | Explains addressing, dereferencing, and nil pointer panics | Knows pointer syntax but cannot trace mutations | Cannot explain why pointers exist |

**Passing threshold**: 3/5 average across all dimensions

## Retake policy

If you miss a category, return to the relevant lesson, revise your answer, and write one sentence explaining what changed in your understanding.

## NEXT UP

[Module 04 — Functions, Errors, and Data Semantics](../../../04-functions-errors-data-semantics/README.md)
