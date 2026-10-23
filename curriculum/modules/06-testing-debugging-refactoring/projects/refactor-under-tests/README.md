# Refactor Under Tests

## Learning objective

Master the disciplined practice of characterization testing and safe code refactoring in Go. Learn how to wrap monolithic, untested legacy code in comprehensive automated tests to lock in existing behavior before breaking it down into clean, cohesive domain types without introducing regressions.

## Why this matters

In professional software engineering, you rarely write software from a blank slate. More frequently, you inherit tangled, monolithic functions written years ago. Refactoring without characterization tests is reckless guessing; refactoring *under* automated characterization tests is predictable, verified engineering.

## Prerequisites

Complete all Module 06 lessons (01–13) on table-driven tests, subtests, test coverage analysis (`go test -cover`), benchmark profiling, and regression safety.

## Available concepts

This project draws on table-driven testing (`testing.T`), subtests (`t.Run`), characterization test suites, domain refactoring, struct modeling, type safety with custom enum types, and `go test -race`.

## Project overview

The `_starter/` directory contains a monolithic legacy function `ProcessInvoice(raw string) (string, int, error)` that parses CSV order rows, computes discounts and sales taxes, and formats output strings.

Your task follows the standard professional refactoring lifecycle:
1. **Characterize:** Run and expand the characterization tests in `_starter/main_test.go` to lock down all calculation rules, customer tier discounts, state tax rates, and edge-case errors.
2. **Decompose:** Refactor the messy function into clean, decoupled domain structures:
   - `CustomerTier`: type-safe enum (`GOLD`, `SILVER`, `BRONZE`, `STANDARD`) with `DiscountRate() float64`.
   - `LineItem`: represents item quantities, prices, and line subtotals.
   - `Invoice`: aggregates line items, computes subtotal, discounts, state taxes, and grand totals.
   - `ParseInvoice(raw string) (*Invoice, error)`: handles validation and CSV parsing.
3. **Verify:** Ensure `ProcessInvoice` delegates cleanly to the new domain models and passes every characterization test without modification.

The `_solution/` directory provides the completed, clean domain architecture.

## Specifications

### Customer Tier Discounts
- `GOLD`: 15% discount (`0.15`)
- `SILVER`: 10% discount (`0.10`)
- `BRONZE`: 5% discount (`0.05`)
- `STANDARD`: 0% discount

### State Sales Tax
- `CA`: 9.0% (`0.09`)
- `NY`: 8.0% (`0.08`)
- `TX`: 6.25% (`0.0625`)
- Other states: 5.0% (`0.05`)

## Tasks

### Task 1 — Characterize legacy behavior

Run the characterization suite in `_starter/`:
```bash
go test -v .
```
Verify that existing calculations and error cases (empty inputs, negative quantities, negative prices) are properly locked down.

### Task 2 — Extract domain models

In `_starter/main.go`, replace monolithic parsing logic with domain types:
- Define `type CustomerTier string` with method `DiscountRate() float64`.
- Define `type LineItem struct` with method `Subtotal() int`.
- Define `type Invoice struct` with methods `Subtotal() int`, `DiscountCents() int`, `TaxCents() int`, and `TotalCents() int`.

### Task 3 — Separate parsing from calculation

- Implement `ParseInvoice(raw string) (*Invoice, error)`.
- Re-implement `ProcessInvoice(raw string) (string, int, error)` as a thin wrapper that parses the raw payload and calls `invoice.TotalCents()`.

### Task 4 — Prove regression-free refactoring

Run `go test -v -race .`. All tests must pass cleanly without modifying the test suite.

## Verification

Run all characterization and regression tests:

```bash
go test -v -race ./curriculum/modules/06-testing-debugging-refactoring/projects/refactor-under-tests/_starter
```

Expected output:
```
=== RUN   TestProcessInvoice_ValidCalculations
--- PASS: TestProcessInvoice_ValidCalculations (0.00s)
=== RUN   TestProcessInvoice_ValidationErrors
--- PASS: TestProcessInvoice_ValidationErrors (0.00s)
PASS
```

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|---|---|---|---|
| **Characterization Coverage** | Tests cover all tiers, state taxes, multi-line invoices, and malformed inputs | Tests cover happy path with minimal edge cases | Incomplete tests allowing regression bugs |
| **Domain Architecture** | Monolithic function decomposed into cohesive types (`Invoice`, `LineItem`, `CustomerTier`); clean single-responsibility methods | Refactored into smaller functions, but structs not used | Code remains monolithic |
| **Regression Safety** | All characterization tests pass without changes; backwards compatibility preserved | Tests pass but public function signature was modified | Tests fail after refactoring |
| **Code Cleanliness** | Magic numbers and string literals replaced with typed constants; idiomatic error handling | Some magic numbers remain | Poor formatting or unidiomatic error handling |

## NEXT UP

[Module 07: CLI, Files, and Config](../../07-cli-files-json-config/README.md)
