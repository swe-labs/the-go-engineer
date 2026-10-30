package main

import (
	"testing"
)

func TestDiffAndReviewEvaluation(t *testing.T) {
	diff := []DiffLine{
		{LineNumber: 1, Type: " ", Content: "package main"},
		{LineNumber: 2, Type: "-", Content: "var counter = 0"},
		{LineNumber: 2, Type: "+", Content: "var counter atomic.Int64"},
	}

	if len(diff) != 3 {
		t.Fatalf("expected 3 diff lines, got %d", len(diff))
	}
	if diff[1].Type != "-" || diff[2].Type != "+" {
		t.Errorf("diff change markers mismatch")
	}

	// Non-blocking suggestions
	comments := []ReviewComment{
		{
			LineNumber: 2,
			Reviewer:   "reviewer1",
			Comment:    "Good use of sync/atomic!",
			IsBlocking: false,
		},
	}

	hasBlocking := false
	for _, c := range comments {
		if c.IsBlocking {
			hasBlocking = true
		}
	}
	if hasBlocking {
		t.Errorf("expected no blocking comments")
	}

	// Add blocking comment
	comments = append(comments, ReviewComment{
		LineNumber: 2,
		Reviewer:   "security-auditor",
		Comment:    "Memory alignment check required for 32-bit platforms",
		IsBlocking: true,
	})

	for _, c := range comments {
		if c.IsBlocking {
			hasBlocking = true
			break
		}
	}
	if !hasBlocking {
		t.Errorf("expected blocking comment to be flagged")
	}
}

func TestMainExecution(t *testing.T) {
	main()
}

