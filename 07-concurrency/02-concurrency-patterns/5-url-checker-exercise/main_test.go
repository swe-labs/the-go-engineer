package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestExerciseOutput(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Process exited with error (this might be expected): %v", err)
	}

	outStr := string(output)
	if !strings.Contains(outStr, "=== URL Health Checker ===") {
		t.Errorf("Expected output to contain '%s', got:\n%s", "=== URL Health Checker ===", outStr)
	}
}
