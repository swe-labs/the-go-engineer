# Quality and Performance Milestone Guidance

## What Counts As Stage Completion

You should be able to prove that code behaves correctly and measure where time or allocations are
actually going.

## Milestones

### `TE.4` benchmarking

This proves that you can compare behavior with evidence instead of intuition alone.
It is the point where testing and performance start to overlap.

### `PR.2` live pprof endpoint

This proves that you can inspect a live process and reason about where work is happening.

## Bridge Path Check

If you are moving quickly through this stage, make sure you can still explain:

- why table-driven tests make behavior easier to compare
- why handler tests are better than full servers for many correctness checks
- why benchmark numbers must be interpreted, not merely copied
- why `pprof` answers different questions than unit tests

## Ready To Move On

Move to [7 Architecture](../07-architecture.md) when you can complete the current path and explain
both the correctness and measurement choices in plain language.
