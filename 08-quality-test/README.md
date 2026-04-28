# 08 Quality & Testing

## Mission

This stage extends both testing and profiling so learners can prove behavior and inspect runtime cost with more than one tool. It elevates testing from basic assertions to comprehensive verification and introduces performance measurement.

By the end of this stage, a learner should be able to:

- write table-driven unit tests, use mocks, and leverage testing seams
- test HTTP handlers and use golden files for output verification
- measure CPU, memory, and allocs through profiling and benchmarking
- understand escape analysis and memory layout
- prove that a program is both correct and reasonably efficient

## Stage Map

| Track | Surface | Core Job |
| --- | --- | --- |
| `TE.1-10` | Testing | Cover unit tests, fuzzing, seams, mocks, integration tests, and golden files. |
| `PR.1-6` | Profiling & Perf | Cover CPU, memory, escape analysis, benchmark discipline, and memory layout. |

## Why This Stage Exists Now

The learner already knows:

- how to build concurrent data pipelines
- how to interact with the network and database
- how to write basic Go code

That is enough to start asking engineering questions like:

- how do we prove this complex concurrent logic is actually correct?
- how do we catch regressions automatically?
- why is the API running slowly, and how do we find the bottleneck?

## Suggested Learning Flow

1. Complete the `TE` track to build a foundation in correctness and testing discipline.
2. Move to the `PR` track to understand runtime cost, profiling, and benchmarking.
3. The stage ends with proof surfaces for both correctness and cost-awareness.

## Next Step

After this section, continue to [09 Architecture & Security](../09-architecture).
