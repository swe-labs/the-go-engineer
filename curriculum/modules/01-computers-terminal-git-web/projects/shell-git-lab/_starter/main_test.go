package main

import (
	"testing"
)

func TestNewRepo(t *testing.T) {
	r := NewRepo()
	if r == nil {
		t.Fatal("NewRepo returned nil")
	}
}

func TestGitInit(t *testing.T) {
	r := NewRepo()
	r.GitInit()
	if r.Current != "main" {
		t.Errorf("expected current branch 'main', got %q", r.Current)
	}
	if _, ok := r.Branches["main"]; !ok {
		t.Error("expected 'main' branch to exist after init")
	}
}

func TestCreateFile(t *testing.T) {
	r := NewRepo()
	r.CreateFile("test.txt", "hello")
	if !r.FileExists("test.txt") {
		t.Error("expected file to exist after CreateFile")
	}
}

func TestGitAddNonExistent(t *testing.T) {
	r := NewRepo()
	err := r.GitAdd("nonexistent.txt")
	if err == nil {
		t.Error("expected error when adding non-existent file")
	}
}

func TestGitAddAndCommit(t *testing.T) {
	r := NewRepo()
	r.GitInit()
	r.CreateFile("main.go", "package main")
	err := r.GitAdd("main.go")
	if err != nil {
		t.Fatalf("GitAdd failed: %v", err)
	}
	err = r.GitCommit("Add main.go")
	if err != nil {
		t.Fatalf("GitCommit failed: %v", err)
	}
	if len(r.Commits) != 1 {
		t.Errorf("expected 1 commit, got %d", len(r.Commits))
	}
	if len(r.Staged) != 0 {
		t.Errorf("expected 0 staged files after commit, got %d", len(r.Staged))
	}
}

func TestGitCommitNothingStaged(t *testing.T) {
	r := NewRepo()
	r.GitInit()
	err := r.GitCommit("empty commit")
	if err == nil {
		t.Error("expected error when committing with nothing staged")
	}
}

func TestGitBranch(t *testing.T) {
	r := NewRepo()
	r.GitInit()
	r.CreateFile("a.go", "a")
	r.GitAdd("a.go")
	r.GitCommit("first")
	err := r.GitBranch("feature")
	if err != nil {
		t.Fatalf("GitBranch failed: %v", err)
	}
	if _, ok := r.Branches["feature"]; !ok {
		t.Error("expected 'feature' branch to exist")
	}
}

func TestGitBranchDuplicate(t *testing.T) {
	r := NewRepo()
	r.GitInit()
	err := r.GitBranch("main")
	if err == nil {
		t.Error("expected error when creating duplicate branch")
	}
}

func TestGitCheckout(t *testing.T) {
	r := NewRepo()
	r.GitInit()
	r.GitBranch("dev")
	err := r.GitCheckout("dev")
	if err != nil {
		t.Fatalf("GitCheckout failed: %v", err)
	}
	if r.Current != "dev" {
		t.Errorf("expected current branch 'dev', got %q", r.Current)
	}
}

func TestGitCheckoutNonExistent(t *testing.T) {
	r := NewRepo()
	r.GitInit()
	err := r.GitCheckout("nonexistent")
	if err == nil {
		t.Error("expected error when checking out non-existent branch")
	}
}

func TestGitMerge(t *testing.T) {
	r := NewRepo()
	r.GitInit()
	r.CreateFile("a.go", "a")
	r.GitAdd("a.go")
	r.GitCommit("first")

	r.GitBranch("feature")
	r.GitCheckout("feature")
	r.CreateFile("b.go", "b")
	r.GitAdd("b.go")
	r.GitCommit("feature work")

	r.GitCheckout("main")
	err := r.GitMerge("feature")
	if err != nil {
		t.Fatalf("GitMerge failed: %v", err)
	}
	if len(r.Commits) < 2 {
		t.Errorf("expected at least 2 commits after merge, got %d", len(r.Commits))
	}
}

func TestGitMergeNonExistent(t *testing.T) {
	r := NewRepo()
	r.GitInit()
	err := r.GitMerge("nonexistent")
	if err == nil {
		t.Error("expected error when merging non-existent branch")
	}
}

func TestListFiles(t *testing.T) {
	r := NewRepo()
	r.CreateFile("a.txt", "a")
	r.CreateFile("b.txt", "b")
	files := r.ListFiles()
	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d", len(files))
	}
}

func TestReadFile(t *testing.T) {
	r := NewRepo()
	r.CreateFile("hello.txt", "world")
	content, err := r.ReadFile("hello.txt")
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if content != "world" {
		t.Errorf("expected 'world', got %q", content)
	}
}

func TestReadFileNonExistent(t *testing.T) {
	r := NewRepo()
	_, err := r.ReadFile("missing.txt")
	if err == nil {
		t.Error("expected error when reading non-existent file")
	}
}
