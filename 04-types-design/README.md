# 04 Types & Design

## Mission

This stage teaches how to model the physical world into Go's type system using structs, interfaces, composition, and text representations.

By the end of this stage, a learner should be able to:

- define and instantiate structs for state holding
- model behavior constraints with interfaces
- compose types without relying on class inheritance
- manipulate and format text effectively
- design type-safe systems that are easy to test

## Stage Map

| Track | Surface | Core Job |
| --- | --- | --- |
| `TI.1-10` | Types & Interfaces | Build the core type-modeling path from basic structs to the payroll processor project. |
| `TI.11-15` | Advanced Types | Cover the stretch path: empty interface, assertions, type switches, and constraints. |
| `CO.1-3` | Composition | Teach embedding and interface composition without classic OOP inheritance. |
| `ST.1-6` | Strings & Text | Teach string internals, builders, parsing, formatting, and the config parser project. |

## Why This Stage Exists Now

The learner already knows:

- functions as behavior boundaries
- basic data structures (slices, maps)
- explicit error handling

That is enough to start asking engineering questions like:

- how do we group state into meaningful domain objects?
- how do we write functions that accept multiple types of structs?
- how do we share behavior without deep class hierarchies?

## Suggested Learning Flow

1. Complete `TI.1` to `TI.10` to build the core type-modeling path.
2. Complete `TI.11` to `TI.15` to understand advanced interface mechanics.
3. Finish the `CO` (composition) and `ST` (strings/text) tracks to round out design skills.

## Next Step

After this section, continue to [05 Packages & IO](../05-packages-io).
