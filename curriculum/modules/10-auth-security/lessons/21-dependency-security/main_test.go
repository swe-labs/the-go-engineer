package main

import (
	"os"
	"testing"
)

func TestHashModule(t *testing.T) {
	tmpFile := "test_hash_mod"
	content := "hello world"
	os.WriteFile(tmpFile, []byte(content), 0644)
	defer os.Remove(tmpFile)

	h, err := hashModule(tmpFile)
	if err != nil {
		t.Fatalf("hashModule error: %v", err)
	}

	if len(h) != 64 {
		t.Errorf("expected 64 hex chars, got %d", len(h))
	}
}

func TestParseGoSum(t *testing.T) {
	sumContent := `github.com/gorilla/mux v1.8.1 h1:TuMF1mMtQ==
github.com/lib/pq v1.10.9 h1:abc123==
golang.org/x/crypto v0.17.0 h1:xyz789==
`
	tmpFile := "test_go_sum"
	os.WriteFile(tmpFile, []byte(sumContent), 0644)
	defer os.Remove(tmpFile)

	entries, err := parseGoSum(tmpFile)
	if err != nil {
		t.Fatalf("parseGoSum error: %v", err)
	}

	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}

	tests := []struct {
		index int
		want  string
	}{
		{index: 0, want: "github.com/gorilla/mux"},
		{index: 1, want: "github.com/lib/pq"},
		{index: 2, want: "golang.org/x/crypto"},
	}

	for _, tc := range tests {
		entry := entries[tc.index]
		if entry.Path != tc.want {
			t.Errorf("entries[%d].Path = %q, want %q", tc.index, entry.Path, tc.want)
		}
		if entry.Version == "" {
			t.Errorf("entries[%d].Version is empty", tc.index)
		}
		if entry.Hash == "" {
			t.Errorf("entries[%d].Hash is empty", tc.index)
		}
	}
}

func TestParseGoSumEmpty(t *testing.T) {
	tmpFile := "test_go_sum_empty"
	os.WriteFile(tmpFile, []byte(""), 0644)
	defer os.Remove(tmpFile)

	entries, err := parseGoSum(tmpFile)
	if err != nil {
		t.Fatalf("parseGoSum error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestMinifyDeps(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  int
	}{
		{
			name:  "removes unused- prefix deps",
			input: []string{"github.com/gorilla/mux", "unused-old-pkg", "github.com/lib/pq", "unused-debug"},
			want:  2,
		},
		{
			name:  "all deps kept",
			input: []string{"github.com/gorilla/mux", "github.com/lib/pq"},
			want:  2,
		},
		{
			name:  "empty input",
			input: []string{},
			want:  0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := minifyDeps(tc.input)
			if len(got) != tc.want {
				t.Errorf("minifyDeps returned %d deps, want %d", len(got), tc.want)
			}
		})
	}
}
