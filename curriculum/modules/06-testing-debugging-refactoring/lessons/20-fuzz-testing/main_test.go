package main

import (
	"testing"
)

func TestReverse(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "olleh"},
		{"world", "dlrow"},
		{"", ""},
		{"a", "a"},
		{"ab", "ba"},
		{"héllo", "olléh"},
		{"12345", "54321"},
	}
	for _, tc := range tests {
		got, err := Reverse(tc.input)
		if err != nil {
			t.Errorf("Reverse(%q) returned error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("Reverse(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestReverseInvalidUTF8(t *testing.T) {
	_, err := Reverse("\xff\xfe")
	if err == nil {
		t.Error("expected error for invalid UTF-8")
	}
}

func FuzzReverse(f *testing.F) {
	seeds := []string{"hello", "world", "12345", "", "a", "abc", "héllo"}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		first, err := Reverse(s)
		if err != nil {
			return
		}
		second, err := Reverse(first)
		if err != nil {
			t.Errorf("Reverse(Reverse(%q)) error: %v", s, err)
		}
		if second != s {
			t.Errorf("Reverse(Reverse(%q)) = %q, want %q", s, second, s)
		}
	})
}
