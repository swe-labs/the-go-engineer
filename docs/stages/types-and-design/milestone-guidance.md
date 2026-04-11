# Types and Design Milestone Guidance

The `2 Types and Design` stage has three proof surfaces.

These milestones show whether a learner can model and shape data intentionally instead of only
copying patterns.

## Milestone Backbone

| ID | Surface | What It Proves |
| --- | --- | --- |
| `TI.6` | payroll processor project | you can model domain data with structs, methods, interfaces, and a small generic helper |
| `CO.3` | bank account project | you can use composition and embedding deliberately without confusing them for inheritance |
| `ST.6` | config parser project | you can parse, transform, and render text-driven workflows cleanly |

## How To Use These Milestones

### Full Path

Complete all three in order.
They are the strongest proof that the stage really stuck.

### Bridge Path

If you move faster through the lessons, use these milestones as the non-skippable proof surfaces.
Do not claim the stage is done if you only skimmed the lessons and skipped the milestone work.

### Targeted Path

If you enter the stage late, choose the milestone that best matches your real gap first, then
check backward honestly:

- weak on structs, methods, interfaces, or generics: start with `TI.6`
- weak on reuse and embedding: start with `CO.3`
- weak on parsing, formatting, or rendering text: start with `ST.6`

## Honest Readiness Signals

You are likely ready to leave this stage when:

- you can complete the milestone without copying the solution line by line
- you can explain why the design is shaped that way
- you can change a type, boundary, or output shape without the whole program collapsing
- you can name the trade-off you made instead of only saying that the code works

## If A Milestone Feels Too Hard

That is usually a routing signal, not a failure signal.

Go back one layer:

- `TI.6` too hard: revisit structs, methods, interfaces, and generic helpers
- `CO.3` too hard: revisit named-field composition and embedding basics
- `ST.6` too hard: revisit strings, formatting, regex, and templates
