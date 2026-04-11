# 6 Quality and Performance

## Purpose

`6 Quality and Performance` teaches how to trust and measure engineering work.

## Who This Is For

- learners who can already build meaningful systems
- developers who want stronger testing and profiling instincts

## Mental Model

Correctness and speed are not guesses.
They are proven through tests, benchmarks, profiling, and repeatable checks.

## Why This Stage Exists

This stage is where learners stop treating quality work as cleanup and start treating it as
engineering proof.

The goal is to build the habit of measuring, testing, profiling, and validating before making
claims about correctness or performance.

## What You Should Learn Here

- unit and integration testing
- HTTP client testing
- profiling and performance inspection
- benchmark-driven improvement
- environment-aware test workflows

## Stage Shape

This stage currently has one live public section with two promoted tracks plus additional reference
surfaces:

1. `testing`
   - the live beta path for correctness, table-driven tests, handler tests, and benchmarking
2. `profiling`
   - the live beta path for CPU profiles and `pprof`
3. `http-client-testing` and `4-testcontainers`
   - valuable reference surfaces that stay available while the public beta path is still more
     selective

That means the stage is honest about what learners should complete now while still exposing the
broader Section `13` inventory.

## Current Source Content

- [13-quality-and-performance/testing](../../13-quality-and-performance/testing/)
- [13-quality-and-performance/http-client-testing](../../13-quality-and-performance/http-client-testing/)
- [13-quality-and-performance/profiling](../../13-quality-and-performance/profiling/)
- [13-quality-and-performance/4-testcontainers](../../13-quality-and-performance/4-testcontainers/)

## Stage Support Docs

Use these support docs when you want the beta-stage view without digging through the full Section
`13` inventory:

- [Quality and Performance support index](./quality-and-performance/README.md)
- [Stage map](./quality-and-performance/stage-map.md)
- [Milestone guidance](./quality-and-performance/milestone-guidance.md)

## Where This Stage Starts

Start with [Section 13: Quality and Performance](../../13-quality-and-performance/).

Inside that section, begin with the testing track first.
That keeps the stage grounded in proof of correctness before it moves into measurement and
performance investigation.

## Recommended Order

Use this order for the current beta-facing path:

1. complete the Testing track from `TE.1` through `TE.4`
2. complete the Profiling track from `PR.1` through `PR.2`
3. use `http-client-testing` and `4-testcontainers` as optional reference material while they
   remain outside the live beta backbone

## Path Guidance

### Full Path

Work through the testing path first, then move into profiling.
This keeps the stage ordered around "prove it works" before "measure how it behaves."

### Bridge Path

You can move faster if unit tests, handlers, and profiling vocabulary are already somewhat
familiar, but do not skip:

- `TE.1`
- `TE.3`
- `TE.4`
- `PR.1`
- `PR.2`

Those are the main proof surfaces that show you can both verify behavior and investigate
performance honestly.

### Targeted Path

If you enter this stage with a narrow goal:

- start with the Testing track if your gap is correctness and test design
- start with the Profiling track if your gap is "why is this slow?" investigation
- come back and complete both tracks before treating the stage as finished

## Stage Milestones

The current live milestone backbone is:

- `TE.4` benchmarking
- `PR.2` live pprof endpoint

## Finish This Stage When

- you can design tests that prove behavior instead of only chasing coverage
- you can benchmark and profile before making performance claims
- you can improve code with measurement instead of guesswork
- you understand when full-environment tests are worth the cost

More concretely:

- you can explain why testability starts with design, not with test framework tricks
- you can interpret benchmark results instead of only reading the nanosecond number
- you can use `pprof` to find hot paths before proposing optimizations
- you understand which Section `13` surfaces are the current beta path and which remain reference
  material

## Next Stage

Move to [7 Architecture](./07-architecture.md).
