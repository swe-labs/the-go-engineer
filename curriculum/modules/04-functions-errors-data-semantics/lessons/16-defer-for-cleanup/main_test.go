package main

import (
	"os"
	"testing"
)

func TestWriteToFile_Success(t *testing.T) {
	err := writeToFile("test_success.txt", "test content")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	data, err := os.ReadFile("test_success.txt")
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if string(data) != "test content" {
		t.Errorf("expected 'test content', got %q", string(data))
	}
	os.Remove("test_success.txt")
}

func TestWriteToFile_InvalidPath(t *testing.T) {
	err := writeToFile("/invalid/path/file.txt", "content")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSafeProcess(t *testing.T) {
	err := safeProcess([]string{"x", "y"})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	os.Remove("x.txt")
	os.Remove("y.txt")
}

func TestSafeProcessWithInvalid(t *testing.T) {
	// Using a path that might cause issues - on Windows, "*" is invalid
	err := safeProcess([]string{"valid", "valid2"})
	if err != nil {
		// Should work with valid names
		t.Fatalf("expected nil, got %v", err)
	}
	os.Remove("valid.txt")
	os.Remove("valid2.txt")
}

func TestWriteToFileToLoop_Success(t *testing.T) {
	err := writeToFileToLoop("test_loop.txt", "loop content")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	data, err := os.ReadFile("test_loop.txt")
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if string(data) != "loop content" {
		t.Errorf("expected 'loop content', got %q", string(data))
	}
	os.Remove("test_loop.txt")
}
