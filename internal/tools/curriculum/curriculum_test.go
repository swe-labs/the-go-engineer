package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIndexOf(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		target   string
		expected int
	}{
		{
			name:     "found at start",
			slice:    []string{"id", "title", "order"},
			target:   "id",
			expected: 0,
		},
		{
			name:     "found in middle",
			slice:    []string{"id", "title", "order"},
			target:   "title",
			expected: 1,
		},
		{
			name:     "found at end",
			slice:    []string{"id", "title", "order"},
			target:   "order",
			expected: 2,
		},
		{
			name:     "not found",
			slice:    []string{"id", "title", "order"},
			target:   "nonexistent",
			expected: -1,
		},
		{
			name:     "empty slice",
			slice:    []string{},
			target:   "any",
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := indexOf(tt.slice, tt.target)
			if result != tt.expected {
				t.Fatalf("expected indexOf(%v, %q) = %d, got %d", tt.slice, tt.target, tt.expected, result)
			}
		})
	}
}

func TestExtractItemKeysInOrder(t *testing.T) {
	rawJSON := []byte(`{
		"schema_version": "3.1.0",
		"items": [
			{
				"id": "item-01",
				"title": "First Item",
				"order": 1
			},
			{
				"id": "item-02",
				"order": 2,
				"title": "Second Item"
			}
		]
	}`)

	entries := extractItemKeysInOrder(rawJSON)
	if len(entries) != 2 {
		t.Fatalf("expected 2 itemKeyEntries, got %d", len(entries))
	}

	// Verify first item
	if entries[0].id != "item-01" {
		t.Errorf("entry 0: expected id 'item-01', got %q", entries[0].id)
	}
	expectedKeys0 := []string{"id", "title", "order"}
	if len(entries[0].keys) != len(expectedKeys0) {
		t.Fatalf("entry 0: expected %d keys, got %d", len(expectedKeys0), len(entries[0].keys))
	}
	for i, k := range expectedKeys0 {
		if entries[0].keys[i] != k {
			t.Errorf("entry 0 key[%d]: expected %q, got %q", i, k, entries[0].keys[i])
		}
	}

	// Verify second item (preserving explicit byte ordering)
	if entries[1].id != "item-02" {
		t.Errorf("entry 1: expected id 'item-02', got %q", entries[1].id)
	}
	expectedKeys1 := []string{"id", "order", "title"}
	if len(entries[1].keys) != len(expectedKeys1) {
		t.Fatalf("entry 1: expected %d keys, got %d", len(expectedKeys1), len(entries[1].keys))
	}
	for i, k := range expectedKeys1 {
		if entries[1].keys[i] != k {
			t.Errorf("entry 1 key[%d]: expected %q, got %q", i, k, entries[1].keys[i])
		}
	}
}

func TestResolvePath(t *testing.T) {
	// 1. Absolute path returns directly
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd failed: %v", err)
	}
	abs := filepath.Join(wd, "any-test-file.txt")
	if resolved := resolvePath(abs); resolved != abs {
		t.Errorf("expected absolute path %q to resolve to itself, got %q", abs, resolved)
	}

	// 2. Relative path to existing file
	relFile := "path.core.json"
	resolved := resolvePath(relFile)
	if _, err := os.Stat(resolved); err != nil {
		t.Errorf("expected resolvePath(%q) to point to existing file, got %q: %v", relFile, resolved, err)
	}
}

func TestCanonicalValidationGatesExecution(t *testing.T) {
	// Reset errorCount before testing gates
	errorCount = 0

	// Run individual gate functions
	validateSchema()
	if errorCount != 0 {
		t.Errorf("validateSchema failed with %d errors", errorCount)
	}

	validateGraph()
	if errorCount != 0 {
		t.Errorf("validateGraph failed with %d errors", errorCount)
	}

	validateNoPlaceholderZM()
	if errorCount != 0 {
		t.Errorf("validateNoPlaceholderZM failed with %d errors", errorCount)
	}

	validateAssessmentCompleteness()
	if errorCount != 0 {
		t.Errorf("validateAssessmentCompleteness failed with %d errors", errorCount)
	}

	validateUniqueTags()
	if errorCount != 0 {
		t.Errorf("validateUniqueTags failed with %d errors", errorCount)
	}

	validateCrossrefs()
	if errorCount != 0 {
		t.Errorf("validateCrossrefs failed with %d errors", errorCount)
	}

	validateConceptOwnership()
	if errorCount != 0 {
		t.Errorf("validateConceptOwnership failed with %d errors", errorCount)
	}

	validateCognitiveLoad()
	if errorCount != 0 {
		t.Errorf("validateCognitiveLoad failed with %d errors", errorCount)
	}

	validateOperationalFailures()
	if errorCount != 0 {
		t.Errorf("validateOperationalFailures failed with %d errors", errorCount)
	}

	validatePhysicalArtifacts()
	if errorCount != 0 {
		t.Errorf("validatePhysicalArtifacts failed with %d errors", errorCount)
	}

	validateAssessmentArtifacts()
	if errorCount != 0 {
		t.Errorf("validateAssessmentArtifacts failed with %d errors", errorCount)
	}

	validateStarterParity()
	if errorCount != 0 {
		t.Errorf("validateStarterParity failed with %d errors", errorCount)
	}
}
