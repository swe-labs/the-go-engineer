# 02 Language Basics

## Mission

This section teaches the first real Go fundamentals after installation and tooling.

By the end of this section, a learner should be able to:

- Declare and update variables confidently
- Use constants for fixed values
- Model enum-like values with named types and iota
- Combine those pieces into one small, readable program

## Section Map

| ID | Lesson | What It Unlocks |
| --- | --- | --- |
| `LB.1` | [variables](./1-variables) | Introduces typed values, zero values, and declaration styles |
| `LB.2` | [constants](./2-constants) | Separates immutable values from ordinary variables |
| `LB.3` | [enums with iota](./3-enums) | Shows how Go models enum-like values |
| `LB.4` | [application logger](./4-application-logger) | Combines the section into one milestone exercise |

## Zero-Magic Boundary

This section stays at the foundations level.

It does not try to teach:
- Functions as a design tool
- Control flow in depth
- Data structures

It only teaches enough about types to make the learner ready for control flow.

## How To Use This Section

For each lesson:

1. Read the `README.md`
2. Open `main.go`
3. Run the lesson
4. Make one small change
5. Run it again

## Finish This Section When

- Variables feel familiar instead of confusing
- The learner can explain zero values
- The learner can create constants and understand why they matter
- The learner can explain iota in their own words

## Next Step

Continue to [03-control-flow](../03-control-flow/) after this section.