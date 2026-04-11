# Section 13: Quality and Performance

## Mission

This section teaches you how Go teams prove correctness and investigate performance without treating
testing and profiling as optional cleanup work.

By the end of the live v2 slice, you should be comfortable:

- writing small unit tests and table-driven tests with confidence
- testing handlers and output-producing code without starting real servers
- reading and writing benchmarks that reveal allocation and throughput differences
- capturing CPU profiles and exposing live pprof endpoints safely
- explaining why correctness work and performance work belong in the same engineering toolbox

Section 13 still contains additional reference material beyond this live slice. Those surfaces
remain available while the first public testing and profiling tracks are migrated.

## Beta Stage Ownership

This section belongs to [6 Quality and Performance](../docs/stages/06-quality-and-performance.md).

Within the beta public shell, it is the complete live source section for that stage:

- the Testing track is the primary correctness path
- the Profiling track is the primary performance investigation path
- the other Section `13` surfaces remain reference material while the public beta route stays
  selective

## Who Should Start Here

### Full Path

Start here after completing Section 12 in order.

### Bridge Path

You can move faster if you already understand:

- packages, interfaces, and dependency injection
- HTTP handlers and basic I/O surfaces
- why allocations and cleanup matter in long-running programs

Even on the bridge path, do not skip `TE.1` or `PR.1`.
They establish the basic testing and profiling habits the rest of the slice depends on.

### Targeted Path

This section is a multi-track v2 slice.
You can choose the track that matches your immediate goal:

- Testing track for correctness, design for testability, and benchmarks
- Profiling track for CPU profiles, live pprof endpoints, and performance investigation

## Section Map

| Track | Entry | Milestone | Focus |
| --- | --- | --- | --- |
| Testing | [TE.1 unit testing](./testing) | `TE.4` | unit tests, table-driven tests, handler tests, and benchmarking |
| Profiling | [PR.1 CPU profile](./profiling) | `PR.2` | offline CPU profiles and live pprof inspection |

## Suggested Order

1. Work through the Testing track first if you want stronger correctness habits.
2. Move into the Profiling track once you can explain what "correct but slow" looks like.
3. Use the legacy reference surfaces after the live slice if you need deeper specialty testing patterns.

## Section Milestones

This live v2 slice has two promoted outputs:

- `TE.4` benchmarking
- `PR.2` live pprof endpoint

The following surfaces remain available as legacy reference material for later alpha work:

- `http-client-testing`
- `4-testcontainers`

If you can complete the live slice and explain:

- why testability starts with code design, not with a mock library
- why benchmarks must be interpreted with allocations and setup costs in mind
- why pprof is the tool for "where is the time really going?" questions

then you are ready to move into the application architecture work in Section 14.

## Pilot Role In V2

This live v2 slice keeps the current `13-quality-and-performance` layout intact while promoting
the main public testing and profiling path:

- `TE.1` through `TE.4` form the live testing track
- `PR.1` and `PR.2` form the live profiling track
- `http-client-testing` and `4-testcontainers` remain legacy reference surfaces

## References

1. [Package testing](https://pkg.go.dev/testing)
2. [Go Blog: Profiling Go Programs](https://go.dev/blog/pprof)
3. [Package net/http/pprof](https://pkg.go.dev/net/http/pprof)

## Next Step

After you finish the track or milestone you care about here, continue to
[Section 14: Application Architecture](../14-application-architecture).

In the beta shell, that means you are moving from
[6 Quality and Performance](../docs/stages/06-quality-and-performance.md)
to
[7 Architecture](../docs/stages/07-architecture.md).
