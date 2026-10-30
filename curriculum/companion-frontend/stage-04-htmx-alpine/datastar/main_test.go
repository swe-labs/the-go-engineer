package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
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

	// Verify Datastar integration markers
	expectedMarkers := []string{
		"/static/datastar.js",
		"data-signals=",
		"data-init=\"@get('/events')\"",
		"data-on-click=\"@post('/tasks?type=Data+Ingestion')\"",
		"id=\"metrics-panel\"",
		"id=\"tasks-table-body\"",
		"id=\"log-stream\"",
	}

	for _, marker := range expectedMarkers {
		if !strings.Contains(body, marker) {
			t.Errorf("expected dashboard HTML to contain marker %q", marker)
		}
	}

	// Verify HTMX and Alpine are NOT loaded in this variant
	forbiddenMarkers := []string{"htmx.min.js", "alpine.min.js"}
	for _, marker := range forbiddenMarkers {
		if strings.Contains(body, marker) {
			t.Errorf("unexpected script marker in Datastar variant: %s", marker)
		}
	}
}

func TestDashboardServer_EnqueueTaskDatastar(t *testing.T) {
	server, _ := setupTestDashboard(t)

	req := httptest.NewRequest(http.MethodPost, "/tasks?type=Security+Audit&duration=5", nil)
	rr := httptest.NewRecorder()
	server.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on enqueuing task, got %d: %s", rr.Code, rr.Body.String())
	}

	contentType := rr.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/event-stream") {
		t.Fatalf("expected Content-Type text/event-stream for Datastar response, got %s", contentType)
	}

	body := rr.Body.String()

	// 1. Verify datastar-patch-elements event prepending task row
	if !strings.Contains(body, "event: datastar-patch-elements") {
		t.Errorf("expected response to contain 'event: datastar-patch-elements'")
	}
	if !strings.Contains(body, "data: selector #tasks-table-body") || !strings.Contains(body, "data: mode prepend") {
		t.Errorf("expected prepend target for #tasks-table-body, got:\n%s", body)
	}
	if !strings.Contains(body, "Security Audit") {
		t.Errorf("expected task title 'Security Audit' in elements payload")
	}

	// 2. Verify datastar-patch-signals closing the modal
	if !strings.Contains(body, "event: datastar-patch-signals") {
		t.Errorf("expected response to contain 'event: datastar-patch-signals'")
	}
	if !strings.Contains(body, `data: signals {"modalOpen": false}`) {
		t.Errorf("expected signals patch to close modal")
	}
}

func TestDashboardServer_CancelTask(t *testing.T) {
	server, pool := setupTestDashboard(t)

	task := pool.Enqueue("Long Backup Job", 30*time.Second)

	cancelURL := fmt.Sprintf("/tasks/%d/cancel", task.ID)
	req := httptest.NewRequest(http.MethodPost, cancelURL, nil)
	rr := httptest.NewRecorder()
	server.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on cancellation, got %d: %s", rr.Code, rr.Body.String())
	}

	body := rr.Body.String()
	if !strings.Contains(body, "event: datastar-patch-elements") {
		t.Errorf("expected datastar-patch-elements event in cancel response")
	}
	if !strings.Contains(body, fmt.Sprintf("data: selector #task-%d", task.ID)) {
		t.Errorf("expected selector #task-%d in cancel response", task.ID)
	}
	if !strings.Contains(body, "CANCELLED") {
		t.Errorf("expected CANCELLED status in response body")
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
	if !strings.Contains(body, "event: datastar-patch-elements") || !strings.Contains(body, "data: selector #log-stream") {
		t.Errorf("expected Datastar SSE log event payload in body, got: %s", body)
	}
}

func setupTestDashboardTempl(t *testing.T) (*DashboardServer, *WorkerPool) {
	t.Helper()

	broker := NewSSEBroker()
	pool := NewWorkerPool(2, broker, nil) // pure templ mode
	server := NewDashboardServer(pool, broker, nil, defaultStaticFS)

	return server, pool
}

func TestDashboardServer_DatastarTemplMode(t *testing.T) {
	server, pool := setupTestDashboardTempl(t)

	// 1. Root dashboard render
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	server.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for templ Datastar dashboard, got %d", rr.Code)
	}
	body := rr.Body.String()
	for _, marker := range []string{"/static/datastar.js", "data-signals=", "data-init=\"@get('/events')\"", "id=\"metrics-panel\"", "id=\"tasks-table-body\"", "id=\"log-stream\""} {
		if !strings.Contains(body, marker) {
			t.Errorf("expected templ Datastar dashboard to contain %q", marker)
		}
	}

	// 2. Enqueue task via POST /tasks
	reqPost := httptest.NewRequest(http.MethodPost, "/tasks?type=Templ+Datastar+Job&duration=4", nil)
	rrPost := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrPost, reqPost)

	if rrPost.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on enqueuing task in templ Datastar mode, got %d", rrPost.Code)
	}
	postBody := rrPost.Body.String()
	if !strings.Contains(postBody, "event: datastar-patch-elements") || !strings.Contains(postBody, "Templ Datastar Job") {
		t.Errorf("expected datastar-patch-elements with task row, got: %s", postBody)
	}
	if !strings.Contains(postBody, `data: signals {"modalOpen": false}`) {
		t.Errorf("expected modal close signal patch, got: %s", postBody)
	}

	// 3. Cancel task
	task := pool.Enqueue("Long Datastar Templ Task", 20*time.Second)
	reqCancel := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/tasks/%d/cancel", task.ID), nil)
	rrCancel := httptest.NewRecorder()
	server.Handler().ServeHTTP(rrCancel, reqCancel)

	if rrCancel.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on cancel in templ Datastar mode, got %d", rrCancel.Code)
	}
	cancelBody := rrCancel.Body.String()
	if !strings.Contains(cancelBody, "event: datastar-patch-elements") || !strings.Contains(cancelBody, "badge-cancelled") {
		t.Errorf("expected cancelled row patch in response, got: %s", cancelBody)
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

