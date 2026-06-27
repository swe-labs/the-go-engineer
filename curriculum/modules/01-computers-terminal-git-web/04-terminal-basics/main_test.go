package main

import (
	"testing"
)

func TestSimulateShellEchoHello(t *testing.T) {
	got := simulateShell("echo hello")
	if got != "hello" {
		t.Errorf("simulateShell(\"echo hello\") = %q; want %q", got, "hello")
	}
}

func TestSimulateShellEchoHelloWorld(t *testing.T) {
	got := simulateShell("echo hello world")
	if got != "hello world" {
		t.Errorf("simulateShell(\"echo hello world\") = %q; want %q", got, "hello world")
	}
}

func TestSimulateShellEchoEmpty(t *testing.T) {
	got := simulateShell("echo")
	if got != "" {
		t.Errorf("simulateShell(\"echo\") = %q; want %q", got, "")
	}
}

func TestSimulateShellEchoSpaced(t *testing.T) {
	got := simulateShell("echo   spaced")
	if got != "spaced" {
		t.Errorf("simulateShell(\"echo   spaced\") = %q; want %q", got, "spaced")
	}
}

func TestSimulateShellHello(t *testing.T) {
	got := simulateShell("hello")
	want := "Hello from the Go shell!"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSimulateShellUnknown(t *testing.T) {
	got := simulateShell("foobar")
	want := "command not found: foobar"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSimulateShellEmpty(t *testing.T) {
	got := simulateShell("")
	if got != "" {
		t.Errorf("expected empty output for empty input, got %q", got)
	}
}

func TestSimulateShellMultipleWords(t *testing.T) {
	got := simulateShell("echo   hello    world")
	if got != "hello world" {
		t.Errorf("got %q, want %q", got, "hello world")
	}
}
