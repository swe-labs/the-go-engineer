# 03 Control Flow

## Mission

This section teaches learners how a Go program chooses what to do next and how it repeats work without turning into duplicated code.

By the end of this section, a learner should be able to:

- branch with `if`, `else if`, and `else`
- repeat work with `for`
- choose between multiple cases with `switch`
- stop or skip loop work with `break` and `continue`
- schedule cleanup work with `defer`
- combine those ideas into one small runnable milestone

## Why This Section Exists Now

The learner already knows basic values and expressions.

That is enough to ask the next engineering questions:

- how does a program choose one path over another?
- how does it repeat the same operation without copy-pasting code?
- how does it stop early when the current path is no longer useful?

Those are control-flow questions.

## Zero-Magic Boundary

This section intentionally stays inside:

- values
- comparisons
- counters
- strings
- small list previews used only as loop targets

It does **not** formally teach:

- data-structure design
- helper-function design
- reusable error-return patterns
- package structure

You may see a tiny list preview inside a loop example.
That is a tool, not the topic.
Data structures are taught properly in the next section.

## Section Ownership

This section belongs to [02 Language Basics](../README.md).

## Section Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `CF.1` | Lesson | [if / else](./1-if-else) | Introduces branching and decision-making. | entry |
| `CF.2` | Lesson | [for basics](./2-for-basics) | Teaches Go's only loop keyword and repeated work. | `CF.1` |
| `CF.3` | Lesson | [break / continue](./3-break-continue) | Teaches early exit and selective skipping inside loops. | `CF.2` |
| `CF.4` | Lesson | [switch](./4-switch) | Teaches readable multi-branch decision logic. | `CF.1`, `CF.2`, `CF.3` |
| `CF.5` | Lesson | [defer basics](./5-defer-basics) | Introduces cleanup scheduling and return-path safety. | `CF.4` |
| `CF.6` | Lesson | [defer use cases](./6-defer-use-cases) | Shows file, mutex, and cleanup-shaped defer patterns. | `CF.5` |
| `CF.7` | Exercise | [pricing checkout](./7-pricing-checkout) | Combines branching, looping, switching, and defer in one milestone. | `CF.1` through `CF.6` |

## Suggested Learning Flow

1. Read each lesson `README.md` first.
2. Open `main.go` only after you understand the lesson mission.
3. Run each lesson and compare the output with the explanation.
4. Change one thing at a time using the `Try It` prompts.
5. Attempt the milestone starter before opening the finished solution.

## Section Milestone

`CF.7` is the live milestone for this section.

You are ready for the next section when you can explain:

- when `if` is enough and when `switch` reads better
- why Go only needs `for`
- when `break` stops the whole loop
- when `continue` skips only the current iteration
- why `defer` ensures cleanup regardless of return path
- why the small `range` example stays a preview instead of turning into a data-structures lesson

## Next Step

After `CF.7`, continue to [04 Data Structures](../04-data-structures).
