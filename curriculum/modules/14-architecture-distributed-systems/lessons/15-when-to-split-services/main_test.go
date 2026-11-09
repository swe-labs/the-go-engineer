package main

import "testing"

func TestAnalyzeSplitContextSmallTeam(t *testing.T) {
	ctx := BoundedContext{
		Name:     "Billing",
		OwnsData: []string{"invoices"},
		TeamSize: 3,
	}
	a := AnalyzeSplitContext(ctx)
	if a.RecommendedAction != "keep as part of monolith" {
		t.Errorf("expected monolith, got %s", a.RecommendedAction)
	}
}

func TestAnalyzeSplitContextLargeTeam(t *testing.T) {
	ctx := BoundedContext{
		Name:     "Orders",
		OwnsData: []string{"orders", "items", "payments", "refunds", "shipments"},
		TeamSize: 10,
	}
	a := AnalyzeSplitContext(ctx)
	if len(a.Reasons) == 0 {
		t.Fatal("expected reasons for split")
	}
}

func TestAnalyzeSplitContextManyEntities(t *testing.T) {
	ctx := BoundedContext{
		Name:     "Catalogue",
		OwnsData: []string{"a", "b", "c", "d", "e", "f"},
		TeamSize: 4,
	}
	a := AnalyzeSplitContext(ctx)
	if a.RecommendedAction == "keep as part of monolith" {
		t.Log("single reason may not trigger split")
	}
}

func TestEvaluateCouplingHigh(t *testing.T) {
	c := ServiceCoupling{
		ServiceA:   "A",
		ServiceB:   "B",
		SharedData: []string{"shared_table"},
		CallPerSec: 200,
	}
	result := EvaluateCoupling(c)
	if result == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestEvaluateCouplingLow(t *testing.T) {
	c := ServiceCoupling{
		ServiceA:   "A",
		ServiceB:   "B",
		SharedData: nil,
		CallPerSec: 5,
	}
	result := EvaluateCoupling(c)
	if result == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestStrategiesExist(t *testing.T) {
	if len(Strategies) != 3 {
		t.Errorf("expected 3 strategies, got %d", len(Strategies))
	}
	if _, ok := Strategies["strangler-fig"]; !ok {
		t.Error("expected strangler-fig strategy")
	}
	if _, ok := Strategies["database-separation"]; !ok {
		t.Error("expected database-separation strategy")
	}
	if _, ok := Strategies["event-driven-split"]; !ok {
		t.Error("expected event-driven-split strategy")
	}
}

func TestAnalyzeSplitContextTable(t *testing.T) {
	tests := []struct {
		name         string
		ctx          BoundedContext
		wantMonolith bool
	}{
		{
			name:         "small_team_few_entities",
			ctx:          BoundedContext{Name: "Auth", OwnsData: []string{"users"}, TeamSize: 2},
			wantMonolith: true,
		},
		{
			name:         "big_team_many_entities",
			ctx:          BoundedContext{Name: "Analytics", OwnsData: []string{"a", "b", "c", "d", "e", "f"}, TeamSize: 12},
			wantMonolith: false,
		},
		{
			name:         "medium_team",
			ctx:          BoundedContext{Name: "Notifications", OwnsData: []string{"templates", "logs"}, TeamSize: 5},
			wantMonolith: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a := AnalyzeSplitContext(tc.ctx)
			isMonolith := a.RecommendedAction == "keep as part of monolith"
			if isMonolith != tc.wantMonolith {
				t.Errorf("wantMonolith=%v, got action=%q", tc.wantMonolith, a.RecommendedAction)
			}
		})
	}
}
