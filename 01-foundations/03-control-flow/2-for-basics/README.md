# CF.2 For Basics

## Mission

Learn how Go repeats work.

Go has only one loop keyword: `for`.
That is enough for counted loops, condition-based loops, and a small `range` preview.

## Why This Lesson Exists Now

After branching, the next question is:

"How do I do the same kind of work more than once without copy-pasting code?"

That is the loop question.

## Run Instructions

```bash
go run ./01-foundations/03-control-flow/2-for-basics
```

## Code Walkthrough

### `for i := 1; i <= 5; i++ { ... }`

This is the classic counted loop.

It has three parts:

- start value
- condition for continuing
- step after each iteration

The program uses it when it knows the count-driven shape of the work.

### `countdown := 3` and `for countdown > 0 { ... }`

This is the condition-only form.
It keeps running while the condition stays true.

That is why people often say Go's `for` can also behave like a `while` loop.

### `words := []string{"go", "learn", "repeat"}`

This is a small preview of a collection.
Do not worry about slice internals yet.

For this lesson, all you need to know is:

- `words` holds several values
- `range` visits them one by one

The data-structures section teaches the collection mechanics properly.
Here, `range` is only a loop shape preview.

### `for _, word := range words { ... }`

This preview keeps only the current value because the lesson goal is repetition, not indexed data
access.

## Common Mistakes

- off-by-one loop conditions
- forgetting that the loop condition is checked before the next round
- trying to understand all slice mechanics here instead of treating the list as a preview tool

## Try It

1. Change the counted loop to stop at `3` instead of `5`.
2. Set `countdown := 5` and rerun it.
3. Add one more word to the preview list and inspect the `range` output.

## Why This Matters In Real Software

Loops are how programs:

- process records
- retry work
- build summaries
- scan input one piece at a time

## Next Step

Continue to `CF.3` break / continue.
