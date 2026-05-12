# Track A: Testing

## Mission

Master the tools and patterns required to build a reliable Go codebase. This track moves from basic unit tests to advanced patterns like fuzzing and golden files, ensuring your systems are resilient to edge cases and regressions.

## Track Map

| ID | Topic | Surface | Why It Matters |
| --- | --- | --- | --- |
| `TE.1` | **Unit Testing** | [`./01-unit-testing`](./01-unit-testing) | Basic `testing.T` and test file structure. |
| `TE.2` | **Table-Driven Tests** | [`./02-table-driven-tests`](./02-table-driven-tests) | The idiomatic Go pattern for scaling test cases. |
| `TE.3` | **HTTP Testing** | [`./03-http-handler-testing`](./03-http-handler-testing) | Testing handlers with `httptest`. |
| `TE.4` | **Benchmarking** | [`./04-benchmarks`](./04-benchmarks) | Measuring `ns/op` and `allocs/op`. |
| `TE.5` | **Sub-tests & Cleanup** | [`./05-sub-tests-and-cleanup`](./05-sub-tests-and-cleanup) | Granular reporting and reliable teardown. |
| `TE.6` | **Fuzz Testing** | [`./06-fuzz-testing`](./06-fuzz-testing) | Automating edge-case discovery. |
| `TE.7` | **Interfaces for Testability**| [`./07-interfaces-for-testability`](./07-interfaces-for-testability) | Creating seams for isolation. |
| `TE.8` | **Mocking Patterns** | [`./08-mocking-with-interfaces`](./08-mocking-with-interfaces) | Controlling dependency behavior. |
| `TE.9` | **Integration Testing** | [`./09-integration-tests`](./09-integration-tests) | Verifying real component boundaries. |
| `TE.10`| **Golden Files** | [`./10-golden-files`](./10-golden-files) | Managing large test outputs. |

## Suggested Order

1. **Foundations**: `TE.1` -> `TE.3` (Basic logic and HTTP).
2. **Performance**: `TE.4` (Benchmarks).
3. **Advanced Flow**: `TE.5` -> `TE.6` (Lifecycle and Fuzzing).
4. **Architecture**: `TE.7` -> `TE.8` (Decoupling and Mocking).
5. **System Quality**: `TE.9` -> `TE.10` (Integration and Goldens).

## Track Milestone

You have mastered this track when you can:
- Explain why table-driven tests are preferred over multiple test functions.
- Use `httptest.NewRecorder` to verify handler outputs.
- Write a Fuzz test that finds a crashing input.
- Create a mock that allows testing a service without a database.

## Next Up

Continue to the [Profiling track](../01-profiling) to measure what your code costs.
