# TE.4 Benchmarking

## Mission

Use `testing.B` to compare implementation choices and understand how allocations show up in
performance work.

This surface is the testing-track output for Stage 08.

## Files

- [benchmarks_test.go](./benchmarks_test.go): benchmark examples for string building, slice growth, and lookup patterns

## Run Instructions

```bash
go test -bench=. -benchmem ./08-quality-test/01-quality-and-performance/testing/benchmarks
```

## Success Criteria

You should be able to:

- read `ns/op`, `B/op`, and `allocs/op`
- explain why benchmark setup should stay outside the timed region
- compare two approaches without confusing correctness tests with performance tests


## Prerequisites

You should be comfortable with Go syntax, basic data structures, and the control flow mechanics covered in earlier sections.

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

Performance assumptions are dangerous until they are measured. In production, small inefficiencies in a hot path — such as a tight loop that processes millions of events per second — can compound into significant infrastructure costs and latency spikes. `testing.B` is the standard tool Go engineers use to prove that an optimization actually works before merging it. However, micro-benchmarks can be misleading if they do not reflect production data distributions. For example, a lookup function might be fast for 10 items but degrade exponentially at 10,000 items. The `-benchmem` flag is particularly critical because in Go, memory allocation pressure (tracked by `allocs/op`) is often the hidden cause of high CPU usage, as excessive garbage generation forces the garbage collector to consume CPU cycles that should belong to the application. Teams that write benchmarks for critical paths catch performance regressions in CI before they reach production.

## Thinking Questions

1. Why must you call `b.ResetTimer()` if your benchmark requires expensive setup before the actual work begins?
2. If you run a benchmark on your local laptop, why might the results not accurately predict the performance of the code running in a constrained Docker container?
3. How can compiler optimizations like loop unrolling or dead-code elimination make a micro-benchmark look artificially fast, and how do you prevent it?
4. If a function is extremely fast (low `ns/op`) but causes high `allocs/op`, what impact will that have on a long-running server?

## Next Step

After benchmarking, continue to [PR.1 CPU profile](../../profiling) or back to the
[Testing track](../README.md).

