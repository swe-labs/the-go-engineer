# 04 Data Structures

## Mission

This section teaches the data structures that turn basic control flow into problem solving.

By the end of this section, a learner should be able to:

- explain when arrays, slices, maps, and pointers are the right tool
- reason about copy behavior versus shared state
- follow the machine-level consequences of mutation and addressing
- combine slices, maps, and pointers in a small milestone project

## Why This Section Exists Now

The learner already knows values, branching, loops, and cleanup.

That is enough to ask stronger questions:

- how should a program store ordered data?
- when is key-based lookup better than scanning?
- what changes when two variables share access to the same underlying state?

Those are data-structure questions.

## Zero-Magic Boundary

This section teaches:

- arrays
- slices
- maps
- pointers
- slice sharing and mutation

It does **not** jump ahead into:

- package-level architecture
- concurrency safety
- profiling strategy
- advanced memory tuning

Those topics come later, after the learner has earned stronger system context.

## Section Ownership

This section belongs to [02 Language Basics](../README.md).

## Suggested Learning Flow

1. Start with arrays and slices so ordered data feels concrete.
2. Move into maps when lookup by key becomes the better mental model.
3. Use pointers after you can already describe the value being changed.
4. Finish with the contact-manager milestone before moving on.

## Section Milestone

`DS.6` is the milestone for this section.

You are ready for the next section when you can explain:

- why slices are usually the real ordered-data tool in Go
- why maps trade order for lookup speed
- why pointers let you update the original value instead of a copy
- how these concepts combine in one small directory-style program

## Next Step

After `DS.6`, continue to [03 Functions & Errors](../../03-functions-errors).
