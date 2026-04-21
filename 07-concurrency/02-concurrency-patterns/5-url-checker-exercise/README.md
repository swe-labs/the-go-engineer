# CP.5 URL Health Checker

## Mission

Build a concurrent URL checker that caps request fan-out, keeps request latency visible, and
reuses HTTP clients safely.

This exercise is the live capstone surface for Stage 07.

## Prerequisites

Complete these first:

- `CP.1` errgroup basics
- `CP.2` errgroup with context cancellation
- `CP.3` sync.Pool
- `CP.4` bounded pipeline exercise

## What You Will Build

Implement a health checker that:

1. issues concurrent HTTP HEAD requests with `errgroup.WithContext`
2. caps active checks with `SetLimit`
3. collects the result for each URL, including status and latency
4. sorts results by latency for a stable final report
5. reuses HTTP clients through a `sync.Pool`

## Files

- [main.go](./main.go): complete solution with teaching comments
- [_starter/main.go](./_starter/main.go): starter file with TODOs and requirements

## Run Instructions

Run the completed solution:

```bash
go run ./07-concurrency/02-concurrency-patterns/5-url-checker-exercise
```

Run the starter:

```bash
go run ./07-concurrency/02-concurrency-patterns/5-url-checker-exercise/_starter
```

## Success Criteria

Your finished solution should:

- bound concurrent HTTP checks instead of launching them without limits
- keep request construction context-aware
- sort output by latency instead of emitting random goroutine order
- reuse pooled clients safely without carrying stale state between checks

## Note

The current example uses live HTTP endpoints, so it expects network access when you run the full
solution.

## Next Step

After you complete this exercise, continue to the [Stage 07 overview](../README.md) or move to
[Stage 08: Quality and Performance](../../../08-quality-test/01-quality-and-performance).

