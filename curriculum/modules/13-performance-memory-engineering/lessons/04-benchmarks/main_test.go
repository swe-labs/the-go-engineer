package main

import (
	"strings"
	"testing"
)

func BenchmarkConcatPlus(b *testing.B) {
	parts := []string{"hello", " ", "world", " ", "from", " ", "Go"}
	for i := 0; i < b.N; i++ {
		concatPlus(parts)
	}
}

func BenchmarkConcatBuilder(b *testing.B) {
	parts := []string{"hello", " ", "world", " ", "from", " ", "Go"}
	for i := 0; i < b.N; i++ {
		concatBuilder(parts)
	}
}

func BenchmarkConcatBuilderWithSize(b *testing.B) {
	parts := []string{"hello", " ", "world", " ", "from", " ", "Go"}
	totalSize := 0
	for _, p := range parts {
		totalSize += len(p)
	}
	for i := 0; i < b.N; i++ {
		var buf strings.Builder
		buf.Grow(totalSize)
		for _, p := range parts {
			buf.WriteString(p)
		}
		_ = buf.String()
	}
}

func BenchmarkStringBuilderVsPlus(b *testing.B) {
	parts := []string{"hello", " ", "world", " ", "from", " ", "Go"}
	b.Run("concatPlus", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			concatPlus(parts)
		}
	})
	b.Run("concatBuilder", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			concatBuilder(parts)
		}
	})
}

func TestConcatCorrectness(t *testing.T) {
	parts := []string{"a", "b", "c"}
	expected := "abc"

	got := concatPlus(parts)
	if got != expected {
		t.Errorf("concatPlus = %q; want %q", got, expected)
	}

	got = concatBuilder(parts)
	if got != expected {
		t.Errorf("concatBuilder = %q; want %q", got, expected)
	}
}

func TestConcatTableDriven(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{"empty", []string{}, ""},
		{"single", []string{"x"}, "x"},
		{"multiple", []string{"hello", " ", "world"}, "hello world"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := concatPlus(tc.input); got != tc.expected {
				t.Errorf("concatPlus = %q; want %q", got, tc.expected)
			}
			if got := concatBuilder(tc.input); got != tc.expected {
				t.Errorf("concatBuilder = %q; want %q", got, tc.expected)
			}
		})
	}
}
