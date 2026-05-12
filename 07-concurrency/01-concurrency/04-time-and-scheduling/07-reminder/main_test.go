package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestReminderOutput(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run failed: %v\nOutput: %s", err, string(output))
	}

	outStr := string(output)
	if !strings.Contains(outStr, "=== Console Reminder ===") {
		t.Errorf("Expected output to contain '=== Console Reminder ===', got:\n%s", outStr)
	}
}
