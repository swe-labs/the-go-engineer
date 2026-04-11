# Section 07: Strings and Text

## Mission

This section teaches you how Go handles text at the byte, rune, pattern, and rendered-output
levels.

By the end of Section 07, you should be comfortable reading and writing:

- immutable string operations and builder-based output
- formatted text with `fmt`
- Unicode-aware text iteration
- regex-based parsing and extraction
- template-driven text rendering
- small text-processing tools with runnable behavior and tests

## Beta Stage Ownership

This section belongs to [2 Types and Design](../docs/stages/02-types-and-design.md).

Within the beta public shell, it is the third and final part of that stage:

1. Section 05 `types-and-interfaces`
2. Section 06 `composition`
3. Section 07 `strings-and-text`

## Who Should Start Here

### Full Path

Start here after completing Section 06 in order.

### Bridge Path

You can move faster if you already understand:

- structs and methods at a basic level
- loops and conditionals
- maps and slices

Even on the bridge path, do not skip `ST.1` and `ST.3`.
They prevent a lot of subtle text-handling mistakes later.

### Targeted Path

If you are here mainly for parsing and rendering text, review these first:

- `ST.1` strings
- `ST.2` formatting
- `ST.4` regex

## Section Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `ST.1` | Lesson | [strings](./1-strings) | Introduces immutable strings, common helpers, and efficient text building. | entry |
| `ST.2` | Lesson | [formatting](./2-formatting-string) | Shows how `fmt` shapes readable output and reusable string formatting. | `ST.1` |
| `ST.3` | Lesson | [unicode and runes](./3-unicode) | Explains why text iteration by byte and by rune can produce very different results. | `ST.1`, `ST.2` |
| `ST.4` | Lesson | [regex](./4-regex) | Adds pattern-based extraction and parsing for structured text work. | `ST.1`, `ST.3` |
| `ST.5` | Lesson | [text templates](./5-text-template) | Introduces controlled text rendering instead of ad hoc concatenation. | `ST.1`, `ST.2` |
| `ST.6` | Exercise | [config parser project](./6-config-parser) | Combines parsing, regex, and template-driven reporting in one milestone. | `ST.1`, `ST.2`, `ST.4`, `ST.5` |

## Suggested Order

1. Work through `ST.1` to `ST.5` in order.
2. Complete `ST.6` without copying the finished solution line by line.

## Section Milestone

`ST.6` is the current live milestone for this pilot section.

If you can complete it and explain:

- why strings are immutable and why that matters for text-heavy loops
- why regex should be compiled once and reused instead of rebuilt per line
- why templates give you cleaner output generation than manual print sprawl

then you are ready to move into modules and packages in Section 08.

## Pilot Role In V2

This live v2 slice keeps the current Section 07 paths and `ST.*` ids stable while upgrading the
learner-facing structure:

- `ST.1` through `ST.5` are the core lessons
- `ST.6` is the milestone exercise

That keeps the section useful for current learners while the broader v2 migration continues.

## Legacy To Pilot Mapping

This pilot intentionally avoids breaking the current Section 07 filesystem layout.

- `ST.1` through `ST.5` keep their existing lesson directories
- `ST.6` stays at `07-strings-and-text/6-config-parser`, now treated as the section milestone
  surface

## References

1. [Strings, bytes, runes and characters in Go](https://go.dev/blog/strings)
2. [Package text/template](https://pkg.go.dev/text/template)

## Next Step

After `ST.6`, you have completed the core milestone path for
[2 Types and Design](../docs/stages/02-types-and-design.md).

From there, move to [3 Modules and IO](../docs/stages/03-modules-and-io.md).
The first source section in that next stage is
[Section 08: Modules and Packages](../08-modules-and-packages).
