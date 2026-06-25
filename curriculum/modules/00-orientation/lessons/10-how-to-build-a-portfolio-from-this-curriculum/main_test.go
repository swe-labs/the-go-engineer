package main

import "testing"

func TestScoreReadiness_CompleteProject(t *testing.T) {
	p := []PortfolioProject{
		{
			Name:        "Test Project",
			Description: "A test project",
			Evidence:    []string{"main.go", "test.go"},
			Status:      "complete",
		},
	}
	scores := scoreReadiness(p)
	if scores["Test Project"] != 5 {
		t.Errorf("complete project should score 5, got %d", scores["Test Project"])
	}
}

func TestScoreReadiness_PlannedProjectNoEvidence(t *testing.T) {
	p := []PortfolioProject{
		{
			Name:        "Future Project",
			Description: "",
			Evidence:    nil,
			Status:      "planned",
		},
	}
	scores := scoreReadiness(p)
	if scores["Future Project"] != 0 {
		t.Errorf("planned project with no info should score 0, got %d", scores["Future Project"])
	}
}

func TestScoreReadiness_InProgressWithEvidence(t *testing.T) {
	p := []PortfolioProject{
		{
			Name:        "API Server",
			Description: "A REST API",
			Evidence:    []string{"design.md"},
			Status:      "in-progress",
		},
	}
	scores := scoreReadiness(p)
	if scores["API Server"] < 2 {
		t.Errorf("in-progress with description should score at least 2, got %d", scores["API Server"])
	}
}

func TestScoreReadiness_MultipleProjects(t *testing.T) {
	projects := []PortfolioProject{
		{Name: "A", Description: "desc", Evidence: []string{"e1", "e2"}, Status: "complete"},
		{Name: "B", Description: "", Evidence: nil, Status: "planned"},
		{Name: "C", Description: "desc", Evidence: []string{"e1"}, Status: "in-progress"},
	}
	scores := scoreReadiness(projects)
	if len(scores) != 3 {
		t.Errorf("expected 3 projects scored, got %d", len(scores))
	}
	if scores["A"] != 5 {
		t.Errorf("project A should score 5, got %d", scores["A"])
	}
	if scores["B"] != 0 {
		t.Errorf("project B should score 0, got %d", scores["B"])
	}
}

func TestScoreReadiness_EvidenceCapped(t *testing.T) {
	p := []PortfolioProject{
		{
			Name:        "Big Project",
			Description: "desc",
			Evidence:    []string{"a", "b", "c", "d", "e"},
			Status:      "complete",
		},
	}
	scores := scoreReadiness(p)
	if scores["Big Project"] > 5 {
		t.Errorf("score should be capped at 5, got %d", scores["Big Project"])
	}
}

func TestScoreReadiness_EmptyProjects(t *testing.T) {
	scores := scoreReadiness(nil)
	if len(scores) != 0 {
		t.Errorf("empty projects should return empty map, got %d items", len(scores))
	}
}
