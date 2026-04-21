# CT.5 Timeout-Aware API Client

## Mission

Build a small HTTP client that uses `context.WithTimeout` to enforce deadlines and fails clearly
when a request takes too long.

This exercise is the Context track milestone for Stage 07.

## Prerequisites

Complete these first:

- `CT.1` background
- `CT.2` with cancel
- `CT.3` with timeout
- `CT.4` with value

## What You Will Build

Implement a client that:

1. creates a timeout-bound context for each outbound request
2. attaches that context to the request with `http.NewRequestWithContext`
3. returns a useful wrapped error when the deadline expires
4. demonstrates both a successful request and a deliberately timed-out request

## Files

- [main.go](./main.go): complete solution with teaching comments
- [_starter/main.go](./_starter/main.go): starter file with TODOs and requirements

## Run Instructions

Run the completed solution:

```bash
go run ./07-concurrency/01-concurrency/context/5-timeout-client
```

Run the starter:

```bash
go run ./07-concurrency/01-concurrency/context/5-timeout-client/_starter
```

## Success Criteria

Your finished solution should:

- use `context.WithTimeout` instead of making an unbounded HTTP call
- attach the context to the request itself, not just the outer function
- distinguish timeout failures from other request errors
- keep the timeout behavior visible in runnable output

## Note

The current example makes live HTTP requests, so it expects network access when you run the full
solution.

## Next Step

After you complete this exercise, continue back to the [Context track](../README.md) or the
[Stage 07 overview](../../README.md).
