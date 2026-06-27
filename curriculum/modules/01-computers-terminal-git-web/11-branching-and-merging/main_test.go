package main

import (
	"strings"
	"testing"
)

func TestInitRepo(t *testing.T) {
	initRepo()
	if commitSeq != 0 {
		t.Errorf("expected commitSeq 0 after init, got %d", commitSeq)
	}
}

func TestCommitAndLog(t *testing.T) {
	initRepo()
	mainBranch, _ := commit("", "Initial commit")
	mainBranch, _ = commit(mainBranch, "Second commit")

	log := logBranch(mainBranch)
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
	initRepo()
	mainBranch, _ := commit("", "Initial commit")
	featureBranch := branchFrom(mainBranch)
	featureBranch, _ = commit(featureBranch, "Feature work")

	merged, ok := merge(featureBranch, mainBranch)
	if !ok {
		t.Fatal("expected successful fast-forward merge")
	}
	mainBranch = merged
	log := logBranch(mainBranch)
	if !strings.Contains(log[0], "Feature work") {
		t.Errorf("expected feature work in main log after merge, got %q", log[0])
	}
}
