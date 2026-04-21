# Stage 07: Concurrency Patterns

## Mission

This track teaches the first production-shaped concurrency patterns that sit on top of the goroutine and context fundamentals from Stage 07.

By the end of this track, you should be comfortable:

- replacing bare `sync.WaitGroup` coordination with `errgroup` when work can fail
- using `errgroup.WithContext` to cancel sibling work on the first failure
- using `sync.Pool` deliberately when short-lived allocations become a real hotspot
- building bounded concurrent pipelines instead of launching unbounded background work
- combining concurrency control and operational safety in small real exercises

## Stage Ownership

This track belongs to [07 Concurrency](../README.md).

## Track Map

| Surface | Entry | Milestone | Focus |
| --- | --- | --- | --- |
| errgroup path | [CP.1 errgroup](./1-errgroup) | `CP.5` | errgroup, cancellation, sync.Pool, and bounded concurrency |
| Worker pool reference | [6 worker pool](./6-worker-pool) | `6-worker-pool` | worker-pool design and shutdown boundaries |

## Suggested Order

1. Work through `CP.1`, `CP.2`, and `CP.3` in order.
2. Complete `CP.4` as the first bounded-concurrency exercise.
3. Complete `CP.5` as the capstone exercise for this pattern track.
4. Use `6-worker-pool` as supporting reference material after the canonical path.

## Next Step

After `CP.5`, return to [07 Concurrency](../README.md) if you want to review the stage, or continue to [08 Quality & Testing](../../08-quality-test).
