package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestGreet(t *testing.T) {
	var buf bytes.Buffer
	Greet(&buf, "Alice")
	want := "Hello, Alice!\n"
	if buf.String() != want {
		t.Errorf("Greet = %q; want %q", buf.String(), want)
	}
}

func TestGreetEmpty(t *testing.T) {
	var buf bytes.Buffer
	Greet(&buf, "")
	want := "Hello, !\n"
	if buf.String() != want {
		t.Errorf("Greet(\"\") = %q; want %q", buf.String(), want)
	}
}

func TestNewLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, "TEST: ")
	logger.Println("hello world")
	if !strings.Contains(buf.String(), "TEST: hello world") {
		t.Errorf("log = %q; want containing %q", buf.String(), "TEST: hello world")
	}
}

func TestWriteTable(t *testing.T) {
	var buf bytes.Buffer
	WriteTable(&buf, [][2]string{
		{"Alice", "30"},
		{"Bob", "25"},
	})
	out := buf.String()
	if !strings.Contains(out, "Name    | Age") {
		t.Error("missing header")
	}
	if !strings.Contains(out, "Alice") {
		t.Error("missing Alice")
	}
	if !strings.Contains(out, "Bob") {
		t.Error("missing Bob")
	}
}

func TestWriteTableEmpty(t *testing.T) {
	var buf bytes.Buffer
	WriteTable(&buf, [][2]string{})
	out := buf.String()
	if !strings.Contains(out, "Name    | Age") {
		t.Error("missing header for empty table")
	}
	if !strings.Contains(out, "--------|-----") {
		t.Error("missing separator for empty table")
	}
}

func TestWriteTableAlignment(t *testing.T) {
	var buf bytes.Buffer
	WriteTable(&buf, [][2]string{
		{"Eve", "28"},
	})
	out := buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 3 {
		t.Fatalf("expected at least 3 lines, got %d", len(lines))
	}
	lastLine := lines[2]
	if !strings.HasPrefix(lastLine, "Eve") {
		t.Errorf("expected line to start with %q, got %q", "Eve", lastLine)
	}
}
