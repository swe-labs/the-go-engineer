package main

import (
	"testing"
)

func TestSimulateShellEcho(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"echo hello", "hello"},
		{"echo hello world", "hello world"},
		{"echo", ""},
		{"echo   spaced", "spaced"},
	}

	for _, tt := range tests {
		result := simulateShell(tt.input)
		if result.Output != tt.want {
			t.Errorf("simulateShell(%q).Output = %q; want %q", tt.input, result.Output, tt.want)
		}
		if result.Input != tt.input {
			t.Errorf("simulateShell(%q).Input = %q; want %q", tt.input, result.Input, tt.input)
		}
	}
}

func TestSimulateShellHello(t *testing.T) {
	result := simulateShell("hello")
	if result.Output != "Hello from the Go shell!" {
		t.Errorf("got %q, want %q", result.Output, "Hello from the Go shell!")
	}
}

func TestSimulateShellUnknown(t *testing.T) {
	result := simulateShell("foobar")
	want := "command not found: foobar"
	if result.Output != want {
		t.Errorf("got %q, want %q", result.Output, want)
	}
}

func TestSimulateShellEmpty(t *testing.T) {
	result := simulateShell("")
	if result.Output != "" {
		t.Errorf("expected empty output for empty input, got %q", result.Output)
	}
}

func TestSimulateShellMultipleWords(t *testing.T) {
	result := simulateShell("echo   hello    world")
	if result.Output != "hello world" {
		t.Errorf("got %q, want %q", result.Output, "hello world")
	}
}
