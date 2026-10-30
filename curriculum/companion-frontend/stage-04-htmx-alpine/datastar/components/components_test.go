package components_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/rasel9t6/the-go-engineer/curriculum/companion-frontend/stage-04-htmx-alpine/datastar/components"
)

func renderToString(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		t.Fatalf("failed to render templ component: %v", err)
	}
	return buf.String()
}

func TestDatastar_MetricsCard(t *testing.T) {
	m := components.MetricsSnapshot{
		ActiveWorkers: 3,
		TotalWorkers:  4,
		QueueDepth:    2,
		CompletedJobs: 15,
		Goroutines:    24,
	}

	out := renderToString(t, components.MetricsCard(m))

	expected := []string{
		"id=\"metrics-panel\"",
		"3 / 4",
		"text-primary",
		"15",
		"24",
	}

	for _, s := range expected {
		if !strings.Contains(out, s) {
			t.Errorf("expected Datastar MetricsCard to contain %q, got:\n%s", s, out)
		}
	}
}

func TestDatastar_TaskRow(t *testing.T) {
	task := &components.Task{
		ID:          205,
		Type:        "Index Reconciliation",
		Status:      "RUNNING",
		StatusLower: "running",
		WorkerID:    3,
		Duration:    4 * time.Second,
		DurationStr: "4.0s",
	}

	out := renderToString(t, components.TaskRow(task))

	if !strings.Contains(out, "id=\"task-205\"") {
		t.Errorf("expected task row ID, got: %s", out)
	}
	if !strings.Contains(out, "@post") || !strings.Contains(out, "/tasks/205/cancel") {
		t.Errorf("expected Datastar @post cancel action, got: %s", out)
	}
}

func TestDatastar_Dashboard(t *testing.T) {
	m := components.MetricsSnapshot{
		ActiveWorkers: 2,
		TotalWorkers:  4,
		QueueDepth:    1,
		CompletedJobs: 9,
		Goroutines:    12,
	}
	tasks := []*components.Task{
		{
			ID:          1,
			Type:        "Data Ingestion",
			Status:      "COMPLETED",
			StatusLower: "completed",
			DurationStr: "1.2s",
		},
	}

	out := renderToString(t, components.Dashboard(m, tasks))

	markers := []string{
		"/static/datastar.js",
		"data-signals=",
		"data-init=\"@get('/events')\"",
		"data-on-click=\"@post('/tasks?type=Data+Ingestion')\"",
		"id=\"metrics-panel\"",
		"id=\"tasks-table-body\"",
		"id=\"task-1\"",
		"id=\"log-stream\"",
	}

	for _, marker := range markers {
		if !strings.Contains(out, marker) {
			t.Errorf("expected Datastar dashboard to contain %q", marker)
		}
	}

	// Verify HTMX is not present
	if strings.Contains(out, "hx-") {
		t.Errorf("unexpected hx- attribute in Datastar template")
	}
}
