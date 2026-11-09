package main

import (
	"errors"
	"testing"
)

func TestSalariedEmployee_CalculatePay(t *testing.T) {
	// Salary: $6,000.00 (600,000 cents), Benefits: $400.00 (40,000 cents), Tax: 20%
	emp, err := NewSalariedEmployee("SAL-1", "Alice", 600000, 40000, 0.20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bd, err := emp.CalculatePay()
	if err != nil {
		t.Fatalf("calculate pay failed: %v", err)
	}

	if bd.GrossCents != 600000 {
		t.Fatalf("expected gross 600000, got %d", bd.GrossCents)
	}
	// Tax: 600000 * 0.20 = 120000 cents + 40000 benefits = 160000 deductions
	if bd.DeductionsCents != 160000 {
		t.Fatalf("expected deductions 160000, got %d", bd.DeductionsCents)
	}
	if bd.NetCents != 440000 {
		t.Fatalf("expected net 440000, got %d", bd.NetCents)
	}
}

func TestHourlyEmployee_RegularAndOvertime(t *testing.T) {
	// $50.00/hr (5000 cents/hr), Tax: 10%
	testCases := []struct {
		name      string
		hours     float64
		wantGross int
		wantTax   int
		wantNet   int
	}{
		{
			name:      "regular hours 30h",
			hours:     30.0,
			wantGross: 150000, // 30 * 5000
			wantTax:   15000,
			wantNet:   135000,
		},
		{
			name:      "exact 40h regular cap",
			hours:     40.0,
			wantGross: 200000, // 40 * 5000
			wantTax:   20000,
			wantNet:   180000,
		},
		{
			name:      "overtime 50h (40 reg + 10 OT @ 1.5x)",
			hours:     50.0,
			wantGross: 275000, // 200000 + (10 * 5000 * 1.5 = 75000)
			wantTax:   27500,
			wantNet:   247500,
		},
		{
			name:      "zero hours",
			hours:     0.0,
			wantGross: 0,
			wantTax:   0,
			wantNet:   0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			emp, err := NewHourlyEmployee("HOU-1", "Bob", 5000, tc.hours, 0.10)
			if err != nil {
				t.Fatalf("unexpected error creating hourly employee: %v", err)
			}

			bd, err := emp.CalculatePay()
			if err != nil {
				t.Fatalf("calculate pay failed: %v", err)
			}

			if bd.GrossCents != tc.wantGross {
				t.Fatalf("expected gross %d, got %d", tc.wantGross, bd.GrossCents)
			}
			if bd.DeductionsCents != tc.wantTax {
				t.Fatalf("expected deductions %d, got %d", tc.wantTax, bd.DeductionsCents)
			}
			if bd.NetCents != tc.wantNet {
				t.Fatalf("expected net %d, got %d", tc.wantNet, bd.NetCents)
			}
		})
	}
}

func TestContractor_NoDeductions(t *testing.T) {
	// $100.00/hr (10,000 cents/hr), 25 hours
	c, err := NewContractor("CON-1", "Carol", 10000, 25.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bd, err := c.CalculatePay()
	if err != nil {
		t.Fatalf("calculate pay failed: %v", err)
	}

	expectedGross := 250000
	if bd.GrossCents != expectedGross {
		t.Fatalf("expected gross %d, got %d", expectedGross, bd.GrossCents)
	}
	if bd.DeductionsCents != 0 {
		t.Fatalf("expected 0 deductions for 1099 contractor, got %d", bd.DeductionsCents)
	}
	if bd.NetCents != expectedGross {
		t.Fatalf("expected net %d, got %d", expectedGross, bd.NetCents)
	}
}

func TestValidationErrors(t *testing.T) {
	// Negative pay rate
	_, err := NewSalariedEmployee("S1", "Name", -1000, 0, 0.1)
	if !errors.Is(err, ErrInvalidPayRate) {
		t.Fatalf("expected ErrInvalidPayRate, got %v", err)
	}

	// Negative hours
	_, errHour := NewHourlyEmployee("H1", "Name", 5000, -5.0, 0.1)
	if !errors.Is(errHour, ErrNegativeHours) {
		t.Fatalf("expected ErrNegativeHours, got %v", errHour)
	}

	_, errCon := NewContractor("C1", "Name", 5000, -10.0)
	if !errors.Is(errCon, ErrNegativeHours) {
		t.Fatalf("expected ErrNegativeHours, got %v", errCon)
	}
}

func TestProcessBatch_PolymorphicPayroll(t *testing.T) {
	engine := NewPayrollEngine()

	salaried, _ := NewSalariedEmployee("S1", "Alice", 500000, 20000, 0.10) // Gross: 500,000, Ded: 70,000, Net: 430,000
	hourly, _ := NewHourlyEmployee("H1", "Bob", 5000, 40.0, 0.10)          // Gross: 200,000, Ded: 20,000, Net: 180,000
	contractor, _ := NewContractor("C1", "Carol", 10000, 10.0)             // Gross: 100,000, Ded: 0, Net: 100,000

	batch := []Payable{salaried, hourly, contractor}
	summary, err := engine.ProcessBatch(batch)
	if err != nil {
		t.Fatalf("batch processing failed: %v", err)
	}

	if summary.TotalEmployees != 3 {
		t.Fatalf("expected 3 employees, got %d", summary.TotalEmployees)
	}

	expectedGross := 500000 + 200000 + 100000  // 800,000
	expectedDed := 70000 + 20000 + 0           // 90,000
	expectedNet := expectedGross - expectedDed // 710,000

	if summary.TotalGrossCents != expectedGross {
		t.Fatalf("expected total gross %d, got %d", expectedGross, summary.TotalGrossCents)
	}
	if summary.TotalDeductionsCents != expectedDed {
		t.Fatalf("expected total deductions %d, got %d", expectedDed, summary.TotalDeductionsCents)
	}
	if summary.TotalNetCents != expectedNet {
		t.Fatalf("expected total net %d, got %d", expectedNet, summary.TotalNetCents)
	}
}
