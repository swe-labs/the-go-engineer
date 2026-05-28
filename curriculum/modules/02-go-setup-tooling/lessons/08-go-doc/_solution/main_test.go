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

func TestSolutionNamesLesson(t *testing.T) {
	output := render(card())
	if !strings.Contains(output, "go doc") {
		t.Fatalf("solution output should include lesson title; got:\n%s", output)
	}
}
