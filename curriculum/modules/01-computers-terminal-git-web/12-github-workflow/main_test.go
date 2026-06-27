package main

import "testing"

func TestCreatePR(t *testing.T) {
	id := createPR("Fix bug", "Fixes the nil pointer issue", "bugfix", "main")
	if id != 1 {
		t.Errorf("expected id 1, got %d", id)
	}
	if prTitle != "Fix bug" {
		t.Errorf("expected title 'Fix bug', got %q", prTitle)
	}
	if prStatus != statusOpen {
		t.Errorf("expected status open, got %s", prStatusString(prStatus))
	}
}

func TestPRWorkflow(t *testing.T) {
	createPR("Add feature", "New feature", "feature", "main")
	approvePR()
	ok, _ := mergePR()
	if !ok {
		t.Fatal("expected successful merge on approved PR")
	}
	if prStatus != statusMerged {
		t.Errorf("expected merged, got %s", prStatusString(prStatus))
	}
}

func TestPRMergeRequiresApproval(t *testing.T) {
	createPR("Unreviewed change", "Quick fix", "hotfix", "main")
	ok, msg := mergePR()
	if ok {
		t.Fatal("expected merge to fail without approval")
	}
	if msg != "PR must be approved before merging" {
		t.Errorf("unexpected message: %q", msg)
	}
}

func TestPRStatusTransitions(t *testing.T) {
	createPR("Test", "Desc", "src", "tgt")
	if prStatus != statusOpen {
		t.Errorf("expected open, got %s", prStatusString(prStatus))
	}
	requestChanges()
	if prStatus != statusChangesRequested {
		t.Errorf("expected changes requested, got %s", prStatusString(prStatus))
	}
	approvePR()
	if prStatus != statusApproved {
		t.Errorf("expected approved, got %s", prStatusString(prStatus))
	}
	closePR()
	if prStatus != statusApproved {
		t.Errorf("close should not affect approved PRs; got %s", prStatusString(prStatus))
	}
}
