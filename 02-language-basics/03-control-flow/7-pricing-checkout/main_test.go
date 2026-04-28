package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestPricingCheckoutOutput(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run failed: %v\nOutput: %s", err, string(output))
	}

	outStr := string(output)
	
	expectedLines := []string{
		"Processing checkout:",
		"TSHIRT: 20.00",
		"MUG: 12.50",
		"HAT: 18.00",
		"BOOK promo: 25.99 -> 23.39",
		"skip KEYBOARD: unknown item",
		"subtotal: 73.89",
	}

	for _, line := range expectedLines {
		if !strings.Contains(outStr, line) {
			t.Errorf("Expected output to contain %q, but it didn't.\nFull output:\n%s", line, outStr)
		}
	}
}
