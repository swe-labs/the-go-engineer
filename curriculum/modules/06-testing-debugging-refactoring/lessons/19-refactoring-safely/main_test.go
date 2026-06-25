package main

import "testing"

func TestProcess(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"empty", []string{}, ""},
		{"single", []string{"alice"}, "Alice"},
		{"multiple", []string{"alice", "bob", "charlie"}, "Alice,Bob,Charlie"},
		{"already mixed case", []string{"aLiCe"}, "Alice"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := process(tc.input)
			if got != tc.want {
				t.Errorf("process(%v) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
