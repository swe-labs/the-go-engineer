# Track A: Testing

## Mission

This track teaches you how to prove behavior in Go with small, explicit tests instead of treating
tests as a separate framework or a giant afterthought.

## Track Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `TE.1` | Lesson | [unit testing](./user) | Introduces `testing.T`, basic assertions, and test file structure. | entry |
| `TE.2` | Lesson | [table-driven tests](./user) | Shows the idiomatic Go pattern for structured test cases and sub-tests. | `TE.1` |
| `TE.3` | Lesson | [HTTP handler testing](./user) | Tests handlers with `httptest` instead of a real server. | `TE.1`, `TE.2` |
| `TE.4` | Lesson | [benchmarking](./benchmarks) | Uses `testing.B` and `-benchmem` to compare performance choices. | `TE.1`, `TE.2` |

## Suggested Order

1. Work through `TE.1`, `TE.2`, and `TE.3` in order.
2. Complete `TE.4` once the test-design mindset feels natural.

## Track Milestone

`TE.4` is the current testing-track output.

If you can explain:

- why table-driven tests are the default Go testing pattern
- why `httptest` is better than spinning up a real server for unit-level handler tests
- why `b.ReportAllocs()` matters when reading benchmark output

then the testing part of Stage 08 is doing its job.

## Next Step

After `TE.4`, continue to the [Profiling track](../profiling) or back to the
[Stage 08 overview](../README.md).
