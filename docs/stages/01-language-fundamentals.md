# 1 Language Fundamentals

## Purpose

`1 Language Fundamentals` turns first contact with Go into actual programming fluency.

## Who This Is For

- beginners who finished `0 Foundation`
- experienced developers who want a disciplined Go fundamentals pass

## Mental Model

This stage teaches the building blocks of code and control.
You are learning how values move, how decisions are made, how collections are shaped, and how
functions and errors create boundaries.

## Why This Stage Exists

This is the early backbone of the whole curriculum.

The goal is not only to "cover the basics."
The goal is to build enough fluency that later stages stop feeling like code-copying and start
feeling like real engineering work.

## What You Should Learn Here

- variables, values, and basic expressions
- branching and loops
- arrays, slices, maps, and pointer-aware mutation
- function design and parameter/result thinking
- errors as explicit control flow instead of hidden failure

## Stage Shape

This stage has four connected parts:

1. `language-basics`
   - variables, constants, named values, and the first milestone exercise
2. `control-flow`
   - loops, branching, and readable decision logic
3. `data-structures`
   - arrays, slices, maps, pointers, and mutation-aware thinking
4. `functions-and-errors`
   - function boundaries, multiple returns, custom errors, defer, and recover

The stage is intentionally ordered from "what is a value?" to "what is a function boundary?".

## Current Source Content

- [01-core-foundations/language-basics](../../01-core-foundations/language-basics/)
- [01-foundations/03-control-flow](../../01-foundations/03-control-flow/)
- [03-data-structures](../../03-data-structures/)
- [04-functions-and-errors](../../04-functions-and-errors/)

## Stage Support Docs

Use these support docs when you want the beta-stage view without digging through four section
READMEs:

- [Language Fundamentals support index](./language-fundamentals/README.md)
- [Stage map](./language-fundamentals/stage-map.md)
- [Milestone guidance](./language-fundamentals/milestone-guidance.md)

## Where This Stage Starts

This stage starts at [01-core-foundations/language-basics](../../01-core-foundations/language-basics/).

The `getting-started` track from Section `01` belongs to
[0 Foundation](./00-foundation.md), not to this stage.

Practical split:

- `GT.1` to `GT.4` = `0 Foundation`
- `LB.1` to `LB.4` = start of `1 Language Fundamentals`

## Recommended Order

Use this order across the source sections:

1. [01-core-foundations/language-basics](../../01-core-foundations/language-basics/)
2. [01-foundations/03-control-flow](../../01-foundations/03-control-flow/)
3. [03-data-structures](../../03-data-structures/)
4. [04-functions-and-errors](../../04-functions-and-errors/)

Do not skip ahead from `language-basics` into functions and errors too early.
Control flow and data structures make the later function and error lessons much easier to
understand.

## Path Guidance

### Full Path

Complete the four stage parts in order and do the milestone exercise in each section.

### Bridge Path

You can move faster through the opening repetition if Go syntax feels familiar already, but do not
skip:

- `LB.4`
- `CF.5`
- `DS.6`
- `FE.9`

Those are the main proof surfaces that show you can actually use the stage, not just read it.

### Targeted Path

If you enter this stage late, choose the section that matches your actual gap and review the
earlier milestones honestly before claiming you can skip them.

Examples:

- weak on branching or loops: start with `control-flow`
- weak on collections and mutation: start with `data-structures`
- weak on idiomatic Go failure handling: start with `functions-and-errors`

## Stage Milestones

The current milestone backbone is:

- `LB.4` application logger
- `CF.5` pricing checkout
- `DS.6` contact manager
- `FE.9` error handling project

If you can complete those four milestone surfaces honestly, you have the practical foundation for
the next beta stage.

## Finish This Stage When

- you can read and write small Go programs without copying line by line
- you understand how control flow and data structures interact
- you can design simple functions with clear error behavior
- you can complete beginner milestone exercises without guessing

More concretely:

- you can explain the difference between values, collections, and function boundaries
- you can reason about loops, maps, slices, and pointers without panic
- you understand why Go treats errors as values instead of exceptions
- you can complete the stage milestones without relying on line-by-line imitation

## Next Stage

Move to [2 Types and Design](./02-types-and-design.md).
