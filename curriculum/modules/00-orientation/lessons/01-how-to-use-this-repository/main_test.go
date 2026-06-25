package main

import (
	"testing"
)

func TestLookupFound(t *testing.T) {
	tests := []struct {
		name     string
		wantName string
		wantOK   bool
	}{
		{name: "metadata/", wantName: "metadata/", wantOK: true},
		{name: "curriculum/", wantName: "curriculum/", wantOK: true},
		{name: "tools/", wantName: "tools/", wantOK: true},
		{name: "docs/", wantName: "docs/", wantOK: true},
		{name: "dist/", wantName: "dist/", wantOK: true},
		{name: "metadata", wantName: "metadata/", wantOK: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Lookup(tc.name)
			if got == nil {
				t.Fatalf("Lookup(%q) = nil, want entry", tc.name)
			}
			if got.Name != tc.wantName {
				t.Errorf("Lookup(%q).Name = %q, want %q", tc.name, got.Name, tc.wantName)
			}
			if got.Purpose == "" {
				t.Error("Purpose must not be empty")
			}
			if len(got.Contents) == 0 {
				t.Error("Contents must not be empty")
			}
		})
	}
}

func TestLookupNotFound(t *testing.T) {
	got := Lookup("nonexistent/")
	if got != nil {
		t.Fatalf("Lookup(%q) = %v, want nil", "nonexistent/", got)
	}
}

func TestRepoDirsNotEmpty(t *testing.T) {
	if len(repoDirs) == 0 {
		t.Fatal("repoDirs must have at least one entry")
	}
}
