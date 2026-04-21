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

## Next Step

After `PR.2`, continue to the [Stage 08 overview](../README.md) or move to
[Stage 09: Application Architecture](../../../../09-architecture).
