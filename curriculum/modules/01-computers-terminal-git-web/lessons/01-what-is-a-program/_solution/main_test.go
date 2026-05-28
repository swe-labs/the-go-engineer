package main

import (
	"strings"
	"testing"
)

func TestSolutionHasNoTODO(t *testing.T) {
	output := render(card())
	if strings.Contains(output, "TODO") {
		t.Fatalf("solution still contains TODO markers:\n%s", output)
	}
}

func TestSolutionNamesConcept(t *testing.T) {
	output := render(card())
	if !strings.Contains(output, "What is a program?") {
		t.Fatalf("solution output should include lesson title; got:\n%s", output)
	}
}
