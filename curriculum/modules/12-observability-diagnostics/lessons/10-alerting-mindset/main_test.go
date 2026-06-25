package main

import (
	"testing"
	"time"
)

func TestSLIErrorRate(t *testing.T) {
	tests := []struct {
		name            string
		total           int
		errors          int
		expectedRate    float64
		expectedSuccess float64
	}{
		{"no errors", 1000, 0, 0.0, 1.0},
		{"1% errors", 1000, 10, 0.01, 0.99},
		{"50% errors", 100, 50, 0.5, 0.5},
		{"all errors", 100, 100, 1.0, 0.0},
		{"zero requests", 0, 0, 0.0, 1.0},
	}

	for _, tc := range tests {
		sli := SLI{TotalRequests: tc.total, ErrorRequests: tc.errors}
		if got := sli.ErrorRate(); got != tc.expectedRate {
			t.Errorf("%s: ErrorRate() = %f, want %f", tc.name, got, tc.expectedRate)
		}
		if got := sli.SuccessRate(); got != tc.expectedSuccess {
			t.Errorf("%s: SuccessRate() = %f, want %f", tc.name, got, tc.expectedSuccess)
		}
	}
}

func closeEnough(a, b float64) bool {
	const epsilon = 1e-9
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < epsilon
}

func TestComputeBurnRate(t *testing.T) {
	slo := SLOConfig{
		Window:      30 * 24 * time.Hour,
		Target:      0.99,
		ErrorBudget: 0.01,
	}

	tests := []struct {
		name         string
		total        int
		errors       int
		since        time.Duration
		expectedBurn float64
	}{
		{"within budget", 10000, 50, 1 * time.Hour, 0.5},
		{"at budget", 10000, 100, 1 * time.Hour, 1.0},
		{"double budget", 10000, 200, 1 * time.Hour, 2.0},
		{"ten times budget", 10000, 1000, 1 * time.Hour, 10.0},
		{"zero traffic", 0, 0, 1 * time.Hour, 0.0},
	}

	for _, tc := range tests {
		got := computeBurnRate(tc.errors, tc.total, tc.since, slo)
		if !closeEnough(got, tc.expectedBurn) {
			t.Errorf("%s: computeBurnRate() = %f, want %f", tc.name, got, tc.expectedBurn)
		}
	}
}

func TestEvaluateAlert(t *testing.T) {
	slo := SLOConfig{
		Window:      30 * 24 * time.Hour,
		Target:      0.99,
		ErrorBudget: 0.01,
	}

	alerts := []BurnRateAlert{
		{Severity: "critical", Threshold: 10, Window: 1 * time.Hour},
		{Severity: "warning", Threshold: 2, Window: 6 * time.Hour},
	}

	tests := []struct {
		name           string
		total          int
		errors         int
		since          time.Duration
		expectedAlerts []string
	}{
		{"normal", 10000, 50, 1 * time.Hour, []string{}},
		{"warning level", 10000, 300, 1 * time.Hour, []string{"warning"}},
		{"critical level", 10000, 1500, 1 * time.Hour, []string{"critical", "warning"}},
	}

	for _, tc := range tests {
		sli := SLI{TotalRequests: tc.total, ErrorRequests: tc.errors}
		got := evaluateAlert(sli, tc.since, slo, alerts)

		if len(got) != len(tc.expectedAlerts) {
			t.Errorf("%s: got %d alerts, want %d: %v", tc.name, len(got), len(tc.expectedAlerts), got)
			continue
		}
		for i, severity := range tc.expectedAlerts {
			if got[i] != severity {
				t.Errorf("%s: alert[%d] = %s, want %s", tc.name, i, got[i], severity)
			}
		}
	}
}

func TestBudgetTracker(t *testing.T) {
	slo := SLOConfig{
		Window:      30 * 24 * time.Hour,
		Target:      0.99,
		ErrorBudget: 0.01,
	}

	total := 10000
	errors := 50
	burn := computeBurnRate(errors, total, 1*time.Hour, slo)

	if !closeEnough(burn, 0.5) {
		t.Errorf("expected burn rate 0.5 for 0.5%% errors with 1%% budget, got %f", burn)
	}
}
