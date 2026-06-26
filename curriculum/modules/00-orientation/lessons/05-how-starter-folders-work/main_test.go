package main

import (
	"testing"
)

func TestCompareTasksStarterIncomplete(t *testing.T) {
	starter := Task{Name: "", Desc: "", Complete: false}
	solution := Task{Name: "Greet the user", Desc: "Write a greeting function", Complete: true}
	c := CompareTasks(starter, solution)

	if !c.MissingName {
		t.Error("expected MissingName = true for empty starter name")
	}
	if !c.MissingDesc {
		t.Error("expected MissingDesc = true for empty starter desc")
	}
	if c.IsComplete {
		t.Error("expected IsComplete = false for incomplete starter")
	}
}

func TestCompareTasksStarterComplete(t *testing.T) {
	starter := Task{Name: "Greet the user", Desc: "Write a greeting function", Complete: true}
	solution := Task{Name: "Greet the user", Desc: "Write a greeting function", Complete: true}
	c := CompareTasks(starter, solution)

	if c.MissingName {
		t.Error("expected MissingName = false when names match")
	}
	if c.MissingDesc {
		t.Error("expected MissingDesc = false when descs match")
	}
	if !c.IsComplete {
		t.Error("expected IsComplete = true when starter is complete")
	}
}

func TestCompareTasksPartialComplete(t *testing.T) {
	tests := []struct {
		name         string
		starter      Task
		solution     Task
		wantMissName bool
		wantMissDesc bool
		wantComplete bool
	}{
		{
			name:         "all fields missing",
			starter:      Task{Name: "", Desc: "", Complete: false},
			solution:     Task{Name: "Task A", Desc: "Do something", Complete: true},
			wantMissName: true,
			wantMissDesc: true,
			wantComplete: false,
		},
		{
			name:         "name only missing",
			starter:      Task{Name: "", Desc: "Do something", Complete: true},
			solution:     Task{Name: "Task A", Desc: "Do something", Complete: true},
			wantMissName: true,
			wantMissDesc: false,
			wantComplete: true,
		},
		{
			name:         "desc only missing",
			starter:      Task{Name: "Task A", Desc: "", Complete: true},
			solution:     Task{Name: "Task A", Desc: "Do something", Complete: true},
			wantMissName: false,
			wantMissDesc: true,
			wantComplete: true,
		},
		{
			name:         "complete match",
			starter:      Task{Name: "Task A", Desc: "Do something", Complete: true},
			solution:     Task{Name: "Task A", Desc: "Do something", Complete: true},
			wantMissName: false,
			wantMissDesc: false,
			wantComplete: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := CompareTasks(tc.starter, tc.solution)
			if c.MissingName != tc.wantMissName {
				t.Errorf("MissingName = %v, want %v", c.MissingName, tc.wantMissName)
			}
			if c.MissingDesc != tc.wantMissDesc {
				t.Errorf("MissingDesc = %v, want %v", c.MissingDesc, tc.wantMissDesc)
			}
			if c.IsComplete != tc.wantComplete {
				t.Errorf("IsComplete = %v, want %v", c.IsComplete, tc.wantComplete)
			}
		})
	}
}

func TestComparisonReport(t *testing.T) {
	tests := []struct {
		name     string
		starter  Task
		solution Task
		wantPart string
	}{
		{
			name:     "incomplete reports missing",
			starter:  Task{Name: "", Desc: "", Complete: false},
			solution: Task{Name: "A", Desc: "B", Complete: true},
			wantPart: "incomplete",
		},
		{
			name:     "complete reports complete",
			starter:  Task{Name: "A", Desc: "B", Complete: true},
			solution: Task{Name: "A", Desc: "B", Complete: true},
			wantPart: "complete",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := CompareTasks(tc.starter, tc.solution)
			report := c.Report()
			if len(report) == 0 {
				t.Error("Report should not be empty")
			}
		})
	}
}
