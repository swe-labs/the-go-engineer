# GC.7 Concurrent Downloader

## Mission

Build a downloader that launches work concurrently, limits the number of active downloads, and
reports results without sharing mutable state between workers.

This exercise is the Goroutines track milestone for Stage 07.

## Prerequisites

Complete these first:

- `GC.1` goroutines
- `GC.2` WaitGroups
- `GC.3` channels
- `GC.4` buffered channels
- `GC.5` closing channels
- `GC.6` pipeline project

## What You Will Build

Implement a downloader that:

1. launches one goroutine per target URL
2. uses a `sync.WaitGroup` to know when all workers are done
3. limits active downloads with a buffered semaphore channel
4. sends success and failure results back through one result channel
5. writes downloaded files to disk and prints a summary at the end

## Files

- [main.go](./main.go): complete solution with teaching comments
- [_starter/main.go](./_starter/main.go): starter file with TODOs and requirements

## Run Instructions

Run the completed solution:

```bash
go run ./07-concurrency/01-concurrency/goroutines/7-downloader
```

Run the starter:

```bash
go run ./07-concurrency/01-concurrency/goroutines/7-downloader/_starter
```

## Success Criteria

Your finished solution should:

- coordinate all downloads with a `sync.WaitGroup`
- limit concurrency with a channel-based semaphore
- aggregate results through a channel instead of shared global state
- clean up partial files on failed downloads
- report both successes and failures clearly

## Note

The current example uses real HTTP downloads, so it expects network access when you run the full
solution.



























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

Bounded concurrency is one of the most important production patterns in Go. Without a semaphore limiting active downloads, a service that receives a burst of 10,000 URLs to fetch would launch 10,000 goroutines simultaneously, exhausting file descriptors, overwhelming the network stack, and likely getting rate-limited or blocked by upstream servers. The semaphore channel pattern this exercise teaches — `sem := make(chan struct{}, maxConcurrency)` — is the idiomatic Go approach used in production crawlers, asset pipelines, and batch processing systems. The result channel pattern is equally critical: sending results through a channel instead of writing to a shared slice eliminates data races that only manifest under production load. Real download systems also need to handle partial file cleanup (as this exercise requires), retry logic with exponential backoff, content-length validation to detect truncated downloads, and context cancellation so that a shutdown signal stops all in-flight downloads instead of leaving orphaned goroutines writing to disk.

## Thinking Questions

1. What happens if the semaphore channel capacity is set to 1? What about equal to the number of URLs?
2. Why is it safer to send results through a channel than to append to a shared slice protected by a mutex?
3. If a download fails halfway through writing a file, what state does the filesystem end up in, and how does cleanup prevent corruption?
4. How would you add a global timeout that cancels all remaining downloads if the total operation exceeds a deadline?

## Next Step

After you complete this exercise, continue back to the [Goroutines track](../README.md) or the
[Stage 07 overview](../../README.md).


