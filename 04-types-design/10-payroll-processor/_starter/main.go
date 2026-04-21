// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package main

import "fmt"

// ============================================================================
// Section 6: Types & Interfaces — Payroll Processor (Exercise Starter)
// Level: Intermediate
// ============================================================================
//
// EXERCISE: Build a Polymorphic Payroll Processor
//
// REQUIREMENTS:
//  1. [ ] Define a `Payable` interface with `CalculatePay() float64` and `String() string`
//  2. [ ] Implement `SalariedEmployee` — pay = AnnualSalary / 12
//  3. [ ] Implement `HourlyEmployee` — pay = HourlyRate * HoursWorked
//  4. [ ] Implement `CommissionEmployee` — pay = BaseSalary + (CommissionRate * SalesAmount)
//  5. [ ] Create a `ProcessPayroll(employees []Payable)` function that prints each
//         employee's info, monthly pay, and calculates the total payroll cost
//  6. [ ] Add one small generic helper for reusable employee summary output
//
// HINTS:
//   - Any type implementing CalculatePay() and String() satisfies Payable
//   - Use fmt.Sprintf in String() to format employee details
//   - Range over the []Payable slice to process all employees generically
//   - A small generic helper can constrain on `fmt.Stringer`
//
// RUN: go run ./04-types-design/10-payroll-processor/_starter
// SOLUTION: See the main.go file in the parent directory
// ============================================================================

// TODO: Define the Payable interface

// TODO: Define SalariedEmployee, HourlyEmployee, CommissionEmployee structs

// TODO: Implement CalculatePay() and String() for each type

// TODO: Implement ProcessPayroll(employees []Payable)

func main() {
	fmt.Println("=== Payroll Processor Exercise ===")
	fmt.Println()
	fmt.Println("TODO: Implement your payroll system!")
	fmt.Println("See the REQUIREMENTS above for what to build.")
	fmt.Println()
	fmt.Println("When finished, compare your solution with ../main.go")
}
