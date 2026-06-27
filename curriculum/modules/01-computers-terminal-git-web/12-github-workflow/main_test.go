package main

import (
	"testing"
)

func TestStartPR(t *testing.T) {
	pr := StartPR("Fix bug", "Fixes the nil pointer issue", "bugfix", "main")
	if pr == nil {
		t.Fatal("StartPR returned nil")
	}
	if pr.Title != "Fix bug" {
		t.Errorf("expected title 'Fix bug', got %q", pr.Title)
	}
	if pr.Status != PROpen {
		t.Errorf("expected status PROpen, got %v", pr.Status)
	}
}

func TestPRWorkflow(t *testing.T) {
	repo := &GitHubRepo{Name: "test", Branches: map[string]bool{"main": true}, Protected: false}
	pr := StartPR("Add feature", "New feature", "feature", "main")

	pr.Approve()
	ok, _ := pr.Merge(repo)
	if !ok {
		t.Fatal("expected successful merge on approved PR")
	}
	if pr.Status != PRMerged {
		t.Errorf("expected PRMerged, got %v", pr.Status)
	}
}

func TestPRMergeRequiresApproval(t *testing.T) {
	repo := &GitHubRepo{Name: "test", Branches: map[string]bool{"main": true}, Protected: false}
	pr := StartPR("Unreviewed change", "Quick fix", "hotfix", "main")

	ok, msg := pr.Merge(repo)
	if ok {
		t.Fatal("expected merge to fail without approval")
	}
	if msg != "PR must be approved before merging" {
		t.Errorf("unexpected message: %q", msg)
	}
}

func TestPRStatusTransitions(t *testing.T) {
	pr := StartPR("Test", "Desc", "src", "tgt")
	if pr.Status != PROpen {
		t.Errorf("expected PROpen, got %v", pr.Status)
	}
	pr.RequestChanges()
	if pr.Status != PRChangesRequested {
		t.Errorf("expected PRChangesRequested, got %v", pr.Status)
	}
	pr.Approve()
	if pr.Status != PRApproved {
		t.Errorf("expected PRApproved, got %v", pr.Status)
	}
	pr.Close()
	if pr.Status != PRApproved {
		t.Errorf("Close should not affect merged or approved PRs; got %v", pr.Status)
	}
}
