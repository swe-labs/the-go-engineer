# Section 11: Concurrency

## Mission

This section teaches you how Go coordinates work across goroutines, channels, contexts, and
timers without hiding the costs and failure modes behind "just add concurrency" magic.

By the end of the live v2 slice, you should be comfortable:

- starting goroutines deliberately instead of sprinkling `go` blindly
- coordinating work with `sync.WaitGroup` and channels
- applying cancellation and deadlines with `context.Context`
- using timers and tickers without leaking background work
- building small concurrent tools that stay readable under load

Section 11 still contains additional concurrency material beyond this live slice. Those surfaces
remain available as legacy reference lessons while the first milestone path is migrated.

## Beta Stage Ownership

This section belongs to [5 Concurrency System](../docs/stages/05-concurrency-system.md).

Within the beta public shell, it is the first and foundation-heavy half of that stage:

1. Section 11 `concurrency`
2. Section 12 `concurrency-patterns`

That means Section 11 is where learners build the core coordination vocabulary before the stage
moves into higher-pressure patterns in Section 12.

## Who Should Start Here

### Full Path

Start here after completing Section 10 in order.

### Bridge Path

You can move faster if you already understand:

- explicit error handling
- methods, interfaces, and small packages
- I/O surfaces like HTTP requests and file operations
- why cleanup matters around resources and long-running work

Even on the bridge path, do not skip the first lesson in any track.
Those entry points establish the vocabulary the later exercises assume.

### Targeted Path

This section is the second multi-track pilot in v2.
You can choose the track that matches your immediate goal:

- Goroutines track for worker coordination and channel flow
- Context track for cancellation and timeout control
- Time track for timers, tickers, and reminder-style scheduling

## Section Map

| Track | Entry | Milestone | Focus |
| --- | --- | --- | --- |
| Goroutines | [GC.1 goroutines](./goroutines) | `GC.7` | goroutines, WaitGroups, channels, and a bounded downloader |
| Context | [CT.1 background](./context) | `CT.5` | context roots, cancellation, timeouts, and timeout-aware HTTP calls |
| Time & Scheduling | [TM.1 time basics](./time-and-scheduling) | `TM.7` | time values, formatting, timers, tickers, and a console reminder |

## Suggested Order

1. Complete the Goroutines track if you want the strongest concurrency foundations first.
2. Complete the Context track if you work with HTTP, APIs, or long-running I/O.
3. Complete the Time track if you want stronger deadline and scheduling intuition.
4. Use the legacy reference lessons after the milestone path if you want deeper coverage in the same section.

## Section Milestones

This live v2 slice has three milestone surfaces:

- `GC.7` concurrent downloader
- `CT.5` timeout-aware API client
- `TM.7` console reminder

The following lessons remain available as legacy reference surfaces for later alpha work:

- `GC.8` race conditions
- `GC.9` select deep dive
- `GC.10` sync primitives
- `TM.4` random numbers
- `TM.5` scheduler
- `TM.6` timezones

If you can complete the three milestone exercises and explain:

- why goroutine fan-out still needs coordination boundaries
- why request-scoped cancellation should reach the HTTP client layer
- why timers and tickers need explicit lifecycle management

then you are ready to move into the higher-order concurrency patterns in Section 12.

## Pilot Role In V2

This live v2 slice keeps the current `11-concurrency` layout intact while upgrading the learner-facing flow:

- the section now has one top-level guide
- each track has a clearer milestone path
- the milestone exercises have explicit README contracts
- only the first practical slice of each track is promoted into `curriculum.v2.json`

That keeps Section 11 useful now without pretending the entire mega-section is fully migrated in one wave.

## Legacy To Pilot Mapping

- `GC.1` through `GC.7` stay in `11-concurrency/goroutines/*`
- `CT.1` through `CT.5` stay in `11-concurrency/context/*`
- `TM.1`, `TM.2`, `TM.3`, and `TM.7` stay in `11-concurrency/time-and-scheduling/*`
- later goroutine and time lessons remain available, but they are not yet promoted into the live v2 graph

## References

1. [Effective Go: Concurrency](https://go.dev/doc/effective_go#concurrency)
2. [Go Concurrency Patterns: Context](https://go.dev/blog/context)
3. [Package time](https://pkg.go.dev/time)

## Next Step

After you finish the track or milestone you care about here, continue to
[Section 12: Concurrency Patterns](../12-concurrency-patterns).

In the beta shell, that keeps you inside
[5 Concurrency System](../docs/stages/05-concurrency-system.md)
before you move on to
[6 Quality and Performance](../docs/stages/06-quality-and-performance.md).
