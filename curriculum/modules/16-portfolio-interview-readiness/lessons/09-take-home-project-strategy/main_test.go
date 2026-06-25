package main

import (
	"strings"
	"testing"
)

func TestProjectScore_TotalScore(t *testing.T) {
	tests := []struct {
		name     string
		score    ProjectScore
		expected float64
	}{
		{
			name: "perfect project",
			score: ProjectScore{
				ProjectName: "Perfect",
				Categories: []ScoreCategory{
					{Name: "Testing", Weight: 1.0, Score: 10},
					{Name: "Doc", Weight: 0.8, Score: 10},
				},
			},
			expected: 100.0,
		},
		{
			name: "average project",
			score: ProjectScore{
				ProjectName: "Average",
				Categories: []ScoreCategory{
					{Name: "Testing", Weight: 1.0, Score: 5},
					{Name: "Doc", Weight: 0.8, Score: 5},
				},
			},
			expected: 50.0,
		},
		{
			name: "no categories",
			score: ProjectScore{
				ProjectName: "Empty",
				Categories:  []ScoreCategory{},
			},
			expected: 0.0,
		},
		{
			name: "single category",
			score: ProjectScore{
				ProjectName: "Minimal",
				Categories: []ScoreCategory{
					{Name: "Testing", Weight: 1.0, Score: 7},
				},
			},
			expected: 70.0,
		},
		{
			name: "uneven weights",
			score: ProjectScore{
				ProjectName: "Uneven",
				Categories: []ScoreCategory{
					{Name: "Testing", Weight: 2.0, Score: 8},
					{Name: "Doc", Weight: 0.5, Score: 4},
				},
			},
			expected: 72.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.score.TotalScore()
			if got != tt.expected {
				t.Errorf("TotalScore() = %.1f, want %.1f", got, tt.expected)
			}
		})
	}
}

func TestProjectScore_Summary(t *testing.T) {
	ps := ProjectScore{
		ProjectName: "Test",
		Categories: []ScoreCategory{
			{Name: "Testing", Weight: 1.0, Score: 8, Comment: "good"},
		},
	}
	s := ps.Summary()
	if s == "" {
		t.Error("Summary() returned empty string")
	}
	if !strings.Contains(s, "Test") {
		t.Error("Summary() should contain project name")
	}
}
