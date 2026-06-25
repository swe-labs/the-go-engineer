package main

import "testing"

func TestCleanAndSplit(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:    "empty string",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "whitespace only",
			input:   "   \t\n  ",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "short words filtered",
			input:   "a an the",
			want:    []string{"the"},
			wantErr: false,
		},
		{
			name:    "punctuation stripped",
			input:   "Hello, World! Go is   amazing.",
			want:    []string{"hello", "world", "amazing"},
			wantErr: false,
		},
		{
			name:    "case normalization",
			input:   "HELLO hello Hello",
			want:    []string{"hello", "hello", "hello"},
			wantErr: false,
		},
		{
			name:    "numbers preserved",
			input:   "go1 go2 go3",
			want:    []string{"go1", "go2", "go3"},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CleanAndSplit(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("len=%d; want len=%d; got=%v", len(got), len(tc.want), got)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("words[%d]=%q; want %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}

func TestWordCount(t *testing.T) {
	tests := []struct {
		name   string
		words  []string
		expect map[string]int
	}{
		{
			name:   "nil slice",
			words:  nil,
			expect: map[string]int{},
		},
		{
			name:   "empty slice",
			words:  []string{},
			expect: map[string]int{},
		},
		{
			name:   "single word",
			words:  []string{"hello"},
			expect: map[string]int{"hello": 1},
		},
		{
			name:   "repeated words",
			words:  []string{"a", "b", "a", "c", "b", "a"},
			expect: map[string]int{"a": 3, "b": 2, "c": 1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := WordCount(tc.words)
			if len(got) != len(tc.expect) {
				t.Fatalf("len=%d; want %d; got %v", len(got), len(tc.expect), got)
			}
			for k, v := range tc.expect {
				if got[k] != v {
					t.Errorf("freq[%q]=%d; want %d", k, got[k], v)
				}
			}
		})
	}
}
