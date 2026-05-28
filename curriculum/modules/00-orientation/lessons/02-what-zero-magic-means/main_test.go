package main

import (
	"strings"
	"testing"
)

func TestCardHasRequiredFields(t *testing.T) {
	c := card()

	if c.ID != "core-00-02" {
		t.Fatalf("ID = %q, want core-00-02", c.ID)
	}
	if strings.TrimSpace(c.Title) == "" {
		t.Fatal("Title must not be empty")
	}
	if strings.TrimSpace(c.Mission) == "" {
		t.Fatal("Mission must not be empty")
	}
	if strings.TrimSpace(c.Proof) == "" {
		t.Fatal("Proof must not be empty")
	}
	if strings.TrimSpace(c.NextStep) == "" {
		t.Fatal("NextStep must not be empty")
	}
}

func TestSummaryMentionsProofAndNextStep(t *testing.T) {
	summary := card().summary()

	for _, want := range []string{"Go Engineer Orientation", "core-00-02", "Proof:", "Next:"} {
		if !strings.Contains(summary, want) {
			t.Fatalf("summary does not contain %q\nsummary:\n%s", want, summary)
		}
	}
}
