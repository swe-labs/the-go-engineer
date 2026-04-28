# TM.7 Console Reminder

## Mission

Build a small reminder app that counts down with a ticker and fires a one-shot reminder with
`time.AfterFunc`.

This exercise is the Time and Scheduling track milestone for Stage 07.

## Prerequisites

Complete these first:

- `TM.1` time basics
- `TM.2` formatting
- `TM.3` timers and tickers

## What You Will Build

Implement a reminder that:

1. reads a duration and message from the command line
2. schedules the reminder with `time.AfterFunc`
3. displays a countdown with `time.NewTicker`
4. stops the ticker cleanly when the reminder fires
5. exits without leaving background timer resources running

## Files

- [main.go](./main.go): complete solution with teaching comments
- [_starter/main.go](./_starter/main.go): starter file with TODOs and requirements

## Run Instructions

Run the completed solution:

```bash
go run ./07-concurrency/01-concurrency/time-and-scheduling/7-reminder 5 "Take a break!"
```

Run the starter:

```bash
go run ./07-concurrency/01-concurrency/time-and-scheduling/7-reminder/_starter 5 "Take a break!"
```

## Success Criteria

Your finished solution should:

- convert the CLI seconds argument into a `time.Duration`
- use `time.AfterFunc` for the final reminder
- use `time.NewTicker` for the countdown loop
- stop the ticker cleanly to avoid leaks
- keep the countdown and reminder messages readable


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

Timer and ticker management is a frequent source of goroutine leaks in production Go services. Every `time.NewTicker` that is not stopped continues to fire on its internal goroutine forever, consuming memory and CPU even after the code that created it has moved on. In production systems, tickers drive health check intervals, metric flush cycles, cache eviction sweeps, and retry backoff loops. The pattern of stopping the ticker cleanly — which this exercise enforces — is critical because a service that leaks one ticker per request will eventually exhaust memory after enough traffic. `time.AfterFunc` introduces a different risk: the callback runs on its own goroutine, so any shared state it accesses must be synchronized. Production schedulers (cron jobs, task queues, SLA monitors) all build on these same primitives. Teams that understand timer lifecycle management avoid the class of bugs where a service works perfectly in development but degrades over hours in production as leaked timers accumulate.

## Thinking Questions

1. What happens if you forget to call `ticker.Stop()` before the function returns? How would you detect this leak in a running service?
2. Why does `time.AfterFunc` run its callback on a separate goroutine, and what synchronization issues does this create if the callback modifies shared state?
3. If you need a ticker that fires "at least once per second" but the handler sometimes takes longer than a second, what behavior does `time.NewTicker` produce?
4. How would you implement a ticker that adjusts its interval dynamically based on system load?

## Next Step

After you complete this exercise, continue back to the [Time and Scheduling track](../README.md)
or the [Stage 07 overview](../../README.md).

