package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func setupTestDashboard(t *testing.T) (*DashboardServer, *WorkerPool) {
	t.Helper()

	tmpl, err := template.ParseFS(defaultTemplatesFS, "templates/*.html")
	if err != nil {
		t.Fatalf("failed to parse test templates: %v", err)
	}

	broker := NewSSEBroker()
	pool := NewWorkerPool(2, broker, tmpl)
	server := NewDashboardServer(pool, broker, tmpl, defaultStaticFS)

	return server, pool
}

func TestDashboardServer_RootRender(t *testing.T) {
	server, _ := setupTestDashboard(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	server.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 for root dashboard, got %d", rr.Code)
	}

	body := rr.Body.String()

	// Verify HTMX SSE integration markers
	expectedMarkers := []string{
		"hx-ext=\"sse\"",
		"sse-connect=\"/events\"",
		"id=\"metrics-panel\"",
		"id=\"tasks-table-body\"",
		"id=\"log-stream\"",
		"/static/htmx.min.js",
		"/static/alpine.min.js",
	}

	for _, marker := range expectedMarkers {
		if !strings.Contains(body, marker) {
			t.Errorf("expected dashboard HTML to contain marker %q", marker)
		}
	}
}

func TestDashboardServer_EnqueueTaskAndOOB(t *testing.T) {
	server, _ := setupTestDashboard(t)

	form := url.Values{
		"type":     {"Security Audit"},
		"duration": {"5"},
	}

	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	server.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on enqueuing task, got %d: %s", rr.Code, rr.Body.String())
	}

	body := rr.Body.String()

	// 1. Verify primary task row partial
	if !strings.Contains(body, "<tr id=\"task-") || !strings.Contains(body, "Security Audit") {
		t.Errorf("expected rendered partial table row for new task")
	}

	// 2. Verify HTMX Out-of-Band (OOB) swap for metrics card
	if !strings.Contains(body, "hx-swap-oob=\"true\"") {
		t.Errorf("expected response to contain out-of-band metrics swap attribute")
	}
}

func TestDashboardServer_CancelTask(t *testing.T) {
	server, pool := setupTestDashboard(t)

	// Enqueue long-running task
	task := pool.Enqueue("Slow Backup Job", 30*time.Second)

	cancelURL := fmt.Sprintf("/tasks/%d/cancel", task.ID)
	req := httptest.NewRequest(http.MethodPost, cancelURL, nil)
	rr := httptest.NewRecorder()
	server.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on cancellation, got %d: %s", rr.Code, rr.Body.String())
	}

	body := rr.Body.String()
	if !strings.Contains(body, "badge-cancelled") && !strings.Contains(body, "CANCELLED") {
		t.Errorf("expected cancelled badge in response, got: %s", body)
	}
}

func TestDashboardServer_MetricsEndpoint(t *testing.T) {
	server, _ := setupTestDashboard(t)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()
	server.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for /metrics, got %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "id=\"metrics-panel\"") || !strings.Contains(body, "Active Workers") {
		t.Errorf("expected metrics fragment in response, got: %s", body)
	}
}

func TestDashboardServer_SSEStream(t *testing.T) {
	server, _ := setupTestDashboard(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := httptest.NewRequest(http.MethodGet, "/events", nil).WithContext(ctx)
	rr := httptest.NewRecorder()

	// Run handler in goroutine and cancel after receiving initial stream
	done := make(chan struct{})
	go func() {
		server.Handler().ServeHTTP(rr, req)
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()
	<-done

	contentType := rr.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/event-stream") {
		t.Errorf("expected Content-Type 'text/event-stream', got %q", contentType)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "event: log") || !strings.Contains(body, "data:") {
		t.Errorf("expected SSE event payload in body, got: %s", body)
	}
}

func setupTestDashboardTempl(t *testing.T) (*DashboardServer, *WorkerPool) {
	t.Helper()

	broker := NewSSEBroker()
	pool := NewWorkerPool(2, broker, nil) // pure templ mode
	server := NewDashboardServer(pool, broker, nil, defaultStaticFS)

	return server, pool
}

func TestDashboardServer_TemplMode(t *testing.T) {
	server, pool := setupTestDashboardTempl(t)

	// 1. Root dashboard render
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	server.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for templ root dashboard, got %d", rr.Code)
	}
	body := rr.Body.String()
	for _, marker := range []string{"hx-ext=\"sse\"", "sse-connect=\"/events\"", "id=\"metrics-panel\"", "id=\"tasks-table-body\"", "id=\"log-stream\""} {
		if !strings.Contains(body, marker) {
			t.Errorf("expected templ dashboard to contain %q", marker)
		}
	}

	// 2. Enqueue task via POST /tasks
	form := url.Values{"type": {"Templ Benchmark"}, "duration": {"4"}}
	reqPost := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(form.Encode()))
	reqPost.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rrPost := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrPost, reqPost)

	if rrPost.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on enqueuing task in templ mode, got %d", rrPost.Code)
	}
	postBody := rrPost.Body.String()
	if !strings.Contains(postBody, "Templ Benchmark") || !strings.Contains(postBody, "hx-swap-oob=\"true\"") {
		t.Errorf("expected templ response to include task row and OOB swap, got %s", postBody)
	}

	// 3. Cancel task
	task := pool.Enqueue("Long Running Templ Task", 20*time.Second)
	reqCancel := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/tasks/%d/cancel", task.ID), nil)
	rrCancel := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrCancel, reqCancel)

	if rrCancel.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on cancel in templ mode, got %d", rrCancel.Code)
	}
	if !strings.Contains(rrCancel.Body.String(), "badge-cancelled") {
		t.Errorf("expected badge-cancelled in templ cancel response, got %s", rrCancel.Body.String())
	}

	// 4. GET /metrics
	reqMetrics := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rrMetrics := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrMetrics, reqMetrics)
	if rrMetrics.Code != http.StatusOK || !strings.Contains(rrMetrics.Body.String(), "id=\"metrics-panel\"") {
		t.Errorf("expected valid metrics panel in templ mode, got %s", rrMetrics.Body.String())
	}

	// 5. GET /tasks list
	reqTasks := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rrTasks := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrTasks, reqTasks)
	if rrTasks.Code != http.StatusOK || !strings.Contains(rrTasks.Body.String(), "<tr id=\"task-") {
		t.Errorf("expected task rows list in templ mode, got %s", rrTasks.Body.String())
	}
}

