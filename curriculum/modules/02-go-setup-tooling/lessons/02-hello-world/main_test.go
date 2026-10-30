package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestFormatGreeting(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"World", "Hello, World!"},
		{"Gopher", "Hello, Gopher!"},
		{"", "Hello, World!"},
	}

	for _, tt := range tests {
		got := FormatGreeting(tt.input)
		if got != tt.expected {
			t.Errorf("FormatGreeting(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestWriteGreeting(t *testing.T) {
	var buf bytes.Buffer
	err := WriteGreeting(&buf, "Go Engineer")
	if err != nil {
		t.Fatalf("WriteGreeting returned error: %v", err)
	}

	output := strings.TrimSpace(buf.String())
	expected := "Hello, Go Engineer!"
	if output != expected {
		t.Errorf("got %q, want %q", output, expected)
	}
}

func TestMainExecution(t *testing.T) {
	main()
}

