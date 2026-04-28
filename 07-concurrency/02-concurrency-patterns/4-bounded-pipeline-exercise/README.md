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


## Mental Model

Think of this as the conceptual blueprint. The components interact by exchanging state, defining clear boundaries between what is requested and what is provided.

## Visual Model

Visualizing this process involves tracing the execution path from the input entry point, through the processing layers, and out to the final output or side effect.

## Machine View

At the hardware level, this translates into specific memory allocations, CPU instruction cycles, and OS-level system calls to manage resources efficiently.

## Solution Walkthrough

The solution demonstrates a complete implementation, proving the concept by bridging the individual requirements into a single, cohesive executable.

## Try It

Run the code locally. Modify the inputs, toggle the conditions, and observe how the output shifts. Experimentation is the fastest way to cement your understanding.

## Verification Surface

The correctness of this component is proven by its associated test suite. We verify boundaries, handle edge cases, and ensure performance constraints are met.

## In Production

Bounded pipelines with early cancellation are the backbone of batch processing in production Go services. Image processors, data transformation jobs, and ETL pipelines all need to limit concurrent work to avoid overwhelming downstream systems (databases, object storage, external APIs). The `errgroup.SetLimit` pattern ensures that a batch of 100,000 items does not launch 100,000 goroutines — it processes them in controlled waves. Early cancellation via `errgroup.WithContext` is equally critical: if item 50 out of 10,000 reveals a corrupt input, there is no value in processing the remaining 9,950 items. The `sync.Pool` pattern prevents allocation pressure that would otherwise trigger frequent garbage collection pauses under high throughput. In production, teams measure the difference: a pipeline that allocates a fresh 64 KB buffer per item generates gigabytes of garbage per minute at scale, while one that reuses pooled buffers keeps GC pause times under a millisecond.

## Thinking Questions

1. What happens to in-flight goroutines when `errgroup.WithContext` cancels the context due to an error? Do they stop immediately or finish their current work?
2. Why must you reset a `bytes.Buffer` obtained from `sync.Pool` before using it, and what bug occurs if you skip the reset?
3. If you set the concurrency limit too low, the pipeline is slow. If you set it too high, you overwhelm downstream systems. How would you determine the right limit for a production workload?
4. How would you modify this pipeline to continue processing remaining items after a failure instead of cancelling, while still collecting and reporting all errors?

## Next Step

After you complete this exercise, continue to [CP.5 URL Health Checker](../5-url-checker-exercise)
or back to the [Stage 07 overview](../README.md).


