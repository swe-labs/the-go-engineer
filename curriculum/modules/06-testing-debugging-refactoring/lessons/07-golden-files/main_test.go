package main

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

func TestFormatPerson(t *testing.T) {
	tests := []struct {
		name   string
		person string
		age    int
	}{
		{name: "alice", person: "Alice", age: 30},
		{name: "bob", person: "Bob", age: 25},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := FormatPerson(tc.person, tc.age)
			golden := filepath.Join("testdata", tc.name+".golden")

			if *update {
				os.MkdirAll("testdata", 0755)
				if err := os.WriteFile(golden, []byte(got), 0644); err != nil {
					t.Fatalf("writing golden file: %v", err)
				}
				return
			}

			wantBytes, err := os.ReadFile(golden)
			if err != nil {
				t.Fatalf("reading golden file: %v", err)
			}
			want := string(wantBytes)
			if got != want {
				t.Errorf("output mismatch\n got: %q\nwant: %q", got, want)
			}
		})
	}
}

func TestGenerateReport(t *testing.T) {
	tests := []struct {
		name  string
		items []string
	}{
		{name: "empty", items: []string{}},
		{name: "single", items: []string{"apples"}},
		{name: "multiple", items: []string{"apples", "bananas", "cherries"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := GenerateReport(tc.items)
			golden := filepath.Join("testdata", "report-"+tc.name+".golden")

			if *update {
				os.MkdirAll("testdata", 0755)
				if err := os.WriteFile(golden, []byte(got), 0644); err != nil {
					t.Fatalf("writing golden file: %v", err)
				}
				return
			}

			wantBytes, err := os.ReadFile(golden)
			if err != nil {
				t.Fatalf("reading golden file: %v", err)
			}
			want := string(wantBytes)
			if got != want {
				t.Errorf("output mismatch for %q\n got: %q\nwant: %q", tc.name, got, want)
			}
		})
	}
}
