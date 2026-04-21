// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestContactDirectoryOutput(t *testing.T) {
	cmd := exec.Command("go", "run", ".")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run failed: %v\n%s", err, string(out))
	}

	output := string(out)

	expected := []string{
		"Duplicate add skipped for Alice Wonderland.",
		"Found Bob at index 1 with phone 333-4444",
		"Updated Bob through pointer: 333-9999",
		"Persisted Bob phone: 333-9999",
		"Zack not found.",
	}

	for _, fragment := range expected {
		if !strings.Contains(output, fragment) {
			t.Fatalf("expected output to contain %q\nfull output:\n%s", fragment, output)
		}
	}
}
