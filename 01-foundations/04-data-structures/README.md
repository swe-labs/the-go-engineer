# 04 Data Structures

## Mission

This section teaches learners how Go stores and moves collections of data in memory, then shows how
that affects mutation and performance.

By the end of this section, you should be comfortable with:

- arrays as fixed-size value types
- slices as Go's main dynamic collection
- maps for keyed lookup
- pointers for mutation and nil-aware references
- combining those pieces in one small in-memory program without leaning on later-section abstractions

## Zero-Magic Rule

This section intentionally stops before later-section abstractions like helper-function design,
struct-heavy modeling, methods, and package layering.

That means the section milestone should prove slices, maps, and pointers directly rather than
smuggling in ideas from later sections.

## Beta Stage Ownership

This section belongs to [02 Language Basics](../../docs/stages/02-language-basics.md).

Inside the new repo architecture, it is the fourth foundations section:

1. `01-getting-started`
2. `02-language-basics`
3. `03-control-flow`
4. `04-data-structures`
5. `05-functions-and-errors`

## Who Should Start Here

### Full Path

Start here after completing [03-control-flow](../03-control-flow/).

### Bridge Path

If loops, branching, and multiple return values already feel comfortable, begin at `DS.1`.

### Targeted Path

If you only want the live milestone, review these first:

- `DS.1` arrays
- `DS.2` slices
- `DS.3` maps
- `DS.4` pointers
- `DS.5` slice sharing and capacity

## Section Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `DS.1` | Lesson | [arrays](./1-array) | Introduces fixed-size collections and value-copy behavior. | entry |
| `DS.2` | Lesson | [slices](./2-slices) | Teaches Go's primary dynamic collection type and the length/capacity model. | `DS.1` |
| `DS.3` | Lesson | [maps](./3-maps) | Introduces keyed lookup and the comma-ok pattern for presence checks. | `DS.2` |
| `DS.4` | Lesson | [pointers](./4-pointers) | Shows how Go models mutation and shared access without pointer arithmetic. | `DS.1`, `DS.2` |
| `DS.5` | Lesson | [slice sharing and capacity](./5-slices-2) | Explains shared backing arrays and the mutation traps that come with sub-slices. | `DS.2`, `DS.4` |
| `DS.6` | Exercise | [contact directory](./6-contact-manager) | Combines slices, maps, and pointers in one runnable milestone without jumping ahead to later abstractions. | `DS.1`, `DS.2`, `DS.3`, `DS.4`, `DS.5` |

## Suggested Order

1. Work through `DS.1` to `DS.5` in order.
2. For each lesson, read the lesson `README.md` first and then run `main.go`.
3. Complete `DS.6` without copying the finished solution line by line.

## Section Milestone

`DS.6` is the current live milestone for this section.

If you can complete it and explain:

- why slices and maps serve different jobs
- why pointers matter when updates must stick
- why shared backing arrays can surprise you when working with sub-slices
- why the milestone avoids helper-function and struct-heavy design on purpose at this stage

then you are ready to move into `05-functions-and-errors`.

## References

1. [Go Blog: Go Slices - usage and internals](https://go.dev/blog/slices-intro)
2. [Go Blog: Go Maps in Action](https://go.dev/blog/maps)
3. [Effective Go: Data](https://go.dev/doc/effective_go#data)

## Next Step

After `DS.6`, continue to [05-functions-and-errors](../05-functions-and-errors/).
That keeps the new `01-foundations` sequence intact while the broader beta shell still groups both
sections under [02 Language Basics](../../docs/stages/02-language-basics.md).
