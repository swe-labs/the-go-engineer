# 2 Types and Design

## Purpose

`2 Types and Design` moves learners from syntax into modeling.

## Who This Is For

- learners who already understand basic Go flow and data handling
- developers who want better instincts for structs, interfaces, composition, and data modeling

## Mental Model

This stage is about choosing shapes, not just writing lines.
Good Go design comes from representing data clearly and choosing boundaries that stay simple under
change.

## Why This Stage Exists

This is where the curriculum stops being only about writing correct code and starts being about
choosing better structures.

The goal is to help learners move from "I can make this compile" to "I can model this problem
cleanly and explain why this shape fits."

## What You Should Learn Here

- structs, methods, and interfaces
- composition over inheritance-style thinking
- text and data transformation
- type choices that support maintainability
- design trade-offs in small systems

## Stage Shape

This stage has three connected parts:

1. `types-and-interfaces`
   - structs, methods, interfaces, and the first serious domain-modeling milestone
2. `composition`
   - explicit reuse, embedding, promoted behavior, and method shadowing
3. `strings-and-text`
   - text transformation, parsing, rendering, and data-shaping through string-heavy workflows

The stage is intentionally ordered from "how do I model a type?" to "how do I compose and
transform real data with those types?"

## Current Source Content

- [05-types-and-interfaces](../../05-types-and-interfaces/)
- [06-composition](../../06-composition/)
- [07-strings-and-text](../../07-strings-and-text/)

## Stage Support Docs

Use these support docs when you want the beta-stage view without digging through three section
READMEs:

- [Types and Design support index](./types-and-design/README.md)
- [Stage map](./types-and-design/stage-map.md)
- [Milestone guidance](./types-and-design/milestone-guidance.md)

## Where This Stage Starts

This stage starts at [05-types-and-interfaces](../../05-types-and-interfaces/).

If `1 Language Fundamentals` gave you working fluency, this stage teaches you how to turn that
fluency into cleaner program structure.

## Recommended Order

Use this order across the source sections:

1. [05-types-and-interfaces](../../05-types-and-interfaces/)
2. [06-composition](../../06-composition/)
3. [07-strings-and-text](../../07-strings-and-text/)

Do not skip directly from structs into text processing too early.
Composition and embedding make the modeling decisions in the final part easier to understand.

## Path Guidance

### Full Path

Complete the three stage parts in order and do the milestone exercise in each section.

### Bridge Path

You can move faster if structs, methods, and interfaces already feel familiar, but do not skip:

- `TI.6`
- `CO.3`
- `ST.6`

Those are the main proof surfaces that show you can actually model and shape data, not just read
the lessons.

### Targeted Path

If you enter this stage late, start with the part that matches your actual design gap and then
check the earlier milestones honestly.

Examples:

- weak on types, methods, and contracts: start with `types-and-interfaces`
- weak on reuse and embedding: start with `composition`
- weak on parsing, formatting, and rendering text data: start with `strings-and-text`

## Stage Milestones

The current milestone backbone is:

- `TI.6` payroll processor project
- `CO.3` bank account project
- `ST.6` config parser project

If you can complete those three milestone surfaces honestly, you have the practical foundation for
the next beta stage.

## Finish This Stage When

- you can model a small domain with structs and interfaces honestly
- you know when composition helps more than clever abstraction
- you can transform and validate text or record-style data cleanly
- your type choices feel intentional instead of accidental

More concretely:

- you can explain why a struct, method, or interface exists in a design
- you can use composition without slipping into inheritance-shaped thinking
- you can parse and render small text-driven workflows without ad hoc sprawl
- you can complete the stage milestones without guessing at the type boundaries

## Next Stage

Move to [3 Modules and IO](./03-modules-and-io.md).
