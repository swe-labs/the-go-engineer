package main

import "testing"

func TestHashContent(t *testing.T) {
	got := HashContent("hello")
	if got == "" {
		t.Error("HashContent() returned empty string")
	}
}

func TestHashContentEmpty(t *testing.T) {
	got := HashContent("")
	if got == "" {
		t.Error("HashContent('') returned empty string")
	}
}

func TestHashDeterministic(t *testing.T) {
	h1 := HashContent("hello world")
	h2 := HashContent("hello world")
	if h1 != h2 {
		t.Errorf("HashContent not deterministic: %s != %s", h1, h2)
	}
}

func TestDifferentContentDifferentHash(t *testing.T) {
	h1 := HashContent("hello")
	h2 := HashContent("world")
	if h1 == h2 {
		t.Error("different content should produce different hashes")
	}
}
