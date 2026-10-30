package main

import (
	"testing"
)

func TestHashContent(t *testing.T) {
	// SHA-1 of empty string: da39a3ee5e6b4b0d3255bfef95601890afd80709
	emptyHash := HashContent("")
	expectedEmpty := "da39a3ee5e6b4b0d3255bfef95601890afd80709"
	if emptyHash != expectedEmpty {
		t.Errorf("expected empty hash %q, got %q", expectedEmpty, emptyHash)
	}

	h1 := HashContent("hello world")
	h2 := HashContent("hello world")
	h3 := HashContent("hello world!")

	if h1 != h2 {
		t.Errorf("expected deterministic hashing: %q != %q", h1, h2)
	}
	if h1 == h3 {
		t.Errorf("different content should produce different hashes")
	}
}

func TestNewBlob(t *testing.T) {
	content := "package main\n\nfunc main() {}\n"
	blob := NewBlob(content)

	if blob.Content != content {
		t.Errorf("blob content mismatch")
	}
	if blob.Hash != HashContent(content) {
		t.Errorf("blob hash %q does not match computed hash %q", blob.Hash, HashContent(content))
	}

	// Content deduplication check
	blob2 := NewBlob(content)
	if blob.Hash != blob2.Hash {
		t.Errorf("identical content should produce identical object hashes")
	}
}

func TestNewCommit(t *testing.T) {
	// 1. Initial root commit (no parent)
	rootCommit := NewCommit("Initial commit", "")
	if rootCommit.Parent != "" {
		t.Errorf("expected root commit to have no parent")
	}
	if rootCommit.Message != "Initial commit" {
		t.Errorf("expected commit message 'Initial commit'")
	}
	if rootCommit.Hash == "" {
		t.Errorf("commit hash should not be empty")
	}

	// 2. Child commit referencing root
	childCommit := NewCommit("Add feature", rootCommit.Hash)
	if childCommit.Parent != rootCommit.Hash {
		t.Errorf("expected child parent %q, got %q", rootCommit.Hash, childCommit.Parent)
	}
	if childCommit.Hash == rootCommit.Hash {
		t.Errorf("child commit should have a distinct hash from root commit")
	}
}

func TestMainExecution(t *testing.T) {
	main()
}

