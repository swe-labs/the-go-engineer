package main

import (
	"testing"
)

func TestValidateWorkflowValid(t *testing.T) {
	w := Workflow{
		Name: "Test",
		On:   "pull_request",
		Jobs: []Job{
			{
				Name:   "build",
				RunsOn: "ubuntu-latest",
				Steps:  []Step{{Name: "checkout", Uses: "actions/checkout@v4"}},
			},
		},
	}
	errs := ValidateWorkflow(w)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidateWorkflowMissingName(t *testing.T) {
	w := Workflow{On: "push", Jobs: []Job{{Name: "j", RunsOn: "ubuntu", Steps: []Step{{Name: "s", Run: "echo hi"}}}}}
	errs := ValidateWorkflow(w)
	if len(errs) != 1 || errs[0].Error() != "name: workflow name is required" {
		t.Errorf("unexpected errors: %v", errs)
	}
}

func TestValidateWorkflowMissingTrigger(t *testing.T) {
	w := Workflow{Name: "test", Jobs: []Job{{Name: "j", RunsOn: "ubuntu", Steps: []Step{{Name: "s", Run: "echo hi"}}}}}
	errs := ValidateWorkflow(w)
	if len(errs) != 1 || errs[0].Error() != "on: trigger event is required" {
		t.Errorf("unexpected errors: %v", errs)
	}
}

func TestValidateWorkflowMissingJobName(t *testing.T) {
	w := Workflow{Name: "test", On: "push", Jobs: []Job{{RunsOn: "ubuntu", Steps: []Step{{Name: "s", Run: "echo hi"}}}}}
	errs := ValidateWorkflow(w)
	matched := false
	for _, e := range errs {
		if e.Error() == "jobs[0].name: job name is required" {
			matched = true
		}
	}
	if !matched {
		t.Errorf("expected job name error, got %v", errs)
	}
}

func TestValidateWorkflowStepNeedsAction(t *testing.T) {
	w := Workflow{Name: "test", On: "push", Jobs: []Job{{Name: "j", RunsOn: "ubuntu", Steps: []Step{{Name: "empty"}}}}}
	errs := ValidateWorkflow(w)
	matched := false
	for _, e := range errs {
		if e.Error() == "jobs[0].steps[0]: step must have run or uses" {
			matched = true
		}
	}
	if !matched {
		t.Errorf("expected step action error, got %v", errs)
	}
}

func TestValidateWorkflowTimeout(t *testing.T) {
	w := Workflow{Name: "test", On: "push", Jobs: []Job{{
		Name: "j", RunsOn: "ubuntu",
		Steps: []Step{{Name: "s", Run: "echo hi", Timeout: 400}},
	}}}
	errs := ValidateWorkflow(w)
	matched := false
	for _, e := range errs {
		if e.Error() == "jobs[0].steps[0].timeout: timeout exceeds 360 minutes" {
			matched = true
		}
	}
	if !matched {
		t.Errorf("expected timeout error, got %v", errs)
	}
}

func TestGoCacheKey(t *testing.T) {
	key := GoCacheKey("1.22", "ubuntu", "abcdef1234567890abcdef")
	want := "go-cache-1.22-ubuntu-abcdef12"
	if key != want {
		t.Errorf("got %q, want %q", key, want)
	}
}

func TestGenerateCacheStep(t *testing.T) {
	s := GenerateCacheStep("mykey", "myrestore", "/path/to/cache")
	if s.Name != "Cache Go modules" {
		t.Errorf("unexpected name: %s", s.Name)
	}
	if s.Uses != "actions/cache@v4" {
		t.Errorf("unexpected uses: %s", s.Uses)
	}
	if s.With["key"] != "mykey" {
		t.Errorf("unexpected key: %s", s.With["key"])
	}
}

func TestExpandMatrixEmpty(t *testing.T) {
	combos := ExpandMatrix(nil)
	if len(combos) != 1 {
		t.Errorf("expected 1 combo, got %d", len(combos))
	}
}

func TestExpandMatrixSingle(t *testing.T) {
	m := map[string][]string{"go": {"1.22"}}
	combos := ExpandMatrix(m)
	if len(combos) != 1 {
		t.Fatalf("expected 1 combo, got %d", len(combos))
	}
	if combos[0]["go"] != "1.22" {
		t.Errorf("expected go=1.22, got %s", combos[0]["go"])
	}
}

func TestExpandMatrixCrossProduct(t *testing.T) {
	m := map[string][]string{
		"go": {"1.21", "1.22"},
		"os": {"linux", "windows"},
	}
	combos := ExpandMatrix(m)
	if len(combos) != 4 {
		t.Fatalf("expected 4 combos, got %d", len(combos))
	}
}

func TestSuggestRunner(t *testing.T) {
	tests := []struct {
		goVer, osSuffix, want string
	}{
		{"1.22", "", "ubuntu-latest"},
		{"1.22", "win-x64", "windows-latest"},
		{"1.22", "mac-x64", "macos-latest"},
		{"1.22", "linux-x64", "ubuntu-latest"},
	}
	for _, tc := range tests {
		got := SuggestRunner(tc.goVer, tc.osSuffix)
		if got != tc.want {
			t.Errorf("SuggestRunner(%q, %q) = %q, want %q", tc.goVer, tc.osSuffix, got, tc.want)
		}
	}
}
