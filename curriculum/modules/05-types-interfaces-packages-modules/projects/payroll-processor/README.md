# Payroll Processor

## Learning objective

Build a polymorphic batch payroll calculation engine in Go. Master domain modeling with Go interfaces, implicit interface satisfaction, custom typed error handling, overtime compensation rules, and multi-employee batch aggregation.

## Why this matters

Real-world financial systems handle diverse employment contracts (full-time salaried, hourly shift workers with overtime rules, 1099 contractors with distinct tax withholding). In Go, interface polymorphism decouples the calculation engine from specific employee implementations without bulky class hierarchies.

## Prerequisites

Complete all Module 05 lessons (01–23) on interfaces, method sets, struct embedding, implicit satisfaction, error wrapping, and package boundaries.

## Available concepts

This project draws on structs, custom types, `Payable` interface polymorphism, error handling (`errors.Is`), slice aggregation, and table-driven testing.

## Project overview

The `_starter/` directory provides type definitions, model layouts, and calculation stubs. Your job is to implement compensation calculations for each employee type and batch aggregation in `PayrollEngine` so that:

- `SalariedEmployee` calculates gross monthly salary, applies benefit deductions and tax withholding.
- `HourlyEmployee` calculates standard pay up to 40 hours and applies 1.5x overtime multiplier for hours exceeding 40.
- `Contractor` receives full gross pay with zero deductions (1099 independent contractor).
- `PayrollEngine.ProcessBatch` polymorphically evaluates a slice of `[]Payable` and aggregates total gross, deductions, and net compensation.
- Validation errors (negative hours, invalid rates) are caught and reported cleanly.

## Compensation Specifications

| Type | Gross Calculation | Deductions | Net Pay |
|---|---|---|---|
| **Salaried** | `monthlySalaryCents` | `benefitsCents + (gross * taxRate)` | `Gross - Deductions` |
| **Hourly (≤ 40h)** | `hours * hourlyRateCents` | `gross * taxRate` | `Gross - Deductions` |
| **Hourly (> 40h)** | `(40 * rate) + ((hours - 40) * rate * 1.5)` | `gross * taxRate` | `Gross - Deductions` |
| **Contractor** | `hoursBilled * hourlyRateCents` | `0` (Zero withholding) | `Gross` |

## Tasks

### Task 1 — Implement SalariedEmployee

In `_starter/main.go`:
- Implement `CalculatePay() (PayBreakdown, error)` for `*SalariedEmployee`:
  - Calculate tax deduction: `int(float64(gross) * s.taxRate)`.
  - Sum deductions: `s.benefitsDeductionCents + tax`.
  - Return `PayBreakdown` with Gross, Deductions, and Net pay.

### Task 2 — Implement HourlyEmployee with overtime

- Implement `CalculatePay() (PayBreakdown, error)` for `*HourlyEmployee`:
  - Reject negative hours with `ErrNegativeHours`.
  - For hours up to 40, pay at standard `hourlyRateCents`.
  - For hours exceeding 40, calculate regular pay for 40 hours plus 1.5x the rate for hours beyond 40.
  - Apply tax deduction and return `PayBreakdown`.

### Task 3 — Implement Contractor

- Implement `CalculatePay() (PayBreakdown, error)` for `*Contractor`:
  - Reject negative hours with `ErrNegativeHours`.
  - Multiply `hoursBilled` by `hourlyRateCents`. Deductions must be 0.

### Task 4 — Implement PayrollEngine batch processing

- Implement `ProcessBatch(employees []Payable) (PayrollSummary, error)`:
  - Guard against nil employee entries.
  - Iterate through `employees` polymorphically calling `CalculatePay()`.
  - Aggregate `TotalGrossCents`, `TotalDeductionsCents`, and `TotalNetCents`.
  - Aggregate individual errors if any employee fails calculation.

## Verification

Run all automated unit tests:

```bash
go test -v ./curriculum/modules/05-types-interfaces-packages-modules/projects/payroll-processor/_starter
```

Expected output:
```
=== RUN   TestSalariedEmployee_CalculatePay
--- PASS: TestSalariedEmployee_CalculatePay (0.00s)
=== RUN   TestHourlyEmployee_RegularAndOvertime
--- PASS: TestHourlyEmployee_RegularAndOvertime (0.00s)
=== RUN   TestContractor_NoDeductions
--- PASS: TestContractor_NoDeductions (0.00s)
=== RUN   TestValidationErrors
--- PASS: TestValidationErrors (0.00s)
=== RUN   TestProcessBatch_PolymorphicPayroll
--- PASS: TestProcessBatch_PolymorphicPayroll (0.00s)
PASS
```

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|---|---|---|---|
| **Polymorphic Design** | All employee types satisfy `Payable` implicitly; engine processes mixed slice with zero type switches | Interface used but concrete type switches present in engine | Concrete types hardcoded; no interface polymorphism |
| **Calculation Accuracy** | Regular, overtime (1.5x), benefits, and contractor zero-withholding calculated with cent precision | Calculations mostly accurate with off-by-one or rounding drift | Overtime rate incorrect or deductions miscalculated |
| **Validation & Error Handling** | Negative hours and invalid rates return sentinel errors; batch aggregation collects errors | Validation present but returns generic untyped errors | Negative hours accepted or causes panics |
| **Code Structure & Idioms** | Idiomatic struct definitions, integer cent arithmetic (no float currency errors), clean constructors | Mostly clean with minor unidiomatic constructs | Floating-point currency calculations introducing precision bugs |

## NEXT UP

[Module 06: Testing & Debugging](../../06-testing-debugging-refactoring/README.md)
