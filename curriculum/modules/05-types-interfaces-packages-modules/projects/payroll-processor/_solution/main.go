package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidPayRate = errors.New("pay rate must be strictly positive")
	ErrNegativeHours  = errors.New("hours cannot be negative")
	ErrNilEmployee    = errors.New("employee cannot be nil")
)

// EmployeeInfo contains general employee identity fields.
type EmployeeInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// PayBreakdown captures the gross, deductions, and net compensation.
type PayBreakdown struct {
	EmployeeID      string `json:"employee_id"`
	GrossCents      int    `json:"gross_cents"`
	DeductionsCents int    `json:"deductions_cents"`
	NetCents        int    `json:"net_cents"`
}

// Payable is the polymorphic interface satisfied by all compensable entities.
type Payable interface {
	EmployeeInfo() EmployeeInfo
	CalculatePay() (PayBreakdown, error)
}

// SalariedEmployee receives fixed monthly compensation with standard benefits deductions.
type SalariedEmployee struct {
	info                   EmployeeInfo
	monthlySalaryCents     int
	benefitsDeductionCents int
	taxRate                float64
}

// NewSalariedEmployee creates a new SalariedEmployee.
func NewSalariedEmployee(id, name string, salaryCents, benefitsCents int, taxRate float64) (*SalariedEmployee, error) {
	if salaryCents <= 0 {
		return nil, ErrInvalidPayRate
	}
	if benefitsCents < 0 || taxRate < 0 || taxRate > 1.0 {
		return nil, errors.New("invalid deductions or tax rate")
	}

	return &SalariedEmployee{
		info:                   EmployeeInfo{ID: id, Name: name, Role: "Salaried"},
		monthlySalaryCents:     salaryCents,
		benefitsDeductionCents: benefitsCents,
		taxRate:                taxRate,
	}, nil
}

func (s *SalariedEmployee) EmployeeInfo() EmployeeInfo {
	return s.info
}

func (s *SalariedEmployee) CalculatePay() (PayBreakdown, error) {
	gross := s.monthlySalaryCents
	tax := int(float64(gross) * s.taxRate)
	deductions := s.benefitsDeductionCents + tax
	if deductions > gross {
		deductions = gross
	}

	return PayBreakdown{
		EmployeeID:      s.info.ID,
		GrossCents:      gross,
		DeductionsCents: deductions,
		NetCents:        gross - deductions,
	}, nil
}

// HourlyEmployee is compensated based on hours worked with 1.5x overtime for hours > 40.
type HourlyEmployee struct {
	info            EmployeeInfo
	hourlyRateCents int
	hoursWorked     float64
	taxRate         float64
}

// NewHourlyEmployee creates a new HourlyEmployee.
func NewHourlyEmployee(id, name string, hourlyRateCents int, hoursWorked, taxRate float64) (*HourlyEmployee, error) {
	if hourlyRateCents <= 0 {
		return nil, ErrInvalidPayRate
	}
	if hoursWorked < 0 {
		return nil, ErrNegativeHours
	}
	if taxRate < 0 || taxRate > 1.0 {
		return nil, errors.New("invalid tax rate")
	}

	return &HourlyEmployee{
		info:            EmployeeInfo{ID: id, Name: name, Role: "Hourly"},
		hourlyRateCents: hourlyRateCents,
		hoursWorked:     hoursWorked,
		taxRate:         taxRate,
	}, nil
}

func (h *HourlyEmployee) EmployeeInfo() EmployeeInfo {
	return h.info
}

func (h *HourlyEmployee) CalculatePay() (PayBreakdown, error) {
	if h.hoursWorked < 0 {
		return PayBreakdown{}, ErrNegativeHours
	}

	var gross float64
	if h.hoursWorked <= 40.0 {
		gross = h.hoursWorked * float64(h.hourlyRateCents)
	} else {
		regularPay := 40.0 * float64(h.hourlyRateCents)
		overtimeHours := h.hoursWorked - 40.0
		overtimePay := overtimeHours * float64(h.hourlyRateCents) * 1.5
		gross = regularPay + overtimePay
	}

	grossCents := int(gross)
	taxCents := int(float64(grossCents) * h.taxRate)
	netCents := grossCents - taxCents

	return PayBreakdown{
		EmployeeID:      h.info.ID,
		GrossCents:      grossCents,
		DeductionsCents: taxCents,
		NetCents:        netCents,
	}, nil
}

