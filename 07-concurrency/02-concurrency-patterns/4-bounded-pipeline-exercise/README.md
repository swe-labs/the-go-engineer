# CP.4 Bounded Pipeline Exercise

## Mission

Build a bounded concurrent pipeline that stops on the first failure and reuses large temporary
buffers instead of allocating a fresh one for every work item.

This exercise is the first promoted exercise surface for Stage 07.

## Prerequisites

Complete these first:

- `CP.1` errgroup basics
- `CP.2` errgroup with context cancellation
- `CP.3` sync.Pool

## What You Will Build

Implement an image-processing style batch job that:

1. launches one unit of work per image ID
2. caps concurrency with `g.SetLimit(4)`
3. cancels the remaining work on the first failure
4. reuses large `bytes.Buffer` instances through a `sync.Pool`

## Files

- [main.go](./main.go): complete solution with teaching comments
- [_starter/main.go](./_starter/main.go): starter file with TODOs and requirements

## Run Instructions

Run the completed solution:

```bash
go run ./07-concurrency/02-concurrency-patterns/4-bounded-pipeline-exercise
```

Run the starter:

```bash
go run ./07-concurrency/02-concurrency-patterns/4-bounded-pipeline-exercise/_starter
```

## Success Criteria

Your finished solution should:

- use `errgroup.WithContext` instead of raw goroutine launch loops
- bound concurrency with `SetLimit`
- stop sibling work quickly when one image fails
- return pooled buffers cleanly without leaking stale data

## Next Step

After you complete this exercise, continue to [CP.5 URL Health Checker](../5-url-checker-exercise)
or back to the [Stage 07 overview](../README.md).

