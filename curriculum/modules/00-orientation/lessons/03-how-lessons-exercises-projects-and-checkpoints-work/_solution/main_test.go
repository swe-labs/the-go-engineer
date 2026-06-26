package main

import "testing"

func TestStepsAreInLearningOrder(t *testing.T) {
	got := steps()
	want := []string{"read", "run", "try", "test", "reflect", "checkpoint"}

	if len(got) != len(want) {
		t.Fatalf("len(steps) = %d, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i].Name != want[i] {
			t.Fatalf("step %d = %q, want %q", i, got[i].Name, want[i])
		}
	}
}

func TestStepsHavePurposes(t *testing.T) {
	for _, step := range steps() {
		if step.Purpose == "" || step.Purpose == "TODO" {
			t.Fatalf("step %q does not have a completed purpose", step.Name)
		}
	}
}
