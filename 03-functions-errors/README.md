# 03 Functions & Errors

## Mission

This stage teaches how to model behavior boundaries and failure states. It moves beyond basic control flow to treat both functions and errors as explicit, manageable values.

By the end of this stage, a learner should be able to:

- define function signatures with multiple returns
- handle errors explicitly as values instead of exceptions
- use defer for safe resource cleanup
- implement first-class functions and closures
- understand panic and recover boundaries

## Stage Map

| Track | Surface | Core Job |
| --- | --- | --- |
| `FE.1-6` | Core Functions | Teach basic signatures, multiple returns, named returns, and defer. |
| `FE.7` | Order Summary | Capstone exercise combining function basics. |
| `FE.8-9` | Advanced Functions | Teach first-class functions, callbacks, and closure state capture. |
| `FE.10` | Panic & Recover | Teach the final boundary for unrecoverable errors. |

## Why This Stage Exists Now

The learner already knows:

- variables, types, and basic operators
- control flow (if, for, switch)
- basic data structures (arrays, slices, maps)

That is enough to start asking engineering questions like:

- how do we group logic into reusable blocks?
- how do we signal failure back to the caller safely?
- how do we ensure resources are closed when a function exits?

## Suggested Learning Flow

1. Follow the numeric sequence from `FE.1` to `FE.6`.
2. Complete the `FE.7` exercise to prove your understanding.
3. Finish the advanced topics in `FE.8` to `FE.10`.

## Next Step

After this section, continue to [04 Types & Design](../04-types-design).
