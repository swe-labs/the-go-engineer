package main

import (
	"testing"
)

func TestNewLearningCycleHasSixSteps(t *testing.T) {
	lc := NewLearningCycle("test")
	if len(lc.Steps) != 6 {
		t.Fatalf("got %d steps, want 6", len(lc.Steps))
	}
}

func TestLearningCycleStepNames(t *testing.T) {
	lc := NewLearningCycle("test")
	expected := []string{"Read", "Run", "Try", "Test", "Reflect", "Checkpoint"}
	for i, s := range lc.Steps {
		if s.Name != expected[i] {
			t.Errorf("step %d name = %q, want %q", i, s.Name, expected[i])
		}
	}
}

func TestLearningCycleAdvance(t *testing.T) {
	lc := NewLearningCycle("test")
	if lc.Complete {
		t.Fatal("new cycle should not be complete")
	}
	for i := 0; i < 6; i++ {
		if !lc.Advance() {
			t.Fatalf("Advance returned false at step %d", i)
		}
	}
	if !lc.Complete {
		t.Fatal("cycle should be complete after 6 advances")
	}
	if lc.Advance() {
		t.Fatal("Advance on empty cycle should return false")
	}
}

func TestLearningCycleHasActionsAndArtifacts(t *testing.T) {
	lc := NewLearningCycle("test")
	for _, s := range lc.Steps {
		if s.Action == "" {
			t.Errorf("step %d (%s) has empty Action", s.Index, s.Name)
		}
		if s.Artifact == "" {
			t.Errorf("step %d (%s) has empty Artifact", s.Index, s.Name)
		}
	}
}

func TestLearningCycleIndexesAreSequential(t *testing.T) {
	lc := NewLearningCycle("test")
	for i, s := range lc.Steps {
		if s.Index != i+1 {
			t.Errorf("step %d has Index %d, want %d", i, s.Index, i+1)
		}
	}
}
