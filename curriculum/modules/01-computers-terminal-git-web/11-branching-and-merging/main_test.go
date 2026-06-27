package main

import (
	"strings"
	"testing"
)

func TestNewSimRepo(t *testing.T) {
	r := NewSimRepo()
	if r == nil {
		t.Fatal("NewSimRepo() returned nil")
	}
	if len(r.Commits) != 0 {
		t.Errorf("expected 0 commits, got %d", len(r.Commits))
	}
	if len(r.Branches) != 0 {
		t.Errorf("expected 0 branches, got %d", len(r.Branches))
	}
}

func TestCommitAndLog(t *testing.T) {
	r := NewSimRepo()
	r.Commit("main", "Initial commit")
	r.Commit("main", "Second commit")

	log := r.Log("main")
	if len(log) != 2 {
		t.Fatalf("expected 2 log entries, got %d", len(log))
	}
	if !strings.Contains(log[0], "Second commit") {
		t.Errorf("expected first log entry to contain 'Second commit', got %q", log[0])
	}
	if !strings.Contains(log[1], "Initial commit") {
		t.Errorf("expected second log entry to contain 'Initial commit', got %q", log[1])
	}
}

func TestBranchAndMergeFastForward(t *testing.T) {
	r := NewSimRepo()
	r.Commit("main", "Initial commit")
	r.Branch("feature", "main")
	r.Commit("feature", "Feature work")

	hash, ok := r.Merge("feature", "main")
	if !ok {
		t.Fatal("expected successful fast-forward merge")
	}
	if hash != "" {
		t.Logf("Fast-forward merge completed")
	}
	log := r.Log("main")
	if !strings.Contains(log[0], "Feature work") {
		t.Errorf("expected feature work in main log after merge, got %q", log[0])
	}
}
