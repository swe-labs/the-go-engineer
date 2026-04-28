package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestTimeoutClientOutput(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run failed: %v\nOutput: %s", err, string(output))
	}

	outStr := string(output)
	if !strings.Contains(outStr, "=== Timeout-Aware API Client ===") {
		t.Errorf("Expected output to contain '=== Timeout-Aware API Client ===', got:\n%s", outStr)
	}
}
