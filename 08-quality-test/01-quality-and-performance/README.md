# Stage 08: Quality and Performance

## Mission

This section teaches how Go teams prove correctness and investigate performance without treating testing and profiling as optional cleanup work.

By the end of this section, you should be comfortable:

- writing small unit tests and table-driven tests with confidence
- testing handlers and output-producing code without starting real servers
- reading and writing benchmarks that reveal allocation and throughput differences
- capturing CPU profiles and exposing live pprof endpoints safely
- explaining why correctness work and performance work belong in the same engineering toolbox

## Stage Ownership

This section belongs to [08 Quality & Testing](../README.md).

## Section Map

| Track | Entry | Milestone | Focus |
| --- | --- | --- | --- |
| Testing | [TE.1 unit testing](./testing) | `TE.10` | unit tests, sub-tests, fuzzing, seams, integration tests, and golden files |
| Profiling | [PR.1 CPU profile](./profiling) | `PR.6` | CPU, memory, escape analysis, benchmarks, and memory layout |

## Suggested Order

1. Work through the Testing track first.
2. Move into the Profiling track once you can explain what "correct but slow" looks like.
3. Use the additional reference surfaces when you want deeper testing variations.

## Next Step

After Stage 08, continue to [09 Architecture & Security](../../09-architecture).
