# CF.3 Break / Continue

## Mission

Learn how to control a loop after it has already started.

Sometimes a loop should stop completely.
Sometimes it should skip only the current item.

That is what `break` and `continue` are for.

## Why This Lesson Exists Now

Once a learner can write a loop, the next question is:

"What if I do not want every iteration to behave the same way?"

That is the purpose of loop control.

## Run Instructions

```bash
go run ./01-foundations/03-control-flow/3-break-continue
```

## Code Walkthrough

### `for i := 1; i <= 10; i++ { ... }`

The loop counts from `1` to `10`.
Inside the loop, we add two control rules.

### `if i%2 == 0 { continue }`

`continue` skips the rest of the current iteration.

That means:

- even numbers are seen
- but they are not printed as part of the main result

The loop then moves forward to the next value.

### `if i == 7 { break }`

`break` stops the loop completely.

Once this line runs, the program does not continue to `8`, `9`, or `10`.

## Common Mistakes

- using `break` when you only meant to skip one value
- using `continue` and then wondering why later code inside the same loop never runs
- putting loop-control checks in an order that hides the intended behavior

## Try It

1. Move the `break` check before the `continue` check and compare the output.
2. Change the stop number from `7` to `9`.
3. Remove `continue` and rerun the lesson.

## Why This Matters In Real Software

Loop control appears in real programs when you need to:

- stop once the target is found
- skip invalid records
- stop early when further work is pointless

## Next Step

Continue to `CF.4` switch.
