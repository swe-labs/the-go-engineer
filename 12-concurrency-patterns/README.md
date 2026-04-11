# Section 12: Concurrency Patterns

## Mission

This section teaches you the first production-shaped concurrency patterns that sit on top of the
goroutine and context fundamentals from Section 11.

By the end of the live v2 slice, you should be comfortable:

- replacing bare `sync.WaitGroup` coordination with `errgroup` when work can fail
- using `errgroup.WithContext` to cancel sibling work on the first failure
- using `sync.Pool` deliberately when short-lived allocations become a real hotspot
- building bounded concurrent pipelines instead of launching unbounded background work
- combining concurrency control and operational safety in small real exercises

## Beta Stage Ownership

This section belongs to [5 Concurrency System](../docs/stages/05-concurrency-system.md).

Within the beta public shell, it is the second and pattern-heavy half of that stage:

1. Section 11 `concurrency`
2. Section 12 `concurrency-patterns`

That means this section should be read as the place where concurrency primitives become
production-shaped patterns instead of as a completely separate topic.

## Who Should Start Here

### Full Path

Start here after completing Section 11 in order.

### Bridge Path

You can move faster if you already understand:

- goroutines, channels, and WaitGroups
- context cancellation and timeouts
- why temporary allocations can pressure the garbage collector

Even on the bridge path, do not skip `CP.1` or `CP.2`.
They establish the control-flow model the exercises rely on.

## Current Section Map

| Surface | Status | Entry | Milestone | Focus |
| --- | --- | --- | --- | --- |
| errgroup path | Live v2 slice | `CP.1` | `CP.5` | errgroup, cancellation, sync.Pool, and bounded concurrency |
| Worker pool reference | Legacy reference | `6-worker-pool/` | `6-worker-pool/` | robust worker-pool design and shutdown boundaries |

## Suggested Order

1. Work through `CP.1`, `CP.2`, and `CP.3` in order.
2. Complete `CP.4` as the first bounded-concurrency exercise.
3. Complete `CP.5` as the live section capstone exercise.
4. Use `6-worker-pool` as legacy reference material after the live slice.

## Section Milestones

This live v2 slice has two promoted exercise surfaces:

- `CP.4` bounded pipeline exercise
- `CP.5` URL health checker

If you can complete them and explain:

- why `errgroup` is safer than a bare WaitGroup when goroutines can fail
- why `errgroup.WithContext` is the clean stop signal for sibling work
- why `sync.Pool` only makes sense after you can explain the allocation problem it solves
- why bounded concurrency is a systems-design choice, not just a syntax trick

then you are ready to move into testing, quality, and performance work in Section 13.

## Pilot Role In V2

This live v2 slice keeps the current `12-concurrency-patterns` layout intact while promoting the
main errgroup and pool path into the public curriculum graph:

- `CP.1` through `CP.3` are the core lessons
- `CP.4` and `CP.5` are the live exercises
- `6-worker-pool` remains a legacy reference surface for later alpha work

## References

1. [golang.org/x/sync/errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup)
2. [Package sync](https://pkg.go.dev/sync)
3. [Go Blog: Profiling Go Programs](https://go.dev/blog/pprof)

## Next Step

After `CP.5`, you have completed the current live milestone path for
[5 Concurrency System](../docs/stages/05-concurrency-system.md).

From there, move to [6 Quality and Performance](../docs/stages/06-quality-and-performance.md).
The first source section in that next stage is
[Section 13: Quality and Performance](../13-quality-and-performance).
