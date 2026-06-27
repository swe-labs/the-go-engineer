package main

import "testing"

func TestStageExistingFile(t *testing.T) {
	working := []string{"a.txt"}
	staging := []string{}
	result, err := Stage("a.txt", working, staging)
	if err != nil {
		t.Errorf("Stage() error = %v", err)
	}
	if len(result) != 1 || result[0] != "a.txt" {
		t.Errorf("Stage() = %v, want [a.txt]", result)
	}
}

func TestStageNonExistentFileReturnsError(t *testing.T) {
	staging := []string{}
	_, err := Stage("b.txt", []string{"a.txt"}, staging)
	if err == nil {
		t.Error("Stage() expected error for non-existent file")
	}
}

func TestStageAlreadyStaged(t *testing.T) {
	working := []string{"c.txt"}
	staging := []string{"c.txt"}
	result, err := Stage("c.txt", working, staging)
	if err != nil {
		t.Errorf("Stage() error = %v", err)
	}
	if len(result) != 1 {
		t.Errorf("Stage() should not duplicate, got %d entries", len(result))
	}
}

func TestCommitWithStagedFiles(t *testing.T) {
	staging := []string{"f1.txt"}
	history := []string{}
	history, err := Commit("feat: add f1", staging, history)
	if err != nil {
		t.Errorf("Commit() error = %v", err)
	}
	if len(history) != 1 {
		t.Errorf("History length = %d, want 1", len(history))
	}
}

func TestCommitNothingStagedReturnsError(t *testing.T) {
	history := []string{}
	_, err := Commit("empty commit", []string{}, history)
	if err == nil {
		t.Error("Commit() expected error for empty staging")
	}
}

func TestUnstage(t *testing.T) {
	staging := []string{"x.txt"}
	result := Unstage("x.txt", staging)
	if len(result) != 0 {
		t.Errorf("expected 0 staged files after unstage, got %d", len(result))
	}
}

func TestMultipleCommits(t *testing.T) {
	staging := []string{"a.txt"}
	history := []string{}
	history, _ = Commit("first: a", staging, history)

	staging = []string{"b.txt"}
	history, _ = Commit("second: b", staging, history)

	if len(history) != 2 {
		t.Errorf("expected 2 commits, got %d", len(history))
	}
	if history[0] != "first: a" {
		t.Errorf("commit 0 = %q, want %q", history[0], "first: a")
	}
	if history[1] != "second: b" {
		t.Errorf("commit 1 = %q, want %q", history[1], "second: b")
	}
}
