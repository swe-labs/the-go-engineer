# Project Assessment — Payroll Processor

## Purpose

Evaluate your ability to build a polymorphic batch processing engine using Go interfaces. This assessment verifies that you can model diverse domain types behind a common interface, implement precise financial calculations with integer cent arithmetic, handle validation errors with sentinel types, and aggregate results across a heterogeneous collection.

## Targets

- `project-payroll-processor`

## Scope

This assessment covers:

- Implicit interface satisfaction with the `Payable` interface
- Struct definitions for distinct employee types (salaried, hourly, contractor)
- Method implementations with domain-specific business rules
- Integer cent arithmetic to avoid floating-point currency bugs
- Sentinel errors (`ErrNegativeHours`) and `errors.Is` matching
- Batch processing a `[]Payable` slice polymorphically
- Table-driven testing for financial calculations

## Evidence required

Submit:

1. Your completed `_starter/main.go` with all employee types and `PayrollEngine.ProcessBatch` implemented
2. Output of `go test -v .` showing all tests passing
3. A brief explanation of how `ProcessBatch` iterates over `[]Payable` without knowing the concrete types — what Go mechanism makes this work?

## Time expectation

Plan for 60 to 90 minutes.

## Rubric

| Criterion | Weight | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------|--------------|----------|----------------|
| **Polymorphic design** | 30% | All employee types satisfy `Payable` implicitly; engine processes mixed `[]Payable` slice with zero type switches or assertions | Interface used but `ProcessBatch` contains type switches to dispatch | Concrete types hardcoded; no interface polymorphism |
| **Calculation accuracy** | 30% | Regular pay, overtime (1.5x for hours > 40), benefits deductions, and contractor zero-withholding calculated with exact cent precision | Calculations mostly accurate with minor rounding errors | Overtime multiplier incorrect or floating-point currency bugs |
| **Validation and error handling** | 20% | Negative hours and nil employees return sentinel errors; `errors.Is` works on returned errors; batch continues after individual failures | Validation present but generic string errors returned | Negative hours accepted silently or cause panics |
| **Code structure** | 20% | Idiomatic struct definitions, integer cent arithmetic throughout, clean constructors with validation | Mostly clean with minor issues (e.g., float currency math) | Monolithic function or floating-point currency calculations |

**Passing score**: 75%

## Review questions

1. Why does Go not require an `implements` keyword? How does the compiler know that `HourlyEmployee` satisfies `Payable`?
2. Walk through the overtime calculation for an hourly employee with 52 hours at $20.00/hour. Show the math step by step using cents.
3. Why does the project use `int` (cents) for currency instead of `float64`? What specific bug would `float64` introduce?
4. What happens inside `ProcessBatch` when one employee in the batch returns an error? Does the batch stop or continue?
5. How would you add a new employee type (e.g., `InternEmployee` with capped hours) without modifying `ProcessBatch`?

## Passing standard

You pass when all tests pass, your batch engine processes mixed employee types polymorphically, and you can explain interface satisfaction and cent arithmetic mechanically.

## Retake policy

If you score below 75%, review Module 05 lessons on interfaces, method sets, and type polymorphism, then resubmit.

## NEXT UP

[Module 05 Checkpoint](../checkpoint/README.md)
