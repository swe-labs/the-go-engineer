# Language Fundamentals Milestone Guidance

The `1 Language Fundamentals` stage has four proof surfaces.

These are not just exercises to check off.
Together, they show whether a learner can actually use the stage instead of only reading it.

## Milestone Backbone

| ID | Surface | What It Proves |
| --- | --- | --- |
| `LB.4` | application logger | you can write a small Go program with named values, simple structure, and readable output |
| `CF.5` | pricing checkout | you can combine branching, looping, switch-based rules, and safe skipping in one runnable flow |
| `DS.6` | contact manager | you can use slices, maps, and pointers together without losing track of mutation |
| `FE.9` | error handling project | you can model failures explicitly and use functions and cleanup in a disciplined way |

## How To Use These Milestones

### Full Path

Complete all four in order.
They are the strongest proof that the stage really stuck.

### Bridge Path

If you move faster through the lessons, use these milestones as the non-skippable proof surfaces.
Do not claim the stage is done if you only skimmed the lessons and skipped the milestone work.

### Targeted Path

If you enter the stage late, choose the milestone that best matches your real gap first, then
check backward honestly:

- weak on basic Go program shape: start with `LB.4`
- weak on decision logic and loops: start with `CF.5`
- weak on collections and mutation: start with `DS.6`
- weak on function contracts and failures: start with `FE.9`

## Honest Readiness Signals

You are likely ready to leave this stage when:

- you can complete the milestone without copying the solution line by line
- you can explain why the solution is shaped that way
- you can make a small variation without the whole program falling apart
- you can describe your own mistakes clearly instead of guessing

## If A Milestone Feels Too Hard

That is usually a routing signal, not a failure signal.

Go back one layer:

- `LB.4` too hard: revisit `language-basics`
- `CF.5` too hard: revisit `if`, `for`, `break / continue`, and `switch`
- `DS.6` too hard: revisit slices, maps, and pointers
- `FE.9` too hard: revisit multiple returns, custom errors, wrapping, and `defer`
