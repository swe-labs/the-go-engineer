package main

import "testing"

func TestAnalyzeGaps_Basic(t *testing.T) {
	plan := CareerPlan{
		Name:  "Test",
		Role:  "senior",
		Track: "ic",
		Skills: []Skill{
			{Name: "Go", Category: "language", Current: 5, Desired: 9, Priority: "high"},
			{Name: "Testing", Category: "testing", Current: 8, Desired: 9, Priority: "medium"},
		},
	}
	ga := AnalyzeGaps(plan)
	if ga.TotalGap != 5 {
		t.Errorf("TotalGap = %d, want 5", ga.TotalGap)
	}
	if len(ga.Gaps) != 2 {
		t.Errorf("expected 2 gaps, got %d", len(ga.Gaps))
	}
}

func TestAnalyzeGaps_NoGaps(t *testing.T) {
	plan := CareerPlan{
		Name:  "Perfect",
		Role:  "senior",
		Track: "ic",
		Skills: []Skill{
			{Name: "Go", Category: "language", Current: 10, Desired: 10, Priority: "high"},
		},
	}
	ga := AnalyzeGaps(plan)
	if ga.TotalGap != 0 {
		t.Errorf("TotalGap = %d, want 0", ga.TotalGap)
	}
	if ga.ReadinessPct != 100.0 {
		t.Errorf("ReadinessPct = %.0f, want 100", ga.ReadinessPct)
	}
}

func TestAnalyzeGaps_EmptySkills(t *testing.T) {
	plan := CareerPlan{
		Name:   "Empty",
		Role:   "junior",
		Track:  "ic",
		Skills: []Skill{},
	}
	ga := AnalyzeGaps(plan)
	if ga.TotalGap != 0 {
		t.Errorf("TotalGap = %d, want 0", ga.TotalGap)
	}
	if ga.ReadinessPct != 0.0 {
		t.Errorf("ReadinessPct = %.0f, want 0", ga.ReadinessPct)
	}
}

func TestComputeUrgency(t *testing.T) {
	tests := []struct {
		gap      int
		priority string
		want     string
	}{
		{5, "high", "critical"},
		{6, "high", "critical"},
		{3, "high", "high"},
		{1, "high", "medium"},
		{0, "high", "none"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := computeUrgency(tt.gap, tt.priority)
			if got != tt.want {
				t.Errorf("computeUrgency(%d, %s) = %s, want %s", tt.gap, tt.priority, got, tt.want)
			}
		})
	}
}

func TestGenerateRoadmap(t *testing.T) {
	plan := CareerPlan{
		Name:  "Test User",
		Role:  "mid",
		Track: "ic",
		Skills: []Skill{
			{Name: "Go", Category: "language", Current: 6, Desired: 8, Priority: "high"},
		},
	}
	ga := AnalyzeGaps(plan)
	roadmap := GenerateRoadmap(ga)
	if len(roadmap) < 50 {
		t.Error("roadmap too short")
	}
}
