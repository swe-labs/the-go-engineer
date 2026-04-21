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

## Next Step

After you complete this exercise, continue back to the [Goroutines track](../README.md) or the
[Stage 07 overview](../../README.md).
