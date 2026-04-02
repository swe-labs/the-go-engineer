// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"math"
	"testing"
)

// ============================================================================
// Tests for: Payroll Processor (Interface Polymorphism)
// ============================================================================
//
// These tests verify CalculatePay() returns the correct amount
// for each employee type, proving the interface contract is satisfied.
//
// RUN: go test -v ./05-types-and-interfaces/6-payroll-processor
// ============================================================================

func TestSalariedEmployeePay(t *testing.T) {
	emp := SalariedEmployee{Name: "Alice", AnnualSalary: 120000}
	got := emp.CalculatePay()
	want := 10000.0 // 120000 / 12

	if got != want {
		t.Errorf("SalariedEmployee.CalculatePay() = %.2f, want %.2f", got, want)
	}
}

func TestHourlyEmployeePay(t *testing.T) {
	emp := HourlyEmployee{Name: "Bob", HourlyRate: 25.0, HoursWorked: 160.0}
	got := emp.CalculatePay()
	want := 4000.0 // 25 * 160

	if got != want {
		t.Errorf("HourlyEmployee.CalculatePay() = %.2f, want %.2f", got, want)
	}
}

func TestCommissionEmployeePay(t *testing.T) {
	emp := CommissionEmployee{
		Name:           "Charlie",
		BaseSalary:     3000.0,
		CommissionRate: 0.05,
		SalesAmount:    50000.0,
	}
	got := emp.CalculatePay()
	want := 5500.0 // 3000 + (0.05 * 50000)

	if got != want {
		t.Errorf("CommissionEmployee.CalculatePay() = %.2f, want %.2f", got, want)
	}
}

func TestPayableInterface(t *testing.T) {
	// Verify all types satisfy the Payable interface at compile-time
	// by assigning them to an interface variable.
	employees := []Payable{
		SalariedEmployee{Name: "A", AnnualSalary: 60000},
		HourlyEmployee{Name: "B", HourlyRate: 20, HoursWorked: 80},
		CommissionEmployee{Name: "C", BaseSalary: 2000, CommissionRate: 0.10, SalesAmount: 10000},
	}

	expectedPays := []float64{5000.0, 1600.0, 3000.0}

	for i, emp := range employees {
		got := emp.CalculatePay()
		if math.Abs(got-expectedPays[i]) > 0.01 {
			t.Errorf("Employee[%d].CalculatePay() = %.2f, want %.2f", i, got, expectedPays[i])
		}
	}
}