// Contractor is an independent 1099 contractor with gross pay equal to net pay (zero withholding).
type Contractor struct {
	info            EmployeeInfo
	hourlyRateCents int
	hoursBilled     float64
}

// NewContractor creates a new Contractor.
func NewContractor(id, name string, hourlyRateCents int, hoursBilled float64) (*Contractor, error) {
	if hourlyRateCents <= 0 {
		return nil, ErrInvalidPayRate
	}
	if hoursBilled < 0 {
		return nil, ErrNegativeHours
	}

	return &Contractor{
		info:            EmployeeInfo{ID: id, Name: name, Role: "Contractor"},
		hourlyRateCents: hourlyRateCents,
		hoursBilled:     hoursBilled,
	}, nil
}

func (c *Contractor) EmployeeInfo() EmployeeInfo {
	return c.info
}

func (c *Contractor) CalculatePay() (PayBreakdown, error) {
	if c.hoursBilled < 0 {
		return PayBreakdown{}, ErrNegativeHours
	}

	gross := int(c.hoursBilled * float64(c.hourlyRateCents))
	return PayBreakdown{
		EmployeeID:      c.info.ID,
		GrossCents:      gross,
		DeductionsCents: 0, // No withholding for contractors
		NetCents:        gross,
	}, nil
}

// PayrollSummary aggregates results across a payroll cycle.
type PayrollSummary struct {
	TotalEmployees       int            `json:"total_employees"`
	TotalGrossCents      int            `json:"total_gross_cents"`
	TotalDeductionsCents int            `json:"total_deductions_cents"`
	TotalNetCents        int            `json:"total_net_cents"`
	Breakdowns           []PayBreakdown `json:"breakdowns"`
}

// PayrollEngine coordinates batch payroll calculations.
type PayrollEngine struct{}

// NewPayrollEngine initializes a PayrollEngine.
func NewPayrollEngine() *PayrollEngine {
	return &PayrollEngine{}
}

// ProcessBatch evaluates all employees polymorphically and aggregates totals.
func (e *PayrollEngine) ProcessBatch(employees []Payable) (PayrollSummary, error) {
	var summary PayrollSummary
	var errMsgs []string

	for idx, emp := range employees {
		if emp == nil {
			errMsgs = append(errMsgs, fmt.Sprintf("employee at index %d is nil", idx))
			continue
		}

		bd, err := emp.CalculatePay()
		if err != nil {
			errMsgs = append(errMsgs, fmt.Sprintf("employee %s (%s): %v", emp.EmployeeInfo().ID, emp.EmployeeInfo().Name, err))
			continue
		}

		summary.TotalEmployees++
		summary.TotalGrossCents += bd.GrossCents
		summary.TotalDeductionsCents += bd.DeductionsCents
		summary.TotalNetCents += bd.NetCents
		summary.Breakdowns = append(summary.Breakdowns, bd)
	}

	if len(errMsgs) > 0 {
		return summary, fmt.Errorf("payroll calculation errors: %s", strings.Join(errMsgs, "; "))
	}

	return summary, nil
}

func main() {
	engine := NewPayrollEngine()

	salaried, _ := NewSalariedEmployee("EMP-01", "Alice Smith", 800000, 50000, 0.20)
	hourly, _ := NewHourlyEmployee("EMP-02", "Bob Jones", 4000, 45.0, 0.15) // 40h reg + 5h OT
	contractor, _ := NewContractor("CON-01", "Carol White", 7500, 30.0)

	team := []Payable{salaried, hourly, contractor}
	summary, err := engine.ProcessBatch(team)
	if err != nil {
		fmt.Printf("Batch error: %v\n", err)
		return
	}

	fmt.Printf("Processed %d employees:\n", summary.TotalEmployees)
	fmt.Printf("Total Gross: $%0.2f\n", float64(summary.TotalGrossCents)/100)
	fmt.Printf("Total Deductions: $%0.2f\n", float64(summary.TotalDeductionsCents)/100)
	fmt.Printf("Total Net: $%0.2f\n", float64(summary.TotalNetCents)/100)
}
