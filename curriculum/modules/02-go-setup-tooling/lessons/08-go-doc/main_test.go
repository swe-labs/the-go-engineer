package main

import (
	"strings"
	"testing"
)

func TestCalculator_Add(t *testing.T) {
	calc := Calculator{Base: 100}
	got := calc.Add(25)
	want := 125
	if got != want {
		t.Errorf("calc.Add(25) = %d, want %d", got, want)
	}
}

func TestCleanDocString(t *testing.T) {
	raw := "// Package math provides mathematical constants.\n// It is part of the standard library."
	cleaned := CleanDocString(raw)

	if !strings.HasPrefix(cleaned, "Package math") {
		t.Errorf("cleaned doc string should start with symbol name: %q", cleaned)
	}
	if strings.Contains(cleaned, "//") {
		t.Errorf("cleaned doc string should not contain slashes: %q", cleaned)
	}
}

func TestMainExecution(t *testing.T) {
	main()
}

