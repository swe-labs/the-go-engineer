package components_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/rasel9t6/the-go-engineer/curriculum/companion-frontend/stage-04-htmx-alpine/components"
)

func renderToString(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		t.Fatalf("failed to render templ component: %v", err)
	}
	return buf.String()
}

func TestHTMX_MetricsCard(t *testing.T) {
	m := components.MetricsSnapshot{
		ActiveWorkers: 2,
		TotalWorkers:  4,
		QueueDepth:    8, // > 5 should trigger text-amber
		CompletedJobs: 42,
		Goroutines:    18,
		OOB:           true,
	}

	out := renderToString(t, components.MetricsCard(m))

	expected := []string{
		"id=\"metrics-panel\"",
		"hx-swap-oob=\"true\"",
		"2 / 4",
		"text-amber",
		"42",
		"18",
		"Worker Pool Capacity",
	}

	for _, s := range expected {
		if !strings.Contains(out, s) {
			t.Errorf("expected MetricsCard to contain %q, got:\n%s", s, out)
		}
	}
}

func TestHTMX_TaskRow(t *testing.T) {
	taskActive := &components.Task{
		ID:          101,
		Type:        "Security Audit",
		Status:      "RUNNING",
		StatusLower: "running",
		WorkerID:    2,
		Duration:    5 * time.Second,
		DurationStr: "5.0s",
	}

	outActive := renderToString(t, components.TaskRow(taskActive))

	if !strings.Contains(outActive, "id=\"task-101\"") {
		t.Errorf("expected task row ID marker, got: %s", outActive)
	}
	if !strings.Contains(outActive, "Security Audit") {
		t.Errorf("expected task type, got: %s", outActive)
	}
	if !strings.Contains(outActive, "Worker-2") {
		t.Errorf("expected Worker-2, got: %s", outActive)
	}
	if !strings.Contains(outActive, "hx-post=\"/tasks/101/cancel\"") {
		t.Errorf("expected cancel button for running task, got: %s", outActive)
	}

	taskDone := &components.Task{
		ID:          102,
		Type:        "Backup",
		Status:      "COMPLETED",
		StatusLower: "completed",
		WorkerID:    1,
		DurationStr: "2.1s",
	}

	outDone := renderToString(t, components.TaskRow(taskDone))
	if !strings.Contains(outDone, "Done") {
		t.Errorf("expected 'Done' label for completed task, got: %s", outDone)
	}
	if strings.Contains(outDone, "btn-action-cancel") {
		t.Errorf("unexpected cancel button for completed task")
	}
}

func TestHTMX_Dashboard(t *testing.T) {
	m := components.MetricsSnapshot{
		ActiveWorkers: 1,
		TotalWorkers:  2,
		QueueDepth:    0,
		CompletedJobs: 5,
		Goroutines:    10,
	}
	tasks := []*components.Task{
		{
			ID:          1,
			Type:        "Data Ingestion",
			Status:      "QUEUED",
			StatusLower: "queued",
			DurationStr: "3.0s",
		},
	}

	out := renderToString(t, components.Dashboard(m, tasks))

	markers := []string{
		"hx-ext=\"sse\"",
		"sse-connect=\"/events\"",
		"id=\"metrics-panel\"",
		"id=\"tasks-table-body\"",
		"id=\"task-1\"",
		"Data Ingestion",
		"id=\"log-stream\"",
		"/static/htmx.min.js",
		"/static/alpine.min.js",
	}

	for _, marker := range markers {
		if !strings.Contains(out, marker) {
			t.Errorf("expected dashboard HTML to contain %q", marker)
		}
	}
}
