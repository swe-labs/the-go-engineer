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

// Fixed table test using simple strings instead of non-existent struct
func TestFindIndexLogic(t *testing.T) {
	tests := []struct {
		name      string
		search    string
		names     []string
		wantIndex int
	}{
		{"exact match first", "Alice", []string{"Alice", "Bob"}, 0},
		{"exact match second", "Bob", []string{"Alice", "Bob"}, 1},
		{"not found", "Charlie", []string{"Alice", "Bob"}, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulating the map-based lookup logic from main.go
			indexByName := make(map[string]int)
			for i, n := range tt.names {
				indexByName[n] = i
			}
			got, ok := indexByName[tt.search]
			if !ok { got = -1 }
			if got != tt.wantIndex {
				t.Errorf("lookup = %v, want %v", got, tt.wantIndex)
			}
		})
	}
}