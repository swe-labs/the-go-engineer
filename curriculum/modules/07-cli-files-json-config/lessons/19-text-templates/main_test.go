package main

import (
	"strings"
	"testing"
)

func TestRenderTasks(t *testing.T) {
	tasks := []Task{
		{Name: "Task A", Priority: 1, Done: true},
		{Name: "Task B", Priority: 2, Done: false},
		{Name: "Task C", Priority: 3, Done: false},
	}
	output, err := RenderTasks(tasks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "[X]") {
		t.Error("expected completed marker [X]")
	}
	if !strings.Contains(output, "[ ]") {
		t.Error("expected incomplete marker [ ]")
	}
	if !strings.Contains(output, "[HIGH]") {
		t.Error("expected [HIGH] label")
	}
	if !strings.Contains(output, "[MED]") {
		t.Error("expected [MED] label")
	}
	if !strings.Contains(output, "[LOW]") {
		t.Error("expected [LOW] label")
	}
	if !strings.Contains(output, "Task A") || !strings.Contains(output, "Task B") || !strings.Contains(output, "Task C") {
		t.Error("missing task names in output")
	}
}

func TestRenderTasksEmpty(t *testing.T) {
	output, err := RenderTasks(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output != "" {
		t.Errorf("expected empty output, got %q", output)
	}
}

func TestRenderTasksLabelFn(t *testing.T) {
	if label(1) != "[HIGH]" {
		t.Errorf("expected [HIGH], got %s", label(1))
	}
	if label(2) != "[MED]" {
		t.Errorf("expected [MED], got %s", label(2))
	}
	if label(3) != "[LOW]" {
		t.Errorf("expected [LOW], got %s", label(3))
	}
	if label(99) != "[UNKNOWN]" {
		t.Errorf("expected [UNKNOWN], got %s", label(99))
	}
}

func TestCompiles(t *testing.T) {
}
