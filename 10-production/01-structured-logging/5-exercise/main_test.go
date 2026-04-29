package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestStructuredLoggingRedaction(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run failed: %v", err)
	}

	outStr := string(output)

	// Ensure PII is redacted
	if !strings.Contains(outStr, `"[REDACTED]"`) {
		t.Errorf("Expected output to contain redacted fields")
	}

	// Ensure actual sensitive data is NOT present
	if strings.Contains(outStr, "supersecret123") {
		t.Errorf("Password was leaked in logs")
	}
	if strings.Contains(outStr, "4111-1111-1111-1111") {
		t.Errorf("Credit card was leaked in logs")
	}

	// Ensure non-sensitive data is present
	if !strings.Contains(outStr, "jdoe") {
		t.Errorf("Username was incorrectly redacted")
	}
}
