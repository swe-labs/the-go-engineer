package main

import (
	"strings"
	"testing"
)

func TestCardHasRequiredFields(t *testing.T) {
	c := card()

	if c.ID != "core-02-03" {
		t.Fatalf("ID = %q, want core-02-03", c.ID)
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
	if strings.TrimSpace(c.CommandPurpose) == "" {
		t.Fatal("CommandPurpose must not be empty")
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

func TestSummaryContainsToolingSignals(t *testing.T) {
	summary := card().summary()

	for _, want := range []string{"core-02-03", "Mental model:", "Machine view:", "Purpose:", "Common mistake:", "Commands:", "Next:"} {
		if !strings.Contains(summary, want) {
			t.Fatalf("summary missing %q\nsummary:\n%s", want, summary)
		}
	}
}
