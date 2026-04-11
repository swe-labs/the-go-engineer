# 5 Concurrency System

## Purpose

`5 Concurrency System` teaches how to think about concurrent work without turning programs into
chaos.

## Who This Is For

- learners who already understand application boundaries
- Go developers who want stronger intuition for goroutines, context, scheduling, and bounded work

## Mental Model

Concurrency is about coordination, cancellation, and limits.
The goal is not to add goroutines everywhere.
The goal is to structure concurrent work so it can start, stop, and fail in controlled ways.

## Why This Stage Exists

This stage is where "I can launch goroutines" turns into "I can reason about concurrent systems."

The goal is not to make learners memorize channel tricks.
The goal is to build real judgment about coordination, cancellation, deadlines, bounded work, and
failure propagation.

## What You Should Learn Here

- goroutine fundamentals
- context and cancellation
- timers, deadlines, and scheduling
- worker-pool and fan-out/fan-in style patterns
- bounded failure handling in concurrent systems

## Stage Shape

This stage currently has one foundation-heavy section plus one higher-order pattern section:

1. Section `11` `concurrency`
   - the live public foundation layer for goroutines, context, and time-based coordination
2. Section `12` `concurrency-patterns`
   - the live public pattern layer for `errgroup`, bounded pipelines, and pressure-aware
     coordination

That means the stage is intentionally split into "learn the primitives" and "apply the
production-shaped patterns" instead of pretending those are the same difficulty level.

## Current Source Content

- [11-concurrency/context](../../11-concurrency/context/)
- [11-concurrency/goroutines](../../11-concurrency/goroutines/)
- [11-concurrency/time-and-scheduling](../../11-concurrency/time-and-scheduling/)
- [12-concurrency-patterns](../../12-concurrency-patterns/)

## Stage Support Docs

Use these support docs when you want the beta-stage view without digging through the full Section
`11` and Section `12` inventories:

- [Concurrency System support index](./concurrency-system/README.md)
- [Stage map](./concurrency-system/stage-map.md)
- [Milestone guidance](./concurrency-system/milestone-guidance.md)

## Where This Stage Starts

Start with [Section 11: Concurrency](../../11-concurrency/).

That is the public entry to this stage because it teaches the coordination vocabulary that Section
`12` depends on.
Section `12` stays inside the same stage, but it should come after the foundation tracks.

## Recommended Order

Use this order for the current beta-facing path:

1. complete the Goroutines track from `GC.1` through `GC.7`
2. complete the Context track from `CT.1` through `CT.5`
3. complete the Time and Scheduling track from `TM.1` through `TM.7`
4. move to Section `12` and complete the `CP.1` through `CP.5` path
5. use the remaining concurrency reference surfaces after the live milestones if you want deeper
   coverage

## Path Guidance

### Full Path

Work through Section `11` first, then move into Section `12`.
This stage is designed to build intuition before adding higher-pressure concurrency patterns.

### Bridge Path

You can move faster if goroutines, channels, timeouts, and HTTP cancellation already feel familiar,
but do not skip:

- `GC.1`
- `GC.4`
- `CT.1`
- `CT.3`
- `TM.1`
- `TM.3`
- `CP.1`
- `CP.4`

Those are the main proof surfaces that keep the stage grounded in coordination discipline instead
of concurrency folklore.

### Targeted Path

If you enter this stage with a narrow goal:

- start with the Goroutines track if your gap is basic coordination
- start with the Context track if your gap is cancellation and timeout control
- start with the Time track if your gap is timers, deadlines, and scheduling
- move to Section `12` once you can explain bounded work and cancellation boundaries without hand
  waving

## Stage Milestones

The current live milestone backbone is:

- `GC.7` concurrent downloader
- `CT.5` timeout-aware API client
- `TM.7` console reminder
- `CP.5` URL health checker

`CP.4` remains an important promoted exercise inside this stage because it is the first place where
bounded-concurrency pressure becomes explicit.

## Finish This Stage When

- you can explain when concurrency helps and when it hurts
- you can wire cancellation and timeouts into real flows
- you can build bounded worker or pipeline patterns
- you can spot common goroutine leaks and coordination mistakes

More concretely:

- you can explain why goroutines need ownership and shutdown boundaries
- you can propagate context cancellation through I/O and sibling work
- you can use timers and tickers without leaking background behavior
- you can explain why bounded concurrency is a systems decision, not just a syntax pattern
- you can complete `CP.5` without treating `errgroup` as magic

## Next Stage

Move to [6 Quality and Performance](./06-quality-and-performance.md).
