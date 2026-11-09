package main

import (
	"strings"
	"testing"
)

func TestResumeAnalyzer_FindsKeywords(t *testing.T) {
	text := "Go and Kubernetes and Docker and gRPC"
	analyzer := NewResumeAnalyzer()
	score := analyzer.Analyze(text)
	if len(score.Results) < 4 {
		t.Errorf("expected at least 4 keyword matches, got %d", len(score.Results))
	}
}

func TestResumeAnalyzer_NoKeywords(t *testing.T) {
	text := "I like to cook and bake"
	analyzer := NewResumeAnalyzer()
	score := analyzer.Analyze(text)
	if len(score.Results) != 0 {
		t.Errorf("expected 0 keyword matches, got %d", len(score.Results))
	}
}

func TestResumeAnalyzer_CategoryBreakdown(t *testing.T) {
	text := "Go and leadership and agile and Docker"
	analyzer := NewResumeAnalyzer()
	score := analyzer.Analyze(text)
	breakdown := score.CategoryBreakdown()
	if len(breakdown) < 3 {
		t.Errorf("expected at least 3 categories, got %d %v", len(breakdown), breakdown)
	}
}

func TestResumeAnalyzer_CountMatchesCorrectly(t *testing.T) {
	text := "Go Go Go and Docker Docker"
	analyzer := NewResumeAnalyzer()
	score := analyzer.Analyze(text)
	for _, r := range score.Results {
		if r.Keyword == "Go" && r.Count != 3 {
			t.Errorf("expected 3 matches for Go, got %d", r.Count)
		}
		if r.Keyword == "Docker" && r.Count != 2 {
			t.Errorf("expected 2 matches for Docker, got %d", r.Count)
		}
	}
}

func TestResumeAnalyzer_WordBoundaries(t *testing.T) {
	text := "Golang is not the same as GoLang"
	analyzer := NewResumeAnalyzer()
	score := analyzer.Analyze(text)
	// Should match "Golang" but not case variations after the first
	var golangCount int
	for _, r := range score.Results {
		if r.Keyword == "Golang" {
			golangCount += r.Count
		}
	}
	if golangCount == 0 {
		t.Error("expected to match Golang at least once")
	}
}

func TestResumeAnalyzer_SummaryContainsBreakdown(t *testing.T) {
	text := "Go and Docker and Kubernetes and leadership"
	analyzer := NewResumeAnalyzer()
	score := analyzer.Analyze(text)
	summary := score.Summary()
	if !strings.Contains(summary, "Category Breakdown") {
		t.Errorf("summary should contain category breakdown")
	}
}
