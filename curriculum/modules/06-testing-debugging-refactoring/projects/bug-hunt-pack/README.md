# Bug Hunt Pack

## Learning objective

Diagnose, reproduce, and fix subtle production bugs in Go codebases. Master concurrency data race diagnosis using the Go race detector (`-race`), nil map initialization panics, slice boundary clamp defects, and silent error variable shadowing.

## Why this matters

Writing new code is only half of a software engineer's job; debugging subtle edge cases and race conditions in existing systems is what separates junior developers from production-ready engineers. In Go, mistakes like data races, writing to nil maps, and shadowing errors with `:=` can slip past compilers and lead to catastrophic production failures.

## Prerequisites

Complete all Module 06 lessons (01–13) on unit testing with `testing.T`, table-driven tests, subtests, race condition detection with `-race`, Delve debugger basics, and error tracing.

## Available concepts

This project draws on `go test`, `go test -race`, `sync.Mutex`, `sync.RWMutex`, slice bounds expressions, map zero-values, variable shadowing, and regression test suites.

## Project overview

The `_starter/` directory contains four deliberate, realistic Go bugs. The test suite in `_starter/main_test.go` catches these bugs and will fail or panic until they are diagnosed and resolved:

1. **Bug 1: Data Race in Counter** — `SafeCounter.Inc()` mutates an integer without synchronization, failing under `go test -race`.
2. **Bug 2: Nil Map Panic** — `KeyValueCache` panics when accessed as a zero-value struct (`var cache KeyValueCache`) because maps must be initialized with `make()`.
3. **Bug 3: Slice Bounds Out of Range** — `PaginateSlice` panics when `offset + limit` exceeds the slice length instead of clamping to the end.
4. **Bug 4: Shadowed Error Variable** — `ValidateAndFormatUser` shadows an error variable using short declaration (`:=`), causing formatting errors to be silently swallowed.

The `_solution/` directory provides clean, idiomatic reference fixes.

## Tasks

### Task 1 — Diagnose & fix Bug 1 (Data Race)

Run `go test -race .` in `_starter/`. Identify the unsynchronized access in `SafeCounter`.
- Add a `sync.Mutex` to protect `c.count`.
- Lock and unlock around both increments and reads.

### Task 2 — Diagnose & fix Bug 2 (Nil Map Panic)

Run `go test -run TestBug2 .`. Notice the panic on `cache.Set()`.
- Ensure `c.data` is lazily initialized if nil (`if c.data == nil { c.data = make(...) }`).

### Task 3 — Diagnose & fix Bug 3 (Slice Bounds Panic)

Run `go test -run TestBug3 .`.
- Clamp `offset` (minimum 0) and `end := offset + limit`. If `end > len(items)`, set `end = len(items)`.

### Task 4 — Diagnose & fix Bug 4 (Shadowed Error)

Run `go test -run TestBug4 .`.
- In `ValidateAndFormatUser`, declare `var formatted string` separately from `err` so `formatName`'s returned error is not shadowed by `:=`.

## Verification

Run the entire test suite with the Go race detector enabled:

```bash
go test -v -race ./curriculum/modules/06-testing-debugging-refactoring/projects/bug-hunt-pack/_starter
```

Expected output:
```
=== RUN   TestBug1_ConcurrentCounter
--- PASS: TestBug1_ConcurrentCounter (0.00s)
=== RUN   TestBug2_NilMapCache
--- PASS: TestBug2_NilMapCache (0.00s)
=== RUN   TestBug3_SlicePaginationBounds
--- PASS: TestBug3_SlicePaginationBounds (0.00s)
=== RUN   TestBug4_ShadowedError
--- PASS: TestBug4_ShadowedError (0.00s)
PASS
```

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|---|---|---|---|
| **Diagnosis & Root Cause** | All 4 bugs identified correctly; fixes address root causes rather than symptoms | 3 bugs fixed; 1 bug masked with workaround | Bugs unresolved or panicking |
| **Race Detector Cleanliness** | Counter is completely thread-safe; passes `go test -race` with zero warnings | Mutex used but lock contention or unlock missing | Data race remains in code |
| **Defensive Programming** | Zero-value structs handled gracefully; slice bounds safely clamped | Fixed for happy path only; edge cases still panic | Nil map or out-of-range panics persist |
| **Error Propagation** | Error shadowing eliminated; original error and sentinel types propagated cleanly | Shadowing resolved but error type lost | Errors ignored or swallowed |

## NEXT UP

[Refactor Under Tests](../refactor-under-tests/README.md)
