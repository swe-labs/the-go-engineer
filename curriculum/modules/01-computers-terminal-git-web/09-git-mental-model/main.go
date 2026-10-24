package main

import (
	"crypto/sha1"
	"fmt"
)

// HashContent returns the SHA-1 hash of content as a hex string (like Git object IDs).
func HashContent(content string) string {
	h := sha1.Sum([]byte(content))
	return fmt.Sprintf("%x", h)
}

// Blob represents a Git-like blob object.
type Blob struct {
	Hash    string
	Content string
}

// NewBlob creates a blob by hashing the content (like git hash-object).
func NewBlob(content string) Blob {
	return Blob{Hash: HashContent(content), Content: content}
}

// Commit represents a Git-like commit with a message, hash, and parent.
type Commit struct {
	Hash    string
	Message string
	Parent  string
}

// NewCommit creates a commit with a hash derived from its content.
func NewCommit(message, parentHash string) Commit {
	content := fmt.Sprintf("commit: %s\nparent: %s", message, parentHash)
	return Commit{Hash: HashContent(content), Message: message, Parent: parentHash}
}

func main() {
	blob := NewBlob("hello world")
	fmt.Printf("Blob hash: %s\n", blob.Hash)

	commit1 := NewCommit("Initial commit", "")
	fmt.Printf("Commit 1 hash: %s\n", commit1.Hash)

	commit2 := NewCommit("Add feature", commit1.Hash)
	fmt.Printf("Commit 2 hash: %s\n", commit2.Hash)

	sameBlob := NewBlob("hello world")
	fmt.Printf("Same content: %s == %s: %v\n", blob.Hash, sameBlob.Hash, blob.Hash == sameBlob.Hash)
}
