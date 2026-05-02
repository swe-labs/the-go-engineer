// Copyright (c) 2026 Rasel Hossen

// ============================================================================
// Section 04: Types and Design
// Title: Payroll Processor Project
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
// WHAT YOU'LL LEARN:
//   - Designing behavioral contracts using Go interfaces.
//   - Implementing polymorphic processing logic for heterogeneous types.
//   - Embedding standard library interfaces (fmt.Stringer) into custom contracts.
//   - Utilizing generic helpers to process interface-satisfied types.
//
// WHY THIS MATTERS:
//   - In a production system, you often need to process varied data
//     models (Salaried vs Hourly employees) through a single workflow.
//     Interfaces allow you to decouple the "what" (calculating pay) from
//     the "how" (specific salary formulas), resulting in highly
//     extensible and testable codebases.
//
// RUN:
//   go run ./04-types-design/10-payroll-processor
//
// KEY TAKEAWAY:
//   - Interfaces enable uniform treatment of distinct types through shared contracts.
// ============================================================================

// See LICENSE for usage terms.

package main

//

import "fmt"

// Section 04: Types & Design - Payroll Processor

// Payable defines the behavioral contract for any entity that requires a payroll calculation.
// It embeds fmt.Stringer to ensure all payable entities can be logged/printed.
type Payable interface {
	fmt.Stringer
	CalculatePay() float64
}

// SalariedEmployee represents an employee with a fixed annual salary.
type SalariedEmployee struct {
	Name         string
	AnnualSalary float64
}

// CalculatePay computes the monthly salary.
func (se SalariedEmployee) CalculatePay() float64 {
	return se.AnnualSalary / 12.0
}

func (se SalariedEmployee) String() string {
	return fmt.Sprintf("Salaried: %s (Annual: $%.2f)", se.Name, se.AnnualSalary)
}

// HourlyEmployee represents an employee paid by the hour.
type HourlyEmployee struct {
	Name        string
	HourlyRate  float64
	HoursWorked float64
}

// CalculatePay computes pay based on hours worked.
func (he HourlyEmployee) CalculatePay() float64 {
	return he.HourlyRate * he.HoursWorked
}

func (he HourlyEmployee) String() string {
	return fmt.Sprintf("Hourly: %s (Rate: $%.2f/hr, Hours: %.1f)", he.Name, he.HourlyRate, he.HoursWorked)
}

// CommissionEmployee represents an employee with a base salary plus sales commission.
type CommissionEmployee struct {
	Name           string
	BaseSalary     float64
	CommissionRate float64
	SalesAmount    float64
}

// CalculatePay computes total pay including base salary and commission.
func (ce CommissionEmployee) CalculatePay() float64 {
	return ce.BaseSalary + (ce.CommissionRate * ce.SalesAmount)
}

func (ce CommissionEmployee) String() string {
	return fmt.Sprintf("Commission: %s (Base: $%.2f, CommRate: %.2f%%, Sales: $%.2f)",
		ce.Name, ce.BaseSalary, ce.CommissionRate*100, ce.SalesAmount)
}

// PrintEmployeeSummary outputs the textual representation of any Stringer.
func PrintEmployeeSummary[P fmt.Stringer](employee P) {
	fmt.Printf("  - Processing: %s\n", employee)
}

// ProcessPayroll iterates over a list of Payable entities and prints the payroll summary.
func ProcessPayroll(employees []Payable) {
	fmt.Println("--- Generating Payroll Summary ---")
	totalPayroll := 0.0
	for _, emp := range employees {
		PrintEmployeeSummary(emp)
		pay := emp.CalculatePay()
		fmt.Printf("    Calculated Monthly Pay: $%.2f\n", pay)
		totalPayroll += pay
	}
	fmt.Printf("\nAggregate Monthly Liability: $%.2f\n", totalPayroll)
	fmt.Println("----------------------------------")
}

func main() {
	fmt.Println("=== Case Study: Payroll Processor ===")
	fmt.Println()

	// 1. Heterogeneous Data.
	// We instantiate different concrete types that store distinct data fields.
	salEmp := SalariedEmployee{Name: "Alice Wonderland", AnnualSalary: 72000.00}
	hrEmp := HourlyEmployee{Name: "Bob The Builder", HourlyRate: 25.00, HoursWorked: 160.0}
	comEmp := CommissionEmployee{Name: "Charlie Chaplin", BaseSalary: 2000.00, CommissionRate: 0.10, SalesAmount: 15000.00}

	// 2. Uniform Treatment.
	// Because all these types satisfy the Payable interface, we can store
	// them in a single slice and process them using the same logic.
	payrollList := []Payable{
		salEmp,
		hrEmp,
		comEmp,
		HourlyEmployee{Name: "Diana Prince", HourlyRate: 30.00, HoursWorked: 150.0},
	}

	ProcessPayroll(payrollList)

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: CO.1 -> 04-types-design/16-composition")
	fmt.Println("Run    : go run ./04-types-design/16-composition")
	fmt.Println("Current: TI.10 (payroll-processor)")
	fmt.Println("---------------------------------------------------")
}
