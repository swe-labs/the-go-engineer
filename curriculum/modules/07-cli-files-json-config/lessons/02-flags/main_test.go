package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestFlagDefaults(t *testing.T) {
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("exec failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "Hello, World!") {
		t.Errorf("expected default greeting, got: %s", out)
	}
}

func TestFlagName(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-name=Alice")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("exec failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "Hello, Alice!") {
		t.Errorf("expected 'Hello, Alice!', got: %s", out)
	}
}

func TestFlagCount(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-name=Bob", "-count=3")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("exec failed: %v\n%s", err, out)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	greetingCount := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "Hello,") {
			greetingCount++
		}
	}
	if greetingCount != 3 {
		t.Errorf("expected 3 greeting lines with -count=3, got %d: %s", greetingCount, out)
	}
}

func TestFlagVerbose(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-verbose", "-name=Test")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("exec failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "(verbose)") {
		t.Errorf("expected verbose output, got: %s", out)
	}
}

func TestNonFlagArgs(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-name=Test", "extra1", "extra2")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("exec failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "extra1 extra2") {
		t.Errorf("expected non-flag args in output, got: %s", out)
	}
}

func TestFlagHelp(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-help")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		// -help exits with status 2 by default, that's fine
	}
	if !strings.Contains(string(out), "Usage") {
		t.Errorf("expected usage in help output, got: %s", out)
	}
}
