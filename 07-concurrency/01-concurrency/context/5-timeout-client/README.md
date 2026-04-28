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

Every outbound HTTP call in a production service must have a timeout. Without one, a single slow upstream dependency can exhaust all available goroutines and connection pool slots, causing the entire service to hang — a failure mode known as cascading failure or "grey failure." The `context.WithTimeout` pattern this exercise teaches is the standard way Go services enforce request deadlines. In production, timeouts are typically configured per-endpoint rather than globally, because a health check that should complete in 100ms has very different timeout requirements than a batch data export that legitimately takes 30 seconds. Teams that do not distinguish between timeout errors and other request failures end up retrying requests that will never succeed, amplifying load on an already-struggling upstream. The pattern of attaching the context to the request itself — not just checking it in the outer function — ensures that the HTTP client, DNS resolver, TLS handshake, and response body read all respect the same deadline.

## Thinking Questions

1. Why should the context be attached to the HTTP request with `http.NewRequestWithContext` instead of checking `ctx.Done()` in a separate goroutine?
2. If an upstream service is consistently slow, should your client retry with the same timeout, retry with a longer timeout, or stop retrying entirely? What factors determine the right strategy?
3. What happens to the TCP connection when a context timeout fires while the server is still writing the response body?
4. How would you set different timeouts for different API endpoints in the same service without duplicating client code?

## Next Step

After you complete this exercise, continue back to the [Context track](../README.md) or the
[Stage 07 overview](../../README.md).


