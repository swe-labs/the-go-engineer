# Types and Design Stage Map

This stage turns basic Go fluency into more intentional modeling and design choices.

It takes learners from structs and interfaces through composition and embedding, then into text
transformation and rendering workflows that depend on those design choices.

## Stage Flow

1. `types-and-interfaces`
   - source: [05-types-and-interfaces](../../../05-types-and-interfaces/)
   - core job: teach structs, methods, interfaces, and small generic helpers
   - milestone: `TI.6` payroll processor project
2. `composition`
   - source: [06-composition](../../../06-composition/)
   - core job: teach explicit reuse, embedding, promoted behavior, and shadowing
   - milestone: `CO.3` bank account project
3. `strings-and-text`
   - source: [07-strings-and-text](../../../07-strings-and-text/)
   - core job: teach parsing, formatting, rendering, and text-driven workflow design
   - milestone: `ST.6` config parser project

## What Each Part Adds

### `types-and-interfaces`

This is where the learner starts modeling data and behavior with clear type boundaries instead of
only writing procedural code.

### `composition`

This is where the learner starts reusing behavior explicitly without sliding into inheritance
thinking.

### `strings-and-text`

This is where the learner starts shaping real inputs and outputs through parsing and rendering
instead of ad hoc text handling.

## Recommended Full-Path Order

1. Finish the Section `05` milestone path first.
2. Move through Section `06` and complete `CO.3`.
3. Move through Section `07` and complete `ST.6`.

## Bridge-Path Reminder

If structs, methods, and interfaces already feel familiar, you can move faster through the early
repetition.
What you should not skip is proof:

- `TI.6`
- `CO.3`
- `ST.6`

## Exit Condition

You are ready for `3 Modules and IO` when you can finish the three stage milestones honestly and
explain how type design, reuse, and text-driven workflows connect together.
