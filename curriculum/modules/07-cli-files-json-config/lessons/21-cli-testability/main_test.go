package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunLineCount(t *testing.T) {
	input := "hello world\nfoo bar baz\nline three\n"
	var buf bytes.Buffer
	code := Run(strings.NewReader(input), &buf, []string{"-l"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	got := strings.TrimSpace(buf.String())
	if got != "3" {
		t.Errorf("expected 3 lines, got %q", got)
	}
}

func TestRunWordCount(t *testing.T) {
	input := "hello world\nfoo bar baz\n"
	var buf bytes.Buffer
	code := Run(strings.NewReader(input), &buf, []string{"-w"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	got := strings.TrimSpace(buf.String())
	if got != "5" {
		t.Errorf("expected 5 words, got %q", got)
	}
}

func TestRunCharCount(t *testing.T) {
	input := "hello world\nfoo bar baz\n"
	var buf bytes.Buffer
	code := Run(strings.NewReader(input), &buf, []string{"-c"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	got := strings.TrimSpace(buf.String())
	if got != "24" {
		t.Errorf("expected 24 chars, got %q", got)
	}
}

func TestRunDefaultIsLines(t *testing.T) {
	input := "a\nb\nc\n"
	var buf bytes.Buffer
	code := Run(strings.NewReader(input), &buf, []string{})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	got := strings.TrimSpace(buf.String())
	if got != "3" {
		t.Errorf("expected 3 lines (default), got %q", got)
	}
}

func TestRunEmptyInput(t *testing.T) {
	var buf bytes.Buffer
	code := Run(strings.NewReader(""), &buf, []string{"-l"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	got := strings.TrimSpace(buf.String())
	if got != "0" {
		t.Errorf("expected 0, got %q", got)
	}
}

func TestRunBadFlag(t *testing.T) {
	var buf bytes.Buffer
	code := Run(strings.NewReader(""), &buf, []string{"--unknown"})
	if code != 1 {
		t.Errorf("expected exit code 1 for bad flag, got %d", code)
	}
}

func TestCompiles(t *testing.T) {
}
