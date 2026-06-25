package main

import "testing"

func TestEvaluateServiceSmallTeam(t *testing.T) {
	dec := EvaluateService("Auth", 3, 100, true)
	if dec.Recommendation != "start with monolith" {
		t.Errorf("expected monolith recommendation for small team, got %q", dec.Recommendation)
	}
}

func TestEvaluateServiceLargeTeamHighThroughput(t *testing.T) {
	dec := EvaluateService("Orders", 10, 50000, false)
	if dec.Recommendation != "consider microservice" {
		t.Errorf("expected microservice recommendation for large team, got %q", dec.Recommendation)
	}
}

func TestEvaluateServiceLatencySensitive(t *testing.T) {
	sensitive := EvaluateService("Payments", 6, 5000, true)
	notSensitive := EvaluateService("Payments", 6, 5000, false)
	if sensitive.Score >= notSensitive.Score {
		t.Error("latency-sensitive should have lower score than non-sensitive")
	}
}

func TestEstimateRequestLatency(t *testing.T) {
	est := EstimateRequestLatency(1, 50, 5)
	if est.MonolithMs < 0 || est.MicroserviceMs < 0 {
		t.Error("latency should be non-negative")
	}
}

func TestTradeoffsExist(t *testing.T) {
	if len(MicroserviceTradeoffs) == 0 {
		t.Fatal("expected tradeoffs list")
	}
	names := make(map[string]bool)
	for _, tr := range MicroserviceTradeoffs {
		if tr.Name == "" {
			t.Error("tradeoff name must not be empty")
		}
		if names[tr.Name] {
			t.Errorf("duplicate tradeoff name: %s", tr.Name)
		}
		names[tr.Name] = true
	}
}

func TestEvaluateServiceTable(t *testing.T) {
	tests := []struct {
		name       string
		teamSize   int
		throughput int
		latSen     bool
		wantRec    string
	}{
		{
			name:       "small_team_low_throughput",
			teamSize:   2,
			throughput: 100,
			latSen:     true,
			wantRec:    "start with monolith",
		},
		{
			name:       "large_team_high_throughput",
			teamSize:   12,
			throughput: 100000,
			latSen:     false,
			wantRec:    "consider microservice",
		},
		{
			name:       "medium_team",
			teamSize:   6,
			throughput: 5000,
			latSen:     false,
			wantRec:    "start with monolith",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dec := EvaluateService(tc.name, tc.teamSize, tc.throughput, tc.latSen)
			if dec.Recommendation != tc.wantRec {
				t.Errorf("want %q, got %q (score=%d)", tc.wantRec, dec.Recommendation, dec.Score)
			}
		})
	}
}
