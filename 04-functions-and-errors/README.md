# Section 04: Functions and Errors

## Mission

This section teaches you to turn functions into explicit contracts and failures into explicit
values.

By the end of Section 04, you should be comfortable reading and writing:

- small function signatures
- closures and variadic helpers
- multi-return functions
- custom and wrapped errors
- defer-based cleanup
- panic and recover boundaries

## Who Should Start Here

### Full Path

Start here after completing Sections 01-03 in order.

### Bridge Path

You can move faster if you already understand:

- basic Go syntax
- control flow
- slices and maps
- pointers at a basic level

Even on the bridge path, do not skip the error-handling lessons.
They are the hinge of the section.

### Targeted Path

If you are here only for error handling, review these first:

- `FE.1` functions
- `FE.3` variadic functions
- `FE.4` multiple return values

## Section Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `FE.1` | Lesson | [functions](./1-function) | Builds the basic function contract model: inputs, outputs, and readable behavior. | entry |
| `FE.2` | Lesson | [closures and recursion](./2-function-2) | Shows how functions can capture state and call themselves when the problem shape needs it. | `FE.1` |
| `FE.3` | Lesson | [variadic functions](./3-variadic-func) | Introduces flexible call shapes without hiding what the function really accepts. | `FE.1`, `FE.2` |
| `FE.4` | Lesson | [multiple return values](./4-function-multi-values) | Introduces the `(value, error)` style that drives the rest of the section. | `FE.1`, `FE.3` |
| `FE.5` | Lesson | [custom errors](./5-custom-error) | Shows how to attach stable meaning to failures instead of relying on fragile strings. | `FE.4` |
| `FE.6` | Lesson | [error wrapping](./5b-error-wrapping) | Adds context to low-level failures while preserving inspectable error identity. | `FE.5` |
| `FE.7` | Lesson | [defer](./6-defer) | Teaches disciplined cleanup and "always run this at the end" behavior. | `FE.1`, `FE.5` |
| `FE.8` | Lesson | [panic and recover](./7-panic-recover) | Shows where panic belongs, where it does not, and how recover works at a boundary. | `FE.7` |
| `FE.9` | Exercise | [error handling project](./8-error-handling) | Combines custom errors, inspectable failures, and deferred cleanup into one milestone exercise. | `FE.1`, `FE.4`, `FE.5`, `FE.6`, `FE.7`, `FE.8` |
| `FE.10` | Stretch Lesson | [functional options pattern](./9-functional-options) | A useful configuration pattern once the section fundamentals feel solid. | `FE.1`, `FE.2` |

## Suggested Order

1. Work through `FE.1` to `FE.8` in order.
2. Complete `FE.9` without copying the finished solution line by line.
3. Use `FE.10` as a stretch lesson before moving to the next section.

## Section Milestone

`FE.9` is the current live milestone for this pilot section.

If you can complete it and explain:

- why custom errors beat string matching
- why wrapped errors keep context useful
- why defer helps cleanup but does not replace ordinary error returns

then you are ready to move into more structured type and interface design in Section 05.

## Pilot Role In V2

This first live v2 pilot keeps the current Section 04 paths and `FE.*` ids stable while adding a
clearer learner-facing structure:

- `FE.1` through `FE.8` are the core lessons
- `FE.9` is the milestone exercise
- `FE.10` is an optional stretch pattern lesson

That keeps the section honest for current learners while the wider v2 migration grows around it.

## References

1. [Error handling and Go](https://go.dev/blog/error-handling-and-go)
2. [Defer, Panic, and Recover](https://go.dev/blog/defer-panic-and-recover)

## Next Step

After `FE.9`, continue to [Section 05: Types and Interfaces](../05-types-and-interfaces).
If you want one more pattern lesson before moving on, finish `FE.10` first.
