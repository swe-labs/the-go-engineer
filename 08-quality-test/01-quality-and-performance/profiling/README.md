# Track B: Profiling

## Mission

This track teaches you how to investigate CPU and runtime behavior with Go's built-in profiling
tooling instead of guessing from intuition alone.

## Track Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `PR.1` | Lesson | [CPU profile](./1-cpu-profile) | Introduces offline CPU profiling and reading pprof output. | entry |
| `PR.2` | Lesson | [live pprof endpoint](./3-http-pprof) | Exposes live profiling on an internal port for production-like inspection. | `PR.1` |

## Suggested Order

1. Start with `PR.1` to understand what a profile file actually contains.
2. Move to `PR.2` once you can explain why live profiling must stay off public ports.

## Track Milestone

`PR.2` is the current profiling-track output.

If you can explain:

- why pprof answers different questions than a benchmark
- why `flat` and `cum` time tell different stories
- why live pprof endpoints should live on an internal-only port

then the profiling part of Stage 08 is doing its job.

## Next Step

After `PR.2`, continue back to the [Stage 08 overview](../README.md) or move on to
[Stage 09: Application Architecture](../../../09-architecture).
