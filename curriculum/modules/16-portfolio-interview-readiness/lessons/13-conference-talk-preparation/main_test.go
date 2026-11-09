package main

import (
	"strings"
	"testing"
)

func TestEvaluateProposal_Basic(t *testing.T) {
	p := TalkProposal{
		Title:       "Test Talk",
		Abstract:    "This talk covers real-world lessons from Go concurrency. First we'll cover channels, then goroutines, finally patterns. You will learn practical patterns.",
		DurationMin: 30,
	}
	s := EvaluateProposal(p)
	if s.Total <= 0 {
		t.Errorf("Total score should be positive, got %d", s.Total)
	}
}

func TestEvaluateProposal_EmptyAbstract(t *testing.T) {
	p := TalkProposal{
		Title:    "Empty",
		Abstract: "",
	}
	s := EvaluateProposal(p)
	if s.Abstract != 3 {
		t.Errorf("Abstract score for empty should be 3, got %d", s.Abstract)
	}
}

func TestEvaluateProposal_Novelty(t *testing.T) {
	tests := []struct {
		name     string
		abstract string
		expected int
	}{
		{"high novelty", "This novel deep dive covers unexpected internals of Go channels with real-world lessons", 9},
		{"medium novelty", "Real-world patterns for Go channels", 6},
		{"low novelty", "How to use fmt.Println", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := TalkProposal{Title: "Test", Abstract: tt.abstract}
			s := EvaluateProposal(p)
			if s.Novelty != tt.expected {
				t.Errorf("Novelty = %d, want %d", s.Novelty, tt.expected)
			}
		})
	}
}

func TestEvaluateProposal_Feedback(t *testing.T) {
	p := TalkProposal{
		Title:       "Bad Talk",
		Abstract:    "a",
		DurationMin: 20,
		HasDemo:     true,
	}
	s := EvaluateProposal(p)
	if len(s.Feedback) == 0 {
		t.Error("expected feedback for bad proposal")
	}
}

func TestProposalScore_Summary(t *testing.T) {
	p := TalkProposal{Title: "SummaryTest", Abstract: "good abstract with real-world lessons from Go"}
	s := EvaluateProposal(p)
	summary := s.Summary()
	if !strings.Contains(summary, "SummaryTest") {
		t.Error("summary should contain title")
	}
	if !strings.Contains(summary, "Total Score") {
		t.Error("summary should contain total score")
	}
}
