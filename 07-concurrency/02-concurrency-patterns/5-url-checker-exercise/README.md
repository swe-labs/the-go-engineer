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

Health checking is a fundamental production operations pattern. Every load balancer, service mesh, and container orchestrator uses health checks to decide whether to route traffic to a given instance. The concurrent fan-out pattern this exercise teaches — bounded HTTP checks with result aggregation — is exactly how production monitoring systems like Prometheus blackbox exporter, uptime robots, and internal SLA monitors work. Sorting results by latency is not cosmetic: it reveals which endpoints are degrading before they fail completely, giving operations teams early warning. In production, health checkers run continuously on intervals and feed alerting systems. The `sync.Pool` for HTTP clients matters at scale because creating a new `http.Client` per check means establishing a new TCP connection (and TLS handshake) every time, adding hundreds of milliseconds of latency and wasting system resources. Reusing clients allows connection pooling via `http.Transport`, which keeps established connections warm.

## Thinking Questions

1. Why does this exercise use HTTP HEAD requests instead of GET requests for health checking, and what information do you lose by using HEAD?
2. If one of the checked URLs has a DNS resolution failure that takes 30 seconds to timeout, how does that affect the other concurrent checks?
3. How would you modify this tool to run as a continuous background service that checks URLs every 60 seconds and alerts when latency exceeds a threshold?
4. What is the difference between reusing an `http.Client` via `sync.Pool` and simply sharing a single global `http.Client` across all goroutines?

## Next Step

After you complete this exercise, continue to the [Stage 07 overview](../README.md) or move to
[Stage 08: Quality and Performance](../../../08-quality-test/01-quality-and-performance).


