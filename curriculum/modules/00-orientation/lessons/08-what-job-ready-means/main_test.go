package main

import (
	"strings"
	"testing"
)

func TestIdentifyGaps_NoGaps(t *testing.T) {
	jr := JobReadiness{
		Categories: []SkillCategory{
			{Name: "Implementation", Evidence: []string{"built an app"}},
			{Name: "Testing", Evidence: []string{"wrote tests"}},
		},
	}
	gaps := identifyGaps(jr)
	if len(gaps) != 0 {
		t.Errorf("expected no gaps, got %v", gaps)
	}
}

func TestIdentifyGaps_AllGaps(t *testing.T) {
	jr := JobReadiness{
		Categories: []SkillCategory{
			{Name: "Implementation"},
			{Name: "Testing"},
		},
	}
	gaps := identifyGaps(jr)
	if len(gaps) != 2 {
		t.Errorf("expected 2 gaps, got %d: %v", len(gaps), gaps)
	}
}

func TestIdentifyGaps_SomeGaps(t *testing.T) {
	jr := JobReadiness{
		Categories: []SkillCategory{
			{Name: "Implementation", Evidence: []string{"app"}},
			{Name: "Testing", Evidence: []string{}},
			{Name: "Debugging", Evidence: []string{"fix"}},
			{Name: "Communication", Evidence: []string{}},
		},
	}
	gaps := identifyGaps(jr)
	if len(gaps) != 2 {
		t.Fatalf("expected 2 gaps, got %d: %v", len(gaps), gaps)
	}
	if gaps[0] != "Testing" || gaps[1] != "Communication" {
		t.Errorf("expected [Testing Communication], got %v", gaps)
	}
}

func TestReadinessSummary_Full(t *testing.T) {
	jr := JobReadiness{
		Categories: []SkillCategory{
			{Name: "A", Evidence: []string{"x"}},
			{Name: "B", Evidence: []string{"y"}},
		},
	}
	summary := readinessSummary(jr)
	if !strings.Contains(summary, "job-ready") {
		t.Errorf("expected job-ready message, got: %s", summary)
	}
}

func TestReadinessSummary_Partial(t *testing.T) {
	jr := JobReadiness{
		Categories: []SkillCategory{
			{Name: "A", Evidence: []string{"x"}},
			{Name: "B", Evidence: []string{}},
			{Name: "C", Evidence: []string{"z"}},
			{Name: "D", Evidence: []string{}},
		},
	}
	summary := readinessSummary(jr)
	if !strings.Contains(summary, "2/4") {
		t.Errorf("expected '2/4' in summary, got: %s", summary)
	}
}

func TestReadinessSummary_Empty(t *testing.T) {
	jr := JobReadiness{}
	summary := readinessSummary(jr)
	if !strings.Contains(summary, "No skill categories defined") {
		t.Errorf("expected empty message, got: %s", summary)
	}
}

func TestReadinessSummary_NoEvidence(t *testing.T) {
	jr := JobReadiness{
		Categories: []SkillCategory{
			{Name: "A"},
			{Name: "B"},
		},
	}
	summary := readinessSummary(jr)
	if !strings.Contains(summary, "0/2") {
		t.Errorf("expected '0/2' in summary, got: %s", summary)
	}
}
