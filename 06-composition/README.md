# Section 06: Composition

## Mission

This section teaches you to build larger Go types by combining smaller ones, then use embedding as
an intentional shortcut instead of pretending Go has inheritance.

By the end of Section 06, you should be comfortable reading and writing:

- named-field composition that keeps ownership explicit
- embedded structs with promoted fields and methods
- shadowed methods that change behavior at the outer type boundary
- small domain models that reuse behavior without class hierarchies

## Beta Stage Ownership

This section belongs to [2 Types and Design](../docs/stages/02-types-and-design.md).

Within the beta public shell, it is the second of three connected parts:

1. Section 05 `types-and-interfaces`
2. Section 06 `composition`
3. Section 07 `strings-and-text`

## Who Should Start Here

### Full Path

Start here after completing Section 05 in order.

### Bridge Path

You can move faster if you already understand:

- structs and methods
- pointer receivers
- interfaces as behavior contracts

Even on the bridge path, do not skip `CO.1`.
It creates the contrast that makes embedding easier to understand.

### Targeted Path

If you are here mainly for embedding and promotion, review these first:

- `TI.1` structs
- `TI.2` methods

## Section Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `CO.1` | Lesson | [composition](./06-composition-and-embedding/1-composition) | Introduces has-a modeling with named fields so reuse stays explicit. | entry |
| `CO.2` | Lesson | [embedding](./06-composition-and-embedding/2-embedding) | Adds promoted fields and methods while clarifying that embedding is still composition. | `CO.1` |
| `CO.3` | Exercise | [bank account project](./06-composition-and-embedding/3-bank-account) | Combines composition, embedding, promoted methods, and method shadowing in one milestone. | `CO.1`, `CO.2` |

## Suggested Order

1. Work through `CO.1` and `CO.2` in order.
2. Complete `CO.3` without copying the finished solution line by line.

## Section Milestone

`CO.3` is the current live milestone for this pilot section.

If you can complete it and explain:

- why composition gives Go explicit reuse without inheritance
- why embedding promotes behavior but does not change the underlying type identity
- why method shadowing should be used deliberately instead of casually

then you are ready to move into strings, formatting, and text processing in Section 07.

## Pilot Role In V2

This live v2 slice keeps the current Section 06 paths and `CO.*` ids stable while upgrading the
learner-facing structure:

- `CO.1` and `CO.2` are the core lessons
- `CO.3` is the milestone exercise

That keeps the section useful for current learners while the broader v2 migration continues.

## Legacy To Pilot Mapping

This pilot intentionally avoids breaking the current Section 06 filesystem layout.

- `CO.1` and `CO.2` stay in `06-composition/06-composition-and-embedding`
- `CO.3` stays at `06-composition/06-composition-and-embedding/3-bank-account`, now treated as
  the section milestone surface

## References

1. [Effective Go: Embedding](https://go.dev/doc/effective_go#embedding)

## Next Step

After `CO.3`, continue to [Section 07: Strings and Text](../07-strings-and-text).
In the beta shell, that is still part of
[2 Types and Design](../docs/stages/02-types-and-design.md).
