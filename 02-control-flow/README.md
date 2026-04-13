# Section 02: Control Flow

## Migration Note

The canonical rebuilt learner-facing implementation for this section now lives at:

[01-foundations/03-control-flow](../01-foundations/03-control-flow/)

This `02-control-flow` folder is retained temporarily as a legacy source surface while the new
top-level architecture is rolled out.

## Mission

This section teaches learners how Go expresses decisions and repetition without adding extra
syntax they do not need.

By the end of Section 02, you should be comfortable with:

- `for` as Go's only loop keyword
- `if` for ordinary branching and guard clauses
- `switch` for readable multi-branch logic
- combining loops, branching, map lookups, and small helpers in one exercise

## Beta Stage Ownership

This section belongs to [1 Language Fundamentals](../docs/stages/01-language-fundamentals.md).

Within the beta public shell, it is the second of four connected parts:

1. Section 01 `language-basics`
2. Section 02 `control-flow`
3. Section 03 `data-structures`
4. Section 04 `functions-and-errors`

## Who Should Start Here

### Full Path

Start here after completing Section 01.

### Bridge Path

If you already understand variables, constants, and enum-style values in Go, you can begin at
`CF.1` directly.

### Targeted Path

If you only want the live milestone, review these first:

- `CF.1` for loop
- `CF.2` if / else
- `CF.3` switch

## Section Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `CF.1` | Lesson | [for loop](./1-for-loop) | Introduces the only loop keyword Go uses and the different shapes it can take. | entry |
| `CF.2` | Lesson | [if / else](./2-if-else) | Teaches ordinary branching, early decisions, and the comma-ok pattern for safe checks. | `CF.1` |
| `CF.3` | Lesson | [switch](./3-switch) | Shows cleaner multi-branch logic, tagless switches, and type switches. | `CF.1`, `CF.2` |
| `CF.4` | Exercise | [pricing calculator](./4-pricing-calculator) | Combines loops, branching, map lookups, suffix handling, and subtotal calculation in one runnable milestone. | `CF.1`, `CF.2`, `CF.3` |

## Suggested Order

1. Work through `CF.1`, `CF.2`, and `CF.3` in order.
2. Complete `CF.4` without copying the solution line by line.

## Section Milestone

`CF.4` is the current live milestone for this section.

If you can complete it and explain:

- why Go only needs `for`
- why the comma-ok pattern is safer than assuming a map lookup succeeded
- why `switch` often reads more clearly than long `if / else if` chains

then you are ready to move into data structures in Section 03.

## References

1. [Effective Go: Control Structures](https://go.dev/doc/effective_go#control-structures)
2. [A Tour of Go: Flow Control Statements](https://go.dev/tour/flowcontrol/1)

## Next Step

After `CF.4`, continue to Section 03: Data Structures.
In the beta shell, that keeps you inside
[1 Language Fundamentals](../docs/stages/01-language-fundamentals.md).
