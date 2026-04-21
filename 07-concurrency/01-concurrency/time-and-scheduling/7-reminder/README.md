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

## Next Step

After you complete this exercise, continue back to the [Time and Scheduling track](../README.md)
or the [Stage 07 overview](../../README.md).
