package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestRESTAPIOutput(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run failed: %v\nOutput: %s", err, string(output))
	}

	outStr := string(output)
	if !strings.Contains(outStr, "=== HS.10 REST API ===") {
		t.Errorf("Expected output to contain '=== HS.10 REST API ===', got:\n%s", outStr)
	}
}
