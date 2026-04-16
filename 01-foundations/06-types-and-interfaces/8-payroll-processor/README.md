# TI.8 Payroll Processor Project

## Mission

Build a small payroll system that proves multiple concrete employee types can share one behavior contract through interfaces instead of inheritance.

This exercise is the Section 06 milestone. It is where structs, methods, interfaces, and a small generic helper come together in one runnable artifact with tests.

## Why This Lesson Exists Now

You have learned all the individual pieces: structs for data, methods for behavior, interfaces for contracts, and generics for reusable helpers. Now it is time to put them together in one realistic application.

## Prerequisites

- `TI.1` structs
- `TI.2` methods
- `TI.3` interfaces
- `TI.4` Stringer
- `TI.5` generics

## What You Will Build

Implement a small payroll processor that:

1. defines a `Payable` interface with behavior contracts
2. models multiple concrete employee types with structs
3. attaches pay-calculation behavior through methods
4. processes a mixed `[]Payable` collection polymorphically
5. includes one small generic helper for reusable summary/report output
6. verifies core pay calculations with tests

## Visual Model

```text
┌─────────────────────────┐
│ Payable interface       │
├─────────────────────────┤
│ CalculatePay() float64 │
│ String() string         │
└─────────────────────────┘
         ▲
         │ implements
    ┌────┴────────┬──────────┐
    │             │          │
┌───┴───┐   ┌────┴────┐  ┌──┴─────────┐
│Salaried│   │ Hourly  │  │Commission  │
└───────┘   └─────────┘  └────────────┘
```

## Machine View

The payroll processor uses polymorphism to process different employee types through a common interface. At runtime, Go uses dynamic dispatch (iTables) to look up which concrete type's CalculatePay method to execute.

## Files

- [main.go](./main.go): complete solution with teaching comments
- [payroll_test.go](./payroll_test.go): tests for the pay-calculation contract
- [_starter/main.go](./_starter/main.go): starter file with TODOs and requirements

## Run Instructions

Run the completed solution:

```bash
go run ./01-foundations/06-types-and-interfaces/8-payroll-processor
```

Run the tests:

```bash
go test ./01-foundations/06-types-and-interfaces/8-payroll-processor
```

Run the starter:

```bash
go run ./01-foundations/06-types-and-interfaces/8-payroll-processor/_starter
```

## Code Walkthrough

### `type Payable interface { ... }`

This interface defines the contract: any type with `CalculatePay()` and `String()` methods satisfies Payable. Go checks this implicitly—no "implements" keyword needed.

### `type SalariedEmployee struct { ... }`

A salaried employee has an annual salary. CalculatePay divides by 12 for monthly pay.

### `type HourlyEmployee struct { ... }`

An hourly employee has an hourly rate and hours worked. CalculatePay multiplies them.

### `type CommissionEmployee struct { ... }`

A commission employee has a base salary plus commission on sales. CalculatePay adds base + (rate × sales).

### `PrintEmployeeSummary` function

This is a generic function with type parameter P constrained to `fmt.Stringer`. The employee parameter P can be any type that implements String().

### `func ProcessPayroll(employees []Payable)`

This function accepts a slice of any types that satisfy Payable. It iterates and calls CalculatePay polymorphically—one function handles all employee types.

## Try It

1. Add a new employee type (e.g., Contractor) and see if ProcessPayroll handles it without changes.
2. Change one CalculatePay to use a pointer receiver and verify it still works.
3. Add a new field to one employee type and update its String() method.

## Success Criteria

Your finished solution should:

- use interfaces to describe the payroll behavior boundary
- allow multiple employee types to satisfy the same contract cleanly
- keep receiver choices intentional and easy to explain
- include one small generic helper without turning the exercise into a generics showcase
- pass the provided tests

## Common Failure Modes

- building separate payroll logic for each employee type instead of using one interface
- using pointer receivers or value receivers without being able to explain why
- treating generics as the main point of the exercise instead of a small reuse helper
- tightly coupling `ProcessPayroll` to one concrete employee struct

## Verification Surface

Use these proof surfaces together:

1. `go run ./01-foundations/06-types-and-interfaces/8-payroll-processor`
2. `go run ./01-foundations/06-types-and-interfaces/8-payroll-processor/_starter`
3. `go test ./01-foundations/06-types-and-interfaces/8-payroll-processor`

## Production Relevance

This exercise simulates real-world scenarios where different employee types (hourly, salaried, contractors) all need to be processed through a common payroll system. The interface-based approach is how real Go applications handle polymorphic scenarios.

## Next Step

After you complete this exercise, move to `TI.7` advanced generics for a stretch lesson, or continue to the next section if you are ready to move on.