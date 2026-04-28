# PR.2 Live pprof Endpoint

## Mission

Expose live profiling endpoints on an internal-only server so you can inspect CPU, heap, mutex,
and goroutine behavior while a service is under load.

This surface is the profiling-track output for Stage 08.

## Prerequisites

Complete these first:

- `PR.1` CPU profile

## Files

- [main.go](./main.go): runnable pprof demo server with a separate public API mux and admin port

## Run Instructions

```bash
go run ./08-quality-test/01-quality-and-performance/profiling/3-http-pprof
```

Then inspect the running service with commands such as:

```bash
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=5
go tool pprof http://localhost:6060/debug/pprof/heap
```

## Success Criteria

You should be able to:

- explain why pprof handlers are registered through the blank import
- describe the two-port pattern for keeping pprof off the public API
- use CPU and heap endpoints to inspect a live service safely


## Mental Model

Think of this as the conceptual blueprint. The components interact by exchanging state, defining clear boundaries between what is requested and what is provided.

## Visual Model

Visualizing this process involves tracing the execution path from the input entry point, through the processing layers, and out to the final output or side effect.

## Machine View

At the hardware level, this translates into specific memory allocations, CPU instruction cycles, and OS-level system calls to manage resources efficiently.

## Code Walkthrough

We step through the code sequentially, examining how the interfaces are satisfied, where the errors are checked, and how the core loop manages control flow.

## Try It

Run the code locally. Modify the inputs, toggle the conditions, and observe how the output shifts. Experimentation is the fastest way to cement your understanding.

## In Production

Continuous profiling is how mature engineering organizations debug performance degradation in real-world environments. Synthetic load tests rarely capture the exact behavior of production traffic, so the ability to pull a CPU or heap profile from a live, struggling service without stopping it is invaluable. However, exposing `net/http/pprof` directly on a public-facing API port is a critical security vulnerability: it leaks internal execution state, memory footprints, and potentially sensitive variables to the internet. The two-port pattern shown in this exercise — binding the main application to one port (e.g., 8080) and the diagnostic/admin handlers to a separate, internal-only port (e.g., 6060) — is the industry standard mitigation. This ensures profiling is accessible only to operators or internal metric scrapers, never to external users.

## Thinking Questions

1. Why does importing `net/http/pprof` with the blank identifier (`_`) automatically register handlers, and why is this considered a risky pattern in library code?
2. If you notice a memory leak in production, which pprof profile (heap, goroutine, allocs) would you check first, and what specifically would you look for?
3. What is the performance overhead of leaving the HTTP pprof endpoints active in production continuously?
4. How would you secure a pprof endpoint if you could not use a separate port and had to serve it on the main public router?

## Next Step

After `PR.2`, continue to the [Stage 08 overview](../README.md) or move to
[Stage 09: Application Architecture](../../../../09-architecture).

