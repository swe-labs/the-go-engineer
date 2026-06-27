package main

import "testing"

func TestHashContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantLen int
	}{
		{"empty string", "", 40},
		{"short string", "hello", 40},
		{"longer string", "The quick brown fox jumps over the lazy dog", 40},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HashContent(tt.content)
			if len(got) != tt.wantLen {
				t.Errorf("HashContent() length = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestHashDeterministic(t *testing.T) {
	h1 := HashContent("hello world")
	h2 := HashContent("hello world")
	if h1 != h2 {
		t.Errorf("HashContent not deterministic: %s != %s", h1, h2)
	}
}

func TestNewBlob(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"empty blob", ""},
		{"text blob", "hello world"},
		{"multiline blob", "line1\nline2\nline3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blob := NewBlob(tt.content)
			if blob.Content != tt.content {
				t.Errorf("NewBlob().Content = %q, want %q", blob.Content, tt.content)
			}
			expectedHash := HashContent(tt.content)
			if blob.Hash != expectedHash {
				t.Errorf("NewBlob().Hash = %s, want %s", blob.Hash, expectedHash)
			}
		})
	}
}

func TestNewCommit(t *testing.T) {
	tests := []struct {
		name       string
		message    string
		parentHash string
	}{
		{"root commit with no parent", "Initial commit", ""},
		{"child commit with parent", "Add README", "abc123def4567890123456789012345678901234"},
		{"merge commit", "Merge branch 'main'", "parent1hash00000000000000000000000000000000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commit := NewCommit(tt.message, tt.parentHash)
			if commit.Message != tt.message {
				t.Errorf("NewCommit().Message = %q, want %q", commit.Message, tt.message)
			}
			if commit.Parent != tt.parentHash {
				t.Errorf("NewCommit().Parent = %q, want %q", commit.Parent, tt.parentHash)
			}
			if commit.Hash == "" {
				t.Error("NewCommit().Hash is empty")
			}
		})
	}
}

func TestDifferentContentDifferentHash(t *testing.T) {
	h1 := HashContent("hello")
	h2 := HashContent("world")
	if h1 == h2 {
		t.Error("different content should produce different hashes")
	}
}
