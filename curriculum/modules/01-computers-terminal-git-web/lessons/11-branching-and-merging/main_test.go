package main

import (
	"fmt"
	"testing"
)

func TestBranch_PointerMovementAndHistory(t *testing.T) {
	c0 := &Commit{
		Hash:    "c000000",
		Message: "Initial commit",
		Changes: map[string]string{"README.md": "v1.0"},
	}

	mainBranch := &Branch{Name: "main", Tip: c0}
	featureBranch := &Branch{Name: "feature", Tip: mainBranch.Tip}

	if featureBranch.Tip != mainBranch.Tip {
		t.Errorf("new branch should point to current HEAD commit")
	}

	// Commit on feature branch
	c1 := &Commit{
		Hash:    "c111111",
		Message: "Add feature logic",
		Parent:  c0,
		Changes: map[string]string{"feature.go": "package feature"},
	}
	featureBranch.Tip = c1

	if mainBranch.Tip.Hash != "c000000" {
		t.Errorf("main branch pointer should remain untouched")
	}
	if featureBranch.Tip.Hash != "c111111" {
		t.Errorf("feature branch pointer should advance to new commit")
	}
	if featureBranch.Tip.Parent != c0 {
		t.Errorf("feature tip parent should point back to c0")
	}
}

func TestConflictDetectionLogic(t *testing.T) {
	baseCommit := &Commit{
		Hash:    "base123",
		Message: "Base commit",
		Changes: map[string]string{"app.go": "original", "config.go": "env: dev"},
	}

	// Case 1: Non-conflicting edits (different files modified)
	branchAChanges := map[string]string{"app.go": "modified in A"}
	branchBChanges := map[string]string{"config.go": "env: prod"}

	hasConflict := false
	for file := range branchAChanges {
		if _, touchedInB := branchBChanges[file]; touchedInB {
			hasConflict = true
		}
	}
	if hasConflict {
		t.Errorf("expected clean merge when orthogonal files are modified")
	}

	// Case 2: Conflicting edits (same file modified on both branches)
	branchCChanges := map[string]string{"app.go": "modified differently in C"}
	for file := range branchAChanges {
		if _, touchedInC := branchCChanges[file]; touchedInC {
			hasConflict = true
		}
	}
	if !hasConflict {
		t.Errorf("expected conflict detected when same file modified concurrently")
	}

	// Verify merge commit generation
	mergeCommit := &Commit{
		Hash:    "merge999",
		Message: fmt.Sprintf("Merge branch 'feature' into main"),
		Parent:  baseCommit,
	}
	if mergeCommit.Parent != baseCommit {
		t.Errorf("expected merge commit parent reference to be set")
	}
}

func TestMainExecution(t *testing.T) {
	main()
}

