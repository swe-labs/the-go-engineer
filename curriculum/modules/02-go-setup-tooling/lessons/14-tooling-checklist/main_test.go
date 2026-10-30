package main

import (
	"testing"
)

func TestRunToolingChecklist(t *testing.T) {
	checks := RunToolingChecklist()
	if len(checks) < 4 {
		t.Fatalf("expected at least 4 checks in checklist, got %d", len(checks))
	}

	for _, c := range checks {
		if !c.Passed {
			t.Errorf("expected check %q to pass, details: %s", c.Name, c.Details)
		}
		if c.Details == "" {
			t.Errorf("check %q should include diagnostic details", c.Name)
		}
	}
}

func TestMainExecution(t *testing.T) {
	main()
}
