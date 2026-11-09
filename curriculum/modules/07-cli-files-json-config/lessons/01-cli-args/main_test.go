package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCliArgsWithName(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "Alice")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("exec failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "Hello Alice") {
		t.Errorf("expected output to contain 'Hello Alice', got: %s", out)
	}
}

func TestCliArgsWithNameAndAge(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "Bob", "30")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("exec failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "Hello Bob, age 30") {
		t.Errorf("expected output to contain 'Hello Bob, age 30', got: %s", out)
	}
}

func TestCliArgsInvalidAge(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "Bob", "abc")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("exec failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "Invalid age") {
		t.Errorf("expected error for invalid age, got: %s", out)
	}
}

func TestCliArgsNoArgsShowsUsage(t *testing.T) {
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("exec failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "Usage:") {
		t.Errorf("expected usage message, got: %s", out)
	}
}

func TestOsArgsLength(t *testing.T) {
	// os.Args[0] is always the program name
	if len(os.Args) == 0 {
		t.Error("os.Args must have at least one element (the program name)")
	}
}

func TestOsArgsFirstIsProgram(t *testing.T) {
	if !strings.HasSuffix(os.Args[0], ".test") && !strings.Contains(os.Args[0], "go-build") {
		t.Logf("os.Args[0] = %q (expected program path during 'go test')", os.Args[0])
	}
}
