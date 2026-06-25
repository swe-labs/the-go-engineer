package main

import "testing"

func TestDefaultEstimator_BaseValues(t *testing.T) {
	p := ProjectScope{Name: "tiny", NumFeatures: 0, TeamSize: 1}
	r := EstimateProject(p, DefaultEstimator())
	if r.EstimatedWeeks < 1 || r.EstimatedWeeks > 4 {
		t.Errorf("expected tiny project ~2 weeks, got %.0f", r.EstimatedWeeks)
	}
}

func TestEstimateProject_ComplexIncreases(t *testing.T) {
	simple := ProjectScope{Name: "simple", NumFeatures: 1, TeamSize: 1}
	complex := ProjectScope{Name: "complex", NumFeatures: 10, NumIntegrations: 4, HasAuth: true, HasDatabase: true, HasExternalAPI: true, NeedsUI: true, TeamSize: 4}

	est := DefaultEstimator()
	r1 := EstimateProject(simple, est)
	r2 := EstimateProject(complex, est)

	if r2.EstimatedWeeks <= r1.EstimatedWeeks {
		t.Errorf("expected complex project to take longer than simple: %0.f <= %0.f", r2.EstimatedWeeks, r1.EstimatedWeeks)
	}
}

func TestEstimateProject_ComplexityLevels(t *testing.T) {
	tests := []struct {
		name    string
		scope   ProjectScope
		wantMin ComplexityLevel
		wantMax ComplexityLevel
	}{
		{"trivial", ProjectScope{Name: "t", NumFeatures: 0, TeamSize: 1}, Trivial, Easy},
		{"moderate", ProjectScope{Name: "m", NumFeatures: 3, HasDatabase: true, TeamSize: 1}, Easy, Moderate},
		{"complex", ProjectScope{Name: "c", NumFeatures: 8, HasAuth: true, HasDatabase: true, HasExternalAPI: true, NeedsUI: true, TeamSize: 3}, Complex, VeryComplex},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := EstimateProject(tc.scope, DefaultEstimator())
			if r.Complexity < tc.wantMin || r.Complexity > tc.wantMax {
				t.Errorf("expected complexity between %v and %v, got %v", tc.wantMin, tc.wantMax, r.Complexity)
			}
		})
	}
}

func TestEstimateProject_RiskFactorsAdded(t *testing.T) {
	p := ProjectScope{Name: "test", HasAuth: true, HasDatabase: true, HasExternalAPI: true, NeedsUI: true}
	r := EstimateProject(p, DefaultEstimator())
	if len(r.RiskFactors) < 3 {
		t.Errorf("expected at least 3 risk factors, got %d", len(r.RiskFactors))
	}
}
