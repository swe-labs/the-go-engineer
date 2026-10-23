# Project Assessment — Bug Hunt Pack

## Purpose

Evaluate your ability to diagnose, reproduce, and fix realistic production bugs in Go. This assessment verifies that you can use the Go race detector, identify nil map panics, fix slice bounds errors, and detect shadowed error variables — skills that distinguish production-ready engineers from beginners.

## Targets

- `project-bug-hunt-pack`

## Scope

This assessment covers:

- Concurrency data race detection with `go test -race`
- `sync.Mutex` for protecting shared state
- Map zero-value behavior and lazy initialization
- Slice bounds checking and safe clamping
- Variable shadowing with `:=` in nested scopes
- Regression test suites to prevent reintroduction

## Evidence required

Submit:

1. Your fixed `_starter/main.go` with all 4 bugs resolved
2. Output of `go test -v -race .` showing all 4 tests passing with zero race warnings
3. For each bug, a brief note (1–2 sentences) explaining: what the symptom was, what the root cause was, and how you fixed it

## Time expectation

Plan for 60 to 90 minutes.

## Rubric

| Criterion | Weight | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------|--------------|----------|----------------|
| **Diagnosis and root cause** | 30% | All 4 bugs identified correctly; fixes address mechanical root causes (not symptoms) | 3 bugs fixed; 1 bug masked with a workaround | Bugs unresolved or still panicking |
| **Race detector cleanliness** | 25% | Counter is completely thread-safe with proper `Mutex` usage; `go test -race` produces zero warnings | Mutex used but lock contention exists or unlock misplaced | Data race persists under concurrent access |
| **Defensive programming** | 25% | Zero-value structs handled gracefully with lazy `make()`; slice bounds clamped safely for all inputs (empty, oversized offset, negative) | Fixed for happy path only; edge inputs still panic | Nil map or out-of-range panics persist |
| **Error propagation** | 20% | Shadowed error variable eliminated; original error returned cleanly through all code paths | Shadowing resolved but error type information lost | Errors swallowed silently due to shadowing |

**Passing score**: 75%

## Review questions

1. How did you discover Bug 1 (data race)? What does `go test -race` actually check at runtime? Does it catch all possible races?
2. Explain Bug 2: why does `var cache KeyValueCache` panic on `cache.Set()`? What is the zero value of a `map` field inside a struct?
3. For Bug 3 (slice bounds): show the calculation you used to clamp `end` to `len(items)`. What happens when `offset` is greater than `len(items)`?
4. Explain Bug 4: write two lines of code that demonstrate how `:=` creates a new variable in a nested scope instead of assigning to the outer one. Why does Go allow this?
5. After fixing all 4 bugs, how do you ensure they never come back? What role do regression tests play?

## Passing standard

You pass when all 4 tests pass with `-race`, your fixes address root causes (not symptoms), and you can explain each bug's mechanical cause.

## Retake policy

If you score below 75%, review Module 06 lessons on race detection, debugging, and variable scoping, then resubmit with corrections.

## NEXT UP

[Project Assessment: Refactor Under Tests](../refactor-under-tests/README.md)
