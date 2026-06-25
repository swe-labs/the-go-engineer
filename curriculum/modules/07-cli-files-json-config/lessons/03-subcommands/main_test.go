package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestSubcommandGreet(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "greet", "-name=Alice")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("exec failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "Hello, Alice!") {
		t.Errorf("expected greeting, got: %s", out)
	}
}

func TestSubcommandGreetDefault(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "greet")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("exec failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "Hello, World!") {
		t.Errorf("expected default greeting, got: %s", out)
	}
}

func TestSubcommandAdd(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "add", "-a=10", "-b=3")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("exec failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "10 + 3 = 13") {
		t.Errorf("expected addition result, got: %s", out)
	}
}

func TestSubcommandUnknown(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "unknown")
	cmd.Dir = "."
	out, _ := cmd.CombinedOutput()
	if !strings.Contains(string(out), "Unknown command") {
		t.Errorf("expected unknown command error, got: %s", out)
	}
}

func TestSubcommandNoArgs(t *testing.T) {
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = "."
	out, _ := cmd.CombinedOutput()
	if !strings.Contains(string(out), "Usage") {
		t.Errorf("expected usage output, got: %s", out)
	}
}

func TestSubcommandFlagsOnlyAfterCommand(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "greet", "-name=Bob")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("exec failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "Hello, Bob!") {
		t.Errorf("expected 'Hello, Bob!', got: %s", out)
	}
}
