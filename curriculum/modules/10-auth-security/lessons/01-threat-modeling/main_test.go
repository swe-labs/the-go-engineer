package main

import (
	"testing"
)

func TestBuildThreatModel_ReturnsThreats(t *testing.T) {
	model := BuildThreatModel(
		[]string{"POST /login"},
		[]string{"password_hash"},
	)
	if len(model) == 0 {
		t.Fatal("expected threats, got none")
	}
	for ep, threats := range model {
		if len(threats) == 0 {
			t.Errorf("endpoint %s has no threats", ep)
		}
		for _, thr := range threats {
			if thr.Category == "" {
				t.Errorf("threat %s has empty category", thr.ID)
			}
			if thr.Severity < 1 || thr.Severity > 10 {
				t.Errorf("threat %s severity %d out of range", thr.ID, thr.Severity)
			}
		}
	}
}

func TestBuildThreatModel_MultipleEndpoints(t *testing.T) {
	model := BuildThreatModel(
		[]string{"POST /login", "GET /documents", "DELETE /admin/users"},
		[]string{"password_hash", "document_content", "admin_session"},
	)
	expectedEndpoints := 3
	if len(model) != expectedEndpoints {
		t.Errorf("expected %d endpoints, got %d", expectedEndpoints, len(model))
	}
}

func TestAnalyzeThreats(t *testing.T) {
	assets := []Asset{
		{ID: "a1", OwnerID: "u1"},
		{ID: "a2", OwnerID: "u2"},
	}
	threats := analyzeThreats(assets)
	if len(threats) != 6 {
		t.Errorf("expected 6 threats for 2 assets, got %d", len(threats))
	}
}

func TestPrioritizeBySeverity(t *testing.T) {
	threats := []Threat{
		{ID: "low", Severity: 1},
		{ID: "high", Severity: 10},
		{ID: "mid", Severity: 5},
	}
	ordered := prioritizeBySeverity(threats)
	if ordered[0].ID != "high" || ordered[2].ID != "low" {
		t.Errorf("expected high, mid, low; got %s, %s, %s", ordered[0].ID, ordered[1].ID, ordered[2].ID)
	}
}
