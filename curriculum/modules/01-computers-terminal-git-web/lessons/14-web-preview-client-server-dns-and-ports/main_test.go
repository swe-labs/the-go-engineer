package main

import (
	"strings"
	"testing"
)

func TestCardHasRequiredFields(t *testing.T) {
	c := card()

	if c.ID != "core-01-14" {
		t.Fatalf("ID = %q, want core-01-14", c.ID)
	}
	if strings.TrimSpace(c.Title) == "" {
		t.Fatal("Title must not be empty")
	}
	if strings.TrimSpace(c.MentalModel) == "" {
		t.Fatal("MentalModel must not be empty")
	}
	if strings.TrimSpace(c.MachineView) == "" {
		t.Fatal("MachineView must not be empty")
	}
	if strings.TrimSpace(c.CommonMistake) == "" {
		t.Fatal("CommonMistake must not be empty")
	}
	if strings.TrimSpace(c.Fix) == "" {
		t.Fatal("Fix must not be empty")
	}
	if len(c.Commands) == 0 {
		t.Fatal("Commands must not be empty")
	}
	if strings.TrimSpace(c.NextStep) == "" {
		t.Fatal("NextStep must not be empty")
	}
}

func TestSummaryContainsProofSignals(t *testing.T) {
	summary := card().summary()

	for _, want := range []string{"core-01-14", "Mental model:", "Machine view:", "Common mistake:", "Fix:", "Next:"} {
		if !strings.Contains(summary, want) {
			t.Fatalf("summary missing %q\nsummary:\n%s", want, summary)
		}
	}
}
