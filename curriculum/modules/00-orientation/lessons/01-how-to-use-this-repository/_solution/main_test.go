package main

import (
	"strings"
	"testing"
)

func TestPurposesCoversRequiredFolders(t *testing.T) {
	got := render(purposes())

	for _, folder := range []string{"metadata/", "curriculum/", "tools/", "docs/", "dist/"} {
		if !strings.Contains(got, folder) {
			t.Fatalf("rendered output missing %s\n%s", folder, got)
		}
	}
}

func TestDistWarnsAgainstHandEditing(t *testing.T) {
	got := render(purposes())

	if !strings.Contains(strings.ToLower(got), "never hand-edit") {
		t.Fatalf("dist purpose should warn against hand editing\n%s", got)
	}
}
