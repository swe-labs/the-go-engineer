// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================
// Section 5: Types & Interfaces — Payroll Processor (Exercise)
// Level: Intermediate
// ============================================================================
//
// RUN: go run ./05-types-and-interfaces/6-payroll-processor
// ============================================================================

import "fmt"

// Payable interface demonstrates Go's implicit interface satisfaction.
// Any struct that implements BOTH `String()` and `CalculatePay()` is automatically a Payable.
type Payable interface {
	fmt.Stringer           // Embedding another interface (Requires String() string)
	CalculatePay() float64 // Calculates monthly pay
}

type SalariedEmployee struct {
	Name         string
	AnnualSalary float64
}

// CalculatePay uses a VALUE RECEIVER `(se SalariedEmployee)`.
// We use a value receiver here because calculating pay does not mutate the struct.
// The Go compiler passes a copy of the struct into this function.
func (se SalariedEmployee) CalculatePay() float64 {
	return se.AnnualSalary / 12.0
}

func (se SalariedEmployee) String() string {
	return fmt.Sprintf("Salaried: %s (Annual: $%.2f)", se.Name, se.AnnualSalary)
}

type HourlyEmployee struct {
	Name        string
	HourlyRate  float64
	HoursWorked float64 // Hours worked in the month
}

func (he HourlyEmployee) CalculatePay() float64 {
	return he.HourlyRate * he.HoursWorked
}

func (he HourlyEmployee) String() string {
	return fmt.Sprintf("Hourly: %s (Rate: $%.2f/hr, Hours: %.1f)", he.Name, he.HourlyRate, he.HoursWorked)
}

type CommissionEmployee struct {
	Name           string
	BaseSalary     float64 // Monthly base
	CommissionRate float64 // e.g., 0.05 for 5%
	SalesAmount    float64
}

func (ce CommissionEmployee) CalculatePay() float64 {
	return ce.BaseSalary + (ce.CommissionRate * ce.SalesAmount)
}

func (ce CommissionEmployee) String() string {
	return fmt.Sprintf("Commission: %s (Base: $%.2f, CommRate: %.2f%%, Sales: $%.2f)",
		ce.Name, ce.BaseSalary, ce.CommissionRate*100, ce.SalesAmount)
}

// PrintEmployeeSummary uses Go 1.18+ Generics.
// `[P fmt.Stringer]` constraints the type P to anything that implements a String() method.
func PrintEmployeeSummary[P fmt.Stringer](employee P) {
	fmt.Printf("  - Processing: %s\n", employee) // Implicitly invokes employee.String()
}

// ProcessPayroll accepts a slice of the Payable interface.
// This is pure polymorphism. At runtime, Go uses "Dynamic Dispatch" (iTables)
// to look up which underlying struct's CalculatePay method to execute.
func ProcessPayroll(employees []Payable) {
	fmt.Println("\n--- Processing Payroll ---")
	totalPayroll := 0.0
	for _, emp := range employees {
		PrintEmployeeSummary(emp)

		// The magic of Interfaces: We don't care if emp is hourly or salaried!
		pay := emp.CalculatePay()
		fmt.Printf("    Monthly Pay: $%.2f\n", pay)
		totalPayroll += pay
	}
	fmt.Printf("\nTotal Monthly Payroll: $%.2f\n", totalPayroll)
	fmt.Println("--------------------------")
}

func main() {

	fmt.Println("Welcome to the Payroll Processor!")

	salEmp := SalariedEmployee{Name: "Alice Wonderland", AnnualSalary: 72000.00}
	hrEmp := HourlyEmployee{Name: "Bob The Builder", HourlyRate: 25.00, HoursWorked: 160.0}
	comEmp := CommissionEmployee{Name: "Charlie Chaplin", BaseSalary: 2000.00, CommissionRate: 0.10, SalesAmount: 15000.00}

	payrollList := []Payable{
		salEmp,
		hrEmp,
		comEmp,
		HourlyEmployee{Name: "Diana Prince", HourlyRate: 30.00, HoursWorked: 150.0},
	}

	ProcessPayroll(payrollList)

}
