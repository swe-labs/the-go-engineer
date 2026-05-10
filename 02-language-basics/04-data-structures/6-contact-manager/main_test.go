// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestContactDirectory(t *testing.T) {
	// Execute the program
	cmd := exec.Command("go", "run", ".")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run failed: %v\n%s", err, string(out))
	}

	output := string(out)

	// Table-driven verification of output fragments
	tests := []struct {
		name     string
		fragment string
	}{
		{"duplicate check", "Duplicate add skipped for Alice Wonderland."},
		{"listing alice", "1. Alice Wonderland | alice@example.com | 111-2222"},
		{"listing bob", "2. Bob The Builder | bob@example.com | 333-4444"},
		{"lookup bob", "Found Bob at index 1 with phone 333-4444"},
		{"update pointer", "Updated Bob through pointer: 333-9999"},
		{"persistence check", "Persisted Bob phone: 333-9999"},
		{"negative lookup", "Zack not found."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(output, tt.fragment) {
				t.Errorf("expected output to contain %q", tt.fragment)
			}
		})
	}
}
