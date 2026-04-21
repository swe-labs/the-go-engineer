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

## Next Step

After benchmarking, continue to [PR.1 CPU profile](../../profiling) or back to the
[Testing track](../README.md).
