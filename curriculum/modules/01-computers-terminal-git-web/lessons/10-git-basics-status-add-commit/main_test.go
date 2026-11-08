package main

import (
	"reflect"
	"testing"
)

func TestRepo_StageAndCommit(t *testing.T) {
	repo := NewRepo()
	if repo == nil {
		t.Fatalf("expected repo instance, got nil")
	}

	// 1. Stage non-existent file -> should fail
	err := repo.Stage("non-existent.txt")
	if err == nil {
		t.Errorf("expected error staging non-existent file")
	}

	// 2. Commit with empty stage -> should fail
	err = repo.Commit("empty commit")
	if err == nil {
		t.Errorf("expected error committing with empty staging area")
	}

	// 3. Add files to working directory
	repo.Working["README.md"] = "# Project"
	repo.Working["main.go"] = "package main"

	// 4. Stage files
	if err := repo.Stage("README.md"); err != nil {
		t.Fatalf("unexpected error staging README.md: %v", err)
	}
	if err := repo.Stage("main.go"); err != nil {
		t.Fatalf("unexpected error staging main.go: %v", err)
	}

	expectedStaged := []string{"README.md", "main.go"}
	actualStaged := repo.StageNames()
	if !reflect.DeepEqual(actualStaged, expectedStaged) {
		t.Errorf("expected staged %v, got %v", expectedStaged, actualStaged)
	}

	// 5. Unstage a file
	repo.Unstage("main.go")
	if len(repo.StageNames()) != 1 || repo.StageNames()[0] != "README.md" {
		t.Errorf("expected only README.md staged after unstage")
	}

	// 6. Commit staged changes
	err = repo.Commit("feat: initial readme")
	if err != nil {
		t.Fatalf("commit failed: %v", err)
	}

	if len(repo.History) != 1 {
		t.Fatalf("expected 1 commit in history, got %d", len(repo.History))
	}
	if repo.History[0].Message != "feat: initial readme" {
		t.Errorf("commit message mismatch: %q", repo.History[0].Message)
	}
	if repo.History[0].Files["README.md"] != "# Project" {
		t.Errorf("committed file content mismatch")
	}

	// Staging area should be empty after commit
	if len(repo.StageNames()) != 0 {
		t.Errorf("staging area should be cleared after commit")
	}
}

func TestMainExecution(t *testing.T) {
	main()
}
