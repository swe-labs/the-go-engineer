package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestGRPCServiceOutput(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run failed: %v\nOutput: %s", err, string(output))
	}

	outStr := string(output)
	if !strings.Contains(outStr, "=== API.9 gRPC Service ===") {
		t.Errorf("Expected output to contain '=== API.9 gRPC Service ===', got:\n%s", outStr)
	}
}
