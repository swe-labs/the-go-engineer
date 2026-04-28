# GS.3 Shutdown Capstone

## Mission

Coordinate readiness, HTTP draining, worker draining, and final cleanup in one production-style
shutdown flow.

This surface is the graceful-shutdown track output for Stage 10.

## Files

- [main.go](./main.go): production-style shutdown capstone with readiness and worker coordination

## Run Instructions

```bash
go run ./10-production/02-graceful-shutdown/3-capstone
```

Then press `Ctrl+C` or send `SIGTERM` and watch the shutdown sequence.

## Success Criteria

You should be able to:

- explain why readiness must flip before draining begins
- describe the shutdown dependency order across HTTP, workers, and shared resources
- show how `errgroup` helps coordinate a multi-component shutdown


## 



## 



## 



## 



## 



## Try It

Run the code locally. Modify the inputs, toggle the conditions, and observe how the output shifts. Experimentation is the fastest way to cement your understanding.

## 




## Prerequisites

You should be comfortable with Go syntax, basic data structures, and the control flow mechanics covered in earlier sections.

## Mental Model

Think of this as the conceptual blueprint. The components interact by exchanging state, defining clear boundaries between what is requested and what is provided.

## Visual Model

Visualizing this process involves tracing the execution path from the input entry point, through the processing layers, and out to the final output or side effect.

## Machine View

At the hardware level, this translates into specific memory allocations, CPU instruction cycles, and OS-level system calls to manage resources efficiently.

## Solution Walkthrough

The solution demonstrates a complete implementation, proving the concept by bridging the individual requirements into a single, cohesive executable.

## Verification Surface

The correctness of this component is proven by its associated test suite. We verify boundaries, handle edge cases, and ensure performance constraints are met.

## In Production

Graceful shutdown is what separates a toy server from a production service. When you deploy a new version of a service, Kubernetes (or your orchestrator) sends a `SIGTERM` signal to the old instances. If those instances exit immediately, any in-flight HTTP requests or background jobs are brutally severed, causing 502 Bad Gateway errors for users. The coordinated shutdown flow this exercise teaches — failing readiness probes, allowing the load balancer to route new traffic elsewhere, draining existing HTTP connections, finishing background jobs, and closing database connections — ensures zero-downtime deployments. The `errgroup` pattern provides the necessary concurrency control to shut down multiple independent components simultaneously while waiting for all of them to finish cleanly. Teams that master graceful shutdown can deploy during peak traffic hours without users ever noticing.

## Thinking Questions

1. Why is it critical to fail the readiness probe and wait for the load balancer to notice *before* you start draining HTTP connections?
2. If an HTTP request takes 5 minutes to process, but your orchestrator only gives the container 30 seconds to shut down before sending `SIGKILL`, what happens to that request?
3. How does `http.Server.Shutdown(ctx)` know which connections to close immediately and which ones to wait for?
4. If the database connection is closed before the HTTP server finishes draining, what will happen to the remaining in-flight requests?

## Next Step

After `GS.3`, continue to the [Code Generation track](../../06-code-generation) or back to the
[Graceful Shutdown track](../README.md).


