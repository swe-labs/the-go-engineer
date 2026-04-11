# Section 05: Types and Interfaces

## Mission

This section teaches you to model data with structs, attach behavior with methods, and define
behavior contracts with interfaces instead of inheritance.

By the end of Section 05, you should be comfortable reading and writing:

- structs that model real domain data
- methods with clear receiver intent
- interfaces that describe behavior boundaries
- small generic helpers built on interface-aware constraints
- polymorphic code that works across multiple concrete types

## Beta Stage Ownership

This section belongs to [2 Types and Design](../docs/stages/02-types-and-design.md).

Within the beta public shell, it is the first of three connected parts:

1. Section 05 `types-and-interfaces`
2. Section 06 `composition`
3. Section 07 `strings-and-text`

## Who Should Start Here

### Full Path

Start here after completing Section 04 in order.

### Bridge Path

You can move faster if you already understand:

- functions and multiple return values
- explicit error handling
- basic pointer usage
- slices and maps

Even on the bridge path, do not skip `TI.3` through `TI.6`.
Those lessons carry the main engineering value of the section.

### Targeted Path

If you are here mainly for interfaces and polymorphism, review these first:

- `TI.1` structs
- `TI.2` methods

## Section Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `TI.1` | Lesson | [structs](./1-struct) | Introduces Go's main data-modeling tool without hidden class behavior. | entry |
| `TI.2` | Lesson | [methods](./2-methods) | Shows how behavior attaches to types and why receiver choice matters. | `TI.1` |
| `TI.3` | Lesson | [interfaces](./3-interfaces) | Teaches implicit contracts and polymorphism without inheritance. | `TI.1`, `TI.2` |
| `TI.4` | Lesson | [Stringer](./4-stringer) | Connects custom types to a standard-library interface in a practical way. | `TI.2`, `TI.3` |
| `TI.5` | Lesson | [generics](./5-generics) | Introduces reusable, type-safe helper patterns built on constraints. | `TI.3`, `TI.4` |
| `TI.6` | Exercise | [payroll processor project](./6-payroll-processor) | Combines structs, methods, interfaces, and a small generic helper into one milestone. | `TI.1`, `TI.2`, `TI.3`, `TI.4`, `TI.5` |
| `TI.7` | Stretch Lesson | [advanced generics](./7-advanced-generics) | Pushes generic collection helpers further once the core section feels solid. | `TI.5` |

## Suggested Order

1. Work through `TI.1` to `TI.5` in order.
2. Complete `TI.6` without copying the finished solution line by line.
3. Use `TI.7` as a stretch lesson before moving to the next section.

## Section Milestone

`TI.6` is the current live milestone for this pilot section.

If you can complete it and explain:

- why structs and methods replace class-style design in Go
- why interfaces describe behavior better than inheritance hierarchies
- why a small generic helper can improve reuse without making the API harder to read

then you are ready to move into composition and embedding in Section 06.

## Pilot Role In V2

This live v2 slice keeps the current Section 05 paths and `TI.*` ids stable while upgrading the
learner-facing structure:

- `TI.1` through `TI.5` are the core lessons
- `TI.6` is the milestone exercise
- `TI.7` is an optional stretch lesson

That keeps the section useful for current learners while the broader v2 migration continues.

## Legacy To Pilot Mapping

This pilot intentionally avoids breaking the current Section 05 filesystem layout.

- `TI.1` through `TI.5` keep their existing lesson directories
- `TI.6` stays at `05-types-and-interfaces/6-payroll-processor`, now treated as the section
  milestone surface
- `TI.7` stays at `05-types-and-interfaces/7-advanced-generics` as an optional stretch lesson

## References

1. [Interfaces and other types](https://go.dev/doc/effective_go#interfaces_and_types)
2. [An Introduction to Generics](https://go.dev/blog/intro-generics)

## Next Step

After `TI.6`, continue to [Section 06: Composition](../06-composition).
In the beta shell, that keeps you inside
[2 Types and Design](../docs/stages/02-types-and-design.md).

If you want one more stretch lesson before moving on, finish `TI.7` first.
