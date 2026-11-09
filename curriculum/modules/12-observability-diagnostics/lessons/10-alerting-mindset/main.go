package main

import (
	"fmt"
	"time"
)

type SLOConfig struct {
	Window      time.Duration
	Target      float64
	ErrorBudget float64
}

type BurnRateAlert struct {
	Severity  string
	Threshold float64
	Window    time.Duration
}

type SLI struct {
	TotalRequests int
	ErrorRequests int
}

func (s SLI) ErrorRate() float64 {
	if s.TotalRequests == 0 {
		return 0
	}
	return float64(s.ErrorRequests) / float64(s.TotalRequests)
}

func (s SLI) SuccessRate() float64 {
	return 1 - s.ErrorRate()
}

func computeBurnRate(errors, total int, since time.Duration, slo SLOConfig) float64 {
	if total == 0 {
		return 0
	}
	actualErrorRate := float64(errors) / float64(total)
	maxAllowedErrorRate := 1 - slo.Target
	if maxAllowedErrorRate == 0 {
		return 0
	}
	return actualErrorRate / maxAllowedErrorRate
}

func evaluateAlert(sli SLI, since time.Duration, slo SLOConfig, alerts []BurnRateAlert) []string {
	var firing []string
	burnRate := computeBurnRate(sli.ErrorRequests, sli.TotalRequests, since, slo)
	for _, alert := range alerts {
		if burnRate >= alert.Threshold {
			firing = append(firing, alert.Severity)
		}
	}
	return firing
}

func main() {
	slo := SLOConfig{
		Window:      30 * 24 * time.Hour,
		Target:      0.999,
		ErrorBudget: 0.001,
	}

	alerts := []BurnRateAlert{
		{Severity: "critical", Threshold: 10, Window: 1 * time.Hour},
		{Severity: "warning", Threshold: 2, Window: 6 * time.Hour},
	}

	scenarios := []struct {
		name   string
		total  int
		errors int
		since  time.Duration
	}{
		{"normal operation", 10000, 2, 1 * time.Hour},
		{"brief spike", 1000, 50, 1 * time.Hour},
		{"sustained degradation", 100000, 500, 6 * time.Hour},
		{"major outage", 5000, 200, 30 * time.Minute},
	}

	fmt.Printf("SLO: %.1f%%, Error Budget: %.1f%%\n\n", slo.Target*100, slo.ErrorBudget*100)

	for _, s := range scenarios {
		sli := SLI{TotalRequests: s.total, ErrorRequests: s.errors}
		rate := sli.ErrorRate() * 100
		burn := computeBurnRate(s.errors, s.total, s.since, slo)
		firing := evaluateAlert(sli, s.since, slo, alerts)

		fmt.Printf("[%s]\n", s.name)
		fmt.Printf("  Requests: %d, Errors: %d, Error rate: %.2f%%\n", s.total, s.errors, rate)
		fmt.Printf("  Burn rate: %.1fx\n", burn)
		if len(firing) == 0 {
			fmt.Println("  Status: OK (no alerts)")
		} else {
			fmt.Printf("  Alerts firing: %v\n", firing)
		}
		fmt.Println()
	}
}
