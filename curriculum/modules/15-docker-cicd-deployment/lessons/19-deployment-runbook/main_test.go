package main

import (
	"strings"
	"testing"
)

func TestNewRunbook(t *testing.T) {
	rb := NewRunbook("svc", "v1")
	if rb.ServiceName != "svc" {
		t.Errorf("expected svc, got %s", rb.ServiceName)
	}
	if rb.Version != "v1" {
		t.Errorf("expected v1, got %s", rb.Version)
	}
	if len(rb.Steps) != 0 {
		t.Errorf("expected 0 steps, got %d", len(rb.Steps))
	}
}

func TestAddStep(t *testing.T) {
	rb := NewRunbook("svc", "v1")
	rb.AddStep(RunbookStep{Phase: PhasePreDeploy, Name: "check", Required: true})
	if len(rb.Steps) != 1 {
		t.Fatalf("expected 1 step, got %d", len(rb.Steps))
	}
	if rb.Steps[0].Name != "check" {
		t.Errorf("unexpected step name: %s", rb.Steps[0].Name)
	}
}

func TestAddContact(t *testing.T) {
	rb := NewRunbook("svc", "v1")
	rb.AddContact("oncall@x.com")
	if len(rb.Contacts) != 1 {
		t.Errorf("expected 1 contact, got %d", len(rb.Contacts))
	}
}

func TestGenerateChecklist(t *testing.T) {
	rb := NewRunbook("myservice", "v2.0.0")
	rb.AddContact("ops@co.com")
	rb.AddStep(RunbookStep{Phase: PhasePreDeploy, Name: "Check CI"})
	rb.AddStep(RunbookStep{Phase: PhaseDeploy, Name: "Deploy binary"})
	rb.AddStep(RunbookStep{Phase: PhasePostDeploy, Name: "Verify health"})
	rb.AddStep(RunbookStep{Phase: PhaseRollback, Name: "Rollback"})
	checklist := rb.GenerateChecklist()
	if !strings.Contains(checklist, "myservice") {
		t.Error("checklist should contain service name")
	}
	if !strings.Contains(checklist, "PRE-DEPLOY") {
		t.Error("checklist should contain PRE-DEPLOY phase")
	}
	if !strings.Contains(checklist, "POST-DEPLOY") {
		t.Error("checklist should contain POST-DEPLOY phase")
	}
	if !strings.Contains(checklist, "Check CI") {
		t.Error("checklist should contain step name")
	}
	if !strings.Contains(checklist, "ops@co.com") {
		t.Error("checklist should contain contacts")
	}
}

func TestValidateValid(t *testing.T) {
	rb := NewRunbook("svc", "v1")
	rb.AddContact("oncall@x.com")
	rb.AddStep(RunbookStep{Phase: PhasePreDeploy, Name: "pre"})
	rb.AddStep(RunbookStep{Phase: PhaseDeploy, Name: "deploy"})
	rb.AddStep(RunbookStep{Phase: PhasePostDeploy, Name: "post"})
	errs := rb.Validate()
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidateMissingName(t *testing.T) {
	rb := NewRunbook("svc", "v1")
	rb.AddContact("oncall@x.com")
	rb.AddStep(RunbookStep{Phase: PhasePreDeploy, Name: ""})
	errs := rb.Validate()
	found := false
	for _, e := range errs {
		if strings.Contains(e.Error(), "step name is required") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected step name error, got %v", errs)
	}
}

func TestValidateMissingServiceName(t *testing.T) {
	rb := NewRunbook("", "v1")
	rb.AddContact("c@x.com")
	rb.AddStep(RunbookStep{Phase: PhasePreDeploy, Name: "pre"})
	rb.AddStep(RunbookStep{Phase: PhaseDeploy, Name: "deploy"})
	rb.AddStep(RunbookStep{Phase: PhasePostDeploy, Name: "post"})
	errs := rb.Validate()
	found := false
	for _, e := range errs {
		if strings.Contains(e.Error(), "service name is required") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected service name error, got %v", errs)
	}
}

func TestValidateMissingVersion(t *testing.T) {
	rb := NewRunbook("svc", "")
	rb.AddContact("c@x.com")
	rb.AddStep(RunbookStep{Phase: PhasePreDeploy, Name: "pre"})
	rb.AddStep(RunbookStep{Phase: PhaseDeploy, Name: "deploy"})
	rb.AddStep(RunbookStep{Phase: PhasePostDeploy, Name: "post"})
	errs := rb.Validate()
	found := false
	for _, e := range errs {
		if strings.Contains(e.Error(), "version is required") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected version error, got %v", errs)
	}
}

func TestValidateMissingContact(t *testing.T) {
	rb := NewRunbook("svc", "v1")
	rb.AddStep(RunbookStep{Phase: PhasePreDeploy, Name: "pre"})
	rb.AddStep(RunbookStep{Phase: PhaseDeploy, Name: "deploy"})
	rb.AddStep(RunbookStep{Phase: PhasePostDeploy, Name: "post"})
	errs := rb.Validate()
	found := false
	for _, e := range errs {
		if strings.Contains(e.Error(), "incident contact is required") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected contact error, got %v", errs)
	}
}

func TestValidateMissingPhases(t *testing.T) {
	rb := NewRunbook("svc", "v1")
	rb.AddContact("c@x.com")
	rb.AddStep(RunbookStep{Phase: PhasePreDeploy, Name: "pre"})
	errs := rb.Validate()
	foundDeploy := false
	foundPost := false
	for _, e := range errs {
		if strings.Contains(e.Error(), "deploy step is required") {
			foundDeploy = true
		}
		if strings.Contains(e.Error(), "post-deploy step is required") {
			foundPost = true
		}
	}
	if !foundDeploy {
		t.Errorf("expected deploy step error")
	}
	if !foundPost {
		t.Errorf("expected post-deploy step error")
	}
}

func TestExecutePhase(t *testing.T) {
	rb := NewRunbook("svc", "v1")
	rb.AddContact("c@x.com")
	rb.AddStep(RunbookStep{Phase: PhasePreDeploy, Name: "check1"})
	rb.AddStep(RunbookStep{Phase: PhasePreDeploy, Name: "check2"})
	ok, executed := rb.ExecutePhase(PhasePreDeploy)
	if !ok {
		t.Error("expected execute to succeed")
	}
	if len(executed) != 2 {
		t.Errorf("expected 2 executed steps, got %d", len(executed))
	}
}

func TestExecutePhaseEmpty(t *testing.T) {
	rb := NewRunbook("svc", "v1")
	rb.AddContact("c@x.com")
	ok, executed := rb.ExecutePhase(PhaseDeploy)
	if !ok {
		t.Error("expected execute to succeed")
	}
	if len(executed) != 0 {
		t.Errorf("expected 0 executed steps, got %d", len(executed))
	}
}

func TestGenerateChecklistIncludesPhases(t *testing.T) {
	rb := NewRunbook("svc", "v1")
	rb.AddContact("c@x.com")
	rb.AddStep(RunbookStep{Phase: PhasePreDeploy, Name: "pre"})
	rb.AddStep(RunbookStep{Phase: PhaseDeploy, Name: "deploy"})
	rb.AddStep(RunbookStep{Phase: PhasePostDeploy, Name: "post"})
	rb.AddStep(RunbookStep{Phase: PhaseRollback, Name: "rollback"})
	checklist := rb.GenerateChecklist()
	for _, section := range []string{"PRE-DEPLOY", "DEPLOY", "POST-DEPLOY", "ROLLBACK"} {
		if !strings.Contains(checklist, section) {
			t.Errorf("checklist missing section: %s", section)
		}
	}
}
