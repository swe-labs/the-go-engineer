# Project Assessment — Refactor Under Tests

## Purpose

Evaluate your ability to safely refactor legacy Go code using characterization tests. This assessment verifies that you can lock down existing behavior with automated tests before decomposing monolithic functions into clean domain types — proving regression-free refactoring through a disciplined engineering workflow.

## Targets

- `project-refactor-under-tests`

## Scope

This assessment covers:

- Characterization testing to capture existing behavior before refactoring
- Table-driven tests with `testing.T` and `t.Run` subtests
- Test coverage analysis with `go test -cover`
- Domain modeling with custom struct types
- Type-safe enumerations using `type CustomerTier string`
- Method-based calculations (separation of parsing from computation)
- The characterization → decompose → verify refactoring lifecycle
- Integer cent arithmetic for financial calculations

## Evidence required

Submit:

1. Your refactored `_starter/main.go` with domain types (`CustomerTier`, `LineItem`, `Invoice`) extracted from the monolithic function
2. Output of `go test -v -race .` showing all characterization tests passing without modification
3. Output of `go test -cover .` showing coverage percentage
4. A brief explanation of why you wrote characterization tests before refactoring — what would happen if you refactored first?

## Time expectation

Plan for 90 to 120 minutes.

## Rubric

| Criterion | Weight | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------|--------------|----------|----------------|
| **Characterization coverage** | 25% | Tests cover all customer tiers (GOLD/SILVER/BRONZE/STANDARD), all state tax rates (CA/NY/TX/other), multi-line invoices, and malformed inputs | Tests cover happy path with minimal edge cases | Incomplete tests allowing regression bugs |
| **Domain architecture** | 30% | Monolithic function decomposed into cohesive types (`Invoice`, `LineItem`, `CustomerTier`); each type has single-responsibility methods | Refactored into smaller functions but no struct types used | Code remains monolithic after "refactoring" |
| **Regression safety** | 25% | All original characterization tests pass without any modifications; backwards compatibility preserved for `ProcessInvoice` signature | Tests pass but public function signature was modified | Tests fail after refactoring |
| **Code cleanliness** | 20% | Magic numbers and string literals replaced with typed constants; idiomatic error handling with `(result, error)` returns | Some magic numbers remain; error handling mostly clean | Poor formatting, unidiomatic error handling, or panics |

**Passing score**: 75%

## Review questions

1. What is a characterization test? How does it differ from a unit test you write for new code?
2. Walk through your `CustomerTier.DiscountRate()` method. Why is `type CustomerTier string` better than using raw string constants?
3. Your original `ProcessInvoice` was monolithic. After refactoring, how many types and methods exist? What principle guided the decomposition?
4. Why must characterization tests be written before refactoring — not after? What guarantee do they provide?
5. Show me the `go test -cover` output. Which lines are uncovered? Are those lines reachable?

## Passing standard

You pass when all characterization tests pass without modification, your domain types are cohesive and single-purpose, and you can explain the refactoring lifecycle mechanically.

## Retake policy

If you score below 75%, review Module 06 lessons on characterization testing, test-driven refactoring, and domain modeling, then resubmit.

## NEXT UP

[Module 06 Checkpoint](../checkpoint/README.md)
