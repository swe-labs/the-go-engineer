package main

import "testing"

func TestAnalyzeProject_HighQuality(t *testing.T) {
	p := Project{
		Name:        "go-url-shortener",
		Description: "A URL shortening service built with Go, Redis, and gRPC",
		Language:    "Go",
		Stars:       245,
		HasTests:    true,
		HasDocs:     true,
		HasCI:       true,
		LinesOfCode: 8500,
	}
	a := AnalyzeProject(p)
	if a.Total < 20 {
		t.Errorf("expected high score for quality project, got %d", a.Total)
	}
}

func TestAnalyzeProject_LowQuality(t *testing.T) {
	p := Project{
		Name:        "hello-world-cli",
		Description: "Simple CLI tool in Go",
		Language:    "Go",
		Stars:       3,
		HasTests:    false,
		HasDocs:     false,
		HasCI:       false,
		LinesOfCode: 120,
	}
	a := AnalyzeProject(p)
	if a.Total > 15 {
		t.Errorf("expected low score for minimal project, got %d", a.Total)
	}
}

func TestAnalyzeProject_AllDimensionsPresent(t *testing.T) {
	p := Project{Name: "test", Description: "test", Language: "Go", Stars: 0, HasTests: false, HasDocs: false, HasCI: false, LinesOfCode: 100}
	a := AnalyzeProject(p)
	if len(a.Scores) != 6 {
		t.Errorf("expected 6 dimensions, got %d", len(a.Scores))
	}
}

func TestAnalyzeProject_CommunityAdoption(t *testing.T) {
	tests := []struct {
		name       string
		stars      int
		wantMin    int
		wantReason string
	}{
		{"high adoption", 500, 5, "significant community adoption"},
		{"moderate adoption", 50, 4, "some adoption"},
		{"low adoption", 5, 3, "initial traction"},
		{"no adoption", 0, 2, "personal or experimental"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := Project{Name: "test", Description: "test", Language: "Go", Stars: tc.stars}
			a := AnalyzeProject(p)
			var realScore Score
			for _, s := range a.Scores {
				if s.Dimension == RealWorldUse {
					realScore = s
					break
				}
			}
			if realScore.Score < tc.wantMin {
				t.Errorf("with %d stars, got score %d, want at least %d", tc.stars, realScore.Score, tc.wantMin)
			}
		})
	}
}
